package collector

import (
	"context"

	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type Collector interface {
	Collect(ctx context.Context) ([]telemetry.RawEvent, error)
}
