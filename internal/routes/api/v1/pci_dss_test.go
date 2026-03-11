package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"barca-strategos/internal/compliance/pci_dss"
	"barca-strategos/pkg/types"
)

type PCIDSSAPITestSuite struct {
	suite.Suite
	app      *fiber.App
	pciSvc   *pci_dss.Service
	broadcast func([]byte)
}

func (suite *PCIDSSAPITestSuite) SetupSuite() {
	suite.app = fiber.New()
	suite.pciSvc = pci_dss.New()
	suite.broadcast = func(data []byte) {
		// Mock broadcast function
	}
	suite.pciSvc.SetBroadcast(suite.broadcast)
	
	// Register routes
	RegisterPCIDSSRoutes(suite.app, suite.pciSvc, suite.broadcast)
}

func (suite *PCIDSSAPITestSuite) TestAssessPCICompliance() {
	// Create test scope
	scope := types.PCIScope{
		Departments: []string{"Payment Processing", "Security"},
		Systems:     []string{"Payment Gateway", "Database"},
		Processes:   []string{"Card Processing"},
	}
	
	// Convert to JSON
	scopeJSON, _ := json.Marshal(scope)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/assess", bytes.NewReader(scopeJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var assessment types.PCIAssessment
	err = json.NewDecoder(resp.Body).Decode(&assessment)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("PCI DSS", assessment.Framework)
	suite.Equal("4.0", assessment.Version)
	suite.Greater(assessment.OverallScore, 0.0)
	suite.NotEmpty(assessment.AssessmentID)
	suite.NotZero(assessment.NextAssessmentDate)
}

func (suite *PCIDSSAPITestSuite) TestAssessPCIComplianceInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/assess", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIRequirements() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(requirements), 0)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIRequirement() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements/1.1", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var requirement types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirement)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("1.1", requirement.ID)
	suite.NotEmpty(requirement.Title)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIRequirementNotFound() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements/NOT_FOUND", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var requirement types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirement)
	suite.NoError(err)
	
	// Verify response (should return mock data)
	suite.Equal("NOT_FOUND", requirement.ID)
}

func (suite *PCIDSSAPITestSuite) TestUpdatePCIRequirement() {
	// Create update payload
	update := types.PCIControlUpdate{
		UpdateType: types.UpdateTypeStatus,
		UpdatedBy:  "test-user",
		Timestamp:  suite.pciSvc.(*pci_dss.Service).(*pci_dss.Service).requirements["1.1"].LastReviewDate,
		Notes:      "compliant",
		Evidence:   []string{"New evidence document"},
	}
	
	// Convert to JSON
	updateJSON, _ := json.Marshal(update)
	
	// Create request
	req := httptest.NewRequest("PUT", "/api/v1/compliance/pci_dss/requirement/1.1", bytes.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestUpdatePCIRequirementInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("PUT", "/api/v1/compliance/pci_dss/requirement/1.1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIRequirementStatus() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirement/1.1/status", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var status types.PCIControlStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	suite.NoError(err)
	
	// Verify response
	suite.NotEmpty(status)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIRequirementStatusNotFound() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirement/NOT_FOUND/status", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestGeneratePCIReport() {
	// Create assessment payload
	assessment := types.PCIAssessment{
		AssessmentID:          "test-assessment",
		Timestamp:             suite.pciSvc.(*pci_dss.Service).(*pci_dss.Service).requirements["1.1"].LastReviewDate,
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
		NextAssessmentDate:     suite.pciSvc.(*pci_dss.Service).(*pci_dss.Service).requirements["1.1"].LastReviewDate.AddDate(1, 0, 0),
	}
	
	// Convert to JSON
	assessmentJSON, _ := json.Marshal(assessment)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/report", bytes.NewReader(assessmentJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var report types.PCIReport
	err = json.NewDecoder(resp.Body).Decode(&report)
	suite.NoError(err)
	
	// Verify response
	suite.Equal(assessment.AssessmentID, report.ReportID)
	suite.NotEmpty(report.ReportContent)
}

func (suite *PCIDSSAPITestSuite) TestGeneratePCIReportInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/report", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIStats() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/stats", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var stats types.PCIStats
	err = json.NewDecoder(resp.Body).Decode(&stats)
	suite.NoError(err)
	
	// Verify response
	suite.GreaterOrEqual(stats.TotalRequirements, 0)
	suite.GreaterOrEqual(stats.AverageComplianceScore, 0.0)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIAssessments() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/assessments", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var assessments []types.PCIAssessment
	err = json.NewDecoder(resp.Body).Decode(&assessments)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(assessments), 0)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIAssessment() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/assessments/assessment-1", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var assessment types.PCIAssessment
	err = json.NewDecoder(resp.Body).Decode(&assessment)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("assessment-1", assessment.AssessmentID)
	suite.Equal("PCI DSS", assessment.Framework)
}

