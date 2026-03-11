package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/compliance/evidence"
)

// RegisterEvidenceRoutes adds evidence export endpoints.
func RegisterEvidenceRoutes(app *fiber.App, evSvc *evidence.Service) {
	app.Get("/api/v1/compliance/evidence", func(c *fiber.Ctx) error {
		format := c.Query("format", "json")
		switch format {
		case "pdf":
			data, err := evSvc.ExportPDF(c.Context())
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			c.Set("Content-Type", "application/pdf")
			return c.Send(data)
		case "json":
			data, err := evSvc.ExportJSON(c.Context())
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			c.Set("Content-Type", "application/json")
			return c.Send(data)
		default:
			return c.Status(400).JSON(fiber.Map{"error": "unsupported format"})
		}
	})
}
