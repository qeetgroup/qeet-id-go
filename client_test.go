package qeetid

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

// serviceNames lists the Client's service fields.
func serviceNames(c *Client) []string {
	var names []string
	v := reflect.ValueOf(c).Elem()
	for i := range v.NumField() {
		f := v.Type().Field(i)
		if f.IsExported() && f.Type.Kind() == reflect.Pointer {
			names = append(names, f.Name)
		}
	}
	return names
}

func TestNewClientWiresEveryService(t *testing.T) {
	c, err := NewClient(WithAPIKey("qk_test.secret"))
	if err != nil {
		t.Fatal(err)
	}
	v := reflect.ValueOf(c).Elem()
	for _, name := range serviceNames(c) {
		if v.FieldByName(name).IsNil() {
			t.Errorf("Client.%s is nil", name)
		}
	}
	if c.BaseURL() != DefaultBaseURL {
		t.Errorf("BaseURL() = %q, want %q", c.BaseURL(), DefaultBaseURL)
	}
}

func TestNewClientRejectsBadOptions(t *testing.T) {
	cases := []struct {
		name string
		opts []Option
		want string
	}{
		{"empty API key", []Option{WithAPIKey(" ")}, "the key is empty"},
		{"two credentials", []Option{WithAPIKey("k"), WithTokenSource(StaticToken("t"))}, "use one credential option"},
		{"plain http to a remote host", []Option{WithBaseURL("http://api.id.qeet.in")}, "must use https"},
		{"relative base URL", []Option{WithBaseURL("/v1")}, "not an absolute"},
		{"base URL with query", []Option{WithBaseURL("https://api.id.qeet.in/?x=1")}, "must not carry"},
		{"base URL with user info", []Option{WithBaseURL("https://user:pass@api.id.qeet.in")}, "must not carry"},
		{"reserved header", []Option{WithHeader("authorization", "x")}, "sets Authorization itself"},
		{"header with line break", []Option{WithHeader("X-Trace", "a\r\nb")}, "line break"},
		{"user agent with line break", []Option{WithUserAgentSuffix("app\n1")}, "line break"},
		{"bad retry delays", []Option{WithRetryDelay(time.Second, time.Millisecond)}, "WithRetryDelay"},
		{"client credentials without secret", []Option{WithClientCredentials("id", "")}, "required"},
		{"nil token source", []Option{WithTokenSource(nil)}, "nil"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewClient(tc.opts...)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("NewClient error = %v, want one containing %q", err, tc.want)
			}
		})
	}
}

func TestNewClientAcceptsLoopbackHTTP(t *testing.T) {
	for _, u := range []string{"http://localhost:4001", "http://127.0.0.1:8080/", "http://api.dev.id.qeet.localhost", "http://[::1]:9000"} {
		c, err := NewClient(WithBaseURL(u))
		if err != nil {
			t.Errorf("NewClient(WithBaseURL(%q)) = %v", u, err)
			continue
		}
		if strings.HasSuffix(c.BaseURL(), "/") {
			t.Errorf("BaseURL() = %q keeps the trailing slash", c.BaseURL())
		}
	}
}

func TestRequestHeaders(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{}`)}}
	c := newFakeClient(t, f,
		WithUserAgentSuffix("billing-worker/2.3.0"),
		WithHeader("X-Trace", "trace-1"),
	)
	_, err := c.Users.Get(context.Background(), "u 1/x", WithRequestHeader("X-Request-Tag", "tag-1"))
	if err != nil {
		t.Fatal(err)
	}
	got := f.calls()[0]
	if got.path != "/v1/users/u%201%2Fx" {
		t.Errorf("path = %q, want the user ID escaped", got.path)
	}
	ua := got.header.Get("User-Agent")
	if !strings.HasPrefix(ua, "qeet-id-go/"+Version+" (go") || !strings.HasSuffix(ua, " billing-worker/2.3.0") {
		t.Errorf("User-Agent = %q", ua)
	}
	if got.header.Get("Accept") != "application/json" {
		t.Errorf("Accept = %q", got.header.Get("Accept"))
	}
	if got.header.Get("X-Trace") != "trace-1" || got.header.Get("X-Request-Tag") != "tag-1" {
		t.Errorf("custom headers = %v", got.header)
	}
}

func TestEmptyPathArgumentIsRefused(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{}`)}}
	c := newFakeClient(t, f)
	_, err := c.Roles.Update(context.Background(), "tenant-1", "", UpdateRoleParams{})
	if err == nil || !strings.Contains(err.Error(), "roleID must not be empty") {
		t.Fatalf("err = %v, want roleID must not be empty", err)
	}
	if len(f.calls()) != 0 {
		t.Fatal("a request was sent")
	}
}

func TestNoCredentials(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{"valid":true}`)}}
	c, err := NewClient(WithHTTPClient(&http.Client{Transport: roundTripFunc(f.roundTrip)}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Users.Get(context.Background(), "u1"); !errors.Is(err, ErrNoCredentials) {
		t.Fatalf("authenticated call error = %v, want ErrNoCredentials", err)
	}
	// Public endpoints need none.
	if _, err := c.Credentials.Verify(context.Background(), VerifyCredentialParams{Credential: "eyJ"}); err != nil {
		t.Fatalf("public call: %v", err)
	}
	if n := len(f.calls()); n != 1 {
		t.Fatalf("sent %d requests, want 1", n)
	}
}

func TestTokenSourceCredential(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{}`)}}
	c2, err := NewClient(
		WithTokenSource(StaticToken("eyJ.agent.token")),
		WithBaseURL("https://api.example.test"),
		WithHTTPClient(&http.Client{Transport: roundTripFunc(f.roundTrip)}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c2.Tenants.Get(context.Background(), "t1"); err != nil {
		t.Fatal(err)
	}
	if got := f.calls()[0].header.Get("Authorization"); got != "Bearer eyJ.agent.token" {
		t.Errorf("Authorization = %q", got)
	}

	failing := TokenSourceFunc(func(context.Context) (string, error) { return "", errors.New("vault sealed") })
	c3, err := NewClient(WithTokenSource(failing), WithHTTPClient(&http.Client{Transport: roundTripFunc(f.roundTrip)}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c3.Tenants.Get(context.Background(), "t1")
	if err == nil || !strings.Contains(err.Error(), "vault sealed") {
		t.Fatalf("err = %v, want the token source's error", err)
	}
	if n := len(f.calls()); n != 1 {
		t.Fatalf("a request was sent without a token: %d calls", n)
	}
}

func TestLoggerNeverSeesCredentials(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{
		answer(http.StatusServiceUnavailable, `{"error":{"code":"internal","message":"x"}}`),
		answer(http.StatusOK, `{}`, "X-Request-Id", "req-42"),
	}}
	c := newFakeClient(t, f, WithLogger(logger))
	if _, err := c.Users.Get(context.Background(), "u1"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, secret := range []string{"qk_test.secret", "ApiKey", "Authorization"} {
		if strings.Contains(out, secret) {
			t.Errorf("log output contains %q:\n%s", secret, out)
		}
	}
	for _, want := range []string{"qeetid: request", "qeetid: retrying request", "request_id=req-42", "path=/v1/users/u1"} {
		if !strings.Contains(out, want) {
			t.Errorf("log output lacks %q:\n%s", want, out)
		}
	}
}
