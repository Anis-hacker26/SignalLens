package app

import (
	"context"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

func TestRunStopsWhenContextIsCancelled(t *testing.T) {
	cfg := config.Config{
		AgentID:    "test-agent-001",
		BackendURL: "http://localhost:3000",
	}

	application := New(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- application.Run(ctx)
	}()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("application did not stop after context cancellation")
	}
}
