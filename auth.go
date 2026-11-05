package qeetid

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenSource supplies the bearer token sent with each request. Implementations
// must be safe for concurrent use and should cache tokens rather than fetch one
// per call.
//
// A golang.org/x/oauth2 token source can be adapted in a few lines:
//
//	qeetid.WithTokenSource(qeetid.TokenSourceFunc(func(context.Context) (string, error) {
//		t, err := ts.Token() // an oauth2.TokenSource
//		if err != nil {
//			return "", err
//		}
//		return t.AccessToken, nil
//	}))
type TokenSource interface {
	// Token returns a valid access token.
	Token(ctx context.Context) (string, error)
}

// TokenSourceFunc adapts a function to the TokenSource interface.
type TokenSourceFunc func(ctx context.Context) (string, error)

// Token calls f(ctx).
func (f TokenSourceFunc) Token(ctx context.Context) (string, error) { return f(ctx) }

// StaticToken returns a TokenSource that always returns token. Use it for a
// token whose lifetime you manage yourself.
func StaticToken(token string) TokenSource {
	return TokenSourceFunc(func(context.Context) (string, error) { return token, nil })
}

// credential adds authentication to an outgoing request.
type credential interface {
	authorize(ctx context.Context, req *http.Request) error
}

type apiKeyCredential struct{ key string }

func (k apiKeyCredential) authorize(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "ApiKey "+k.key)
	return nil
}

type bearerCredential struct{ ts TokenSource }

func (b bearerCredential) authorize(ctx context.Context, req *http.Request) error {
	tok, err := b.ts.Token(ctx)
	if err != nil {
		return err
	}
	if tok == "" {
		return errors.New("token source returned an empty token")
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	return nil
}

// tokenRefreshMargin is how long before expiry a cached client credentials
// token is replaced, so that a request never leaves with a token about to lapse.
const tokenRefreshMargin = time.Minute

// clientCredentialsSource obtains and caches service principal tokens from the
// client credentials endpoint (POST /v1/oauth/token).
type clientCredentialsSource struct {
	c            *Client
	clientID     string
	clientSecret string
	scope        string

	mu     sync.Mutex
	token  string
	expiry time.Time
}

func newClientCredentialsSource(c *Client, cfg *clientCredentialsConfig) *clientCredentialsSource {
	return &clientCredentialsSource{
		c:            c,
		clientID:     cfg.clientID,
		clientSecret: cfg.clientSecret,
		scope:        strings.Join(cfg.scopes, " "),
	}
}

// Token returns the cached token, fetching a new one when there is none or it
// is about to expire. Concurrent callers share one fetch.
func (s *clientCredentialsSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.token != "" && time.Now().Add(tokenRefreshMargin).Before(s.expiry) {
		return s.token, nil
	}
	tok, err := s.c.clientCredentialsToken(ctx, s.clientID, s.clientSecret, s.scope)
	if err != nil {
		return "", err
	}
	s.token, s.expiry = tok.AccessToken, tok.Expiry()
	return s.token, nil
}

// clientCredentialsToken runs the client credentials grant. It is the one
// token request the core needs; OAuthService exposes it as ClientCredentials.
func (c *Client) clientCredentialsToken(ctx context.Context, clientID, clientSecret, scope string, opts ...RequestOption) (*Token, error) {
	form := url.Values{"grant_type": {"client_credentials"}}
	if scope != "" {
		form.Set("scope", scope)
	}
	tok, err := call[Token](ctx, c, &request{
		method:       http.MethodPost,
		path:         "/v1/oauth/token",
		form:         form,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, opts)
	if err != nil {
		return nil, err
	}
	tok.received = time.Now()
	return tok, nil
}
