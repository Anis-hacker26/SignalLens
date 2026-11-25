package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
)

func TestHTTPTransportSendProcessEvidence(t *testing.T) {
	var received map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/evidence" {
			t.Fatalf("expected /api/evidence, got %s", r.URL.Path)
		}

		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			t.Fatalf("expected application/json, got %s", contentType)
		}

		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(HTTPTransportConfig{
		BackendURL:   server.URL,
		AgentID:      "agent-test-001",
		AgentVersion: "0.1.0",
		OSName:       "Windows",
		OSVersion:    "10.0.26200",
	})

	if err != nil {
		t.Fatalf("create transport: %v", err)
	}

	observedAt := time.Date(
		2026,
		11,
		13,
		10,
		30,
		0,
		0,
		time.UTC,
	)

	collectedAt := observedAt.Add(2 * time.Second)

	item := evidence.Evidence{
		ID: "evidence-process-001",
		Event: event.Event{
			ID:            "event-process-001",
			Type:          event.TypeProcess,
			Timestamp:     observedAt,
			Host:          "SIGNALLENS-TEST",
			Source:        "Microsoft-Windows-Sysmon",
			SourceEventID: "1",
			Process: &event.ProcessDetails{
				ProcessGuid:       "{process-guid}",
				ProcessID:         1234,
				ParentProcessID:   567,
				ParentProcessGuid: "{parent-guid}",
				ImagePath:         `C:\Windows\System32\cmd.exe`,
				CommandLine:       `cmd.exe /c whoami`,
				CurrentDirectory:  `C:\Windows\System32`,
				User:              `TEST\User`,
				IntegrityLevel:    "High",
				TerminalSessionID: 1,
				FileVersion:       "10.0.1",
				Description:       "Windows Command Processor",
				Product:           "Microsoft Windows",
				Company:           "Microsoft Corporation",
				OriginalFileName:  "cmd.exe",
			},
		},
		CollectedAt:   collectedAt,
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "1",
	}

	if err := transport.Send(context.Background(), item); err != nil {
		t.Fatalf("send evidence: %v", err)
	}

	assertJSONField(t, received, "agentId", "agent-test-001")
	assertJSONField(t, received, "hostname", "SIGNALLENS-TEST")
	assertJSONField(t, received, "osName", "Windows")
	assertJSONField(t, received, "osVersion", "10.0.26200")
	assertJSONField(t, received, "agentVersion", "0.1.0")

	assertJSONField(t, received, "type", "PROCESS")
	assertJSONField(t, received, "source", "Microsoft-Windows-Sysmon")
	assertJSONField(t, received, "sourceEventId", "1")

	assertJSONField(t, received, "processGuid", "{process-guid}")
	assertJSONField(t, received, "processId", float64(1234))
	assertJSONField(t, received, "parentProcessId", float64(567))
	assertJSONField(t, received, "imagePath", `C:\Windows\System32\cmd.exe`)
	assertJSONField(t, received, "commandLine", `cmd.exe /c whoami`)
	assertJSONField(t, received, "user", `TEST\User`)
	assertJSONField(t, received, "integrityLevel", "High")
}

func TestHTTPTransportSendNetworkEvidence(t *testing.T) {
	var received map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/evidence" {
			t.Fatalf("expected /api/evidence, got %s", r.URL.Path)
		}

		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(HTTPTransportConfig{
		BackendURL: server.URL,
		AgentID:    "agent-test-002",
		OSName:     "Windows",
		OSVersion:  "10.0.26200",
	})

	if err != nil {
		t.Fatalf("create transport: %v", err)
	}

	item := evidence.Evidence{
		ID: "evidence-network-001",
		Event: event.Event{
			ID:            "event-network-001",
			Type:          event.TypeNetwork,
			Timestamp:     time.Now().UTC(),
			Host:          "SIGNALLENS-TEST",
			Source:        "Microsoft-Windows-Sysmon",
			SourceEventID: "3",
			Network: &event.NetworkDetails{
				ProcessGuid:         "{process-guid}",
				ProcessID:           9236,
				Protocol:            "tcp",
				Initiated:           true,
				SourceIP:            "10.37.16.143",
				SourcePort:          62793,
				DestinationIP:       "13.107.5.93",
				DestinationPort:     443,
				SourceHostname:      "SIGNALLENS-TEST",
				DestinationHostname: "",
				SourceIsIPv6:        false,
				DestinationIsIPv6:   false,
			},
		},
		CollectedAt:   time.Now().UTC(),
		Source:        "Microsoft-Windows-Sysmon",
		SourceEventID: "3",
	}

	if err := transport.Send(context.Background(), item); err != nil {
		t.Fatalf("send evidence: %v", err)
	}

	assertJSONField(t, received, "agentId", "agent-test-002")
	assertJSONField(t, received, "type", "NETWORK")
	assertJSONField(t, received, "sourceEventId", "3")
	assertJSONField(t, received, "processGuid", "{process-guid}")
	assertJSONField(t, received, "processId", float64(9236))
	assertJSONField(t, received, "protocol", "tcp")
	assertJSONField(t, received, "initiated", true)
	assertJSONField(t, received, "sourceIp", "10.37.16.143")
	assertJSONField(t, received, "sourcePort", float64(62793))
	assertJSONField(t, received, "destinationIp", "13.107.5.93")
	assertJSONField(t, received, "destinationPort", float64(443))
}

func TestHTTPTransportSendReturnsErrorForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "backend failure", http.StatusInternalServerError)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(HTTPTransportConfig{
		BackendURL: server.URL,
		AgentID:    "agent-test-003",
	})

	if err != nil {
		t.Fatalf("create transport: %v", err)
	}

	item := evidence.Evidence{
		ID: "evidence-error-001",
		Event: event.Event{
			ID:        "event-error-001",
			Type:      event.TypeProcess,
			Timestamp: time.Now().UTC(),
			Host:      "SIGNALLENS-TEST",
			Source:    "Microsoft-Windows-Sysmon",
			Process:   &event.ProcessDetails{},
		},
		CollectedAt: time.Now().UTC(),
		Source:      "Microsoft-Windows-Sysmon",
	}

	err = transport.Send(context.Background(), item)

	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}

	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Fatalf("expected HTTP 500 error, got %v", err)
	}
}

func assertJSONField(
	t *testing.T,
	payload map[string]any,
	field string,
	expected any,
) {
	t.Helper()

	actual, ok := payload[field]
	if !ok {
		t.Fatalf("expected JSON field %q to exist", field)
	}

	if actual != expected {
		t.Fatalf(
			"field %q: expected %v, got %v",
			field,
			expected,
			actual,
		)
	}
}
