package app

import (
	"context"

	"github.com/Anis-hacker26/SignalLens/agent/internal/component"
	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
	"github.com/Anis-hacker26/SignalLens/agent/internal/logging"
)

type App struct {
	config     config.Config
	logger     logging.Logger
	components []component.Component
}

func New(
	cfg config.Config,
	logger logging.Logger,
	components ...component.Component,
) *App {
	return &App{
		config:     cfg,
		logger:     logger,
		components: components,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("application starting")

	for _, component := range a.components {
		if err := component.Start(ctx); err != nil {
			a.logger.Error("component start failed: %v", err)
			return err
		}
	}

	<-ctx.Done()

	a.logger.Info("application stopping")

	for i := len(a.components) - 1; i >= 0; i-- {
		if err := a.components[i].Stop(ctx); err != nil {
			a.logger.Error("component stop failed: %v", err)
			return err
		}
	}

	return nil
}
