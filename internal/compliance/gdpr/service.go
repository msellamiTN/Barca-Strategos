package gdpr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service manages GDPR compliance
type Service struct {
	gdpr       map[string]*types.GDPR
	principles map[string]*types.Principle
	processes  map[string]*types.Process
	rights     map[string]*types.Right
	breaches   map[string]*types.Breach
	dpias      map[string]*types.DPIA
	requests   map[string]*types.Request
}

// New creates a GDPR service
func New() *Service {
	return &Service{
		gdpr:       make(map[string]*types.GDPR),
		principles: make(map[string]*types.Principle),
		processes:  make(map[string]*types.Process),
		rights:     make(map[string]*types.Right),
		breaches:   make(map[string]*types.Breach),
		dpias:      make(map[string]*types.DPIA),
		requests:   make(map[string]*types.Request),
	}
}

// CreateGDPR creates a new GDPR implementation
func (s *Service) CreateGDPR(ctx context.Context, name, description, version string) (*types.GDPR, error) {
	gdpr := &types.GDPR{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Version:     version,
		Status:      types.GDPRStatusDraft,
		Principles:  []types.Principle{},
		Processes:   []types.Process{},
		Rights:      []types.Right{},
		Breaches:    []types.Breach{},
		Assessments: []types.DPIA{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.gdpr[gdpr.ID] = gdpr
	log.Printf("gdpr: created GDPR %s", gdpr.ID)
	
	// Initialize standard principles and rights
	s.initializeStandardPrinciples(gdpr.ID)
	s.initializeStandardRights(gdpr.ID)
	
	return gdpr, nil
}

// GetGDPR retrieves a GDPR by ID
func (s *Service) GetGDPR(ctx context.Context, id string) (*types.GDPR, error) {
	gdpr, ok := s.gdpr[id]
	if !ok {
		return nil, ErrGDPRNotFound
	}
	return gdpr, nil
}

// ListGDPRs returns all GDPR instances
func (s *Service) ListGDPRs(ctx context.Context) ([]*types.GDPR, error) {
	var list []*types.GDPR
	for _, gdpr := range s.gdpr {
		list = append(list, gdpr)
	}
	return list, nil
}

// CreateProcess creates a new data processing activity
func (s *Service) CreateProcess(ctx context.Context, gdprID, name, description, purpose string, legalBasis types.LegalBasis, owner string) (*types.Process, error) {
	process := &types.Process{
		ID:            uuid.New().String(),
		Name:          name,
		Description:   description,
		Purpose:       purpose,
		LegalBasis:    legalBasis,
		DataCategories: []types.DataCategory{},
		DataSubjects:   []types.DataSubject{},
		Processors:    []types.Processor{},
		Controllers:   []types.Controller{},
		Transfers:     []types.Transfer{},
		Retention:     types.Retention{},
		Security:      types.Security{},
		Status:        types.ProcessStatusActive,
		Owner:         owner,
		NextReview:    time.Now().UTC().AddDate(0, 12, 0),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	
	s.processes[process.ID] = process
	
	// Add to GDPR
	if gdpr, ok := s.gdpr[gdprID]; ok {
		gdpr.Processes = append(gdpr.Processes, *process)
		gdpr.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("gdpr: created process %s for GDPR %s", process.ID, gdprID)
	return process, nil
}

// UpdateProcessStatus updates process status
func (s *Service) UpdateProcessStatus(ctx context.Context, processID string, status types.ProcessStatus) error {
	process, ok := s.processes[processID]
	if !ok {
		return ErrProcessNotFound
	}
	
	process.Status = status
	process.UpdatedAt = time.Now().UTC()
	
	log.Printf("gdpr: updated process %s status to %s", processID, status)
	return nil
}

// CreateRequest creates a new data subject request
func (s *Service) CreateRequest(ctx context.Context, rightID string, requestType types.RequestType, dataSubject, identity, description string, priority types.Priority, assignedTo string) (*types.Request, error) {
	request := &types.Request{
		ID:          uuid.New().String(),
		RightID:     rightID,
		Type:        requestType,
		DataSubject: dataSubject,
		Identity:    identity,
		Description: description,
		Status:      types.RequestStatusReceived,
		Priority:    priority,
		AssignedTo:  assignedTo,
		DueDate:     time.Now().UTC().AddDate(0, 0, 30), // 30 days
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.requests[request.ID] = request
	
	// Add to right
	if right, ok := s.rights[rightID]; ok {
		right.Requests = append(right.Requests, *request)
		right.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("gdpr: created request %s for right %s", request.ID, rightID)
	return request, nil
}

// UpdateRequestStatus updates request status
func (s *Service) UpdateRequestStatus(ctx context.Context, requestID string, status types.RequestStatus, resolution string) error {
	request, ok := s.requests[requestID]
	if !ok {
		return ErrRequestNotFound
	}
	
	request.Status = status
	request.Resolution = resolution
	request.UpdatedAt = time.Now().UTC()
	
	if status == types.RequestStatusCompleted {
		now := time.Now().UTC()
		request.CompletedAt = &now
	}
	
	log.Printf("gdpr: updated request %s status to %s", requestID, status)
	return nil
}

// CreateBreach creates a new data breach
func (s *Service) CreateBreach(ctx context.Context, gdprID, title, description string, breachDate time.Time, breachType types.BreachType, categories []types.DataCategory, affected int, cause string, owner string) (*types.Breach, error) {
	breach := &types.Breach{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Date:        breachDate,
		Detected:    time.Now().UTC(),
		Type:        breachType,
		Categories:  categories,
		Affected:    affected,
		Cause:       cause,
		Impact: types.Impact{
			Level:        types.ImpactLevelMedium,
			Likelihood:   types.LikelihoodPossible,
			Severity:     types.SeverityMedium,
			Assessment:   "Initial assessment - to be updated",
			Consequences: []string{},
		},
		Notification: types.Notification{
			Required:  false,
			Deadline:  breachDate.AddDate(0, 0, 72), // 72 hours
			Authority: "",
			DataSubjects: false,
			Content:   "",
			Sent:      false,
		},
		Remediation: "",
		Status:      types.BreachStatusOpen,
		Owner:       owner,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.breaches[breach.ID] = breach
	
	// Add to GDPR
	if gdpr, ok := s.gdpr[gdprID]; ok {
		gdpr.Breaches = append(gdpr.Breaches, *breach)
		gdpr.UpdatedAt = time.Now().UTC()
	}
	
	// Check if notification is required
	s.checkNotificationRequirement(breach)
	
	log.Printf("gdpr: created breach %s for GDPR %s", breach.ID, gdprID)
	return breach, nil
}

// UpdateBreachStatus updates breach status
func (s *Service) UpdateBreachStatus(ctx context.Context, breachID string, status types.BreachStatus, remediation string) error {
	breach, ok := s.breaches[breachID]
	if !ok {
		return ErrBreachNotFound
	}
	
	breach.Status = status
	breach.Remediation = remediation
	breach.UpdatedAt = time.Now().UTC()
	
	log.Printf("gdpr: updated breach %s status to %s", breachID, status)
	return nil
}

// CreateDPIA creates a new Data Protection Impact Assessment
func (s *Service) CreateDPIA(ctx context.Context, gdprID, title, description, processID, controller, assessor string) (*types.DPIA, error) {
	dpia := &types.DPIA{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		ProcessID:   processID,
		Controller:  controller,
		Assessor:    assessor,
		Status:      types.DPIAStatusRequired,
		Risk: types.Risk{
			Level:       types.RiskLevelMedium,
			Likelihood:  types.LikelihoodPossible,
			Impact:      types.Impact{},
			Description: "Initial risk assessment",
			Mitigation:  "",
		},
		Measures:    []types.Measure{},
		Outcome:     "",
		Approval:    "",
		ReviewDate:  time.Now().UTC().AddDate(0, 12, 0),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.dpias[dpia.ID] = dpia
	
	// Add to GDPR
	if gdpr, ok := s.gdpr[gdprID]; ok {
		gdpr.Assessments = append(gdpr.Assessments, *dpia)
		gdpr.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("gdpr: created DPIA %s for GDPR %s", dpia.ID, gdprID)
	return dpia, nil
}

// UpdateDPIAStatus updates DPIA status
func (s *Service) UpdateDPIAStatus(ctx context.Context, dpiaID string, status types.DPIAStatus, outcome string, approval string) error {
	dpia, ok := s.dpias[dpiaID]
	if !ok {
		return ErrDPIANotFound
	}
	
	dpia.Status = status
	dpia.Outcome = outcome
	dpia.Approval = approval
	dpia.UpdatedAt = time.Now().UTC()
	
	log.Printf("gdpr: updated DPIA %s status to %s", dpiaID, status)
	return nil
}

// GetComplianceScore calculates overall GDPR compliance score
func (s *Service) GetComplianceScore(ctx context.Context, gdprID string) (float64, error) {
	gdpr, ok := s.gdpr[gdprID]
	if !ok {
		return 0, ErrGDPRNotFound
	}
	
	if len(gdpr.Principles) == 0 {
		return 0, nil
	}
	
	totalScore := 0.0
	for _, principle := range gdpr.Principles {
		switch principle.Status {
		case types.PrincipleStatusCompliant:
			totalScore += 1.0
		case types.PrincipleStatusValidated:
			totalScore += 0.8
		case types.PrincipleStatusImplemented:
			totalScore += 0.6
		case types.PrincipleStatusPartiallyImplemented:
			totalScore += 0.3
		}
	}
	
	return (totalScore / float64(len(gdpr.Principles))) * 100, nil
}

// BroadcastUpdates sends GDPR updates via WebSocket
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type": "gdpr_update",
		"gdpr": s.gdpr,
	})
	broadcast(data)
}

// initializeStandardPrinciples loads GDPR standard principles
func (s *Service) initializeStandardPrinciples(gdprID string) {
	standardPrinciples := []struct {
		title       string
		description string
		article     string
		category    types.PrincipleCategory
	}{
		{"Lawfulness, fairness and transparency", "Process personal data lawfully, fairly and in a transparent manner", "Article 5(1)(a)", types.PrincipleCategoryLawfulness},
		{"Purpose limitation", "Collect personal data for specified, explicit and legitimate purposes", "Article 5(1)(b)", types.PrincipleCategoryPurpose},
		{"Data minimisation", "Collect and process only the personal data that is necessary for your purposes", "Article 5(1)(c)", types.PrincipleCategoryDataMinimization},
		{"Accuracy", "Keep personal data accurate and up to date", "Article 5(1)(d)", types.PrincipleCategoryAccuracy},
		{"Storage limitation", "Keep personal data for no longer than is necessary", "Article 5(1)(e)", types.PrincipleCategoryStorageLimit},
		{"Integrity and confidentiality", "Process personal data in a manner that ensures security", "Article 5(1)(f)", types.PrincipleCategorySecurity},
		{"Accountability", "Be responsible for and demonstrate compliance with these principles", "Article 5(2)", types.PrincipleCategoryAccountability},
	}
	
	for _, principle := range standardPrinciples {
		principleObj := &types.Principle{
			ID:          uuid.New().String(),
			Title:       principle.title,
			Description: principle.description,
			Article:     principle.article,
			Category:    principle.category,
			Status:      types.PrincipleStatusNotImplemented,
			Controls:    []types.Control{},
			Owner:       "system",
			NextReview:  time.Now().UTC().AddDate(0, 12, 0),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		
		s.principles[principleObj.ID] = principleObj
		
		// Add to GDPR
		if gdpr, ok := s.gdpr[gdprID]; ok {
			gdpr.Principles = append(gdpr.Principles, *principleObj)
			gdpr.UpdatedAt = time.Now().UTC()
		}
	}
}

// initializeStandardRights loads GDPR standard rights
func (s *Service) initializeStandardRights(gdprID string) {
	standardRights := []struct {
		title       string
		description string
		article     string
		category    types.RightCategory
	}{
		{"Right to be informed", "The right to be informed about the collection and use of personal data", "Article 13 & 14", types.RightCategoryTransparency},
		{"Right of access", "The right to access personal data", "Article 15", types.RightCategoryAccess},
		{"Right to rectification", "The right to have inaccurate personal data rectified", "Article 16", types.RightCategoryRectification},
		{"Right to erasure", "The right to have personal data erased", "Article 17", types.RightCategoryErasure},
		{"Right to restrict processing", "The right to restrict processing of personal data", "Article 18", types.RightCategoryRestriction},
		{"Right to data portability", "The right to obtain and reuse personal data", "Article 20", types.RightCategoryPortability},
		{"Right to object", "The right to object to processing of personal data", "Article 21", types.RightCategoryObjection},
		{"Rights in relation to automated decision making and profiling", "Rights related to automated decision making", "Article 22", types.RightCategoryAutomatedDecision},
	}
	
	for _, right := range standardRights {
		rightObj := &types.Right{
			ID:          uuid.New().String(),
			Title:       right.title,
			Description: right.description,
			Article:     right.article,
			Category:    right.category,
			Status:      types.RightStatusNotImplemented,
			Procedures:  []types.Procedure{},
			Requests:    []types.Request{},
			Owner:       "system",
			NextReview:  time.Now().UTC().AddDate(0, 12, 0),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		
		s.rights[rightObj.ID] = rightObj
		
		// Add to GDPR
		if gdpr, ok := s.gdpr[gdprID]; ok {
			gdpr.Rights = append(gdpr.Rights, *rightObj)
			gdpr.UpdatedAt = time.Now().UTC()
		}
	}
}

// checkNotificationRequirement checks if breach notification is required
func (s *Service) checkNotificationRequirement(breach *types.Breach) {
	// Check if notification is required based on risk
	if breach.Impact.Level == types.ImpactLevelHigh || breach.Impact.Level == types.ImpactLevelCritical {
		breach.Notification.Required = true
		breach.Notification.DataSubjects = true
	}
	
	// Check for special category data
	for _, category := range breach.Categories {
		if category == types.DataCategorySpecial || category == types.DataCategorySensitive {
			breach.Notification.Required = true
			break
		}
	}
}

// Errors
var (
	ErrGDPRNotFound     = fmt.Errorf("gdpr not found")
	ErrProcessNotFound  = fmt.Errorf("process not found")
	ErrRequestNotFound  = fmt.Errorf("request not found")
	ErrBreachNotFound    = fmt.Errorf("breach not found")
	ErrDPIANotFound     = fmt.Errorf("dpia not found")
)
