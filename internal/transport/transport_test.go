package transport

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestTransport_RetriesOn429ThenSucceeds(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL, MaxRetries: 2})
	var out struct {
		OK bool `json:"ok"`
	}
	err := tr.Get(context.Background(), "/v1/check", RequestOptions{}, &out)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !out.OK {
		t.Fatal("expected ok=true")
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("calls = %d, want 2 (one retry)", got)
	}
}

func TestTransport_ErrorMapping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_123")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":"not_found","message":"nope"}}`))
	}))
	defer srv.Close()

	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL})
	err := tr.Get(context.Background(), "/v1/users/missing", RequestOptions{}, nil)
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("want *Error, got %T (%v)", err, err)
	}
	if e.Status != 404 || e.Code != "not_found" || e.Message != "nope" || e.RequestID != "req_123" {
		t.Fatalf("unexpected error: %+v", e)
	}
	if !e.IsNotFound() {
		t.Fatal("IsNotFound() should be true")
	}
}

func TestTransport_RateLimitCarriesRetryAfter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"code":"rate_limited","message":"slow down"}}`))
	}))
	defer srv.Close()

	// MaxRetries=0 so the 429 surfaces immediately instead of being retried.
	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL, MaxRetries: 0})
	err := tr.Get(context.Background(), "/v1/users", RequestOptions{}, nil)
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("want *Error, got %T (%v)", err, err)
	}
	if !e.IsRateLimited() || e.RetryAfterSeconds != 7 {
		t.Fatalf("unexpected: %+v", e)
	}
}

func TestTransport_NoRetryOnPost5xx(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"code":"server_error","message":"boom"}}`))
	}))
	defer srv.Close()

	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL, MaxRetries: 3})
	err := tr.Post(context.Background(), "/v1/users", map[string]string{"email": "a@b.com"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("calls = %d, want 1 (non-idempotent POST must not retry on 5xx)", got)
	}
}

func TestTransport_SendsAuthUserAgentAndCustomHeaders(t *testing.T) {
	var auth, ua, trace, accept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		ua = r.Header.Get("User-Agent")
		trace = r.Header.Get("X-Trace")
		accept = r.Header.Get("Accept")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	tr := New(Options{
		APIKey:  "qk_abc.def",
		BaseURL: srv.URL,
		// Authorization here must be ignored (reserved header wins).
		Headers: http.Header{"X-Trace": {"t1"}, "Authorization": {"nope"}},
	})
	if err := tr.Get(context.Background(), "/v1/users/u1", RequestOptions{}, nil); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if auth != "ApiKey qk_abc.def" {
		t.Fatalf("Authorization = %q, want ApiKey scheme (custom header must not override)", auth)
	}
	if !strings.HasPrefix(ua, "qeet-id-go/") {
		t.Fatalf("User-Agent = %q, want qeet-id-go/ prefix", ua)
	}
	if trace != "t1" {
		t.Fatalf("X-Trace = %q, want t1", trace)
	}
	if accept != "application/json" {
		t.Fatalf("Accept = %q", accept)
	}
}

func TestTransport_UserAgentOverridePrepended(t *testing.T) {
	var ua string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL, UserAgent: "myapp/1.2.0"})
	if err := tr.Get(context.Background(), "/v1/users/u1", RequestOptions{}, nil); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !strings.HasPrefix(ua, "myapp/1.2.0 qeet-id-go/") {
		t.Fatalf("User-Agent = %q, want myapp/1.2.0 prefix retained", ua)
	}
}

func TestTransport_RawBodySkipsJSONMarshal(t *testing.T) {
	var gotBody, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		_, _ = w.Write([]byte(`{"succeeded":1}`))
	}))
	defer srv.Close()

	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL})
	err := tr.Do(context.Background(), http.MethodPost, "/v1/users/bulk/import", RequestOptions{
		RawBody:        []byte("email,name\na@b.com,A\n"),
		RawContentType: "text/csv",
	}, nil)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotContentType != "text/csv" {
		t.Fatalf("Content-Type = %q, want text/csv", gotContentType)
	}
	if gotBody != "email,name\na@b.com,A\n" {
		t.Fatalf("body = %q", gotBody)
	}
}

func TestDoForm_BasicAuthAndErrorMapping(t *testing.T) {
	var gotUser, gotPass, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"code":"invalid_client","message":"bad creds"}}`))
	}))
	defer srv.Close()

	tr := New(Options{APIKey: "qk_x.y", BaseURL: srv.URL})
	err := tr.DoForm(context.Background(), http.MethodPost, "/v1/oauth/token", nil, "client1", "secret1", nil)
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("want *Error, got %T (%v)", err, err)
	}
	if !e.IsUnauthorized() || e.Code != "invalid_client" {
		t.Fatalf("unexpected: %+v", e)
	}
	if gotUser != "client1" || gotPass != "secret1" {
		t.Fatalf("basic auth = %q/%q", gotUser, gotPass)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q", gotContentType)
	}
}
