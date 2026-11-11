package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anis-hacker26/SignalLens/agent/internal/app"
	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

func main() {
	log.Println("SignalLens agent starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf("agent configured: %s", cfg.AgentID)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application := app.New(cfg)

	if err := application.Run(ctx); err != nil {
		log.Fatalf("agent stopped with error: %v", err)
	}

	log.Println("SignalLens agent shutting down")
}
