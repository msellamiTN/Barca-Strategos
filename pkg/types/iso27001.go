package types

import "time"

// ISO27001 represents the ISO 27001 Information Security Management System
type ISO27001 struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Version         string              `json:"version"`
	Status          ISO27001Status      `json:"status"`
	Controls        []Control           `json:"controls"`
	Clauses         []Clause            `json:"clauses"`
	RiskAssessment  RiskAssessment      `json:"risk_assessment"`
	AuditTrail      []AuditEntry        `json:"audit_trail"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// ISO27001Status represents the implementation status
type ISO27001Status string

const (
	ISO27001StatusDraft        ISO27001Status = "draft"
	ISO27001StatusActive       ISO27001Status = "active"
	ISO27001StatusUnderReview  ISO27001Status = "under_review"
	ISO27001StatusApproved     ISO27001Status = "approved"
	ISO27001StatusNonCompliant ISO27001Status = "non_compliant"
)

// Control represents an ISO 27001 control
type Control struct {
	ID          string       `json:"id"`
	Number      string       `json:"number"`      // e.g., "A.9.1.1"
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    ControlCategory `json:"category"`
	Status      ControlStatus `json:"status"`
	Owner       string       `json:"owner"`
	Evidence    []Evidence   `json:"evidence"`
	Tests       []ControlTest `json:"tests"`
	LastTested  *time.Time   `json:"last_tested,omitempty"`
	NextReview  time.Time    `json:"next_review"`
}

// ControlCategory groups controls by domain
type ControlCategory string

const (
	CategoryAccessControl       ControlCategory = "access_control"
	CategoryInformationSecurity ControlCategory = "information_security"
	CategoryPhysicalSecurity    ControlCategory = "physical_security"
	CategoryOperationsSecurity  ControlCategory = "operations_security"
	CategoryCommunicationsSecurity ControlCategory = "communications_security"
	CategoryAcquisition        ControlCategory = "acquisition"
	CategoryIncidentManagement  ControlCategory = "incident_management"
	CategoryBusinessContinuity ControlCategory = "business_continuity"
	CategoryCompliance         ControlCategory = "compliance"
)

// ControlStatus represents control implementation status
type ControlStatus string

const (
	ControlStatusNotImplemented ControlStatus = "not_implemented"
	ControlStatusPartiallyImplemented ControlStatus = "partially_implemented"
	ControlStatusImplemented    ControlStatus = "implemented"
	ControlStatusTested         ControlStatus = "tested"
	ControlStatusEffective      ControlStatus = "effective"
)

// Clause represents an ISO 27001 clause
type Clause struct {
	Number      string    `json:"number"`      // e.g., "9.1"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Requirements []Requirement `json:"requirements"`
	Status      ClauseStatus `json:"status"`
	Owner       string    `json:"owner"`
	ReviewDate  time.Time `json:"review_date"`
}

// Requirement represents a clause requirement
type Requirement struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Controls    []string `json:"controls"` // Control IDs
	Status      RequirementStatus `json:"status"`
	Evidence    []string `json:"evidence"` // Evidence IDs
}

// RequirementStatus represents requirement fulfillment status
type RequirementStatus string

const (
	RequirementStatusNotMet    RequirementStatus = "not_met"
	RequirementStatusPartiallyMet RequirementStatus = "partially_met"
	RequirementStatusMet       RequirementStatus = "met"
	RequirementStatusExceeded   RequirementStatus = "exceeded"
)

// ClauseStatus represents clause implementation status
type ClauseStatus string

const (
	ClauseStatusNotStarted   ClauseStatus = "not_started"
	ClauseStatusInProgress   ClauseStatus = "in_progress"
	ClauseStatusCompleted    ClauseStatus = "completed"
	ClauseStatusUnderReview  ClauseStatus = "under_review"
	ClauseStatusApproved     ClauseStatus = "approved"
)

// RiskAssessment represents ISO 27001 risk assessment
type RiskAssessment struct {
	ID           string         `json:"id"`
	ISMSID       string         `json:"isms_id"`
	Risks        []ISMSRisk     `json:"risks"`
	Methodology  string         `json:"methodology"`
	Assessor     string         `json:"assessor"`
	AssessmentDate time.Time    `json:"assessment_date"`
	NextReview   time.Time      `json:"next_review"`
}

// ISMSRisk represents a risk in the ISMS
type ISMSRisk struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Threat       Threat        `json:"threat"`
	Vulnerability Vulnerability `json:"vulnerability"`
	Asset        Asset         `json:"asset"`
	Likelihood   Likelihood    `json:"likelihood"`
	Impact       Impact        `json:"impact"`
	RiskLevel    RiskLevel     `json:"risk_level"`
	Treatment    RiskTreatment `json:"treatment"`
	Owner        string        `json:"owner"`
	Status       RiskStatus    `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// Threat represents a threat scenario
type Threat struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    ThreatCategory `json:"category"`
	Source      ThreatSource `json:"source"`
	Capability  ThreatCapability `json:"capability"`
	Motivation  ThreatMotivation `json:"motivation"`
}

// ThreatCategory classifies threat types
type ThreatCategory string

