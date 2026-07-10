package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// APIKey is a server-side management-API key. There is no per-key Get or
// Rotate in the backend — only Create, Delete (revoke), and List.
type APIKey struct {
	ID         string   `json:"id"`
	TenantID   string   `json:"tenant_id"`
	UserID     string   `json:"user_id,omitempty"`
	Name       string   `json:"name"`
	Prefix     string   `json:"prefix"`
	Scopes     []string `json:"scopes,omitempty"`
	ExpiresAt  string   `json:"expires_at,omitempty"`
	LastUsedAt string   `json:"last_used_at,omitempty"`
	RevokedAt  string   `json:"revoked_at,omitempty"`
	CreatedAt  string   `json:"created_at"`
}

// CreateAPIKeyInput — TenantID and Name are required. Unlike most other
// resources, TenantID here is genuinely required in the body (not derived
// from the caller's own key) — you're creating a key for a tenant, so it
// can't be implicit. ExpiresAt is an RFC3339 timestamp, not a day offset.
type CreateAPIKeyInput struct {
	TenantID  string   `json:"tenant_id"`
	UserID    string   `json:"user_id,omitempty"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
}

// CreateAPIKeyResult carries the plaintext secret, shown once at creation.
type CreateAPIKeyResult struct {
	APIKey  APIKey `json:"api_key"`
	Secret  string `json:"secret"`
	Warning string `json:"warning,omitempty"`
}

// APIKeysService manages server-side API keys.
type APIKeysService struct{ t *transport.Transport }

func (s *APIKeysService) Create(ctx context.Context, in CreateAPIKeyInput) (*CreateAPIKeyResult, error) {
	var out CreateAPIKeyResult
	err := s.t.Post(ctx, "/v1/api-keys", in, &out)
	return &out, err
}

// Delete revokes an API key immediately.
func (s *APIKeysService) Delete(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/api-keys/"+url.PathEscape(id), nil)
}

func (s *APIKeysService) List(ctx context.Context, tenantID string) ([]APIKey, error) {
	var env marshal.Envelope[APIKey]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/api-keys", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
