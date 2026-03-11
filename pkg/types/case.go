package types

import "time"

// Case represents an incident case.
type Case struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      CaseStatus  `json:"status"`
	Severity    Severity    `json:"severity"`
	AssigneeID  string      `json:"assignee_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ResolvedAt *time.Time  `json:"resolved_at,omitempty"`
	SLA         *SLA        `json:"sla,omitempty"`
}

// CaseStatus represents lifecycle.
type CaseStatus string

const (
	CaseStatusNew        CaseStatus = "new"
	CaseStatusInProgress CaseStatus = "in_progress"
	CaseStatusResolved   CaseStatus = "resolved"
	CaseStatusClosed     CaseStatus = "closed"
)

// Severity represents impact.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// SLA tracks service level agreement.
type SLA struct {
	ResponseDue time.Time `json:"response_due"`
	ResolveDue  time.Time `json:"resolve_due"`
}

// Action represents an automated response action.
type Action struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"case_id"`
	Type        ActionType `json:"type"`
	Status       string    `json:"status"`
	ExecutedAt  time.Time `json:"executed_at"`
	ExecutedBy  string    `json:"executed_by"`
	Details     string    `json:"details,omitempty"`
}

// ActionType represents action types.
type ActionType string

const (
	ActionIsolateHost   ActionType = "isolate_host"
	ActionBlockIP       ActionType = "block_ip"
	ActionDisableUser   ActionType = "disable_user"
	ActionRunPlaybook   ActionType = "run_playbook"
	ActionCreateTicket  ActionType = "create_ticket"
)
