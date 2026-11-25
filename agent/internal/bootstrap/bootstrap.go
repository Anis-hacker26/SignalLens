package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/app"
	"github.com/Anis-hacker26/SignalLens/agent/internal/collector"
	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
	"github.com/Anis-hacker26/SignalLens/agent/internal/logging"
	"github.com/Anis-hacker26/SignalLens/agent/internal/pipeline"
	"github.com/Anis-hacker26/SignalLens/agent/internal/platform"
	"github.com/Anis-hacker26/SignalLens/agent/internal/runtime"
	"github.com/Anis-hacker26/SignalLens/agent/internal/state"
	"github.com/Anis-hacker26/SignalLens/agent/internal/transport"
)

const runtimeInterval = 10 * time.Second

func BuildRuntime(cfg config.Config) (*runtime.Runtime, error) {
	sysmonCollector := collector.NewSysmonCollector()

	cursorStore, err := state.NewFileCursorStore(
		cfg.StatePath,
	)
	if err != nil {
		return nil, err
	}

	defaultPipeline, err := pipeline.NewDefault()
	if err != nil {
		return nil, err
	}

	osInfo, err := platform.GetOSInfo()
	if err != nil {
		return nil, err
	}

	httpTransport, err := transport.NewHTTPTransport(
		transport.HTTPTransportConfig{
			BackendURL:   cfg.BackendURL,
			AgentID:      cfg.AgentID,
			AgentVersion: cfg.AgentVersion,
			OSName:       osInfo.Name,
			OSVersion:    osInfo.Version,
		},
	)
	if err != nil {
		return nil, err
	}

	return runtime.New(
		sysmonCollector,
		defaultPipeline,
		cursorStore,
		runtimeInterval,
		func(ctx context.Context, evidenceBatch []evidence.Evidence) error {
			for _, item := range evidenceBatch {
				if err := httpTransport.Send(ctx, item); err != nil {
					return fmt.Errorf(
						"send evidence %s: %w",
						item.ID,
						err,
					)
				}
			}

			return nil
		},
	)
}

func NewApplication(cfg config.Config) *app.App {
	stdLogger := log.New(
		os.Stdout,
		"signallens-agent ",
		log.LstdFlags,
	)

	logger := logging.NewStandardLogger(stdLogger)

	runtimeComponent, err := BuildRuntime(cfg)
	if err != nil {
		logger.Error("failed to build runtime: %v", err)
		return nil
	}

	return app.New(cfg, logger, runtimeComponent)
}
