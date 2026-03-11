package types

import (
	"time"
)

// SOC 2 Type II Compliance Types

type SOC2Control struct {
	ID                   string              `json:"id"`
	Title                string              `json:"title"`
	Description           string              `json:"description"`
	Category             SOC2ControlCategory `json:"category"`
	Subcategories        []string            `json:"subcategories"`
	Objective            string              `json:"objective"`
	ControlType          SOC2ControlType    `json:"control_type"`
	Status               SOC2ControlStatus   `json:"status"`
	ImplementationDate   *time.Time          `json:"implementation_date,omitempty"`
	LastReviewDate       *time.Time          `json:"last_review_date,omitempty"`
	Evidence             []string            `json:"evidence"`
	Owner                string              `json:"owner"`
	RiskLevel            RiskLevel           `json:"risk_level"`
}

type SOC2ControlCategory string

const (
	SOC2ControlCategoryGovernance            SOC2ControlCategory = "governance"
	SOC2ControlCategoryAssetManagement       SOC2ControlCategory = "asset_management"
	SOC2ControlCategoryAccessControl         SOC2ControlCategory = "access_control"
	SOC2ControlCategoryOperational           SOC2ControlCategory = "operational"
	SOC2ControlCategoryIncidentResponse      SOC2ControlCategory = "incident_response"
	SOC2ControlCategoryVulnerabilityManagement SOC2ControlCategory = "vulnerability_management"
	SOC2ControlCategoryDisasterRecovery     SOC2ControlCategory = "disaster_recovery"
	SOC2ControlCategoryTestEvaluation        SOC2ControlCategory = "test_evaluation"
	SOC2ControlCategoryCommunicationsSecurity SOC2ControlCategory = "communications_security"
)

type SOC2ControlType string

const (
	SOC2ControlTypeOrganizational SOC2ControlType = "organizational"
	SOC2ControlTypeTechnical      SOC2ControlType = "technical"
	SOC2ControlTypeOperational    SOC2ControlType = "operational"
)

type SOC2ControlStatus string

const (
	SOC2ControlStatusNotImplemented        SOC2ControlStatus = "not_implemented"
	SOC2ControlStatusPartiallyImplemented   SOC2ControlStatus = "partially_implemented"
	SOC2ControlStatusImplemented           SOC2ControlStatus = "implemented"
	SOC2ControlStatusCompliant             SOC2ControlStatus = "compliant"
)

type SOC2ControlUpdate struct {
	UpdateType UpdateType `json:"update_type"`
	UpdatedBy  string    `json:"updated_by"`
	Timestamp  time.Time `json:"timestamp"`
	Notes      string    `json:"notes"`
	Evidence   []string  `json:"evidence"`
}

type SOC2TestType string

const (
	SOC2TestTypePenetrationTest        SOC2TestType = "penetration_test"
	SOC2TestTypeInternalTest           SOC2TestType = "internal_test"
	SOC2TestTypeRedTeamTest            SOC2TestType = "red_team_test"
	SOC2TestTypeExternalTest           SOC2TestType = "external_test"
	SOC2TestTypeVulnerabilityValidation SOC2TestType = "vulnerability_validation"
)

type SOC2Finding struct {
	Severity      FindingSeverity `json:"severity"`
	Description   string          `json:"description"`
	Recommendation string         `json:"recommendation"`
	EvidenceGaps  []string        `json:"evidence_gaps"`
}

type SOC2Recommendation struct {
	Priority       RecommendationPriority `json:"priority"`
	Title          string                  `json:"title"`
	Description     string                  `json:"description"`
	Findings       []SOC2Finding           `json:"findings"`
	EstimatedEffort string                 `json:"estimated_effort"`
	Owner          string                  `json:"owner"`
}

type SOC2Assessment struct {
	AssessmentID        string                `json:"assessment_id"`
	Timestamp           time.Time             `json:"timestamp"`
	Framework           string                `json:"framework"`
	Version             string                `json:"version"`
	Scope               SOC2Scope             `json:"scope"`
	OverallScore        float64               `json:"overall_score"`
	ControlAssessments  []SOC2ControlAssessment `json:"control_assessments"`
	Findings            []SOC2Finding         `json:"findings"`
	Recommendations     []SOC2Recommendation  `json:"recommendations"`
	LastAssessed        time.Time             `json:"last_assessed"`
}

type SOC2ControlAssessment struct {
	ControlID        string          `json:"control_id"`
	ControlTitle     string          `json:"control_title"`
	Category         SOC2ControlCategory `json:"category"`
	ComplianceScore  float64         `json:"compliance_score"`
	Status           SOC2ControlStatus `json:"status"`
	Findings         []SOC2Finding    `json:"findings"`
	Recommendations  []string        `json:"recommendations"`
	LastAssessed     time.Time       `json:"last_assessed"`
}

type SOC2Scope struct {
	Departments []string `json:"departments"`
	Systems     []string `json:"systems"`
	Processes   []string `json:"processes"`
}

type SOC2Stats struct {
	TotalControls                int     `json:"total_controls"`
	CompliantControls           int     `json:"compliant_controls"`
	ImplementedControls         int     `json:"implemented_controls"`
	PartiallyImplementedControls int     `json:"partially_implemented_controls"`
	NotImplementedControls      int     `json:"not_implemented_controls"`
	AverageComplianceScore      float64 `json:"average_compliance_score"`
	TotalIncidents              uint64  `json:"total_incidents"`
	AverageDetectionTimeMinutes float64 `json:"average_detection_time_minutes"`
	AverageResponseTimeMinutes  float64 `json:"average_response_time_minutes"`
	LastIncidentDate            *time.Time `json:"last_incident_date,omitempty"`
}

type SOC2Database struct {
	Assessments []SOC2Assessment          `json:"assessments"`
	Reports     map[string]SOC2Report    `json:"reports"`
}

type SOC2Report struct {
	ReportID      string       `json:"report_id"`
	GeneratedAt   time.Time    `json:"generated_at"`
	Assessment    SOC2Assessment `json:"assessment"`
	ReportContent string       `json:"report_content"`
	Format        string       `json:"format"`
}

// Framework Components

type SOC2Framework struct {
	Controls map[string]SOC2Control `json:"controls"`
}

type SOC2Config struct {
	FrameworkConfig  FrameworkConfig  `json:"framework_config"`
	SecurityConfig   SecurityConfig   `json:"security_config"`
	AuditConfig      AuditConfig      `json:"audit_config"`
	MonitoringConfig MonitoringConfig `json:"monitoring_config"`
	MetricsConfig    MetricsConfig    `json:"metrics_config"`
}

type SecurityEngine struct {
	// Implementation details
}

type AuditManager struct {
	// Implementation details
}

type MetricsCollector struct {
	// Implementation details
}

// Error types

type SOC2Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e SOC2Error) Error() string {
	return e.Message
}
