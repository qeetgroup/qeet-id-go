package qeetid

import (
	"context"
	"net/http"
	"time"
)

// OAuthResourcesService manages the resource servers (RFC 8707 resource indicators) tokens can be
// issued for.
type OAuthResourcesService struct{ c *Client }

// List lists a tenant's resources.
//
// API: GET /v1/tenants/{tenantID}/oidc/resources
func (s *OAuthResourcesService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[OAuthResource], error) {
	return call[Page[OAuthResource]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oidc/resources",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create registers a resource.
//
// API: POST /v1/tenants/{tenantID}/oidc/resources
func (s *OAuthResourcesService) Create(ctx context.Context, tenantID string, params CreateOAuthResourceParams, opts ...RequestOption) (*OAuthResource, error) {
	return call[OAuthResource](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/oidc/resources",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Delete removes a resource.
//
// API: DELETE /v1/tenants/{tenantID}/oidc/resources/{id}
func (s *OAuthResourcesService) Delete(ctx context.Context, tenantID, resourceID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/oidc/resources/{resourceID}",
		pathArgs: []string{tenantID, resourceID},
	}, nil, opts)
}

// OAuthResource is a resource server that access tokens can be issued for (RFC 8707).
type OAuthResource struct {
	ID             string    `json:"id,omitempty"`
	TenantID       string    `json:"tenant_id,omitempty"`
	Identifier     string    `json:"identifier,omitempty"`
	Name           string    `json:"name,omitempty"`
	ServerClientID string    `json:"server_client_id,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitzero"`
}

// CreateOAuthResourceParams is the request body of OAuthResourcesService.Create.
type CreateOAuthResourceParams struct {
	Identifier     string `json:"identifier,omitempty"`
	Name           string `json:"name,omitempty"`
	ServerClientID string `json:"server_client_id,omitempty"`
}
