package jwt

import (
	"encoding/json"
	"slices"
	"strings"
	"time"
)

// Claims are the verified claims of a Qeet ID access token.
type Claims struct {
	// Issuer is the iss claim.
	Issuer string
	// Subject is the sub claim: the person's user ID, or the service
	// principal's ID for a service token.
	Subject string
	// Audience lists the aud claim's values.
	Audience []string
	// ExpiresAt, NotBefore and IssuedAt are the exp, nbf and iat claims.
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
	// ID is the jti claim.
	ID string

	// UserID is the person the token was issued for; empty for a service
	// token.
	UserID string
	// TenantID is the tenant the token is scoped to.
	TenantID string
	// SessionID is the session the token belongs to (the sid claim). It is the
	// SessionID of a session.revoked webhook event.
	SessionID string
	// Scope lists the granted scopes, separated by spaces.
	Scope string
	// ActorType is "service" for a service principal's token and "agent" for
	// an AI agent's; empty for a person's.
	ActorType string
	// AgentID identifies the agent of an agent token.
	AgentID string
	// Actor is the delegate named by a token exchange's act claim.
	Actor *Actor
	// ClientID is the application an OAuth access token was issued to.
	ClientID string
	// ExchangeDepth counts the token exchanges that produced the token.
	ExchangeDepth int
	// Custom holds the claims an authentication hook added, from the token's
	// claims object.
	Custom map[string]any

	// Raw is the token's decoded payload, for claims not listed above.
	Raw json.RawMessage
}

// Actor is the acting party of a delegated token (RFC 8693 section 4.1).
type Actor struct {
	// Subject is the actor's subject.
	Subject string `json:"sub"`
}

// Scopes returns the granted scopes.
func (c *Claims) Scopes() []string { return strings.Fields(c.Scope) }

// HasScope reports whether scope was granted.
func (c *Claims) HasScope(scope string) bool { return slices.Contains(c.Scopes(), scope) }

// HasAudience reports whether audience is among the token's audiences.
func (c *Claims) HasAudience(audience string) bool { return slices.Contains(c.Audience, audience) }

func parseClaims(payloadB64 string) (*Claims, error) {
	var raw struct {
		Issuer        string          `json:"iss"`
		Subject       string          `json:"sub"`
		Audience      json.RawMessage `json:"aud"`
		ExpiresAt     json.Number     `json:"exp"`
		NotBefore     json.Number     `json:"nbf"`
		IssuedAt      json.Number     `json:"iat"`
		ID            string          `json:"jti"`
		UserID        string          `json:"user_id"`
		TenantID      string          `json:"tenant_id"`
		SessionID     string          `json:"sid"`
		Scope         string          `json:"scope"`
		ActorType     string          `json:"actor_type"`
		AgentID       string          `json:"agent_id"`
		Actor         *Actor          `json:"act"`
		ClientID      string          `json:"client_id"`
		ExchangeDepth int             `json:"exchange_depth"`
		Custom        map[string]any  `json:"claims"`
	}
	if err := decodeSegment(payloadB64, &raw); err != nil {
		return nil, err
	}
	aud, err := audiences(raw.Audience)
	if err != nil {
		return nil, ErrMalformed
	}
	exp, err1 := numericDate(raw.ExpiresAt)
	nbf, err2 := numericDate(raw.NotBefore)
	iat, err3 := numericDate(raw.IssuedAt)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, ErrMalformed
	}
	payload, _ := decodeRaw(payloadB64)
	return &Claims{
		Issuer: raw.Issuer, Subject: raw.Subject, Audience: aud,
		ExpiresAt: exp, NotBefore: nbf, IssuedAt: iat, ID: raw.ID,
		UserID: raw.UserID, TenantID: raw.TenantID, SessionID: raw.SessionID, Scope: raw.Scope,
		ActorType: raw.ActorType, AgentID: raw.AgentID, Actor: raw.Actor, ClientID: raw.ClientID,
		ExchangeDepth: raw.ExchangeDepth, Custom: raw.Custom, Raw: payload,
	}, nil
}

// audiences decodes an aud claim, which is a string or an array of strings.
func audiences(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var one string
	if err := json.Unmarshal(raw, &one); err == nil {
		return []string{one}, nil
	}
	var many []string
	if err := json.Unmarshal(raw, &many); err != nil {
		return nil, err
	}
	return many, nil
}

// numericDate decodes a NumericDate claim (seconds since the Unix epoch,
// possibly fractional).
func numericDate(n json.Number) (time.Time, error) {
	if n == "" {
		return time.Time{}, nil
	}
	f, err := n.Float64()
	if err != nil {
		return time.Time{}, err
	}
	secs := int64(f)
	return time.Unix(secs, int64((f-float64(secs))*1e9)), nil
}

func decodeRaw(seg string) (json.RawMessage, error) {
	var m json.RawMessage
	err := decodeSegment(seg, &m)
	return m, err
}
