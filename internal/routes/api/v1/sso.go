package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/auth"
)

// RegisterSSORoutes adds SSO endpoints.
func RegisterSSORoutes(app *fiber.App, sso *auth.SSOHandler) {
	app.Get("/auth/saml", func(c *fiber.Ctx) error {
		// Convert Fiber request to http.Request for SSO handler
		req, _ := http.NewRequest(c.Method(), c.OriginalURL(), nil)
		req.Header = c.GetReqHeaders()
		w := &fiberResponseWriter{c: c}
		sso.InitiateSAML(w, req)
		return nil
	})

	app.Post("/auth/saml/callback", func(c *fiber.Ctx) error {
		req, _ := http.NewRequest(c.Method(), c.OriginalURL(), c.Body())
		req.Header = c.GetReqHeaders()
		req.ParseForm()
		w := &fiberResponseWriter{c: c}
		sso.ConsumeSAML(w, req)
		return nil
	})

	app.Get("/auth/oidc", func(c *fiber.Ctx) error {
		req, _ := http.NewRequest(c.Method(), c.OriginalURL(), nil)
		req.Header = c.GetReqHeaders()
		w := &fiberResponseWriter{c: c}
		sso.InitiateOIDC(w, req)
		return nil
	})

	app.Get("/auth/oidc/callback", func(c *fiber.Ctx) error {
		req, _ := http.NewRequest(c.Method(), c.OriginalURL(), nil)
		req.Header = c.GetReqHeaders()
		w := &fiberResponseWriter{c: c}
		sso.OIDCCallback(w, req)
		return nil
	})
}
