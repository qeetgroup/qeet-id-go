package transport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/qeetgroup/qeet-id-go/internal/constants"
)

// DoForm executes a form-encoded request — used by OAuth token/introspection/
// revocation endpoints, which authenticate via OIDC client credentials
// (HTTP Basic, optional) rather than the management API's ApiKey header. No
// retry: these are one-shot grant/introspection calls, not idempotent reads.
func (t *Transport) DoForm(ctx context.Context, method, path string, form url.Values, basicUser, basicPass string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, t.baseURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("qeetid: build request: %w", err)
	}
	req.Header.Set(constants.HeaderContentType, "application/x-www-form-urlencoded")
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set(constants.HeaderUserAgent, t.userAgent)
	if basicUser != "" {
		req.SetBasicAuth(basicUser, basicPass)
	}

	res, err := t.hc.Do(req)
	if err != nil {
		return &Error{Code: "network_error", Message: err.Error()}
	}
	defer res.Body.Close()
	return readResponse(res, out)
}
