package pci_dss

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"barca-strategos/pkg/types"
)

// Mock broadcast function for testing
func mockBroadcast(data []byte) {
	// Mock implementation - in real tests this would verify the broadcast was called
}

func TestPCIDSSService_AssessPCICompliance(t *testing.T) {
	// Create service
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create assessment scope
	scope := &types.PCIScope{
		Departments: []string{"Payment Processing", "Security"},
		Systems:     []string{"Payment Gateway", "Database"},
		Processes:   []string{"Card Processing"},
	}
	
	// Conduct assessment
	assessment, err := service.AssessPCICompliance(context.Background(), scope)
	assert.NoError(t, err)
	assert.NotNil(t, assessment)
	assert.Equal(t, "PCI DSS", assessment.Framework)
	assert.Equal(t, "4.0", assessment.Version)
	assert.Greater(t, assessment.OverallScore, 0.0)
	assert.NotEmpty(t, assessment.RequirementAssessments)
	assert.NotEmpty(t, assessment.AssessmentID)
	assert.NotZero(t, assessment.NextAssessmentDate)
}

func TestPCIDSSService_GetRequirementStatus(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Test existing requirement
	status, err := service.GetRequirementStatus(context.Background(), "1.1")
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, types.PCIControlStatusImplemented, *status)
	
	// Test non-existent requirement
	status, err = service.GetRequirementStatus(context.Background(), "NON_EXISTENT")
	assert.Error(t, err)
	assert.Nil(t, status)
}

func TestPCIDSSService_UpdateRequirement(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create update
	update := &types.PCIControlUpdate{
		UpdateType: types.UpdateTypeStatus,
		UpdatedBy:  "test-user",
		Timestamp: time.Now().UTC(),
		Notes:      "compliant",
		Evidence:   []string{"New evidence document"},
	}
	
	// Update requirement
	err = service.UpdateRequirement(context.Background(), "1.1", update)
	assert.NoError(t, err)
}

func TestPCIDSSService_GeneratePCIReport(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create mock assessment
	assessment := &types.PCIAssessment{
		AssessmentID:          uuid.New().String(),
		Timestamp:             time.Now().UTC(),
		Framework:             "PCI DSS",
		Version:               "4.0",
		Scope: types.PCIScope{
			Departments: []string{"Payment Processing", "Security"},
			Systems:     []string{"Payment Gateway", "Database"},
			Processes:   []string{"Card Processing"},
		},
		OverallScore:          0.85,
		RequirementAssessments: []types.PCIRequirementAssessment{},
		Findings:              []types.PCIFinding{},
		Recommendations:       []types.PCIRecommendation{},
		NextAssessmentDate:     time.Now().UTC().AddDate(1, 0, 0),
	}
	
	// Generate report
	report, err := service.GeneratePCIReport(context.Background(), assessment)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, assessment.AssessmentID, report.ReportID)
	assert.NotEmpty(t, report.ReportContent)
}

