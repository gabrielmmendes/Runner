package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

// textHandler emits concise human-readable lines: "LEVEL  message  key=value ...".
// The version/command attrs added via With are kept but rendered compactly so
// terminal output stays readable while still carrying structured context.
type textHandler struct {
	w     io.Writer
	mu    *sync.Mutex
	opts  *slog.HandlerOptions
	attrs []slog.Attr
}

func newTextHandler(w io.Writer, opts *slog.HandlerOptions) *textHandler {
	return &textHandler{w: w, mu: &sync.Mutex{}, opts: opts}
}

func (h *textHandler) Enabled(_ context.Context, l slog.Level) bool {
	min := slog.LevelInfo
	if h.opts != nil && h.opts.Level != nil {
		min = h.opts.Level.Level()
	}
	return l >= min
}

func (h *textHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	b.WriteString(levelLabel(r.Level))
	b.WriteString("  ")
	b.WriteString(r.Message)
	for _, a := range h.attrs {
		writeAttr(&b, a)
	}
	r.Attrs(func(a slog.Attr) bool {
		writeAttr(&b, a)
		return true
	})
	b.WriteByte('\n')
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.w, b.String())
	return err
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	n := *h
	n.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &n
}

func (h *textHandler) WithGroup(_ string) slog.Handler { return h }

func writeAttr(b *strings.Builder, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}
	fmt.Fprintf(b, "  %s=%v", a.Key, a.Value.Any())
}

func levelLabel(l slog.Level) string {
	switch {
	case l < slog.LevelInfo:
		return "DEBUG"
	case l < slog.LevelWarn:
		return "INFO "
	case l < slog.LevelError:
		return "WARN "
	default:
		return "ERROR"
	}
}
