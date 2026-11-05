// Package version holds the module's release version and the User-Agent built
// from it, so the root package and its sub-packages report the same value.
package version

import (
	"runtime"
)

// SDK is the release version of the module. Keep it equal to the tag being
// released; it is sent in the User-Agent header.
const SDK = "0.1.0"

// UserAgent identifies the SDK, its version, the Go runtime and the platform,
// for example "qeet-id-go/0.1.0 (go1.27.1; linux/amd64)".
func UserAgent() string {
	return "qeet-id-go/" + SDK + " (" + runtime.Version() + "; " + runtime.GOOS + "/" + runtime.GOARCH + ")"
}
