package qeetid

import (
	"context"
	"net/http"
	"time"
)

// OIDCClientsService manages the applications that sign people in with Qeet ID over OpenID Connect.
type OIDCClientsService struct{ c *Client }

// List lists a tenant's applications.
//
// API: GET /v1/tenants/{tenantID}/oidc/clients
func (s *OIDCClientsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[OIDCClient], error) {
	return call[Page[OIDCClient]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc/clients",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create registers an application. The result carries its secret, shown only once.
//
// API: POST /v1/tenants/{tenantID}/oidc/clients
func (s *OIDCClientsService) Create(ctx context.Context, tenantID string, params CreateOIDCClientParams, opts ...RequestOption) (*CreateOIDCClientResult, error) {
	return call[CreateOIDCClientResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/oidc/clients",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Get returns an application.
//
// API: GET /v1/tenants/{tenantID}/oidc/clients/{id}
func (s *OIDCClientsService) Get(ctx context.Context, tenantID string, clientID string, opts ...RequestOption) (*OIDCClient, error) {
	return call[OIDCClient](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc/clients/{clientID}",
		pathArgs: []string{tenantID, clientID},
	}, opts)
}

// Update changes an application. Fields left nil are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/oidc/clients/{id}
func (s *OIDCClientsService) Update(ctx context.Context, tenantID string, clientID string, params UpdateOIDCClientParams, opts ...RequestOption) (*OIDCClient, error) {
	return call[OIDCClient](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/oidc/clients/{clientID}",
		pathArgs: []string{tenantID, clientID},
		body:     params,
	}, opts)
}

// Delete deletes an application.
//
// API: DELETE /v1/tenants/{tenantID}/oidc/clients/{id}
func (s *OIDCClientsService) Delete(ctx context.Context, tenantID string, clientID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/oidc/clients/{clientID}",
		pathArgs: []string{tenantID, clientID},
	}, nil, opts)
}

// RotateSecret replaces an application's secret. The result carries the new secret, shown only once.
//
// API: POST /v1/tenants/{tenantID}/oidc/clients/{id}/rotate-secret
func (s *OIDCClientsService) RotateSecret(ctx context.Context, tenantID string, clientID string, opts ...RequestOption) (*RotatedClientSecret, error) {
	return call[RotatedClientSecret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/oidc/clients/{clientID}/rotate-secret",
		pathArgs: []string{tenantID, clientID},
	}, opts)
}

// AuthFailures lists recent wrong-secret attempts for an application and whether it is locked out.
//
// API: GET /v1/tenants/{tenantID}/oidc/clients/{id}/auth-failures
func (s *OIDCClientsService) AuthFailures(ctx context.Context, tenantID string, clientID string, opts ...RequestOption) (*ClientAuthFailures, error) {
	return call[ClientAuthFailures](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc/clients/{clientID}/auth-failures",
		pathArgs: []string{tenantID, clientID},
	}, opts)
}

// ListShadowAI lists applications with live grants that look like unmanaged AI agents.
//
// API: GET /v1/tenants/{tenantID}/oidc/clients/shadow-ai
func (s *OIDCClientsService) ListShadowAI(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[ShadowAICandidate], error) {
	return call[Page[ShadowAICandidate]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc/clients/shadow-ai",
		pathArgs: []string{tenantID},
	}, opts)
}

// AcknowledgeShadowAI records that an administrator reviewed a shadow-AI candidate.
//
// API: POST /v1/tenants/{tenantID}/oidc/clients/{id}/review
func (s *OIDCClientsService) AcknowledgeShadowAI(ctx context.Context, tenantID string, clientID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/oidc/clients/{clientID}/review",
		pathArgs: []string{tenantID, clientID},
	}, nil, opts)
}

