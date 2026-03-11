# Rust-to-Go Types Mapping

This document maps the Rust types from `old_rust/src/common/types.rs` to the new Go types in `pkg/types/`.

## User Management

| Rust | Go |
|------|----|
| `UserId(pub String)` | `User` struct with `ID string` |
| `User` | `User` struct with `ID`, `Username`, `Email`, `Role`, `CreatedAt` |
| `UserRole` enum | `UserRole` type with const values (`RoleAdmin`, `RoleSecurityAnalyst`, `RoleComplianceOfficer`, `RoleRiskManager`) |

## Monitoring

| Rust | Go |
|------|----|
| `MonitoringConfig` | `MonitoringConfig` struct with JSON tags |
| `ComplianceMonitor` | `ComplianceMonitor` struct with `LastRun time.Time` |
| `MonitorStatus` enum | `MonitorStatus` type with const values (`MonitorStatusHealthy`, `MonitorStatusDegraded`, etc.) |

## Compliance

| Rust | Go |
|------|----|
| `FindingSeverity` enum | `FindingSeverity` type with const values (`FindingSeverityLow`, etc.) |
| `RecommendationPriority` enum | `RecommendationPriority` type with const values (`RecommendationPriorityLow`, etc.) |

## Security

| Rust | Go |
|------|----|
| `ThreatSeverity` enum | `ThreatSeverity` type with const values |
| `EngineStatus` enum | `EngineStatus` type with const values |
| `DataClassification` enum | `DataClassification` type with const values |

## GUI

| Rust | Go |
|------|----|
| `InterfaceComplexity` enum | `InterfaceComplexity` type with const values |
| `LayoutType` enum | `LayoutType` type with const values |
| `ThemeType` enum | `ThemeType` type with const values |

## Usage Example

```go
import "github.com/barca-strategos/phoenix/pkg/types"

user := types.User{
    ID:        "123",
    Username:  "alice",
    Email:     "alice@example.com",
    Role:      types.RoleSecurityAnalyst,
    CreatedAt: time.Now(),
}
```

All Go types include JSON tags for API serialization and validation helpers where applicable.
