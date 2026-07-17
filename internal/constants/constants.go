// Package constants holds default configuration values and wire-protocol
// constants shared across internal/transport, internal/auth, and the root
// package. Kept dependency-free (stdlib only) so it can sit at the bottom of
// the internal import graph without risk of a cycle.
package constants

import "time"

const (
	// DefaultBaseURL is used when Config.BaseURL is empty.
	DefaultBaseURL = "https://api.id.qeet.in"
	// DefaultTimeout is the per-request timeout when no HTTPClient is supplied.
	DefaultTimeout = 10 * time.Second
	// DefaultMaxRetries is the retry budget for 429/5xx on idempotent requests.
	DefaultMaxRetries = 2
	// MaxResponseBytes caps how much of a response body is read into memory.
	MaxResponseBytes = 1 << 20 // 1 MiB
)

// HTTP header names used on every request/response.
const (
	HeaderAuthorization = "Authorization"
	HeaderAccept        = "Accept"
	HeaderContentType   = "Content-Type"
	HeaderUserAgent     = "User-Agent"
	HeaderRequestID     = "X-Request-Id"
	HeaderRetryAfter    = "Retry-After"
)

// Webhook delivery headers set by the Qeet ID dispatcher.
const (
	HeaderWebhookSignature = "X-Qeet-Signature"
	HeaderWebhookEvent     = "X-Qeet-Event"
)

// JWKSCacheTTL is how long a fetched JWKS document is trusted before a
// scheduled refresh. An unknown kid triggers an out-of-band refresh, but no
// more often than JWKSRefreshCooldown (see below).
const JWKSCacheTTL = 5 * time.Minute

// JWKSRefreshCooldown bounds how often an unknown/garbage kid may trigger an
// out-of-band JWKS refetch. Without it, tokens bearing random kids force one
// upstream JWKS fetch per request — a cheap amplification vector. A genuinely
// rotated-in key is still picked up within this window.
const JWKSRefreshCooldown = 1 * time.Minute
