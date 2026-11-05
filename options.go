package qeetid

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

// Default client settings. Each can be changed with an Option.
const (
	// DefaultBaseURL is the Qeet ID API the client talks to unless WithBaseURL
	// says otherwise.
	DefaultBaseURL = "https://api.id.qeet.in"
	// DefaultTimeout bounds each HTTP attempt. Retries get a fresh timeout; the
	// context passed to a method bounds the call as a whole.
	DefaultTimeout = 30 * time.Second
	// DefaultMaxRetries is how many times a failed request is retried when the
	// failure is retryable and retrying is safe (see WithMaxRetries).
	DefaultMaxRetries = 2
	// DefaultMinRetryDelay and DefaultMaxRetryDelay bound the exponential
	// backoff between retries when the server does not send Retry-After.
	DefaultMinRetryDelay = 500 * time.Millisecond
	DefaultMaxRetryDelay = 8 * time.Second
	// MaxRetryAfter is the longest Retry-After the client waits for on its own.
	// A longer one is returned to the caller as an error that carries it.
	MaxRetryAfter = 60 * time.Second
)

// Option configures a Client. Pass options to NewClient.
type Option interface {
	applyOption(*clientConfig)
}

type optionFunc func(*clientConfig)

func (f optionFunc) applyOption(c *clientConfig) { f(c) }

type clientConfig struct {
	baseURL         string
	httpClient      *http.Client
	timeout         time.Duration
	maxRetries      int
	minRetryDelay   time.Duration
	maxRetryDelay   time.Duration
	logger          *slog.Logger
	userAgentSuffix string
	headers         http.Header
	credentials     []credentialSpec
	errs            []error
}

// credentialSpec records one credential option; NewClient refuses more than one.
type credentialSpec struct {
	name   string
	apiKey string
	ts     TokenSource
	cc     *clientCredentialsConfig
}

type clientCredentialsConfig struct {
	clientID     string
	clientSecret string
	scopes       []string
}

func defaultConfig() clientConfig {
	return clientConfig{
		baseURL:       DefaultBaseURL,
		timeout:       DefaultTimeout,
		maxRetries:    DefaultMaxRetries,
		minRetryDelay: DefaultMinRetryDelay,
		maxRetryDelay: DefaultMaxRetryDelay,
	}
}

// WithAPIKey authenticates every request with an API key, sent as
// "Authorization: ApiKey <key>". Keep API keys on servers; never ship one to a
// browser or a mobile app.
func WithAPIKey(key string) Option {
	return optionFunc(func(c *clientConfig) {
		if strings.TrimSpace(key) == "" {
			c.errs = append(c.errs, errors.New("qeetid: WithAPIKey: the key is empty"))
			return
		}
		c.credentials = append(c.credentials, credentialSpec{name: "WithAPIKey", apiKey: key})
	})
}

// WithClientCredentials authenticates every request with an access token that
// the client obtains for a service principal with the OAuth client credentials
// grant (POST /v1/oauth/token). The client caches the token and fetches a new
// one shortly before it expires. scopes, when given, are sent with the grant,
// but the server ignores them today: the token carries every scope the service
// principal holds.
func WithClientCredentials(clientID, clientSecret string, scopes ...string) Option {
	return optionFunc(func(c *clientConfig) {
		if clientID == "" || clientSecret == "" {
			c.errs = append(c.errs, errors.New("qeetid: WithClientCredentials: the client ID and secret are required"))
			return
		}
		c.credentials = append(c.credentials, credentialSpec{name: "WithClientCredentials", cc: &clientCredentialsConfig{
			clientID: clientID, clientSecret: clientSecret, scopes: append([]string(nil), scopes...),
		}})
	})
}

// WithTokenSource authenticates every request with a bearer token from ts, sent
// as "Authorization: Bearer <token>". Use it for tokens the client does not
// obtain itself, such as an agent's token or one from golang.org/x/oauth2.
func WithTokenSource(ts TokenSource) Option {
	return optionFunc(func(c *clientConfig) {
		if ts == nil {
			c.errs = append(c.errs, errors.New("qeetid: WithTokenSource: the token source is nil"))
			return
		}
		c.credentials = append(c.credentials, credentialSpec{name: "WithTokenSource", ts: ts})
	})
}

// WithBaseURL sets the API the client talks to (default DefaultBaseURL). The
// URL must use https; http is accepted only for a loopback host such as
// localhost, so that credentials never cross a network unencrypted.
func WithBaseURL(rawURL string) Option {
	return optionFunc(func(c *clientConfig) { c.baseURL = rawURL })
}

// WithHTTPClient sets the HTTP client used for every request, for example one
// with a custom transport, proxy or TLS settings. The client's own Timeout,
// if any, applies on top of WithTimeout.
func WithHTTPClient(hc *http.Client) Option {
	return optionFunc(func(c *clientConfig) { c.httpClient = hc })
}

// WithTimeout sets how long each HTTP attempt may take (default
// DefaultTimeout). A value of zero or less removes the per-attempt limit; the
// context passed to each method still applies.
func WithTimeout(d time.Duration) Option {
	return optionFunc(func(c *clientConfig) { c.timeout = d })
}