const (
	ThreatCategoryMalicious    ThreatCategory = "malicious"
	ThreatCategoryAccidental  ThreatCategory = "accidental"
	ThreatCategoryEnvironmental ThreatCategory = "environmental"
	ThreatCategoryStructural   ThreatCategory = "structural"
)

// ThreatSource identifies threat origins
type ThreatSource string

const (
	ThreatSourceInternal     ThreatSource = "internal"
	ThreatSourceExternal     ThreatSource = "external"
	ThreatSourcePartner      ThreatSource = "partner"
	ThreatSourceSupplier     ThreatSource = "supplier"
)

// ThreatCapability rates threat capability
type ThreatCapability string

const (
	ThreatCapabilityLow      ThreatCapability = "low"
	ThreatCapabilityMedium   ThreatCapability = "medium"
	ThreatCapabilityHigh     ThreatCapability = "high"
	ThreatCapabilityAdvanced ThreatCapability = "advanced"
)

// ThreatMotivation identifies threat motivations
type ThreatMotivation string

const (
	ThreatMotivationFinancial   ThreatMotivation = "financial"
	ThreatMotivationEspionage   ThreatMotivation = "espionage"
	ThreatMotivationIdeological ThreatMotivation = "ideological"
	ThreatMotivationEgo         ThreatMotivation = "ego"
	ThreatMotivationAccidental  ThreatMotivation = "accidental"
)

// RiskLevel represents risk evaluation
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// RiskTreatment defines risk handling approach
type RiskTreatment string

const (
	RiskTreatmentAccept      RiskTreatment = "accept"
	RiskTreatmentMitigate    RiskTreatment = "mitigate"
	RiskTreatmentTransfer    RiskTreatment = "transfer"
	RiskTreatmentAvoid       RiskTreatment = "avoid"
)

// AuditEntry represents audit trail entries
type AuditEntry struct {
	ID          string      `json:"id"`
	Timestamp   time.Time   `json:"timestamp"`
	User        string      `json:"user"`
	Action      AuditAction `json:"action"`
	Resource    string      `json:"resource"`
	Details     string      `json:"details"`
	IPAddress   string      `json:"ip_address"`
	UserAgent   string      `json:"user_agent"`
}

// AuditAction types for audit logging
type AuditAction string

const (
	AuditActionCreate    AuditAction = "create"
	AuditActionUpdate    AuditAction = "update"
	AuditActionDelete    AuditAction = "delete"
	AuditActionView      AuditAction = "view"
	AuditActionExport    AuditAction = "export"
	AuditActionTest      AuditAction = "test"
	AuditActionApprove   AuditAction = "approve"
	AuditActionReject    AuditAction = "reject"
)

// ControlTest represents control testing
type ControlTest struct {
	ID          string       `json:"id"`
	ControlID   string       `json:"control_id"`
	TestType    TestType     `json:"test_type"`
	Description string       `json:"description"`
	Procedure   string       `json:"procedure"`
	Evidence    []string     `json:"evidence"`
	Result      TestResult   `json:"result"`
	Score       float64      `json:"score"`
	Tester      string       `json:"tester"`
	TestDate    time.Time    `json:"test_date"`
	NextTest    time.Time    `json:"next_test"`
	Findings    []Finding    `json:"findings"`
}

// TestType represents testing methodologies
type TestType string

const (
	TestTypeWalkthrough   TestType = "walkthrough"
	TestTypeInquiry       TestType = "inquiry"
	TestTypeObservation   TestType = "observation"
	TestTypeInspection    TestType = "inspection"
	TestTypeReperformance TestType = "reperformance"
	TestTypeAnalytical    TestType = "analytical"
)

// TestResult represents test outcomes
type TestResult string

const (
	TestResultPass       TestResult = "pass"
	TestResultFail       TestResult = "fail"
	TestResultPartial    TestResult = "partial"
	TestResultNotTested  TestResult = "not_tested"
)

// Finding represents audit findings
type Finding struct {
	ID          string       `json:"id"`
	TestID      string       `json:"test_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    FindingCategory `json:"category"`
	Severity    FindingSeverity `json:"severity"`
	Status      FindingStatus `json:"status"`
	Remediation string       `json:"remediation"`
	Owner       string       `json:"owner"`
	DueDate     time.Time    `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
}

// FindingCategory classifies finding types
type FindingCategory string

const (
	FindingCategoryControl    FindingCategory = "control"
	FindingCategoryProcess    FindingCategory = "process"
	FindingCategoryPolicy     FindingCategory = "policy"
	FindingCategoryTechnical  FindingCategory = "technical"
	FindingCategoryPhysical   FindingCategory = "physical"
)

// FindingSeverity rates finding severity
type FindingSeverity string

const (
	FindingSeverityLow      FindingSeverity = "low"
	FindingSeverityMedium   FindingSeverity = "medium"
	FindingSeverityHigh     FindingSeverity = "high"
	FindingSeverityCritical FindingSeverity = "critical"
)

// FindingStatus represents finding resolution status
type FindingStatus string

const (
	FindingStatusOpen        FindingStatus = "open"
	FindingStatusInProgress  FindingStatus = "in_progress"
	FindingStatusResolved    FindingStatus = "resolved"
	FindingStatusAccepted    FindingStatus = "accepted"
	FindingStatusClosed      FindingStatus = "closed"
)
