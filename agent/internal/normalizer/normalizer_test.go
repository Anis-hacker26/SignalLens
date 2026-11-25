package normalizer

import (
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

func TestSysmonNormalizerProcessEvent(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	timestamp := time.Date(
		2026,
		11,
		13,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: timestamp,
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>process event</Event>",
	}

	normalized, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalized.Type != event.TypeProcess {
		t.Fatalf(
			"expected event type %q, got %q",
			event.TypeProcess,
			normalized.Type,
		)
	}

	if !normalized.Timestamp.Equal(timestamp) {
		t.Fatalf("unexpected timestamp: %v", normalized.Timestamp)
	}

	if normalized.Host != "TEST-HOST" {
		t.Fatalf("unexpected host: %q", normalized.Host)
	}

	if normalized.Source != "Microsoft-Windows-Sysmon" {
		t.Fatalf("unexpected source: %q", normalized.Source)
	}

	if normalized.SourceEventID != "1" {
		t.Fatalf("unexpected source event ID: %q", normalized.SourceEventID)
	}

	if normalized.ID == "" {
		t.Fatal("expected normalized event ID")
	}
}

func TestSysmonNormalizerNetworkEvent(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	raw := telemetry.RawEvent{
		EventID:   3,
		Timestamp: time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>network event</Event>",
	}

	normalized, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalized.Type != event.TypeNetwork {
		t.Fatalf(
			"expected event type %q, got %q",
			event.TypeNetwork,
			normalized.Type,
		)
	}

	if normalized.SourceEventID != "3" {
		t.Fatalf("unexpected source event ID: %q", normalized.SourceEventID)
	}
}

func TestSysmonNormalizerRejectsUnsupportedEventID(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	raw := telemetry.RawEvent{
		EventID:   22,
		Timestamp: time.Now(),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>DNS event</Event>",
	}

	_, err := normalizer.Normalize(raw)
	if err == nil {
		t.Fatal("expected unsupported event ID error")
	}
}

func TestSysmonNormalizerRejectsWrongProvider(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Now(),
		Provider:  "Other-Provider",
		Computer:  "TEST-HOST",
		Message:   "<Event>event</Event>",
	}

	_, err := normalizer.Normalize(raw)
	if err == nil {
		t.Fatal("expected provider validation error")
	}
}

func TestSysmonNormalizerRejectsMissingTimestamp(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	raw := telemetry.RawEvent{
		EventID:  1,
		Provider: "Microsoft-Windows-Sysmon",
		Computer: "TEST-HOST",
		Message:  "<Event>event</Event>",
	}

	_, err := normalizer.Normalize(raw)
	if err == nil {
		t.Fatal("expected timestamp validation error")
	}
}

func TestSysmonNormalizerRejectsMissingComputer(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Now(),
		Provider:  "Microsoft-Windows-Sysmon",
		Message:   "<Event>event</Event>",
	}

	_, err := normalizer.Normalize(raw)
	if err == nil {
		t.Fatal("expected computer validation error")
	}
}

