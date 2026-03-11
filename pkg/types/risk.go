package types

import "time"

// Risk represents a risk entry.
type Risk struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Likelihood  Likelihood     `json:"likelihood"`
	Impact      Impact         `json:"impact"`
	Score       float64        `json:"score"`
	Status      RiskStatus     `json:"status"`
	Mitigation  *Mitigation    `json:"mitigation,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Likelihood represents probability.
type Likelihood string

const (
	LikelihoodRare     Likelihood = "rare"
	LikelihoodUnlikely Likelihood = "unlikely"
	LikelihoodPossible Likelihood = "possible"
	LikelihoodLikely   Likelihood = "likely"
	LikelihoodCertain  Likelihood = "certain"
)

// Impact represents severity.
type Impact string

const (
	ImpactNegligible Impact = "negligible"
	ImpactMinor     Impact = "minor"
	ImpactModerate  Impact = "moderate"
	ImpactMajor     Impact = "major"
	ImpactCatastrophic Impact = "catastrophic"
)

// RiskStatus represents lifecycle.
type RiskStatus string

const (
	RiskStatusOpen        RiskStatus = "open"
	RiskStatusInProgress RiskStatus = "in_progress"
	RiskStatusMitigated   RiskStatus = "mitigated"
	RiskStatusAccepted   RiskStatus = "accepted"
)

// Mitigation tracks remediation.
type Mitigation struct {
	Plan       string    `json:"plan"`
	Owner      string    `json:"owner"`
	DueDate    time.Time `json:"due_date"`
	Status     string    `json:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
