package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/siem"
)

// RegisterSIEMRoutes adds SIEM ingestion and alert endpoints.
func RegisterSIEMRoutes(app *fiber.App, siemSvc *siem.Service, broadcast func([]byte)) {
	app.Post("/api/v1/siem/ingest/:source", func(c *fiber.Ctx) error {
		source := c.Params("source")
		if err := siemSvc.IngestAlert(c.Context(), source, c.Body()); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		siemSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "alert ingested"})
	})

	app.Get("/api/v1/siem/alerts", func(c *fiber.Ctx) error {
		alerts, err := siemSvc.ListAlerts(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(alerts)
	})
}
