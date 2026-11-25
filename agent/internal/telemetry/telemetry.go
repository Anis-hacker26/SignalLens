package telemetry

import "time"

type RawEvent struct {
	RecordID  uint64
	EventID   uint32
	Timestamp time.Time
	Provider  string
	Computer  string
	Message   string
}
