package alerting

import (
	"context"
	"log"
	"time"

	"github.com/barca-strategos/phoenix/internal/config"
)

// Service manages alert evaluation and notifications.
type Service struct {
	cfg config.Config
}

// New creates an Alerting service.
func New(cfg config.Config) *Service {
	return &Service{cfg: cfg}
}

// Run starts the alerting background loop.
func (s *Service) Run(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.evaluate(); err != nil {
				log.Printf("alerting evaluation error: %v", err)
			}
		}
	}
}

func (s *Service) evaluate() error {
	// TODO: fetch events from monitoring, apply rules, send notifications
	log.Println("alerting: evaluating alerts...")
	return nil
}
