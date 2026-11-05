package qeetid

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// handRouteCases covers the methods written by hand rather than from the
// contract tables.
var handRouteCases = []routeCase{
	{name: "OAuth.ClientCredentials", method: "POST", path: "/v1/oauth/token", public: true, basic: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuth.ClientCredentials(ctx, ClientCredentialsParams{ClientID: "id", ClientSecret: "s"})
		return err
	}},
	{name: "OAuth.ExchangeCode", method: "POST", path: "/v1/oauth/token-code", public: true, basic: true, call: exchangeCode},
	{name: "OAuth.RefreshToken", method: "POST", path: "/v1/oauth/token-code", public: true, basic: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuth.RefreshToken(ctx, RefreshTokenParams{ClientID: "id", ClientSecret: "s", RefreshToken: "rt"})
		return err
	}},
	{name: "OAuth.ExchangeToken", method: "POST", path: "/v1/oauth/token-code", public: true, basic: true, call: exchangeToken},
	{name: "OAuth.AuthorizeDevice", method: "POST", path: "/v1/oauth/device_authorization", public: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuth.AuthorizeDevice(ctx, DeviceAuthorizationParams{ClientID: "tv"})
		return err
	}},
	{name: "OAuth.PollDeviceToken", method: "POST", path: "/v1/oauth/token-code", public: true, call: pollDevice},
	{name: "OAuth.AuthorizeBackchannel", method: "POST", path: "/v1/oauth/bc-authorize", public: true, basic: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuth.AuthorizeBackchannel(ctx, BackchannelAuthorizationParams{ClientID: "k", ClientSecret: "s", LoginHint: "ada@example.com"})
		return err
	}},
	{name: "OAuth.PollBackchannelToken", method: "POST", path: "/v1/oauth/token-code", public: true, basic: true, call: pollBackchannel},
	{name: "OAuth.Introspect", method: "POST", path: "/oauth/introspect", public: true, basic: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuth.Introspect(ctx, IntrospectParams{ClientID: "id", ClientSecret: "s", Token: "t"})
		return err
	}},
	{name: "OAuth.Revoke", method: "POST", path: "/oauth/revoke", public: true, basic: true, call: revoke},
	{name: "OAuth.Discover", method: "GET", path: "/.well-known/openid-configuration", public: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuth.Discover(ctx)
		return err
	}},
	{name: "LDAP.Authenticate", method: "POST", path: "/ldap/conn-1/authenticate", public: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Authenticate(ctx, "conn-1", LDAPAuthenticateParams{Username: "u", Password: "p"})
		return err
	}},
	{name: "Users.Import", method: "POST", path: "/v1/users/bulk/import", call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.Import(ctx, ImportSourceCognito, []byte("email\na@x.in\n"))
		return err
	}},
}

func TestRoutes(t *testing.T) {
	for _, tc := range append(append([]routeCase(nil), routeCases...), handRouteCases...) {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{"items":[]}`)}}
			c := newFakeClient(t, f)
			if err := tc.call(context.Background(), c); err != nil {
				t.Fatalf("call: %v", err)
			}
			calls := f.calls()
			if len(calls) != 1 {
				t.Fatalf("sent %d requests, want 1", len(calls))
			}
			got := calls[0]
			if got.method != tc.method || got.path != tc.path {
				t.Errorf("sent %s %s, want %s %s", got.method, got.path, tc.method, tc.path)
			}
			auth := got.header.Get("Authorization")
			switch {
			case tc.basic && !strings.HasPrefix(auth, "Basic "):
				t.Errorf("Authorization = %q, want the OAuth client's Basic credentials", auth)
			case tc.basic:
			case tc.public && auth != "":
				t.Errorf("public endpoint got Authorization %q, want none", auth)
			case !tc.public && auth != "ApiKey qk_test.secret":
				t.Errorf("Authorization = %q, want the API key", auth)
			}
			if key := got.header.Get("Idempotency-Key"); (key != "") != tc.idempotent {
				t.Errorf("Idempotency-Key = %q, want one: %v", key, tc.idempotent)
			}
		})
	}
}

// TestRoutesCoverEveryService fails when a service has no route case, which
// happens when a method is added by hand without a test.
func TestRoutesCoverEveryService(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	for _, tc := range append(append([]routeCase(nil), routeCases...), handRouteCases...) {
		for i := range len(tc.name) {
			if tc.name[i] == '.' {
				covered[tc.name[:i]] = true
				break
			}
		}
	}
	for _, name := range serviceNames(c) {
		if !covered[name] {
			t.Errorf("service %s has no route case", name)
		}
	}
}
