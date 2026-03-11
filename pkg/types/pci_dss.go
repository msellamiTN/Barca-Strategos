package types

import (
	"time"
)

// PCI DSS Compliance Types

type PCIRequirement struct {
	ID                   string                `json:"id"`
	Title                string                `json:"title"`
	Description           string                `json:"description"`
	Category             PCIControlCategory    `json:"category"`
	Subcategories        []string              `json:"subcategories"`
	Objective            string                `json:"objective"`
	ControlType          PCIControlType        `json:"control_type"`
	Status               PCIControlStatus      `json:"status"`
	ImplementationDate   *time.Time            `json:"implementation_date,omitempty"`
	LastReviewDate       *time.Time            `json:"last_review_date,omitempty"`
	Evidence             []string              `json:"evidence"`
	Owner                string                `json:"owner"`
	RiskLevel            RiskLevel             `json:"risk_level"`
}

type PCIControlCategory string

const (
	PCIControlCategoryNetworkSecurity      PCIControlCategory = "network_security"
	PCIControlCategorySystemConfiguration  PCIControlCategory = "system_configuration"
	PCIControlCategoryDataProtection       PCIControlCategory = "data_protection"
	PCIControlCategoryMalwareProtection   PCIControlCategory = "malware_protection"
	PCIControlCategorySecureDevelopment   PCIControlCategory = "secure_development"
	PCIControlCategoryAccessControl        PCIControlCategory = "access_control"
	PCIControlCategoryPhysicalSecurity     PCIControlCategory = "physical_security"
	PCIControlCategoryMonitoring           PCIControlCategory = "monitoring"
	PCIControlCategoryTesting              PCIControlCategory = "testing"
	PCIControlCategoryPolicyManagement     PCIControlCategory = "policy_management"
)

type PCIControlType string

const (
	PCIControlTypeOrganizational PCIControlType = "organizational"
	PCIControlTypeTechnical      PCIControlType = "technical"
	PCIControlTypeOperational    PCIControlType = "operational"
	PCIControlTypePhysical       PCIControlType = "physical"
)

type PCIControlStatus string

const (
	PCIControlStatusNotImplemented        PCIControlStatus = "not_implemented"
	PCIControlStatusPartiallyImplemented   PCIControlStatus = "partially_implemented"
	PCIControlStatusImplemented           PCIControlStatus = "implemented"
	PCIControlStatusCompliant             PCIControlStatus = "compliant"
)

type PCIControlUpdate struct {
	UpdateType UpdateType `json:"update_type"`
	UpdatedBy  string    `json:"updated_by"`
	Timestamp  time.Time `json:"timestamp"`
	Notes      string    `json:"notes"`
	Evidence   []string  `json:"evidence"`
}

type PCIFinding struct {
	Severity       FindingSeverity `json:"severity"`
	Description    string          `json:"description"`
	Recommendation string          `json:"recommendation"`
	EvidenceGaps   []string        `json:"evidence_gaps"`
}

type PCIRecommendation struct {
	Priority        RecommendationPriority `json:"priority"`
	Title           string                  `json:"title"`
	Description      string                  `json:"description"`
	Findings        []PCIFinding            `json:"findings"`
	EstimatedEffort string                  `json:"estimated_effort"`
	Owner           string                  `json:"owner"`
}

type PCIAssessment struct {
	AssessmentID          string                     `json:"assessment_id"`
	Timestamp             time.Time                  `json:"timestamp"`
	Framework             string                     `json:"framework"`
	Version               string                     `json:"version"`
	Scope                 PCIScope                   `json:"scope"`
	OverallScore          float64                    `json:"overall_score"`
	RequirementAssessments []PCIRequirementAssessment `json:"requirement_assessments"`
	Findings              []PCIFinding               `json:"findings"`
	Recommendations       []PCIRecommendation        `json:"recommendations"`
	NextAssessmentDate    time.Time                  `json:"next_assessment_date"`
}

type PCIRequirementAssessment struct {
	RequirementID       string            `json:"requirement_id"`
	RequirementTitle    string            `json:"requirement_title"`
	Category            PCIControlCategory `json:"category"`
	ComplianceScore     float64           `json:"compliance_score"`
	Status              PCIControlStatus  `json:"status"`
	Findings            []PCIFinding       `json:"findings"`
	Recommendations     []string          `json:"recommendations"`
	LastAssessed        time.Time         `json:"last_assessed"`
}

type PCIScope struct {
	Departments []string `json:"departments"`
	Systems     []string `json:"systems"`
	Processes   []string `json:"processes"`
}

type PCIStats struct {
	TotalRequirements                int     `json:"total_requirements"`
	CompliantRequirements           int     `json:"compliant_requirements"`
	ImplementedRequirements         int     `json:"implemented_requirements"`
	PartiallyImplementedRequirements int     `json:"partially_implemented_requirements"`
	NotImplementedRequirements      int     `json:"not_implemented_requirements"`
	AverageComplianceScore          float64 `json:"average_compliance_score"`
	TotalIncidents                  uint64  `json:"total_incidents"`
	AverageDetectionTimeMinutes     float64 `json:"average_detection_time_minutes"`
	AverageResponseTimeMinutes      float64 `json:"average_response_time_minutes"`
	LastIncidentDate                *time.Time `json:"last_incident_date,omitempty"`
}

type PCIDatabase struct {
	Assessments  []PCIAssessment     `json:"assessments"`
	Reports      map[string]PCIReport `json:"reports"`
	MetricsStore  map[string]interface{} `json:"metrics_store"`
}

type PCIReport struct {
	ReportID      string       `json:"report_id"`
	GeneratedAt   time.Time    `json:"generated_at"`
	Assessment    PCIAssessment `json:"assessment"`
	ReportContent string       `json:"report_content"`
	Format        string       `json:"format"`
}

// Framework Components

type PCIFramework struct {
	Requirements map[string]PCIRequirement `json:"requirements"`
}

type PCIConfig struct {
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

type PCIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e PCIError) Error() string {
	return e.Message
}
