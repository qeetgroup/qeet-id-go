package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type SAMLConnection struct {
	ID               string            `json:"id"`
	TenantID         string            `json:"tenant_id"`
	Name             string            `json:"name"`
	Enabled          bool              `json:"enabled"`
	IdpEntityID      string            `json:"idp_entity_id,omitempty"`
	IdpSSOURL        string            `json:"idp_sso_url,omitempty"`
	IdpCertificate   string            `json:"idp_certificate,omitempty"`
	SpEntityID       string            `json:"sp_entity_id,omitempty"`
	SpACSURL         string            `json:"sp_acs_url,omitempty"`
	AttributeMapping map[string]string `json:"attribute_mapping,omitempty"`
	CreatedAt        string            `json:"created_at"`
	UpdatedAt        string            `json:"updated_at,omitempty"`
}

type CreateSAMLConnectionInput struct {
	Name             string            `json:"name"`
	IdpEntityID      string            `json:"idp_entity_id,omitempty"`
	IdpSSOURL        string            `json:"idp_sso_url,omitempty"`
	IdpCertificate   string            `json:"idp_certificate,omitempty"`
	AttributeMapping map[string]string `json:"attribute_mapping,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
}

type UpdateSAMLConnectionInput struct {
	Name             *string           `json:"name,omitempty"`
	IdpEntityID      *string           `json:"idp_entity_id,omitempty"`
	IdpSSOURL        *string           `json:"idp_sso_url,omitempty"`
	IdpCertificate   *string           `json:"idp_certificate,omitempty"`
	AttributeMapping map[string]string `json:"attribute_mapping,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
}

type SAMLTestResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// SAMLService manages SAML SSO connections.
type SAMLService struct{ t *transport.Transport }

func (s *SAMLService) Create(ctx context.Context, tenantID string, in CreateSAMLConnectionInput) (*SAMLConnection, error) {
	var out SAMLConnection
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml", in, &out)
	return &out, err
}

func (s *SAMLService) Get(ctx context.Context, tenantID, id string) (*SAMLConnection, error) {
	var out SAMLConnection
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *SAMLService) Update(ctx context.Context, tenantID, id string, in UpdateSAMLConnectionInput) (*SAMLConnection, error) {
	var out SAMLConnection
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *SAMLService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml/"+url.PathEscape(id), nil)
}

func (s *SAMLService) Test(ctx context.Context, tenantID, id string) (*SAMLTestResult, error) {
	var out SAMLTestResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml/"+url.PathEscape(id)+"/test", struct{}{}, &out)
	return &out, err
}

func (s *SAMLService) List(ctx context.Context, tenantID string) ([]SAMLConnection, error) {
	var env marshal.Envelope[SAMLConnection]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/saml", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
