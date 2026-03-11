package types

import "time"

// MonitoringConfig holds runtime monitoring settings.
type MonitoringConfig struct {
	CheckIntervalSeconds uint64  `json:"check_interval_seconds"`
	EnableRealTime       bool    `json:"enable_real_time"`
	AlertThreshold       float64 `json:"alert_threshold"`
}

// ComplianceMonitor tracks a compliance framework status.
type ComplianceMonitor struct {
	MonitorID string          `json:"monitor_id"`
	Framework string          `json:"framework"`
	Status    MonitorStatus    `json:"status"`
	LastRun   time.Time       `json:"last_run"`
}

// MonitorStatus represents the health of a monitor.
type MonitorStatus string

const (
	MonitorStatusHealthy   MonitorStatus = "healthy"
	MonitorStatusDegraded  MonitorStatus = "degraded"
	MonitorStatusCritical  MonitorStatus = "critical"
	MonitorStatusDisabled  MonitorStatus = "disabled"
)
