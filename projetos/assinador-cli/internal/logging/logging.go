// Package logging provides structured (JSON) or human-readable logging
// for the CLIs, configurable via --log-format and --log-level.
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gabrielmmendes/runner/internal/version"
)

// Format selects the output encoding.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

var logger *slog.Logger

// ParseLevel maps a textual level to slog.Level. Unknown values default to info.
func ParseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("log-level inválido: %q (use debug|info|warn|error)", s)
	}
}

// ParseFormat validates the format string.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "text":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return FormatText, fmt.Errorf("log-format inválido: %q (use text|json)", s)
	}
}

// Init builds the global logger. Logs go to stderr so stdout stays reserved
// for command output (e.g. the JSON signing response).
func Init(format Format, level slog.Level) {
	var h slog.Handler
	w := io.Writer(os.Stderr)
	opts := &slog.HandlerOptions{Level: level}
	if format == FormatJSON {
		h = slog.NewJSONHandler(w, opts)
	} else {
		h = newTextHandler(w, opts)
	}
	logger = slog.New(h).With(slog.String("version", version.Version))
}

// L returns the configured logger, initializing a default text logger on first
// use if Init was never called.
func L() *slog.Logger {
	if logger == nil {
		Init(FormatText, slog.LevelInfo)
	}
	return logger
}

// WithCommand returns a logger tagged with the running command name.
func WithCommand(name string) *slog.Logger {
	return L().With(slog.String("command", name))
}
