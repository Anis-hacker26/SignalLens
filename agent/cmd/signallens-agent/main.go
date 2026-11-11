package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

func main() {
	log.Println("SignalLens agent starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf("agent configured: %s", cfg.AgentID)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	<-signalChannel

	log.Println("SignalLens agent shutting down")
}
