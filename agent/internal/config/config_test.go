package config

import (
	"testing"
)

func TestLoadRequiresAgentID(t *testing.T) {
	t.Setenv("SIGNALLENS_AGENT_ID", "")
	t.Setenv("SIGNALLENS_BACKEND_URL", "http://localhost:3000")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when SIGNALLENS_AGENT_ID is missing")
	}
}

func TestLoadRequiresBackendURL(t *testing.T) {
	t.Setenv("SIGNALLENS_AGENT_ID", "test-agent")
	t.Setenv("SIGNALLENS_BACKEND_URL", "")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when SIGNALLENS_BACKEND_URL is missing")
	}
}

func TestLoadSuccess(t *testing.T) {
	t.Setenv("SIGNALLENS_AGENT_ID", "test-agent-001")
	t.Setenv("SIGNALLENS_BACKEND_URL", "http://localhost:3000")
	t.Setenv("SIGNALLENS_AGENT_VERSION", "0.1.0")
	t.Setenv("SIGNALLENS_LOG_LEVEL", "info")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.AgentID != "test-agent-001" {
		t.Fatalf("expected AgentID test-agent-001, got: %s", cfg.AgentID)
	}

	if cfg.BackendURL != "http://localhost:3000" {
		t.Fatalf("expected BackendURL http://localhost:3000, got: %s", cfg.BackendURL)
	}

	if cfg.AgentVersion != "0.1.0" {
		t.Fatalf("expected AgentVersion 0.1.0, got: %s", cfg.AgentVersion)
	}

	if cfg.LogLevel != "info" {
		t.Fatalf("expected LogLevel info, got: %s", cfg.LogLevel)
	}
}
