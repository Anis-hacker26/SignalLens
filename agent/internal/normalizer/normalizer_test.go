package normalizer

import (
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type fakeNormalizer struct {
	event event.Event
	err   error
}

func (n *fakeNormalizer) Normalize(
	raw telemetry.RawEvent,
) (event.Event, error) {
	return n.event, n.err
}

func TestNormalizerReturnsEvent(t *testing.T) {
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

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: timestamp,
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "Process Create:",
	}

	expected := event.Event{
		ID:            "event-001",
		Type:          event.TypeProcess,
		Timestamp:     timestamp,
		Host:          "TEST-HOST",
		Source:        "sysmon",
		SourceEventID: "1",
	}

	normalizer := &fakeNormalizer{
		event: expected,
	}

	result, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.ID != "event-001" {
		t.Fatalf("expected event ID event-001, got %s", result.ID)
	}

	if result.Type != event.TypeProcess {
		t.Fatalf(
			"expected process event, got %s",
			result.Type,
		)
	}

	if result.Source != "sysmon" {
		t.Fatalf(
			"expected source sysmon, got %s",
			result.Source,
		)
	}

	if result.SourceEventID != "1" {
		t.Fatalf(
			"expected source event ID 1, got %s",
			result.SourceEventID,
		)
	}
}
