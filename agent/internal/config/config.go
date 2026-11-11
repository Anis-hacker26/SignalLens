package config

import (
	"fmt"
	"os"
)

type Config struct {
	AgentID      string
	BackendURL   string
	AgentVersion string
	LogLevel     string
}

func Load() (Config, error) {
	agentID := os.Getenv("SIGNALLENS_AGENT_ID")
	backendURL := os.Getenv("SIGNALLENS_BACKEND_URL")

	if agentID == "" {
		return Config{}, fmt.Errorf("SIGNALLENS_AGENT_ID is required")
	}

	if backendURL == "" {
		return Config{}, fmt.Errorf("SIGNALLENS_BACKEND_URL is required")
	}

	return Config{
		AgentID:      agentID,
		BackendURL:   backendURL,
		AgentVersion: os.Getenv("SIGNALLENS_AGENT_VERSION"),
		LogLevel:     os.Getenv("SIGNALLENS_LOG_LEVEL"),
	}, nil
}