// OIDCClient is an application that signs people in with Qeet ID over OpenID Connect.
type OIDCClient struct {
	ID             string         `json:"id,omitempty"`
	TenantID       string         `json:"tenant_id,omitempty"`
	ClientID       string         `json:"client_id,omitempty"`
	Type           OIDCClientType `json:"type,omitempty"`
	Name           string         `json:"name,omitempty"`
	RedirectURIs   []string       `json:"redirect_uris,omitempty"`
	PostLogoutURIs []string       `json:"post_logout_uris,omitempty"`
	GrantTypes     []string       `json:"grant_types,omitempty"`
	Scopes         []string       `json:"scopes,omitempty"`
	CreatedAt      time.Time      `json:"created_at,omitzero"`
	// When this confidential application's exemption from sending a PKCE code_challenge ends; null for
	// none.
	PKCEExemptUntil *time.Time `json:"pkce_exempt_until,omitempty"`
	// Receives a signed logout token (OpenID Connect Back-Channel Logout, with sid) when someone the
	// application signed in signs out; "" for none.
	BackchannelLogoutURI string `json:"backchannel_logout_uri,omitempty"`
	// Loaded by the hosted signed-out page (with iss and sid) when someone the application signed in
	// signs out; "" for none.
	FrontchannelLogoutURI string `json:"frontchannel_logout_uri,omitempty"`
	// The application's logo, an https URL, shown on the hosted sign-in page; "" for none (set "" to
	// clear).
	LogoURI string `json:"logo_uri,omitempty"`
	// The application's home page, an https URL; the hosted sign-in page shows its host. "" for none (set
	// "" to clear).
	ClientURI string `json:"client_uri,omitempty"`
	// The origins a browser application calls the token, userinfo, revocation and introspection endpoints
	// from, which CORS then allows (without credentials): those the application set, else the origins of
	// its http(s) redirect URIs.
	WebOrigins []string `json:"web_origins,omitempty"`
}

// OIDCClientType is a string enumeration used by OIDCClient.Type.
type OIDCClientType string

// Values of OIDCClientType. The server may add values; handle an unknown one gracefully.
const (
	OIDCClientTypePublic       OIDCClientType = "public"
	OIDCClientTypeConfidential OIDCClientType = "confidential"
)

