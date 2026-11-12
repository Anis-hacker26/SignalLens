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

type Event struct {
	ID            string
	Type          Type
	Timestamp     time.Time
	Host          string
	Source        string
	SourceEventID string
}
