package runtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/collector"
	"github.com/Anis-hacker26/SignalLens/agent/internal/component"
	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
	"github.com/Anis-hacker26/SignalLens/agent/internal/pipeline"
	"github.com/Anis-hacker26/SignalLens/agent/internal/state"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type Runtime struct {
	collector collector.Collector
	pipeline  *pipeline.Pipeline
	interval  time.Duration
	state     state.CursorStore

	onEvidence func(context.Context, []evidence.Evidence) error

	cancel context.CancelFunc
	done   chan struct{}
}

var _ component.Component = (*Runtime)(nil)

func New(
	eventCollector collector.Collector,
	eventPipeline *pipeline.Pipeline,
	cursorStore state.CursorStore,
	interval time.Duration,
	onEvidence func(context.Context, []evidence.Evidence) error,
) (*Runtime, error) {
	if eventCollector == nil {
		return nil, fmt.Errorf("collector is required")
	}

	if eventPipeline == nil {
		return nil, fmt.Errorf("pipeline is required")
	}

	if cursorStore == nil {
		return nil, fmt.Errorf("cursor store is required")
	}

	if interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than zero")
	}

	return &Runtime{
		collector:  eventCollector,
		pipeline:   eventPipeline,
		interval:   interval,
		state:      cursorStore,
		onEvidence: onEvidence,
		done:       make(chan struct{}),
	}, nil
}

func (r *Runtime) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is required")
	}

	if r.cancel != nil {
		return errors.New("runtime is already running")
	}

	recordID, err := r.state.Load(ctx)
	if err != nil {
		return fmt.Errorf("load cursor state: %w", err)
	}

	if initializer, ok := r.collector.(collector.CursorInitializer); ok {
		initializer.SetCursor(recordID)
	}

	runCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel

	go func() {
		defer close(r.done)

		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()

		r.collectOnce(runCtx)

		for {
			select {
			case <-runCtx.Done():
				return

			case <-ticker.C:
				r.collectOnce(runCtx)
			}
		}
	}()

	return nil
}

func (r *Runtime) Stop(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is required")
	}

	if r.cancel == nil {
		return nil
	}

	r.cancel()

	select {
	case <-r.done:
		r.cancel = nil
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Runtime) Run(ctx context.Context) ([]evidence.Evidence, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rawEvents, err := r.collector.Collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("collect events: %w", err)
	}

	if len(rawEvents) == 0 {
		results := []evidence.Evidence{}

		if r.onEvidence != nil {
			if err := r.onEvidence(ctx, results); err != nil {
				return nil, fmt.Errorf("deliver evidence: %w", err)
			}
		}

		return results, nil
	}

	results, err := r.pipeline.ProcessBatch(rawEvents)
	if err != nil {
		return nil, fmt.Errorf("process events: %w", err)
	}

	if results == nil {
		results = []evidence.Evidence{}
	}

	if r.onEvidence != nil {
		if err := r.onEvidence(ctx, results); err != nil {
			return nil, fmt.Errorf("deliver evidence: %w", err)
		}
	}

	latestRecordID := highestRecordID(rawEvents)

	if latestRecordID > 0 {
		if err := r.state.Save(ctx, latestRecordID); err != nil {
			return nil, fmt.Errorf(
				"save cursor state: %w",
				err,
			)
		}

		if committer, ok := r.collector.(collector.CursorCommitter); ok {
			committer.Commit(latestRecordID)
		}
	}

	return results, nil
}

func (r *Runtime) collectOnce(ctx context.Context) {
	if _, err := r.Run(ctx); err != nil {
		fmt.Printf("runtime collection cycle failed: %v\n", err)
	}
}

func highestRecordID(events []telemetry.RawEvent) uint64 {
	var highest uint64

	for _, event := range events {
		if event.RecordID > highest {
			highest = event.RecordID
		}
	}

	return highest
}
