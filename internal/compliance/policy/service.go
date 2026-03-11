package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"barca-strategos/pkg/types"
)

type Service struct {
	policies           map[string]*types.Policy
	workflows          map[string]*types.ApprovalWorkflow
	complianceRecords  map[string]*types.PolicyComplianceStatus
	categories         map[string]*types.PolicyCategory
	broadcast          func([]byte)
}

func New() *Service {
	return &Service{
		policies:          make(map[string]*types.Policy),
		workflows:         make(map[string]*types.ApprovalWorkflow),
		complianceRecords: make(map[string]*types.PolicyComplianceStatus),
		categories:        make(map[string]*types.PolicyCategory),
	}
}

func (s *Service) SetBroadcast(broadcast func([]byte)) {
	s.broadcast = broadcast
}

func (s *Service) Initialize(ctx context.Context) error {
	// Initialize policy categories
	s.initializeCategories()
	
	// Start background monitoring
	go s.startBackgroundMonitoring()
	
	return nil
}

func (s *Service) CreatePolicy(ctx context.Context, request *types.PolicyRequest) (*types.Policy, error) {
	// Validate policy request
	if err := s.validatePolicyRequest(request); err != nil {
		return nil, err
	}

	// Get category
	category, exists := s.categories[request.CategoryID]
	if !exists {
		return nil, types.PolicyError{
			Code:    "INVALID_CATEGORY",
			Message: "Invalid policy category",
		}
	}

	// Create policy
	policy := &types.Policy{
		ID:          uuid.New().String(),
		Title:       request.Title,
		Description: request.Description,
		Content:     request.Content,
		Category:    *category,
		Version:     "1.0",
		Status:      types.PolicyStatusDraft,
		Owner:       request.Owner,
		Approvers:   request.Approvers,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Tags:        request.Tags,
	}

	s.policies[policy.ID] = policy
	s.broadcastPolicyUpdate()
	
	return policy, nil
}

func (s *Service) GetPolicy(ctx context.Context, policyID string) (*types.Policy, error) {
	policy, exists := s.policies[policyID]
	if !exists {
		return nil, types.PolicyError{
			Code:    "POLICY_NOT_FOUND",
			Message: fmt.Sprintf("Policy %s not found", policyID),
		}
	}
	return policy, nil
}

func (s *Service) ListPolicies(ctx context.Context) ([]*types.Policy, error) {
	policies := make([]*types.Policy, 0, len(s.policies))
	for _, policy := range s.policies {
		policies = append(policies, policy)
	}
	return policies, nil
}

func (s *Service) UpdatePolicy(ctx context.Context, policyID string, request *types.PolicyUpdateRequest) error {
	policy, exists := s.policies[policyID]
	if !exists {
		return types.PolicyError{
			Code:    "POLICY_NOT_FOUND",
			Message: fmt.Sprintf("Policy %s not found", policyID),
		}
	}

	// Validate update request
	if err := s.validateUpdateRequest(policy, request); err != nil {
		return err
	}

	// Create new version
	policy.Version = s.incrementVersion(policy.Version)
	policy.UpdatedAt = time.Now().UTC()

	if request.TitleChange != nil {
		policy.Title = *request.TitleChange
	}

	policy.Content = request.ContentChanges

	s.policies[policyID] = policy
	s.broadcastPolicyUpdate()

	// If policy is published, create review workflow
	if policy.Status == types.PolicyStatusPublished {
		s.SubmitForApproval(ctx, policyID, []string{"policy_committee"})
	}

	return nil
}

func (s *Service) SubmitForApproval(ctx context.Context, policyID string, approvers []string) error {
	policy, exists := s.policies[policyID]
	if !exists {
		return types.PolicyError{
			Code:    "POLICY_NOT_FOUND",
			Message: fmt.Sprintf("Policy %s not found", policyID),
		}
	}

	// Validate policy status
	if policy.Status != types.PolicyStatusDraft {
		return types.PolicyError{
			Code:    "INVALID_STATUS",
			Message: "Policy must be in draft status",
		}
	}

	// Create approval workflow
	workflow := &types.ApprovalWorkflow{
		ID:        uuid.New().String(),
		PolicyID:  policyID,
		Approvers: s.createApprovers(approvers),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour), // 7 days
		Status:    string(types.WorkflowStatusPending),
	}

	s.workflows[workflow.ID] = workflow
	
	// Update policy status
	policy.Status = types.PolicyStatusPendingApproval
	policy.UpdatedAt = time.Now().UTC()
	
	s.broadcastPolicyUpdate()
	
	return nil
}

