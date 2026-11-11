package app

import (
	"context"

	"github.com/Anis-hacker26/SignalLens/agent/internal/component"
	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

type App struct {
	config     config.Config
	components []component.Component
}

func New(cfg config.Config, components ...component.Component) *App {
	return &App{
		config:     cfg,
		components: components,
	}
}

func (a *App) Run(ctx context.Context) error {
	for _, component := range a.components {
		if err := component.Start(ctx); err != nil {
			return err
		}
	}

	<-ctx.Done()

	for i := len(a.components) - 1; i >= 0; i-- {
		if err := a.components[i].Stop(ctx); err != nil {
			return err
		}
	}

	return nil
}
