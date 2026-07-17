package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
	"github.com/qeetgroup/qeet-id-go/internal/validation"
)

// Provider is a tenant's registered OAuth2 endpoint config for one third-party
// service. ClientSecret is never returned by the API.
type Provider struct {
	ID           string `json:"id"`
	Provider     string `json:"provider"`
	ClientID     string `json:"client_id"`
	AuthorizeURL string `json:"authorize_url"`
	TokenURL     string `json:"token_url"`
	Scopes       string `json:"scopes"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type RegisterProviderInput struct {
	Provider     string `json:"provider"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthorizeURL string `json:"authorize_url"`
	TokenURL     string `json:"token_url"`
	Scopes       string `json:"scopes,omitempty"`
}

// GrantMeta is the non-secret view of a connected account.
type GrantMeta struct {
	Provider          string `json:"provider"`
	ExternalAccountID string `json:"external_account_id,omitempty"`
	Scope             string `json:"scope,omitempty"`
	ExpiresAt         string `json:"expires_at,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// AccessTokenResult is the wire shape of the live-access-token response.
type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
}

// TokenVaultService is a per-tenant encrypted store for third-party OAuth
// tokens (Slack, GitHub, Google, or any custom OAuth2 provider an admin
// registers) — distinct from VaultService, which is a generic encrypted-secrets
// store. A user connects their account once via a standard authorization-code
// ceremony; from then on a caller (typically an AI agent or backend
// integration acting on that user's behalf) asks for a live access token via
// GetAccessToken and never sees — or needs to handle — the underlying refresh
// token.
//
// The browser-redirect OAuth-dance endpoints (GET
// /v1/vault/tokens/{provider}/connect, which starts the ceremony, and the GET
// /v1/vault/tokens/callback return leg) are end-user browser flows and are
// intentionally not wrapped here — the same exclusion reasoning used throughout
// this SDK for login/signup/consent-redirect endpoints.
type TokenVaultService struct{ t *transport.Transport }

// RegisterProvider registers (or, on conflict, updates) a third-party OAuth
// provider config for the tenant.
func (s *TokenVaultService) RegisterProvider(ctx context.Context, tenantID string, in RegisterProviderInput) (*Provider, error) {
	var out Provider
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/vault/tokens/providers", in, &out)
	return &out, err
}

// ListProviders lists the tenant's registered third-party OAuth providers.
func (s *TokenVaultService) ListProviders(ctx context.Context, tenantID string) ([]Provider, error) {
	var env marshal.Envelope[Provider]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/vault/tokens/providers", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// DeleteProvider removes a third-party OAuth provider config from the tenant.
func (s *TokenVaultService) DeleteProvider(ctx context.Context, tenantID, provider string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/vault/tokens/providers/"+url.PathEscape(provider), nil)
}

// ListGrants lists the caller's own connected-account grants. Tenant and user
// are resolved from the authenticated principal (never from a path segment or
// body), matching the implicit-scoping pattern used elsewhere (e.g.
// WebhooksService.Create/Get/Delete).
func (s *TokenVaultService) ListGrants(ctx context.Context) ([]GrantMeta, error) {
	var env marshal.Envelope[GrantMeta]
	err := s.t.Get(ctx, "/v1/vault/tokens", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// GetAccessToken returns a live access token for provider, scoped to the
// authenticated caller's own tenant/user context, transparently refreshing it
// first if it's expired (or about to be) and a refresh token is on file. The
// raw refresh token itself is never returned.
func (s *TokenVaultService) GetAccessToken(ctx context.Context, provider string) (*AccessTokenResult, error) {
	if err := validation.Required("provider", provider); err != nil {
		return nil, err
	}
	var out AccessTokenResult
	err := s.t.Get(ctx, "/v1/vault/tokens/"+url.PathEscape(provider)+"/access-token", transport.RequestOptions{}, &out)
	return &out, err
}

// Disconnect disconnects the authenticated caller's own connected account for
// provider.
func (s *TokenVaultService) Disconnect(ctx context.Context, provider string) error {
	if err := validation.Required("provider", provider); err != nil {
		return err
	}
	return s.t.Delete(ctx, "/v1/vault/tokens/"+url.PathEscape(provider), nil)
}
