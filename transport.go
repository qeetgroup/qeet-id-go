package qeetid

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	mrand "math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// maxResponseBytes caps how much of a response body the client reads.
const maxResponseBytes = 32 << 20

// request describes one API call. Generated service methods fill it in; do
// sends it.
type request struct {
	method   string
	path     string   // may hold {name} placeholders, filled from pathArgs in order
	pathArgs []string // unescaped values; empty values are refused
	query    url.Values

	body        any        // sent as JSON when non-nil
	form        url.Values // sent as application/x-www-form-urlencoded when non-nil
	rawBody     []byte     // sent as is with contentType when non-nil
	contentType string

	idempotent bool // the endpoint honours the Idempotency-Key header
	public     bool // the endpoint takes no credentials

	// OAuth client authentication (HTTP Basic) for the OAuth endpoints. When
	// clientID is set the client's own credential is not sent.
	clientID, clientSecret string
	// bearer, when set, is sent as the bearer token instead of the client's
	// credential (for example a person's token for the userinfo endpoint).
	bearer string
}

// call sends req and decodes a successful JSON answer into a new T.
func call[T any](ctx context.Context, c *Client, req *request, opts []RequestOption) (*T, error) {
	var out T
	if err := c.do(ctx, req, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// do sends req, retrying when that is safe, and decodes a successful answer
// into out (which may be nil). An error answer becomes an *APIError or
// *OAuthError; a failure to send becomes a wrapped error.
func (c *Client) do(ctx context.Context, req *request, out any, opts []RequestOption) error {
	if ctx == nil {
		return errors.New("qeetid: nil context")
	}
	var rc requestConfig
	for _, o := range opts {
		if o != nil {
			o.applyRequestOption(&rc)
		}
	}
	if err := errors.Join(rc.errs...); err != nil {
		return err
	}
	path, err := expandPath(req.path, req.pathArgs)
	if err != nil {
		return err
	}
	if rc.idempotencyKey != "" && !req.idempotent {
		return fmt.Errorf("qeetid: %s %s does not accept an idempotency key", req.method, req.path)
	}
	if !req.public && req.clientID == "" && req.bearer == "" && c.cred == nil {
		return ErrNoCredentials
	}
	idemKey := rc.idempotencyKey
	if req.idempotent && idemKey == "" && c.maxRetries > 0 {
		idemKey = newIdempotencyKey()
	}
	payload, contentType, err := encodeBody(req)
	if err != nil {
		return err
	}
	target := c.baseURL + path
	if len(req.query) > 0 {
		target += "?" + req.query.Encode()
	}
	replaySafe := safeToReplay(req.method) || idemKey != ""

	for attempt := 0; ; attempt++ {
		start := time.Now()
		status, header, body, sendErr := c.send(ctx, req, target, payload, contentType, idemKey, rc.headers)
		requestID := ""
		if header != nil {
			requestID = header.Get("X-Request-Id")
		}
		c.logAttempt(ctx, req.method, path, status, time.Since(start), attempt, requestID, sendErr)

		if sendErr != nil {
			if ctx.Err() != nil || !replaySafe || attempt >= c.maxRetries || !retryableSendError(sendErr) {
				return requestError(req.method, path, sendErr)
			}
			if !c.sleep(ctx, c.backoff(attempt), req.method, path, attempt, "network error") {
				return requestError(req.method, path, sendErr)
			}
			continue
		}

		if status >= 200 && status < 300 {
			if out == nil || len(bytes.TrimSpace(body)) == 0 || status == http.StatusNoContent {
				return nil
			}
			if err := json.Unmarshal(body, out); err != nil {
				return requestError(req.method, path, fmt.Errorf("decode response: %w", err))
			}
			return nil
		}

		apiErr := errorFromResponse(status, header, body)
		wait, retry := c.retryPlan(apiErr, status, replaySafe, idemKey != "", attempt)
		if !retry {
			return apiErr
		}
		if !c.sleep(ctx, wait, req.method, path, attempt, strconv.Itoa(status)) {
			return apiErr
		}
	}
}

// send makes one HTTP attempt and reads the whole answer.
func (c *Client) send(ctx context.Context, req *request, target string, payload []byte, contentType, idemKey string, extra http.Header) (int, http.Header, []byte, error) {
	actx := ctx
	if c.timeout > 0 {
		var cancel context.CancelFunc
		actx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	hreq, err := http.NewRequestWithContext(actx, req.method, target, body)
	if err != nil {
		return 0, nil, nil, err
	}
	for k, vs := range c.headers {
		for _, v := range vs {
			hreq.Header.Add(k, v)
		}
	}
	for k, vs := range extra {
		for _, v := range vs {
			hreq.Header.Add(k, v)
		}
	}
	hreq.Header.Set("Accept", "application/json")
	hreq.Header.Set("User-Agent", c.userAgent)
	if contentType != "" {
		hreq.Header.Set("Content-Type", contentType)
	}
	if idemKey != "" {
		hreq.Header.Set("Idempotency-Key", idemKey)
	}
	switch {
	case req.clientID != "":
		// RFC 6749 §2.3.1: each value is form-encoded before Basic encoding,
		// and the server decodes it. The URL-safe secrets it issues are
		// unchanged by this.
		hreq.SetBasicAuth(url.QueryEscape(req.clientID), url.QueryEscape(req.clientSecret))
	case req.bearer != "":
		hreq.Header.Set("Authorization", "Bearer "+req.bearer)
	case !req.public && c.cred != nil:
		if err := c.cred.authorize(ctx, hreq); err != nil {
			return 0, nil, nil, &credentialError{err: err}
		}
	}

	resp, err := c.hc.Do(hreq)
	if err != nil {
		var uerr *url.Error
		if errors.As(err, &uerr) {
			// The URL is already in the error text the caller sees; keep the cause.
			err = uerr.Err
		}
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("read response: %w", err)
	}
	if len(data) > maxResponseBytes {
		return 0, nil, nil, fmt.Errorf("response larger than %d bytes", maxResponseBytes)
	}
	return resp.StatusCode, resp.Header, data, nil
}

// retryPlan decides whether an error answer is retried and after how long.
//
// A request the server refused before acting on it, which it marks with 429 or
// 503 and a Retry-After header, is retried whatever its method. Any other
// retryable status is retried only when the request is safe to replay: a GET,
// HEAD or DELETE, or a request with an idempotency key. A 401 never is,
// including auth.session_revoked: the session has ended and a retry cannot
// succeed.
func (c *Client) retryPlan(err error, status int, replaySafe, keyed bool, attempt int) (time.Duration, bool) {
	if attempt >= c.maxRetries {
		return 0, false
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return 0, false // OAuth errors are answers, not transient failures
	}
	if status == http.StatusUnauthorized || apiErr.Code == CodeSessionRevoked {
		return 0, false
	}
	if apiErr.RetryAfter > MaxRetryAfter {
		return 0, false
	}
	switch {
	case (status == http.StatusTooManyRequests || status == http.StatusServiceUnavailable) && apiErr.RetryAfter > 0:
		return apiErr.RetryAfter, true
	case status == http.StatusConflict && apiErr.Code == CodeIdempotencyInProgress && keyed:
		return c.backoff(attempt), true
	case replaySafe && retryableStatus(status):
		if apiErr.RetryAfter > 0 {
			return apiErr.RetryAfter, true
		}
		return c.backoff(attempt), true
	}
	return 0, false
}

func retryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout, http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

// retryableSendError reports whether a failure to send is worth retrying: a
// network failure, but not a cancelled context or a credential that could not
// be obtained.
func retryableSendError(err error) bool {
	var credErr *credentialError
	return !errors.Is(err, context.Canceled) && !errors.As(err, &credErr)
}

// credentialError reports that the credential for a request could not be
// obtained, for example because the token endpoint refused the client.
type credentialError struct{ err error }

func (e *credentialError) Error() string { return "obtain credentials: " + e.err.Error() }
func (e *credentialError) Unwrap() error { return e.err }

// safeToReplay reports whether a request with this method can be sent again
// without risk of acting twice.
func safeToReplay(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodDelete:
		return true
	}
	return false
}

// backoff returns the delay before retry attempt+1: a random value between
// half and all of min×2^attempt, capped at max.
func (c *Client) backoff(attempt int) time.Duration {
	d := c.maxRetryDelay
	if attempt < 30 {
		if e := c.minRetryDelay << attempt; e > 0 && e < d {
			d = e
		}
	}
	half := d / 2
	return half + mrand.N(half+1) // #nosec G404 -- jitter, not a secret
}

// sleep waits d before the next attempt. It returns false without waiting when
// the context ends first or its deadline would pass during the wait.
func (c *Client) sleep(ctx context.Context, d time.Duration, method, path string, attempt int, reason string) bool {
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < d {
		return false
	}
	if c.logger != nil {
		c.logger.LogAttrs(ctx, slog.LevelWarn, "qeetid: retrying request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("attempt", attempt+1),
			slog.String("reason", reason),
			slog.Duration("wait", d),
		)
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func (c *Client) logAttempt(ctx context.Context, method, path string, status int, took time.Duration, attempt int, requestID string, err error) {
	if c.logger == nil || !c.logger.Enabled(ctx, slog.LevelDebug) {
		return
	}
	attrs := []slog.Attr{
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.Duration("duration", took),
		slog.Int("attempt", attempt+1),
	}
	if requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	c.logger.LogAttrs(ctx, slog.LevelDebug, "qeetid: request", attrs...)
}

// expandPath fills a path template's {name} placeholders with args, in order,
// escaping each. An empty argument is refused, because it would address a
// different endpoint.
func expandPath(template string, args []string) (string, error) {
	if !strings.Contains(template, "{") {
		return template, nil
	}
	var b strings.Builder
	rest := template
	for i := 0; ; i++ {
		open := strings.IndexByte(rest, '{')
		if open < 0 {
			b.WriteString(rest)
			break
		}
		closing := strings.IndexByte(rest[open:], '}')
		if closing < 0 {
			return "", fmt.Errorf("qeetid: malformed path template %q", template)
		}
		name := rest[open+1 : open+closing]
		if i >= len(args) {
			return "", fmt.Errorf("qeetid: missing value for %s in %q", name, template)
		}
		if strings.TrimSpace(args[i]) == "" {
			return "", fmt.Errorf("qeetid: %s must not be empty", name)
		}
		b.WriteString(rest[:open])
		b.WriteString(url.PathEscape(args[i]))
		rest = rest[open+closing+1:]
	}
	return b.String(), nil
}

func encodeBody(req *request) ([]byte, string, error) {
	switch {
	case req.rawBody != nil:
		return req.rawBody, req.contentType, nil
	case req.form != nil:
		return []byte(req.form.Encode()), "application/x-www-form-urlencoded", nil
	case req.body != nil:
		b, err := json.Marshal(req.body)
		if err != nil {
			return nil, "", fmt.Errorf("qeetid: encode request body: %w", err)
		}
		return b, "application/json", nil
	}
	return nil, "", nil
}

// newIdempotencyKey returns a random version 4 UUID.
func newIdempotencyKey() string {
	var b [16]byte
	_, _ = rand.Read(b[:]) // crypto/rand.Read never fails on supported platforms
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	h := hex.EncodeToString(b[:])
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}
