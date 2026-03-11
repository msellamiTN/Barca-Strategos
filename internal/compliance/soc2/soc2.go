package soc2

import (
	"context"
	"log"
	"time"
)

// Service runs SOC2 compliance assessments.
type Service struct{}

// New creates a SOC2 service.
func New() *Service {
	return &Service{}
}

// RunBackgroundMonitoring runs periodic SOC2 checks.
func (s *Service) RunBackgroundMonitoring(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.assess(); err != nil {
				log.Printf("soc2 assessment error: %v", err)
			}
		}
	}
}

func (s *Service) assess() error {
	log.Println("soc2: running background assessment")
	// TODO: evaluate controls, calculate score, store findings
	return nil
}
