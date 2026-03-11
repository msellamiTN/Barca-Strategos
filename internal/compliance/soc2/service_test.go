package soc2

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"barca-strategos/pkg/types"
)

// Mock broadcast function for testing
func mockBroadcast(data []byte) {
	// Mock implementation - in real tests this would verify the broadcast was called
}

func TestSOC2Service_AssessSOC2Compliance(t *testing.T) {
	// Create service
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create assessment scope
	scope := &types.SOC2Scope{
		Departments: []string{"IT", "Security"},
		Systems:     []string{"Phoenix Core", "Database"},
		Processes:   []string{"Incident Response"},
	}
	
	// Conduct assessment
	assessment, err := service.AssessSOC2Compliance(context.Background(), scope)
	assert.NoError(t, err)
	assert.NotNil(t, assessment)
	assert.Equal(t, "SOC 2", assessment.Framework)
	assert.Equal(t, "2017", assessment.Version)
	assert.Greater(t, assessment.OverallScore, 0.0)
	assert.NotEmpty(t, assessment.ControlAssessments)
	assert.NotEmpty(t, assessment.AssessmentID)
}

func TestSOC2Service_GetControlStatus(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Test existing control
	status, err := service.GetControlStatus(context.Background(), "CC1.1")
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, types.SOC2ControlStatusImplemented, *status)
	
	// Test non-existent control
	status, err = service.GetControlStatus(context.Background(), "NON_EXISTENT")
	assert.Error(t, err)
	assert.Nil(t, status)
}

func TestSOC2Service_UpdateControl(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create update
	update := &types.SOC2ControlUpdate{
		UpdateType: types.UpdateTypeStatus,
		UpdatedBy:  "test-user",
		Timestamp: time.Now().UTC(),
		Notes:      "compliant",
		Evidence:   []string{"New evidence document"},
	}
	
	// Update control
	err = service.UpdateControl(context.Background(), "CC1.1", update)
	assert.NoError(t, err)
}

func TestSOC2Service_GenerateSOC2Report(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create mock assessment
	assessment := &types.SOC2Assessment{
		AssessmentID:        uuid.New().String(),
		Timestamp:           time.Now().UTC(),
		Framework:           "SOC 2",
		Version:             "2017",
		Scope: types.SOC2Scope{
			Departments: []string{"IT", "Security"},
			Systems:     []string{"Phoenix Core", "Database"},
			Processes:   []string{"Incident Response"},
		},
		OverallScore:        0.85,
		ControlAssessments: []types.SOC2ControlAssessment{},
		Findings:            []types.SOC2Finding{},
		Recommendations:     []types.SOC2Recommendation{},
		LastAssessed:        time.Now().UTC(),
	}
	
	// Generate report
	report, err := service.GenerateSOC2Report(context.Background(), assessment)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, assessment.AssessmentID, report.ReportID)
	assert.NotEmpty(t, report.ReportContent)
}

