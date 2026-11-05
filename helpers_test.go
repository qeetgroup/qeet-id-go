package qeetid

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// roundTripFunc fakes the network: tests answer requests in process.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// response builds an HTTP response with a JSON body and optional header
// key/value pairs.
func response(status int, body string, header ...string) *http.Response {
	h := http.Header{"Content-Type": {"application/json"}}
	for i := 0; i+1 < len(header); i += 2 {
		h.Set(header[i], header[i+1])
	}
	return &http.Response{
		StatusCode: status,
		Header:     h,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// recorded is a request as the fake server saw it.
type recorded struct {
	method string
	path   string
	query  string
	header http.Header
	body   string
}

// fakeAPI answers requests from a script, one answer per call (the last answer
// repeats), and records every request.
type fakeAPI struct {
	mu       sync.Mutex
	answers  []func(*http.Request) *http.Response
	requests []recorded
}

func (f *fakeAPI) roundTrip(r *http.Request) (*http.Response, error) {
	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
	}
	f.mu.Lock()
	n := len(f.requests)
	f.requests = append(f.requests, recorded{
		method: r.Method,
		path:   r.URL.EscapedPath(),
		query:  r.URL.RawQuery,
		header: r.Header.Clone(),
		body:   string(body),
	})
	answer := f.answers[min(n, len(f.answers)-1)]
	f.mu.Unlock()
	return answer(r), nil
}

func (f *fakeAPI) calls() []recorded {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]recorded(nil), f.requests...)
}

func answer(status int, body string, header ...string) func(*http.Request) *http.Response {
	return func(*http.Request) *http.Response { return response(status, body, header...) }
}

// newFakeClient returns a client whose requests the fake answers. Retry
// delays are short; tests that wait on Retry-After run inside synctest.
func newFakeClient(t *testing.T, f *fakeAPI, opts ...Option) *Client {
	t.Helper()
	base := []Option{
		WithAPIKey("qk_test.secret"),
		WithBaseURL("https://api.example.test"),
		WithHTTPClient(&http.Client{Transport: roundTripFunc(f.roundTrip)}),
		WithRetryDelay(time.Millisecond, 4*time.Millisecond),
	}
	c, err := NewClient(append(base, opts...)...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}
