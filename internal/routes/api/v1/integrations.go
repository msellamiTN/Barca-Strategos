package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/integrations"
)

// RegisterIntegrationRoutes adds ticketing integration endpoints.
func RegisterIntegrationRoutes(app *fiber.App, serviceNow *integrations.ServiceNow, jira *integrations.Jira) {
	app.Post("/api/v1/integrations/servicenow/ticket", func(c *fiber.Ctx) error {
		var payload struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Severity    string `json:"severity"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		ticketID, err := serviceNow.CreateTicket(payload.Title, payload.Description, payload.Severity)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"ticket_id": ticketID})
	})

	app.Post("/api/v1/integrations/jira/issue", func(c *fiber.Ctx) error {
		var payload struct {
			Project     string `json:"project"`
			Summary     string `json:"summary"`
			Description string `json:"description"`
			Priority    string `json:"priority"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		issueKey, err := jira.CreateIssue(payload.Project, payload.Summary, payload.Description, payload.Priority)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"issue_key": issueKey})
	})
}
