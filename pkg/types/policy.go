package types

import (
	"time"
)

// Policy Management Types

type Policy struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Content      string         `json:"content"`
	Category     PolicyCategory `json:"category"`
	Version      string         `json:"version"`
	Status       PolicyStatus   `json:"status"`
	Owner        string         `json:"owner"`
	Approvers    []string       `json:"approvers"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	PublishedAt  *time.Time     `json:"published_at,omitempty"`
	ReviewDate   *time.Time     `json:"review_date,omitempty"`
	Tags         []string       `json:"tags"`
}

type PolicyCategory struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Subcategories     []string `json:"subcategories"`
	ApprovalRequired  bool     `json:"approval_required"`
	ReviewPeriodDays  uint32   `json:"review_period_days"`
}

type PolicyStatus string

const (
	PolicyStatusDraft           PolicyStatus = "draft"
	PolicyStatusPendingApproval PolicyStatus = "pending_approval"
	PolicyStatusApproved        PolicyStatus = "approved"
	PolicyStatusPublished       PolicyStatus = "published"
	PolicyStatusUnderReview     PolicyStatus = "under_review"
	PolicyStatusArchived        PolicyStatus = "archived"
)

type PolicyRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	CategoryID  string   `json:"category_id"`
	Owner       string   `json:"owner"`
	Approvers   []string `json:"approvers"`
	Tags        []string `json:"tags"`
	Priority    string   `json:"priority"`
}

type PolicyPriority string

const (
	PolicyPriorityLow      PolicyPriority = "low"
	PolicyPriorityMedium   PolicyPriority = "medium"
	PolicyPriorityHigh     PolicyPriority = "high"
	PolicyPriorityCritical PolicyPriority = "critical"
)

type PolicyUpdateRequest struct {
	TitleChange   *string `json:"title_change,omitempty"`
	ContentChanges string  `json:"content_changes"`
	Reason        string  `json:"reason"`
	UpdatedBy     string  `json:"updated_by"`
}

type ApprovalWorkflow struct {
	ID        string      `json:"id"`
	PolicyID  string      `json:"policy_id"`
	Approvers []Approver  `json:"approvers"`
	CreatedAt time.Time   `json:"created_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	Status    string      `json:"status"`
}

type Approver struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	Role         string              `json:"role"`
	Decision     *ApprovalDecision   `json:"decision,omitempty"`
	DecisionDate *time.Time          `json:"decision_date,omitempty"`
	Comments     *string             `json:"comments,omitempty"`
}

type ApprovalDecision string

const (
	ApprovalDecisionApproved       ApprovalDecision = "approved"
	ApprovalDecisionRejected       ApprovalDecision = "rejected"
	ApprovalDecisionRequiresChanges ApprovalDecision = "requires_changes"
)

type WorkflowStatus string

const (
	WorkflowStatusPending  WorkflowStatus = "pending"
	WorkflowStatusApproved WorkflowStatus = "approved"
	WorkflowStatusRejected WorkflowStatus = "rejected"
	WorkflowStatusExpired  WorkflowStatus = "expired"
)

type PolicyComplianceStatus struct {
	PolicyID         string              `json:"policy_id"`
	PolicyTitle      string              `json:"policy_title"`
	ComplianceScore  float64             `json:"compliance_score"`
	LastChecked      time.Time           `json:"last_checked"`
	ComplianceIssues []ComplianceIssue   `json:"compliance_issues"`
	NextReviewDate   time.Time           `json:"next_review_date"`
}

type ComplianceIssue struct {
	ID                   string              `json:"id"`
	Description          string              `json:"description"`
	Severity             ComplianceSeverity  `json:"severity"`
	AffectedDepartments  []string            `json:"affected_departments"`
	RemediationRequired  bool                `json:"remediation_required"`
}

type ComplianceSeverity string

const (
	ComplianceSeverityLow      ComplianceSeverity = "low"
	ComplianceSeverityMedium   ComplianceSeverity = "medium"
	ComplianceSeverityHigh     ComplianceSeverity = "high"
	ComplianceSeverityCritical ComplianceSeverity = "critical"
)

type PolicyComplianceReport struct {
	ReportID            string                    `json:"report_id"`
	GeneratedAt         time.Time                 `json:"generated_at"`
	Scope               PolicyScope               `json:"scope"`
	PolicyCompliances   []PolicyComplianceStatus  `json:"policy_compliances"`
	OverallCompliance  float64                   `json:"overall_compliance"`
	Recommendations     []PolicyRecommendation    `json:"recommendations"`
}

type PolicyScope struct {
	Categories   []string       `json:"categories"`
	Departments  []string       `json:"departments"`
	Status       []PolicyStatus `json:"status"`
}

type PolicyRecommendation struct {
	Priority    RecommendationPriority `json:"priority"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	ActionItems []string               `json:"action_items"`
	Owner       string                 `json:"owner"`
	Timeline    string                 `json:"timeline"`
}

type RecommendationPriority string

const (
	RecommendationPriorityLow      RecommendationPriority = "low"
	RecommendationPriorityMedium   RecommendationPriority = "medium"
	RecommendationPriorityHigh     RecommendationPriority = "high"
	RecommendationPriorityCritical RecommendationPriority = "critical"
)

type PolicyStats struct {
	TotalPolicies           int     `json:"total_policies"`
	DraftPolicies           int     `json:"draft_policies"`
	PendingApproval         int     `json:"pending_approval"`
	PublishedPolicies       int     `json:"published_policies"`
	ArchivedPolicies        int     `json:"archived_policies"`
	AverageComplianceScore  float64 `json:"average_compliance_score"`
	PoliciesDueForReview    int     `json:"policies_due_for_review"`
	OverdueReviews          int     `json:"overdue_reviews"`
}

type PolicyDatabase struct {
	Policies           map[string]Policy               `json:"policies"`
	Workflows          map[string]ApprovalWorkflow    `json:"workflows"`
	ComplianceRecords  map[string]PolicyComplianceStatus `json:"compliance_records"`
}

// Framework Components

type PolicyEngine struct {
	Categories map[string]PolicyCategory `json:"categories"`
}

type PolicyComplianceManager struct {
	// Implementation details
}

type PolicyApprovalManager struct {
	// Implementation details
}

type PolicyDistributionManager struct {
	// Implementation details
}

type PolicyMonitoringManager struct {
	// Implementation details
}

// Configuration types

type PolicyConfig struct {
	EngineConfig       EngineConfig       `json:"engine_config"`
	ComplianceConfig   ComplianceConfig   `json:"compliance_config"`
	ApprovalConfig     ApprovalConfig     `json:"approval_config"`
	DistributionConfig DistributionConfig `json:"distribution_config"`
	MonitoringConfig  MonitoringConfig  `json:"monitoring_config"`
}

type EngineConfig struct {
	// Engine configuration
}

type ComplianceConfig struct {
	// Compliance configuration
}

type ApprovalConfig struct {
	// Approval configuration
}

type DistributionConfig struct {
	// Distribution configuration
}

type MonitoringConfig struct {
	// Monitoring configuration
}

// Error types

type PolicyError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e PolicyError) Error() string {
	return e.Message
}

// Approval result

type ApprovalResult struct {
	IsApproved bool   `json:"is_approved"`
	Message    string `json:"message"`
}
