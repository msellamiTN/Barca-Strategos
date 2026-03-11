package vendor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"barca-strategos/pkg/types"
)

type Service struct {
	vendors      map[string]*types.Vendor
	assessments  map[string][]*types.VendorAssessment
	categories   map[string]*types.VendorCategory
	broadcast    func([]byte)
}

func New() *Service {
	return &Service{
		vendors:     make(map[string]*types.Vendor),
		assessments: make(map[string][]*types.VendorAssessment),
		categories:  make(map[string]*types.VendorCategory),
	}
}

func (s *Service) SetBroadcast(broadcast func([]byte)) {
	s.broadcast = broadcast
}

func (s *Service) Initialize(ctx context.Context) error {
	// Initialize vendor categories
	s.initializeCategories()
	
	// Start background monitoring
	go s.startBackgroundMonitoring()
	
	return nil
}

func (s *Service) AddVendor(ctx context.Context, request *types.VendorRequest) (*types.Vendor, error) {
	// Validate vendor request
	if err := s.validateVendorRequest(request); err != nil {
		return nil, err
	}

	// Get category
	category, exists := s.categories[request.CategoryID]
	if !exists {
		return nil, types.VendorError{
			Code:    "INVALID_CATEGORY",
			Message: "Invalid vendor category",
		}
	}

	// Create vendor
	vendor := &types.Vendor{
		ID:               uuid.New().String(),
		Name:             request.Name,
		Description:      request.Description,
		Category:         *category,
		ContactEmail:     request.ContactEmail,
		ContactPhone:     request.ContactPhone,
		Address:          request.Address,
		ServicesOffered:  request.ServicesOffered,
		ContractStartDate: request.ContractStartDate,
		ContractEndDate:  request.ContractEndDate,
		Status:           types.VendorStatusActive,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	s.vendors[vendor.ID] = vendor
	
	// Schedule initial assessment
	go s.scheduleInitialAssessment(vendor.ID)
	
	s.broadcastVendorUpdate()
	
	return vendor, nil
}

func (s *Service) GetVendor(ctx context.Context, vendorID string) (*types.Vendor, error) {
	vendor, exists := s.vendors[vendorID]
	if !exists {
		return nil, types.VendorError{
			Code:    "VENDOR_NOT_FOUND",
			Message: fmt.Sprintf("Vendor %s not found", vendorID),
		}
	}
	return vendor, nil
}

func (s *Service) ListVendors(ctx context.Context) ([]*types.Vendor, error) {
	vendors := make([]*types.Vendor, 0, len(s.vendors))
	for _, vendor := range s.vendors {
		vendors = append(vendors, vendor)
	}
	return vendors, nil
}

func (s *Service) UpdateVendor(ctx context.Context, vendorID string, request *types.VendorUpdateRequest) error {
	vendor, exists := s.vendors[vendorID]
	if !exists {
		return types.VendorError{
			Code:    "VENDOR_NOT_FOUND",
			Message: fmt.Sprintf("Vendor %s not found", vendorID),
		}
	}

	// Validate update request
	if err := s.validateUpdateRequest(vendor, request); err != nil {
		return err
	}

	// Update vendor
	if request.NameChange != nil {
		vendor.Name = *request.NameChange
	}
	if request.DescriptionChange != nil {
		vendor.Description = *request.DescriptionChange
	}
	if request.ContactEmailChange != nil {
		vendor.ContactEmail = *request.ContactEmailChange
	}
	if request.ContactPhoneChange != nil {
		vendor.ContactPhone = request.ContactPhoneChange
	}
	if request.AddressChange != nil {
		vendor.Address = *request.AddressChange
	}
	if request.ServicesChange != nil {
		vendor.ServicesOffered = *request.ServicesChange
	}
	if request.ContractEndDateChange != nil {
		vendor.ContractEndDate = request.ContractEndDateChange
	}
	
	vendor.UpdatedAt = time.Now().UTC()
	s.vendors[vendorID] = vendor
	s.broadcastVendorUpdate()

	// Trigger reassessment if significant changes
	if request.RequiresReassessment {
		go s.scheduleReassessment(vendorID)
	}

	return nil
}

func (s *Service) ConductVendorAssessment(ctx context.Context, vendorID string) (*types.VendorAssessment, error) {
	vendor, exists := s.vendors[vendorID]
	if !exists {
		return nil, types.VendorError{
			Code:    "VENDOR_NOT_FOUND",
			Message: fmt.Sprintf("Vendor %s not found", vendorID),
		}
	}

	// Perform risk assessment
	riskAssessment := s.performRiskAssessment(vendor)
	
	// Perform compliance assessment
	complianceAssessment := s.performComplianceAssessment(vendor)
	
	// Calculate overall vendor risk score
	overallRiskScore := s.calculateOverallRiskScore(&riskAssessment, &complianceAssessment)
	
	// Determine vendor risk level
	riskLevel := s.determineVendorRiskLevel(overallRiskScore)
	
	// Generate recommendations
	recommendations := s.generateVendorRecommendations(&riskAssessment, &complianceAssessment)
	
	assessment := &types.VendorAssessment{
		AssessmentID:         uuid.New().String(),
		VendorID:             vendorID,
		VendorName:           vendor.Name,
		Timestamp:            time.Now().UTC(),
		RiskAssessment:       riskAssessment,
		ComplianceAssessment: complianceAssessment,
		OverallRiskScore:     overallRiskScore,
		RiskLevel:            riskLevel,
		Recommendations:      recommendations,
		NextAssessmentDate:   time.Now().UTC().AddDate(0, 0, s.getAssessmentFrequency(&riskLevel)),
	}

	// Store assessment
	if _, exists := s.assessments[vendorID]; !exists {
		s.assessments[vendorID] = make([]*types.VendorAssessment, 0)
	}
	s.assessments[vendorID] = append(s.assessments[vendorID], assessment)
	
	// Update vendor risk level
	vendor.RiskLevel = &riskLevel
	vendor.LastAssessed = &assessment.Timestamp
	vendor.UpdatedAt = time.Now().UTC()
	
	s.broadcastVendorUpdate()
	
	return assessment, nil
}

func (s *Service) GetVendorRiskStatus(ctx context.Context, vendorID string) (*types.VendorRiskStatus, error) {
	vendor, exists := s.vendors[vendorID]
	if !exists {
		return nil, types.VendorError{
			Code:    "VENDOR_NOT_FOUND",
			Message: fmt.Sprintf("Vendor %s not found", vendorID),
		}
	}

	latestAssessment, err := s.getLatestAssessment(vendorID)
	if err != nil {
		return nil, err
	}

	criticalIssues := s.getCriticalIssues(latestAssessment)
	monitoringStatus := s.getMonitoringStatus(vendorID)

	return &types.VendorRiskStatus{
		VendorID:         vendorID,
		VendorName:       vendor.Name,
		RiskLevel:        latestAssessment.RiskLevel,
		RiskScore:        latestAssessment.OverallRiskScore,
		LastAssessed:     latestAssessment.Timestamp,
		NextAssessment:   latestAssessment.NextAssessmentDate,
		CriticalIssues:   criticalIssues,
		MonitoringStatus: monitoringStatus,
	}, nil
}

func (s *Service) MonitorVendorCompliance(ctx context.Context, vendorID string) (*types.VendorComplianceStatus, error) {
	vendor, exists := s.vendors[vendorID]
	if !exists {
		return nil, types.VendorError{
			Code:    "VENDOR_NOT_FOUND",
			Message: fmt.Sprintf("Vendor %s not found", vendorID),
		}
	}

	latestAssessment, err := s.getLatestAssessment(vendorID)
	if err != nil {
		// Return default compliance status if no assessment exists
		return &types.VendorComplianceStatus{
			VendorID:        vendorID,
			ComplianceScore: 0.8, // Default score
			LastChecked:     time.Now().UTC(),
			ComplianceIssues: []types.ComplianceIssue{},
			NextReviewDate:   time.Now().UTC().AddDate(0, 0, 90),
		}, nil
	}

	return &types.VendorComplianceStatus{
		VendorID:         vendorID,
		ComplianceScore:  latestAssessment.ComplianceAssessment.ComplianceScore,
		LastChecked:      latestAssessment.Timestamp,
		ComplianceIssues: latestAssessment.ComplianceAssessment.ComplianceIssues,
		NextReviewDate:   latestAssessment.NextAssessmentDate,
	}, nil
}

func (s *Service) GenerateVendorReport(ctx context.Context, scope *types.VendorScope) (*types.VendorRiskReport, error) {
	vendors := s.getVendorsByScope(scope)
	
	var vendorAssessments []types.VendorAssessmentPair
	for _, vendor := range vendors {
		if assessment, err := s.getLatestAssessment(vendor.ID); err == nil {
			vendorAssessments = append(vendorAssessments, types.VendorAssessmentPair{
				Vendor:     *vendor,
				Assessment: *assessment,
			})
		}
	}

	overallRiskProfile := s.calculateOverallRiskProfile(vendorAssessments)
	highRiskVendors := s.identifyHighRiskVendors(vendorAssessments)

	report := &types.VendorRiskReport{
		ReportID:           uuid.New().String(),
		GeneratedAt:        time.Now().UTC(),
		Scope:              *scope,
		TotalVendors:       len(vendorAssessments),
		VendorAssessments:  vendorAssessments,
		OverallRiskProfile: overallRiskProfile,
		HighRiskVendors:    highRiskVendors,
		Recommendations:    s.generatePortfolioRecommendations(vendorAssessments),
	}

	return report, nil
}

func (s *Service) GetVendorStats(ctx context.Context) (*types.VendorStats, error) {
	stats := &types.VendorStats{
		TotalVendors: len(s.vendors),
	}

	var totalRiskScore float64
	var riskScoreCount int

	for _, vendor := range s.vendors {
		switch vendor.Status {
		case types.VendorStatusActive:
			stats.ActiveVendors++
		}

		if vendor.RiskLevel != nil {
			switch *vendor.RiskLevel {
			case types.VendorRiskLevelCritical:
				stats.CriticalRiskVendors++
			case types.VendorRiskLevelHigh:
				stats.HighRiskVendors++
			case types.VendorRiskLevelMedium:
				stats.MediumRiskVendors++
			case types.VendorRiskLevelLow:
				stats.LowRiskVendors++
			}

			// Get latest assessment for risk score
			if assessment, err := s.getLatestAssessment(vendor.ID); err == nil {
				totalRiskScore += assessment.OverallRiskScore
				riskScoreCount++
			}
		}
	}

	if riskScoreCount > 0 {
		stats.AverageRiskScore = totalRiskScore / float64(riskScoreCount)
	}

	// Count vendors due for assessment
	now := time.Now().UTC()
	for _, vendor := range s.vendors {
		if assessment, err := s.getLatestAssessment(vendor.ID); err == nil {
			if assessment.NextAssessmentDate.Before(now) {
				stats.VendorsDueForAssessment++
			}
		}
	}

	stats.OverdueAssessments = stats.VendorsDueForAssessment

	return stats, nil
}

// Private methods

func (s *Service) initializeCategories() {
	categories := []types.VendorCategory{
		{
			ID:                     "CLOUD",
			Name:                   "Cloud Service Providers",
			Description:            "Vendors providing cloud infrastructure and services",
			RiskFactors:            []string{"Data sovereignty", "Service availability", "Data access controls", "Encryption standards"},
			BaseRiskScore:          0.6,
			AssessmentFrequencyDays: 180,
		},
		{
			ID:                     "SOFTWARE",
			Name:                   "Software Vendors",
			Description:            "Vendors providing software products and solutions",
			RiskFactors:            []string{"Software security", "Update management", "License compliance", "Support availability"},
			BaseRiskScore:          0.4,
			AssessmentFrequencyDays: 365,
		},
		{
			ID:                     "CONSULTING",
			Name:                   "Consulting Services",
			Description:            "Vendors providing consulting and professional services",
			RiskFactors:            []string{"Access to sensitive data", "Consultant credentials", "Confidentiality agreements", "Quality of deliverables"},
			BaseRiskScore:          0.5,
			AssessmentFrequencyDays: 365,
		},
		{
			ID:                     "INFRASTRUCTURE",
			Name:                   "Infrastructure Providers",
			Description:            "Vendors providing physical infrastructure and facilities",
			RiskFactors:            []string{"Physical security", "Disaster recovery", "Environmental controls", "Access management"},
			BaseRiskScore:          0.7,
			AssessmentFrequencyDays: 180,
		},
		{
			ID:                     "FINANCIAL",
			Name:                   "Financial Services",
			Description:            "Vendors providing financial and payment services",
			RiskFactors:            []string{"PCI DSS compliance", "Financial stability", "Regulatory compliance", "Transaction security"},
			BaseRiskScore:          0.8,
			AssessmentFrequencyDays: 90,
		},
	}

	for _, category := range categories {
		s.categories[category.ID] = &category
	}
}

func (s *Service) validateVendorRequest(request *types.VendorRequest) error {
	if request.Name == "" {
		return types.VendorError{
			Code:    "VALIDATION_ERROR",
			Message: "Vendor name is required",
		}
	}

	if request.ContactEmail == "" {
		return types.VendorError{
			Code:    "VALIDATION_ERROR",
			Message: "Contact email is required",
		}
	}

	if len(request.ServicesOffered) == 0 {
		return types.VendorError{
			Code:    "VALIDATION_ERROR",
			Message: "Services offered is required",
		}
	}

	return nil
}

func (s *Service) validateUpdateRequest(vendor *types.Vendor, request *types.VendorUpdateRequest) error {
	if vendor.Status == types.VendorStatusTerminated {
		return types.VendorError{
			Code:    "INVALID_STATUS",
			Message: "Cannot update terminated vendor",
		}
	}

	return nil
}

func (s *Service) performRiskAssessment(vendor *types.Vendor) types.VendorRiskAssessment {
	// Create risk factors based on category
	var riskFactors []types.RiskFactor
	for _, factorName := range vendor.Category.RiskFactors {
		riskFactors = append(riskFactors, types.RiskFactor{
			ID:                uuid.New().String(),
			Name:              factorName,
			Description:       fmt.Sprintf("Risk factor: %s", factorName),
			Severity:          types.RiskFactorSeverityMedium,
			Score:             vendor.Category.BaseRiskScore,
			MitigationRequired: vendor.Category.BaseRiskScore > 0.6,
		})
	}

	return types.VendorRiskAssessment{
		RiskScore:             vendor.Category.BaseRiskScore,
		RiskFactors:           riskFactors,
		FinancialStability:    0.8,
		OperationalCapability: 0.85,
		SecurityPosture:       0.9,
		ReputationScore:       0.8,
	}
}

func (s *Service) performComplianceAssessment(vendor *types.Vendor) types.VendorComplianceAssessment {
	// Mock compliance assessment
	return types.VendorComplianceAssessment{
		ComplianceScore:     0.85,
		ComplianceIssues:    []types.ComplianceIssue{},
		Certifications:      []types.Certification{},
		RegulatoryAdherence: 0.9,
		PolicyCompliance:    0.8,
	}
}

func (s *Service) calculateOverallRiskScore(riskAssessment *types.VendorRiskAssessment, complianceAssessment *types.VendorComplianceAssessment) float64 {
	riskWeight := 0.6
	complianceWeight := 0.4
	
	return (riskAssessment.RiskScore * riskWeight) + (complianceAssessment.ComplianceScore * complianceWeight)
}

func (s *Service) determineVendorRiskLevel(riskScore float64) types.VendorRiskLevel {
	if riskScore >= 0.8 {
		return types.VendorRiskLevelCritical
	} else if riskScore >= 0.6 {
		return types.VendorRiskLevelHigh
	} else if riskScore >= 0.4 {
		return types.VendorRiskLevelMedium
	} else if riskScore >= 0.2 {
		return types.VendorRiskLevelLow
	} else {
		return types.VendorRiskLevelMinimal
	}
}

func (s *Service) getAssessmentFrequency(riskLevel *types.VendorRiskLevel) int {
	switch *riskLevel {
	case types.VendorRiskLevelCritical:
		return 30
	case types.VendorRiskLevelHigh:
		return 60
	case types.VendorRiskLevelMedium:
		return 90
	case types.VendorRiskLevelLow:
		return 180
	case types.VendorRiskLevelMinimal:
		return 365
	default:
		return 180
	}
}

func (s *Service) generateVendorRecommendations(riskAssessment *types.VendorRiskAssessment, complianceAssessment *types.VendorComplianceAssessment) []types.VendorRecommendation {
	var recommendations []types.VendorRecommendation

	// Risk-based recommendations
	if riskAssessment.RiskScore > 0.7 {
		recommendations = append(recommendations, types.VendorRecommendation{
			Priority:    types.RecommendationPriorityHigh,
			Title:       "Address High Vendor Risk",
			Description: "Vendor poses significant risk requiring immediate attention",
			ActionItems: []string{
				"Implement additional controls",
				"Increase monitoring frequency",
				"Consider alternative vendors",
			},
			Owner:    "Vendor Manager",
			Timeline: "30 days",
		})
	}

	// Compliance-based recommendations
	if complianceAssessment.ComplianceScore < 0.8 {
		recommendations = append(recommendations, types.VendorRecommendation{
			Priority:    types.RecommendationPriorityMedium,
			Title:       "Improve Vendor Compliance",
			Description: "Vendor compliance gaps need to be addressed",
			ActionItems: []string{
				"Request compliance documentation",
				"Schedule compliance review",
				"Update contractual requirements",
			},
			Owner:    "Compliance Officer",
			Timeline: "60 days",
		})
	}

	return recommendations
}

func (s *Service) getLatestAssessment(vendorID string) (*types.VendorAssessment, error) {
	assessments, exists := s.assessments[vendorID]
	if !exists || len(assessments) == 0 {
		return nil, types.VendorError{
			Code:    "ASSESSMENT_NOT_FOUND",
			Message: fmt.Sprintf("No assessment found for vendor %s", vendorID),
		}
	}
	
	// Return the latest assessment
	latest := assessments[len(assessments)-1]
	return latest, nil
}

func (s *Service) getCriticalIssues(assessment *types.VendorAssessment) []string {
	var issues []string

	// Check for critical risk factors
	for _, factor := range assessment.RiskAssessment.RiskFactors {
		if factor.Severity == types.RiskFactorSeverityCritical {
			issues = append(issues, factor.Description)
		}
	}

	// Check for critical compliance issues
	for _, issue := range assessment.ComplianceAssessment.ComplianceIssues {
		if issue.Severity == types.ComplianceSeverityCritical {
			issues = append(issues, issue.Description)
		}
	}

	return issues
}

func (s *Service) getMonitoringStatus(vendorID string) types.MonitoringStatus {
	return types.MonitoringStatus{
		IsActive:      true,
		LastCheck:     time.Now().UTC(),
		AlertsCount:   0,
		IssuesDetected: 0,
	}
}

func (s *Service) getVendorsByScope(scope *types.VendorScope) []*types.Vendor {
	var vendors []*types.Vendor
	for _, vendor := range s.vendors {
		if (len(scope.Categories) == 0 || s.containsString(scope.Categories, vendor.Category.ID)) &&
			(len(scope.RiskLevels) == 0 || s.containsRiskLevel(scope.RiskLevels, vendor.RiskLevel)) &&
			(len(scope.Status) == 0 || s.containsVendorStatus(scope.Status, vendor.Status)) {
			vendors = append(vendors, vendor)
		}
	}
	return vendors
}

func (s *Service) calculateOverallRiskProfile(vendorAssessments []types.VendorAssessmentPair) types.VendorRiskProfile {
	if len(vendorAssessments) == 0 {
		return types.VendorRiskProfile{
			OverallRiskScore:   0.0,
			CriticalVendors:    0,
			HighRiskVendors:    0,
			MediumRiskVendors:  0,
			LowRiskVendors:     0,
			MinimalRiskVendors: 0,
			RiskDistribution:   make(map[string]int),
		}
	}

	var criticalCount, highCount, mediumCount, lowCount, minimalCount int
	var totalScore float64

	for _, pair := range vendorAssessments {
		switch pair.Assessment.RiskLevel {
		case types.VendorRiskLevelCritical:
			criticalCount++
		case types.VendorRiskLevelHigh:
			highCount++
		case types.VendorRiskLevelMedium:
			mediumCount++
		case types.VendorRiskLevelLow:
			lowCount++
		case types.VendorRiskLevelMinimal:
			minimalCount++
		}
		totalScore += pair.Assessment.OverallRiskScore
	}

	riskDistribution := map[string]int{
		"Critical": criticalCount,
		"High":     highCount,
		"Medium":   mediumCount,
		"Low":      lowCount,
		"Minimal":  minimalCount,
	}

	return types.VendorRiskProfile{
		OverallRiskScore:   totalScore / float64(len(vendorAssessments)),
		CriticalVendors:    criticalCount,
		HighRiskVendors:    highCount,
		MediumRiskVendors:  mediumCount,
		LowRiskVendors:     lowCount,
		MinimalRiskVendors: minimalCount,
		RiskDistribution:   riskDistribution,
	}
}

func (s *Service) identifyHighRiskVendors(vendorAssessments []types.VendorAssessmentPair) []types.VendorAssessmentPair {
	var highRisk []types.VendorAssessmentPair
	for _, pair := range vendorAssessments {
		if pair.Assessment.RiskLevel == types.VendorRiskLevelCritical || pair.Assessment.RiskLevel == types.VendorRiskLevelHigh {
			highRisk = append(highRisk, pair)
		}
	}
	return highRisk
}

func (s *Service) generatePortfolioRecommendations(vendorAssessments []types.VendorAssessmentPair) []types.VendorRecommendation {
	var recommendations []types.VendorRecommendation

	highRiskCount := 0
	for _, pair := range vendorAssessments {
		if pair.Assessment.RiskLevel == types.VendorRiskLevelCritical || pair.Assessment.RiskLevel == types.VendorRiskLevelHigh {
			highRiskCount++
		}
	}

	if highRiskCount > 0 {
		recommendations = append(recommendations, types.VendorRecommendation{
			Priority:    types.RecommendationPriorityHigh,
			Title:       "Address High-Risk Vendors",
			Description: fmt.Sprintf("%d vendors require immediate risk mitigation", highRiskCount),
			ActionItems: []string{
				"Review high-risk vendor contracts",
				"Implement additional monitoring",
				"Develop contingency plans",
				"Consider vendor diversification",
			},
			Owner:    "CPO",
			Timeline: "60 days",
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

func (s *Service) containsRiskLevel(slice []types.VendorRiskLevel, item *types.VendorRiskLevel) bool {
	if item == nil {
		return false
	}
	for _, s := range slice {
		if s == *item {
			return true
		}
	}
	return false
}

func (s *Service) containsVendorStatus(slice []types.VendorStatus, item types.VendorStatus) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *Service) scheduleInitialAssessment(vendorID string) {
	time.Sleep(7 * 24 * time.Hour) // Wait 7 days
	s.ConductVendorAssessment(context.Background(), vendorID)
}

func (s *Service) scheduleReassessment(vendorID string) {
	time.Sleep(14 * 24 * time.Hour) // Wait 14 days
	s.ConductVendorAssessment(context.Background(), vendorID)
}

func (s *Service) startBackgroundMonitoring() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		// Monitor vendor status
		s.monitorVendorStatus()
		
		// Monitor vendor compliance
		s.monitorVendorCompliance()
		
		// Schedule pending assessments
		s.schedulePendingAssessments()
	}
}

func (s *Service) monitorVendorStatus() {
	// Monitor vendor health and status changes
	for _, vendor := range s.vendors {
		// Check vendor health
		_ = s.getMonitoringStatus(vendor.ID)
	}
}

func (s *Service) monitorVendorCompliance() {
	// Monitor vendor compliance status
	for _, vendor := range s.vendors {
		_, _ = s.MonitorVendorCompliance(context.Background(), vendor.ID)
	}
}

func (s *Service) schedulePendingAssessments() {
	// Schedule assessments that are due
	now := time.Now().UTC()
	for _, vendor := range s.vendors {
		if assessment, err := s.getLatestAssessment(vendor.ID); err == nil {
			if assessment.NextAssessmentDate.Before(now) {
				go s.ConductVendorAssessment(context.Background(), vendor.ID)
			}
		}
	}
}

func (s *Service) broadcastVendorUpdate() {
	if s.broadcast != nil {
		vendors := make([]*types.Vendor, 0, len(s.vendors))
		for _, vendor := range s.vendors {
			vendors = append(vendors, vendor)
		}
		
		data := map[string]interface{}{
			"type":    "vendor_update",
			"vendors": vendors,
		}
		
		if jsonData, err := json.Marshal(data); err == nil {
			s.broadcast(jsonData)
		}
	}
}
