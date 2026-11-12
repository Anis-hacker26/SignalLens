package normalizer

import (
	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

type Normalizer interface {
	Normalize(raw any) (event.Event, error)
}
