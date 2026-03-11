package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/monitoring/siem"
	"github.com/barca-strategos/phoenix/internal/monitoring/incident"
)

// RegisterMonitoringRoutes adds monitoring endpoints.
func RegisterMonitoringRoutes(app *fiber.App) {
	siemSvc := siem.New()
	incSvc := incident.New()

	app.Post("/api/v1/monitoring/events", func(c *fiber.Ctx) error {
		var raw map[string]interface{}
		if err := c.BodyParser(&raw); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		enriched, err := siemSvc.EnrichSecurityEvent(c.Context(), raw)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(enriched)
	})

	app.Post("/api/v1/monitoring/incidents", func(c *fiber.Ctx) error {
		var payload struct {
			Type     string `json:"type"`
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		inc, err := incSvc.CreateIncident(c.Context(), payload.Type, payload.Severity, payload.Summary)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(inc)
	})
}
