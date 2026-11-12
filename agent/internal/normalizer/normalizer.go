package normalizer

import (
	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type Normalizer interface {
	Normalize(raw telemetry.RawEvent) (event.Event, error)
}
