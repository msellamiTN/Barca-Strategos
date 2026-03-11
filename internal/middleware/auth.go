package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/auth"
)

// JWTAuth creates a middleware that validates JWT tokens.
func JWTAuth(authSvc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
		}
		claims, err := authSvc.ValidateToken(authHeader)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		c.SetUserContext(context.WithValue(c.UserContext(), "claims", claims))
		return c.Next()
	}
}

// RequireRole creates a middleware that checks user role.
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value("claims").(*auth.Claims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthenticated"})
		}
		if string(claims.Role) != requiredRole && string(claims.Role) != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "insufficient role"})
		}
		return c.Next()
	}
}
