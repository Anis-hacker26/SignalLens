package collector

import (
	"context"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type fakeCollector struct {
	events []telemetry.RawEvent
	err    error
}

func (c *fakeCollector) Collect(
	ctx context.Context,
) ([]telemetry.RawEvent, error) {
	return c.events, c.err
}

func TestCollectorReturnsRawEvents(t *testing.T) {
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

	expected := []telemetry.RawEvent{
		{
			EventID:   1,
			Timestamp: timestamp,
			Provider:  "Microsoft-Windows-Sysmon",
			Computer:  "TEST-HOST",
			Message:   "Process Create:",
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

	if events[0].EventID != 1 {
		t.Fatalf("expected event ID 1, got %d", events[0].EventID)
	}

	if events[0].Provider != "Microsoft-Windows-Sysmon" {
		t.Fatalf(
			"expected Sysmon provider, got %s",
			events[0].Provider,
		)
	}

	if events[0].Computer != "TEST-HOST" {
		t.Fatalf(
			"expected computer TEST-HOST, got %s",
			events[0].Computer,
		)
	}
}
