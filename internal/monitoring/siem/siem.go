package siem

import (
	"context"
	"log"
	"time"
)

// Service integrates with SIEM systems for enrichment.
type Service struct{}

// New creates a SIEM integration service.
func New() *Service {
	return &Service{}
}

// EnrichSecurityEvent enriches a raw security event with SIEM context.
func (s *Service) EnrichSecurityEvent(ctx context.Context, rawEvent map[string]interface{}) (map[string]interface{}, error) {
	// Placeholder: enrich with threat intel, asset data, etc.
	enriched := make(map[string]interface{})
	for k, v := range rawEvent {
		enriched[k] = v
	}
	enriched["siem_enriched_at"] = time.Now().UTC()
	log.Println("siem: enriched security event")
	return enriched, nil
}
