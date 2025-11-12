package logging

import (
	"log/slog"
	"testing"
)

func TestInitReturnsLogger(t *testing.T) {
    t.Run("returns non-nil logger", func(t *testing.T) {
        logger := Init("debug")
        if logger == nil {
            t.Fatal("expected logger instance")
        }
    })
}

func TestInitWithVariousLevels(t *testing.T) {
    levels := []string{"info", "warn", "warning", "error", "DEBUG", "unknown"}
    for _, lvl := range levels {
        if Init(lvl) == nil {
            t.Fatalf("Init(%q) returned nil", lvl)
        }
    }
}

func TestParseLevel(t *testing.T) {
	tcases := []struct {
		name  string
		input string
		want  slog.Level
	}{
		{"debug", "DEBUG", slog.LevelDebug},
		{"warn", "warn", slog.LevelWarn},
		{"warning alias", " Warning ", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"default info", "unknown", slog.LevelInfo},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseLevel(tc.input)
			if got != tc.want {
				t.Fatalf("parseLevel(%q)=%v want %v", tc.input, got, tc.want)
			}
		})
	}
}
