package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/asset"
)

// RegisterAssetRoutes adds asset management endpoints.
func RegisterAssetRoutes(app *fiber.App, assetSvc *asset.Service, broadcast func([]byte)) {
	app.Post("/api/v1/assets", func(c *fiber.Ctx) error {
		var payload struct {
			Name        string                `json:"name"`
			IPAddress  string                `json:"ip_address"`
			Owner       string                `json:"owner"`
			Type        asset.AssetType       `json:"type"`
			Criticality asset.AssetCriticality `json:"criticality"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		a, err := assetSvc.CreateAsset(c.Context(), payload.Name, payload.IPAddress, payload.Owner, payload.Type, payload.Criticality)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		assetSvc.BroadcastUpdates(broadcast)
		return c.JSON(a)
	})

	app.Get("/api/v1/assets", func(c *fiber.Ctx) error {
		assets, err := assetSvc.ListAssets(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(assets)
	})

	app.Post("/api/v1/assets/:id/vulnerabilities", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Vulnerabilities []asset.Vulnerability `json:"vulnerabilities"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := assetSvc.ImportVulnerabilities(c.Context(), id, payload.Vulnerabilities); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		assetSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "vulnerabilities imported"})
	})

	app.Get("/api/v1/assets/:id/risk", func(c *fiber.Ctx) error {
		id := c.Params("id")
		score, err := assetSvc.CalculateRiskScore(c.Context(), id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"risk_score": score})
	})
}
