package qeetid

import (
	"net/http"
	"time"
)

// Config configures the client. APIKey is required (a server-side `qk_…` key).
type Config struct {
	APIKey string

	// BaseURL overrides the default https://api.id.qeet.in — set it to point
	// at a self-hosted instance.
	BaseURL string

	// HTTPClient overrides the transport entirely (custom proxying, TLS
	// config, etc.). If set, Timeout below is ignored — configure the
	// timeout on the client you provide.
	HTTPClient *http.Client

	// Timeout sets the per-request timeout when HTTPClient is nil (default
	// 10s). Ignored if HTTPClient is provided.
	Timeout time.Duration

	// MaxRetries is the retry budget for 429/5xx on idempotent requests
	// (default 2).
	MaxRetries int

	// Headers are sent on every management-API request — an escape hatch for
	// tracing headers or forward-compatible options. They cannot override the
	// Authorization, Accept, Content-Type, or User-Agent headers.
	Headers http.Header

	// UserAgent, if set, is prepended to the SDK's own User-Agent string
	// (e.g. "myapp/1.2.0") rather than replacing it — API-side observability
	// can still attribute traffic to this SDK version.
	UserAgent string

	// Logger is an optional per-request observability hook. Nil is a no-op.
	Logger Logger
}
