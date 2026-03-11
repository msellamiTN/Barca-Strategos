package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/reporting"
)

// RegisterReportingRoutes adds reporting endpoints.
func RegisterReportingRoutes(app *fiber.App, reportSvc *reporting.Service) {
	app.Get("/api/v1/reports/summary", func(c *fiber.Ctx) error {
		data, err := reportSvc.GenerateExecutiveSummary(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set("Content-Type", "application/json")
		return c.Send(data)
	})

	app.Get("/api/v1/reports/:type/pdf", func(c *fiber.Ctx) error {
		reportType := c.Params("type")
		data, err := reportSvc.GeneratePDFReport(c.Context(), reportType)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", "attachment; filename=report.pdf")
		return c.Send(data)
	})
}