func TestSOC2Service_GetSOC2Stats(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Get stats
	stats, err := service.GetSOC2Stats(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.GreaterOrEqual(t, stats.TotalControls, 0)
	assert.GreaterOrEqual(t, stats.AverageComplianceScore, 0.0)
}

func TestSOC2Service_LoadSOC2Controls(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Verify controls are loaded
	controls := service.controls
	assert.Greater(t, len(controls), 0)
	
	// Verify specific controls exist
	assert.Contains(t, controls, "CC1.1")
	assert.Contains(t, controls, "CC2.1")
	assert.Contains(t, controls, "CC3.2")
	assert.Contains(t, controls, "CC4.1")
	assert.Contains(t, controls, "CC5.1")
	assert.Contains(t, controls, "CC6.1")
	assert.Contains(t, controls, "CC7.1")
	assert.Contains(t, controls, "CC8.1")
	assert.Contains(t, controls, "CC9.1")
}

func TestSOC2Service_AssessSOC2Control(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Get a control for assessment
	control := service.controls["CC1.1"]
	assert.NotNil(t, control)
	
	// Assess the control
	assessment := service.assessSOC2Control(control)
	assert.NotNil(t, assessment)
	assert.Equal(t, "CC1.1", assessment.ControlID)
	assert.Equal(t, "Governance", assessment.ControlTitle)
	assert.Greater(t, assessment.ComplianceScore, 0.0)
	assert.NotEmpty(t, assessment.LastAssessed)
}

func TestSOC2Service_AssessImplementation(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test different control statuses
	testCases := []struct {
		status       types.SOC2ControlStatus
		expectedScore float64
	}{
		{types.SOC2ControlStatusCompliant, 1.0},
		{types.SOC2ControlStatusImplemented, 0.8},
		{types.SOC2ControlStatusPartiallyImplemented, 0.6},
		{types.SOC2ControlStatusNotImplemented, 0.0},
	}
	
	for _, tc := range testCases {
		control := &types.SOC2Control{
			Status: tc.status,
			Evidence: []string{"Test evidence"},
		}
		
		score := service.assessImplementation(control)
		assert.Equal(t, tc.expectedScore, score)
	}
}

func TestSOC2Service_AssessEffectiveness(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test different risk levels
	testCases := []struct {
		riskLevel   types.RiskLevel
		expectedScore float64
	}{
		{types.RiskLevelLow, 0.9},
		{types.RiskLevelMedium, 0.8},
		{types.RiskLevelHigh, 0.7},
		{types.RiskLevelCritical, 0.6},
	}
	
	for _, tc := range testCases {
		control := &types.SOC2Control{
			RiskLevel: tc.riskLevel,
		}
		
		score := service.assessEffectiveness(control)
		assert.Equal(t, tc.expectedScore, score)
	}
}

func TestSOC2Service_AssessFindings(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test with low compliance score
	control := &types.SOC2Control{
		ID:          "TEST",
		Title:       "Test Control",
		Description: "Test Description",
		Category:    types.SOC2ControlCategoryOperational,
		ControlType:  types.SOC2ControlTypeTechnical,
		Status:       types.SOC2ControlStatusNotImplemented,
		Evidence:     []string{},
		Owner:        "Test Owner",
		RiskLevel:    types.RiskLevelHigh,
	}
	
	findings := service.assessFindings(control, 0.3)
		assert.Len(t, findings, 1)
		assert.Equal(t, types.FindingSeverityHigh, findings[0].Severity)
		assert.Contains(t, findings[0].Description, "SOC 2 control")
		assert.NotEmpty(t, findings[0].Recommendation)
}

func TestSOC2Service_AssessRecommendations(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test with low compliance score
	control := &types.SOC2Control{
		ID:          "TEST",
		Title:       "Test Control",
		Description: "Test Description",
		Category:    types.SOC2ControlCategoryOperational,
		ControlType:  types.SOC2ControlTypeTechnical,
		Status:       types.SOC2ControlStatusNotImplemented,
		Evidence:     []string{},
		Owner:        "Test Owner",
		RiskLevel:    types.RiskLevelHigh,
	}
	
	recommendations := service.assessRecommendations(control, 0.3)
	assert.Len(t, recommendations, 1)
	assert.Equal(t, "Implement Test control completely", recommendations[0])
}

func TestSOC2Service_CalculateSOC2Score(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test with various compliance scores
	testCases := []struct {
		scores      []float64
		expectedScore float64
	}{
		{1.0, 0.8, 0.9, 0.85},
		{0.0, 0.0, 0.0, 0.0},
		{0.5, 0.6, 0.7, 0.6},
		{0.9, 0.8, 0.9, 0.85},
	}
	
	for _, tc := range testCases {
		score := service.calculateSOC2Score([]types.SOC2ControlAssessment{
			{ComplianceScore: tc.scores[0]},
			{ComplianceScore: tc.scores[1]},
			{ComplianceScore: tc.scores[2]},
			{ComplianceScore: tc.scores[3]},
		})
		assert.Equal(t, tc.expectedScore, score)
	}
}

func TestSOC2Service_GenerateSOC2Findings(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Create assessments with findings
	assessments := []types.SOC2ControlAssessment{
		{
			ControlID:       "CC1.1",
			ControlTitle:    "Governance",
			ComplianceScore:  0.3,
			Status:          types.SOC2ControlStatusNotImplemented,
			Findings: []types.SOC2Finding{
				{
					Severity:       types.FindingSeverityHigh,
					Description:     "Test finding",
					Recommendation: "Test recommendation",
					EvidenceGaps:  []string{"Evidence gap"},
				},
			},
			Recommendations: []string{"Test recommendation"},
			LastAssessed:    time.Now().UTC(),
		},
		{
			ControlID:       "CC2.1",
			ControlTitle:    "Asset Inventory",
			ComplianceScore: 0.7,
			Status:          types.SOC2ControlStatusImplemented,
			Findings: []types.SOC2Finding{},
			Recommendations: []string{},
			LastAssessed:    time.Now().UTC(),
		},
	}
	
	findings := service.generateSOC2Findings(assessments)
	assert.Len(t, findings, 1)
	assert.Equal(t, types.FindingSeverityHigh, findings[0].Severity)
}

func TestSOC2Service_GenerateSOC2Recommendations(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Create findings with different severities
	findings := []types.SOC2Finding{
		{
			Severity:       types.FindingSeverityCritical,
			Description:     "Critical finding",
			Recommendation: "Critical recommendation",
			EvidenceGaps:  []string{"Critical gap"},
		},
		{
			Severity:       types.FindingSeverityHigh,
			Description:     "High finding",
			Recommendation: "High recommendation",
			EvidenceGaps:  []string{"High gap"},
		},
		{
			Severity:       types.FindingSeverityMedium,
			Description:     "Medium finding",
			Recommendation: "Medium recommendation",
			EvidenceGaps:  []string{"Medium gap"},
		},
		{
			Severity:       types.FindingSeverityLow,
			Description:     "Low finding",
			Recommendation: "Low recommendation",
			EvidenceGaps:  []string{"Low gap"},
		},
	}
	
	recommendations := service.generateSOC2Recommendations(findings)
	assert.Len(t, recommendations, 4)
	assert.Equal(t, types.RecommendationPriorityCritical, recommendations[0].Priority)
	assert.Equal(t, "Address Critical SOC 2 Issues", recommendations[0].Title)
	assert.Equal(t, recommendations[0].Findings[0].Severity, types.FindingSeverityCritical)
}

func TestSOC2Service_IdentifyEvidenceGaps(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test governance control
	control := &types.SOC2Control{
		ID:          "CC1.1",
		Category:    types.SOC2ControlCategoryGovernance,
		Evidence:     []string{"Board meeting minutes"},
	}
	
	gaps := service.identifySOC2EvidenceGaps(control)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "Missing governance evidence", gaps[0])
	
	// Test access control control
	control = &types.SOC2Control{
		ID:          "CC3.2",
		Category:    types.SOC2ControlCategoryAccessControl,
		Evidence:     []string{"User access logs"},
	}
	
	gaps = service.identifySOC2EvidenceGaps(control)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "Missing authentication evidence", gaps[0])
	
	// Test operational control with penetration testing
	control = &types.SOC2Control{
		ID:          "CC8.1",
		Category:    types.SOC2ControlCategoryTestEvaluation,
		Title:       "Penetration Testing",
		Evidence:     []string{},
	}
	
	gaps = service.identifySOC2EvidenceGaps(control)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "Missing penetration test evidence", gaps[0])
	
	// Test control with evidence
	control = &types.SOC2Control{
		ID:          "CC2.1",
		Category:    types.SOC2ControlCategoryAssetManagement,
		Evidence:     []string{"Asset registry", "CMDB"},
	}
	
	gaps = service.identifySOC2EvidenceGaps(control)
	assert.Len(t, gaps, 0)
}

