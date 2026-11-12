package app

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
	"github.com/Anis-hacker26/SignalLens/agent/internal/logging"
)

func newTestLogger() logging.Logger {
	return logging.NewStandardLogger(log.Default())
}

func TestRunStopsWhenContextIsCancelled(t *testing.T) {
	cfg := config.Config{
		AgentID:    "test-agent-001",
		BackendURL: "http://localhost:3000",
	}

	logger := newTestLogger()
	application := New(cfg, logger)

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

type testComponent struct {
	started chan struct{}
	stopped chan struct{}
}

func newTestComponent() *testComponent {
	return &testComponent{
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (c *testComponent) Start(ctx context.Context) error {
	close(c.started)
	return nil
}

func (c *testComponent) Stop(ctx context.Context) error {
	close(c.stopped)
	return nil
}

func TestRunStartsAndStopsComponents(t *testing.T) {
	cfg := config.Config{
		AgentID:    "test-agent-001",
		BackendURL: "http://localhost:3000",
	}

	component := newTestComponent()
	logger := newTestLogger()
	application := New(cfg, logger, component)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- application.Run(ctx)
	}()

	select {
	case <-component.started:
	case <-time.After(time.Second):
		t.Fatal("component was not started")
	}

	cancel()

	select {
	case <-component.stopped:
	case <-time.After(time.Second):
		t.Fatal("component was not stopped")
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("application did not stop after context cancellation")
	}
}
