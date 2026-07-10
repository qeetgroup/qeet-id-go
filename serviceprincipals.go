package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// ServicePrincipal is a machine identity for client-credentials (M2M) auth.
// Description is accepted on create but never returned on any response.
type ServicePrincipal struct {
	ID         string   `json:"id"`
	TenantID   string   `json:"tenant_id"`
	Name       string   `json:"name"`
	Scopes     []string `json:"scopes"`
	DisabledAt string   `json:"disabled_at,omitempty"`
	CreatedAt  string   `json:"created_at"`
}

type CreateServicePrincipalInput struct {
	TenantID    string   `json:"tenant_id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

// CreateServicePrincipalResult carries the plaintext client secret, shown
// once at creation.
type CreateServicePrincipalResult struct {
	ServicePrincipal ServicePrincipal `json:"service_principal"`
	ClientID         string           `json:"client_id"`
	ClientSecret     string           `json:"client_secret"`
	Warning          string           `json:"warning,omitempty"`
}

// ServicePrincipalsService manages machine identities. There is no Get or
// Update anywhere in the backend — only create/disable/list.
type ServicePrincipalsService struct{ t *transport.Transport }

func (s *ServicePrincipalsService) Create(ctx context.Context, in CreateServicePrincipalInput) (*CreateServicePrincipalResult, error) {
	var out CreateServicePrincipalResult
	err := s.t.Post(ctx, "/v1/service-principals", in, &out)
	return &out, err
}

// Disable revokes a service principal (called "disable" server-side; there
// is no re-enable — create a new one instead).
func (s *ServicePrincipalsService) Disable(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/service-principals/"+url.PathEscape(id), nil)
}

func (s *ServicePrincipalsService) List(ctx context.Context, tenantID string) ([]ServicePrincipal, error) {
	var env marshal.Envelope[ServicePrincipal]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/service-principals", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
