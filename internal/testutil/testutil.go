// Package testutil holds the shared httptest server + request-recorder used
// by every resource's client_test.go, so each test file doesn't reimplement
// the same "spin up a server, record what hit it, return canned JSON"
// boilerplate.
package testutil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Recorder captures the last request an httptest server received.
type Recorder struct {
	Method string
	Path   string
	Body   string
	Query  string
}

// NewServer starts an httptest.Server that records each request into rec and
// replies with the given response body (as-is, so pass valid JSON). It is
// registered for cleanup on t.Cleanup.
func NewServer(t *testing.T, rec *Recorder, response string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		rec.Method, rec.Path, rec.Body, rec.Query = r.Method, r.URL.Path, string(b), r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// NewTransport builds an internal/transport.Transport pointed at srv, for
// tests that want to call resource methods directly without going through
// the root package's Client/Config.
func NewTransport(srv *httptest.Server) *transport.Transport {
	return transport.New(transport.Options{APIKey: "qk_test.secret", BaseURL: srv.URL})
}

// NewRecordingTransport is the common one-liner: server + recorder +
// transport, wired together.
func NewRecordingTransport(t *testing.T, rec *Recorder, response string) *transport.Transport {
	t.Helper()
	return NewTransport(NewServer(t, rec, response))
}
