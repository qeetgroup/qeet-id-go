// Package transport is the shared HTTP execution engine for every resource
// in the root package: auth-header injection, JSON (de)serialisation, typed
// errors, timeouts, and retry/backoff on 429/5xx. It is the single owner of
// the concrete Error type — the root package re-exports it as a type alias
// (see errors.go) so that internal/auth (and every resource file) can
// construct/return errors without importing the root package, which would
// create an import cycle (root imports every resource file; a resource file
// or internal/auth importing root back would cycle).
package transport

import (
	"net/http"
	"strings"
	"time"

	"github.com/qeetgroup/qeet-id-go/internal/constants"
)

// Options configures a Transport. APIKey is required.
type Options struct {
	APIKey     string
	BaseURL    string        // default constants.DefaultBaseURL
	HTTPClient *http.Client  // if nil, built from Timeout
	Timeout    time.Duration // used only when HTTPClient is nil; default constants.DefaultTimeout
	MaxRetries int           // default constants.DefaultMaxRetries
	// Headers are sent on every request — an escape hatch for tracing
	// headers or forward-compatible options. They cannot override the
	// reserved Authorization/Accept/Content-Type/User-Agent headers.
	Headers http.Header
	// UserAgent, if set, is prepended to the SDK's own User-Agent string.
	UserAgent string
	// Logger is an optional per-request observability hook. Nil is a no-op.
	Logger Logger
}

// Transport is the shared HTTP client every resource's *Service embeds.
// Safe for concurrent use.
type Transport struct {
	apiKey     string
	baseURL    string
	hc         *http.Client
	maxRetries int
	userAgent  string
	headers    http.Header
	logger     Logger
}

// New builds a Transport. Panics are never used for config errors — an
// empty APIKey surfaces as an authentication failure on the first request,
// matching how every other config mistake is discovered.
func New(opts Options) *Transport {
	base := opts.BaseURL
	if base == "" {
		base = constants.DefaultBaseURL
	}
	base = strings.TrimRight(base, "/")

	hc := opts.HTTPClient
	if hc == nil {
		timeout := opts.Timeout
		if timeout <= 0 {
			timeout = constants.DefaultTimeout
		}
		hc = &http.Client{Timeout: timeout}
	}
	maxRetries := opts.MaxRetries
	if maxRetries <= 0 {
		maxRetries = constants.DefaultMaxRetries
	}
	ua := UserAgent()
	if opts.UserAgent != "" {
		ua = opts.UserAgent + " " + ua
	}
	return &Transport{
		apiKey:     opts.APIKey,
		baseURL:    base,
		hc:         hc,
		maxRetries: maxRetries,
		userAgent:  ua,
		headers:    opts.Headers,
		logger:     opts.Logger,
	}
}

// BaseURL returns the configured (trimmed) base URL — used by resources
// that need to derive a related endpoint (Sessions' JWKS URL, Discovery).
func (t *Transport) BaseURL() string { return t.baseURL }

// HTTPClient returns the underlying *http.Client — used by resources that
// bypass the ApiKey-authed request path (Sessions' JWKS fetch, OAuth's
// form-encoded grants with Basic auth).
func (t *Transport) HTTPClient() *http.Client { return t.hc }
