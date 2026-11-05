package qeetid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Sentinel errors. An *APIError or *OAuthError matches the sentinels that fit
// its status and code, so callers can test for a class of failure with
// errors.Is without parsing codes:
//
//	if errors.Is(err, qeetid.ErrNotFound) { ... }
var (
	// ErrInvalidRequest matches a request the server refused as malformed or
	// invalid (HTTP 400 or 422).
	ErrInvalidRequest = errors.New("qeetid: invalid request")
	// ErrUnauthorized matches a request whose credentials are missing, wrong or
	// no longer valid (HTTP 401), and an OAuth invalid_client error.
	ErrUnauthorized = errors.New("qeetid: unauthorized")
	// ErrSessionRevoked matches a 401 with code auth.session_revoked: the
	// session the access token belongs to has ended. Sign the person in again;
	// retrying with the same token cannot succeed, and the client never does.
	ErrSessionRevoked = errors.New("qeetid: session revoked")
	// ErrPaymentRequired matches a request the tenant's plan does not cover
	// (HTTP 402), such as a feature outside the plan or a full seat count.
	ErrPaymentRequired = errors.New("qeetid: payment required")
	// ErrForbidden matches a request the credentials are not allowed to make
	// (HTTP 403).
	ErrForbidden = errors.New("qeetid: forbidden")
	// ErrNotFound matches a request for something that does not exist, or that
	// the caller may not see (HTTP 404).
	ErrNotFound = errors.New("qeetid: not found")
	// ErrConflict matches a request that conflicts with the current state
	// (HTTP 409), such as a duplicate or an idempotency key still in use.
	ErrConflict = errors.New("qeetid: conflict")
	// ErrRateLimited matches a request refused by rate limiting (HTTP 429).
	// APIError.RetryAfter says how long to wait.
	ErrRateLimited = errors.New("qeetid: rate limited")
	// ErrUnavailable matches a request the server could not take on now
	// (HTTP 503), for example because it is overloaded. APIError.RetryAfter
	// says how long to wait when the server said.
	ErrUnavailable = errors.New("qeetid: service unavailable")
	// ErrNoCredentials is returned, without sending the request, when an
	// authenticated method is called on a client created without WithAPIKey,
	// WithClientCredentials or WithTokenSource.
	ErrNoCredentials = errors.New("qeetid: no credentials configured")
)

// Error codes the client itself acts on, and others commonly worth handling.
// APIError.Code holds the code of any error; the server documents the full set
// per endpoint, and a shipped code never changes.
const (
	CodeSessionRevoked            = "auth.session_revoked"
	CodeUnauthorized              = "unauthorized"
	CodeForbidden                 = "forbidden"
	CodeNotFound                  = "not_found"
	CodeConflict                  = "conflict"
	CodeValidationFailed          = "validation_failed"
	CodeTooManyRequests           = "too_many_requests"
	CodeServerOverloaded          = "server.overloaded"
	CodeInternal                  = "internal"
	CodeIdempotencyKeyInvalid     = "idempotency.key_invalid"
	CodeIdempotencyKeyReused      = "idempotency.key_reused"
	CodeIdempotencyInProgress     = "idempotency.in_progress"
	CodeIdempotencyOutcomeUnknown = "idempotency.outcome_unknown"
)

// APIError is the error the client returns when the Qeet ID API answers with
// an error status. It carries the server's error envelope:
//
//	{"error": {"code": "...", "message": "...", "request_id": "...", ...}}
//
// Use errors.As to read it and errors.Is to match the sentinel errors above.
type APIError struct {
	// StatusCode is the HTTP status of the response.
	StatusCode int
	// Code is the stable machine code, such as "not_found" or
	// "auth.session_revoked". Branch on it, never on Message. It is empty when
	// the response carried no error envelope, such as a proxy's error page.
	Code string
	// Message is a human-readable message that is safe to show to end users.
	Message string
	// Detail is optional developer context. It is never localized.
	Detail string
	// TranslationKey is the key a client can localize Message with.
	TranslationKey string
	// Retryable is the server's hint that the request may be tried again
	// later. It describes the user's situation, such as a session that may be
	// renewed; the client decides on its own which requests it retries.
	Retryable bool
	// Fields lists per-field validation errors, typically on a 422.
	Fields []FieldError
	// Metadata carries structured, error-specific context, such as
	// {"remaining_attempts": 2}.
	Metadata map[string]any
	// RequestID identifies the request in Qeet ID's logs. Quote it when you
	// contact support.
	RequestID string
	// RetryAfter is how long the server asked the caller to wait, from the
	// Retry-After header, or zero when it did not say.
	RetryAfter time.Duration
}

