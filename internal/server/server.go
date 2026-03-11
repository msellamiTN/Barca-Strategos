package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/barca-strategos/phoenix/internal/config"
	"github.com/barca-strategos/phoenix/internal/database"
	"github.com/barca-strategos/phoenix/internal/routes"
)

// Server wraps the Fiber app, configuration, and database connections.
type Server struct {
	app    *fiber.App
	cfg    config.Config
	pgDB   *database.DB
	redis  *database.Redis
}

// New initializes a Fiber server with routes, middleware, and database connections.
func New(cfg config.Config) (*Server, error) {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize database connections
	pgDB, err := database.New(cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	if err := pgDB.InitSchema(context.Background()); err != nil {
		log.Printf("warning: schema init failed: %v", err)
	}

	redis, err := database.New(cfg.Redis.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	routes.Register(app)

	return &Server{
		app:   app,
		cfg:   cfg,
		pgDB:  pgDB,
		redis: redis,
	}, nil
}

// Start runs the Fiber server on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.app.Listen(addr)
}

// Close gracefully shuts down database connections.
func (s *Server) Close() error {
	if err := s.pgDB.Close(); err != nil {
		log.Printf("postgres close error: %v", err)
	}
	if err := s.redis.Close(); err != nil {
		log.Printf("redis close error: %v", err)
	}
	return nil
}
