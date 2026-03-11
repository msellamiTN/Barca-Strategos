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
	"barca-strategos/internal/compliance/soc2"
	"barca-strategos/pkg/types"
)

type SOC2APITestSuite struct {
	suite.Suite
	app      *fiber.App
	soc2Svc  *soc2.Service
	broadcast func([]byte)
}

func (suite *SOC2APITestSuite) SetupSuite() {
	suite.app = fiber.New()
	suite.soc2Svc = soc2.New()
	suite.broadcast = func(data []byte) {
		// Mock broadcast function
	}
	suite.soc2Svc.SetBroadcast(suite.broadcast)
	
	// Register routes
	RegisterSOC2Routes(suite.app, suite.soc2Svc, suite.broadcast)
}

func (suite *SOC2APITestSuite) TestAssessSOC2Compliance() {
	// Create test scope
	scope := types.SOC2Scope{
		Departments: []string{"IT", "Security"},
		Systems:     []string{"Phoenix Core", "Database"},
		Processes:   []string{"Incident Response"},
	}
	
	// Convert to JSON
	scopeJSON, _ := json.Marshal(scope)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/assess", bytes.NewReader(scopeJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var assessment types.SOC2Assessment
	err = json.NewDecoder(resp.Body).Decode(&assessment)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("SOC 2", assessment.Framework)
	suite.Equal("2017", assessment.Version)
	suite.Greater(assessment.OverallScore, 0.0)
	suite.NotEmpty(assessment.AssessmentID)
}

func (suite *SOC2APITestSuite) TestAssessSOC2ComplianceInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/assess", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestGetSOC2Controls() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/controls", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var controls []types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&controls)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(controls), 0)
}

func (suite *SOC2APITestSuite) TestGetSOC2Control() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/controls/CC1.1", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var control types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&control)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("CC1.1", control.ID)
	suite.NotEmpty(control.Title)
}

func (suite *SOC2APITestSuite) TestGetSOC2ControlNotFound() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/controls/NOT_FOUND", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var control types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&control)
	suite.NoError(err)
	
	// Verify response (should return mock data)
	suite.Equal("NOT_FOUND", control.ID)
}

func (suite *SOC2APITestSuite) TestUpdateSOC2Control() {
	// Create update payload
	update := types.SOC2ControlUpdate{
		UpdateType: types.UpdateTypeStatus,
		UpdatedBy:  "test-user",
		Timestamp:  suite.soc2Svc.(*soc2.Service).(*soc2.Service).controls["CC1.1"].LastReviewDate,
		Notes:      "compliant",
		Evidence:   []string{"New evidence document"},
	}
	
	// Convert to JSON
	updateJSON, _ := json.Marshal(update)
	
	// Create request
	req := httptest.NewRequest("PUT", "/api/v1/compliance/soc2/control/CC1.1", bytes.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestUpdateSOC2ControlInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("PUT", "/api/v1/compliance/soc2/control/CC1.1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestGetSOC2ControlStatus() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/control/CC1.1/status", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var status types.SOC2ControlStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	suite.NoError(err)
	
	// Verify response
	suite.NotEmpty(status)
}

func (suite *SOC2APITestSuite) TestGetSOC2ControlStatusNotFound() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/control/NOT_FOUND/status", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestGenerateSOC2Report() {
	// Create assessment payload
	assessment := types.SOC2Assessment{
		AssessmentID:        "test-assessment",
		Timestamp:           suite.soc2Svc.(*soc2.Service).(*soc2.Service).controls["CC1.1"].LastReviewDate,
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
		LastAssessed:        suite.soc2Svc.(*soc2.Service).(*soc2.Service).controls["CC1.1"].LastReviewDate,
	}
	
	// Convert to JSON
	assessmentJSON, _ := json.Marshal(assessment)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/report", bytes.NewReader(assessmentJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var report types.SOC2Report
	err = json.NewDecoder(resp.Body).Decode(&report)
	suite.NoError(err)
	
	// Verify response
	suite.Equal(assessment.AssessmentID, report.ReportID)
	suite.NotEmpty(report.ReportContent)
}

func (suite *SOC2APITestSuite) TestGenerateSOC2ReportInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/report", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestGetSOC2Stats() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/stats", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var stats types.SOC2Stats
	err = json.NewDecoder(resp.Body).Decode(&stats)
	suite.NoError(err)
	
	// Verify response
	suite.GreaterOrEqual(stats.TotalControls, 0)
	suite.GreaterOrEqual(stats.AverageComplianceScore, 0.0)
}

