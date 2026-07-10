package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// AdminLink is a delegated-admin-portal link: a time-boxed token that lets
// a tenant's own IT admin configure SAML/SCIM without a Qeet Group console
// account.
type AdminLink struct {
	ID           string   `json:"id"`
	TenantID     string   `json:"tenant_id"`
	Capabilities []string `json:"capabilities"` // "saml" and/or "scim"
	CreatedBy    string   `json:"created_by,omitempty"`
	ExpiresAt    string   `json:"expires_at"`
	RevokedAt    string   `json:"revoked_at,omitempty"`
	LastUsedAt   string   `json:"last_used_at,omitempty"`
	CreatedAt    string   `json:"created_at"`
}

// CreateAdminLinkInput — TTLSeconds is clamped server-side to [15min, 7d],
// defaulting to 24h when zero.
type CreateAdminLinkInput struct {
	Capabilities []string `json:"capabilities"`
	TTLSeconds   int      `json:"ttl_seconds,omitempty"`
}

// CreateAdminLinkResult carries the plaintext token and shareable URL, both
// shown only once at creation.
type CreateAdminLinkResult struct {
	Link  AdminLink `json:"link"`
	Token string    `json:"token"`
	URL   string    `json:"url"`
}

// AdminLinksService manages delegated admin-portal links (renamed from the
// earlier "AdminPortal" naming — this manages links/tokens, not a portal
// itself). The public token-based portal session
// (/admin-portal/{token}/...) is for the external IT admin's own browser
// session and isn't wrapped here.
type AdminLinksService struct{ t *transport.Transport }

func (s *AdminLinksService) Create(ctx context.Context, tenantID string, in CreateAdminLinkInput) (*CreateAdminLinkResult, error) {
	var out CreateAdminLinkResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/admin-portal/links", in, &out)
	return &out, err
}

func (s *AdminLinksService) List(ctx context.Context, tenantID string) ([]AdminLink, error) {
	var env marshal.Envelope[AdminLink]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/admin-portal/links", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *AdminLinksService) Revoke(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/admin-portal/links/"+url.PathEscape(id), nil)
}
