package evidence

import (
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

func TestEvidenceCreation(t *testing.T) {
	eventTimestamp := time.Date(
		2026,
		time.November,
		12,
		19,
		30,
		0,
		0,
		time.UTC,
	)

	collectedAt := time.Date(
		2026,
		time.November,
		12,
		19,
		30,
		1,
		0,
		time.UTC,
	)

	expectedEvent := event.Event{
		ID:            "event-001",
		Type:          event.TypeProcess,
		Timestamp:     eventTimestamp,
		Host:          "TEST-HOST",
		Source:        "sysmon",
		SourceEventID: "1",
	}

	evidence := Evidence{
		ID:            "evidence-001",
		Event:         expectedEvent,
		CollectedAt:   collectedAt,
		Source:        "sysmon",
		SourceEventID: "1",
	}

	if evidence.ID != "evidence-001" {
		t.Fatalf("expected evidence ID evidence-001, got %s", evidence.ID)
	}

	if evidence.Event.ID != "event-001" {
		t.Fatalf("expected event ID event-001, got %s", evidence.Event.ID)
	}

	if evidence.Event.Type != event.TypeProcess {
		t.Fatalf("expected process event, got %s", evidence.Event.Type)
	}

	if !evidence.Event.Timestamp.Equal(eventTimestamp) {
		t.Fatalf(
			"expected event timestamp %v, got %v",
			eventTimestamp,
			evidence.Event.Timestamp,
		)
	}

	if !evidence.CollectedAt.Equal(collectedAt) {
		t.Fatalf(
			"expected collected timestamp %v, got %v",
			collectedAt,
			evidence.CollectedAt,
		)
	}

	if evidence.Source != "sysmon" {
		t.Fatalf("expected source sysmon, got %s", evidence.Source)
	}

	if evidence.SourceEventID != "1" {
		t.Fatalf(
			"expected source event ID 1, got %s",
			evidence.SourceEventID,
		)
	}
}
