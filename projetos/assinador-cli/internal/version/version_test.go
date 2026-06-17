package version

import (
	"strings"
	"testing"
)

func TestStringIncludesAllIdentifiers(t *testing.T) {
	Version, Commit, Date = "v1.2.3", "abc1234", "2026-06-16T00:00:00Z"
	got := String()
	for _, want := range []string{"v1.2.3", "abc1234", "2026-06-16T00:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Errorf("String()=%q missing %q", got, want)
		}
	}
}

func TestStringDefaults(t *testing.T) {
	Version, Commit, Date = "dev", "none", "unknown"
	if got := String(); got != "dev (none, built unknown)" {
		t.Errorf("unexpected default String()=%q", got)
	}
}
