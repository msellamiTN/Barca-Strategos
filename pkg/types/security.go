package types

// ThreatSeverity represents the severity level of a threat.
type ThreatSeverity string

const (
	ThreatSeverityLow      ThreatSeverity = "low"
	ThreatSeverityMedium   ThreatSeverity = "medium"
	ThreatSeverityHigh     ThreatSeverity = "high"
	ThreatSeverityCritical ThreatSeverity = "critical"
)

// EngineStatus indicates the status of a security engine.
type EngineStatus string

const (
	EngineStatusRunning EngineStatus = "running"
	EngineStatusStopped EngineStatus = "stopped"
	EngineStatusDegraded EngineStatus = "degraded"
)

// DataClassification defines data sensitivity levels.
type DataClassification string

const (
	DataClassificationPublic    DataClassification = "public"
	DataClassificationInternal  DataClassification = "internal"
	DataClassificationConfidential DataClassification = "confidential"
	DataClassificationRestricted DataClassification = "restricted"
)