func TestPCIDSSService_GetPCIStats(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Get stats
	stats, err := service.GetPCIStats(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.GreaterOrEqual(t, stats.TotalRequirements, 0)
	assert.GreaterOrEqual(t, stats.AverageComplianceScore, 0.0)
}

func TestPCIDSSService_LoadPCIRequirements(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Verify requirements are loaded
	requirements := service.requirements
	assert.Greater(t, len(requirements), 0)
	
	// Verify specific requirements exist
	assert.Contains(t, requirements, "1.1")
	assert.Contains(t, requirements, "2.1")
	assert.Contains(t, requirements, "3.1")
	assert.Contains(t, requirements, "4.1")
	assert.Contains(t, requirements, "5.1")
	assert.Contains(t, requirements, "6.1")
	assert.Contains(t, requirements, "7.1")
	assert.Contains(t, requirements, "8.1")
	assert.Contains(t, requirements, "9.1")
	assert.Contains(t, requirements, "10.1")
	assert.Contains(t, requirements, "11.1")
	assert.Contains(t, requirements, "12.1")
}

func TestPCIDSSService_AssessPCIRequirement(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Get a requirement for assessment
	requirement := service.requirements["1.1"]
	assert.NotNil(t, requirement)
	
	// Assess the requirement
	assessment := service.assessPCIRequirement(requirement)
	assert.NotNil(t, assessment)
	assert.Equal(t, "1.1", assessment.RequirementID)
	assert.Equal(t, "Network Security Controls", assessment.RequirementTitle)
	assert.Greater(t, assessment.ComplianceScore, 0.0)
	assert.NotEmpty(t, assessment.LastAssessed)
}

func TestPCIDSSService_AssessImplementation(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test different requirement statuses
	testCases := []struct {
		status       types.PCIControlStatus
		expectedScore float64
	}{
		{types.PCIControlStatusCompliant, 1.0},
		{types.PCIControlStatusImplemented, 0.8},
		{types.PCIControlStatusPartiallyImplemented, 0.6},
		{types.PCIControlStatusNotImplemented, 0.0},
	}
	
	for _, tc := range testCases {
		requirement := &types.PCIRequirement{
			Status: tc.status,
			Evidence: []string{"Test evidence"},
		}
		
		score := service.assessImplementation(requirement)
		assert.Equal(t, tc.expectedScore, score)
	}
}

func TestPCIDSSService_AssessEffectiveness(t *testing.T) {
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
		requirement := &types.PCIRequirement{
			RiskLevel: tc.riskLevel,
		}
		
		score := service.assessEffectiveness(requirement)
		assert.Equal(t, tc.expectedScore, score)
	}
}

func TestPCIDSSService_AssessFindings(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test with low compliance score
	requirement := &types.PCIRequirement{
		ID:          "TEST",
		Title:       "Test Requirement",
		Description: "Test Description",
		Category:    types.PCIControlCategoryNetworkSecurity,
		ControlType:  types.PCIControlTypeTechnical,
		Status:       types.PCIControlStatusNotImplemented,
		Evidence:     []string{},
		Owner:        "Test Owner",
		RiskLevel:    types.RiskLevelHigh,
	}
	
	findings := service.assessFindings(requirement, 0.3)
	assert.Len(t, findings, 1)
	assert.Equal(t, types.FindingSeverityCritical, findings[0].Severity)
	assert.Contains(t, findings[0].Description, "PCI DSS requirement")
	assert.NotEmpty(t, findings[0].Recommendation)
}

func TestPCIDSSService_AssessRecommendations(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test with low compliance score
	requirement := &types.PCIRequirement{
		ID:          "TEST",
		Title:       "Test Requirement",
		Description: "Test Description",
		Category:    types.PCIControlCategoryNetworkSecurity,
		ControlType:  types.PCIControlTypeTechnical,
		Status:       types.PCIControlStatusNotImplemented,
		Evidence:     []string{},
		Owner:        "Test Owner",
		RiskLevel:    types.RiskLevelHigh,
	}
	
	recommendations := service.assessRecommendations(requirement, 0.3)
	assert.Len(t, recommendations, 1)
	assert.Equal(t, "Implement Test requirement completely", recommendations[0])
}

func TestPCIDSSService_CalculatePCIScore(t *testing.T) {
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
		score := service.calculatePCIScore([]types.PCIRequirementAssessment{
			{ComplianceScore: tc.scores[0]},
			{ComplianceScore: tc.scores[1]},
			{ComplianceScore: tc.scores[2]},
			{ComplianceScore: tc.scores[3]},
		})
		assert.Equal(t, tc.expectedScore, score)
	}
}

