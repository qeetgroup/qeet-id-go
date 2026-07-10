package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type OIDCClient struct {
	ID                      string   `json:"id"`
	TenantID                string   `json:"tenant_id,omitempty"`
	Name                    string   `json:"name"`
	ClientID                string   `json:"client_id"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	Scopes                  []string `json:"scopes"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at,omitempty"`
}

type CreateOIDCClientInput struct {
	Name                    string   `json:"name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types,omitempty"`
	Scopes                  []string `json:"scopes,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
}

type UpdateOIDCClientInput struct {
	Name                    *string  `json:"name,omitempty"`
	RedirectURIs            []string `json:"redirect_uris,omitempty"`
	GrantTypes              []string `json:"grant_types,omitempty"`
	Scopes                  []string `json:"scopes,omitempty"`
	TokenEndpointAuthMethod *string  `json:"token_endpoint_auth_method,omitempty"`
}

type OIDCRotateSecretResult struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// ShadowAIClient is an OIDC client with live grants that isn't registered
// as a managed agent/service principal — a candidate for Shadow-AI review.
type ShadowAIClient struct {
	ClientID    string `json:"client_id"`
	Name        string `json:"name"`
	GrantCount  int    `json:"grant_count,omitempty"`
	LastUsedAt  string `json:"last_used_at,omitempty"`
	FirstSeenAt string `json:"first_seen_at,omitempty"`
}

// OIDCService manages OIDC clients. All CRUD is tenant-scoped — the previous
// SDK version targeted a top-level /v1/oidc/clients path that only supports
// register (POST) and delete in the current API; Get/List/Update/RotateSecret
// live exclusively under /v1/tenants/{tenantID}/oidc/clients and would have
// 404ed.
type OIDCService struct{ t *transport.Transport }

func (s *OIDCService) Create(ctx context.Context, tenantID string, in CreateOIDCClientInput) (*OIDCClient, error) {
	var out OIDCClient
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients", in, &out)
	return &out, err
}

func (s *OIDCService) Get(ctx context.Context, tenantID, id string) (*OIDCClient, error) {
	var out OIDCClient
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *OIDCService) Update(ctx context.Context, tenantID, id string, in UpdateOIDCClientInput) (*OIDCClient, error) {
	var out OIDCClient
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *OIDCService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients/"+url.PathEscape(id), nil)
}

func (s *OIDCService) RotateSecret(ctx context.Context, tenantID, id string) (*OIDCRotateSecretResult, error) {
	var out OIDCRotateSecretResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients/"+url.PathEscape(id)+"/rotate-secret", struct{}{}, &out)
	return &out, err
}

func (s *OIDCService) List(ctx context.Context, tenantID string) ([]OIDCClient, error) {
	var env marshal.Envelope[OIDCClient]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// ListShadowAI lists OIDC clients with live grants that aren't registered as
// managed agents/service principals — flagging unmanaged AI/automation access.
func (s *OIDCService) ListShadowAI(ctx context.Context, tenantID string) ([]ShadowAIClient, error) {
	var env marshal.Envelope[ShadowAIClient]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients/shadow-ai", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// Review marks a shadow-AI client as reviewed.
func (s *OIDCService) Review(ctx context.Context, tenantID, id string) error {
	return s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oidc/clients/"+url.PathEscape(id)+"/review", struct{}{}, nil)
}
