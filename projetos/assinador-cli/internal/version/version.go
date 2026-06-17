// Package version exposes build identifiers injected via -ldflags at build
// time. Defaults apply for local `go run`/`go build` without ldflags.
package version

import "fmt"

var (
	// Version is the semantic tag (e.g. v0.2.0) or "dev" for untagged builds.
	Version = "dev"
	// Commit is the short git SHA of the build (e.g. 2f88d67).
	Commit = "none"
	// Date is the build timestamp in RFC 3339 (UTC).
	Date = "unknown"
)

// String returns a traceable identifier combining tag, short SHA and date,
// e.g. "v0.2.0 (2f88d67, built 2026-06-16T12:00:00Z)".
func String() string {
	return fmt.Sprintf("%s (%s, built %s)", Version, Commit, Date)
}
