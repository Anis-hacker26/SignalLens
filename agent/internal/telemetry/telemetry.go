package telemetry

import "time"

type RawEvent struct {
	EventID   uint32
	Timestamp time.Time
	Provider  string
	Computer  string
	Message   string
}
