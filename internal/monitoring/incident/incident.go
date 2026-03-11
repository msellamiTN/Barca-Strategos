package incident

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Incident represents a security incident.
type Incident struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Severity  string    `json:"severity"`
	Status    string    `json:"status"`
	Summary   string    `json:"summary"`
}

// Service manages incident lifecycle.
type Service struct{}

// New creates an Incident service.
func New() *Service {
	return &Service{}
}

// CreateIncident creates a new incident.
func (s *Service) CreateIncident(ctx context.Context, typ, severity, summary string) (*Incident, error) {
	inc := &Incident{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UTC(),
		Type:      typ,
		Severity:  severity,
		Status:    "open",
		Summary:   summary,
	}
	log.Printf("incident: created incident %s", inc.ID)
	return inc, nil
}
