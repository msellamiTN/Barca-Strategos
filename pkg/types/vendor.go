package types

import (
	"time"
)

// Vendor Risk Management Types

type Vendor struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	Category         VendorCategory  `json:"category"`
	ContactEmail     string          `json:"contact_email"`
	ContactPhone     *string         `json:"contact_phone,omitempty"`
	Address          string          `json:"address"`
	ServicesOffered  []string        `json:"services_offered"`
	ContractStartDate time.Time      `json:"contract_start_date"`
	ContractEndDate  *time.Time      `json:"contract_end_date,omitempty"`
	Status           VendorStatus    `json:"status"`
	RiskLevel        *VendorRiskLevel `json:"risk_level,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	LastAssessed     *time.Time      `json:"last_assessed,omitempty"`
}

type VendorCategory struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	RiskFactors            []string `json:"risk_factors"`
	BaseRiskScore          float64  `json:"base_risk_score"`
	AssessmentFrequencyDays uint32   `json:"assessment_frequency_days"`
}

type VendorStatus string

const (
	VendorStatusActive      VendorStatus = "active"
	VendorStatusUnderReview VendorStatus = "under_review"
	VendorStatusSuspended   VendorStatus = "suspended"
	VendorStatusTerminated  VendorStatus = "terminated"
)

type VendorRiskLevel string

const (
	VendorRiskLevelMinimal  VendorRiskLevel = "minimal"
	VendorRiskLevelLow      VendorRiskLevel = "low"
	VendorRiskLevelMedium   VendorRiskLevel = "medium"
	VendorRiskLevelHigh     VendorRiskLevel = "high"
	VendorRiskLevelCritical VendorRiskLevel = "critical"
)

type VendorRequest struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	CategoryID         string     `json:"category_id"`
	ContactEmail       string     `json:"contact_email"`
	ContactPhone       *string    `json:"contact_phone,omitempty"`
	Address            string     `json:"address"`
	ServicesOffered    []string   `json:"services_offered"`
	ContractStartDate  time.Time  `json:"contract_start_date"`
	ContractEndDate    *time.Time `json:"contract_end_date,omitempty"`
}

type VendorUpdateRequest struct {
	NameChange              *string    `json:"name_change,omitempty"`
	DescriptionChange       *string    `json:"description_change,omitempty"`
	ContactEmailChange      *string    `json:"contact_email_change,omitempty"`
	ContactPhoneChange      *string    `json:"contact_phone_change,omitempty"`
	AddressChange           *string    `json:"address_change,omitempty"`
	ServicesChange          *[]string  `json:"services_change,omitempty"`
	ContractEndDateChange   *time.Time `json:"contract_end_date_change,omitempty"`
	RequiresReassessment    bool       `json:"requires_reassessment"`
}

type VendorAssessment struct {
	AssessmentID         string                      `json:"assessment_id"`
	VendorID             string                      `json:"vendor_id"`
	VendorName           string                      `json:"vendor_name"`
	Timestamp            time.Time                   `json:"timestamp"`
	RiskAssessment       VendorRiskAssessment        `json:"risk_assessment"`
	ComplianceAssessment VendorComplianceAssessment  `json:"compliance_assessment"`
	OverallRiskScore     float64                     `json:"overall_risk_score"`
	RiskLevel            VendorRiskLevel             `json:"risk_level"`
	Recommendations      []VendorRecommendation      `json:"recommendations"`
	NextAssessmentDate   time.Time                   `json:"next_assessment_date"`
}

type VendorRiskAssessment struct {
	RiskScore          float64      `json:"risk_score"`
	RiskFactors        []RiskFactor `json:"risk_factors"`
	FinancialStability float64      `json:"financial_stability"`
	OperationalCapability float64   `json:"operational_capability"`
	SecurityPosture     float64      `json:"security_posture"`
	ReputationScore     float64      `json:"reputation_score"`
}

type RiskFactor struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Severity          RiskFactorSeverity `json:"severity"`
	Score             float64            `json:"score"`
	MitigationRequired bool              `json:"mitigation_required"`
}

type RiskFactorSeverity string

const (
	RiskFactorSeverityLow      RiskFactorSeverity = "low"
	RiskFactorSeverityMedium   RiskFactorSeverity = "medium"
	RiskFactorSeverityHigh     RiskFactorSeverity = "high"
	RiskFactorSeverityCritical RiskFactorSeverity = "critical"
)

type VendorComplianceAssessment struct {
	ComplianceScore     float64           `json:"compliance_score"`
	ComplianceIssues    []ComplianceIssue `json:"compliance_issues"`
	Certifications      []Certification  `json:"certifications"`
	RegulatoryAdherence float64           `json:"regulatory_adherence"`
	PolicyCompliance    float64           `json:"policy_compliance"`
}

type ComplianceIssue struct {
	ID                   string              `json:"id"`
	Description          string              `json:"description"`
	Severity             ComplianceSeverity  `json:"severity"`
	AffectedRegulations  []string            `json:"affected_regulations"`
	RemediationPlan      *string             `json:"remediation_plan,omitempty"`
}

type Certification struct {
	Name        string             `json:"name"`
	Issuer      string             `json:"issuer"`
	ObtainedDate time.Time          `json:"obtained_date"`
	ExpiryDate  time.Time          `json:"expiry_date"`
	Status      CertificationStatus `json:"status"`
}

type CertificationStatus string

const (
	CertificationStatusActive   CertificationStatus = "active"
	CertificationStatusExpired  CertificationStatus = "expired"
	CertificationStatusRevoked  CertificationStatus = "revoked"
	CertificationStatusPending  CertificationStatus = "pending"
)

type VendorRecommendation struct {
	Priority    RecommendationPriority `json:"priority"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	ActionItems []string               `json:"action_items"`
	Owner       string                 `json:"owner"`
	Timeline    string                 `json:"timeline"`
}

