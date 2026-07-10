package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/constants"
)

// RequestOptions configures a single Do call.
type RequestOptions struct {
	Query url.Values
	Body  any

	// RawBody, if non-nil, is sent verbatim with RawContentType instead of
	// JSON-marshaling Body — for non-JSON payloads such as a raw NDJSON/CSV
	// file upload. When set, Body is ignored.
	RawBody        []byte
	RawContentType string

	// Idempotent governs whether a 5xx is retried (GET/DELETE are; POST/PATCH
	// generally aren't, since the server may have already applied a mutation
	// before failing). 429 is always retried regardless of this flag.
	Idempotent bool
}

// buildPayload resolves opts into the raw bytes to send and the Content-Type
// to send them with. RawBody takes precedence over Body.
func buildPayload(opts RequestOptions) (payload []byte, contentType string, err error) {
	if opts.RawBody != nil {
		return opts.RawBody, opts.RawContentType, nil
	}
	if opts.Body == nil {
		return nil, "", nil
	}
	b, err := json.Marshal(opts.Body)
	if err != nil {
		return nil, "", fmt.Errorf("qeetid: marshal request body: %w", err)
	}
	return b, "application/json", nil
}

// newRequest builds one attempt's *http.Request: URL with query, auth/UA/
// custom headers, and body if present. Caller-supplied headers are applied
// first so the reserved headers below always win.
func (t *Transport) newRequest(ctx context.Context, method, path string, opts RequestOptions, payload []byte, contentType string) (*http.Request, error) {
	u := t.baseURL + path
	if len(opts.Query) > 0 {
		u += "?" + opts.Query.Encode()
	}

	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, fmt.Errorf("qeetid: build request: %w", err)
	}

	for k, vv := range t.headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set(constants.HeaderAuthorization, "ApiKey "+t.apiKey)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set(constants.HeaderUserAgent, t.userAgent)
	if contentType != "" {
		req.Header.Set(constants.HeaderContentType, contentType)
	}
	return req, nil
}
