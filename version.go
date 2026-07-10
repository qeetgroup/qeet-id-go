package qeetid

import "github.com/qeetgroup/qeet-id-go/internal/version"

// Version is the released version of this SDK. It is sent on every request
// in the User-Agent header so API-side observability can attribute traffic
// and spot outdated clients.
const Version = version.SDKVersion
