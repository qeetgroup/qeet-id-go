package qeetid

import (
	"context"
	"net/http"
	"time"
)

// OIDCConnectionsService manages OpenID Connect connections to a tenant's identity providers.
type OIDCConnectionsService struct{ c *Client }

// List lists a tenant's OpenID Connect connections.
//
// API: GET /v1/tenants/{tenantID}/oidc-connections
func (s *OIDCConnectionsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[OIDCConnection], error) {
	return call[Page[OIDCConnection]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc-connections",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create creates a connection.
//
// API: POST /v1/tenants/{tenantID}/oidc-connections
func (s *OIDCConnectionsService) Create(ctx context.Context, tenantID string, params CreateOIDCConnectionParams, opts ...RequestOption) (*OIDCConnection, error) {
	return call[OIDCConnection](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/oidc-connections",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Get returns a connection.
//
// API: GET /v1/tenants/{tenantID}/oidc-connections/{id}
func (s *OIDCConnectionsService) Get(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*OIDCConnection, error) {
	return call[OIDCConnection](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc-connections/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// Update changes a connection. Fields left nil are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/oidc-connections/{id}
func (s *OIDCConnectionsService) Update(ctx context.Context, tenantID string, connectionID string, params UpdateOIDCConnectionParams, opts ...RequestOption) (*OIDCConnection, error) {
	return call[OIDCConnection](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/oidc-connections/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// Delete deletes a connection.
//
// API: DELETE /v1/tenants/{tenantID}/oidc-connections/{id}
func (s *OIDCConnectionsService) Delete(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/oidc-connections/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
	}, nil, opts)
}

// OIDCConnection connects a tenant to an OpenID Connect identity provider.
type OIDCConnection struct {
	ID       string `json:"id,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Name     string `json:"name,omitempty"`
	// The iss an ID token must carry.
	Issuer string `json:"issuer,omitempty"`
	// Where the discovery document is read; empty for the issuer's /.well-known/openid-configuration.
	DiscoveryURL string   `json:"discovery_url,omitempty"`
	ClientID     string   `json:"client_id,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	// The provider's own tenant IDs (tid) accepted; empty accepts the issuer alone.
	AllowedIDPTenantIDs []string             `json:"allowed_idp_tenant_ids,omitempty"`
	ClaimMapping        OIDCClaimMapping     `json:"claim_mapping,omitzero"`
	Status              OIDCConnectionStatus `json:"status,omitempty"`
	// The redirect URI to register with the provider.
	RedirectURI string     `json:"redirect_uri,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitzero"`
	UpdatedAt   time.Time  `json:"updated_at,omitzero"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// OIDCClaimMapping is the ID token claims a sign-in reads. An empty name reads none, except for email,
// which is required.
type OIDCClaimMapping struct {
	// The claim holding the email address, such as email, preferred_username or upn. The provider's
	// email_verified is consulted only when this is email.
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	FamilyName string `json:"family_name,omitempty"`
	// The claim listing the person's groups, kept on their linked identity.
	Groups string `json:"groups,omitempty"`
}

// OIDCConnectionStatus is a string enumeration used by OIDCConnection.Status.
type OIDCConnectionStatus string

// Values of OIDCConnectionStatus. The server may add values; handle an unknown one gracefully.
const (
	OIDCConnectionStatusDraft    OIDCConnectionStatus = "draft"
	OIDCConnectionStatusActive   OIDCConnectionStatus = "active"
	OIDCConnectionStatusDisabled OIDCConnectionStatus = "disabled"
)

// CreateOIDCConnectionParams is the request body of OIDCConnectionsService.Create.
type CreateOIDCConnectionParams struct {
	Name                string               `json:"name,omitempty"`
	Issuer              string               `json:"issuer,omitempty"`
	DiscoveryURL        string               `json:"discovery_url,omitempty"`
	ClientID            string               `json:"client_id,omitempty"`
	ClientSecret        string               `json:"client_secret,omitempty"`
	Scopes              []string             `json:"scopes,omitempty"`
	AllowedIDPTenantIDs []string             `json:"allowed_idp_tenant_ids,omitempty"`
	ClaimMapping        *OIDCClaimMapping    `json:"claim_mapping,omitempty"`
	Status              OIDCConnectionStatus `json:"status,omitempty"`
}

// UpdateOIDCConnectionParams is the request body of OIDCConnectionsService.Update.
type UpdateOIDCConnectionParams struct {
	Name         *string `json:"name,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	DiscoveryURL *string `json:"discovery_url,omitempty"`
	ClientID     *string `json:"client_id,omitempty"`
	// A new client secret; absent or empty keeps the stored one.
	ClientSecret        *string               `json:"client_secret,omitempty"`
	Scopes              []string              `json:"scopes,omitempty"`
	AllowedIDPTenantIDs []string              `json:"allowed_idp_tenant_ids,omitempty"`
	ClaimMapping        *OIDCClaimMapping     `json:"claim_mapping,omitempty"`
	Status              *OIDCConnectionStatus `json:"status,omitempty"`
}
