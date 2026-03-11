package nist

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service manages NIST CSF implementation
type Service struct {
	csf         map[string]*types.NISTCSF
	functions   map[string]*types.Function
	subcategories map[string]*types.SubCategory
	profiles    map[string]*types.Profile
	assessments map[string]*types.Assessment
}

// New creates a NIST service
func New() *Service {
	return &Service{
		csf:           make(map[string]*types.NISTCSF),
		functions:     make(map[string]*types.Function),
		subcategories: make(map[string]*types.SubCategory),
		profiles:      make(map[string]*types.Profile),
		assessments:   make(map[string]*types.Assessment),
	}
}

// CreateCSF creates a new NIST CSF implementation
func (s *Service) CreateCSF(ctx context.Context, name, description, version string) (*types.NISTCSF, error) {
	csf := &types.NISTCSF{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Version:     version,
		Status:      types.NISTStatusDraft,
		Functions:   []types.Function{},
		Profiles:    []types.Profile{},
		Assessments: []types.Assessment{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.csf[csf.ID] = csf
	log.Printf("nist: created CSF %s", csf.ID)
	
	// Initialize standard functions
	s.initializeStandardFunctions(csf.ID)
	
	return csf, nil
}

// GetCSF retrieves a CSF by ID
func (s *Service) GetCSF(ctx context.Context, id string) (*types.NISTCSF, error) {
	csf, ok := s.csf[id]
	if !ok {
		return nil, ErrCSFNotFound
	}
	return csf, nil
}

// ListCSFs returns all CSF instances
func (s *Service) ListCSFs(ctx context.Context) ([]*types.NISTCSF, error) {
	var list []*types.NISTCSF
	for _, csf := range s.csf {
		list = append(list, csf)
	}
	return list, nil
}

// CreateFunction creates a new NIST function
func (s *Service) CreateFunction(ctx context.Context, csfID, name, description string, functionType types.FunctionType, owner string) (*types.Function, error) {
	function := &types.Function{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Category:    functionType,
		Categories:  []types.SubCategory{},
		Status:      types.FunctionStatusNotStarted,
		Progress:    0,
		Owner:       owner,
		NextReview:  time.Now().UTC().AddDate(0, 6, 0),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.functions[function.ID] = function
	
	// Add to CSF
	if csf, ok := s.csf[csfID]; ok {
		csf.Functions = append(csf.Functions, *function)
		csf.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("nist: created function %s for CSF %s", function.ID, csfID)
	return function, nil
}

// CreateSubCategory creates a new subcategory
func (s *Service) CreateSubCategory(ctx context.Context, functionID, number, title, description string, priority types.Priority, owner string) (*types.SubCategory, error) {
	subcategory := &types.SubCategory{
		ID:                    uuid.New().String(),
		Number:                number,
		Title:                 title,
		Description:            description,
		InformativeReferences: []string{},
		ImplementationStatus:  types.ImplementationStatusNotImplemented,
		Measurements:          []types.Measurement{},
		Resources:             []types.Resource{},
		Priority:              priority,
		Owner:                 owner,
		DueDate:               time.Now().UTC().AddDate(0, 12, 0),
		CreatedAt:             time.Now().UTC(),
		UpdatedAt:             time.Now().UTC(),
	}
	
	s.subcategories[subcategory.ID] = subcategory
	
	// Add to function
	if function, ok := s.functions[functionID]; ok {
		function.Categories = append(function.Categories, *subcategory)
		function.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("nist: created subcategory %s for function %s", subcategory.ID, functionID)
	return subcategory, nil
}

// UpdateSubCategoryStatus updates subcategory implementation status
func (s *Service) UpdateSubCategoryStatus(ctx context.Context, subCategoryID string, status types.ImplementationStatus) error {
	subcategory, ok := s.subcategories[subCategoryID]
	if !ok {
		return ErrSubCategoryNotFound
	}
	
	subcategory.ImplementationStatus = status
	subcategory.UpdatedAt = time.Now().UTC()
	
	// Update function progress
	s.updateFunctionProgress(subcategory)
	
	log.Printf("nist: updated subcategory %s status to %s", subCategoryID, status)
	return nil
}

// CreateProfile creates a new NIST CSF profile
func (s *Service) CreateProfile(ctx context.Context, csfID, name, description string, profileType types.ProfileType) (*types.Profile, error) {
	profile := &types.Profile{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Type:        profileType,
		Functions:   []types.Function{},
		Scopes:      []types.Scope{},
		Targets:     []types.Target{},
		Status:      types.ProfileStatusDraft,
		Version:     "1.0",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.profiles[profile.ID] = profile
	
	// Add to CSF
	if csf, ok := s.csf[csfID]; ok {
		csf.Profiles = append(csf.Profiles, *profile)
		csf.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("nist: created profile %s for CSF %s", profile.ID, csfID)
	return profile, nil
}

// CreateAssessment creates a new assessment
func (s *Service) CreateAssessment(ctx context.Context, profileID, title, description string, assessmentType types.AssessmentType, method types.AssessmentMethod, scope, assessor string) (*types.Assessment, error) {
	assessment := &types.Assessment{
		ID:          uuid.New().String(),
		ProfileID:   profileID,
		Title:       title,
		Description: description,
		Type:        assessmentType,
		Method:      method,
		Scope:       scope,
		Assessor:    assessor,
		StartDate:   time.Now().UTC(),
		EndDate:     time.Now().UTC().AddDate(0, 0, 30), // 30 days
		Status:      types.AssessmentStatusPlanned,
		Results:     []types.AssessmentResult{},
		Score:       0,
		Findings:    []types.Finding{},
		Recommendations: []types.Recommendation{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	s.assessments[assessment.ID] = assessment
	
	// Add to profile
	if profile, ok := s.profiles[profileID]; ok {
		profile.Assessments = append(profile.Assessments, *assessment)
		profile.UpdatedAt = time.Now().UTC()
	}
	
	log.Printf("nist: created assessment %s for profile %s", assessment.ID, profileID)
	return assessment, nil
}

// UpdateAssessmentResult updates assessment results
func (s *Service) UpdateAssessmentResult(ctx context.Context, assessmentID string, results []types.AssessmentResult, findings []types.Finding, recommendations []types.Recommendation) error {
	assessment, ok := s.assessments[assessmentID]
	if !ok {
		return ErrAssessmentNotFound
	}
	
	assessment.Results = results
	assessment.Findings = findings
	assessment.Recommendations = recommendations
	assessment.Score = s.calculateAssessmentScore(results)
	assessment.Status = types.AssessmentStatusCompleted
	assessment.EndDate = time.Now().UTC()
	assessment.UpdatedAt = time.Now().UTC()
	
	log.Printf("nist: updated assessment %s with score %.2f", assessmentID, assessment.Score)
	return nil
}

// GetMaturityScore calculates overall maturity score
func (s *Service) GetMaturityScore(ctx context.Context, csfID string) (float64, error) {
	csf, ok := s.csf[csfID]
	if !ok {
		return 0, ErrCSFNotFound
	}
	
	if len(csf.Functions) == 0 {
		return 0, nil
	}
	
	totalScore := 0.0
	for _, function := range csf.Functions {
		totalScore += function.Progress
	}
	
	return totalScore / float64(len(csf.Functions)), nil
}

// BroadcastUpdates sends CSF updates via WebSocket
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type": "nist_update",
		"csf":  s.csf,
	})
	broadcast(data)
}

// initializeStandardFunctions loads NIST CSF standard functions
func (s *Service) initializeStandardFunctions(csfID string) {
	standardFunctions := []struct {
		name        string
		description string
		functionType types.FunctionType
	}{
		{"Identify", "Develop an understanding of business context, resources, and risk", types.FunctionTypeIdentify},
		{"Protect", "Implement safeguards to ensure delivery of critical infrastructure services", types.FunctionTypeProtect},
		{"Detect", "Implement activities to identify the occurrence of a cybersecurity event", types.FunctionTypeDetect},
		{"Respond", "Implement activities to take action regarding a detected cybersecurity event", types.FunctionTypeRespond},
		{"Recover", "Implement plans for resilience and restoration of capabilities impaired by a cybersecurity event", types.FunctionTypeRecover},
	}
	
	for _, fn := range standardFunctions {
		s.CreateFunction(context.Background(), csfID, fn.name, fn.description, fn.functionType, "system")
	}
}

// updateFunctionProgress updates function progress based on subcategories
func (s *Service) updateFunctionProgress(subcategory *types.SubCategory) {
	// Find the function containing this subcategory
	for _, function := range s.functions {
		for _, cat := range function.Categories {
			if cat.ID == subcategory.ID {
				// Calculate progress based on implementation status
				total := len(function.Categories)
				if total == 0 {
					return
				}
				
				completed := 0
				for _, cat := range function.Categories {
					switch cat.ImplementationStatus {
					case types.ImplementationStatusImplemented, types.ImplementationStatusTested, types.ImplementationStatusValidated:
						completed++
					case types.ImplementationStatusPartiallyImplemented:
						completed += 0.5
					}
				}
				
				function.Progress = (float64(completed) / float64(total)) * 100
				function.UpdatedAt = time.Now().UTC()
				
				// Update function status based on progress
				switch {
				case function.Progress >= 90:
					function.Status = types.FunctionStatusOptimized
				case function.Progress >= 75:
					function.Status = types.FunctionStatusComplete
				case function.Progress >= 50:
					function.Status = types.FunctionStatusPartiallyComplete
				case function.Progress > 0:
					function.Status = types.FunctionStatusInProgress
				default:
					function.Status = types.FunctionStatusNotStarted
				}
				
				return
			}
		}
	}
}

// calculateAssessmentScore calculates overall assessment score
func (s *Service) calculateAssessmentScore(results []types.AssessmentResult) float64 {
	if len(results) == 0 {
		return 0
	}
	
	totalScore := 0.0
	for _, result := range results {
		// Convert rating level to numeric score (0-5)
		score := 0.0
		switch result.Rating {
		case types.RatingLevelLevel0:
			score = 0
		case types.RatingLevelLevel1:
			score = 1
		case types.RatingLevelLevel2:
			score = 2
		case types.RatingLevelLevel3:
			score = 3
		case types.RatingLevelLevel4:
			score = 4
		case types.RatingLevelLevel5:
			score = 5
		}
		totalScore += score
	}
	
	return (totalScore / float64(len(results))) / 5.0 * 100 // Convert to 0-100 scale
}

// Errors
var (
	ErrCSFNotFound         = fmt.Errorf("csf not found")
	ErrSubCategoryNotFound = fmt.Errorf("subcategory not found")
	ErrAssessmentNotFound  = fmt.Errorf("assessment not found")
)
