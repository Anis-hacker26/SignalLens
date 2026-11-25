//go:build windows

package collector

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
	"golang.org/x/sys/windows"
)

const (
	sysmonChannel  = "Microsoft-Windows-Sysmon/Operational"
	sysmonProvider = "Microsoft-Windows-Sysmon"

	evtQueryChannelPath = 0x1
	evtQueryReverse     = 0x200

	evtRenderRenderXML = 1
)

var (
	wevtapi = windows.NewLazySystemDLL("wevtapi.dll")

	procEvtQuery  = wevtapi.NewProc("EvtQuery")
	procEvtNext   = wevtapi.NewProc("EvtNext")
	procEvtRender = wevtapi.NewProc("EvtRender")
	procEvtClose  = wevtapi.NewProc("EvtClose")
)

type SysmonCollector struct {
	maxEvents    uint32
	lastRecordID uint64
}

func NewSysmonCollector() *SysmonCollector {
	return &SysmonCollector{
		maxEvents: 20,
	}
}

func (c *SysmonCollector) SetCursor(recordID uint64) {
	c.lastRecordID = recordID
}

func (c *SysmonCollector) Commit(recordID uint64) {
	if recordID > c.lastRecordID {
		c.lastRecordID = recordID
	}
}

func filterNewEvents(
	events []telemetry.RawEvent,
	lastRecordID uint64,
) []telemetry.RawEvent {
	results := make([]telemetry.RawEvent, 0, len(events))

	for _, event := range events {
		if event.RecordID <= lastRecordID {
			continue
		}

		results = append(results, event)
	}

	return results
}

func (c *SysmonCollector) Collect(
	ctx context.Context,
) ([]telemetry.RawEvent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

query := `*[System[(EventID=1 or EventID=3)]]`

	channelPtr, err := windows.UTF16PtrFromString(sysmonChannel)
	if err != nil {
		return nil, fmt.Errorf("convert channel name: %w", err)
	}

	queryPtr, err := windows.UTF16PtrFromString(query)
	if err != nil {
		return nil, fmt.Errorf("convert query: %w", err)
	}

	handle, err := evtQuery(
		channelPtr,
		queryPtr,
		evtQueryChannelPath,
	)
	if err != nil {
		return nil, fmt.Errorf("query Sysmon event log: %w", err)
	}
	defer evtClose(handle)

	events := make([]telemetry.RawEvent, 0, c.maxEvents)

	for uint32(len(events)) < c.maxEvents {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		eventHandle, ok, err := evtNext(handle)
		if err != nil {
			return nil, fmt.Errorf("read Sysmon event: %w", err)
		}

if !ok {
    fmt.Printf("DEBUG EvtNext returned no more items\n")
    break
}

		xmlData, err := evtRenderXML(eventHandle)
		evtClose(eventHandle)

		if err != nil {
			return nil, fmt.Errorf("render Sysmon event: %w", err)
		}

		rawEvent, err := parseSysmonEvent(xmlData)
		if err != nil {
			return nil, fmt.Errorf("parse Sysmon event: %w", err)
		}

		if rawEvent.RecordID <= c.lastRecordID {
			continue
		}

		events = append(events, rawEvent)
	}

	return events, nil
}

func evtQuery(
	path *uint16,
	query *uint16,
	flags uint32,
) (windows.Handle, error) {
	r1, _, callErr := procEvtQuery.Call(
		0,
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(query)),
		uintptr(flags),
	)

	handle := windows.Handle(r1)

	if handle == 0 {
		return 0, callErr
	}

	return handle, nil
}

func evtNext(
	query windows.Handle,
) (windows.Handle, bool, error) {
	var eventHandle windows.Handle
	var returned uint32

	r1, _, callErr := procEvtNext.Call(
		uintptr(query),
		1,
		uintptr(unsafe.Pointer(&eventHandle)),
		0,
		0,
		uintptr(unsafe.Pointer(&returned)),
	)

if r1 == 0 {
    if callErr == windows.ERROR_NO_MORE_ITEMS {
        return 0, false, nil
    }

    return 0, false, callErr
}

	if returned != 1 || eventHandle == 0 {
		return 0, false, fmt.Errorf(
			"EvtNext returned an invalid event handle",
		)
	}

	return eventHandle, true, nil
}

