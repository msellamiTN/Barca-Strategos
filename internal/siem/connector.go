package siem

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/barca-strategos/phoenix/pkg/types"
)

// Alert represents a normalized SIEM alert.
type Alert struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Title     string                 `json:"title"`
	Source    string                 `json:"source"`
	Severity  ThreatSeverity         `json:"severity"`
	Details   map[string]interface{} `json:"details"`
}

// Connector abstracts SIEM integrations.
type Connector interface {
	Ingest(ctx context.Context, raw []byte) ([]*Alert, error)
}

// SplunkConnector integrates with Splunk.
type SplunkConnector struct{}

func NewSplunkConnector() *SplunkConnector {
	return &SplunkConnector{}
}

func (s *SplunkConnector) Ingest(ctx context.Context, raw []byte) ([]*Alert, error) {
	var splunkAlert struct {
		ID        string                 `json:"_id"`
		Time      string                 `json:"_time"`
		Message   string                 `json:"_raw"`
		Severity  string                 `json:"_severity"`
		Source    string                 `json:"_source"`
	}
	if err := json.Unmarshal(raw, &splunkAlert); err != nil {
		return nil, err
	}
	ts, _ := time.Parse(time.RFC3339, splunkAlert.Time)
	severity := ThreatSeverityLow
	switch splunkAlert.Severity {
	case "INFO":
		severity = ThreatSeverityLow
	case "WARN":
		severity = ThreatSeverityMedium
	case "ERROR":
		severity = ThreatSeverityHigh
	case "CRITICAL":
		severity = ThreatSeverityCritical
	}
	return []*Alert{
		{
			ID:        splunkAlert.ID,
			Timestamp: ts,
			Title:     splunkAlert.Message,
			Source:    splunkAlert.Source,
			Severity:  severity,
			Details:   map[string]interface{}{"raw": raw},
		},
	}, nil
}

// Service manages SIEM ingestion and enrichment.
type Service struct {
	connectors map[string]Connector
	alerts     []*Alert
}

// New creates a SIEM service.
func New() *Service {
	return &Service{
		connectors: make(map[string]Connector),
		alerts:     []*Alert{},
	}
}

// RegisterConnector adds a SIEM connector.
func (s *Service) RegisterConnector(name string, conn Connector) {
	s.connectors[name] = conn
	log.Printf("siem: registered connector %s", name)
}

// IngestAlert processes an incoming alert.
func (s *Service) IngestAlert(ctx context.Context, source string, raw []byte) error {
	conn, ok := s.connectors[source]
	if !ok {
		return fmt.Errorf("unknown connector: %s", source)
	}
	alerts, err := conn.Ingest(ctx, raw)
	if err != nil {
		return err
	}
	s.alerts = append(s.alerts, alerts...)
	log.Printf("siem: ingested %d alerts from %s", len(alerts), source)
	return nil
}

// ListAlerts returns recent alerts.
func (s *Service) ListAlerts(ctx context.Context) ([]*Alert, error) {
	return s.alerts, nil
}

// BroadcastUpdates sends alert updates via WebSocket.
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":   "siem_alert",
		"alerts": s.alerts,
	})
	broadcast(data)
}
