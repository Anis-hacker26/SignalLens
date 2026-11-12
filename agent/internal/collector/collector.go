package collector

import (
	"context"

	"github.com/Anis-hacker26/SignalLens/agent/internal/event"
)

type Collector interface {
	Collect(ctx context.Context) ([]event.Event, error)
}
