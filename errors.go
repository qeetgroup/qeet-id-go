package qeetid

import "github.com/qeetgroup/qeet-id-go/internal/transport"

// Error is returned by every failed API call. Inspect Status or use the Is*
// helpers; errors.As(err, &qeetid.Error{}) unwraps it. This is a type alias
// for internal/transport.Error, not a new type — every resource's call into
// the transport already returns this concrete value, so no wrapping is
// needed at the boundary.
type Error = transport.Error
