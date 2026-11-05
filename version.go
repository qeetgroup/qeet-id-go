package qeetid

import "github.com/qeetgroup/qeet-id-go/internal/version"

// Version is the release version of this module. The client sends it in the
// User-Agent header so that Qeet ID can attribute traffic to SDK versions.
const Version = version.SDK
