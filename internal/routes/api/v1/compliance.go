package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/compliance/soc2"
	"github.com/barca-strategos/phoenix/internal/compliance/pci"
)

// RegisterComplianceRoutes adds compliance endpoints.
func RegisterComplianceRoutes(app *fiber.App) {
	soc2Svc := soc2.New()
	pciSvc := pci.New()

	app.Get("/api/v1/compliance/soc2/status", func(c *fiber.Ctx) error {
		// TODO: fetch latest SOC2 assessment result
		return c.JSON(fiber.Map{"framework": "SOC2", "status": "in_progress"})
	})

	app.Post("/api/v1/compliance/pci/scan", func(c *fiber.Ctx) error {
		if err := pciSvc.ScanCardholderEnvironment(c.Context()); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "PCI scan initiated"})
	})
}
