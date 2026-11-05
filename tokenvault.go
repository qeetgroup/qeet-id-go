package qeetid

import (
	"context"
	"net/http"
	"time"
)

// TokenVaultService stores third-party OAuth tokens for people and hands out live access tokens on
// their behalf.
type TokenVaultService struct{ c *Client }

// RegisterProvider registers, or replaces, a third-party OAuth provider for the tenant.
//
// API: POST /v1/tenants/{tenantID}/vault/tokens/providers
func (s *TokenVaultService) RegisterProvider(ctx context.Context, tenantID string, params RegisterTokenVaultProviderParams, opts ...RequestOption) (*TokenVaultProvider, error) {
	return call[TokenVaultProvider](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/vault/tokens/providers",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// ListProviders lists the tenant's registered providers.
//
// API: GET /v1/tenants/{tenantID}/vault/tokens/providers
func (s *TokenVaultService) ListProviders(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[TokenVaultProvider], error) {
	return call[Page[TokenVaultProvider]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/vault/tokens/providers",
		pathArgs: []string{tenantID},
	}, opts)
}

// DeleteProvider removes a provider.
//
// API: DELETE /v1/tenants/{tenantID}/vault/tokens/providers/{provider}
func (s *TokenVaultService) DeleteProvider(ctx context.Context, tenantID, provider string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/vault/tokens/providers/{provider}",
		pathArgs: []string{tenantID, provider},
	}, nil, opts)
}

// ListConnections lists the third-party accounts the calling person connected.
//
// API: GET /v1/vault/tokens
func (s *TokenVaultService) ListConnections(ctx context.Context, opts ...RequestOption) (*Page[TokenVaultConnection], error) {
	return call[Page[TokenVaultConnection]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/vault/tokens"}, opts)
}

// Connect begins connecting the calling person's account at a provider. Send the person to
// AuthorizeURL; the provider redirects back to Qeet ID, which stores the tokens.
//
// API: GET /v1/vault/tokens/{provider}/connect
func (s *TokenVaultService) Connect(ctx context.Context, provider string, opts ...RequestOption) (*TokenVaultAuthorization, error) {
	return call[TokenVaultAuthorization](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/vault/tokens/{provider}/connect",
		pathArgs: []string{provider},
	}, opts)
}

// AccessToken returns a live access token for the calling person's connected account, refreshing it
// first when needed. The refresh token is never returned.
//
// API: GET /v1/vault/tokens/{provider}/access-token
func (s *TokenVaultService) AccessToken(ctx context.Context, provider string, opts ...RequestOption) (*TokenVaultAccessToken, error) {
	return call[TokenVaultAccessToken](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/vault/tokens/{provider}/access-token",
		pathArgs: []string{provider},
	}, opts)
}

// Disconnect disconnects the calling person's account at a provider.
//
// API: DELETE /v1/vault/tokens/{provider}
func (s *TokenVaultService) Disconnect(ctx context.Context, provider string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/vault/tokens/{provider}",
		pathArgs: []string{provider},
	}, nil, opts)
}

// RegisterTokenVaultProviderParams is the request body of TokenVaultService.RegisterProvider.
type RegisterTokenVaultProviderParams struct {
	// A short key naming the third-party service, e.g. "slack", "github", "google", or any custom name.
	Provider     string `json:"provider,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	AuthorizeURL string `json:"authorize_url,omitempty"`
	TokenURL     string `json:"token_url,omitempty"`
	Scopes       string `json:"scopes,omitempty"`
}

// TokenVaultProvider is a third-party OAuth provider registered with the token vault. Its client
// secret is never returned.
type TokenVaultProvider struct {
	AuthorizeURL string    `json:"authorize_url,omitempty"`
	ClientID     string    `json:"client_id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitzero"`
	ID           string    `json:"id,omitempty"`
	Provider     string    `json:"provider,omitempty"`
	Scopes       string    `json:"scopes,omitempty"`
	TokenURL     string    `json:"token_url,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitzero"`
}

// TokenVaultConnection is a third-party account a person connected to the token vault, without its
// tokens.
type TokenVaultConnection struct {
	CreatedAt         time.Time `json:"created_at,omitzero"`
	ExpiresAt         time.Time `json:"expires_at,omitzero"`
	ExternalAccountID string    `json:"external_account_id,omitempty"`
	Provider          string    `json:"provider,omitempty"`
	Scope             string    `json:"scope,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitzero"`
}

// TokenVaultAuthorization is the result of TokenVaultService.Connect.
type TokenVaultAuthorization struct {
	AuthorizeURL string `json:"authorize_url,omitempty"`
}

// TokenVaultAccessToken is the result of TokenVaultService.AccessToken.
type TokenVaultAccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
}
