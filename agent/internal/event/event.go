package event

import "time"

type Type string

const (
	TypeProcess    Type = "process"
	TypeNetwork    Type = "network"
	TypeDNS        Type = "dns"
	TypeFile       Type = "file"
	TypePowerShell Type = "powershell"
)

type ProcessDetails struct {
	ProcessGuid       string
	ProcessID         uint32
	ParentProcessID   uint32
	ParentProcessGuid string
	ImagePath         string
	CommandLine       string
	CurrentDirectory  string
	User              string
	IntegrityLevel    string
	TerminalSessionID uint32
	FileVersion       string
	Description       string
	Product           string
	Company           string
	OriginalFileName  string
}

type NetworkDetails struct {
	ProcessGuid         string
	ProcessID           uint32
	Protocol            string
	Initiated           bool
	SourceIP            string
	SourcePort          uint32
	DestinationIP       string
	DestinationPort     uint32
	SourceHostname      string
	DestinationHostname string
	SourceIsIPv6        bool
	DestinationIsIPv6   bool
}

type Event struct {
	ID            string
	Type          Type
	Timestamp     time.Time
	Host          string
	Source        string
	SourceEventID string

	Process *ProcessDetails
	Network *NetworkDetails
}
