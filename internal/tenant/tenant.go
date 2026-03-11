package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Tenant represents a multi-tenant organization.
type Tenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Domain      string    `json:"domain"`
	Plan        PlanType  `json:"plan"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlanType defines subscription tiers.
type PlanType string

const (
	PlanTypeFree    PlanType = "free"
	PlanTypePro     PlanType = "pro"
	PlanTypeEnterprise PlanType = "enterprise"
)

// Service manages tenant isolation.
type Service struct {
	tenants map[string]*Tenant
}

// New creates a Tenant service.
func New() *Service {
	return &Service{
		tenants: make(map[string]*Tenant),
	}
}

// CreateTenant adds a new tenant.
func (s *Service) CreateTenant(ctx context.Context, name, domain string, plan PlanType) (*Tenant, error) {
	now := time.Now().UTC()
	t := &Tenant{
		ID:        uuid.New().String(),
		Name:      name,
		Domain:    domain,
		Plan:      plan,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.tenants[t.ID] = t
	log.Printf("tenant: created tenant %s", t.ID)
	return t, nil
}

// ListTenants returns all tenants.
func (s *Service) ListTenants(ctx context.Context) ([]*Tenant, error) {
	var list []*Tenant
	for _, t := range s.tenants {
		list = append(list, t)
	}
	return list, nil
}

// GetTenantByID retrieves a tenant.
func (s *Service) GetTenantByID(ctx context.Context, id string) (*Tenant, error) {
	t, ok := s.tenants[id]
	if !ok {
		return nil, ErrTenantNotFound
	}
	return t, nil
}

// BroadcastUpdates sends tenant updates via WebSocket.
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":    "tenant_update",
		"tenants": s.tenants,
	})
	broadcast(data)
}

// Errors
var (
	ErrTenantNotFound = fmt.Errorf("tenant not found")
)
