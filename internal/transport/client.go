package transport

import (
	"context"
	"net/http"
	"time"
)

// Do executes a request with auth, JSON (de)serialisation, typed errors, and
// retry/backoff on 429/5xx, decoding a successful response into out (which
// may be nil). This is the one method every resource's *Service calls.
func (t *Transport) Do(ctx context.Context, method, path string, opts RequestOptions, out any) error {
	payload, contentType, err := buildPayload(opts)
	if err != nil {
		return err
	}

	for attempt := 0; ; attempt++ {
		start := time.Now()
		req, err := t.newRequest(ctx, method, path, opts, payload, contentType)
		if err != nil {
			return err
		}

		res, err := t.hc.Do(req)
		if err != nil {
			if opts.Idempotent && attempt < t.maxRetries {
				sleep(ctx, backoff(attempt))
				continue
			}
			t.logRequest(method, path, 0, start, "")
			return &Error{Code: "network_error", Message: err.Error()}
		}

		if shouldRetry(res.StatusCode, opts.Idempotent) && attempt < t.maxRetries {
			wait := retryAfterDuration(res)
			res.Body.Close()
			if wait <= 0 {
				wait = backoff(attempt)
			}
			t.logRequest(method, path, res.StatusCode, start, res.Header.Get("X-Request-Id"))
			sleep(ctx, wait)
			continue
		}

		defer res.Body.Close()
		requestID := res.Header.Get("X-Request-Id")
		err = readResponse(res, out)
		t.logRequest(method, path, res.StatusCode, start, requestID)
		return err
	}
}

// Get is sugar for Do(ctx, http.MethodGet, ...) with Idempotent implied.
func (t *Transport) Get(ctx context.Context, path string, query RequestOptions, out any) error {
	query.Idempotent = true
	return t.Do(ctx, http.MethodGet, path, query, out)
}

// Post is sugar for Do(ctx, http.MethodPost, ...).
func (t *Transport) Post(ctx context.Context, path string, body any, out any) error {
	return t.Do(ctx, http.MethodPost, path, RequestOptions{Body: body}, out)
}

// Patch is sugar for Do(ctx, http.MethodPatch, ...).
func (t *Transport) Patch(ctx context.Context, path string, body any, out any) error {
	return t.Do(ctx, http.MethodPatch, path, RequestOptions{Body: body}, out)
}

// Delete is sugar for Do(ctx, http.MethodDelete, ...) with Idempotent implied.
func (t *Transport) Delete(ctx context.Context, path string, out any) error {
	return t.Do(ctx, http.MethodDelete, path, RequestOptions{Idempotent: true}, out)
}

// Put is sugar for Do(ctx, http.MethodPut, ...).
func (t *Transport) Put(ctx context.Context, path string, body any, out any) error {
	return t.Do(ctx, http.MethodPut, path, RequestOptions{Body: body}, out)
}
