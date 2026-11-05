package qeetid

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestAPIErrorFromEnvelope(t *testing.T) {
	body := `{"error":{"code":"validation_failed","message":"One or more fields are invalid.",
		"detail":"email","translation_key":"validation_failed","retryable":false,
		"fields":[{"field":"email","code":"invalid_email","message":"Must be a valid email address."}],
		"metadata":{"limit":5},"request_id":"req-body"}}`
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusUnprocessableEntity, body, "X-Request-Id", "req-header")}}
	c := newFakeClient(t, f)

	_, err := c.Users.Create(context.Background(), CreateUserParams{TenantID: "t1", Email: "not-an-email"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %T %v, want *APIError", err, err)
	}
	if apiErr.StatusCode != http.StatusUnprocessableEntity || apiErr.Code != CodeValidationFailed || apiErr.RequestID != "req-body" ||
		apiErr.Detail != "email" || len(apiErr.Fields) != 1 || apiErr.Fields[0].Code != "invalid_email" ||
		apiErr.Metadata["limit"] != float64(5) || apiErr.TranslationKey != "validation_failed" {
		t.Errorf("APIError = %+v", apiErr)
	}
	if !errors.Is(err, ErrInvalidRequest) || errors.Is(err, ErrNotFound) {
		t.Errorf("errors.Is mismatch for %v", err)
	}
	want := "qeetid: 422 validation_failed: One or more fields are invalid. (request ID req-body)"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestAPIErrorSentinels(t *testing.T) {
	cases := []struct {
		status int
		code   string
		want   []error
		not    []error
	}{
		{400, "bad_request", []error{ErrInvalidRequest}, []error{ErrUnauthorized}},
		{401, "unauthorized", []error{ErrUnauthorized}, []error{ErrSessionRevoked}},
		{401, "auth.session_revoked", []error{ErrUnauthorized, ErrSessionRevoked}, nil},
		{402, "billing.plan_limit", []error{ErrPaymentRequired}, nil},
		{403, "forbidden", []error{ErrForbidden}, []error{ErrUnauthorized}},
		{404, "not_found", []error{ErrNotFound}, nil},
		{409, "conflict", []error{ErrConflict}, nil},
		{422, "validation_failed", []error{ErrInvalidRequest}, nil},
		{429, "too_many_requests", []error{ErrRateLimited}, nil},
		{503, "server.overloaded", []error{ErrUnavailable}, []error{ErrRateLimited}},
	}
	for _, tc := range cases {
		err := error(&APIError{StatusCode: tc.status, Code: tc.code})
		for _, s := range tc.want {
			if !errors.Is(err, s) {
				t.Errorf("%d %s: errors.Is(%v) = false", tc.status, tc.code, s)
			}
		}
		for _, s := range tc.not {
			if errors.Is(err, s) {
				t.Errorf("%d %s: errors.Is(%v) = true", tc.status, tc.code, s)
			}
		}
	}
}

func TestAPIErrorWithoutEnvelope(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusBadGateway, `<html>bad gateway</html>`, "X-Request-Id", "req-1")}}
	c := newFakeClient(t, f, WithMaxRetries(0))
	_, err := c.Users.Get(context.Background(), "u1")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadGateway || apiErr.Code != "" || apiErr.Message != "Bad Gateway" || apiErr.RequestID != "req-1" {
		t.Fatalf("err = %#v", err)
	}
}

func TestOAuthErrorShape(t *testing.T) {
	cases := []struct {
		status int
		body   string
		code   OAuthErrorCode
		is     error
	}{
		{400, `{"error":"invalid_request","error_description":"grant_type is required."}`, OAuthInvalidRequest, ErrInvalidRequest},
		{401, `{"error":"invalid_client","error_description":"The client could not be authenticated."}`, OAuthInvalidClient, ErrUnauthorized},
		{400, `{"error":"unauthorized_client","error_description":"The service principal is disabled."}`, OAuthUnauthorizedClient, ErrInvalidRequest},
		{400, `{"error":"unsupported_grant_type","error_description":"Only client_credentials is supported."}`, OAuthUnsupportedGrantType, ErrInvalidRequest},
		{400, `{"error":"invalid_grant","error_description":"The code was already used."}`, OAuthInvalidGrant, ErrInvalidRequest},
		{400, `{"error":"invalid_scope"}`, OAuthInvalidScope, ErrInvalidRequest},
		{400, `{"error":"invalid_target","error_description":"resource must be an absolute URI without a fragment"}`, OAuthInvalidTarget, ErrInvalidRequest},
		{400, `{"error":"authorization_pending"}`, OAuthAuthorizationPending, ErrInvalidRequest},
		{400, `{"error":"slow_down"}`, OAuthSlowDown, ErrInvalidRequest},
		{400, `{"error":"expired_token"}`, OAuthExpiredToken, ErrInvalidRequest},
		{400, `{"error":"access_denied"}`, OAuthAccessDenied, ErrInvalidRequest},
	}
	for _, tc := range cases {
		err := errorFromResponse(tc.status, http.Header{"X-Request-Id": {"req-9"}}, []byte(tc.body))
		var oe *OAuthError
		if !errors.As(err, &oe) {
			t.Fatalf("%s: err = %T, want *OAuthError", tc.code, err)
		}
		if oe.Code != tc.code || oe.StatusCode != tc.status || oe.RequestID != "req-9" {
			t.Errorf("%s: OAuthError = %+v", tc.code, oe)
		}
		if !errors.Is(err, tc.is) {
			t.Errorf("%s: errors.Is(%v) = false", tc.code, tc.is)
		}
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Errorf("%s: an OAuth error must not be an *APIError", tc.code)
		}
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	cases := map[string]time.Duration{
		"":                              0,
		"3":                             3 * time.Second,
		"0":                             0,
		"-5":                            0,
		"soon":                          0,
		"Thu, 05 Nov 2026 10:00:30 GMT": 30 * time.Second,
		"Thu, 05 Nov 2026 09:00:00 GMT": 0,
	}
	for in, want := range cases {
		if got := parseRetryAfter(in, now); got != want {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestRetryAfterFromBodyMetadata(t *testing.T) {
	err := errorFromResponse(503, http.Header{}, []byte(`{"error":{"code":"server.overloaded","message":"busy","metadata":{"retry_after":2}}}`))
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.RetryAfter != 2*time.Second {
		t.Fatalf("err = %+v, want RetryAfter 2s", err)
	}
}
