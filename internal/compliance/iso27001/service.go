package iso27001

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service manages ISO 27001 ISMS implementation
type Service struct {
	isms      map[string]*types.ISO27001
	controls  map[string]*types.Control
	clauses   map[string]*types.Clause
	risks     map[string]*types.ISMSRisk
	tests     map[string]*types.ControlTest
	findings  map[string]*types.Finding
}

// New creates an ISO27001 service
func New() *Service {
	return &Service{
		isms:     make(map[string]*types.ISO27001),
		controls: make(map[string]*types.Control),
		clauses:  make(map[string]*types.Clause),
		risks:    make(map[string]*types.ISMSRisk),
		tests:    make(map[string]*types.ControlTest),
		findings: make(map[string]*types.Finding),
	}
}

// CreateISMS creates a new ISO 27001 ISMS
func (s *Service) CreateISMS(ctx context.Context, name, description, version string) (*types.ISO27001, error) {
	isms := &types.ISO27001{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Version:     version,
		Status:      types.ISO27001StatusDraft,
		Controls:    []types.Control{},
		Clauses:     []types.Clause{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.isms[isms.ID] = isms
	log.Printf("iso27001: created ISMS %s", isms.ID)
	
	// Initialize standard controls and clauses
	s.initializeStandardControls(isms.ID)
	s.initializeStandardClauses(isms.ID)
	
	return isms, nil
}

// GetISMS retrieves an ISMS by ID
func (s *Service) GetISMS(ctx context.Context, id string) (*types.ISO27001, error) {
	isms, ok := s.isms[id]
	if !ok {
		return nil, ErrISMSNotFound
	}
	return isms, nil
}

// ListISMS returns all ISMS instances
func (s *Service) ListISMS(ctx context.Context) ([]*types.ISO27001, error) {
	var list []*types.ISO27001
	for _, isms := range s.isms {
		list = append(list, isms)
	}
	return list, nil
}

// UpdateISMS updates ISMS metadata
func (s *Service) UpdateISMS(ctx context.Context, id, name, description, version string) error {
	isms, ok := s.isms[id]
	if !ok {
		return ErrISMSNotFound
	}
	
	if name != "" {
		isms.Name = name
	}
	if description != "" {
		isms.Description = description
	}
	if version != "" {
		isms.Version = version
	}
	
	isms.UpdatedAt = time.Now().UTC()
	log.Printf("iso27001: updated ISMS %s", id)
	return nil
}

// CreateControl creates a new control
func (s *Service) CreateControl(ctx context.Context, ismsID, number, title, description string, category types.ControlCategory, owner string) (*types.Control, error) {
	control := &types.Control{
		ID:          uuid.New().String(),
		Number:      number,
		Title:       title,
		Description: description,
		Category:    category,
		Status:      types.ControlStatusNotImplemented,
		Owner:       owner,
		Evidence:    []types.Evidence{},
		Tests:       []types.ControlTest{},
		NextReview:  time.Now().UTC().AddDate(0, 6, 0), // 6 months
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.controls[control.ID] = control
	
	// Add to ISMS
	if isms, ok := s.isms[ismsID]; ok {
		isms.Controls = append(isms.Controls, *control)
		isms.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("iso27001: created control %s for ISMS %s", control.ID, ismsID)
	return control, nil
}

// UpdateControlStatus updates control implementation status
func (s *Service) UpdateControlStatus(ctx context.Context, controlID string, status types.ControlStatus) error {
	control, ok := s.controls[controlID]
	if !ok {
		return ErrControlNotFound
	}
	
	control.Status = status
	control.UpdatedAt = time.Now().UTC()
	
	log.Printf("iso27001: updated control %s status to %s", controlID, status)
	return nil
}

// AddEvidenceToControl adds evidence to a control
func (s *Service) AddEvidenceToControl(ctx context.Context, controlID string, evidence types.Evidence) error {
	control, ok := s.controls[controlID]
	if !ok {
		return ErrControlNotFound
	}
	
	control.Evidence = append(control.Evidence, evidence)
	control.UpdatedAt = time.Now().UTC()
	
	log.Printf("iso27001: added evidence to control %s", controlID)
	return nil
}

// CreateControlTest creates a new control test
func (s *Service) CreateControlTest(ctx context.Context, controlID string, testType types.TestType, description, procedure, tester string) (*types.ControlTest, error) {
	test := &types.ControlTest{
		ID:          uuid.New().String(),
		ControlID:   controlID,
		TestType:    testType,
		Description: description,
		Procedure:   procedure,
		Evidence:    []string{},
		Result:      types.TestResultNotTested,
		Score:       0,
		Tester:      tester,
		TestDate:    time.Now().UTC(),
		NextTest:    time.Now().UTC().AddDate(0, 12, 0), // 1 year
		Findings:    []types.Finding{},
		CreatedAt:   time.Now().UTC(),
	}
	
	s.tests[test.ID] = test
	
	// Add to control
	if control, ok := s.controls[controlID]; ok {
		control.Tests = append(control.Tests, *test)
		control.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("iso27001: created test %s for control %s", test.ID, controlID)
	return test, nil
}

// UpdateTestResult updates test result and score
func (s *Service) UpdateTestResult(ctx context.Context, testID string, result types.TestResult, score float64, findings []types.Finding) error {
	test, ok := s.tests[testID]
	if !ok {
		return ErrTestNotFound
	}
	
	test.Result = result
	test.Score = score
	test.Findings = findings
	test.UpdatedAt = time.Now().UTC()
	
	// Update control status based on test result
	if control, ok := s.controls[test.ControlID]; ok {
		if result == types.TestResultPass && score >= 0.8 {
			control.Status = types.ControlStatusEffective
		} else if result == types.TestResultFail {
			control.Status = types.ControlStatusNotImplemented
		} else if result == types.TestResultPartial {
			control.Status = types.ControlStatusPartiallyImplemented
		}
		control.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("iso27001: updated test %s result to %s (score: %.2f)", testID, result, score)
	return nil
}

// CreateRisk creates a new ISMS risk
func (s *Service) CreateRisk(ctx context.Context, ismsID, title, description string, threat types.Threat, vulnerability types.Vulnerability, asset types.Asset, likelihood types.Likelihood, impact types.Impact, treatment types.RiskTreatment, owner string) (*types.ISMSRisk, error) {
	riskLevel := s.calculateRiskLevel(likelihood, impact)
	
	risk := &types.ISMSRisk{
		ID:            uuid.New().String(),
		Title:         title,
		Description:   description,
		Threat:        threat,
		Vulnerability: vulnerability,
		Asset:         asset,
		Likelihood:    likelihood,
		Impact:        impact,
		RiskLevel:     riskLevel,
		Treatment:     treatment,
		Owner:         owner,
		Status:        types.RiskStatusOpen,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	
	s.risks[risk.ID] = risk
	
	// Add to ISMS risk assessment
	if isms, ok := s.isms[ismsID]; ok {
		isms.RiskAssessment.Risks = append(isms.RiskAssessment.Risks, *risk)
		isms.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("iso27001: created risk %s for ISMS %s", risk.ID, ismsID)
	return risk, nil
}

// UpdateRiskStatus updates risk status
func (s *Service) UpdateRiskStatus(ctx context.Context, riskID string, status types.RiskStatus) error {
	risk, ok := s.risks[riskID]
	if !ok {
		return ErrRiskNotFound
	}
	
	risk.Status = status
	risk.UpdatedAt = time.Now().UTC()
	
	log.Printf("iso27001: updated risk %s status to %s", riskID, status)
	return nil
}

// GetComplianceScore calculates overall ISMS compliance score
func (s *Service) GetComplianceScore(ctx context.Context, ismsID string) (float64, error) {
	isms, ok := s.isms[ismsID]
	if !ok {
		return 0, ErrISMSNotFound
	}
	
	if len(isms.Controls) == 0 {
		return 0, nil
	}
	
	totalScore := 0.0
	for _, control := range isms.Controls {
		switch control.Status {
		case types.ControlStatusEffective:
			totalScore += 1.0
		case types.ControlStatusTested:
			totalScore += 0.8
		case types.ControlStatusImplemented:
			totalScore += 0.6
		case types.ControlStatusPartiallyImplemented:
			totalScore += 0.3
		}
	}
	
	return (totalScore / float64(len(isms.Controls))) * 100, nil
}

// BroadcastUpdates sends ISMS updates via WebSocket
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":  "iso27001_update",
		"isms":  s.isms,
		"risks": s.risks,
	})
	broadcast(data)
}

// initializeStandardControls loads ISO 27001 standard controls
func (s *Service) initializeStandardControls(ismsID string) {
	standardControls := []struct {
		number      string
		title       string
		description string
		category    types.ControlCategory
	}{
		{"A.9.1.1", "Access Control Policy", "Establish, document, review, and maintain the access control policy", types.CategoryAccessControl},
		{"A.9.2.1", "User Registration and Deregistration", "Formal user registration and deregistration process", types.CategoryAccessControl},
		{"A.9.2.2", "User Access Provisioning", "Manage user access rights according to business requirements", types.CategoryAccessControl},
		{"A.10.1.1", "Cryptographic Controls Policy", "Develop and implement a cryptographic policy", types.CategoryOperationsSecurity},
		{"A.12.1.2", "Change Management", "Control changes to information processing facilities", types.CategoryOperationsSecurity},
		{"A.13.1.1", "Network Security Controls", "Network controls to protect information in networks", types.CategoryCommunicationsSecurity},
		{"A.14.2.5", "Secure System Engineering", "Security engineering principles for system development", types.CategoryAcquisition},
		{"A.16.1.7", "Response to Incidents", "Respond to information security incidents in a timely manner", types.CategoryIncidentManagement},
	}
	
	for _, ctrl := range standardControls {
		s.CreateControl(context.Background(), ismsID, ctrl.number, ctrl.title, ctrl.description, ctrl.category, "system")
	}
}

// initializeStandardClauses loads ISO 27001 standard clauses
func (s *Service) initializeStandardClauses(ismsID string) {
	standardClauses := []struct {
		number       string
		title        string
		description  string
	}{
		{"4", "Context of the Organization", "Understanding the organization and its context"},
		{"5", "Leadership", "Leadership and commitment to information security"},
		{"6", "Planning", "Information security risk assessment and treatment"},
		{"7", "Support", "Resources, competence, awareness, communication, documentation"},
		{"8", "Operation", "Planning, implementation, control of processes, procedures"},
		{"9", "Performance Evaluation", "Monitoring, measurement, internal audit, management review"},
		{"10", "Improvement", "Nonconformities, corrective actions, continual improvement"},
	}
	
	for _, clause := range standardClauses {
		clauseObj := &types.Clause{
			ID:          uuid.New().String(),
			Number:      clause.number,
			Title:       clause.title,
			Description: clause.description,
			Requirements: []types.Requirement{},
			Status:      types.ClauseStatusNotStarted,
			Owner:       "system",
			ReviewDate:  time.Now().UTC().AddDate(1, 0, 0),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		
		s.clauses[clauseObj.ID] = clauseObj
		
		// Add to ISMS
		if isms, ok := s.isms[ismsID]; ok {
			isms.Clauses = append(isms.Clauses, *clauseObj)
			isms.UpdatedAt = time.Now().UTC()
		}
	}
}

// calculateRiskLevel determines risk level from likelihood and impact
func (s *Service) calculateRiskLevel(likelihood types.Likelihood, impact types.Impact) types.RiskLevel {
	likelihoodScores := map[types.Likelihood]int{
		types.LikelihoodRare:     1,
		types.LikelihoodUnlikely: 2,
		types.LikelihoodPossible: 3,
		types.LikelihoodLikely:   4,
		types.LikelihoodCertain:  5,
	}
	
	impactScores := map[types.Impact]int{
		types.ImpactNegligible:     1,
		types.ImpactMinor:          2,
		types.ImpactModerate:       3,
		types.ImpactMajor:          4,
		types.ImpactCatastrophic:   5,
	}
	
	score := likelihoodScores[likelihood] * impactScores[impact]
	
	switch {
	case score >= 20:
		return types.RiskLevelCritical
	case score >= 12:
		return types.RiskLevelHigh
	case score >= 6:
		return types.RiskLevelMedium
	default:
		return types.RiskLevelLow
	}
}

// Errors
var (
	ErrISMSNotFound     = fmt.Errorf("isms not found")
	ErrControlNotFound  = fmt.Errorf("control not found")
	ErrTestNotFound     = fmt.Errorf("test not found")
	ErrRiskNotFound     = fmt.Errorf("risk not found")
)
