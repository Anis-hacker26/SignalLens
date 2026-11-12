package telemetry

import (
	"testing"
	"time"
)

func TestRawEventCreation(t *testing.T) {
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

	raw := RawEvent{
		EventID:   3,
		Timestamp: timestamp,
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "Network connection detected:",
	}

	if raw.EventID != 3 {
		t.Fatalf("expected event ID 3, got %d", raw.EventID)
	}

	if !raw.Timestamp.Equal(timestamp) {
		t.Fatalf(
			"expected timestamp %v, got %v",
			timestamp,
			raw.Timestamp,
		)
	}

	if raw.Provider != "Microsoft-Windows-Sysmon" {
		t.Fatalf(
			"expected provider Microsoft-Windows-Sysmon, got %s",
			raw.Provider,
		)
	}

	if raw.Computer != "TEST-HOST" {
		t.Fatalf(
			"expected computer TEST-HOST, got %s",
			raw.Computer,
		)
	}

	if raw.Message == "" {
		t.Fatal("expected raw event message")
	}
}
