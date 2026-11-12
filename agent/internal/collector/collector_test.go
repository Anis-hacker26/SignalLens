package collector

import (
	"context"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

type fakeCollector struct {
	events []event.Event
	err    error
}

func (c *fakeCollector) Collect(ctx context.Context) ([]event.Event, error) {
	return c.events, c.err
}

func TestCollectorReturnsEvents(t *testing.T) {
	timestamp := time.Date(
		2026,
		time.November,
		12,
		19,
		30,
		0,
		0,
		time.UTC,
	)

	expected := []event.Event{
		{
			ID:            "event-001",
			Type:          event.TypeProcess,
			Timestamp:     timestamp,
			Host:          "TEST-HOST",
			Source:        "test",
			SourceEventID: "1",
		},
	}

	collector := &fakeCollector{
		events: expected,
	}

	events, err := collector.Collect(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].ID != "event-001" {
		t.Fatalf("expected event ID event-001, got %s", events[0].ID)
	}

	if events[0].Type != event.TypeProcess {
		t.Fatalf("expected process event, got %s", events[0].Type)
	}

	if events[0].Source != "test" {
		t.Fatalf("expected source test, got %s", events[0].Source)
	}
}
