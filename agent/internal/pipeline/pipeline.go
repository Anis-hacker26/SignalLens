package pipeline

import (
	"fmt"

	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
	"github.com/Anis-hacker26/SignalLens/agent/internal/normalizer"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type Pipeline struct {
	normalizer normalizer.Normalizer
	builder    *evidence.Builder
}

func New(
	eventNormalizer normalizer.Normalizer,
	builder *evidence.Builder,
) (*Pipeline, error) {
	if eventNormalizer == nil {
		return nil, fmt.Errorf("normalizer is required")
	}

	if builder == nil {
		return nil, fmt.Errorf("evidence builder is required")
	}

	return &Pipeline{
		normalizer: eventNormalizer,
		builder:    builder,
	}, nil
}

func (p *Pipeline) Process(
	raw telemetry.RawEvent,
) (evidence.Evidence, error) {
	normalized, err := p.normalizer.Normalize(raw)
	if err != nil {
		return evidence.Evidence{}, fmt.Errorf(
			"normalize raw event: %w",
			err,
		)
	}

	result, err := p.builder.Build(normalized)
	if err != nil {
		return evidence.Evidence{}, fmt.Errorf(
			"build evidence: %w",
			err,
		)
	}

	return result, nil
}

func (p *Pipeline) ProcessBatch(
	rawEvents []telemetry.RawEvent,
) ([]evidence.Evidence, error) {
	results := make([]evidence.Evidence, 0, len(rawEvents))

	for _, raw := range rawEvents {
		result, err := p.Process(raw)
		if err != nil {
			return nil, fmt.Errorf(
				"process raw event: %w",
				err,
			)
		}

		results = append(results, result)
	}

	return results, nil
}

func NewDefault() (*Pipeline, error) {
	return New(
		normalizer.NewSysmonNormalizer(),
		evidence.NewBuilder(),
	)
}
