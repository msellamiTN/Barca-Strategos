package threatintel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/barca-strategos/phoenix/pkg/types"
	"github.com/barca-strategos/phoenix/internal/siem"
)

// Service manages threat intelligence indicators.
type Service struct {
	iocs map[string]*types.ThreatIntel
}

// New creates a ThreatIntel service.
func New() *Service {
	return &Service{
		iocs: make(map[string]*types.ThreatIntel),
	}
}

// AddIOC adds a threat intelligence indicator.
func (s *Service) AddIOC(ctx context.Context, iocType types.ThreatIOCType, value, source, description string, confidence float64) (*types.ThreatIntel, error) {
	ioc := &types.ThreatIntel{
		IOCType:     iocType,
		Value:       value,
		Source:      source,
		Description: description,
		Confidence:  confidence,
		ValidUntil:  nil,
	}
	s.iocs[value] = ioc
	log.Printf("threatintel: added IOC %s (%s)", value, iocType)
	return ioc, nil
}

// Match checks if a value matches any known IOC.
func (s *Service) Match(ctx context.Context, value string) (*types.ThreatIntel, bool) {
	ioc, ok := s.iocs[value]
	if !ok {
		return nil, false
	}
	// Check validity
	if ioc.ValidUntil != nil && time.Now().UTC().After(*ioc.ValidUntil) {
		return nil, false
	}
	return ioc, true
}

// ListIOCs returns all indicators.
func (s *Service) ListIOCs(ctx context.Context) ([]*types.ThreatIntel, error) {
	var list []*types.ThreatIntel
	for _, ioc := range s.iocs {
		if ioc.ValidUntil == nil || time.Now().UTC().Before(*ioc.ValidUntil) {
			list = append(list, ioc)
		}
	}
	return list, nil
}

// EnrichAlert enriches an alert with threat intel matches.
func (s *Service) EnrichAlert(ctx context.Context, alert *siem.Alert) error {
	// Simple enrichment: check if any field matches IOCs
	if ioc, ok := s.Match(ctx, alert.Details["src_ip"].(string)); ok {
		alert.Details["threat_intel"] = ioc
	}
	return nil
}

// BroadcastUpdates sends IOC updates via WebSocket.
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type": "threatintel_update",
		"iocs": s.iocs,
	})
	broadcast(data)
}
