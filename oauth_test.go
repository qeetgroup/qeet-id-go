package qeetid

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// clientCredentialsAnswer is the server's answer to a client credentials grant
// (developer/principal TokenResponse).
const clientCredentialsAnswer = `{"access_token":"eyJ.service.token","token_type":"Bearer","expires_in":600,
	"expires_at":"2026-11-05T10:10:00Z","scope":"user.read user.write"}`

func TestClientCredentials(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, clientCredentialsAnswer)}}
	c := newFakeClient(t, f)
	tok, err := c.OAuth.ClientCredentials(context.Background(), ClientCredentialsParams{
		ClientID: "sp-1", ClientSecret: "s3cr3t+/=", Scope: "user.read",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tok.ExpiresIn != 600 || tok.AccessToken != "eyJ.service.token" || tok.Scope != "user.read user.write" {
		t.Errorf("token = %+v", tok)
	}
	if want := time.Date(2026, 11, 5, 10, 10, 0, 0, time.UTC); !tok.Expiry().Equal(want) {
		t.Errorf("Expiry() = %v, want %v", tok.Expiry(), want)
	}
	got := f.calls()[0]
	if got.method != http.MethodPost || got.path != "/v1/oauth/token" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	req := &http.Request{Header: got.header}
	user, pass, ok := req.BasicAuth()
	if !ok || user != "sp-1" || pass != "s3cr3t+/=" {
		t.Errorf("Basic credentials = %q %q %v; want them sent as is", user, pass, ok)
	}
	form, _ := url.ParseQuery(got.body)
	if form.Get("grant_type") != "client_credentials" || form.Get("scope") != "user.read" || form.Has("client_secret") {
		t.Errorf("form = %v", form)
	}
	if got.header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q", got.header.Get("Content-Type"))
	}
}

