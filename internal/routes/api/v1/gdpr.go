package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/compliance/gdpr"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// RegisterGDPRRoutes adds GDPR compliance endpoints
func RegisterGDPRRoutes(app *fiber.App, gdprSvc *gdpr.Service, broadcast func([]byte)) {
	// GDPR Management
	app.Post("/api/v1/compliance/gdpr", func(c *fiber.Ctx) error {
		var payload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Version     string `json:"version"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		gdpr, err := gdprSvc.CreateGDPR(c.Context(), payload.Name, payload.Description, payload.Version)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(gdpr)
	})

	app.Get("/api/v1/compliance/gdpr", func(c *fiber.Ctx) error {
		gdpr, err := gdprSvc.ListGDPRs(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(gdpr)
	})

	app.Get("/api/v1/compliance/gdpr/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		gdpr, err := gdprSvc.GetGDPR(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(gdpr)
	})

	// Process Management
	app.Post("/api/v1/compliance/gdpr/processes", func(c *fiber.Ctx) error {
		var payload struct {
			GDPRID      string            `json:"gdpr_id"`
			Name        string            `json:"name"`
			Description string            `json:"description"`
			Purpose     string            `json:"purpose"`
			LegalBasis  types.LegalBasis  `json:"legal_basis"`
			Owner       string            `json:"owner"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		process, err := gdprSvc.CreateProcess(c.Context(), payload.GDPRID, payload.Name, payload.Description, payload.Purpose, payload.LegalBasis, payload.Owner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(process)
	})

	app.Put("/api/v1/compliance/gdpr/processes/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status types.ProcessStatus `json:"status"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := gdprSvc.UpdateProcessStatus(c.Context(), id, payload.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "process status updated"})
	})

	// Data Subject Requests
	app.Post("/api/v1/compliance/gdpr/requests", func(c *fiber.Ctx) error {
		var payload struct {
			RightID     string            `json:"right_id"`
			Type        types.RequestType `json:"type"`
			DataSubject string            `json:"data_subject"`
			Identity    string            `json:"identity"`
			Description string            `json:"description"`
			Priority    types.Priority     `json:"priority"`
			AssignedTo  string            `json:"assigned_to"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		request, err := gdprSvc.CreateRequest(c.Context(), payload.RightID, payload.Type, payload.DataSubject, payload.Identity, payload.Description, payload.Priority, payload.AssignedTo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(request)
	})

	app.Put("/api/v1/compliance/gdpr/requests/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status     types.RequestStatus `json:"status"`
			Resolution string              `json:"resolution"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := gdprSvc.UpdateRequestStatus(c.Context(), id, payload.Status, payload.Resolution); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "request status updated"})
	})

	// Data Breach Management
	app.Post("/api/v1/compliance/gdpr/breaches", func(c *fiber.Ctx) error {
		var payload struct {
			GDPRID      string                `json:"gdpr_id"`
			Title       string                `json:"title"`
			Description string                `json:"description"`
			Date        string                `json:"date"`
			Type        types.BreachType      `json:"type"`
			Categories  []types.DataCategory  `json:"categories"`
			Affected    int                   `json:"affected"`
			Cause       string                `json:"cause"`
			Owner       string                `json:"owner"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		breachDate, err := time.Parse(time.RFC3339, payload.Date)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid date format"})
		}
		breach, err := gdprSvc.CreateBreach(c.Context(), payload.GDPRID, payload.Title, payload.Description, breachDate, payload.Type, payload.Categories, payload.Affected, payload.Cause, payload.Owner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(breach)
	})

	app.Put("/api/v1/compliance/gdpr/breaches/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status      types.BreachStatus `json:"status"`
			Remediation string             `json:"remediation"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := gdprSvc.UpdateBreachStatus(c.Context(), id, payload.Status, payload.Remediation); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "breach status updated"})
	})

	// DPIA Management
	app.Post("/api/v1/compliance/gdpr/dpias", func(c *fiber.Ctx) error {
		var payload struct {
			GDPRID      string `json:"gdpr_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			ProcessID   string `json:"process_id"`
			Controller  string `json:"controller"`
			Assessor    string `json:"assessor"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		dpia, err := gdprSvc.CreateDPIA(c.Context(), payload.GDPRID, payload.Title, payload.Description, payload.ProcessID, payload.Controller, payload.Assessor)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(dpia)
	})

	app.Put("/api/v1/compliance/gdpr/dpias/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status   types.DPIAStatus `json:"status"`
			Outcome  string           `json:"outcome"`
			Approval string           `json:"approval"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := gdprSvc.UpdateDPIAStatus(c.Context(), id, payload.Status, payload.Outcome, payload.Approval); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		gdprSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "dpia status updated"})
	})

	// Compliance Score
	app.Get("/api/v1/compliance/gdpr/:id/score", func(c *fiber.Ctx) error {
		id := c.Params("id")
		score, err := gdprSvc.GetComplianceScore(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"compliance_score": score})
	})
}
