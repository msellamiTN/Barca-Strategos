package case

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service manages incident cases.
type Service struct {
	cases map[string]*types.Case
}

// New creates a Case service.
func New() *Service {
	return &Service{
		cases: make(map[string]*types.Case),
	}
}

// CreateCase creates a new incident case.
func (s *Service) CreateCase(ctx context.Context, title, description string, severity types.Severity, assigneeID string) (*types.Case, error) {
	now := time.Now().UTC()
	c := &types.Case{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      types.CaseStatusNew,
		Severity:    severity,
		AssigneeID:  assigneeID,
		CreatedAt:   now,
		UpdatedAt:   now,
		SLA:         s.calculateSLA(severity, now),
	}
	s.cases[c.ID] = c
	log.Printf("case: created case %s", c.ID)
	return c, nil
}

// ListCases returns all cases.
func (s *Service) ListCases(ctx context.Context) ([]*types.Case, error) {
	var list []*types.Case
	for _, c := range s.cases {
		list = append(list, c)
	}
	return list, nil
}

// UpdateCase updates case status or assignee.
func (s *Service) UpdateCase(ctx context.Context, caseID string, status types.CaseStatus, assigneeID string) error {
	c, ok := s.cases[caseID]
	if !ok {
		return ErrCaseNotFound
	}
	if status != "" {
		c.Status = status
		if status == types.CaseStatusResolved {
			now := time.Now().UTC()
			c.ResolvedAt = &now
		}
	}
	if assigneeID != "" {
		c.AssigneeID = assigneeID
	}
	c.UpdatedAt = time.Now().UTC()
	log.Printf("case: updated case %s", caseID)
	return nil
}

// AddAction adds an automated action to a case.
func (s *Service) AddAction(ctx context.Context, caseID string, actionType types.ActionType, details string, executedBy string) (*types.Action, error) {
	c, ok := s.cases[caseID]
	if !ok {
		return nil, ErrCaseNotFound
	}
	action := &types.Action{
		ID:         uuid.New().String(),
		CaseID:     caseID,
		Type:       actionType,
		Status:      "completed",
		ExecutedAt: time.Now().UTC(),
		ExecutedBy: executedBy,
		Details:    details,
	}
	c.UpdatedAt = time.Now().UTC()
	log.Printf("case: added action %s to case %s", action.ID, caseID)
	return action, nil
}

// calculateSLA sets response and resolution due times based on severity.
func (s *Service) calculateSLA(severity types.Severity, now time.Time) *types.SLA {
	var response, resolve time.Duration
	switch severity {
	case types.SeverityCritical:
		response = 15 * time.Minute
		resolve = 4 * time.Hour
	case types.SeverityHigh:
		response = 1 * time.Hour
		resolve = 8 * time.Hour
	case types.SeverityMedium:
		response = 4 * time.Hour
		resolve = 24 * time.Hour
	case types.SeverityLow:
		response = 24 * time.Hour
		resolve = 72 * time.Hour
	}
	return &types.SLA{
		ResponseDue: now.Add(response),
		ResolveDue:  now.Add(resolve),
	}
}

// BroadcastUpdates sends case updates via WebSocket.
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":  "case_update",
		"cases": s.cases,
	})
	broadcast(data)
}

// Errors
var (
	ErrCaseNotFound = fmt.Errorf("case not found")
)