func (suite *PCIDSSAPITestSuite) TestSearchPCIRequirements() {
	// Create search request
	search := struct {
		Query        string   `json:"query"`
		Categories  []string `json:"categories"`
		Status      []string `json:"status"`
		Types       []string `json:"types"`
	}{
		Query:       "network",
		Categories: []string{"network_security"},
		Status:      []string{"implemented"},
		Types:       []string{"technical"},
	}
	
	// Convert to JSON
	searchJSON, _ := json.Marshal(search)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/search", bytes.NewReader(searchJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(requirements), 0)
}

func (suite *PCIDSSAPITestSuite) TestSearchPCIRequirementsInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/search", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestBulkPCIOperations() {
	// Create bulk operation request
	bulk := struct {
		Action         string   `json:"action"`
		RequirementIDs []string `json:"requirement_ids"`
		Parameters    map[string]interface{} `json:"parameters"`
	}{
		Action:         "assess",
		RequirementIDs: []string{"1.1", "2.1"},
		Parameters:     map[string]interface{}{},
	}
	
	// Convert to JSON
	bulkJSON, _ := json.Marshal(bulk)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/bulk", bytes.NewReader(bulkJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	
	// Verify response
	suite.Contains(response, "results")
}

func (suite *PCIDSSAPITestSuite) TestBulkPCIOperationsInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/bulk", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIRequirementEvidence() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirement/1.1/evidence", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var evidence []string
	err = json.NewDecoder(resp.Body).Decode(&evidence)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(evidence), 0)
}

func (suite *PCIDSSAPITestSuite) TestUploadPCIRequirementEvidence() {
	// Create evidence upload request
	upload := struct {
		Evidence []string `json:"evidence"`
	}{
		Evidence: []string{"New evidence document 1", "New evidence document 2"},
	}
	
	// Convert to JSON
	uploadJSON, _ := json.Marshal(upload)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/requirement/1.1/evidence", bytes.NewReader(uploadJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	
	// Verify response
	suite.Contains(response, "message")
	suite.Contains(response, "count")
}

func (suite *PCIDSSAPITestSuite) TestUploadPCIRequirementEvidenceInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/requirement/1.1/evidence", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIReports() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/reports", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var reports map[string]types.PCIReport
	err = json.NewDecoder(resp.Body).Decode(&reports)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(reports), 0)
}

func (suite *PCIDSSAPITestSuite) TestGetPCIReport() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/reports/report-1", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var report types.PCIReport
	err = json.NewDecoder(resp.Body).Decode(&report)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("report-1", report.ReportID)
	suite.NotEmpty(report.ReportContent)
}

func (suite *PCIDSSAPITestSuite) TestDownloadPCIReport() {
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/reports/report-1/download", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify response headers
	suite.Equal("text/plain", resp.Header.Get("Content-Type"))
	suite.Contains(resp.Header.Get("Content-Disposition"), "attachment")
	suite.Contains(resp.Header.Get("Content-Disposition"), "pci-dss-report-report-1.txt")
}

func (suite *PCIDSSAPITestSuite) TestGetPCIDashboard() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/dashboard", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var dashboard map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	suite.NoError(err)
	
	// Verify response structure
	suite.Contains(dashboard, "overview")
	suite.Contains(dashboard, "by_category")
	suite.Contains(dashboard, "risk_distribution")
	suite.Contains(dashboard, "recent_assessments")
	suite.Contains(dashboard, "upcoming_assessments")
}

