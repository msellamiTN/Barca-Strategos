package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/case"
)

// RegisterCaseRoutes adds case management endpoints.
func RegisterCaseRoutes(app *fiber.App, caseSvc *case.Service, broadcast func([]byte)) {
	app.Post("/api/v1/cases", func(c *fiber.Ctx) error {
		var payload struct {
			Title       string            `json:"title"`
			Description string            `json:"description"`
			Severity   case.Severity     `json:"severity"`
			AssigneeID string            `json:"assignee_id"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		cs, err := caseSvc.CreateCase(c.Context(), payload.Title, payload.Description, payload.Severity, payload.AssigneeID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		caseSvc.BroadcastUpdates(broadcast)
		return c.JSON(cs)
	})

	app.Get("/api/v1/cases", func(c *fiber.Ctx) error {
		cases, err := caseSvc.ListCases(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(cases)
	})

	app.Put("/api/v1/cases/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status     string `json:"status,omitempty"`
			AssigneeID string `json:"assignee_id,omitempty"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		var status case.CaseStatus
		if payload.Status != "" {
			status = case.CaseStatus(payload.Status)
		}
		if err := caseSvc.UpdateCase(c.Context(), id, status, payload.AssigneeID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		caseSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "case updated"})
	})

	app.Post("/api/v1/cases/:id/actions", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Type        case.ActionType `json:"type"`
			Details     string          `json:"details"`
			ExecutedBy  string          `json:"executed_by"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		action, err := caseSvc.AddAction(c.Context(), id, payload.Type, payload.Details, payload.ExecutedBy)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		caseSvc.BroadcastUpdates(broadcast)
		return c.JSON(action)
	})
}
