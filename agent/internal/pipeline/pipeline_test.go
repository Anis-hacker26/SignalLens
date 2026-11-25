package pipeline

import (
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

func TestPipelineProcessesSysmonProcessEvent(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Date(2026, 11, 13, 0, 29, 0, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>process</Event>",
	}

	result, err := pipeline.Process(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID == "" {
		t.Fatal("expected evidence ID")
	}

	if result.Event.ID == "" {
		t.Fatal("expected normalized event ID")
	}

	if result.Event.Type != "process" {
		t.Fatalf(
			"expected process event, got %q",
			result.Event.Type,
		)
	}

	if result.Event.SourceEventID != "1" {
		t.Fatalf(
			"expected source event ID 1, got %q",
			result.Event.SourceEventID,
		)
	}

	if result.Source != "Microsoft-Windows-Sysmon" {
		t.Fatalf("unexpected source: %q", result.Source)
	}

	if result.CollectedAt.IsZero() {
		t.Fatal("expected collected timestamp")
	}
}

func TestPipelineProcessesSysmonNetworkEvent(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:   3,
		Timestamp: time.Date(2026, 11, 13, 0, 30, 0, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>network</Event>",
	}

	result, err := pipeline.Process(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID == "" {
		t.Fatal("expected evidence ID")
	}

	if result.Event.ID == "" {
		t.Fatal("expected normalized event ID")
	}

	if result.Event.Type != "network" {
		t.Fatalf(
			"expected network event, got %q",
			result.Event.Type,
		)
	}

	if result.Event.SourceEventID != "3" {
		t.Fatalf(
			"expected source event ID 3, got %q",
			result.Event.SourceEventID,
		)
	}

	if result.Source != "Microsoft-Windows-Sysmon" {
		t.Fatalf("unexpected source: %q", result.Source)
	}

	if result.CollectedAt.IsZero() {
		t.Fatal("expected collected timestamp")
	}
}

func TestPipelineRejectsUnsupportedSysmonEvent(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:   5,
		Timestamp: time.Date(2026, 11, 13, 0, 31, 0, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>unsupported</Event>",
	}

	_, err = pipeline.Process(raw)
	if err == nil {
		t.Fatal("expected error for unsupported Sysmon event")
	}
}

func TestPipelineRejectsInvalidProvider(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Date(2026, 11, 13, 0, 32, 0, 0, time.UTC),
		Provider:  "Unknown-Provider",
		Computer:  "TEST-HOST",
		Message:   "<Event>invalid-provider</Event>",
	}

	_, err = pipeline.Process(raw)
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestPipelineRejectsMissingTimestamp(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:  1,
		Provider: "Microsoft-Windows-Sysmon",
		Computer: "TEST-HOST",
		Message:  "<Event>missing-timestamp</Event>",
	}

	_, err = pipeline.Process(raw)
	if err == nil {
		t.Fatal("expected error for missing timestamp")
	}
}

func TestPipelineRejectsMissingComputer(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Date(2026, 11, 13, 0, 33, 0, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Message:   "<Event>missing-computer</Event>",
	}

	_, err = pipeline.Process(raw)
	if err == nil {
		t.Fatal("expected error for missing computer")
	}
}

func TestPipelineRejectsMissingProvider(t *testing.T) {
	pipeline, err := NewDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Date(2026, 11, 13, 0, 34, 0, 0, time.UTC),
		Computer:  "TEST-HOST",
		Message:   "<Event>missing-provider</Event>",
	}

	_, err = pipeline.Process(raw)
	if err == nil {
		t.Fatal("expected error for missing provider")
	}
}
