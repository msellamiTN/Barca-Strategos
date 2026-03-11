package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/chat"
)

// RegisterChatRoutes adds chat bot endpoints.
func RegisterChatRoutes(app *fiber.App, slackBot *chat.SlackBot) {
	app.Post("/api/v1/chat/slack", func(c *fiber.Ctx) error {
		// Convert Fiber request to http.Request for Slack handler
		req, err := http.NewRequest(c.Method(), c.OriginalURL(), c.Body())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to create request"})
		}
		req.Header = c.GetReqHeaders()
		req.ParseForm()
		// Use a response writer wrapper
		w := &fiberResponseWriter{c: c}
		slackBot.HandleSlashCommand(w, req)
		return nil
	})
}

// fiberResponseWriter wraps Fiber Ctx to implement http.ResponseWriter
type fiberResponseWriter struct {
	c *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	w.c.Write(b)
	return len(b), nil
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.c.Status(statusCode)
}
