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
	StatePath    string
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

	statePath := os.Getenv("SIGNALLENS_STATE_PATH")

	if statePath == "" {
		var err error

		statePath, err = defaultStateFilePath()
		if err != nil {
			return Config{}, fmt.Errorf(
				"resolve default state path: %w",
				err,
			)
		}
	}

	return Config{
		AgentID:      agentID,
		BackendURL:   backendURL,
		AgentVersion: os.Getenv("SIGNALLENS_AGENT_VERSION"),
		LogLevel:     os.Getenv("SIGNALLENS_LOG_LEVEL"),
		StatePath:    statePath,
	}, nil
}
