package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/tenant"
)

// RegisterTenantRoutes adds multi-tenant management endpoints.
func RegisterTenantRoutes(app *fiber.App, tenantSvc *tenant.Service, broadcast func([]byte)) {
	app.Post("/api/v1/tenants", func(c *fiber.Ctx) error {
		var payload struct {
			Name   string            `json:"name"`
			Domain string            `json:"domain"`
			Plan   tenant.PlanType  `json:"plan"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		t, err := tenantSvc.CreateTenant(c.Context(), payload.Name, payload.Domain, payload.Plan)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		tenantSvc.BroadcastUpdates(broadcast)
		return c.JSON(t)
	})

	app.Get("/api/v1/tenants", func(c *fiber.Ctx) error {
		tenants, err := tenantSvc.ListTenants(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(tenants)
	})

	app.Get("/api/v1/tenants/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		t, err := tenantSvc.GetTenantByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(t)
	})
}