func TestExpiresInWithoutExpiresAt(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{"access_token":"at","token_type":"Bearer","expires_in":300,"refresh_token":"rt"}`)}}
	c := newFakeClient(t, f)
	before := time.Now()
	tok, err := c.OAuth.RefreshToken(context.Background(), RefreshTokenParams{ClientID: "app", ClientSecret: "s", RefreshToken: "rt0"})
	if err != nil {
		t.Fatal(err)
	}
	if tok.ExpiresIn != 300 || !tok.ExpiresAt.IsZero() {
		t.Fatalf("token = %+v", tok)
	}
	if exp := tok.Expiry(); exp.Before(before.Add(300*time.Second)) || exp.After(time.Now().Add(300*time.Second)) {
		t.Errorf("Expiry() = %v, want about 300s from now", exp)
	}
}

// TestTokenEndpointErrors covers each error /v1/oauth/token and
// /v1/oauth/token-code answer with, in the RFC 6749 shape.
func TestTokenEndpointErrors(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		code   OAuthErrorCode
		is     error
		call   func(context.Context, *Client) error
	}{
		{"missing grant type", 400, `{"error":"invalid_request","error_description":"grant_type is required."}`, OAuthInvalidRequest, ErrInvalidRequest, clientCreds},
		{"wrong secret", 401, `{"error":"invalid_client","error_description":"Invalid client credentials."}`, OAuthInvalidClient, ErrUnauthorized, clientCreds},
		{"disabled principal", 400, `{"error":"unauthorized_client","error_description":"This service principal is disabled."}`, OAuthUnauthorizedClient, ErrInvalidRequest, clientCreds},
		{"other grant", 400, `{"error":"unsupported_grant_type","error_description":"Use client_credentials."}`, OAuthUnsupportedGrantType, ErrInvalidRequest, clientCreds},
		{"reused code", 400, `{"error":"invalid_grant","error_description":"The authorization code was already used."}`, OAuthInvalidGrant, ErrInvalidRequest, exchangeCode},
		{"scope beyond subject", 400, `{"error":"invalid_scope","error_description":"scope exceeds the subject token"}`, OAuthInvalidScope, ErrInvalidRequest, exchangeToken},
		{"bad resource", 400, `{"error":"invalid_target","error_description":"resource must be an absolute URI without a fragment"}`, OAuthInvalidTarget, ErrInvalidRequest, exchangeCode},
		{"pending", 400, `{"error":"authorization_pending"}`, OAuthAuthorizationPending, ErrInvalidRequest, pollDevice},
		{"slow down", 400, `{"error":"slow_down"}`, OAuthSlowDown, ErrInvalidRequest, pollDevice},
		{"denied", 400, `{"error":"access_denied","error_description":"address_not_allowed"}`, OAuthAccessDenied, ErrInvalidRequest, pollDevice},
		{"expired", 400, `{"error":"expired_token"}`, OAuthExpiredToken, ErrInvalidRequest, pollBackchannel},
		{"access token revoke", 400, `{"error":"unsupported_token_type"}`, OAuthUnsupportedTokenType, ErrInvalidRequest, revoke},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(tc.status, tc.body, "X-Request-Id", "req-oauth")}}
			c := newFakeClient(t, f)
			err := tc.call(context.Background(), c)
			oe, ok := errors.AsType[*OAuthError](err)
			if !ok {
				t.Fatalf("err = %T %v, want *OAuthError", err, err)
			}
			if oe.Code != tc.code || oe.StatusCode != tc.status || oe.RequestID != "req-oauth" {
				t.Errorf("OAuthError = %+v", oe)
			}
			if !errors.Is(err, tc.is) {
				t.Errorf("errors.Is(err, %v) = false", tc.is)
			}
			if len(f.calls()) != 1 {
				t.Errorf("an OAuth error was retried: %d requests", len(f.calls()))
			}
		})
	}
}

func clientCreds(ctx context.Context, c *Client) error {
	_, err := c.OAuth.ClientCredentials(ctx, ClientCredentialsParams{ClientID: "sp", ClientSecret: "s"})
	return err
}

func exchangeCode(ctx context.Context, c *Client) error {
	_, err := c.OAuth.ExchangeCode(ctx, ExchangeCodeParams{ClientID: "app", ClientSecret: "s", Code: "c", RedirectURI: "https://app.example/cb"})
	return err
}

func exchangeToken(ctx context.Context, c *Client) error {
	_, err := c.OAuth.ExchangeToken(ctx, TokenExchangeParams{ClientID: "app", ClientSecret: "s", SubjectToken: "eyJ", Scope: "admin"})
	return err
}

func pollDevice(ctx context.Context, c *Client) error {
	_, err := c.OAuth.PollDeviceToken(ctx, DeviceTokenParams{ClientID: "tv-app", DeviceCode: "dc"})
	return err
}

func pollBackchannel(ctx context.Context, c *Client) error {
	_, err := c.OAuth.PollBackchannelToken(ctx, BackchannelTokenParams{ClientID: "kiosk", ClientSecret: "s", AuthReqID: "ar"})
	return err
}

func revoke(ctx context.Context, c *Client) error {
	return c.OAuth.Revoke(ctx, RevokeParams{ClientID: "app", ClientSecret: "s", Token: "eyJ.access"})
}

func TestTokenEndpointRateLimitIsAnAPIError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeAPI{answers: []func(*http.Request) *http.Response{
			answer(http.StatusTooManyRequests, `{"error":{"code":"too_many_requests","message":"slow down","metadata":{"retry_after":1}}}`, "Retry-After", "1"),
			answer(http.StatusOK, clientCredentialsAnswer),
		}}
		c := newFakeClient(t, f)
		if err := clientCreds(t.Context(), c); err != nil {
			t.Fatalf("err = %v; a 429 with Retry-After should be retried", err)
		}
		if len(f.calls()) != 2 {
			t.Fatalf("sent %d requests, want 2", len(f.calls()))
		}
	})
}

func TestTokenExchangeForm(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK,
		`{"access_token":"x","token_type":"Bearer","expires_in":120,"issued_token_type":"urn:ietf:params:oauth:token-type:access_token"}`)}}
	c := newFakeClient(t, f)
	tok, err := c.OAuth.ExchangeToken(context.Background(), TokenExchangeParams{
		ClientID: "app", ClientSecret: "s", SubjectToken: "subj", ActorToken: "agent", Scope: "docs.read", Resource: "https://mcp.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tok.IssuedTokenType != TokenTypeAccessToken {
		t.Errorf("IssuedTokenType = %q", tok.IssuedTokenType)
	}
	got := f.calls()[0]
	form, _ := url.ParseQuery(got.body)
	want := map[string]string{
		"grant_type":         "urn:ietf:params:oauth:grant-type:token-exchange",
		"subject_token":      "subj",
		"subject_token_type": TokenTypeAccessToken,
		"actor_token":        "agent",
		"actor_token_type":   TokenTypeAccessToken,
		"scope":              "docs.read",
		"resource":           "https://mcp.example.com",
	}
	for k, v := range want {
		if form.Get(k) != v {
			t.Errorf("form[%s] = %q, want %q", k, form.Get(k), v)
		}
	}
	if got.path != "/v1/oauth/token-code" {
		t.Errorf("path = %q: token exchange is served by the token-code endpoint", got.path)
	}
}

func TestPublicClientSendsClientIDInForm(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK,
		`{"device_code":"dc","user_code":"WDJB-MJHT","verification_uri":"https://id.example/device","expires_in":600,"interval":5}`)}}
	c := newFakeClient(t, f)
	res, err := c.OAuth.AuthorizeDevice(context.Background(), DeviceAuthorizationParams{ClientID: "tv-app", Scope: "openid"})
	if err != nil {
		t.Fatal(err)
	}
	if res.UserCode != "WDJB-MJHT" || res.Interval != 5 {
		t.Errorf("result = %+v", res)
	}
	got := f.calls()[0]
	form, _ := url.ParseQuery(got.body)
	if form.Get("client_id") != "tv-app" || got.header.Get("Authorization") != "" {
		t.Errorf("public client: form %v, Authorization %q", form, got.header.Get("Authorization"))
	}
}

func TestIntrospect(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{
		answer(http.StatusOK, `{"active":true,"token_type":"Bearer","sub":"u1","scope":"openid email","iss":"qeet-id",
			"aud":["qeet-id","https://mcp.example.com"],"client_id":"app","resource":"https://mcp.example.com",
			"exp":1762337400,"iat":1762336800,"tenant_id":"t1","sid":"s1","act":{"sub":"agent-7"}}`),
		answer(http.StatusOK, `{"active":false}`),
		answer(http.StatusUnauthorized, `{"error":{"code":"oidc.client_auth_failed","message":"The client could not be authenticated."}}`),
	}}
	c := newFakeClient(t, f)
	ctx := context.Background()
	params := IntrospectParams{ClientID: "app", ClientSecret: "s", Token: "eyJ"}

	in, err := c.OAuth.Introspect(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	if !in.Active || in.Subject != "u1" || len(in.Audience) != 2 || len(in.Resource) != 1 || in.Actor == nil || in.Actor.Subject != "agent-7" ||
		in.ExpiresAt.Unix() != 1762337400 || in.SessionID != "s1" {
		t.Errorf("introspection = %+v", in)
	}
	if f.calls()[0].path != "/oauth/introspect" {
		t.Errorf("path = %q", f.calls()[0].path)
	}

	in, err = c.OAuth.Introspect(ctx, params)
	if err != nil || in.Active {
		t.Fatalf("inactive token: %+v, %v", in, err)
	}

	_, err = c.OAuth.Introspect(ctx, params)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
	if _, err := c.OAuth.Introspect(ctx, IntrospectParams{ClientID: "spa", Token: "eyJ"}); err == nil {
		t.Error("a public client was allowed to introspect")
	}
}

func TestUserInfoSendsThePersonsToken(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{"sub":"u1","email":"ada@example.com","email_verified":true}`)}}
	c := newFakeClient(t, f)
	info, err := c.OAuth.UserInfo(context.Background(), "eyJ.person")
	if err != nil {
		t.Fatal(err)
	}
	if info.Email != "ada@example.com" || !info.EmailVerified {
		t.Errorf("info = %+v", info)
	}
	if got := f.calls()[0].header.Get("Authorization"); got != "Bearer eyJ.person" {
		t.Errorf("Authorization = %q, want the person's token, not the API key", got)
	}
}