func TestPCIDSSService_GeneratePCIFindings(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Create assessments with findings
	assessments := []types.PCIRequirementAssessment{
		{
			RequirementID:       "1.1",
			RequirementTitle:    "Network Security Controls",
			ComplianceScore:     0.3,
			Status:              types.PCIControlStatusNotImplemented,
			Findings: []types.PCIFinding{
				{
					Severity:       types.FindingSeverityCritical,
					Description:     "Test finding",
					Recommendation: "Test recommendation",
					EvidenceGaps:  []string{"Evidence gap"},
				},
			},
			Recommendations: []string{"Test recommendation"},
			LastAssessed:    time.Now().UTC(),
		},
		{
			RequirementID:       "2.1",
			RequirementTitle:    "Secure Configurations",
			ComplianceScore:     0.7,
			Status:              types.PCIControlStatusPartiallyImplemented,
			Findings: []types.PCIFinding{},
			Recommendations: []string{},
			LastAssessed:    time.Now().UTC(),
		},
	}
	
	findings := service.generatePCIFindings(assessments)
	assert.Len(t, findings, 1)
	assert.Equal(t, types.FindingSeverityCritical, findings[0].Severity)
}

func TestPCIDSSService_GeneratePCIRecommendations(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Create findings with different severities
	findings := []types.PCIFinding{
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
	
	recommendations := service.generatePCIRecommendations(findings)
	assert.Len(t, recommendations, 4)
	assert.Equal(t, types.RecommendationPriorityCritical, recommendations[0].Priority)
	assert.Equal(t, "Address Critical PCI DSS Issues", recommendations[0].Title)
	assert.Equal(t, recommendations[0].Findings[0].Severity, types.FindingSeverityCritical)
}

func TestPCIDSSService_IdentifyEvidenceGaps(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Test data protection control
	requirement := &types.PCIRequirement{
		ID:          "3.1",
		Category:    types.PCIControlCategoryDataProtection,
		Evidence:     []string{"Key management logs"},
	}
	
	gaps := service.identifyEvidenceGaps(requirement)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "Missing encryption evidence", gaps[0])
	
	// Test access control control
	requirement = &types.PCIRequirement{
		ID:          "7.1",
		Category:    types.PCIControlCategoryAccessControl,
		Evidence:     []string{"User access logs"},
	}
	
	gaps = service.identifyEvidenceGaps(requirement)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "Missing access control evidence", gaps[0])
	
	// Test network security control
	requirement = &types.PCIRequirement{
		ID:          "1.1",
		Category:    types.PCIControlCategoryNetworkSecurity,
		Evidence:     []string{},
	}
	
	gaps = service.identifyEvidenceGaps(requirement)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "Missing firewall evidence", gaps[0])
	
	// Test requirement with evidence
	requirement = &types.PCIRequirement{
		ID:          "5.1",
		Category:    types.PCIControlCategoryMalwareProtection,
		Evidence:     []string{"Antivirus reports", "Malware scan logs"},
	}
	
	gaps = service.identifyEvidenceGaps(requirement)
	assert.Len(t, gaps, 0)
}

func TestPCIDSSService_StartBackgroundMonitoring(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Background monitoring should be started automatically
	// In a real test, we would verify background goroutines are started
	// For now, just verify the service is properly initialized
	assert.NotNil(t, service.requirements)
	assert.Greater(t, len(service.requirements), 0)
}

func TestPCIDSSService_MonitorPCIStatus(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Monitor status (this would normally trigger alerts for low compliance)
	service.monitorPCIStatus()
	
	// Verify service is still functional
	assert.NotNil(t, service.requirements)
	assert.Greater(t, len(service.requirements), 0)
}

func TestPCIDSSService_PerformPCIAssessments(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Perform assessments
	service.performPCIAssessments()
	
	// Verify service is still functional
	assert.NotNil(t, service.requirements)
	assert.Greater(t, len(service.requirements), 0)
}

func TestPCIDSSService_CollectPCIMetrics(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Collect metrics
	service.collectPCIMetrics()
	
	// Verify service is still functional
	assert.NotNil(t, service.requirements)
	assert.Greater(t, len(service.requirements), 0)
}

func TestPCIDSSService_TriggerPCIAlert(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create assessment with low score
	assessment := &types.PCIAssessment{
		OverallScore: 0.3,
		AssessmentID: uuid.New().String(),
		Timestamp:  time.Now().UTC(),
	}
	
	// Trigger alert
	service.triggerPCIAlert(assessment)
	
	// Verify service is still functional
	assert.NotNil(t, service.requirements)
	assert.Greater(t, len(service.requirements), 0)
}

