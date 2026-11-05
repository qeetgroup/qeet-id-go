package qeetid

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"testing/synctest"
	"time"
)

const (
	overloaded = `{"error":{"code":"server.overloaded","message":"The service is busy.","retryable":true}}`
	rateLimit  = `{"error":{"code":"too_many_requests","message":"Too many attempts.","retryable":true}}`
	internal   = `{"error":{"code":"internal","message":"Something went wrong.","retryable":true}}`
)

func TestRetryPolicy(t *testing.T) {
	cases := []struct {
		name    string
		answers []func(*http.Request) *http.Response
		call    func(context.Context, *Client) error
		calls   int
		wantErr error
		apiErr  bool // the error must be an *APIError
	}{
		{
			name:    "GET retried after a 503 without Retry-After",
			answers: []func(*http.Request) *http.Response{answer(503, internal), answer(200, `{}`)},
			call:    getUser,
			calls:   2,
		},
		{
			name:    "GET retried after a 500, up to the limit",
			answers: []func(*http.Request) *http.Response{answer(500, internal)},
			call:    getUser,
			calls:   3,
			apiErr:  true,
		},
		{
			name:    "DELETE retried after a 502",
			answers: []func(*http.Request) *http.Response{answer(502, ``), answer(204, ``)},
			call:    func(ctx context.Context, c *Client) error { return c.Groups.Delete(ctx, "g1") },
			calls:   2,
		},
		{
			name:    "POST not retried after a 500",
			answers: []func(*http.Request) *http.Response{answer(500, internal)},
			call:    createRole,
			calls:   1,
			apiErr:  true,
		},
		{
			name:    "POST retried after a 429 with Retry-After",
			answers: []func(*http.Request) *http.Response{answer(429, rateLimit, "Retry-After", "2"), answer(201, `{"id":"r1"}`)},
			call:    createRole,
			calls:   2,
		},
		{
			name:    "POST retried after a 503 overload with Retry-After",
			answers: []func(*http.Request) *http.Response{answer(503, overloaded, "Retry-After", "3"), answer(201, `{"id":"r1"}`)},
			call:    createRole,
			calls:   2,
		},
		{
			name:    "429 without Retry-After is a lockout, not retried",
			answers: []func(*http.Request) *http.Response{answer(429, `{"error":{"code":"auth.account_locked","message":"locked"}}`)},
			call:    getUser,
			calls:   1,
			wantErr: ErrRateLimited,
		},
		{
			name:    "Retry-After beyond MaxRetryAfter is returned, not waited for",
			answers: []func(*http.Request) *http.Response{answer(429, rateLimit, "Retry-After", "3600")},
			call:    getUser,
			calls:   1,
			wantErr: ErrRateLimited,
		},
		{
			name:    "401 never retried",
			answers: []func(*http.Request) *http.Response{answer(401, `{"error":{"code":"unauthorized","message":"Please sign in."}}`)},
			call:    getUser,
			calls:   1,
			wantErr: ErrUnauthorized,
		},
		{
			name: "session revoked never retried, even with Retry-After and a retryable flag",
			answers: []func(*http.Request) *http.Response{answer(401,
				`{"error":{"code":"auth.session_revoked","message":"Your session was revoked.","retryable":true}}`, "Retry-After", "1")},
			call:    getUser,
			calls:   1,
			wantErr: ErrSessionRevoked,
		},
		{
			name:    "404 not retried",
			answers: []func(*http.Request) *http.Response{answer(404, `{"error":{"code":"not_found","message":"x"}}`)},
			call:    getUser,
			calls:   1,
			wantErr: ErrNotFound,
		},
		{
			name: "create with an idempotency key retried after a 500 with the same key",
			answers: []func(*http.Request) *http.Response{
				answer(500, internal),
				answer(201, `{"id":"u1"}`),
			},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Users.Create(ctx, CreateUserParams{TenantID: "t1", Email: "a@example.com"}, WithIdempotencyKey("import-row-17"))
				return err
			},
			calls: 2,
		},
		{
			name:    "create without a key gets a generated one and is retried",
			answers: []func(*http.Request) *http.Response{answer(503, internal), answer(201, `{"id":"g1"}`)},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Groups.Create(ctx, CreateGroupParams{Name: "eng"})
				return err
			},
			calls: 2,
		},
		{
			name: "idempotency key still in progress is retried",
			answers: []func(*http.Request) *http.Response{
				answer(409, `{"error":{"code":"idempotency.in_progress","message":"retry shortly"}}`),
				answer(201, `{"id":"w1"}`),
			},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Webhooks.Create(ctx, CreateWebhookParams{URL: "https://example.com/hook", Events: []string{"user.created"}})
				return err
			},
			calls: 2,
		},
		{
			name:    "idempotency outcome unknown is not retried",
			answers: []func(*http.Request) *http.Response{answer(409, `{"error":{"code":"idempotency.outcome_unknown","message":"check"}}`)},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.APIKeys.Create(ctx, CreateAPIKeyParams{Name: "ci"})
				return err
			},
			calls:   1,
			wantErr: ErrConflict,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				f := &fakeAPI{answers: tc.answers}
				c := newFakeClient(t, f)
				err := tc.call(t.Context(), c)
				wantFailure := tc.wantErr != nil || tc.apiErr
				var apiErr *APIError
				switch {
				case !wantFailure && err != nil:
					t.Fatalf("err = %v", err)
				case wantFailure && err == nil:
					t.Fatal("err = nil, want a failure")
				case tc.apiErr && !errors.As(err, &apiErr):
					t.Fatalf("err = %T, want *APIError", err)
				case tc.wantErr != nil && !errors.Is(err, tc.wantErr):
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				calls := f.calls()
				if len(calls) != tc.calls {
					t.Fatalf("sent %d requests, want %d", len(calls), tc.calls)
				}
				first := calls[0].header.Get("Idempotency-Key")
				for _, r := range calls[1:] {
					if r.header.Get("Idempotency-Key") != first {
						t.Errorf("retry changed Idempotency-Key from %q to %q", first, r.header.Get("Idempotency-Key"))
					}
					if r.body != calls[0].body {
						t.Errorf("retry changed the body")
					}
				}
			})
		})
	}
}

