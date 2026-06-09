package logging

import (
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"":      slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
		"WARN":  slog.LevelWarn,
	}
	for in, want := range cases {
		got, err := ParseLevel(in)
		if err != nil {
			t.Errorf("ParseLevel(%q) erro: %v", in, err)
		}
		if got != want {
			t.Errorf("ParseLevel(%q) = %v, want %v", in, got, want)
		}
	}
	if _, err := ParseLevel("loud"); err == nil {
		t.Error("nível inválido deveria falhar")
	}
}

func TestParseFormat(t *testing.T) {
	for _, in := range []string{"", "text", "json", "JSON"} {
		if _, err := ParseFormat(in); err != nil {
			t.Errorf("ParseFormat(%q) erro: %v", in, err)
		}
	}
	if _, err := ParseFormat("yaml"); err == nil {
		t.Error("formato inválido deveria falhar")
	}
}