func TestPCIDSSService_UpdateStats(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Update stats
	service.updateStats()
	
	// Verify stats are updated
	assert.NotNil(t, service.stats)
	assert.Greater(t, service.stats.TotalRequirements, 0)
	assert.GreaterOrEqual(t, service.stats.AverageComplianceScore, 0.0)
}

func TestPCIDSSService_BroadcastPCIUpdate(t *testing.T) {
	service := New()
	
	// Set mock broadcast
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Broadcast update
	service.broadcastPCIUpdate()
	
	// Verify service is functional
	assert.NotNil(t, service.requirements)
	assert.Greater(t, len(service.requirements), 0)
}

func TestPCIDSSService_HasEvidenceContaining(t *testing.T) {
	service := New()
	
	// Test with evidence containing search term
	requirement := &types.PCIRequirement{
		Evidence: []string{"Test evidence document"},
	}
	
	assert.True(t, service.hasEvidenceContaining(requirement, "evidence"))
	assert.True(t, service.hasEvidenceContaining(requirement, "Test"))
	assert.False(t, service.hasEvidenceContaining(requirement, "nonexistent"))
}

func TestPCIDSSService_GetRequirementsByScope(t *testing.T) {
	service := New()
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Get requirements by scope
	requirements := service.getRequirementsByScope(&types.PCIScope{
		Departments: []string{"Payment Processing", "Security"},
		Systems:     []string{"Payment Gateway", "Database"},
		Processes:   []string{"Card Processing"},
	})
	
	// Should return all requirements for simple scope
	assert.Greater(t, len(requirements), 0)
	assert.Equal(t, len(service.requirements), len(requirements))
}

func TestPCIDSSService_ControlTypeAdjustments(t *testing.T) {
	service := New()
	
	// Test different control types
	testCases := []struct {
		controlType  types.PCIControlType
		expectedScore float64
	}{
		{types.PCIControlTypeOrganizational, 0.0},
		{types.PCIControlTypeTechnical, 0.1},
		{types.PCIControlTypeOperational, 0.0},
		{types.PCIControlTypePhysical, 0.05},
	}
	
	baseScore := 0.8
	
	for _, tc := range testCases {
		requirement := &types.PCIRequirement{
			RiskLevel: types.RiskLevelMedium,
			ControlType: tc.controlType,
		}
		
		score := service.assessEffectiveness(requirement)
		expected := baseScore + tc.expectedScore
		assert.Equal(t, expected, score)
	}
}

func TestPCIDSSService_EvidenceQualityScoring(t *testing.T) {
	service := New()
	
	// Test different evidence counts
	testCases := []struct {
		evidenceCount int
		expectedScore float64
	}{
		{0, 0.5},
		{1, 0.7},
		{2, 0.7},
		{3, 0.9},
		{5, 0.9},
	}
	
	baseStatus := types.PCIControlStatusImplemented
	baseScore := 0.8
	
	for _, tc := range testCases {
		requirement := &types.PCIRequirement{
			Status: baseStatus,
			Evidence: make([]string, tc.evidenceCount),
		}
		
		score := service.assessImplementation(requirement)
		expected := baseScore * tc.expectedScore
		assert.Equal(t, expected, score)
	}
}

