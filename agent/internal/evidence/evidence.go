package evidence

import (
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

type Evidence struct {
	ID            string
	Event         event.Event
	CollectedAt   time.Time
	Source        string
	SourceEventID string
}