func (suite *SOC2APITestSuite) TestGetSOC2Assessments() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/assessments", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var assessments []types.SOC2Assessment
	err = json.NewDecoder(resp.Body).Decode(&assessments)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(assessments), 0)
}

func (suite *SOC2APITestSuite) TestGetSOC2Assessment() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/assessments/assessment-1", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var assessment types.SOC2Assessment
	err = json.NewDecoder(resp.Body).Decode(&assessment)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("assessment-1", assessment.AssessmentID)
	suite.Equal("SOC 2", assessment.Framework)
}

func (suite *SOC2APITestSuite) TestSearchSOC2Controls() {
	// Create search request
	search := struct {
		Query      string   `json:"query"`
		Categories []string `json:"categories"`
		Status     []string `json:"status"`
		Types      []string `json:"types"`
	}{
		Query:      "governance",
		Categories: []string{"governance"},
		Status:     []string{"implemented"},
		Types:      []string{"organizational"},
	}
	
	// Convert to JSON
	searchJSON, _ := json.Marshal(search)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/search", bytes.NewReader(searchJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var controls []types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&controls)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(controls), 0)
}

func (suite *SOC2APITestSuite) TestSearchSOC2ControlsInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/search", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestBulkSOC2Operations() {
	// Create bulk operation request
	bulk := struct {
		Action     string   `json:"action"`
		ControlIDs []string `json:"control_ids"`
		Parameters map[string]interface{} `json:"parameters"`
	}{
		Action:     "assess",
		ControlIDs: []string{"CC1.1", "CC2.1"},
		Parameters: map[string]interface{}{},
	}
	
	// Convert to JSON
	bulkJSON, _ := json.Marshal(bulk)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/bulk", bytes.NewReader(bulkJSON))
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

func (suite *SOC2APITestSuite) TestBulkSOC2OperationsInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/bulk", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestGetSOC2ControlEvidence() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/control/CC1.1/evidence", nil)
	
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

func (suite *SOC2APITestSuite) TestUploadSOC2ControlEvidence() {
	// Create evidence upload request
	upload := struct {
		Evidence []string `json:"evidence"`
	}{
		Evidence: []string{"New evidence document 1", "New evidence document 2"},
	}
	
	// Convert to JSON
	uploadJSON, _ := json.Marshal(upload)
	
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/control/CC1.1/evidence", bytes.NewReader(uploadJSON))
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

func (suite *SOC2APITestSuite) TestUploadSOC2ControlEvidenceInvalidRequest() {
	// Create invalid request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/control/CC1.1/evidence", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestGetSOC2Reports() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/reports", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var reports map[string]types.SOC2Report
	err = json.NewDecoder(resp.Body).Decode(&reports)
	suite.NoError(err)
	
	// Verify response
	suite.Greater(len(reports), 0)
}

func (suite *SOC2APITestSuite) TestGetSOC2Report() {
	// Create request
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/reports/report-1", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Parse response
	var report types.SOC2Report
	err = json.NewDecoder(resp.Body).Decode(&report)
	suite.NoError(err)
	
	// Verify response
	suite.Equal("report-1", report.ReportID)
	suite.NotEmpty(report.ReportContent)
}

func (suite *SOC2APITestSuite) TestDownloadSOC2Report() {
	// Create request
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/reports/report-1/download", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify response headers
	suite.Equal("text/plain", resp.Header.Get("Content-Type"))
	suite.Contains(resp.Header.Get("Content-Disposition"), "attachment")
	suite.Contains(resp.Header.Get("Content-Disposition"), "soc2-report-report-1.txt")
}

func (suite *SOC2APITestSuite) TestMethodNotAllowed() {
	// Test GET on POST endpoint
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/assess", nil)
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *SOC2APITestSuite) TestInvalidJSON() {
	// Create request with invalid JSON
	req := httptest.NewRequest("PUT", "/api/v1/compliance/soc2/control/CC1.1", bytes.NewReader([]byte(`{"invalid": json}`)))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestSOC2APITestSuite(t *testing.T) {
	suite.Run(t, new(SOC2APITestSuite))
}
