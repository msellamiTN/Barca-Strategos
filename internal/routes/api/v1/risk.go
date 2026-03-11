package v1

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/risk"
)

// RegisterRiskRoutes adds risk management endpoints.
func RegisterRiskRoutes(app *fiber.App, riskSvc *risk.Service, broadcast func([]byte)) {
	app.Post("/api/v1/risks", func(c *fiber.Ctx) error {
		var payload struct {
			Title       string            `json:"title"`
			Description string            `json:"description"`
			Likelihood risk.Likelihood   `json:"likelihood"`
			Impact     risk.Impact       `json:"impact"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		r, err := riskSvc.CreateRisk(c.Context(), payload.Title, payload.Description, payload.Likelihood, payload.Impact)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		riskSvc.BroadcastUpdates(broadcast)
		return c.JSON(r)
	})

	app.Get("/api/v1/risks", func(c *fiber.Ctx) error {
		risks, err := riskSvc.ListRisks(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(risks)
	})

	app.Put("/api/v1/risks/:id/mitigation", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Plan   string `json:"plan"`
			Owner  string `json:"owner"`
			DueDate string `json:"due_date"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		due, err := time.Parse(time.RFC3339, payload.DueDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid due_date"})
		}
		if err := riskSvc.UpdateMitigation(c.Context(), id, payload.Plan, payload.Owner, due); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		riskSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "mitigation updated"})
	})
}
