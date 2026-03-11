package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/playbook"
)

// RegisterPlaybookRoutes adds playbook engine endpoints.
func RegisterPlaybookRoutes(app *fiber.App, engine *playbook.Engine) {
	app.Post("/api/v1/playbooks", func(c *fiber.Ctx) error {
		var pb playbook.Playbook
		if err := c.BodyParser(&pb); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		engine.RegisterPlaybook(&pb)
		return c.JSON(fiber.Map{"message": "playbook registered"})
	})

	app.Post("/api/v1/playbooks/:id/execute", func(c *fiber.Ctx) error {
		playbookID := c.Params("id")
		var payload struct {
			Context map[string]interface{} `json:"context"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		exec, err := engine.Execute(c.Context(), playbookID, payload.Context)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(exec)
	})
}