func getUser(ctx context.Context, c *Client) error {
	_, err := c.Users.Get(ctx, "u1")
	return err
}

func createRole(ctx context.Context, c *Client) error {
	_, err := c.Roles.Create(ctx, "t1", CreateRoleParams{Name: "editor"})
	return err
}

func TestRetryHonoursRetryAfter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(429, rateLimit, "Retry-After", "7"), answer(200, `{}`)}}
		c := newFakeClient(t, f)
		start := time.Now()
		if err := getUser(t.Context(), c); err != nil {
			t.Fatal(err)
		}
		if waited := time.Since(start); waited < 7*time.Second || waited > 8*time.Second {
			t.Fatalf("waited %v, want the 7s the server asked for", waited)
		}
	})
}

func TestRetryStopsAtContextDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(429, rateLimit, "Retry-After", "30")}}
		c := newFakeClient(t, f)
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()
		start := time.Now()
		err := getUser(ctx, c)
		if !errors.Is(err, ErrRateLimited) {
			t.Fatalf("err = %v, want the 429", err)
		}
		if time.Since(start) != 0 {
			t.Fatalf("slept %v although the deadline comes first", time.Since(start))
		}
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.RetryAfter != 30*time.Second {
			t.Fatalf("RetryAfter = %v", apiErr.RetryAfter)
		}
	})
}

func TestNetworkErrors(t *testing.T) {
	refused := func(*http.Request) (*http.Response, error) {
		return nil, &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
	}
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		rt := roundTripFunc(func(r *http.Request) (*http.Response, error) { calls++; return refused(r) })
		c, err := NewClient(WithAPIKey("k"), WithHTTPClient(&http.Client{Transport: rt}), WithRetryDelay(time.Millisecond, time.Millisecond))
		if err != nil {
			t.Fatal(err)
		}
		_, err = c.Users.Get(t.Context(), "u1")
		var opErr *net.OpError
		if !errors.As(err, &opErr) {
			t.Fatalf("err = %v, want a *net.OpError in the chain", err)
		}
		if calls != 3 {
			t.Fatalf("GET sent %d times, want 3", calls)
		}
		calls = 0
		_, err = c.Roles.Create(t.Context(), "t1", CreateRoleParams{Name: "x"})
		if err == nil || calls != 1 {
			t.Fatalf("POST sent %d times (err %v), want 1", calls, err)
		}
	})
}

func TestTimeoutPerAttempt(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		slow := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			<-r.Context().Done()
			return nil, r.Context().Err()
		})
		c, err := NewClient(WithAPIKey("k"), WithHTTPClient(&http.Client{Transport: slow}), WithTimeout(2*time.Second), WithMaxRetries(0))
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now()
		_, err = c.Users.Get(t.Context(), "u1")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded", err)
		}
		if time.Since(start) != 2*time.Second {
			t.Fatalf("gave up after %v, want 2s", time.Since(start))
		}
	})
}

func TestIdempotencyKeyOption(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(201, `{}`)}}
	c := newFakeClient(t, f)
	ctx := context.Background()

	if _, err := c.Agents.Create(ctx, "t1", CreateAgentParams{Name: "triage-bot"}, WithIdempotencyKey("agent-triage-v1")); err != nil {
		t.Fatal(err)
	}
	if got := f.calls()[0].header.Get("Idempotency-Key"); got != "agent-triage-v1" {
		t.Errorf("Idempotency-Key = %q", got)
	}

	// A call whose endpoint has no idempotency support refuses the option.
	if _, err := c.Roles.Create(ctx, "t1", CreateRoleParams{Name: "x"}, WithIdempotencyKey("k1")); err == nil {
		t.Error("Roles.Create accepted an idempotency key")
	}
	for _, bad := range []string{"", "has space", string(make([]byte, 256))} {
		if _, err := c.Users.Create(ctx, CreateUserParams{}, WithIdempotencyKey(bad)); err == nil {
			t.Errorf("key %q accepted", bad)
		}
	}
	if n := len(f.calls()); n != 1 {
		t.Fatalf("refused calls were sent: %d requests", n)
	}

	// With retries off, no key is generated: there is nothing to make safe.
	f2 := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(201, `{}`)}}
	c2 := newFakeClient(t, f2, WithMaxRetries(0))
	if _, err := c2.Users.Create(ctx, CreateUserParams{}); err != nil {
		t.Fatal(err)
	}
	if got := f2.calls()[0].header.Get("Idempotency-Key"); got != "" {
		t.Errorf("Idempotency-Key = %q with retries off", got)
	}
}

func TestGeneratedIdempotencyKeysAreUUIDs(t *testing.T) {
	seen := map[string]bool{}
	for range 100 {
		k := newIdempotencyKey()
		if len(k) != 36 || k[14] != '4' || !validIdempotencyKey(k) || seen[k] {
			t.Fatalf("bad key %q", k)
		}
		seen[k] = true
	}
}
