package types

import "time"

// GDPR represents GDPR compliance management
type GDPR struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Status      GDPRStatus  `json:"status"`
	Principles  []Principle `json:"principles"`
	Processes   []Process   `json:"processes"`
	Rights      []Right     `json:"rights"`
	Breaches    []Breach    `json:"breaches"`
	Assessments []DPIA      `json:"assessments"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// GDPRStatus represents implementation status
type GDPRStatus string

const (
	GDPRStatusDraft        GDPRStatus = "draft"
	GDPRStatusActive       GDPRStatus = "active"
	GDPRStatusUnderReview  GDPRStatus = "under_review"
	GDPRStatusApproved     GDPRStatus = "approved"
	GDPRStatusNonCompliant GDPRStatus = "non_compliant"
)

// Principle represents GDPR data protection principles
type Principle struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Article     string            `json:"article"`
	Category    PrincipleCategory `json:"category"`
	Status      PrincipleStatus   `json:"status"`
	Controls    []Control         `json:"controls"`
	Owner       string            `json:"owner"`
	LastReview  *time.Time        `json:"last_review,omitempty"`
	NextReview  time.Time         `json:"next_review"`
}

// PrincipleCategory represents principle categories
type PrincipleCategory string

const (
	PrincipleCategoryLawfulness     PrincipleCategory = "lawfulness"
	PrincipleCategoryFairness       PrincipleCategory = "fairness"
	PrincipleCategoryTransparency   PrincipleCategory = "transparency"
	PrincipleCategoryPurpose        PrincipleCategory = "purpose"
	PrincipleCategoryDataMinimization PrincipleCategory = "data_minimization"
	PrincipleCategoryAccuracy       PrincipleCategory = "accuracy"
	PrincipleCategoryStorageLimit   PrincipleCategory = "storage_limit"
	PrincipleCategorySecurity       PrincipleCategory = "security"
	PrincipleCategoryAccountability PrincipleCategory = "accountability"
)

// PrincipleStatus represents principle implementation status
type PrincipleStatus string

const (
	PrincipleStatusNotImplemented PrincipleStatus = "not_implemented"
	PrincipleStatusPartiallyImplemented PrincipleStatus = "partially_implemented"
	PrincipleStatusImplemented    PrincipleStatus = "implemented"
	PrincipleStatusValidated      PrincipleStatus = "validated"
	PrincipleStatusCompliant      PrincipleStatus = "compliant"
)

// Process represents data processing activities
type Process struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Purpose       string        `json:"purpose"`
	LegalBasis    LegalBasis    `json:"legal_basis"`
	DataCategories []DataCategory `json:"data_categories"`
	DataSubjects   []DataSubject  `json:"data_subjects"`
	Processors    []Processor   `json:"processors"`
	Controllers   []Controller  `json:"controllers"`
	Transfers     []Transfer     `json:"transfers"`
	Retention     Retention     `json:"retention"`
	Security      Security      `json:"security"`
	Status        ProcessStatus `json:"status"`
	Owner         string        `json:"owner"`
	LastReview    *time.Time    `json:"last_review,omitempty"`
	NextReview    time.Time     `json:"next_review"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// LegalBasis represents lawful basis for processing
type LegalBasis string

const (
	LegalBasisConsent          LegalBasis = "consent"
	LegalBasisContract         LegalBasis = "contract"
	LegalBasisLegalObligation  LegalBasis = "legal_obligation"
	LegalBasisVitalInterests    LegalBasis = "vital_interests"
	LegalBasisPublicTask       LegalBasis = "public_task"
	LegalBasisLegitimateInterests LegalBasis = "legitimate_interests"
)

// DataCategory represents types of personal data
type DataCategory string

const (
	DataCategoryPersonal        DataCategory = "personal"
	DataCategorySpecial         DataCategory = "special"
	DataCategorySensitive       DataCategory = "sensitive"
	DataCategoryBiometric       DataCategory = "biometric"
	DataCategoryGenetic         DataCategory = "genetic"
	DataCategoryHealth          DataCategory = "health"
	DataCategoryFinancial       DataCategory = "financial"
	DataCategoryLocation        DataCategory = "location"
	DataCategoryCommunication   DataCategory = "communication"
	DataCategoryBehavioral      DataCategory = "behavioral"
	DataCategoryProfiling       DataCategory = "profiling"
)

