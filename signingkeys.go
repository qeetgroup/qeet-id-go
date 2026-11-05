package qeetid

import (
	"context"
	"net/http"
	"time"
)

// SigningKeysService reports the keys Qeet ID signs tokens with.
type SigningKeysService struct{ c *Client }

// List returns metadata for every key that can still verify a token: the active key, candidates and
// retired keys in their grace period. It carries no key material; the public keys are in the JWKS.
//
// API: GET /v1/oidc/signing-keys
func (s *SigningKeysService) List(ctx context.Context, opts ...RequestOption) ([]SigningKey, error) {
	var out struct {
		V []SigningKey `json:"keys"`
	}
	if err := s.c.do(ctx, &request{method: http.MethodGet, path: "/v1/oidc/signing-keys"}, &out, opts); err != nil {
		return nil, err
	}
	return out.V, nil
}

// SigningKey describes a key Qeet ID signs tokens with. It carries no key material; the public keys
// are in the JWKS.
type SigningKey struct {
	// RFC 7638 JWK thumbprint of the public key.
	Kid      string           `json:"kid,omitempty"`
	Alg      string           `json:"alg,omitempty"`
	Use      string           `json:"use,omitempty"`
	Status   SigningKeyStatus `json:"status,omitempty"`
	RetireAt *time.Time       `json:"retire_at,omitempty"`
	// When the key registry first held the key; null for a key it doesn't hold (one only in
	// configuration).
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// When the key became the active signing key; null for a candidate that never was, or a key the
	// registry doesn't hold.
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

// SigningKeyStatus is a string enumeration used by SigningKey.Status.
type SigningKeyStatus string

// Values of SigningKeyStatus. The server may add values; handle an unknown one gracefully.
const (
	SigningKeyStatusCandidate SigningKeyStatus = "candidate"
	SigningKeyStatusActive    SigningKeyStatus = "active"
	SigningKeyStatusRetired   SigningKeyStatus = "retired"
)
