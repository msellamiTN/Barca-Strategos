package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service manages risk register and scoring.
type Service struct {
	risks map[string]*types.Risk
}

// New creates a Risk service.
func New() *Service {
	return &Service{
		risks: make(map[string]*types.Risk),
	}
}

// CreateRisk adds a new risk and calculates its score.
func (s *Service) CreateRisk(ctx context.Context, title, description string, likelihood types.Likelihood, impact types.Impact) (*types.Risk, error) {
	r := &types.Risk{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Likelihood:  likelihood,
		Impact:      impact,
		Score:       s.calculateScore(likelihood, impact),
		Status:      types.RiskStatusOpen,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	s.risks[r.ID] = r
	log.Printf("risk: created risk %s with score %.2f", r.ID, r.Score)
	return r, nil
}

// ListRisks returns all risks.
func (s *Service) ListRisks(ctx context.Context) ([]*types.Risk, error) {
	var list []*types.Risk
	for _, r := range s.risks {
		list = append(list, r)
	}
	return list, nil
}

// UpdateMitigation sets or updates mitigation for a risk.
func (s *Service) UpdateMitigation(ctx context.Context, riskID, plan, owner string, dueDate time.Time) error {
	r, ok := s.risks[riskID]
	if !ok {
		return fmt.Errorf("risk not found")
	}
	r.Mitigation = &types.Mitigation{
		Plan:    plan,
		Owner:   owner,
		DueDate: dueDate,
		Status:  "planned",
	}
	r.Status = types.RiskStatusInProgress
	r.UpdatedAt = time.Now().UTC()
	log.Printf("risk: mitigation updated for %s", riskID)
	return nil
}

// calculateScore maps likelihood × impact to 1–25.
func (s *Service) calculateScore(l types.Likelihood, i types.Impact) float64 {
	lMap := map[types.Likelihood]float64{
		types.LikelihoodRare:     1,
		types.LikelihoodUnlikely: 2,
		types.LikelihoodPossible: 3,
		types.LikelihoodLikely:   4,
		types.LikelihoodCertain:  5,
	}
	iMap := map[types.Impact]float64{
		types.ImpactNegligible:     1,
		types.ImpactMinor:          2,
		types.ImpactModerate:       3,
		types.ImpactMajor:          4,
		types.ImpactCatastrophic:   5,
	}
	score := lMap[l] * iMap[i]
	return math.Min(score, 25)
}

// BroadcastUpdates sends risk updates via WebSocket.
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":  "risk_update",
		"risks": s.risks,
	})
	broadcast(data)
}
