package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/compliance/iso27001"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// RegisterISO27001Routes adds ISO 27001 compliance endpoints
func RegisterISO27001Routes(app *fiber.App, iso27001Svc *iso27001.Service, broadcast func([]byte)) {
	// ISMS Management
	app.Post("/api/v1/compliance/iso27001/isms", func(c *fiber.Ctx) error {
		var payload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Version     string `json:"version"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		isms, err := iso27001Svc.CreateISMS(c.Context(), payload.Name, payload.Description, payload.Version)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(isms)
	})

	app.Get("/api/v1/compliance/iso27001/isms", func(c *fiber.Ctx) error {
		isms, err := iso27001Svc.ListISMS(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(isms)
	})

	app.Get("/api/v1/compliance/iso27001/isms/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		isms, err := iso27001Svc.GetISMS(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(isms)
	})

	app.Put("/api/v1/compliance/iso27001/isms/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Version     string `json:"version"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := iso27001Svc.UpdateISMS(c.Context(), id, payload.Name, payload.Description, payload.Version); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "isms updated"})
	})

	// Controls Management
	app.Post("/api/v1/compliance/iso27001/controls", func(c *fiber.Ctx) error {
		var payload struct {
			ISMSID      string                `json:"isms_id"`
			Number      string                `json:"number"`
			Title       string                `json:"title"`
			Description string                `json:"description"`
			Category    types.ControlCategory `json:"category"`
			Owner       string                `json:"owner"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		control, err := iso27001Svc.CreateControl(c.Context(), payload.ISMSID, payload.Number, payload.Title, payload.Description, payload.Category, payload.Owner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(control)
	})

	app.Put("/api/v1/compliance/iso27001/controls/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status types.ControlStatus `json:"status"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := iso27001Svc.UpdateControlStatus(c.Context(), id, payload.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "control status updated"})
	})

	app.Post("/api/v1/compliance/iso27001/controls/:id/evidence", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Source      string `json:"source"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		evidence := types.Evidence{
			ID:          generateUUID(),
			Title:       payload.Title,
			Description: payload.Description,
			Source:      payload.Source,
			CreatedAt:   time.Now().UTC(),
		}
		if err := iso27001Svc.AddEvidenceToControl(c.Context(), id, evidence); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "evidence added"})
	})

	// Control Testing
	app.Post("/api/v1/compliance/iso27001/tests", func(c *fiber.Ctx) error {
		var payload struct {
			ControlID   string              `json:"control_id"`
			TestType    types.TestType      `json:"test_type"`
			Description string              `json:"description"`
			Procedure   string              `json:"procedure"`
			Tester      string              `json:"tester"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		test, err := iso27001Svc.CreateControlTest(c.Context(), payload.ControlID, payload.TestType, payload.Description, payload.Procedure, payload.Tester)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(test)
	})

	app.Put("/api/v1/compliance/iso27001/tests/:id/result", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Result   types.TestResult `json:"result"`
			Score    float64              `json:"score"`
			Findings []types.Finding   `json:"findings"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := iso27001Svc.UpdateTestResult(c.Context(), id, payload.Result, payload.Score, payload.Findings); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "test result updated"})
	})

	// Risk Management
	app.Post("/api/v1/compliance/iso27001/risks", func(c *fiber.Ctx) error {
		var payload struct {
			ISMSID       string              `json:"isms_id"`
			Title        string              `json:"title"`
			Description  string              `json:"description"`
			Threat       types.Threat        `json:"threat"`
			Vulnerability types.Vulnerability `json:"vulnerability"`
			Asset        types.Asset          `json:"asset"`
			Likelihood   types.Likelihood     `json:"likelihood"`
			Impact       types.Impact         `json:"impact"`
			Treatment    types.RiskTreatment  `json:"treatment"`
			Owner        string              `json:"owner"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		risk, err := iso27001Svc.CreateRisk(c.Context(), payload.ISMSID, payload.Title, payload.Description, payload.Threat, payload.Vulnerability, payload.Asset, payload.Likelihood, payload.Impact, payload.Treatment, payload.Owner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(risk)
	})

	app.Put("/api/v1/compliance/iso27001/risks/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status types.RiskStatus `json:"status"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := iso27001Svc.UpdateRiskStatus(c.Context(), id, payload.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		iso27001Svc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "risk status updated"})
	})

	// Compliance Score
	app.Get("/api/v1/compliance/iso27001/isms/:id/score", func(c *fiber.Ctx) error {
		id := c.Params("id")
		score, err := iso27001Svc.GetComplianceScore(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"compliance_score": score})
	})
}
