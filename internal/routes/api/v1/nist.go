package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/compliance/nist"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// RegisterNISTRoutes adds NIST CSF compliance endpoints
func RegisterNISTRoutes(app *fiber.App, nistSvc *nist.Service, broadcast func([]byte)) {
	// CSF Management
	app.Post("/api/v1/compliance/nist/csf", func(c *fiber.Ctx) error {
		var payload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Version     string `json:"version"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		csf, err := nistSvc.CreateCSF(c.Context(), payload.Name, payload.Description, payload.Version)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		nistSvc.BroadcastUpdates(broadcast)
		return c.JSON(csf)
	})

	app.Get("/api/v1/compliance/nist/csf", func(c *fiber.Ctx) error {
		csf, err := nistSvc.ListCSFs(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(csf)
	})

	app.Get("/api/v1/compliance/nist/csf/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		csf, err := nistSvc.GetCSF(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(csf)
	})

	// Functions Management
	app.Post("/api/v1/compliance/nist/functions", func(c *fiber.Ctx) error {
		var payload struct {
			CSFID       string              `json:"csf_id"`
			Name        string              `json:"name"`
			Description string              `json:"description"`
			Type        types.FunctionType  `json:"type"`
			Owner       string              `json:"owner"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		function, err := nistSvc.CreateFunction(c.Context(), payload.CSFID, payload.Name, payload.Description, payload.Type, payload.Owner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		nistSvc.BroadcastUpdates(broadcast)
		return c.JSON(function)
	})

	// Subcategories Management
	app.Post("/api/v1/compliance/nist/subcategories", func(c *fiber.Ctx) error {
		var payload struct {
			FunctionID  string              `json:"function_id"`
			Number      string              `json:"number"`
			Title       string              `json:"title"`
			Description string              `json:"description"`
			Priority    types.Priority      `json:"priority"`
			Owner       string              `json:"owner"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		subcategory, err := nistSvc.CreateSubCategory(c.Context(), payload.FunctionID, payload.Number, payload.Title, payload.Description, payload.Priority, payload.Owner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		nistSvc.BroadcastUpdates(broadcast)
		return c.JSON(subcategory)
	})

	app.Put("/api/v1/compliance/nist/subcategories/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Status types.ImplementationStatus `json:"status"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := nistSvc.UpdateSubCategoryStatus(c.Context(), id, payload.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		nistSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "subcategory status updated"})
	})

	// Profiles Management
	app.Post("/api/v1/compliance/nist/profiles", func(c *fiber.Ctx) error {
		var payload struct {
			CSFID       string            `json:"csf_id"`
			Name        string            `json:"name"`
			Description string            `json:"description"`
			Type        types.ProfileType `json:"type"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		profile, err := nistSvc.CreateProfile(c.Context(), payload.CSFID, payload.Name, payload.Description, payload.Type)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		nistSvc.BroadcastUpdates(broadcast)
		return c.JSON(profile)
	})

	// Assessments Management
	app.Post("/api/v1/compliance/nist/assessments", func(c *fiber.Ctx) error {
		var payload struct {
			ProfileID   string                `json:"profile_id"`
			Title       string                `json:"title"`
			Description string                `json:"description"`
			Type        types.AssessmentType  `json:"type"`
			Method      types.AssessmentMethod `json:"method"`
			Scope       string                `json:"scope"`
			Assessor    string                `json:"assessor"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		assessment, err := nistSvc.CreateAssessment(c.Context(), payload.ProfileID, payload.Title, payload.Description, payload.Type, payload.Method, payload.Scope, payload.Assessor)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(assessment)
	})

	app.Put("/api/v1/compliance/nist/assessments/:id/results", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payload struct {
			Results        []types.AssessmentResult `json:"results"`
			Findings       []types.Finding          `json:"findings"`
			Recommendations []types.Recommendation   `json:"recommendations"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := nistSvc.UpdateAssessmentResult(c.Context(), id, payload.Results, payload.Findings, payload.Recommendations); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		nistSvc.BroadcastUpdates(broadcast)
		return c.JSON(fiber.Map{"message": "assessment results updated"})
	})

	// Maturity Score
	app.Get("/api/v1/compliance/nist/csf/:id/maturity", func(c *fiber.Ctx) error {
		id := c.Params("id")
		score, err := nistSvc.GetMaturityScore(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"maturity_score": score})
	})
}