type VendorRiskStatus struct {
	VendorID         string           `json:"vendor_id"`
	VendorName       string           `json:"vendor_name"`
	RiskLevel        VendorRiskLevel  `json:"risk_level"`
	RiskScore        float64          `json:"risk_score"`
	LastAssessed     time.Time        `json:"last_assessed"`
	NextAssessment   time.Time        `json:"next_assessment"`
	CriticalIssues   []string         `json:"critical_issues"`
	MonitoringStatus MonitoringStatus `json:"monitoring_status"`
}

type MonitoringStatus struct {
	IsActive      bool      `json:"is_active"`
	LastCheck     time.Time `json:"last_check"`
	AlertsCount   uint32    `json:"alerts_count"`
	IssuesDetected uint32   `json:"issues_detected"`
}

type VendorComplianceStatus struct {
	VendorID         string              `json:"vendor_id"`
	ComplianceScore  float64             `json:"compliance_score"`
	LastChecked      time.Time           `json:"last_checked"`
	ComplianceIssues []ComplianceIssue   `json:"compliance_issues"`
	NextReviewDate   time.Time           `json:"next_review_date"`
}

type VendorScope struct {
	Categories  []string         `json:"categories"`
	RiskLevels  []VendorRiskLevel `json:"risk_levels"`
	Status      []VendorStatus    `json:"status"`
}

type VendorRiskReport struct {
	ReportID           string                    `json:"report_id"`
	GeneratedAt        time.Time                 `json:"generated_at"`
	Scope              VendorScope               `json:"scope"`
	TotalVendors       int                       `json:"total_vendors"`
	VendorAssessments  []VendorAssessmentPair    `json:"vendor_assessments"`
	OverallRiskProfile VendorRiskProfile        `json:"overall_risk_profile"`
	HighRiskVendors    []VendorAssessmentPair    `json:"high_risk_vendors"`
	Recommendations    []VendorRecommendation    `json:"recommendations"`
}

type VendorAssessmentPair struct {
	Vendor      Vendor          `json:"vendor"`
	Assessment  VendorAssessment `json:"assessment"`
}

type VendorRiskProfile struct {
	OverallRiskScore  float64            `json:"overall_risk_score"`
	CriticalVendors   int                `json:"critical_vendors"`
	HighRiskVendors   int                `json:"high_risk_vendors"`
	MediumRiskVendors int                `json:"medium_risk_vendors"`
	LowRiskVendors    int                `json:"low_risk_vendors"`
	MinimalRiskVendors int                `json:"minimal_risk_vendors"`
	RiskDistribution  map[string]int     `json:"risk_distribution"`
}

type VendorStats struct {
	TotalVendors             int     `json:"total_vendors"`
	ActiveVendors             int     `json:"active_vendors"`
	CriticalRiskVendors       int     `json:"critical_risk_vendors"`
	HighRiskVendors           int     `json:"high_risk_vendors"`
	MediumRiskVendors         int     `json:"medium_risk_vendors"`
	LowRiskVendors            int     `json:"low_risk_vendors"`
	AverageRiskScore          float64 `json:"average_risk_score"`
	VendorsDueForAssessment   int     `json:"vendors_due_for_assessment"`
	OverdueAssessments        int     `json:"overdue_assessments"`
}

type VendorDatabase struct {
	Vendors      map[string]Vendor                    `json:"vendors"`
	Assessments  map[string][]VendorAssessment        `json:"assessments"`
}

// Framework Components

type VendorEngine struct {
	Categories map[string]VendorCategory `json:"categories"`
}

type VendorAssessmentManager struct {
	// Implementation details
}

type VendorMonitoringManager struct {
	// Implementation details
}

type VendorComplianceManager struct {
	// Implementation details
}

type VendorReportingManager struct {
	// Implementation details
}

// Configuration types

type VendorConfig struct {
	EngineConfig       EngineConfig       `json:"engine_config"`
	AssessmentConfig   AssessmentConfig   `json:"assessment_config"`
	MonitoringConfig   MonitoringConfig   `json:"monitoring_config"`
	ComplianceConfig   ComplianceConfig   `json:"compliance_config"`
	ReportingConfig    ReportingConfig    `json:"reporting_config"`
}

type AssessmentConfig struct {
	// Assessment configuration
}

type MonitoringConfig struct {
	// Monitoring configuration
}

type ReportingConfig struct {
	// Reporting configuration
}

// Error types

type VendorError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e VendorError) Error() string {
	return e.Message
}
