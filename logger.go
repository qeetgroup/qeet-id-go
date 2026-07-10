package qeetid

import "github.com/qeetgroup/qeet-id-go/internal/transport"

// Logger is an optional per-request observability hook — invoked once per
// request after the retry loop settles, with the method, path, final status,
// duration, and request ID. Implement it for structured logging or metrics;
// it never affects control flow. A nil Logger (the default) is a no-op —
// there is no OpenTelemetry dependency baked into the core SDK.
type Logger = transport.Logger
