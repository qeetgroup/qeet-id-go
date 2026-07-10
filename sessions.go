package qeetid

import (
	"context"
	"net/http"

	"github.com/qeetgroup/qeet-id-go/internal/auth"
)

// SessionsService verifies ES256 tokens against the issuer's published JWKS.
// After the keys are cached it is fully local — no network round-trip per
// request. It doesn't use the API-key-authed transport at all (JWKS is a
// public endpoint), so it's constructed separately from every other service.
type SessionsService struct {
	verifier *auth.JWKSVerifier
}

func newSessionsService(baseURL string, hc *http.Client) *SessionsService {
	return &SessionsService{verifier: auth.NewJWKSVerifier(baseURL+"/.well-known/jwks.json", hc)}
}

// setJWKSURL repoints verification at a new JWKS endpoint (e.g. one resolved
// from discovery) and drops any cached keys so the next Verify refetches.
func (s *SessionsService) setJWKSURL(u string) {
	s.verifier.SetJWKSURL(u)
}

// Verify checks the token's ES256 signature against the JWKS, then validates
// expiry/issuer/audience. Returns *Claims or an error.
func (s *SessionsService) Verify(ctx context.Context, token string, opts ...VerifyOptions) (*Claims, error) {
	var o VerifyOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	return s.verifier.Verify(ctx, token, o)
}
