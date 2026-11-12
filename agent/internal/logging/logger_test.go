package logging

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestStandardLoggerInfo(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)

	standardLogger := NewStandardLogger(logger)
	standardLogger.Info("agent started: %s", "test-agent")

	output := buffer.String()

	if !strings.Contains(output, "INFO: agent started: test-agent") {
		t.Fatalf("expected info message, got: %q", output)
	}
}

func TestStandardLoggerError(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)

	standardLogger := NewStandardLogger(logger)
	standardLogger.Error("component failed: %s", "collector")

	output := buffer.String()

	if !strings.Contains(output, "ERROR: component failed: collector") {
		t.Fatalf("expected error message, got: %q", output)
	}
}

func TestStandardLoggerDebug(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)

	standardLogger := NewStandardLogger(logger)
	standardLogger.Debug("processing event: %d", 42)

	output := buffer.String()

	if !strings.Contains(output, "DEBUG: processing event: 42") {
		t.Fatalf("expected debug message, got: %q", output)
	}
}