func (s *Service) ApprovePolicy(ctx context.Context, policyID, approverID string, comments *string) error {
	workflow, exists := s.getWorkflowByPolicyID(policyID)
	if !exists {
		return types.PolicyError{
			Code:    "WORKFLOW_NOT_FOUND",
			Message: "Approval workflow not found",
		}
	}

	// Process approval
	approved := s.processApproval(workflow, approverID, comments)

	if approved {
		// Update policy status to approved
		policy := s.policies[policyID]
		policy.Status = types.PolicyStatusApproved
		policy.UpdatedAt = time.Now().UTC()
		
		// Schedule policy publication
		go s.schedulePolicyPublication(policyID)
	}

	s.broadcastPolicyUpdate()
	return nil
}

func (s *Service) PublishPolicy(ctx context.Context, policyID string) error {
	policy, exists := s.policies[policyID]
	if !exists {
		return types.PolicyError{
			Code:    "POLICY_NOT_FOUND",
			Message: fmt.Sprintf("Policy %s not found", policyID),
		}
	}

	// Validate policy status
	if policy.Status != types.PolicyStatusApproved {
		return types.PolicyError{
			Code:    "INVALID_STATUS",
			Message: "Policy must be approved before publishing",
		}
	}

	// Publish policy
	now := time.Now().UTC()
	policy.Status = types.PolicyStatusPublished
	policy.PublishedAt = &now
	policy.UpdatedAt = now
	
	// Set review date based on category
	reviewDate := now.AddDate(0, 0, int(policy.Category.ReviewPeriodDays))
	policy.ReviewDate = &reviewDate

	// Start compliance tracking
	s.startComplianceTracking(policyID)
	
	s.broadcastPolicyUpdate()
	return nil
}

func (s *Service) ArchivePolicy(ctx context.Context, policyID string) error {
	policy, exists := s.policies[policyID]
	if !exists {
		return types.PolicyError{
			Code:    "POLICY_NOT_FOUND",
			Message: fmt.Sprintf("Policy %s not found", policyID),
		}
	}

	// Validate policy status
	if policy.Status != types.PolicyStatusPublished {
		return types.PolicyError{
			Code:    "INVALID_STATUS",
			Message: "Only published policies can be archived",
		}
	}

	// Archive policy
	policy.Status = types.PolicyStatusArchived
	policy.UpdatedAt = time.Now().UTC()
	
	s.broadcastPolicyUpdate()
	return nil
}

func (s *Service) GetPolicyCompliance(ctx context.Context, policyID string) (*types.PolicyComplianceStatus, error) {
	if compliance, exists := s.complianceRecords[policyID]; exists {
		return compliance, nil
	}

	// Return default compliance status
	policy, exists := s.policies[policyID]
	if !exists {
		return nil, types.PolicyError{
			Code:    "POLICY_NOT_FOUND",
			Message: fmt.Sprintf("Policy %s not found", policyID),
		}
	}

	compliance := &types.PolicyComplianceStatus{
		PolicyID:        policyID,
		PolicyTitle:     policy.Title,
		ComplianceScore: 0.85, // Default score
		LastChecked:     time.Now().UTC(),
		ComplianceIssues: []types.ComplianceIssue{},
		NextReviewDate:  time.Now().UTC().Add(90 * 24 * time.Hour),
	}

	s.complianceRecords[policyID] = compliance
	return compliance, nil
}

func (s *Service) GenerateComplianceReport(ctx context.Context, scope *types.PolicyScope) (*types.PolicyComplianceReport, error) {
	policies := s.getPoliciesByScope(scope)
	
	var policyCompliances []types.PolicyComplianceStatus
	for _, policy := range policies {
		compliance, _ := s.GetPolicyCompliance(ctx, policy.ID)
		policyCompliances = append(policyCompliances, *compliance)
	}

	overallCompliance := s.calculateOverallCompliance(policyCompliances)
	
	report := &types.PolicyComplianceReport{
		ReportID:           uuid.New().String(),
		GeneratedAt:        time.Now().UTC(),
		Scope:              *scope,
		PolicyCompliances:  policyCompliances,
		OverallCompliance:  overallCompliance,
		Recommendations:    s.generateComplianceRecommendations(policyCompliances),
	}

	return report, nil
}

