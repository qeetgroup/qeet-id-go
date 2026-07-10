package transport

import (
	"fmt"
	"runtime"

	"github.com/qeetgroup/qeet-id-go/internal/version"
)

// UserAgent identifies the SDK, its version, and the Go runtime — e.g.
// "qeet-id-go/0.1.0 go1.23.0 (darwin/arm64)". Sent on every request so
// API-side observability can attribute traffic and spot outdated clients.
func UserAgent() string {
	return fmt.Sprintf("qeet-id-go/%s %s (%s/%s)", version.SDKVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