// WithMaxRetries sets how many times a failed request is retried (default
// DefaultMaxRetries). Zero turns retries off.
//
// The client retries only when it is safe: a GET or DELETE, or a create call
// sent with an idempotency key, after a network error or a 408, 500, 502, 503
// or 504 answer; and any request the server refused before acting on it, which
// it signals with 429 or 503 and a Retry-After header. A 401 is never retried.
func WithMaxRetries(n int) Option {
	return optionFunc(func(c *clientConfig) {
		if n < 0 {
			n = 0
		}
		c.maxRetries = n
	})
}

// WithRetryDelay sets the bounds of the exponential backoff between retries
// (defaults DefaultMinRetryDelay and DefaultMaxRetryDelay). The delay before
// retry n is a random value between half and all of min×2ⁿ, capped at max. A
// Retry-After header from the server takes precedence.
func WithRetryDelay(minDelay, maxDelay time.Duration) Option {
	return optionFunc(func(c *clientConfig) {
		if minDelay <= 0 || maxDelay < minDelay {
			c.errs = append(c.errs, fmt.Errorf("qeetid: WithRetryDelay: need 0 < min <= max, got %v and %v", minDelay, maxDelay))
			return
		}
		c.minRetryDelay, c.maxRetryDelay = minDelay, maxDelay
	})
}

// WithLogger logs each HTTP attempt at debug level and each retry at warn level
// to l. Log records carry the method, path, status, duration, attempt and the
// server's request ID; they never carry credentials, headers or bodies. The
// default is no logging.
func WithLogger(l *slog.Logger) Option {
	return optionFunc(func(c *clientConfig) { c.logger = l })
}

// WithUserAgentSuffix appends s, such as "myapp/1.4.2", to the User-Agent the
// client sends, so that requests can be traced to your application.
func WithUserAgentSuffix(s string) Option {
	return optionFunc(func(c *clientConfig) {
		if strings.ContainsAny(s, "\r\n") {
			c.errs = append(c.errs, errors.New("qeetid: WithUserAgentSuffix: the suffix contains a line break"))
			return
		}
		c.userAgentSuffix = strings.TrimSpace(s)
	})
}

// WithHeader adds a header to every request, for example a tracing header. It
// cannot set a header the client manages itself: Authorization, Accept,
// Content-Type, Idempotency-Key or User-Agent.
func WithHeader(key, value string) Option {
	return optionFunc(func(c *clientConfig) {
		if err := checkHeader("WithHeader", key, value); err != nil {
			c.errs = append(c.errs, err)
			return
		}
		if c.headers == nil {
			c.headers = http.Header{}
		}
		c.headers.Add(key, value)
	})
}

// RequestOption configures a single call. Pass request options as the last
// arguments of any method.
type RequestOption interface {
	applyRequestOption(*requestConfig)
}

type requestOptionFunc func(*requestConfig)

func (f requestOptionFunc) applyRequestOption(c *requestConfig) { f(c) }

type requestConfig struct {
	idempotencyKey string
	headers        http.Header
	errs           []error
}

// WithIdempotencyKey sends key in the Idempotency-Key header, so that sending
// the same request again with the same key within 24 hours returns the first
// answer instead of acting twice. key must be 1 to 255 printable ASCII
// characters. Only the create calls whose documentation says so accept it;
// any other call returns an error without sending the request.
//
// Those calls send a random key of their own when you give none, which makes
// the client's automatic retries safe. Give your own key to make retries across
// process restarts or queue redeliveries safe as well.
func WithIdempotencyKey(key string) RequestOption {
	return requestOptionFunc(func(c *requestConfig) {
		if !validIdempotencyKey(key) {
			c.errs = append(c.errs, errors.New("qeetid: WithIdempotencyKey: the key must be 1 to 255 printable ASCII characters"))
			return
		}
		c.idempotencyKey = key
	})
}

// WithRequestHeader adds a header to one request. The same headers as for
// WithHeader are reserved.
func WithRequestHeader(key, value string) RequestOption {
	return requestOptionFunc(func(c *requestConfig) {
		if err := checkHeader("WithRequestHeader", key, value); err != nil {
			c.errs = append(c.errs, err)
			return
		}
		if c.headers == nil {
			c.headers = http.Header{}
		}
		c.headers.Add(key, value)
	})
}

var reservedHeaders = map[string]bool{
	"Authorization":   true,
	"Accept":          true,
	"Content-Type":    true,
	"Content-Length":  true,
	"Host":            true,
	"Idempotency-Key": true,
	"User-Agent":      true,
}

func checkHeader(option, key, value string) error {
	canonical := textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(key))
	switch {
	case canonical == "":
		return fmt.Errorf("qeetid: %s: the header name is empty", option)
	case reservedHeaders[canonical]:
		return fmt.Errorf("qeetid: %s: the client sets %s itself", option, canonical)
	case strings.ContainsAny(key+value, "\r\n"):
		return fmt.Errorf("qeetid: %s: header %s contains a line break", option, canonical)
	}
	return nil
}

// validIdempotencyKey mirrors the server's rule: 1 to 255 printable ASCII
// characters, without spaces.
func validIdempotencyKey(key string) bool {
	if key == "" || len(key) > 255 {
		return false
	}
	for i := range len(key) {
		if key[i] < 0x21 || key[i] > 0x7e {
			return false
		}
	}
	return true
}