func (s *Service) GetPolicyStats(ctx context.Context) (*types.PolicyStats, error) {
	stats := &types.PolicyStats{
		TotalPolicies: len(s.policies),
	}

	for _, policy := range s.policies {
		switch policy.Status {
		case types.PolicyStatusDraft:
			stats.DraftPolicies++
		case types.PolicyStatusPendingApproval:
			stats.PendingApproval++
		case types.PolicyStatusPublished:
			stats.PublishedPolicies++
		case types.PolicyStatusArchived:
			stats.ArchivedPolicies++
		}
	}

	// Calculate average compliance score
	if len(s.complianceRecords) > 0 {
		totalScore := 0.0
		for _, compliance := range s.complianceRecords {
			totalScore += compliance.ComplianceScore
		}
		stats.AverageComplianceScore = totalScore / float64(len(s.complianceRecords))
	}

	// Count policies due for review
	now := time.Now().UTC()
	for _, policy := range s.policies {
		if policy.ReviewDate != nil && policy.ReviewDate.Before(now) {
			stats.PoliciesDueForReview++
		}
	}

	stats.OverdueReviews = stats.PoliciesDueForReview

	return stats, nil
}

// Private methods

func (s *Service) initializeCategories() {
	categories := []types.PolicyCategory{
		{
			ID:               "SECURITY",
			Name:             "Security Policies",
			Description:      "Policies related to information security",
			Subcategories:    []string{"Access Control", "Data Protection", "Network Security", "Incident Response"},
			ApprovalRequired: true,
			ReviewPeriodDays: 365,
		},
		{
			ID:               "COMPLIANCE",
			Name:             "Compliance Policies",
			Description:      "Policies related to regulatory compliance",
			Subcategories:    []string{"Regulatory Compliance", "Audit Requirements", "Reporting Standards", "Documentation"},
			ApprovalRequired: true,
			ReviewPeriodDays: 180,
		},
		{
			ID:               "OPERATIONAL",
			Name:             "Operational Policies",
			Description:      "Policies related to business operations",
			Subcategories:    []string{"Business Continuity", "Disaster Recovery", "Change Management", "Service Management"},
			ApprovalRequired: false,
			ReviewPeriodDays: 730,
		},
		{
			ID:               "HR",
			Name:             "HR Policies",
			Description:      "Policies related to human resources",
			Subcategories:    []string{"Employee Conduct", "Training and Awareness", "Privacy and Data Handling", "Remote Work"},
			ApprovalRequired: true,
			ReviewPeriodDays: 365,
		},
	}

	for _, category := range categories {
		s.categories[category.ID] = &category
	}
}

func (s *Service) validatePolicyRequest(request *types.PolicyRequest) error {
	if request.Title == "" {
		return types.PolicyError{
			Code:    "VALIDATION_ERROR",
			Message: "Policy title is required",
		}
	}

	if request.Content == "" {
		return types.PolicyError{
			Code:    "VALIDATION_ERROR",
			Message: "Policy content is required",
		}
	}

	if request.Owner == "" {
		return types.PolicyError{
			Code:    "VALIDATION_ERROR",
			Message: "Policy owner is required",
		}
	}

	return nil
}

func (s *Service) validateUpdateRequest(policy *types.Policy, request *types.PolicyUpdateRequest) error {
	if policy.Status == types.PolicyStatusArchived {
		return types.PolicyError{
			Code:    "INVALID_STATUS",
			Message: "Cannot update archived policy",
		}
	}

	if request.ContentChanges == "" && request.TitleChange == nil {
		return types.PolicyError{
			Code:    "VALIDATION_ERROR",
			Message: "No changes specified",
		}
	}

	return nil
}

func (s *Service) createApprovers(approvers []string) []types.Approver {
	var approverObjs []types.Approver
	for _, approverID := range approvers {
		approverObjs = append(approverObjs, types.Approver{
			ID:    approverID,
			Name:  approverID,
			Email: fmt.Sprintf("%s@company.com", approverID),
			Role:  "Approver",
		})
	}
	return approverObjs
}

func (s *Service) getWorkflowByPolicyID(policyID string) (*types.ApprovalWorkflow, bool) {
	for _, workflow := range s.workflows {
		if workflow.PolicyID == policyID {
			return workflow, true
		}
	}
	return nil, false
}

func (s *Service) processApproval(workflow *types.ApprovalWorkflow, approverID string, comments *string) bool {
	// Simple approval logic - in real implementation, this would be more complex
	for i, approver := range workflow.Approvers {
		if approver.ID == approverID {
			decision := types.ApprovalDecisionApproved
			workflow.Approvers[i].Decision = &decision
			workflow.Approvers[i].DecisionDate = &[]time.Time{time.Now().UTC()}[0]
			workflow.Approvers[i].Comments = comments
			break
		}
	}

	// Check if all approvers have approved
	allApproved := true
	for _, approver := range workflow.Approvers {
		if approver.Decision == nil || *approver.Decision != types.ApprovalDecisionApproved {
			allApproved = false
			break
		}
	}

	if allApproved {
		workflow.Status = string(types.WorkflowStatusApproved)
		return true
	}

	return false
}

