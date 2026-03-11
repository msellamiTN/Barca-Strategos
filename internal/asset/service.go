package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service manages asset inventory and vulnerabilities.
type Service struct {
	assets map[string]*types.Asset
}

// New creates an Asset service.
func New() *Service {
	return &Service{
		assets: make(map[string]*types.Asset),
	}
}

// CreateAsset adds a new asset.
func (s *Service) CreateAsset(ctx context.Context, name, ip, owner string, assetType types.AssetType, criticality types.AssetCriticality) (*types.Asset, error) {
	now := time.Now().UTC()
	asset := &types.Asset{
		ID:          uuid.New().String(),
		Name:        name,
		Type:        assetType,
		IPAddress:  ip,
		Owner:       owner,
		Criticality:  criticality,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.assets[asset.ID] = asset
	log.Printf("asset: created asset %s", asset.ID)
	return asset, nil
}

// ListAssets returns all assets.
func (s *Service) ListAssets(ctx context.Context) ([]*types.Asset, error) {
	var list []*types.Asset
	for _, a := range s.assets {
		list = append(list, a)
	}
	return list, nil
}

// ImportVulnerabilities adds vulnerability data to an asset.
func (s *Service) ImportVulnerabilities(ctx context.Context, assetID string, vulns []types.Vulnerability) error {
	asset, ok := s.assets[assetID]
	if !ok {
		return ErrAssetNotFound
	}
	asset.Vulnerabilities = append(asset.Vulnerabilities, vulns...)
	asset.UpdatedAt = time.Now().UTC()
	log.Printf("asset: imported %d vulnerabilities to %s", len(vulns), assetID)
	return nil
}

// CalculateRiskScore computes a risk score for an asset.
func (s *Service) CalculateRiskScore(ctx context.Context, assetID string) (float64, error) {
	asset, ok := s.assets[assetID]
	if !ok {
		return 0, ErrAssetNotFound
	}
	critWeight := map[types.AssetCriticality]float64{
		types.AssetCriticalityLow:      1,
		types.AssetCriticalityMedium:   2,
		types.AssetCriticalityHigh:     3,
		types.AssetCriticalityCritical: 4,
	}
	vulnWeight := 0.0
	for _, v := range asset.Vulnerabilities {
		switch v.Severity {
		case "low":
			vulnWeight += 1
		case "medium":
			vulnWeight += 2
		case "high":
			vulnWeight += 3
		case "critical":
			vulnWeight += 4
		}
	}
	score := critWeight[asset.Criticality] * (1 + vulnWeight/10)
	return score, nil
}

// BroadcastUpdates sends asset updates via WebSocket.
func (s *Service) BroadcastUpdates(broadcast func([]byte)) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":   "asset_update",
		"assets": s.assets,
	})
	broadcast(data)
}

// Errors
var (
	ErrAssetNotFound = fmt.Errorf("asset not found")
)
