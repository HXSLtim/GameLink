package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Init sets the default slog logger with JSON handler and given level.
func Init(level string) *slog.Logger {
	lvl := parseLevel(level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: false})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(v string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
