package logger

import (
	"log/slog"
	"os"

	"github.com/rshelekhov/golib/logger/handler"
	"github.com/rshelekhov/golib/logger/handler/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var h slog.Handler

	switch env {
	case envLocal:
		h = slogpretty.NewPrettyHandler(os.Stdout, &slogpretty.Options{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	case envDev:
		h = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}))
	case envProd:
		h = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		}))
	}

	h = handler.NewHandlerMiddleware(h)
	log := slog.New(h)

	return log
}
