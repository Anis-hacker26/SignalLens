package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

type Builder struct {
	clock func() time.Time
}

func NewBuilder() *Builder {
	return &Builder{
		clock: time.Now,
	}
}

func (b *Builder) Build(
	normalized event.Event,
) (Evidence, error) {
	if normalized.ID == "" {
		return Evidence{}, fmt.Errorf("event ID is required")
	}

	if normalized.Type == "" {
		return Evidence{}, fmt.Errorf("event type is required")
	}

	if normalized.Timestamp.IsZero() {
		return Evidence{}, fmt.Errorf("event timestamp is required")
	}

	if normalized.Host == "" {
		return Evidence{}, fmt.Errorf("event host is required")
	}

	if normalized.Source == "" {
		return Evidence{}, fmt.Errorf("event source is required")
	}

	if normalized.SourceEventID == "" {
		return Evidence{}, fmt.Errorf("source event ID is required")
	}

	collectedAt := b.clock()

	return Evidence{
		ID:            buildEvidenceID(normalized),
		Event:         normalized,
		CollectedAt:   collectedAt,
		Source:        normalized.Source,
		SourceEventID: normalized.SourceEventID,
	}, nil
}

func buildEvidenceID(e event.Event) string {
	input := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s",
		e.ID,
		e.Type,
		e.Timestamp.UTC().Format(time.RFC3339Nano),
		e.Host,
		e.Source,
		e.SourceEventID,
	)

	hash := sha256.Sum256([]byte(input))

	return hex.EncodeToString(hash[:])
}
