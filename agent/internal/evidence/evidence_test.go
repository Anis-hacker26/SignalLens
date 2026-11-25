package evidence

import (
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

func TestBuilderBuildsEvidence(t *testing.T) {
	expectedCollectedAt := time.Date(
		2026,
		11,
		13,
		0,
		15,
		0,
		0,
		time.UTC,
	)

	builder := &Builder{
		clock: func() time.Time {
			return expectedCollectedAt
		},
	}

	normalized := event.Event{
		ID:            "event-123",
		Type:          event.TypeProcess,
		Timestamp:     time.Date(2026, 11, 13, 0, 10, 0, 0, time.UTC),
		Host:          "TEST-HOST",
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "1",
	}

	result, err := builder.Build(normalized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID == "" {
		t.Fatal("expected evidence ID")
	}

	if result.Event.ID != normalized.ID {
		t.Fatalf("unexpected event ID: %q", result.Event.ID)
	}

	if result.Event.Type != event.TypeProcess {
		t.Fatalf("unexpected event type: %q", result.Event.Type)
	}

	if !result.CollectedAt.Equal(expectedCollectedAt) {
		t.Fatalf(
			"expected collected timestamp %v, got %v",
			expectedCollectedAt,
			result.CollectedAt,
		)
	}

	if result.Source != normalized.Source {
		t.Fatalf("unexpected source: %q", result.Source)
	}

	if result.SourceEventID != normalized.SourceEventID {
		t.Fatalf(
			"unexpected source event ID: %q",
			result.SourceEventID,
		)
	}
}

func TestBuilderProducesDeterministicEvidenceID(t *testing.T) {
	builder := &Builder{
		clock: func() time.Time {
			return time.Now()
		},
	}

	normalized := event.Event{
		ID:            "event-123",
		Type:          event.TypeNetwork,
		Timestamp:     time.Date(2026, 11, 13, 0, 10, 0, 0, time.UTC),
		Host:          "TEST-HOST",
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "3",
	}

	first, err := builder.Build(normalized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := builder.Build(normalized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf(
			"expected deterministic evidence IDs, got %q and %q",
			first.ID,
			second.ID,
		)
	}
}

func TestBuilderRejectsMissingEventID(t *testing.T) {
	builder := NewBuilder()

	normalized := event.Event{
		Type:          event.TypeProcess,
		Timestamp:     time.Now(),
		Host:          "TEST-HOST",
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "1",
	}

	_, err := builder.Build(normalized)
	if err == nil {
		t.Fatal("expected missing event ID error")
	}
}

func TestBuilderRejectsMissingEventType(t *testing.T) {
	builder := NewBuilder()

	normalized := event.Event{
		ID:            "event-123",
		Timestamp:     time.Now(),
		Host:          "TEST-HOST",
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "1",
	}

	_, err := builder.Build(normalized)
	if err == nil {
		t.Fatal("expected missing event type error")
	}
}

func TestBuilderRejectsMissingTimestamp(t *testing.T) {
	builder := NewBuilder()

	normalized := event.Event{
		ID:            "event-123",
		Type:          event.TypeProcess,
		Host:          "TEST-HOST",
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "1",
	}

	_, err := builder.Build(normalized)
	if err == nil {
		t.Fatal("expected missing timestamp error")
	}
}

func TestBuilderRejectsMissingHost(t *testing.T) {
	builder := NewBuilder()

	normalized := event.Event{
		ID:            "event-123",
		Type:          event.TypeProcess,
		Timestamp:     time.Now(),
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "1",
	}

	_, err := builder.Build(normalized)
	if err == nil {
		t.Fatal("expected missing host error")
	}
}

func TestBuilderRejectsMissingSource(t *testing.T) {
	builder := NewBuilder()

	normalized := event.Event{
		ID:            "event-123",
		Type:          event.TypeProcess,
		Timestamp:     time.Now(),
		Host:          "TEST-HOST",
		SourceEventID: "1",
	}

	_, err := builder.Build(normalized)
	if err == nil {
		t.Fatal("expected missing source error")
	}
}

func TestBuilderRejectsMissingSourceEventID(t *testing.T) {
	builder := NewBuilder()

	normalized := event.Event{
		ID:        "event-123",
		Type:      event.TypeProcess,
		Timestamp: time.Now(),
		Host:      "TEST-HOST",
		Source:    "Microsoft-Windows-Sysmon",
	}

	_, err := builder.Build(normalized)
	if err == nil {
		t.Fatal("expected missing source event ID error")
	}
}