func evtRenderXML(
	eventHandle windows.Handle,
) (string, error) {
	var bufferUsed uint32
	var propertyCount uint32

	// First call determines the required buffer size.
	r1, _, err := procEvtRender.Call(
		0,
		uintptr(eventHandle),
		evtRenderRenderXML,
		0,
		0,
		uintptr(unsafe.Pointer(&bufferUsed)),
		uintptr(unsafe.Pointer(&propertyCount)),
	)

	if r1 == 0 {
		if err != nil && err != syscall.Errno(0) {
			if err != windows.ERROR_INSUFFICIENT_BUFFER {
				return "",
					fmt.Errorf(
						"EvtRender size query failed: %w",
						err,
					)
			}
		}
	}

	if bufferUsed == 0 {
		return "",
			fmt.Errorf(
				"EvtRender reported zero buffer size",
			)
	}

	buffer := make([]byte, bufferUsed)

	// Second call renders the actual XML.
	r1, _, err = procEvtRender.Call(
		0,
		uintptr(eventHandle),
		evtRenderRenderXML,
		uintptr(bufferUsed),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bufferUsed)),
		uintptr(unsafe.Pointer(&propertyCount)),
	)

	if r1 == 0 {
		if err != nil && err != syscall.Errno(0) {
			return "",
				fmt.Errorf(
					"EvtRender failed: %w",
					err,
				)
		}

		return "",
			fmt.Errorf(
				"EvtRender failed",
			)
	}

	// EvtRender returns a null-terminated UTF-16 XML string.
	utf16Buffer := unsafe.Slice(
		(*uint16)(unsafe.Pointer(&buffer[0])),
		len(buffer)/2,
	)

	return windows.UTF16ToString(utf16Buffer), nil
}

func evtClose(handle windows.Handle) {
	if handle == 0 {
		return
	}

	procEvtClose.Call(uintptr(handle))
}

type sysmonEventXML struct {
	System struct {
		Provider struct {
			Name string `xml:"Name,attr"`
		} `xml:"Provider"`

		EventID       uint32 `xml:"EventID"`
		EventRecordID uint64 `xml:"EventRecordID"`
		TimeCreated   struct {
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`

		Computer string `xml:"Computer"`
	} `xml:"System"`
}

func parseSysmonEvent(xmlData string) (telemetry.RawEvent, error) {
	var parsed sysmonEventXML

	if err := xml.Unmarshal(
		[]byte(xmlData),
		&parsed,
	); err != nil {
		return telemetry.RawEvent{}, fmt.Errorf(
			"decode event XML: %w",
			err,
		)
	}

	if parsed.System.Provider.Name == "" {
		return telemetry.RawEvent{}, fmt.Errorf(
			"event provider is missing",
		)
	}

	if parsed.System.EventID == 0 {
		return telemetry.RawEvent{}, fmt.Errorf(
			"event ID is missing",
		)
	}

	if parsed.System.TimeCreated.SystemTime == "" {
		return telemetry.RawEvent{}, fmt.Errorf(
			"event timestamp is missing",
		)
	}

	timestamp, err := parseSysmonTimestamp(
		parsed.System.TimeCreated.SystemTime,
	)
	if err != nil {
		return telemetry.RawEvent{}, err
	}

	return telemetry.RawEvent{
		RecordID:  parsed.System.EventRecordID,
		EventID:   parsed.System.EventID,
		Timestamp: timestamp,
		Provider:  parsed.System.Provider.Name,
		Computer:  parsed.System.Computer,
		Message:   strings.TrimSpace(xmlData),
	}, nil
}

func parseSysmonTimestamp(value string) (
	timestamp time.Time,
	err error,
) {
	return time.Parse(
		"2006-01-02T15:04:05.9999999Z",
		value,
	)
}