// DataSubject represents categories of data subjects
type DataSubject string

const (
	DataSubjectCustomers        DataSubject = "customers"
	DataSubjectEmployees         DataSubject = "employees"
	DataSubjectProspects         DataSubject = "prospects"
	DataSubjectSuppliers         DataSubject = "suppliers"
	DataSubjectPartners          DataSubject = "partners"
	DataSubjectWebsiteVisitors   DataSubject = "website_visitors"
	DataSubjectMarketingLeads    DataSubject = "marketing_leads"
	DataSubjectJobApplicants     DataSubject = "job_applicants"
	DataSubjectContractors        DataSubject = "contractors"
	DataSubjectMinors            DataSubject = "minors"
)

// Processor represents data processors
type Processor struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        ProcessorType  `json:"type"`
	Country     string         `json:"country"`
	Contact     string         `json:"contact"`
	Contract    string         `json:"contract"`
	Security    Security       `json:"security"`
	Certifications []string   `json:"certifications"`
	Status      ProcessorStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ProcessorType represents processor types
type ProcessorType string

const (
	ProcessorTypeInternal    ProcessorType = "internal"
	ProcessorTypeExternal    ProcessorType = "external"
	ProcessorTypeCloud       ProcessorType = "cloud"
	ProcessorTypeSaaS        ProcessorType = "saas"
	ProcessorTypeConsultant  ProcessorType = "consultant"
	ProcessorTypePartner     ProcessorType = "partner"
)

// ProcessorStatus represents processor status
type ProcessorStatus string

const (
	ProcessorStatusActive        ProcessorStatus = "active"
	ProcessorStatusInactive      ProcessorStatus = "inactive"
	ProcessorStatusUnderReview   ProcessorStatus = "under_review"
	ProcessorStatusTerminated    ProcessorStatus = "terminated"
)

// Controller represents data controllers
type Controller struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Type        ControllerType  `json:"type"`
	Country     string          `json:"country"`
	Contact     string          `json:"contact"`
	DPO         string          `json:"dpo"`
	Rep         string          `json:"rep"`
	Registration string         `json:"registration"`
	Status      ControllerStatus `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ControllerType represents controller types
type ControllerType string

const (
	ControllerTypePrimary   ControllerType = "primary"
	ControllerTypeSecondary ControllerType = "secondary"
	ControllerTypeJoint    ControllerType = "joint"
	ControllerTypeProcessor ControllerType = "processor"
)

// ControllerStatus represents controller status
type ControllerStatus string

const (
	ControllerStatusActive     ControllerStatus = "active"
	ControllerStatusInactive   ControllerStatus = "inactive"
	ControllerStatusRegistered ControllerStatus = "registered"
	ControllerStatusExempt     ControllerStatus = "exempt"
)

// Transfer represents international data transfers
type Transfer struct {
	ID           string        `json:"id"`
	ProcessID    string        `json:"process_id"`
	Destination  string        `json:"destination"`
	Country      string        `json:"country"`
	Mechanism    TransferMechanism `json:"mechanism"`
	Assessment   string        `json:"assessment"`
	Documentation string       `json:"documentation"`
	Status       TransferStatus `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// TransferMechanism represents transfer mechanisms
type TransferMechanism string

const (
	TransferMechanismAdequacy    TransferMechanism = "adequacy"
	TransferMechanismSCCs        TransferMechanism = "sccs"
	TransferMechanismBCRs        TransferMechanism = "bcrs"
	TransferMechanismDerogation  TransferMechanism = "derogation"
	TransferMechanismBindingRules TransferMechanism = "binding_rules"
)

// TransferStatus represents transfer status
type TransferStatus string

const (
	TransferStatusAuthorized  TransferStatus = "authorized"
	TransferStatusRestricted  TransferStatus = "restricted"
	TransferStatusProhibited  TransferStatus = "prohibited"
	TransferStatusUnderReview TransferStatus = "under_review"
)

