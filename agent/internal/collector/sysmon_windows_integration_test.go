//go:build windows

package collector

import (
	"context"
	"strings"
	"testing"

	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

func TestSysmonCollectorCollectIntegration(t *testing.T) {
	collector := NewSysmonCollector()

	events, err := collector.Collect(context.Background())

t.Logf("collector returned %d events", len(events))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "access is denied") {
			t.Skip("Sysmon integration test requires an elevated Administrator PowerShell")
		}

		t.Fatalf("failed to collect Sysmon events: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("expected at least one Sysmon event")
	}

	for _, event := range events {
		assertValidSysmonRawEvent(t, event)
	}
}

func assertValidSysmonRawEvent(t *testing.T, event telemetry.RawEvent) {
	t.Helper()

	if event.Provider != sysmonProvider {
		t.Fatalf(
			"expected provider %q, got %q",
			sysmonProvider,
			event.Provider,
		)
	}

	if event.EventID != 1 && event.EventID != 3 {
		t.Fatalf(
			"expected Event ID 1 or 3, got %d",
			event.EventID,
		)
	}

	if event.Timestamp.IsZero() {
		t.Fatal("expected event timestamp")
	}

	if event.Computer == "" {
		t.Fatal("expected computer name")
	}

	if event.Message == "" {
		t.Fatal("expected raw event XML")
	}
}
