package utils

import (
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"os"
)

var DefaultLogger = slog.Make(sloghuman.Sink(os.Stdout))
