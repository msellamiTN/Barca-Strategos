package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/agentic/broker"
	"github.com/barca-strategos/phoenix/internal/agentic/memory"
	"github.com/barca-strategos/phoenix/internal/agentic/referee"
)

// RegisterAgenticRoutes adds agentic AI endpoints.
func RegisterAgenticRoutes(app *fiber.App) {
	brk := broker.New()
	mem := memory.New()
	ref := referee.New()

	// Register a sample tool
	brk.RegisterTool(broker.Tool{
		Name:        "resolve_alert",
		Description: "Resolve a security alert",
		Inputs:      map[string]string{"alert_id": "string"},
		SafetyTier:  "medium",
	})

	app.Post("/api/v1/agent/tools/:name", func(c *fiber.Ctx) error {
		toolName := c.Params("name")
		var inputs map[string]interface{}
		if err := c.BodyParser(&inputs); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		result, err := brk.ExecuteTool(c.Context(), toolName, inputs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	app.Get("/api/v1/agent/memory", func(c *fiber.Ctx) error {
		keys := mem.List(c.Context())
		return c.JSON(fiber.Map{"keys": keys})
	})

	app.Post("/api/v1/agent/referee", func(c *fiber.Ctx) error {
		var payload struct {
			Prompt   string `json:"prompt"`
			Response string `json:"response"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		ok, msg := ref.Review(c.Context(), payload.Prompt, payload.Response)
		return c.JSON(fiber.Map{"allowed": ok, "message": msg})
	})
}
