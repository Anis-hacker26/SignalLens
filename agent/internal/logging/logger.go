package logging

import (
	"log"
)

type StandardLogger struct {
	logger *log.Logger
}

func NewStandardLogger(logger *log.Logger) *StandardLogger {
	return &StandardLogger{
		logger: logger,
	}
}

func (l *StandardLogger) Info(message string, args ...any) {
	l.logger.Printf("INFO: "+message, args...)
}

func (l *StandardLogger) Error(message string, args ...any) {
	l.logger.Printf("ERROR: "+message, args...)
}

func (l *StandardLogger) Debug(message string, args ...any) {
	l.logger.Printf("DEBUG: "+message, args...)
}
