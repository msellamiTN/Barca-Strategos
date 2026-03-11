package types

import "time"

// Asset represents an infrastructure asset.
type Asset struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Type        AssetType       `json:"type"`
	IPAddress  string          `json:"ip_address"`
	Owner       string          `json:"owner"`
	Criticality  AssetCriticality  `json:"criticality"`
	Location   string          `json:"location"`
	OS          string          `json:"os"`
	Version     string          `json:"version"`
	Tags        []string        `json:"tags"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// AssetType categorizes assets.
type AssetType string

const (
	AssetTypeServer    AssetType = "server"
	AssetTypeWorkstation AssetType = "workstation"
	AssetTypeNetwork  AssetType = "network"
	AssetTypeDatabase AssetType = "database"
	AssetTypeApplication AssetType = "application"
)

// AssetCriticality defines impact levels.
type AssetCriticality string

const (
	AssetCriticalityLow      AssetCriticality = "low"
	AssetCriticalityMedium   AssetCriticality = "medium"
	AssetCriticalityHigh     AssetCriticality = "high"
	AssetCriticalityCritical AssetCriticality = "critical"
)

// Vulnerability represents a security finding.
type Vulnerability struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Severity    string    `json:"severity"`
	CVSS        string    `json:"cvss"`
	Source      string    `json:"source"`
	DiscoveredAt time.Time `json:"discovered_at"`
	References []string `json:"references"`
}

// ThreatIntel represents a threat intelligence indicator.
type ThreatIntel struct {
	IOCType     ThreatIOCType `json:"ioc_type"`
	Value       string        `json:"value"`
	Confidence  float64       `json:"confidence"`
	Source      string        `json:"source"`
	Description string        `json:"description"`
	ValidUntil  *time.Time     `json:"valid_until"`
}

// ThreatIOCType defines indicator types.
type ThreatIOCType string

const (
	ThreatIOCTypeIP       ThreatIOCType = "ip"
	ThreatIOCTypeDomain   ThreatIOCType = "domain"
	ThreatIOCTypeHash     ThreatIOCType = "hash"
	ThreatIOCTypeURL      ThreatIOCType = "url"
)
