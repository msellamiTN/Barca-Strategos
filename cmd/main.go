package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/barca-strategos/phoenix/internal/config"
	"github.com/barca-strategos/phoenix/internal/server"
	"github.com/barca-strategos/phoenix/internal/monitoring/alerting"
	"github.com/barca-strategos/phoenix/internal/monitoring/threat"
	"github.com/barca-strategos/phoenix/internal/compliance/soc2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background services
	go startBackgroundServices(ctx)

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}

func startBackgroundServices(ctx context.Context) {
	alertSvc := alerting.New(config.Config{})
	threatSvc := threat.New()
	soc2Svc := soc2.New()

	go func() {
		if err := alertSvc.Run(ctx); err != nil {
			log.Printf("alerting service error: %v", err)
		}
	}()

	go func() {
		if err := threatSvc.UpdateFeeds(ctx); err != nil {
			log.Printf("threat service error: %v", err)
		}
	}()

	go func() {
		if err := soc2Svc.RunBackgroundMonitoring(ctx); err != nil {
			log.Printf("soc2 service error: %v", err)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("shutdown signal received")
}
