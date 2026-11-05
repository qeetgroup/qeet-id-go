package qeetid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LDAPAuthenticateParams is the request body of LDAPService.Authenticate.
type LDAPAuthenticateParams struct {
	// Username is the person's directory username.
	Username string `json:"username"`
	// Password is the person's directory password. It is sent to Qeet ID,
	// which checks it with a bind against the directory; it is never logged.
	Password string `json:"password"`
}

// SignInResult is the outcome of a sign-in that may need a second step. It is
// either a *TokenPair, when the person is signed in, or an *MFAChallenge, when
// they must complete a second factor first. Use a type switch:
//
//	switch r := res.(type) {
//	case *qeetid.TokenPair:
//		// signed in: r.AccessToken, r.RefreshToken
//	case *qeetid.MFAChallenge:
//		// complete r.MFAToken at POST /v1/auth/mfa
//	}
type SignInResult interface {
	signInResult()
}

// TokenPair is a signed-in session: an access token and the refresh token
// that renews it.
type TokenPair struct {
	// AccessToken is the ES256-signed access token.
	AccessToken string `json:"access_token"`
	// TokenType is "Bearer".
	TokenType string `json:"token_type"`
	// ExpiresAt is when the access token expires.
	ExpiresAt time.Time `json:"expires_at"`
	// RefreshToken renews the session.
	RefreshToken string `json:"refresh_token"`
	// SessionExpiresAt is when the session ends however active it is.
	SessionExpiresAt time.Time `json:"session_expires_at"`
	// IdleExpiresAt is when the session ends unless it is renewed first. It
	// is never after SessionExpiresAt.
	IdleExpiresAt time.Time `json:"idle_expires_at"`
	// SessionID identifies the session.
	SessionID string `json:"session_id"`
	// UserID identifies the person signed in.
	UserID string `json:"user_id"`
	// TenantID is the tenant the token is scoped to; empty for a session
	// without one.
	TenantID string `json:"tenant_id,omitempty"`
}

// MFAChallenge is a sign-in stopped before its session: the person holds a
// second factor, or must set one up, and completes the challenge at POST
// /v1/auth/mfa with MFAToken.
type MFAChallenge struct {
	// MFARequired is always true on a challenge.
	MFARequired bool `json:"mfa_required"`
	// MFAToken identifies the challenge when completing it.
	MFAToken string `json:"mfa_token"`
	// Methods lists the distinct factor types the person can use, such as
	// "totp", "recovery_code", "email", "sms", "passkey" or "push".
	Methods []string `json:"methods,omitempty"`
	// Factors lists the second factors the person can complete the challenge
	// with. It is empty for an enrollment challenge.
	Factors []LoginFactor `json:"factors,omitempty"`
	// EnrollmentRequired reports that the person must set up an authenticator
	// app and holds no factor yet.
	EnrollmentRequired bool `json:"enrollment_required,omitempty"`
	// EnrollmentDueAt is when the grace period for setting up a factor ends,
	// on an enrollment challenge that may still be put off.
	EnrollmentDueAt *time.Time `json:"enrollment_due_at,omitempty"`
}

// LoginFactor is one way to complete an MFAChallenge.
type LoginFactor struct {
	// Type is "totp", "recovery_code", "email", "sms", "passkey" or "push".
	Type string `json:"type"`
	// ID names an email or SMS factor.
	ID string `json:"id,omitempty"`
	// Destination is an email or SMS factor's masked address or number.
	Destination string `json:"destination,omitempty"`
}

func (*TokenPair) signInResult()    {}
func (*MFAChallenge) signInResult() {}

// Authenticate signs a person in with their directory username and password
// through an LDAP connection. It returns a *TokenPair when the person is signed
// in, or an *MFAChallenge when they must complete a second factor first; see
// SignInResult. A wrong username or password is an *APIError with code
// ldap.invalid_credentials (HTTP 401).
//
// The endpoint is public: the client sends no credentials with this request.
//
// API: POST /ldap/{id}/authenticate
func (s *LDAPService) Authenticate(ctx context.Context, connectionID string, params LDAPAuthenticateParams, opts ...RequestOption) (SignInResult, error) {
	var raw json.RawMessage
	if err := s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/ldap/{connectionID}/authenticate",
		pathArgs: []string{connectionID},
		body:     params,
		public:   true,
	}, &raw, opts); err != nil {
		return nil, err
	}
	return decodeSignInResult(raw)
}

// decodeSignInResult tells a token pair from a second-factor challenge: only a
// challenge carries mfa_token.
func decodeSignInResult(raw json.RawMessage) (SignInResult, error) {
	var probe struct {
		MFAToken    string `json:"mfa_token"`
		MFARequired bool   `json:"mfa_required"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil, fmt.Errorf("qeetid: decode sign-in result: %w", err)
	}
	if probe.MFAToken != "" || probe.MFARequired {
		var ch MFAChallenge
		if err := json.Unmarshal(raw, &ch); err != nil {
			return nil, fmt.Errorf("qeetid: decode sign-in result: %w", err)
		}
		return &ch, nil
	}
	var pair TokenPair
	if err := json.Unmarshal(raw, &pair); err != nil {
		return nil, fmt.Errorf("qeetid: decode sign-in result: %w", err)
	}
	return &pair, nil
}
