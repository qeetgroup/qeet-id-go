package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// SAMLServiceProvider is an external service provider registered to consume
// this tenant's SAML assertions — Qeet ID acting as the IdP. This is the
// mirror image of SAMLConnection (saml.go), where Qeet ID is the SP
// connecting out to an external IdP; the two are genuinely separate
// resources sharing only a URL prefix.
type SAMLServiceProvider struct {
	ID              string `json:"id"`
	TenantID        string `json:"tenant_id"`
	Name            string `json:"name"`
	EntityID        string `json:"entity_id"`
	ACSURL          string `json:"acs_url"`
	NameIDFormat    string `json:"name_id_format,omitempty"`
	NameIDAttribute string `json:"name_id_attribute,omitempty"`
	Certificate     string `json:"certificate,omitempty"`
	Status          string `json:"status"` // draft | active | disabled
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	LastLoginAt     string `json:"last_login_at,omitempty"`
}

type CreateSAMLServiceProviderInput struct {
	Name            string `json:"name"`
	EntityID        string `json:"entity_id"`
	ACSURL          string `json:"acs_url"`
	NameIDFormat    string `json:"name_id_format,omitempty"`
	NameIDAttribute string `json:"name_id_attribute,omitempty"`
	Certificate     string `json:"certificate,omitempty"`
	Status          string `json:"status,omitempty"`
}

type UpdateSAMLServiceProviderInput struct {
	Name            *string `json:"name,omitempty"`
	EntityID        *string `json:"entity_id,omitempty"`
	ACSURL          *string `json:"acs_url,omitempty"`
	NameIDFormat    *string `json:"name_id_format,omitempty"`
	NameIDAttribute *string `json:"name_id_attribute,omitempty"`
	Certificate     *string `json:"certificate,omitempty"`
	Status          *string `json:"status,omitempty"` // draft | active | disabled
}

// SAMLServiceProvidersService manages external SPs registered against this
// tenant's SAML IdP.
type SAMLServiceProvidersService struct{ t *transport.Transport }

func (s *SAMLServiceProvidersService) Create(ctx context.Context, tenantID string, in CreateSAMLServiceProviderInput) (*SAMLServiceProvider, error) {
	var out SAMLServiceProvider
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml-providers", in, &out)
	return &out, err
}

func (s *SAMLServiceProvidersService) Get(ctx context.Context, tenantID, id string) (*SAMLServiceProvider, error) {
	var out SAMLServiceProvider
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml-providers/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *SAMLServiceProvidersService) Update(ctx context.Context, tenantID, id string, in UpdateSAMLServiceProviderInput) (*SAMLServiceProvider, error) {
	var out SAMLServiceProvider
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml-providers/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *SAMLServiceProvidersService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml-providers/"+url.PathEscape(id), nil)
}

func (s *SAMLServiceProvidersService) List(ctx context.Context, tenantID string) ([]SAMLServiceProvider, error) {
	var env marshal.Envelope[SAMLServiceProvider]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml-providers", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
