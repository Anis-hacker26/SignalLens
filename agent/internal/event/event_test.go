package event

import (
	"testing"
	"time"
)

func TestEventCreation(t *testing.T) {
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

	event := Event{
		ID:            "event-001",
		Type:          TypeProcess,
		Timestamp:     timestamp,
		Host:          "TEST-HOST",
		Source:        "sysmon",
		SourceEventID: "1",
	}

	if event.ID != "event-001" {
		t.Fatalf("expected ID event-001, got %s", event.ID)
	}

	if event.Type != TypeProcess {
		t.Fatalf("expected type %s, got %s", TypeProcess, event.Type)
	}

	if !event.Timestamp.Equal(timestamp) {
		t.Fatalf("expected timestamp %v, got %v", timestamp, event.Timestamp)
	}

	if event.Host != "TEST-HOST" {
		t.Fatalf("expected host TEST-HOST, got %s", event.Host)
	}

	if event.Source != "sysmon" {
		t.Fatalf("expected source sysmon, got %s", event.Source)
	}

	if event.SourceEventID != "1" {
		t.Fatalf("expected source event ID 1, got %s", event.SourceEventID)
	}
}

func TestEventTypes(t *testing.T) {
	tests := []struct {
		name string
		typ  Type
		want string
	}{
		{
			name: "process",
			typ:  TypeProcess,
			want: "process",
		},
		{
			name: "network",
			typ:  TypeNetwork,
			want: "network",
		},
		{
			name: "dns",
			typ:  TypeDNS,
			want: "dns",
		},
		{
			name: "file",
			typ:  TypeFile,
			want: "file",
		},
		{
			name: "powershell",
			typ:  TypePowerShell,
			want: "powershell",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.typ) != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, tt.typ)
			}
		})
	}
}
