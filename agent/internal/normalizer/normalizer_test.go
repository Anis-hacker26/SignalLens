package normalizer

import (
	"testing"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

type fakeNormalizer struct {
	event event.Event
	err   error
}

func (n *fakeNormalizer) Normalize(raw any) (event.Event, error) {
	return n.event, n.err
}

func TestNormalizerReturnsEvent(t *testing.T) {
	expected := event.Event{
		ID:            "event-001",
		Type:          event.TypeProcess,
		Host:          "TEST-HOST",
		Source:        "sysmon",
		SourceEventID: "1",
	}

	normalizer := &fakeNormalizer{
		event: expected,
	}

	result, err := normalizer.Normalize("raw-sysmon-event")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.ID != "event-001" {
		t.Fatalf("expected event ID event-001, got %s", result.ID)
	}

	if result.Type != event.TypeProcess {
		t.Fatalf("expected process event, got %s", result.Type)
	}

	if result.Source != "sysmon" {
		t.Fatalf("expected source sysmon, got %s", result.Source)
	}

	if result.SourceEventID != "1" {
		t.Fatalf(
			"expected source event ID 1, got %s",
			result.SourceEventID,
		)
	}
}