func TestPCIDSSService_FindingSeverityCalculation(t *testing.T) {
	service := New()
	
	// Test different compliance scores
	testCases := []struct {
		complianceScore float64
		expectedSeverity types.FindingSeverity
	}{
		{0.9, types.FindingSeverityLow},
		{0.7, types.FindingSeverityLow},
		{0.5, types.FindingSeverityHigh},
		{0.3, types.FindingSeverityCritical},
		{0.1, types.FindingSeverityCritical},
	}
	
	for _, tc := range testCases {
		requirement := &types.PCIRequirement{
			ID:          "TEST",
			Title:       "Test Requirement",
			Description: "Test Description",
			Category:    types.PCIControlCategoryNetworkSecurity,
			ControlType:  types.PCIControlTypeTechnical,
			Status:       types.PCIControlStatusNotImplemented,
			Evidence:     []string{},
			Owner:        "Test Owner",
			RiskLevel:    types.RiskLevelHigh,
		}
		
		findings := service.assessFindings(requirement, tc.complianceScore)
		assert.Len(t, findings, 1)
		assert.Equal(t, tc.expectedSeverity, findings[0].Severity)
	}
}

func TestPCIDSSService_RecommendationPriorityCalculation(t *testing.T) {
	service := New()
	
	// Test different compliance scores
	testCases := []struct {
		complianceScore float64
		expectedPriority types.RecommendationPriority
	}{
		{0.9, types.RecommendationPriorityLow},
		{0.7, types.RecommendationPriorityLow},
		{0.5, types.RecommendationPriorityMedium},
		{0.3, types.RecommendationPriorityHigh},
		{0.1, types.RecommendationPriorityHigh},
	}
	
	for _, tc := range testCases {
		requirement := &types.PCIRequirement{
			ID:          "TEST",
			Title:       "Test Requirement",
			Description: "Test Description",
			Category:    types.PCIControlCategoryNetworkSecurity,
			ControlType:  types.PCIControlTypeTechnical,
			Status:       types.PCIControlStatusNotImplemented,
			Evidence:     []string{},
			Owner:        "Test Owner",
			RiskLevel:    types.RiskLevelHigh,
		}
		
		recommendations := service.assessRecommendations(requirement, tc.complianceScore)
		assert.Len(t, recommendations, 1)
		assert.Equal(t, tc.expectedPriority, recommendations[0])
	}
}

func TestPCIDSSService_NextAssessmentDate(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create assessment scope
	scope := &types.PCIScope{
		Departments: []string{"Payment Processing", "Security"},
		Systems:     []string{"Payment Gateway", "Database"},
		Processes:   []string{"Card Processing"},
	}
	
	// Conduct assessment
	assessment, err := service.AssessPCICompliance(context.Background(), scope)
	assert.NoError(t, err)
	
	// Verify next assessment date is set (1 year from now)
	expectedNextAssessment := time.Now().UTC().AddDate(1, 0, 0)
	actualNextAssessment := assessment.NextAssessment
	
	// Allow for small time differences (within 1 minute)
	timeDiff := actualNextAssessment.Sub(expectedNextAssessment)
	assert.True(t, timeDiff.Abs() < time.Minute)
}

func TestPCIDSSService_ReportContentGeneration(t *testing.T) {
	service := New()
	service.SetBroadcast(mockBroadcast)
	
	// Initialize service
	err := service.Initialize(context.Background())
	assert.NoError(t, err)
	
	// Create mock assessment
	assessment := &types.PCIAssessment{
		AssessmentID:          uuid.New().String(),
		Timestamp:             time.Now().UTC(),
		Framework:             "PCI DSS",
		Version:               "4.0",
		Scope: types.PCIScope{
			Departments: []string{"Payment Processing", "Security"},
			Systems:     []string{"Payment Gateway", "Database"},
			Processes:   []string{"Card Processing"},
		},
		OverallScore:          0.85,
		RequirementAssessments: []types.PCIRequirementAssessment{},
		Findings:              []types.PCIFinding{},
		Recommendations:       []types.PCIRecommendation{},
		NextAssessmentDate:     time.Now().UTC().AddDate(1, 0, 0),
	}
	
	// Generate report content
	content := service.generateReportContent(assessment)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "PCI DSS Compliance Report")
	assert.Contains(t, content, assessment.AssessmentID)
	assert.Contains(t, content, "85%")
	assert.Contains(t, content, assessment.Timestamp.Format("2006-01-02 15:04:05"))
	assert.Contains(t, content, assessment.NextAssessmentDate.Format("2006-01-02"))
}
