package version

import "fmt"

// These variables are meant to be overridden at build time using -ldflags.
// Defaults are useful for local development.
var (
	// Version is the semantic version of the binary, e.g. "v1.0.0".
	Version = "dev"

	// Commit is the git commit hash used for the build.
	Commit = "none"

	// Date is the build date in RFC3339 or similar format.
	Date = "unknown"
)

// Full returns a human-readable version string.
func Full() string {
	return fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
}