// CreateOIDCClientParams is the request body of OIDCClientsService.Create.
type CreateOIDCClientParams struct {
	TenantID       string         `json:"tenant_id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Type           OIDCClientType `json:"type,omitempty"`
	RedirectURIs   []string       `json:"redirect_uris,omitempty"`
	PostLogoutURIs []string       `json:"post_logout_uris,omitempty"`
	// Any of authorization_code, refresh_token, urn:ietf:params:oauth:grant-type:device_code, token
	// exchange (token_exchange or urn:ietf:params:oauth:grant-type:token-exchange; confidential only) and
	// urn:openid:params:grant-type:ciba (confidential only). client_credentials is refused with
	// oidc.grant_type_unsupported: a machine caller is a service principal, which gets its tokens from
	// POST /v1/oauth/token.
	GrantTypes []string `json:"grant_types,omitempty"`
	// The scopes the application may request. offline_access needs no entry: an application registered
	// for refresh_token may always request it. Listing offline_access opts the application in to OpenID
	// Connect Core §11: it then gets a refresh token only for a request granted offline_access, where
	// otherwise every grant to an application registered for refresh_token comes with one.
	Scopes []string `json:"scopes,omitempty"`
	// Receives a signed logout token (OpenID Connect Back-Channel Logout, with sid) when someone the
	// application signed in signs out; "" for none.
	BackchannelLogoutURI string `json:"backchannel_logout_uri,omitempty"`
	// Loaded by the hosted signed-out page (with iss and sid) when someone the application signed in
	// signs out; "" for none.
	FrontchannelLogoutURI string `json:"frontchannel_logout_uri,omitempty"`
	// The application's logo, an https URL, shown on the hosted sign-in page; "" for none (set "" to
	// clear).
	LogoURI string `json:"logo_uri,omitempty"`
	// The application's home page, an https URL; the hosted sign-in page shows its host. "" for none (set
	// "" to clear).
	ClientURI string `json:"client_uri,omitempty"`
	// Origins (scheme, host and optional port: https, or http only for localhost) a browser application
	// calls the OAuth endpoints from. Left out or null: its redirect URIs' origins. An entry that isn't
	// one is refused with 422 oidc.web_origin_invalid.
	WebOrigins []string `json:"web_origins,omitempty"`
}

// CreateOIDCClientResult carries a new application and its client secret, which is shown only once.
type CreateOIDCClientResult struct {
	Client       OIDCClient `json:"client,omitzero"`
	ClientSecret string     `json:"client_secret,omitempty"`
	Warning      string     `json:"warning,omitempty"`
}

// UpdateOIDCClientParams is the request body of OIDCClientsService.Update. Partial (COALESCE-style)
// update; omitted fields are left unchanged.
type UpdateOIDCClientParams struct {
	Name           *string  `json:"name,omitempty"`
	RedirectURIs   []string `json:"redirect_uris,omitempty"`
	PostLogoutURIs []string `json:"post_logout_uris,omitempty"`
	// Any of authorization_code, refresh_token, urn:ietf:params:oauth:grant-type:device_code, token
	// exchange (token_exchange or urn:ietf:params:oauth:grant-type:token-exchange; confidential only) and
	// urn:openid:params:grant-type:ciba (confidential only). client_credentials is refused with
	// oidc.grant_type_unsupported: a machine caller is a service principal, which gets its tokens from
	// POST /v1/oauth/token.
	GrantTypes []string `json:"grant_types,omitempty"`
	// The scopes the application may request. offline_access needs no entry: an application registered
	// for refresh_token may always request it. Listing offline_access opts the application in to OpenID
	// Connect Core §11: it then gets a refresh token only for a request granted offline_access, where
	// otherwise every grant to an application registered for refresh_token comes with one.
	Scopes []string `json:"scopes,omitempty"`
	// Sets (a date-time, within 90 days) or clears (null) a confidential application's PKCE exemption;
	// omitted, unchanged.
	PKCEExemptUntil *time.Time `json:"pkce_exempt_until,omitempty"`
	// Receives a signed logout token (OpenID Connect Back-Channel Logout, with sid) when someone the
	// application signed in signs out; "" for none.
	BackchannelLogoutURI *string `json:"backchannel_logout_uri,omitempty"`
	// Loaded by the hosted signed-out page (with iss and sid) when someone the application signed in
	// signs out; "" for none.
	FrontchannelLogoutURI *string `json:"frontchannel_logout_uri,omitempty"`
	// The application's logo, an https URL, shown on the hosted sign-in page; "" for none (set "" to
	// clear).
	LogoURI *string `json:"logo_uri,omitempty"`
	// The application's home page, an https URL; the hosted sign-in page shows its host. "" for none (set
	// "" to clear).
	ClientURI *string `json:"client_uri,omitempty"`
	// Sets the origins a browser application calls the OAuth endpoints from; null goes back to its
	// redirect URIs' origins; left out, unchanged. An entry that isn't an origin is refused with 422
	// oidc.web_origin_invalid.
	WebOrigins []string `json:"web_origins,omitempty"`
}

// RotatedClientSecret is the result of OIDCClientsService.RotateSecret.
type RotatedClientSecret struct {
	ClientSecret string `json:"client_secret,omitempty"`
	Warning      string `json:"warning,omitempty"`
}

// ClientAuthFailures is the result of OIDCClientsService.AuthFailures.
type ClientAuthFailures struct {
	LockedUntil *time.Time          `json:"locked_until,omitempty"`
	Items       []ClientAuthFailure `json:"items,omitempty"`
}

// ClientAuthFailure is one element of ClientAuthFailures.Items.
type ClientAuthFailure struct {
	At     time.Time `json:"at,omitzero"`
	IP     string    `json:"ip,omitempty"`
	Reason string    `json:"reason,omitempty"`
}

// ShadowAICandidate is one item of the list OIDCClientsService.ListShadowAI returns.
type ShadowAICandidate struct {
	ClientID   string    `json:"client_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitzero"`
	GrantTypes []string  `json:"grant_types,omitempty"`
	ID         string    `json:"id,omitempty"`
	LiveGrants int       `json:"live_grants"`
	Name       string    `json:"name,omitempty"`
}