func TestDiscover(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK,
		`{"issuer":"qeet-id","jwks_uri":"https://api.id.qeet.in/.well-known/jwks.json","token_endpoint":"https://api.id.qeet.in/v1/oauth/token-code",
		"grant_types_supported":["authorization_code","urn:openid:params:grant-type:ciba"],"actor_types_supported":["user","service","agent"]}`)}}
	c := newFakeClient(t, f)
	md, err := c.OAuth.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if md.Issuer != "qeet-id" || !strings.HasSuffix(md.JWKSURI, "/.well-known/jwks.json") || len(md.ActorTypesSupported) != 3 {
		t.Errorf("metadata = %+v", md)
	}
	if f.calls()[0].header.Get("Authorization") != "" {
		t.Error("discovery sent credentials")
	}
}

// TestWithClientCredentials checks that the client fetches a service principal
// token once, shares it across concurrent calls, sends it as a bearer token and
// renews it before it expires.
func TestWithClientCredentials(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var mu sync.Mutex
		tokens := 0
		rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path == "/v1/oauth/token" {
				mu.Lock()
				tokens++
				n := tokens
				mu.Unlock()
				exp := time.Now().Add(10 * time.Minute).UTC().Format(time.RFC3339)
				return response(http.StatusOK, `{"access_token":"tok-`+string(rune('0'+n))+`","token_type":"Bearer","expires_in":600,"expires_at":"`+exp+`"}`), nil
			}
			if r.Header.Get("Authorization") == "" || !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer tok-") {
				return response(http.StatusUnauthorized, `{"error":{"code":"unauthorized","message":"no"}}`), nil
			}
			return response(http.StatusOK, `{"id":"`+r.Header.Get("Authorization")+`"}`), nil
		})
		c, err := NewClient(WithClientCredentials("sp-1", "secret", "user.read"), WithHTTPClient(&http.Client{Transport: rt}))
		if err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				u, err := c.Users.Get(t.Context(), "u1")
				if err != nil || u.ID != "Bearer tok-1" {
					t.Errorf("Get = %+v, %v", u, err)
				}
			})
		}
		wg.Wait()
		if tokens != 1 {
			t.Fatalf("fetched %d tokens for concurrent calls, want 1", tokens)
		}
		time.Sleep(9*time.Minute + 30*time.Second) // inside the refresh margin
		u, err := c.Users.Get(t.Context(), "u1")
		if err != nil || u.ID != "Bearer tok-2" || tokens != 2 {
			t.Fatalf("after expiry: %+v, %v, %d token requests", u, err, tokens)
		}
	})
}

func TestWithClientCredentialsFailure(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusUnauthorized, `{"error":"invalid_client","error_description":"Invalid client credentials."}`)}}
	c, err := NewClient(WithClientCredentials("sp-1", "wrong"), WithHTTPClient(&http.Client{Transport: roundTripFunc(f.roundTrip)}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Users.Get(context.Background(), "u1")
	oe, ok := errors.AsType[*OAuthError](err)
	if !ok || oe.Code != OAuthInvalidClient || !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("err = %v, want the token endpoint's invalid_client", err)
	}
	if strings.Contains(err.Error(), "wrong") {
		t.Errorf("error text leaks the secret: %v", err)
	}
	if n := len(f.calls()); n != 1 {
		t.Errorf("sent %d requests, want only the token request", n)
	}
}
