package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
)

type HTTPTransport struct {
	client       *http.Client
	backendURL   string
	agentID      string
	agentVersion string
	osName       string
	osVersion    string
}

type HTTPTransportConfig struct {
	BackendURL   string
	AgentID      string
	AgentVersion string
	OSName       string
	OSVersion    string
	Client       *http.Client
}

func NewHTTPTransport(cfg HTTPTransportConfig) (*HTTPTransport, error) {
	if strings.TrimSpace(cfg.BackendURL) == "" {
		return nil, fmt.Errorf("backend URL is required")
	}

	if strings.TrimSpace(cfg.AgentID) == "" {
		return nil, fmt.Errorf("agent ID is required")
	}

	client := cfg.Client
	if client == nil {
		client = &http.Client{
			Timeout: 3 * time.Second,
		}
	}

	return &HTTPTransport{
		client:       client,
		backendURL:   strings.TrimRight(cfg.BackendURL, "/"),
		agentID:      cfg.AgentID,
		agentVersion: cfg.AgentVersion,
		osName:       cfg.OSName,
		osVersion:    cfg.OSVersion,
	}, nil
}

type ingestEvidenceRequest struct {
	AgentID      string `json:"agentId"`
	Hostname     string `json:"hostname"`
	OSName       string `json:"osName"`
	OSVersion    string `json:"osVersion"`
	AgentVersion string `json:"agentVersion"`

	Type          string  `json:"type"`
	Source        string  `json:"source"`
	SourceEventID string  `json:"sourceEventId"`
	ObservedAt    string  `json:"observedAt"`
	CollectedAt   string  `json:"collectedAt"`
	CorrelationID *string `json:"correlationId,omitempty"`

	Payload any `json:"payload"`

	ProcessGuid       string `json:"processGuid,omitempty"`
	ProcessID         uint32 `json:"processId,omitempty"`
	ParentProcessID   uint32 `json:"parentProcessId,omitempty"`
	ParentProcessGuid string `json:"parentProcessGuid,omitempty"`
	ImagePath         string `json:"imagePath,omitempty"`
	CommandLine       string `json:"commandLine,omitempty"`
	CurrentDirectory  string `json:"currentDirectory,omitempty"`
	User              string `json:"user,omitempty"`
	IntegrityLevel    string `json:"integrityLevel,omitempty"`
	TerminalSessionID uint32 `json:"terminalSessionId,omitempty"`
	FileVersion       string `json:"fileVersion,omitempty"`
	Description       string `json:"description,omitempty"`
	Product           string `json:"product,omitempty"`
	Company           string `json:"company,omitempty"`
	OriginalFileName  string `json:"originalFileName,omitempty"`

	Protocol            string `json:"protocol,omitempty"`
	Initiated           bool   `json:"initiated,omitempty"`
	SourceIP            string `json:"sourceIp,omitempty"`
	SourcePort          uint32 `json:"sourcePort,omitempty"`
	DestinationIP       string `json:"destinationIp,omitempty"`
	DestinationPort     uint32 `json:"destinationPort,omitempty"`
	SourceHostname      string `json:"sourceHostname,omitempty"`
	DestinationHostname string `json:"destinationHostname,omitempty"`
	SourceIsIPv6        bool   `json:"sourceIsIpv6,omitempty"`
	DestinationIsIPv6   bool   `json:"destinationIsIpv6,omitempty"`
}

func (t *HTTPTransport) Send(ctx context.Context, item evidence.Evidence) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}

	requestBody, err := t.buildRequest(item)
	if err != nil {
		return err
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal evidence request: %w", err)
	}

	url := t.backendURL + "/api/evidence"
	fmt.Printf("DEBUG evidence POST URL: %s\n", url)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create evidence request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := t.client.Do(request)
if err != nil {
    fmt.Printf("DEBUG HTTP CLIENT ERROR: %v\n", err)
    return fmt.Errorf("send evidence request: %w", err)
}
	defer response.Body.Close()

