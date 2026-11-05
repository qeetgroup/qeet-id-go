package qeetid

import (
	"context"
	"net/http"
	"time"
)

// AdminPortalLinksService creates links that let a customer's IT administrator set up single sign-on
// and provisioning without a console account.
type AdminPortalLinksService struct{ c *Client }

// Create creates a link. The result carries its URL, shown only once.
//
// API: POST /v1/tenants/{tenantID}/admin-portal/links
func (s *AdminPortalLinksService) Create(ctx context.Context, tenantID string, params CreateAdminPortalLinkParams, opts ...RequestOption) (*CreatedAdminPortalLink, error) {
	return call[CreatedAdminPortalLink](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/admin-portal/links",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// List lists a tenant's links. Page.Truncated reports a list cut short by the server's cap.
//
// API: GET /v1/tenants/{tenantID}/admin-portal/links
func (s *AdminPortalLinksService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[AdminPortalLink], error) {
	return call[Page[AdminPortalLink]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/admin-portal/links",
		pathArgs: []string{tenantID},
	}, opts)
}

// Revoke revokes a link.
//
// API: DELETE /v1/tenants/{tenantID}/admin-portal/links/{id}
func (s *AdminPortalLinksService) Revoke(ctx context.Context, tenantID, linkID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/admin-portal/links/{linkID}",
		pathArgs: []string{tenantID, linkID},
	}, nil, opts)
}

// CreateAdminPortalLinkParams is the request body of AdminPortalLinksService.Create.
type CreateAdminPortalLinkParams struct {
	Capabilities []AdminPortalCapability `json:"capabilities,omitempty"`
	// How long the link lasts, in seconds, from 900 to 86400. Omitted or 0 means an hour.
	TTLSeconds *int `json:"ttl_seconds,omitempty"`
}

// AdminPortalCapability is a string enumeration used by CreateAdminPortalLinkParams.Capabilities.
type AdminPortalCapability string

// Values of AdminPortalCapability. The server may add values; handle an unknown one gracefully.
const (
	AdminPortalCapabilitySAML               AdminPortalCapability = "saml"
	AdminPortalCapabilitySCIM               AdminPortalCapability = "scim"
	AdminPortalCapabilityDomainVerification AdminPortalCapability = "domain_verification"
)

// CreatedAdminPortalLink is the result of AdminPortalLinksService.Create.
type CreatedAdminPortalLink struct {
	Link AdminPortalLink `json:"link,omitzero"`
	// The raw token — shown only here. Not recoverable afterward.
	Token string `json:"token,omitempty"`
	URL   string `json:"url,omitempty"`
}

// AdminPortalLink is a time-limited link that lets a customer's IT administrator configure single
// sign-on and provisioning.
type AdminPortalLink struct {
	ID           string                  `json:"id,omitempty"`
	TenantID     string                  `json:"tenant_id,omitempty"`
	Capabilities []AdminPortalCapability `json:"capabilities,omitempty"`
	CreatedBy    string                  `json:"created_by,omitempty"`
	ExpiresAt    time.Time               `json:"expires_at,omitzero"`
	RevokedAt    *time.Time              `json:"revoked_at,omitempty"`
	LastUsedAt   *time.Time              `json:"last_used_at,omitempty"`
	CreatedAt    time.Time               `json:"created_at,omitzero"`
}