// FieldError is a per-field validation error inside an APIError.
type FieldError struct {
	// Field names the request field, as the API spells it.
	Field string `json:"field"`
	// Code is a stable reason, such as "required" or "invalid_email".
	Code string `json:"code"`
	// Message is a human-readable explanation.
	Message string `json:"message"`
}

// Error implements the error interface. The text holds the status, code,
// message and request ID; it never holds credentials.
func (e *APIError) Error() string {
	var b strings.Builder
	b.WriteString("qeetid: ")
	b.WriteString(strconv.Itoa(e.StatusCode))
	if e.Code != "" {
		b.WriteString(" ")
		b.WriteString(e.Code)
	}
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	if e.RequestID != "" {
		b.WriteString(" (request ID ")
		b.WriteString(e.RequestID)
		b.WriteString(")")
	}
	return b.String()
}

// Is reports whether e belongs to the class of failure target names, so that
// errors.Is(err, ErrNotFound) and the like work.
func (e *APIError) Is(target error) bool {
	return matchStatus(target, e.StatusCode, e.Code)
}

// sentinelMatchers says which statuses and codes each sentinel error stands
// for. A map keeps the comparison by identity explicit: an Is method receives
// the sentinel itself, never a wrapped error.
var sentinelMatchers = map[error]func(status int, code string) bool{
	ErrInvalidRequest: func(s int, _ string) bool {
		return s == http.StatusBadRequest || s == http.StatusUnprocessableEntity
	},
	ErrUnauthorized:    func(s int, _ string) bool { return s == http.StatusUnauthorized },
	ErrSessionRevoked:  func(_ int, c string) bool { return c == CodeSessionRevoked },
	ErrPaymentRequired: func(s int, _ string) bool { return s == http.StatusPaymentRequired },
	ErrForbidden:       func(s int, _ string) bool { return s == http.StatusForbidden },
	ErrNotFound:        func(s int, _ string) bool { return s == http.StatusNotFound },
	ErrConflict:        func(s int, _ string) bool { return s == http.StatusConflict },
	ErrRateLimited:     func(s int, _ string) bool { return s == http.StatusTooManyRequests },
	ErrUnavailable:     func(s int, _ string) bool { return s == http.StatusServiceUnavailable },
}

func matchStatus(target error, status int, code string) bool {
	match, ok := sentinelMatchers[target]
	return ok && match(status, code)
}

// OAuthErrorCode is an OAuth 2.0 error code (RFC 6749 section 5.2, RFC 8628
// section 3.5 and the related specifications).
type OAuthErrorCode string

// OAuth error codes the Qeet ID OAuth endpoints answer with. The server may
// add codes; compare against these constants and treat others as failures.
const (
	OAuthInvalidRequest       OAuthErrorCode = "invalid_request"
	OAuthInvalidClient        OAuthErrorCode = "invalid_client"
	OAuthInvalidGrant         OAuthErrorCode = "invalid_grant"
	OAuthUnauthorizedClient   OAuthErrorCode = "unauthorized_client"
	OAuthUnsupportedGrantType OAuthErrorCode = "unsupported_grant_type"
	OAuthInvalidScope         OAuthErrorCode = "invalid_scope"
	OAuthInvalidTarget        OAuthErrorCode = "invalid_target"
	OAuthUnsupportedTokenType OAuthErrorCode = "unsupported_token_type"
	OAuthAccessDenied         OAuthErrorCode = "access_denied"
	OAuthAuthorizationPending OAuthErrorCode = "authorization_pending"
	OAuthSlowDown             OAuthErrorCode = "slow_down"
	OAuthExpiredToken         OAuthErrorCode = "expired_token"
	OAuthUnknownUserID        OAuthErrorCode = "unknown_user_id"
)

