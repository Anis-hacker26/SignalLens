package transport

import (
	"context"

	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
)

type Transport interface {
	Send(ctx context.Context, evidence evidence.Evidence) error
}