func (suite *PCIDSSAPITestSuite) TestGetPCIMetrics() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/metrics", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var metrics map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&metrics)
	suite.NoError(err)
	
	// Verify response structure
	suite.Contains(metrics, "compliance_trend")
	suite.Contains(metrics, "requirement_status_counts")
	suite.Contains(metrics, "evidence_coverage")
	suite.Contains(metrics, "incident_metrics")
	suite.Contains(metrics, "remediation_metrics")
	suite.Contains(metrics, "training_metrics")
	suite.Contains(metrics, "vendor_compliance")
	suite.Contains(metrics, "audit_readiness")
	suite.Contains(metrics, "cost_metrics")
}

func (suite *PCIDSSAPITestSuite) TestMethodNotAllowed() {
	// Test GET on POST endpoint
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/assess", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestInvalidJSON() {
	// Create request with invalid JSON
	req := httptest.NewRequest("PUT", "/api/v1/compliance/pci_dss/requirement/1.1", bytes.NewReader([]byte(`{"invalid": json}`)))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *PCIDSSAPITestSuite) TestPCIRequirementCategories() {
	// Get all requirements
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	
	// Verify all categories are present
	categories := make(map[types.PCIControlCategory]bool)
	for _, req := range requirements {
		categories[req.Category] = true
	}
	
	expectedCategories := []types.PCIControlCategory{
		types.PCIControlCategoryNetworkSecurity,
		types.PCIControlCategorySystemConfiguration,
		types.PCIControlCategoryDataProtection,
		types.PCIControlCategoryMalwareProtection,
		types.PCIControlCategorySecureDevelopment,
		types.PCIControlCategoryAccessControl,
		types.PCIControlCategoryPhysicalSecurity,
		types.PCIControlCategoryMonitoring,
		types.PCIControlCategoryTesting,
		types.PCIControlCategoryPolicyManagement,
	}
	
	for _, category := range expectedCategories {
		suite.True(categories[category], "Missing category: %s", category)
	}
}

func (suite *PCIDSSAPITestSuite) TestPCIRequirementRiskLevels() {
	// Get all requirements
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	
	// Verify all risk levels are present
	riskLevels := make(map[types.RiskLevel]bool)
	for _, req := range requirements {
		riskLevels[req.RiskLevel] = true
	}
	
	expectedRiskLevels := []types.RiskLevel{
		types.RiskLevelCritical,
		types.RiskLevelHigh,
		types.RiskLevelMedium,
		types.RiskLevelLow,
	}
	
	for _, riskLevel := range expectedRiskLevels {
		suite.True(riskLevels[riskLevel], "Missing risk level: %s", riskLevel)
	}
}

func (suite *PCIDSSAPITestSuite) TestPCIRequirementControlTypes() {
	// Get all requirements
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	
	// Verify all control types are present
	controlTypes := make(map[types.PCIControlType]bool)
	for _, req := range requirements {
		controlTypes[req.ControlType] = true
	}
	
	expectedControlTypes := []types.PCIControlType{
		types.PCIControlTypeOrganizational,
		types.PCIControlTypeTechnical,
		types.PCIControlTypeOperational,
		types.PCIControlTypePhysical,
	}
	
	for _, controlType := range expectedControlTypes {
		suite.True(controlTypes[controlType], "Missing control type: %s", controlType)
	}
}

func (suite *PCIDSSAPITestSuite) TestPCIRequirementStatuses() {
	// Get all requirements
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	
	// Verify all statuses are present
	statuses := make(map[types.PCIControlStatus]bool)
	for _, req := range requirements {
		statuses[req.Status] = true
	}
	
	expectedStatuses := []types.PCIControlStatus{
		types.PCIControlStatusNotImplemented,
		types.PCIControlStatusPartiallyImplemented,
		types.PCIControlStatusImplemented,
		types.PCIControlStatusCompliant,
	}
	
	for _, status := range expectedStatuses {
		suite.True(statuses[status], "Missing status: %s", status)
	}
}

func TestPCIDSSAPITestSuite(t *testing.T) {
	suite.Run(t, new(PCIDSSAPITestSuite))
}
