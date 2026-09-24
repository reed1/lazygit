package utils

import "os"

// LAZYGIT_UPSTREAM_BEHAVIOR=1 switches off this fork's changes, which is how
// upstream's tests run unmodified.
func UpstreamBehavior() bool {
	return os.Getenv("LAZYGIT_UPSTREAM_BEHAVIOR") == "1"
}
