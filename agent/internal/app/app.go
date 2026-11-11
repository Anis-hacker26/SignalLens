package app

import (
	"context"

	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

type App struct {
	config config.Config
}

func New(cfg config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	<-ctx.Done()

	return nil
}