// Retention represents data retention policies
type Retention struct {
	ID          string       `json:"id"`
	Policy      string       `json:"policy"`
	Period      string       `json:"period"`
	Criteria    string       `json:"criteria"`
	Archival    bool         `json:"archival"`
	Deletion    bool         `json:"deletion"`
	Anonymization bool       `json:"anonymization"`
	Status      RetentionStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// RetentionStatus represents retention status
type RetentionStatus string

const (
	RetentionStatusActive    RetentionStatus = "active"
	RetentionStatusExpired  RetentionStatus = "expired"
	RetentionStatusArchived RetentionStatus = "archived"
	RetentionStatusDeleted  RetentionStatus = "deleted"
)

// Security represents security measures
type Security struct {
	ID          string       `json:"id"`
	Measures    []SecurityMeasure `json:"measures"`
	Encryption  Encryption   `json:"encryption"`
	Access      Access       `json:"access"`
	Monitoring  Monitoring   `json:"monitoring"`
	Incidents   []Incident   `json:"incidents"`
	Status      SecurityStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// SecurityMeasure represents security measures
type SecurityMeasure struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    SecurityCategory `json:"category"`
	Implemented bool           `json:"implemented"`
	Tested      bool           `json:"tested"`
	LastTest    *time.Time     `json:"last_test,omitempty"`
	NextTest    time.Time      `json:"next_test"`
}

// SecurityCategory represents security categories
type SecurityCategory string

const (
	SecurityCategoryTechnical     SecurityCategory = "technical"
	SecurityCategoryOrganizational SecurityCategory = "organizational"
	SecurityCategoryPhysical      SecurityCategory = "physical"
	SecurityCategoryAdministrative SecurityCategory = "administrative"
)

// Encryption represents encryption controls
type Encryption struct {
	AtRest      bool   `json:"at_rest"`
	InTransit   bool   `json:"in_transit"`
	Algorithms  []string `json:"algorithms"`
	KeyManagement string `json:"key_management"`
	LastReview   *time.Time `json:"last_review,omitempty"`
}

// Access represents access controls
type Access struct {
	RBAC        bool     `json:"rbac"`
	MFA         bool     `json:"mfa"`
	LeastPrivilege bool   `json:"least_privilege"`
	Auditing    bool     `json:"auditing"`
	ReviewCycle string   `json:"review_cycle"`
	LastReview  *time.Time `json:"last_review,omitempty"`
}

// Monitoring represents security monitoring
type Monitoring struct {
	Logging     bool     `json:"logging"`
	Alerting    bool     `json:"alerting"`
	Analytics   bool     `json:"analytics"`
	Tools       []string `json:"tools"`
	ReviewCycle string   `json:"review_cycle"`
	LastReview  *time.Time `json:"last_review,omitempty"`
}

// SecurityStatus represents security status
type SecurityStatus string

const (
	SecurityStatusCompliant    SecurityStatus = "compliant"
	SecurityStatusNonCompliant SecurityStatus = "non_compliant"
	SecurityStatusPartially    SecurityStatus = "partially"
	SecurityStatusUnderReview  SecurityStatus = "under_review"
)

// ProcessStatus represents process status
type ProcessStatus string

const (
	ProcessStatusActive      ProcessStatus = "active"
	ProcessStatusInactive    ProcessStatus = "inactive"
	ProcessStatusUnderReview ProcessStatus = "under_review"
	ProcessStatusSuspended   ProcessStatus = "suspended"
	ProcessStatusTerminated  ProcessStatus = "terminated"
)

// Right represents data subject rights
type Right struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Article      string      `json:"article"`
	Category     RightCategory `json:"category"`
	Status       RightStatus  `json:"status"`
	Procedures   []Procedure  `json:"procedures"`
	Requests     []Request    `json:"requests"`
	Owner        string       `json:"owner"`
	LastReview   *time.Time   `json:"last_review,omitempty"`
	NextReview   time.Time    `json:"next_review"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// RightCategory represents right categories
type RightCategory string

const (
	RightCategoryTransparency RightCategory = "transparency"
	RightCategoryAccess      RightCategory = "access"
	RightCategoryRectification RightCategory = "rectification"
	RightCategoryErasure     RightCategory = "erasure"
	RightCategoryPortability RightCategory = "portability"
	RightCategoryObjection   RightCategory = "objection"
	RightCategoryRestriction RightCategory = "restriction"
	RightCategoryAutomatedDecision RightCategory = "automated_decision"
)

// RightStatus represents right implementation status
type RightStatus string

const (
	RightStatusNotImplemented RightStatus = "not_implemented"
	RightStatusPartiallyImplemented RightStatus = "partially_implemented"
	RightStatusImplemented    RightStatus = "implemented"
	RightStatusOperational    RightStatus = "operational"
	RightStatusCompliant      RightStatus = "compliant"
)

// Procedure represents procedures for handling rights
type Procedure struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Steps       []string       `json:"steps"`
	Timeline    string         `json:"timeline"`
	Responsible string         `json:"responsible"`
	Status      ProcedureStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ProcedureStatus represents procedure status
type ProcedureStatus string

const (
	ProcedureStatusDraft        ProcedureStatus = "draft"
	ProcedureStatusActive       ProcedureStatus = "active"
	ProcedureStatusUnderReview  ProcedureStatus = "under_review"
	ProcedureStatusApproved     ProcedureStatus = "approved"
	ProcedureStatusDeprecated   ProcedureStatus = "deprecated"
)

// Request represents data subject requests
type Request struct {
	ID           string       `json:"id"`
	RightID      string       `json:"right_id"`
	Type         RequestType  `json:"type"`
	DataSubject  string       `json:"data_subject"`
	Identity     string       `json:"identity"`
	Description  string       `json:"description"`
	Status       RequestStatus `json:"status"`
	Priority     Priority     `json:"priority"`
	AssignedTo   string       `json:"assigned_to"`
	DueDate      time.Time    `json:"due_date"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
	Resolution   string       `json:"resolution"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// RequestType represents request types
type RequestType string

const (
	RequestTypeAccess         RequestType = "access"
	RequestTypeRectification  RequestType = "rectification"
	RequestTypeErasure       RequestType = "erasure"
	RequestTypePortability    RequestType = "portability"
	RequestTypeRestriction   RequestType = "restriction"
	RequestTypeObjection     RequestType = "objection"
	RequestTypeInformation   RequestType = "information"
	RequestTypeAutomatedDecision RequestType = "automated_decision"
)

// RequestStatus represents request status
type RequestStatus string

const (
	RequestStatusReceived     RequestStatus = "received"
	RequestStatusInProgress   RequestStatus = "in_progress"
	RequestStatusPending      RequestStatus = "pending"
	RequestStatusCompleted    RequestStatus = "completed"
	RequestStatusRejected     RequestStatus = "rejected"
	RequestStatusWithdrawn    RequestStatus = "withdrawn"
)

// Breach represents data breaches
type Breach struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Date         time.Time    `json:"date"`
	Detected     time.Time    `json:"detected"`
	Type         BreachType   `json:"type"`
	Categories   []DataCategory `json:"categories"`
	Affected     int          `json:"affected"`
	Cause        string       `json:"cause"`
	Impact       Impact       `json:"impact"`
	Notification Notification `json:"notification"`
	Remediation  string       `json:"remediation"`
	Status       BreachStatus `json:"status"`
	Owner        string       `json:"owner"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// BreachType represents breach types
type BreachType string

const (
	BreachTypeUnauthorizedAccess BreachType = "unauthorized_access"
	BreachTypeDataLeakage        BreachType = "data_leakage"
	BreachTypeRansomware         BreachType = "ransomware"
	BreachTypePhishing          BreachType = "phishing"
	BreachTypeInsiderThreat     BreachType = "insider_threat"
	BreachTypePhysicalLoss       BreachType = "physical_loss"
	BreachTypeSystemError       BreachType = "system_error"
	BreachTypeThirdParty        BreachType = "third_party"
)

// Impact represents breach impact assessment
type Impact struct {
	Level        ImpactLevel `json:"level"`
	Likelihood   Likelihood  `json:"likelihood"`
	Severity     Severity    `json:"severity"`
	Assessment   string      `json:"assessment"`
	Consequences []string    `json:"consequences"`
}

// ImpactLevel represents impact levels
type ImpactLevel string

const (
	ImpactLevelLow      ImpactLevel = "low"
	ImpactLevelMedium   ImpactLevel = "medium"
	ImpactLevelHigh     ImpactLevel = "high"
	ImpactLevelCritical ImpactLevel = "critical"
)

// Severity represents severity levels
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Notification represents breach notification requirements
type Notification struct {
	Required      bool     `json:"required"`
	Deadline      time.Time `json:"deadline"`
	Authority     string   `json:"authority"`
	DataSubjects  bool     `json:"data_subjects"`
	Content       string   `json:"content"`
	Sent          bool     `json:"sent"`
	SentAt        *time.Time `json:"sent_at,omitempty"`
}

// BreachStatus represents breach status
type BreachStatus string

const (
	BreachStatusOpen        BreachStatus = "open"
	BreachStatusInvestigating BreachStatus = "investigating"
	BreachStatusContainment BreachStatus = "containment"
	BreachStatusResolved    BreachStatus = "resolved"
	BreachStatusClosed      BreachStatus = "closed"
)

// DPIA represents Data Protection Impact Assessments
type DPIA struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ProcessID    string     `json:"process_id"`
	Controller   string     `json:"controller"`
	Assessor     string     `json:"assessor"`
	Status       DPIAStatus `json:"status"`
	Risk         Risk       `json:"risk"`
	Measures     []Measure  `json:"measures"`
	Outcome      string     `json:"outcome"`
	Approval     string     `json:"approval"`
	ReviewDate   time.Time  `json:"review_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// DPIAStatus represents DPIA status
type DPIAStatus string

const (
	DPIAStatusRequired    DPIAStatus = "required"
	DPIAStatusInProgress  DPIAStatus = "in_progress"
	DPIAStatusCompleted   DPIAStatus = "completed"
	DPIAStatusApproved    DPIAStatus = "approved"
	DPIAStatusRejected    DPIAStatus = "rejected"
	DPIAStatusNotRequired DPIAStatus = "not_required"
)

// Risk represents DPIA risk assessment
type Risk struct {
	Level        RiskLevel   `json:"level"`
	Likelihood   Likelihood  `json:"likelihood"`
	Impact       Impact      `json:"impact"`
	Description  string      `json:"description"`
	Mitigation   string      `json:"mitigation"`
}

// Measure represents mitigation measures
type Measure struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Type        MeasureType  `json:"type"`
	Implemented bool         `json:"implemented"`
	Effective  bool         `json:"effective"`
	CreatedAt   time.Time    `json:"created_at"`
}

// MeasureType represents measure types
type MeasureType string

const (
	MeasureTypeTechnical     MeasureType = "technical"
	MeasureTypeOrganizational MeasureType = "organizational"
	MeasureTypeLegal        MeasureType = "legal"
	MeasureTypeProcedural   MeasureType = "procedural"
)

// Incident represents security incidents
type Incident struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Date        time.Time    `json:"date"`
	Type        IncidentType `json:"type"`
	Severity    Severity    `json:"severity"`
	Status      IncidentStatus `json:"status"`
	Affected    []string     `json:"affected"`
	Response    string       `json:"response"`
	Resolution  string       `json:"resolution"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// IncidentType represents incident types
type IncidentType string

const (
	IncidentTypeMalware      IncidentType = "malware"
	IncidentTypePhishing     IncidentType = "phishing"
	IncidentTypeInsider      IncidentType = "insider"
	IncidentTypeDenial       IncidentType = "denial"
	IncidentTypePhysical     IncidentType = "physical"
	IncidentTypeHumanError   IncidentType = "human_error"
	IncidentTypeSystemFailure IncidentType = "system_failure"
)

// IncidentStatus represents incident status
type IncidentStatus string

const (
	IncidentStatusOpen        IncidentStatus = "open"
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusContainment IncidentStatus = "containment"
	IncidentStatusResolved    IncidentStatus = "resolved"
	IncidentStatusClosed      IncidentStatus = "closed"
)