func (s *Service) schedulePolicyPublication(policyID string) {
	time.Sleep(24 * time.Hour) // Wait 24 hours
	policy := s.policies[policyID]
	if policy.Status == types.PolicyStatusApproved {
		s.PublishPolicy(context.Background(), policyID)
	}
}

func (s *Service) startComplianceTracking(policyID string) {
	// Start compliance tracking for published policy
	compliance := &types.PolicyComplianceStatus{
		PolicyID:        policyID,
		PolicyTitle:     s.policies[policyID].Title,
		ComplianceScore: 1.0,
		LastChecked:     time.Now().UTC(),
		ComplianceIssues: []types.ComplianceIssue{},
		NextReviewDate:   time.Now().UTC().Add(90 * 24 * time.Hour),
	}
	s.complianceRecords[policyID] = compliance
}

func (s *Service) incrementVersion(version string) string {
	// Simple version increment
	if version == "1.0" {
		return "2.0"
	}
	return "2.0"
}

func (s *Service) getPoliciesByScope(scope *types.PolicyScope) []*types.Policy {
	var policies []*types.Policy
	for _, policy := range s.policies {
		if (len(scope.Categories) == 0 || s.containsString(scope.Categories, policy.Category.ID)) &&
			(len(scope.Status) == 0 || s.containsStatus(scope.Status, policy.Status)) {
			policies = append(policies, policy)
		}
	}
	return policies
}

func (s *Service) calculateOverallCompliance(compliances []types.PolicyComplianceStatus) float64 {
	if len(compliances) == 0 {
		return 1.0
	}

	totalScore := 0.0
	for _, compliance := range compliances {
		totalScore += compliance.ComplianceScore
	}

	return totalScore / float64(len(compliances))
}

func (s *Service) generateComplianceRecommendations(compliances []types.PolicyComplianceStatus) []types.PolicyRecommendation {
	var recommendations []types.PolicyRecommendation

	nonCompliantCount := 0
	for _, compliance := range compliances {
		if compliance.ComplianceScore < 0.8 {
			nonCompliantCount++
		}
	}

	if nonCompliantCount > 0 {
		recommendations = append(recommendations, types.PolicyRecommendation{
			Priority: types.RecommendationPriorityHigh,
			Title:    "Address Policy Compliance Issues",
			Description: fmt.Sprintf("%d policies require attention for compliance", nonCompliantCount),
			ActionItems: []string{
				"Review non-compliant policies",
				"Update policy documentation",
				"Conduct employee training",
				"Implement monitoring controls",
			},
			Owner:    "Policy Manager",
			Timeline: "30 days",
		})
	}

	return recommendations
}

func (s *Service) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *Service) containsStatus(slice []types.PolicyStatus, item types.PolicyStatus) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *Service) startBackgroundMonitoring() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		// Monitor policy workflows and status changes
		s.monitorPolicyStatus()
		
		// Check policy compliance
		s.checkPolicyCompliance()
		
		// Send review reminders
		s.sendReviewReminders()
	}
}

func (s *Service) monitorPolicyStatus() {
	// Monitor policy workflows and status changes
	for _, workflow := range s.workflows {
		if time.Now().UTC().After(workflow.ExpiresAt) {
			workflow.Status = string(types.WorkflowStatusExpired)
		}
	}
}

func (s *Service) checkPolicyCompliance() {
	// Check compliance for published policies
	for _, policy := range s.policies {
		if policy.Status == types.PolicyStatusPublished {
			// Update compliance status
			if compliance, exists := s.complianceRecords[policy.ID]; exists {
				compliance.LastChecked = time.Now().UTC()
			}
		}
	}
}

func (s *Service) sendReviewReminders() {
	// Send review reminders for policies due for review
	now := time.Now().UTC()
	for _, policy := range s.policies {
		if policy.ReviewDate != nil && policy.ReviewDate.Before(now) {
			// Send review reminder to policy owner
			fmt.Printf("Sending review reminder for policy: %s to owner: %s\n", policy.ID, policy.Owner)
		}
	}
}

func (s *Service) broadcastPolicyUpdate() {
	if s.broadcast != nil {
		policies := make([]*types.Policy, 0, len(s.policies))
		for _, policy := range s.policies {
			policies = append(policies, policy)
		}
		
		data := map[string]interface{}{
			"type":     "policy_update",
			"policies": policies,
		}
		
		if jsonData, err := json.Marshal(data); err == nil {
			s.broadcast(jsonData)
		}
	}
}
