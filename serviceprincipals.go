package qeetid

import (
	"context"
	"net/http"
	"time"
)

// ServicePrincipalsService manages machine identities that get tokens with the OAuth client
// credentials grant.
type ServicePrincipalsService struct{ c *Client }

// Create creates a service principal. The result carries its client secret, shown only once.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/service-principals
func (s *ServicePrincipalsService) Create(ctx context.Context, params CreateServicePrincipalParams, opts ...RequestOption) (*CreateServicePrincipalResult, error) {
	return call[CreateServicePrincipalResult](ctx, s.c, &request{
		method:     http.MethodPost,
		path:       "/v1/service-principals",
		body:       params,
		idempotent: true,
	}, opts)
}

// List lists a tenant's service principals.
//
// API: GET /v1/tenants/{tenantID}/service-principals
func (s *ServicePrincipalsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[ServicePrincipal], error) {
	return call[Page[ServicePrincipal]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/service-principals",
		pathArgs: []string{tenantID},
	}, opts)
}

// Disable disables a service principal. Its tokens stop working at once.
//
// API: DELETE /v1/service-principals/{id}
func (s *ServicePrincipalsService) Disable(ctx context.Context, principalID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/service-principals/{principalID}",
		pathArgs: []string{principalID},
	}, nil, opts)
}

// Enable enables a disabled service principal.
//
// API: POST /v1/tenants/{tenantID}/service-principals/{id}/enable
func (s *ServicePrincipalsService) Enable(ctx context.Context, tenantID, principalID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/service-principals/{principalID}/enable",
		pathArgs: []string{tenantID, principalID},
	}, nil, opts)
}

// RotateSecret replaces a service principal's secret. The old secret keeps working for the overlap
// given.
//
// API: POST /v1/tenants/{tenantID}/service-principals/{id}/rotate-secret
func (s *ServicePrincipalsService) RotateSecret(ctx context.Context, tenantID, principalID string, params RotateSecretParams, opts ...RequestOption) (*RotatedServicePrincipalSecret, error) {
	return call[RotatedServicePrincipalSecret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/service-principals/{principalID}/rotate-secret",
		pathArgs: []string{tenantID, principalID},
		body:     params,
	}, opts)
}

// CreateServicePrincipalParams is the request body of ServicePrincipalsService.Create.
type CreateServicePrincipalParams struct {
	TenantID    string   `json:"tenant_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

// CreateServicePrincipalResult carries a new service principal and its client secret, which is shown
// only once.
type CreateServicePrincipalResult struct {
	ClientID         string           `json:"client_id,omitempty"`
	ClientSecret     string           `json:"client_secret,omitempty"`
	ServicePrincipal ServicePrincipal `json:"service_principal,omitzero"`
	Warning          string           `json:"warning,omitempty"`
}

// ServicePrincipal is a machine identity that gets access tokens with the OAuth client credentials
// grant.
type ServicePrincipal struct {
	CreatedAt   time.Time  `json:"created_at,omitzero"`
	DisabledAt  *time.Time `json:"disabled_at,omitempty"`
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Scopes      []string   `json:"scopes,omitempty"`
	TenantID    string     `json:"tenant_id,omitempty"`
	Description string     `json:"description,omitempty"`
	// When the principal was last issued a token, to the minute.
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	// Present while the secret replaced at its last rotation still works, until then.
	PreviousSecretExpiresAt time.Time `json:"previous_secret_expires_at,omitzero"`
}

// RotateSecretParams is the request body of ServicePrincipalsService.RotateSecret.
type RotateSecretParams struct {
	// How long the old secret goes on working beside the new one. Defaults to 24; 0 stops it at once.
	OverlapHours *int `json:"overlap_hours,omitempty"`
}

// RotatedServicePrincipalSecret is the result of ServicePrincipalsService.RotateSecret.
type RotatedServicePrincipalSecret struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	// When the old secret stops working.
	PreviousSecretExpiresAt time.Time `json:"previous_secret_expires_at,omitzero"`
	Warning                 string    `json:"warning,omitempty"`
}
