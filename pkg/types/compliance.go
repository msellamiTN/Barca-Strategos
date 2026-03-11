package types

// FindingSeverity represents the severity of a compliance finding.
type FindingSeverity string

const (
	FindingSeverityLow    FindingSeverity = "low"
	FindingSeverityMedium FindingSeverity = "medium"
	FindingSeverityHigh   FindingSeverity = "high"
	FindingSeverityCritical FindingSeverity = "critical"
)

// RecommendationPriority represents the priority of a recommendation.
type RecommendationPriority string

const (
	RecommendationPriorityLow    RecommendationPriority = "low"
	RecommendationPriorityMedium RecommendationPriority = "medium"
	RecommendationPriorityHigh   RecommendationPriority = "high"
	RecommendationPriorityCritical RecommendationPriority = "critical"
)
