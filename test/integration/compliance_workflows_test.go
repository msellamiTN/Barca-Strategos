package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"barca-strategos/internal/compliance/soc2"
	"barca-strategos/internal/compliance/pci_dss"
	"barca-strategos/internal/routes/api/v1"
	"barca-strategos/pkg/types"
	"barca-strategos/internal/websocket"
)

type ComplianceWorkflowTestSuite struct {
	suite.Suite
	app        *fiber.App
	soc2Svc    *soc2.Service
	pciSvc     *pci_dss.Service
	hub        *websocket.Hub
	broadcast  func([]byte)
}

func (suite *ComplianceWorkflowTestSuite) SetupSuite() {
	suite.app = fiber.New()
	suite.hub = websocket.NewHub()
	suite.soc2Svc = soc2.New()
	suite.pciSvc = pci_dss.New()
	suite.broadcast = func(data []byte) {
		// Mock broadcast function
	}
	
	suite.soc2Svc.SetBroadcast(suite.broadcast)
	suite.pciSvc.SetBroadcast(suite.broadcast)
	
	// Initialize services
	suite.soc2Svc.Initialize(suite.app.Context())
	suite.pciSvc.Initialize(suite.app.Context())
	
	// Register routes
	v1.RegisterSOC2Routes(suite.app, suite.soc2Svc, suite.broadcast)
	v1.RegisterPCIDSSRoutes(suite.app, suite.pciSvc, suite.broadcast)
	v1.RegisterWebSocketRoutes(suite.app, suite.hub)
}

