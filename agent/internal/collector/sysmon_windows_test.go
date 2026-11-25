//go:build windows

package collector

import (
	"testing"

	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

func TestNewSysmonCollector(t *testing.T) {
	collector := NewSysmonCollector()

	if collector == nil {
		t.Fatal("expected Sysmon collector")
	}
}

func TestNewSysmonCollectorHasDefaultBatchSize(t *testing.T) {
	collector := NewSysmonCollector()

	if collector.maxEvents != 20 {
		t.Fatalf(
			"expected default max events 20, got %d",
			collector.maxEvents,
		)
	}
}

func TestFilterNewEvents(t *testing.T) {
	events := []telemetry.RawEvent{
		{RecordID: 100},
		{RecordID: 99},
		{RecordID: 98},
		{RecordID: 97},
	}

	result := filterNewEvents(events, 98)

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 new events, got %d",
			len(result),
		)
	}

	if result[0].RecordID != 100 {
		t.Fatalf(
			"expected first record ID 100, got %d",
			result[0].RecordID,
		)
	}

	if result[1].RecordID != 99 {
		t.Fatalf(
			"expected second record ID 99, got %d",
			result[1].RecordID,
		)
	}
}
