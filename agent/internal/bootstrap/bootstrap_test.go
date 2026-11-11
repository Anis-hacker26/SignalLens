package bootstrap

import (
	"testing"

	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

func TestNewApplication(t *testing.T) {
	cfg := config.Config{
		AgentID:    "test-agent-001",
		BackendURL: "http://localhost:3000",
	}

	application := NewApplication(cfg)

	if application == nil {
		t.Fatal("expected application to be created")
	}
}
