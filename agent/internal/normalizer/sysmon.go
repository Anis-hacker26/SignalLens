package normalizer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

const sysmonProvider = "Microsoft-Windows-Sysmon"

type SysmonNormalizer struct{}

func NewSysmonNormalizer() *SysmonNormalizer {
	return &SysmonNormalizer{}
}

func (n *SysmonNormalizer) Normalize(
	raw telemetry.RawEvent,
) (event.Event, error) {
	if raw.Provider == "" {
		return event.Event{}, fmt.Errorf("raw event provider is required")
	}

	if raw.Provider != sysmonProvider {
		return event.Event{}, fmt.Errorf(
			"unsupported provider: %q",
			raw.Provider,
		)
	}

	if raw.Timestamp.IsZero() {
		return event.Event{}, fmt.Errorf("raw event timestamp is required")
	}

	if raw.Computer == "" {
		return event.Event{}, fmt.Errorf("raw event computer is required")
	}

	eventType, err := normalizeSysmonEventType(raw.EventID)
	if err != nil {
		return event.Event{}, err
	}

	result := event.Event{
		ID:            buildEventID(raw),
		Type:          eventType,
		Timestamp:     raw.Timestamp,
		Host:          raw.Computer,
		Source:        raw.Provider,
		SourceEventID: fmt.Sprintf("%d", raw.EventID),
	}

	switch raw.EventID {
	case 1:
		details, err := parseProcessDetails(raw.Message)
		if err != nil {
			return event.Event{}, fmt.Errorf(
				"parse Sysmon process event: %w",
				err,
			)
		}

		result.Process = &details

	case 3:
		details, err := parseNetworkDetails(raw.Message)
		if err != nil {
			return event.Event{}, fmt.Errorf(
				"parse Sysmon network event: %w",
				err,
			)
		}

		result.Network = &details
	}

	return result, nil
}

func normalizeSysmonEventType(eventID uint32) (event.Type, error) {
	switch eventID {
	case 1:
		return event.TypeProcess, nil

	case 3:
		return event.TypeNetwork, nil

	default:
		return "", fmt.Errorf(
			"unsupported Sysmon Event ID: %d",
			eventID,
		)
	}
}

type sysmonEventXML struct {
	EventData struct {
		Data []sysmonData `xml:"Data"`
	} `xml:"EventData"`
}

type sysmonData struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

func parseSysmonProcessData(xmlData string) (map[string]string, error) {
	var parsed sysmonEventXML

	if err := xml.Unmarshal([]byte(xmlData), &parsed); err != nil {
		return nil, err
	}

	values := make(map[string]string, len(parsed.EventData.Data))

	for _, item := range parsed.EventData.Data {
		values[item.Name] = item.Value
	}

	return values, nil
}

func parseProcessDetails(xmlData string) (event.ProcessDetails, error) {
	values, err := parseSysmonProcessData(xmlData)
	if err != nil {
		return event.ProcessDetails{}, err
	}

	return event.ProcessDetails{
		ProcessGuid:       values["ProcessGuid"],
		ProcessID:         parseUint32(values["ProcessId"]),
		ParentProcessID:   parseUint32(values["ParentProcessId"]),
		ParentProcessGuid: values["ParentProcessGuid"],
		ImagePath:         values["Image"],
		CommandLine:       values["CommandLine"],
		CurrentDirectory:  values["CurrentDirectory"],
		User:              values["User"],
		IntegrityLevel:    values["IntegrityLevel"],
		TerminalSessionID: parseUint32(values["TerminalSessionId"]),
		FileVersion:       values["FileVersion"],
		Description:       values["Description"],
		Product:           values["Product"],
		Company:           values["Company"],
		OriginalFileName:  values["OriginalFileName"],
	}, nil
}

func parseNetworkDetails(xmlData string) (event.NetworkDetails, error) {
	values, err := parseSysmonProcessData(xmlData)
	if err != nil {
		return event.NetworkDetails{}, err
	}

	return event.NetworkDetails{
		ProcessGuid:         values["ProcessGuid"],
		ProcessID:           parseUint32(values["ProcessId"]),
		Protocol:            values["Protocol"],
		Initiated:           parseBool(values["Initiated"]),
		SourceIP:            values["SourceIp"],
		SourcePort:          parseUint32(values["SourcePort"]),
		DestinationIP:       values["DestinationIp"],
		DestinationPort:     parseUint32(values["DestinationPort"]),
		SourceHostname:      values["SourceHostname"],
		DestinationHostname: values["DestinationHostname"],
		SourceIsIPv6:        parseBool(values["SourceIsIpv6"]),
		DestinationIsIPv6:   parseBool(values["DestinationIsIpv6"]),
	}, nil
}

func parseUint32(value string) uint32 {
	if value == "" {
		return 0
	}

	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0
	}

	return uint32(parsed)
}

func parseBool(value string) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return parsed
}

func buildEventID(raw telemetry.RawEvent) string {
	input := fmt.Sprintf(
		"%s|%d|%s|%s|%s",
		raw.Provider,
		raw.EventID,
		raw.Timestamp.UTC().Format(time.RFC3339Nano),
		raw.Computer,
		raw.Message,
	)

	hash := sha256.Sum256([]byte(input))

	return hex.EncodeToString(hash[:])
}
