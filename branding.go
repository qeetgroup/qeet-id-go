package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type BrandingSettings struct {
	TenantID       string `json:"tenant_id"`
	LogoURL        string `json:"logo_url,omitempty"`
	PrimaryColor   string `json:"primary_color,omitempty"`
	SecondaryColor string `json:"secondary_color,omitempty"`
	CustomDomain   string `json:"custom_domain,omitempty"`
	FaviconURL     string `json:"favicon_url,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type UpdateBrandingInput struct {
	LogoURL        *string `json:"logo_url,omitempty"`
	PrimaryColor   *string `json:"primary_color,omitempty"`
	SecondaryColor *string `json:"secondary_color,omitempty"`
	CustomDomain   *string `json:"custom_domain,omitempty"`
	FaviconURL     *string `json:"favicon_url,omitempty"`
}

// BrandingService manages a tenant's hosted-login branding.
type BrandingService struct{ t *transport.Transport }

func (s *BrandingService) Get(ctx context.Context, tenantID string) (*BrandingSettings, error) {
	var out BrandingSettings
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/branding", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *BrandingService) Update(ctx context.Context, tenantID string, in UpdateBrandingInput) (*BrandingSettings, error) {
	var out BrandingSettings
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/branding", in, &out)
	return &out, err
}