// OAuthError is the error the client returns when an OAuth endpoint answers
// in the RFC 6749 shape {"error": "...", "error_description": "..."}, as the
// token endpoints (/v1/oauth/token and /v1/oauth/token-code), device
// authorization, backchannel authentication and token revocation do.
//
// The device and backchannel polling errors authorization_pending and
// slow_down arrive as an *OAuthError too; see OAuthService.PollDeviceToken.
type OAuthError struct {
	// StatusCode is the HTTP status: 401 for invalid_client, 400 otherwise.
	StatusCode int
	// Code is the OAuth error code.
	Code OAuthErrorCode
	// Description is the server's human-readable explanation.
	Description string
	// URI points to a page about the error, when the server gives one.
	URI string
	// RequestID identifies the request in Qeet ID's logs.
	RequestID string
}

// Error implements the error interface.
func (e *OAuthError) Error() string {
	s := "qeetid: oauth " + string(e.Code)
	if e.Description != "" {
		s += ": " + e.Description
	}
	if e.RequestID != "" {
		s += " (request ID " + e.RequestID + ")"
	}
	return s
}

// Is reports whether e belongs to the class of failure target names. An
// invalid_client error matches ErrUnauthorized; any other OAuth error matches
// ErrInvalidRequest.
func (e *OAuthError) Is(target error) bool {
	if e.Code == OAuthInvalidClient {
		return matchStatus(target, http.StatusUnauthorized, "")
	}
	return matchStatus(target, e.StatusCode, "")
}

// errorFromResponse turns a non-2xx response into an *APIError or, for a body
// in the OAuth error shape, an *OAuthError.
func errorFromResponse(status int, header http.Header, body []byte) error {
	requestID := header.Get("X-Request-Id")
	retryAfter := parseRetryAfter(header.Get("Retry-After"), time.Now())

	var probe struct {
		Error json.RawMessage `json:"error"`
	}
	if json.Unmarshal(body, &probe) == nil && len(probe.Error) > 0 {
		switch probe.Error[0] {
		case '"':
			var o struct {
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
				ErrorURI         string `json:"error_uri"`
			}
			if json.Unmarshal(body, &o) == nil && o.Error != "" {
				return &OAuthError{
					StatusCode:  status,
					Code:        OAuthErrorCode(o.Error),
					Description: o.ErrorDescription,
					URI:         o.ErrorURI,
					RequestID:   requestID,
				}
			}
		case '{':
			var env struct {
				Code           string         `json:"code"`
				Message        string         `json:"message"`
				Detail         string         `json:"detail"`
				TranslationKey string         `json:"translation_key"`
				Retryable      bool           `json:"retryable"`
				Fields         []FieldError   `json:"fields"`
				Metadata       map[string]any `json:"metadata"`
				RequestID      string         `json:"request_id"`
			}
			if json.Unmarshal(probe.Error, &env) == nil {
				if env.RequestID != "" {
					requestID = env.RequestID
				}
				if retryAfter == 0 {
					retryAfter = retryAfterFromMetadata(env.Metadata)
				}
				return &APIError{
					StatusCode:     status,
					Code:           env.Code,
					Message:        env.Message,
					Detail:         env.Detail,
					TranslationKey: env.TranslationKey,
					Retryable:      env.Retryable,
					Fields:         env.Fields,
					Metadata:       env.Metadata,
					RequestID:      requestID,
					RetryAfter:     retryAfter,
				}
			}
		}
	}
	return &APIError{
		StatusCode: status,
		Message:    http.StatusText(status),
		RequestID:  requestID,
		RetryAfter: retryAfter,
	}
}

// parseRetryAfter reads a Retry-After header given in seconds or as an HTTP
// date. It returns zero when the header is absent or unusable.
func parseRetryAfter(v string, now time.Time) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if secs, err := strconv.ParseInt(v, 10, 64); err == nil {
		if secs <= 0 {
			return 0
		}
		if secs > int64(24*time.Hour/time.Second) {
			secs = int64(24 * time.Hour / time.Second)
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := t.Sub(now); d > 0 {
			return d
		}
	}
	return 0
}

func retryAfterFromMetadata(m map[string]any) time.Duration {
	if v, ok := m["retry_after"].(float64); ok && v > 0 && v < 86400 {
		return time.Duration(v * float64(time.Second))
	}
	return 0
}

// requestError wraps a failure to send a request or read its answer, keeping
// the cause available to errors.Is and errors.As (for example
// context.DeadlineExceeded or a *net.OpError).
func requestError(method, path string, err error) error {
	return fmt.Errorf("qeetid: %s %s: %w", method, path, err)
}