if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
    responseBody, _ := io.ReadAll(response.Body)

    fmt.Printf(
        "DEBUG evidence response: HTTP %d: %s\n",
        response.StatusCode,
        string(responseBody),
    )

    return fmt.Errorf("evidence ingestion returned HTTP %d", response.StatusCode)
}

	return nil
}

func (t *HTTPTransport) buildRequest(
	item evidence.Evidence,
) (ingestEvidenceRequest, error) {
	if strings.TrimSpace(item.Event.ID) == "" {
		return ingestEvidenceRequest{}, fmt.Errorf(
			"evidence event ID is required",
		)
	}

	if item.Event.Timestamp.IsZero() {
		return ingestEvidenceRequest{}, fmt.Errorf(
			"evidence timestamp is required",
		)
	}

	request := ingestEvidenceRequest{
		AgentID:      t.agentID,
		Hostname:     item.Event.Host,
		OSName:       t.osName,
		OSVersion:    t.osVersion,
		AgentVersion: t.agentVersion,

		Type:          mapEventType(item.Event.Type),
		Source:        item.Source,
		SourceEventID: item.SourceEventID,
		ObservedAt:    item.Event.Timestamp.UTC().Format(time.RFC3339Nano),
		CollectedAt:   item.CollectedAt.UTC().Format(time.RFC3339Nano),

		Payload: item,
	}

	switch item.Event.Type {
	case event.TypeProcess:
		if item.Event.Process == nil {
			return ingestEvidenceRequest{}, fmt.Errorf(
				"process evidence details are missing",
			)
		}

		process := item.Event.Process

		request.ProcessGuid = process.ProcessGuid
		request.ProcessID = process.ProcessID
		request.ParentProcessID = process.ParentProcessID
		request.ParentProcessGuid = process.ParentProcessGuid
		request.ImagePath = process.ImagePath
		request.CommandLine = process.CommandLine
		request.CurrentDirectory = process.CurrentDirectory
		request.User = process.User
		request.IntegrityLevel = process.IntegrityLevel
		request.TerminalSessionID = process.TerminalSessionID
		request.FileVersion = process.FileVersion
		request.Description = process.Description
		request.Product = process.Product
		request.Company = process.Company
		request.OriginalFileName = process.OriginalFileName

	case event.TypeNetwork:
		if item.Event.Network == nil {
			return ingestEvidenceRequest{}, fmt.Errorf(
				"network evidence details are missing",
			)
		}

		network := item.Event.Network

		request.ProcessGuid = network.ProcessGuid
		request.ProcessID = network.ProcessID
		request.Protocol = network.Protocol
		request.Initiated = network.Initiated
		request.SourceIP = network.SourceIP
		request.SourcePort = network.SourcePort
		request.DestinationIP = network.DestinationIP
		request.DestinationPort = network.DestinationPort
		request.SourceHostname = network.SourceHostname
		request.DestinationHostname = network.DestinationHostname
		request.SourceIsIPv6 = network.SourceIsIPv6
		request.DestinationIsIPv6 = network.DestinationIsIPv6

	case event.TypeDNS,
		event.TypeFile,
		event.TypePowerShell:
		// These event types are supported by the backend contract,
		// but their endpoint-specific normalized fields are not yet
		// represented in the agent Event model.

	default:
		return ingestEvidenceRequest{}, fmt.Errorf(
			"unsupported evidence event type: %s",
			item.Event.Type,
		)
	}

	return request, nil
}

func mapEventType(eventType event.Type) string {
	switch eventType {
	case event.TypeProcess:
		return "PROCESS"
	case event.TypeNetwork:
		return "NETWORK"
	case event.TypeDNS:
		return "DNS"
	case event.TypeFile:
		return "FILE"
	case event.TypePowerShell:
		return "POWERSHELL"
	default:
		return strings.ToUpper(string(eventType))
	}
}
