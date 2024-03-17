package bot

import (
	"context"

	"cdr.dev/slog"
	"github.com/ppaanngggg/option-bot/pkg/util"
)

var Scheduler = &scheduler{
	logger: util.DefaultLogger.With(slog.F("bot", "scheduler")),
}

type scheduler struct {
	logger slog.Logger
}

func (s *scheduler) Run() {
	s.logger.Info(context.Background(), "scheduler started")
}
