package transport

import (
	"context"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
)

type fakeTransport struct {
	sentEvidence *evidence.Evidence
}

func (t *fakeTransport) Send(
	ctx context.Context,
	e evidence.Evidence,
) error {
	t.sentEvidence = &e
	return nil
}

func TestTransportSendsEvidence(t *testing.T) {
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

	expected := evidence.Evidence{
		ID: "evidence-001",
		Event: event.Event{
			ID:            "event-001",
			Type:          event.TypeProcess,
			Timestamp:     timestamp,
			Host:          "TEST-HOST",
			Source:        "sysmon",
			SourceEventID: "1",
		},
		CollectedAt:   timestamp,
		Source:        "sysmon",
		SourceEventID: "1",
	}

	transport := &fakeTransport{}

	err := transport.Send(context.Background(), expected)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if transport.sentEvidence == nil {
		t.Fatal("expected evidence to be sent")
	}

	if transport.sentEvidence.ID != "evidence-001" {
		t.Fatalf(
			"expected evidence ID evidence-001, got %s",
			transport.sentEvidence.ID,
		)
	}

	if transport.sentEvidence.Event.ID != "event-001" {
		t.Fatalf(
			"expected event ID event-001, got %s",
			transport.sentEvidence.Event.ID,
		)
	}

	if transport.sentEvidence.Event.Type != event.TypeProcess {
		t.Fatalf(
			"expected process event, got %s",
			transport.sentEvidence.Event.Type,
		)
	}

	if transport.sentEvidence.Source != "sysmon" {
		t.Fatalf(
			"expected source sysmon, got %s",
			transport.sentEvidence.Source,
		)
	}
}