func TestSysmonNormalizerProducesDeterministicID(t *testing.T) {
	normalizer := NewSysmonNormalizer()

	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message:   "<Event>same event</Event>",
	}

	first, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf(
			"expected deterministic IDs, got %q and %q",
			first.ID,
			second.ID,
		)
	}
}
func TestSysmonNormalizerNormalizesProcessDetails(t *testing.T) {
	raw := telemetry.RawEvent{
		EventID:   1,
		Timestamp: time.Date(2026, 11, 13, 10, 20, 30, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message: `<Event>
			<EventData>
				<Data Name="ProcessGuid">{12345678-1234-1234-1234-123456789abc}</Data>
				<Data Name="ProcessId">1234</Data>
				<Data Name="Image">C:\Windows\System32\cmd.exe</Data>
				<Data Name="CommandLine">cmd.exe /c whoami</Data>
				<Data Name="CurrentDirectory">C:\Users\Test</Data>
				<Data Name="User">TEST-HOST\TestUser</Data>
				<Data Name="IntegrityLevel">Medium</Data>
				<Data Name="TerminalSessionId">1</Data>
				<Data Name="ParentProcessId">5678</Data>
				<Data Name="ParentProcessGuid">{87654321-4321-4321-4321-abcdef123456}</Data>
				<Data Name="FileVersion">10.0.1</Data>
				<Data Name="Description">Windows Command Processor</Data>
				<Data Name="Product">Microsoft Windows</Data>
				<Data Name="Company">Microsoft Corporation</Data>
				<Data Name="OriginalFileName">cmd.exe</Data>
			</EventData>
		</Event>`,
	}

	normalizer := NewSysmonNormalizer()

	result, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Type != event.TypeProcess {
		t.Fatalf("expected process event, got: %s", result.Type)
	}

	if result.Process == nil {
		t.Fatal("expected process details")
	}

	if result.Process.ProcessGuid != "{12345678-1234-1234-1234-123456789abc}" {
		t.Fatalf("unexpected ProcessGuid: %s", result.Process.ProcessGuid)
	}

	if result.Process.ProcessID != 1234 {
		t.Fatalf("unexpected ProcessID: %d", result.Process.ProcessID)
	}

	if result.Process.ParentProcessID != 5678 {
		t.Fatalf("unexpected ParentProcessID: %d", result.Process.ParentProcessID)
	}

	if result.Process.ImagePath != `C:\Windows\System32\cmd.exe` {
		t.Fatalf("unexpected ImagePath: %s", result.Process.ImagePath)
	}

	if result.Process.CommandLine != `cmd.exe /c whoami` {
		t.Fatalf("unexpected CommandLine: %s", result.Process.CommandLine)
	}

	if result.Process.User != `TEST-HOST\TestUser` {
		t.Fatalf("unexpected User: %s", result.Process.User)
	}

	if result.Process.IntegrityLevel != "Medium" {
		t.Fatalf("unexpected IntegrityLevel: %s", result.Process.IntegrityLevel)
	}

	if result.Network != nil {
		t.Fatal("expected network details to be nil")
	}
}

func TestSysmonNormalizerNormalizesNetworkDetails(t *testing.T) {
	raw := telemetry.RawEvent{
		EventID:   3,
		Timestamp: time.Date(2026, 11, 13, 10, 20, 30, 0, time.UTC),
		Provider:  "Microsoft-Windows-Sysmon",
		Computer:  "TEST-HOST",
		Message: `<Event>
			<EventData>
				<Data Name="ProcessGuid">{12345678-1234-1234-1234-123456789abc}</Data>
				<Data Name="ProcessId">4321</Data>
				<Data Name="Protocol">tcp</Data>
				<Data Name="Initiated">true</Data>
				<Data Name="SourceIp">10.0.0.10</Data>
				<Data Name="SourcePort">54321</Data>
				<Data Name="DestinationIp">13.107.5.93</Data>
				<Data Name="DestinationPort">443</Data>
				<Data Name="SourceHostname">TEST-HOST</Data>
				<Data Name="DestinationHostname">example.com</Data>
				<Data Name="SourceIsIpv6">false</Data>
				<Data Name="DestinationIsIpv6">false</Data>
			</EventData>
		</Event>`,
	}

	normalizer := NewSysmonNormalizer()

	result, err := normalizer.Normalize(raw)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Type != event.TypeNetwork {
		t.Fatalf("expected network event, got: %s", result.Type)
	}

	if result.Network == nil {
		t.Fatal("expected network details")
	}

	if result.Network.ProcessGuid != "{12345678-1234-1234-1234-123456789abc}" {
		t.Fatalf("unexpected ProcessGuid: %s", result.Network.ProcessGuid)
	}

	if result.Network.ProcessID != 4321 {
		t.Fatalf("unexpected ProcessID: %d", result.Network.ProcessID)
	}

	if result.Network.Protocol != "tcp" {
		t.Fatalf("unexpected protocol: %s", result.Network.Protocol)
	}

	if !result.Network.Initiated {
		t.Fatal("expected Initiated to be true")
	}

	if result.Network.SourceIP != "10.0.0.10" {
		t.Fatalf("unexpected SourceIP: %s", result.Network.SourceIP)
	}

	if result.Network.SourcePort != 54321 {
		t.Fatalf("unexpected SourcePort: %d", result.Network.SourcePort)
	}

	if result.Network.DestinationIP != "13.107.5.93" {
		t.Fatalf("unexpected DestinationIP: %s", result.Network.DestinationIP)
	}

	if result.Network.DestinationPort != 443 {
		t.Fatalf("unexpected DestinationPort: %d", result.Network.DestinationPort)
	}

	if result.Process != nil {
		t.Fatal("expected process details to be nil")
	}
}
