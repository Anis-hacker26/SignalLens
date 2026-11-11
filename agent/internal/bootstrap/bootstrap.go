package bootstrap

import (
	"github.com/Anis-hacker26/SignalLens/agent/internal/app"
	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
)

func NewApplication(cfg config.Config) *app.App {
	return app.New(cfg)
}
