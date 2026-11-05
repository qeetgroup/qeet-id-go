package qeetid

import (
	"context"
	"net/http"
	"time"
)

// APIKeysService manages API keys.
type APIKeysService struct{ c *Client }

// Create creates an API key. The result carries the key, shown only once.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/api-keys
func (s *APIKeysService) Create(ctx context.Context, params CreateAPIKeyParams, opts ...RequestOption) (*CreateAPIKeyResult, error) {
	return call[CreateAPIKeyResult](ctx, s.c, &request{method: http.MethodPost, path: "/v1/api-keys", body: params, idempotent: true}, opts)
}

// List lists a tenant's API keys.
//
// API: GET /v1/tenants/{tenantID}/api-keys
func (s *APIKeysService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[APIKey], error) {
	return call[Page[APIKey]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/api-keys",
		pathArgs: []string{tenantID},
	}, opts)
}

// Rotate replaces an API key with a new one. The old key keeps working for the overlap given.
//
// API: POST /v1/tenants/{tenantID}/api-keys/{id}/rotate
func (s *APIKeysService) Rotate(ctx context.Context, tenantID string, keyID string, params RotateSecretParams, opts ...RequestOption) (*CreateAPIKeyResult, error) {
	return call[CreateAPIKeyResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/api-keys/{keyID}/rotate",
		pathArgs: []string{tenantID, keyID},
		body:     params,
	}, opts)
}

// Revoke revokes an API key at once.
//
// API: DELETE /v1/api-keys/{id}
func (s *APIKeysService) Revoke(ctx context.Context, keyID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{method: http.MethodDelete, path: "/v1/api-keys/{keyID}", pathArgs: []string{keyID}}, nil, opts)
}

// CreateAPIKeyParams is the request body of APIKeysService.Create.
type CreateAPIKeyParams struct {
	TenantID  string     `json:"tenant_id,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Scopes    []string   `json:"scopes,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKeyResult carries an API key and its secret, which is shown only once.
type CreateAPIKeyResult struct {
	APIKey  APIKey `json:"api_key,omitzero"`
	Secret  string `json:"secret,omitempty"`
	Warning string `json:"warning,omitempty"`
}

// APIKey is an API key. Its secret is shown only when it is created or rotated.
type APIKey struct {
	ID         string     `json:"id,omitempty"`
	TenantID   string     `json:"tenant_id,omitempty"`
	UserID     string     `json:"user_id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Prefix     string     `json:"prefix,omitempty"`
	Scopes     []string   `json:"scopes,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitzero"`
	// The person who created the key, or rotated it into being; null for a key made by a machine or
	// before this was kept.
	CreatedBy string `json:"created_by,omitempty"`
	// The key this one replaced at a rotation.
	RotatedFrom string `json:"rotated_from,omitempty"`
}
