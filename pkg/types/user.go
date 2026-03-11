package types

import (
	"time"
)

// User represents a system user.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRole defines user permissions.
type UserRole string

const (
	RoleAdmin            UserRole = "admin"
	RoleSecurityAnalyst  UserRole = "security_analyst"
	RoleComplianceOfficer UserRole = "compliance_officer"
	RoleRiskManager      UserRole = "risk_manager"
)

// Validate checks if the role is known.
func (r UserRole) Validate() bool {
	switch r {
	case RoleAdmin, RoleSecurityAnalyst, RoleComplianceOfficer, RoleRiskManager:
		return true
	default:
		return false
	}
}
