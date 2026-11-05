package qeetid

import (
	"context"
	"net/http"
	"time"
)

// SocialProvidersService manages sign-in with social providers and the social accounts linked to
// users.
type SocialProvidersService struct{ c *Client }

// List lists a tenant's social providers.
//
// API: GET /v1/tenants/{tenantID}/social/providers
func (s *SocialProvidersService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[SocialProvider], error) {
	return call[Page[SocialProvider]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/social/providers",
		pathArgs: []string{tenantID},
	}, opts)
}

// Upsert creates or replaces the caller's tenant configuration for a provider.
//
// API: POST /v1/social/providers
func (s *SocialProvidersService) Upsert(ctx context.Context, params UpsertSocialProviderParams, opts ...RequestOption) (*SocialProvider, error) {
	return call[SocialProvider](ctx, s.c, &request{method: http.MethodPost, path: "/v1/social/providers", body: params}, opts)
}

// Update changes a provider's configuration. Fields left nil are not changed.
//
// API: PATCH /v1/social/providers/{provider}
func (s *SocialProvidersService) Update(ctx context.Context, provider string, params UpdateSocialProviderParams, opts ...RequestOption) (*SocialProvider, error) {
	return call[SocialProvider](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/social/providers/{provider}",
		pathArgs: []string{provider},
		body:     params,
	}, opts)
}

// Delete removes a provider.
//
// API: DELETE /v1/social/providers/{provider}
func (s *SocialProvidersService) Delete(ctx context.Context, provider string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/social/providers/{provider}",
		pathArgs: []string{provider},
	}, nil, opts)
}

// UnlinkIdentity unlinks a social account from its user.
//
// API: DELETE /v1/social/identities/{id}
func (s *SocialProvidersService) UnlinkIdentity(ctx context.Context, identityID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/social/identities/{identityID}",
		pathArgs: []string{identityID},
	}, nil, opts)
}

// SocialProvider is a tenant's configuration for signing in with a social provider such as Google or
// GitHub.
type SocialProvider struct {
	// oidc for a provider found by discovery_url, else the adapter its name picks.
	Kind         SocialProviderKind `json:"kind,omitempty"`
	AppleTeamID  string             `json:"apple_team_id,omitempty"`
	AppleKeyID   string             `json:"apple_key_id,omitempty"`
	ClientID     string             `json:"client_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitzero"`
	DiscoveryURL string             `json:"discovery_url,omitempty"`
	Enabled      bool               `json:"enabled"`
	ID           string             `json:"id,omitempty"`
	Provider     string             `json:"provider,omitempty"`
	TenantID     string             `json:"tenant_id,omitempty"`
}

// SocialProviderKind is a string enumeration used by SocialProvider.Kind. oidc for a provider found by
// discovery_url, else the adapter its name picks.
type SocialProviderKind string

// Values of SocialProviderKind. The server may add values; handle an unknown one gracefully.
const (
	SocialProviderKindOIDC     SocialProviderKind = "oidc"
	SocialProviderKindGithub   SocialProviderKind = "github"
	SocialProviderKindFacebook SocialProviderKind = "facebook"
	SocialProviderKindX        SocialProviderKind = "x"
	SocialProviderKindDiscord  SocialProviderKind = "discord"
	SocialProviderKindApple    SocialProviderKind = "apple"
)

// UpsertSocialProviderParams is the request body of SocialProvidersService.Upsert.
type UpsertSocialProviderParams struct {
	Provider string `json:"provider,omitempty"`
	ClientID string `json:"client_id,omitempty"`
	// Required for a new provider. On a provider that exists, empty keeps the stored secret. For Apple,
	// the .p8 private key.
	ClientSecret string `json:"client_secret,omitempty"`
	// An OpenID Connect provider's public https discovery URL; ignored for github, facebook, x, discord
	// and apple.
	DiscoveryURL string `json:"discovery_url,omitempty"`
	// Apple's team ID; required for apple.
	AppleTeamID string `json:"apple_team_id,omitempty"`
	// The ID of the key client_secret holds; required for apple.
	AppleKeyID string `json:"apple_key_id,omitempty"`
}

// UpdateSocialProviderParams is the request body of SocialProvidersService.Update.
type UpdateSocialProviderParams struct {
	// false switches the provider off — it is no longer offered and sign-in through it is refused — until
	// switched on again; other saves keep it off.
	Enabled  *bool   `json:"enabled,omitempty"`
	ClientID *string `json:"client_id,omitempty"`
	// Replaces the stored secret; omitted or empty keeps it. For Apple, a new .p8 key, which must sign
	// with the team and key IDs.
	ClientSecret *string `json:"client_secret,omitempty"`
	// An OpenID Connect provider's alone; ignored for the others.
	DiscoveryURL *string `json:"discovery_url,omitempty"`
	AppleTeamID  *string `json:"apple_team_id,omitempty"`
	AppleKeyID   *string `json:"apple_key_id,omitempty"`
}
