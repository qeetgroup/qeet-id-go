package qeetid

import "time"

// Token is an OAuth 2.0 token response from the Qeet ID token endpoints.
type Token struct {
	// AccessToken is the access token, an ES256-signed JWT.
	AccessToken string `json:"access_token"`
	// TokenType is "Bearer".
	TokenType string `json:"token_type"`
	// ExpiresIn is the access token's lifetime in seconds, counted from when
	// the server issued it.
	ExpiresIn int `json:"expires_in"`
	// ExpiresAt is when the access token expires. Only the client credentials
	// endpoint sends it; Expiry works for every grant.
	ExpiresAt time.Time `json:"expires_at,omitzero"`
	// RefreshToken renews the grant without the user, when one was issued.
	RefreshToken string `json:"refresh_token,omitempty"`
	// IDToken is the OpenID Connect ID token, when the grant included the
	// openid scope.
	IDToken string `json:"id_token,omitempty"`
	// Scope lists the granted scopes, separated by spaces.
	Scope string `json:"scope,omitempty"`
	// IssuedTokenType is the type of the issued token in a token exchange
	// (RFC 8693), such as "urn:ietf:params:oauth:token-type:access_token".
	IssuedTokenType string `json:"issued_token_type,omitempty"`

	received time.Time
}

// Expiry returns when the access token expires: ExpiresAt when the server sent
// it, otherwise the time the response was received plus ExpiresIn. It returns
// the zero time when neither is known.
func (t *Token) Expiry() time.Time {
	if !t.ExpiresAt.IsZero() {
		return t.ExpiresAt
	}
	if t.ExpiresIn > 0 && !t.received.IsZero() {
		return t.received.Add(time.Duration(t.ExpiresIn) * time.Second)
	}
	return time.Time{}
}