func (suite *ComplianceWorkflowTestSuite) TestSOC2CompleteWorkflow() {
	// Step 1: Get initial SOC2 stats
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/stats", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var initialStats types.SOC2Stats
	err = json.NewDecoder(resp.Body).Decode(&initialStats)
	suite.NoError(err)
	
	// Step 2: Get all SOC2 controls
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/controls", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var controls []types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&controls)
	suite.NoError(err)
	suite.Greater(len(controls), 0)
	
	// Step 3: Select a control and update it
	control := controls[0]
	update := types.SOC2ControlUpdate{
		UpdateType: types.UpdateTypeStatus,
		UpdatedBy:  "test-user",
		Timestamp:  time.Now().UTC(),
		Notes:      "compliant",
		Evidence:   []string{"Updated evidence document", "Test compliance proof"},
	}
	
	updateJSON, _ := json.Marshal(update)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/compliance/soc2/control/%s", control.ID), bytes.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Step 4: Verify control status was updated
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/compliance/soc2/control/%s/status", control.ID), nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var status types.SOC2ControlStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	suite.NoError(err)
	suite.Equal(types.SOC2ControlStatusCompliant, status)
	
	// Step 5: Conduct SOC2 assessment
	scope := types.SOC2Scope{
		Departments: []string{"IT", "Security"},
		Systems:     []string{"Phoenix Core", "Database"},
		Processes:   []string{"Incident Response"},
	}
	
	scopeJSON, _ := json.Marshal(scope)
	req = httptest.NewRequest("POST", "/api/v1/compliance/soc2/assess", bytes.NewReader(scopeJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var assessment types.SOC2Assessment
	err = json.NewDecoder(resp.Body).Decode(&assessment)
	suite.NoError(err)
	suite.Equal("SOC 2", assessment.Framework)
	suite.Equal("2017", assessment.Version)
	suite.Greater(assessment.OverallScore, 0.0)
	suite.NotEmpty(assessment.AssessmentID)
	
	// Step 6: Generate SOC2 report
	assessmentJSON, _ := json.Marshal(assessment)
	req = httptest.NewRequest("POST", "/api/v1/compliance/soc2/report", bytes.NewReader(assessmentJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var report types.SOC2Report
	err = json.NewDecoder(resp.Body).Decode(&report)
	suite.NoError(err)
	suite.Equal(assessment.AssessmentID, report.ReportID)
	suite.NotEmpty(report.ReportContent)
	
	// Step 7: Verify assessment history
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/assessments", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var assessments []types.SOC2Assessment
	err = json.NewDecoder(resp.Body).Decode(&assessments)
	suite.NoError(err)
	suite.Greater(len(assessments), 0)
	
	// Step 8: Verify updated stats
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/stats", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var updatedStats types.SOC2Stats
	err = json.NewDecoder(resp.Body).Decode(&updatedStats)
	suite.NoError(err)
	
	// Stats should reflect changes
	suite.GreaterOrEqual(updatedStats.AverageComplianceScore, initialStats.AverageComplianceScore)
}

func (suite *ComplianceWorkflowTestSuite) TestPCICompleteWorkflow() {
	// Step 1: Get initial PCI DSS stats
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/stats", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var initialStats types.PCIStats
	err = json.NewDecoder(resp.Body).Decode(&initialStats)
	suite.NoError(err)
	
	// Step 2: Get all PCI DSS requirements
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	suite.Greater(len(requirements), 0)
	
	// Step 3: Select a requirement and update it
	requirement := requirements[0]
	update := types.PCIControlUpdate{
		UpdateType: types.UpdateTypeStatus,
		UpdatedBy:  "test-user",
		Timestamp:  time.Now().UTC(),
		Notes:      "compliant",
		Evidence:   []string{"Updated evidence document", "Test compliance proof"},
	}
	
	updateJSON, _ := json.Marshal(update)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/compliance/pci_dss/requirement/%s", requirement.ID), bytes.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Step 4: Verify requirement status was updated
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/compliance/pci_dss/requirement/%s/status", requirement.ID), nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var status types.PCIControlStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	suite.NoError(err)
	suite.Equal(types.PCIControlStatusCompliant, status)
	
	// Step 5: Conduct PCI DSS assessment
	scope := types.PCIScope{
		Departments: []string{"Payment Processing", "Security"},
		Systems:     []string{"Payment Gateway", "Database"},
		Processes:   []string{"Card Processing"},
	}
	
	scopeJSON, _ := json.Marshal(scope)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/assess", bytes.NewReader(scopeJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var assessment types.PCIAssessment
	err = json.NewDecoder(resp.Body).Decode(&assessment)
	suite.NoError(err)
	suite.Equal("PCI DSS", assessment.Framework)
	suite.Equal("4.0", assessment.Version)
	suite.Greater(assessment.OverallScore, 0.0)
	suite.NotEmpty(assessment.AssessmentID)
	suite.NotZero(assessment.NextAssessmentDate)
	
	// Step 6: Generate PCI DSS report
	assessmentJSON, _ := json.Marshal(assessment)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/report", bytes.NewReader(assessmentJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var report types.PCIReport
	err = json.NewDecoder(resp.Body).Decode(&report)
	suite.NoError(err)
	suite.Equal(assessment.AssessmentID, report.ReportID)
	suite.NotEmpty(report.ReportContent)
	
	// Step 7: Verify assessment history
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/assessments", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var assessments []types.PCIAssessment
	err = json.NewDecoder(resp.Body).Decode(&assessments)
	suite.NoError(err)
	suite.Greater(len(assessments), 0)
	
	// Step 8: Verify updated stats
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/stats", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var updatedStats types.PCIStats
	err = json.NewDecoder(resp.Body).Decode(&updatedStats)
	suite.NoError(err)
	
	// Stats should reflect changes
	suite.GreaterOrEqual(updatedStats.AverageComplianceScore, initialStats.AverageComplianceScore)
	
	// Step 9: Verify dashboard data
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/dashboard", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var dashboard map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	suite.NoError(err)
	suite.Contains(dashboard, "overview")
	suite.Contains(dashboard, "by_category")
	suite.Contains(dashboard, "risk_distribution")
	
	// Step 10: Verify metrics data
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/metrics", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var metrics map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&metrics)
	suite.NoError(err)
	suite.Contains(metrics, "compliance_trend")
	suite.Contains(metrics, "requirement_status_counts")
	suite.Contains(metrics, "evidence_coverage")
}

func (suite *ComplianceWorkflowTestSuite) TestSOC2BulkOperationsWorkflow() {
	// Step 1: Get all controls
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/controls", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var controls []types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&controls)
	suite.NoError(err)
	suite.Greater(len(controls), 0)
	
	// Step 2: Perform bulk assessment on multiple controls
	controlIDs := make([]string, 0, 3)
	for i := 0; i < 3 && i < len(controls); i++ {
		controlIDs = append(controlIDs, controls[i].ID)
	}
	
	bulk := struct {
		Action     string   `json:"action"`
		ControlIDs []string `json:"control_ids"`
		Parameters map[string]interface{} `json:"parameters"`
	}{
		Action:     "assess",
		ControlIDs: controlIDs,
		Parameters: map[string]interface{}{},
	}
	
	bulkJSON, _ := json.Marshal(bulk)
	req = httptest.NewRequest("POST", "/api/v1/compliance/soc2/bulk", bytes.NewReader(bulkJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Contains(response, "results")
	
	// Step 3: Verify bulk operation results
	results, ok := response["results"].([]interface{})
	suite.True(ok)
	suite.Len(results, len(controlIDs))
	
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		suite.True(ok)
		suite.Contains(resultMap, "control_id")
		suite.Contains(resultMap, "success")
		suite.Contains(resultMap, "message")
	}
	
	// Step 4: Perform bulk update operation
	bulkUpdate := struct {
		Action     string   `json:"action"`
		ControlIDs []string `json:"control_ids"`
		Parameters map[string]interface{} `json:"parameters"`
	}{
		Action:     "update",
		ControlIDs: controlIDs,
		Parameters: map[string]interface{}{
			"status":  "compliant",
			"evidence": []string{"Bulk evidence update"},
		},
	}
	
	bulkUpdateJSON, _ := json.Marshal(bulkUpdate)
	req = httptest.NewRequest("POST", "/api/v1/compliance/soc2/bulk", bytes.NewReader(bulkUpdateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	results, ok = response["results"].([]interface{})
	suite.True(ok)
	suite.Len(results, len(controlIDs))
}

func (suite *ComplianceWorkflowTestSuite) TestPCIBulkOperationsWorkflow() {
	// Step 1: Get all requirements
	req := httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirements", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var requirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&requirements)
	suite.NoError(err)
	suite.Greater(len(requirements), 0)
	
	// Step 2: Perform bulk assessment on multiple requirements
	requirementIDs := make([]string, 0, 3)
	for i := 0; i < 3 && i < len(requirements); i++ {
		requirementIDs = append(requirementIDs, requirements[i].ID)
	}
	
	bulk := struct {
		Action         string   `json:"action"`
		RequirementIDs []string `json:"requirement_ids"`
		Parameters    map[string]interface{} `json:"parameters"`
	}{
		Action:         "assess",
		RequirementIDs: requirementIDs,
		Parameters:     map[string]interface{}{},
	}
	
	bulkJSON, _ := json.Marshal(bulk)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/bulk", bytes.NewReader(bulkJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Contains(response, "results")
	
	// Step 3: Verify bulk operation results
	results, ok := response["results"].([]interface{})
	suite.True(ok)
	suite.Len(results, len(requirementIDs))
	
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		suite.True(ok)
		suite.Contains(resultMap, "requirement_id")
		suite.Contains(resultMap, "success")
		suite.Contains(resultMap, "message")
	}
	
	// Step 4: Perform bulk update operation
	bulkUpdate := struct {
		Action         string   `json:"action"`
		RequirementIDs []string `json:"requirement_ids"`
		Parameters    map[string]interface{} `json:"parameters"`
	}{
		Action:         "update",
		RequirementIDs: requirementIDs,
		Parameters:     map[string]interface{}{
			"status":  "compliant",
			"evidence": []string{"Bulk evidence update"},
		},
	}
	
	bulkUpdateJSON, _ := json.Marshal(bulkUpdate)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/bulk", bytes.NewReader(bulkUpdateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	results, ok = response["results"].([]interface{})
	suite.True(ok)
	suite.Len(results, len(requirementIDs))
}

func (suite *ComplianceWorkflowTestSuite) TestEvidenceManagementWorkflow() {
	// Test SOC2 evidence management
	req := httptest.NewRequest("GET", "/api/v1/compliance/soc2/control/CC1.1/evidence", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var soc2Evidence []string
	err = json.NewDecoder(resp.Body).Decode(&soc2Evidence)
	suite.NoError(err)
	
	// Upload new evidence
	upload := struct {
		Evidence []string `json:"evidence"`
	}{
		Evidence: append(soc2Evidence, "New SOC2 evidence", "Additional SOC2 proof"),
	}
	
	uploadJSON, _ := json.Marshal(upload)
	req = httptest.NewRequest("POST", "/api/v1/compliance/soc2/control/CC1.1/evidence", bytes.NewReader(uploadJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify evidence was updated
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/control/CC1.1/evidence", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var updatedSOCEvidence []string
	err = json.NewDecoder(resp.Body).Decode(&updatedSOCEvidence)
	suite.NoError(err)
	suite.Greater(len(updatedSOCEvidence), len(soc2Evidence))
	
	// Test PCI DSS evidence management
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirement/1.1/evidence", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var pciEvidence []string
	err = json.NewDecoder(resp.Body).Decode(&pciEvidence)
	suite.NoError(err)
	
	// Upload new evidence
	pciUpload := struct {
		Evidence []string `json:"evidence"`
	}{
		Evidence: append(pciEvidence, "New PCI evidence", "Additional PCI proof"),
	}
	
	pciUploadJSON, _ := json.Marshal(pciUpload)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/requirement/1.1/evidence", bytes.NewReader(pciUploadJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify evidence was updated
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirement/1.1/evidence", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var updatedPCIEvidence []string
	err = json.NewDecoder(resp.Body).Decode(&updatedPCIEvidence)
	suite.NoError(err)
	suite.Greater(len(updatedPCIEvidence), len(pciEvidence))
}

func (suite *ComplianceWorkflowTestSuite) TestReportingWorkflow() {
	// Step 1: Conduct SOC2 assessment
	scope := types.SOC2Scope{
		Departments: []string{"IT", "Security"},
		Systems:     []string{"Phoenix Core", "Database"},
		Processes:   []string{"Incident Response"},
	}
	
	scopeJSON, _ := json.Marshal(scope)
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/assess", bytes.NewReader(scopeJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var soc2Assessment types.SOC2Assessment
	err = json.NewDecoder(resp.Body).Decode(&soc2Assessment)
	suite.NoError(err)
	
	// Step 2: Generate SOC2 report
	soc2AssessmentJSON, _ := json.Marshal(soc2Assessment)
	req = httptest.NewRequest("POST", "/api/v1/compliance/soc2/report", bytes.NewReader(soc2AssessmentJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var soc2Report types.SOC2Report
	err = json.NewDecoder(resp.Body).Decode(&soc2Report)
	suite.NoError(err)
	
	// Step 3: Download SOC2 report
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/v1/compliance/soc2/reports/%s/download", soc2Report.ReportID), nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("text/plain", resp.Header.Get("Content-Type"))
	suite.Contains(resp.Header.Get("Content-Disposition"), "attachment")
	
	// Step 4: Conduct PCI DSS assessment
	pciScope := types.PCIScope{
		Departments: []string{"Payment Processing", "Security"},
		Systems:     []string{"Payment Gateway", "Database"},
		Processes:   []string{"Card Processing"},
	}
	
	pciScopeJSON, _ := json.Marshal(pciScope)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/assess", bytes.NewReader(pciScopeJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var pciAssessment types.PCIAssessment
	err = json.NewDecoder(resp.Body).Decode(&pciAssessment)
	suite.NoError(err)
	
	// Step 5: Generate PCI DSS report
	pciAssessmentJSON, _ := json.Marshal(pciAssessment)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/report", bytes.NewReader(pciAssessmentJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var pciReport types.PCIReport
	err = json.NewDecoder(resp.Body).Decode(&pciReport)
	suite.NoError(err)
	
	// Step 6: Download PCI DSS report
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/v1/compliance/pci_dss/reports/%s/download", pciReport.ReportID), nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("text/plain", resp.Header.Get("Content-Type"))
	suite.Contains(resp.Header.Get("Content-Disposition"), "attachment")
	
	// Step 7: Verify reports are stored
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/reports", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var soc2Reports map[string]types.SOC2Report
	err = json.NewDecoder(resp.Body).Decode(&soc2Reports)
	suite.NoError(err)
	suite.Contains(soc2Reports, soc2Report.ReportID)
	
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/reports", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var pciReports map[string]types.PCIReport
	err = json.NewDecoder(resp.Body).Decode(&pciReports)
	suite.NoError(err)
	suite.Contains(pciReports, pciReport.ReportID)
}

func (suite *ComplianceWorkflowTestSuite) TestSearchAndFilterWorkflow() {
	// Test SOC2 search
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
	
	searchJSON, _ := json.Marshal(search)
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/search", bytes.NewReader(searchJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var soc2Controls []types.SOC2Control
	err = json.NewDecoder(resp.Body).Decode(&soc2Controls)
	suite.NoError(err)
	suite.Greater(len(soc2Controls), 0)
	
	// Verify search results
	for _, control := range soc2Controls {
		suite.Contains(control.Category, types.SOC2ControlCategoryGovernance)
	}
	
	// Test PCI DSS search
	pciSearch := struct {
		Query        string   `json:"query"`
		Categories  []string `json:"categories"`
		Status      []string `json:"status"`
		Types       []string `json:"types"`
	}{
		Query:       "network",
		Categories:  []string{"network_security"},
		Status:      []string{"implemented"},
		Types:       []string{"technical"},
	}
	
	pciSearchJSON, _ := json.Marshal(pciSearch)
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/search", bytes.NewReader(pciSearchJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	var pciRequirements []types.PCIRequirement
	err = json.NewDecoder(resp.Body).Decode(&pciRequirements)
	suite.NoError(err)
	suite.Greater(len(pciRequirements), 0)
	
	// Verify search results
	for _, requirement := range pciRequirements {
		suite.Contains(requirement.Category, types.PCIControlCategoryNetworkSecurity)
	}
}

func (suite *ComplianceWorkflowTestSuite) TestErrorHandlingWorkflow() {
	// Test invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/compliance/soc2/assess", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	
	req = httptest.NewRequest("POST", "/api/v1/compliance/pci_dss/assess", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	
	// Test non-existent resources
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/control/NOT_FOUND/status", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
	
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/requirement/NOT_FOUND/status", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
	
	// Test method not allowed
	req = httptest.NewRequest("GET", "/api/v1/compliance/soc2/assess", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
	
	req = httptest.NewRequest("GET", "/api/v1/compliance/pci_dss/assess", nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestComplianceWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(ComplianceWorkflowTestSuite))
}
