package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/threatintel"
)

// RegisterThreatIntelRoutes adds threat intel endpoints.
func RegisterThreatIntelRoutes(app *fiber.App, tiSvc *threatintel.Service, broadcast func([]byte)) {
	app.Post("/api/v1/threatintel/iocs", func(c *fiber.Ctx) error {
		var payload struct {
			Type        threatintel.ThreatIOCType `json:"type"`
			Value       string                    `json:"value"`
			Source      string                    `json:"source"`
			Description string                    `json:"description"`
			Confidence  float64                   `json:"confidence"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		ioc, err := tiSvc.AddIOC(c.Context(), payload.Type, payload.Value, payload.Source, payload.Description, payload.Confidence)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		tiSvc.BroadcastUpdates(broadcast)
		return c.JSON(ioc)
	})

	app.Get("/api/v1/threatintel/iocs", func(c *fiber.Ctx) error {
		iocs, err := tiSvc.ListIOCs(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(iocs)
	})

	app.Post("/api/v1/threatintel/match", func(c *fiber.Ctx) error {
		var payload struct {
			Value string `json:"value"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		ioc, ok := tiSvc.Match(c.Context(), payload.Value)
		if !ok {
			return c.JSON(fiber.Map{"match": false})
		}
		return c.JSON(fiber.Map{"match": true, "ioc": ioc})
	})
}
