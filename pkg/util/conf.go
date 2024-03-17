package util

import (
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogjson"
	"github.com/caarlos0/env/v10"
)

var Conf = &struct {
	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
		Port int    `env:"SERVER_PORT" envDefault:"8000"`
	}
	Log struct {
		Level  string `env:"LOG_LEVEL" envDefault:"info"`
		Format string `env:"LOG_FORMAT" envDefault:"human"`
	}
}{}

var DefaultLogger slog.Logger

func init() {
	if err := env.Parse(Conf); err != nil {
		panic(err)
	}

	{
		// set logger
		var level slog.Level
		switch Conf.Log.Level {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			panic("invalid log level")
		}
		var sink slog.Sink
		switch Conf.Log.Format {
		case "json":
			sink = slogjson.Sink(os.Stdout)
		case "human":
			sink = sloghuman.Sink(os.Stdout)
		default:
			panic("invalid log format")
		}

		DefaultLogger = slog.Make(sink).Leveled(level)
	}

}
