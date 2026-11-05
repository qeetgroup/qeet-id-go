package qeetid

import (
	"context"
	"net/http"
	"time"
)

// LDAPService manages connections to a tenant's LDAP or Active Directory servers.
type LDAPService struct{ c *Client }

// List lists a tenant's LDAP connections.
//
// API: GET /v1/tenants/{tenantID}/ldap
func (s *LDAPService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[LDAPConnection], error) {
	return call[Page[LDAPConnection]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/ldap",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create creates a connection.
//
// API: POST /v1/tenants/{tenantID}/ldap
func (s *LDAPService) Create(ctx context.Context, tenantID string, params CreateLDAPConnectionParams, opts ...RequestOption) (*LDAPConnection, error) {
	return call[LDAPConnection](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/ldap",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Get returns a connection.
//
// API: GET /v1/tenants/{tenantID}/ldap/{id}
func (s *LDAPService) Get(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*LDAPConnection, error) {
	return call[LDAPConnection](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// Update changes a connection. Fields left nil are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/ldap/{id}
func (s *LDAPService) Update(ctx context.Context, tenantID string, connectionID string, params UpdateLDAPConnectionParams, opts ...RequestOption) (*LDAPConnection, error) {
	return call[LDAPConnection](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// Delete deletes a connection.
//
// API: DELETE /v1/tenants/{tenantID}/ldap/{id}
func (s *LDAPService) Delete(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
	}, nil, opts)
}

// Test binds to the directory with the connection's settings. A failed bind is returned as an error.
//
// API: POST /v1/tenants/{tenantID}/ldap/{id}/test
func (s *LDAPService) Test(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*LDAPTestResult, error) {
	return call[LDAPTestResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}/test",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// ListGroupMappings lists how directory groups map to tenant groups.
//
// API: GET /v1/tenants/{tenantID}/ldap/{id}/group-mappings
func (s *LDAPService) ListGroupMappings(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*Page[GroupMapping], error) {
	return call[Page[GroupMapping]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}/group-mappings",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// CreateGroupMapping maps a directory group to a tenant group.
//
// API: POST /v1/tenants/{tenantID}/ldap/{id}/group-mappings
func (s *LDAPService) CreateGroupMapping(ctx context.Context, tenantID string, connectionID string, params CreateGroupMappingParams, opts ...RequestOption) (*GroupMapping, error) {
	return call[GroupMapping](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}/group-mappings",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// DeleteGroupMapping removes a group mapping.
//
// API: DELETE /v1/tenants/{tenantID}/ldap/{id}/group-mappings/{mappingID}
func (s *LDAPService) DeleteGroupMapping(ctx context.Context, tenantID string, connectionID string, mappingID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/ldap/{connectionID}/group-mappings/{mappingID}",
		pathArgs: []string{tenantID, connectionID, mappingID},
	}, nil, opts)
}

// LDAPConnection connects a tenant to an LDAP or Active Directory server. The bind password is never
// returned.
type LDAPConnection struct {
	BaseDN         string     `json:"base_dn,omitempty"`
	BindDN         string     `json:"bind_dn,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitzero"`
	EmailAttribute string     `json:"email_attribute,omitempty"`
	ID             string     `json:"id,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	Name           string     `json:"name,omitempty"`
	NameAttribute  string     `json:"name_attribute,omitempty"`
	ServerURL      string     `json:"server_url,omitempty"`
	SkipTLSVerify  bool       `json:"skip_tls_verify"`
	StartTLS       bool       `json:"start_tls"`
	Status         string     `json:"status,omitempty"`
	TenantID       string     `json:"tenant_id,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at,omitzero"`
	UserFilter     string     `json:"user_filter,omitempty"`
}

// CreateLDAPConnectionParams is the request body of LDAPService.Create.
type CreateLDAPConnectionParams struct {
	BaseDN         string `json:"base_dn,omitempty"`
	BindDN         string `json:"bind_dn,omitempty"`
	BindPassword   string `json:"bind_password,omitempty"`
	EmailAttribute string `json:"email_attribute,omitempty"`
	Name           string `json:"name,omitempty"`
	NameAttribute  string `json:"name_attribute,omitempty"`
	ServerURL      string `json:"server_url,omitempty"`
	SkipTLSVerify  bool   `json:"skip_tls_verify"`
	StartTLS       bool   `json:"start_tls"`
	Status         string `json:"status,omitempty"`
	UserFilter     string `json:"user_filter,omitempty"`
}

// UpdateLDAPConnectionParams is the request body of LDAPService.Update.
type UpdateLDAPConnectionParams struct {
	BaseDN         *string `json:"base_dn,omitempty"`
	BindDN         *string `json:"bind_dn,omitempty"`
	BindPassword   *string `json:"bind_password,omitempty"`
	EmailAttribute *string `json:"email_attribute,omitempty"`
	Name           *string `json:"name,omitempty"`
	NameAttribute  *string `json:"name_attribute,omitempty"`
	ServerURL      *string `json:"server_url,omitempty"`
	SkipTLSVerify  *bool   `json:"skip_tls_verify,omitempty"`
	StartTLS       *bool   `json:"start_tls,omitempty"`
	Status         *string `json:"status,omitempty"`
	UserFilter     *string `json:"user_filter,omitempty"`
}

// LDAPTestResult is the result of LDAPService.Test.
type LDAPTestResult struct {
	Ok bool `json:"ok"`
}
