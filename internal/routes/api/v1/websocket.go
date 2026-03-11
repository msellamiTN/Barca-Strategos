package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/websocket"
)

// RegisterWebSocketRoutes adds WebSocket endpoints.
func RegisterWebSocketRoutes(app *fiber.App, hub *websocket.Hub) {
	app.Get("/ws", hub.UpgradeHandler())
}
