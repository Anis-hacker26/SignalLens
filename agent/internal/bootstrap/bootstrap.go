package bootstrap

import (
	"log"

	"github.com/Anis-hacker26/SignalLens/agent/internal/app"
	"github.com/Anis-hacker26/SignalLens/agent/internal/config"
	"github.com/Anis-hacker26/SignalLens/agent/internal/logging"
)

func NewApplication(cfg config.Config) *app.App {
	logger := logging.NewStandardLogger(log.Default())

	return app.New(cfg, logger)
}