func TestSOC2Service_StartBackgroundMonitoring(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Background monitoring should be started automatically
	// In a real test, we would verify background goroutines are started
	// For now, just verify the service is properly initialized
	assert.NotNil(t, service.controls)
	assert.Greater(t, len(service.controls), 0)
}

func TestSOC2Service_MonitorSOC2Status(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Monitor status (this would normally trigger alerts for low compliance)
	service.monitorSOC2Status()
	
	// Verify service is still functional
	assert.NotNil(t, service.controls)
	assert.Greater(t, len(service.controls), 0)
}

func TestSOC2Service_PerformSOC2Assessments(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Perform assessments
	service.performSOC2Assessments()
	
	// Verify service is still functional
	assert.NotNil(t, service.controls)
	assert.Greater(t, len(service.controls), 0)
}

func TestSOC2Service_CollectSOC2Metrics(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t err)
	
	// Collect metrics
	service.collectSOC2Metrics()
	
	// Verify service is still functional
	assert.NotNil(t, service.controls)
	assert.Greater(t, len(service.controls), 0)
}

func TestSOC2Service_TriggerSOC2Alert(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t err)
	
	// Create assessment with low score
	assessment := &types.SOC2Assessment{
		OverallScore: 0.3,
		AssessmentID: uuid.New().String(),
		Timestamp:  time.Now().UTC(),
	}
	
	// Trigger alert
	service.triggerSOC2Alert(assessment)
	
	// Verify service is still functional
	assert.NotNil(t, service.controls)
	assert.Greater(t, len(service.controls), 0)
}

func TestSOC2Service_UpdateStats(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t err)
	
	// Update stats
	service.updateStats()
	
	// Verify stats are updated
	assert.NotNil(t, service.stats)
	assert.Greater(t, service.stats.TotalControls, 0)
	assert.GreaterOrEqual(t, service.stats.AverageComplianceScore, 0.0)
}

func TestSOC2Service_BroadcastSOC2Update(t *testing.T) {
	service := New()
	
	// Set mock broadcast
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err :=	service.Initialize(context.Background())
	assert.NoError(t err)
	
	// Broadcast update
	service.broadcastSOC2Update()
	
	// Verify service is functional
	assert.NotNil(t, service.controls)
	assert.Greater(t, len(service.controls), 0)
}

func TestSOC2Service_HasEvidenceContaining(t *testing.T) {
	service := New()
	
	// Test with evidence containing search term
	control := &types.SOC2Control{
		Evidence: []string{"Test evidence document"},
	}
	
	assert.True(t, service.hasEvidenceContaining(control, "evidence"))
	assert.True(t, service.hasEvidenceContaining(control, "Test"))
	assert.False(t, service.hasEvidenceContaining(control, "nonexistent"))
}

func TestSOC2Service_GetControlsByScope(t *testing.T) {
	service := New()
	
	// Initialize service
	err :=	service.Initialize(context.Background())
	assert.NoError(t err)
	
	// Get controls by scope
	controls := service.getControlsByScope(&types.SOC2Scope{
		Departments: []string{"IT", "Security"},
		Systems:     []string{"Phoenix Core", "Database"},
		Processes:   []string{"Incident Response"},
	})
	
	// Should return all controls for simple scope
	assert.Greater(t, len(controls), 0)
	assert.Equal(t, len(service.controls), len(controls))
}
