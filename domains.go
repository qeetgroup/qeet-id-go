package qeetid

import (
	"context"
	"net/http"
	"time"
)

// DomainsService proves control of email domains and manages the tenant's custom sign-in domain and
// break-glass members.
type DomainsService struct{ c *Client }

// List lists a tenant's domains.
//
// API: GET /v1/tenants/{tenantID}/domains
func (s *DomainsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Domain], error) {
	return call[Page[Domain]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/domains",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create claims a domain. Publish the DNS record the result names, then call Verify.
//
// API: POST /v1/tenants/{tenantID}/domains
func (s *DomainsService) Create(ctx context.Context, tenantID string, params CreateDomainParams, opts ...RequestOption) (*Domain, error) {
	return call[Domain](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/domains",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Verify checks a claimed domain's DNS record and marks the domain verified when it matches.
//
// API: POST /v1/tenants/{tenantID}/domains/{id}/verify
func (s *DomainsService) Verify(ctx context.Context, tenantID, domainID string, opts ...RequestOption) (*Domain, error) {
	return call[Domain](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/domains/{domainID}/verify",
		pathArgs: []string{tenantID, domainID},
	}, opts)
}

// Update changes a verified domain's single sign-on and just-in-time provisioning settings.
//
// API: PATCH /v1/tenants/{tenantID}/domains/{id}
func (s *DomainsService) Update(ctx context.Context, tenantID, domainID string, params UpdateDomainParams, opts ...RequestOption) (*Domain, error) {
	return call[Domain](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/domains/{domainID}",
		pathArgs: []string{tenantID, domainID},
		body:     params,
	}, opts)
}

// Delete removes a domain.
//
// API: DELETE /v1/tenants/{tenantID}/domains/{id}
func (s *DomainsService) Delete(ctx context.Context, tenantID, domainID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/domains/{domainID}",
		pathArgs: []string{tenantID, domainID},
	}, nil, opts)
}

// LoginDomain returns the state of the tenant's custom sign-in domain.
//
// API: GET /v1/tenants/{tenantID}/login-domain
func (s *DomainsService) LoginDomain(ctx context.Context, tenantID string, opts ...RequestOption) (*LoginDomainStatus, error) {
	return call[LoginDomainStatus](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/login-domain",
		pathArgs: []string{tenantID},
	}, opts)
}

// CheckLoginDomain checks the tenant's custom sign-in domain now and returns its state.
//
// API: POST /v1/tenants/{tenantID}/login-domain/check
func (s *DomainsService) CheckLoginDomain(ctx context.Context, tenantID string, opts ...RequestOption) (*LoginDomainStatus, error) {
	return call[LoginDomainStatus](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/login-domain/check",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListBreakGlassMembers lists the members who may sign in with a password even when single sign-on is
// required.
//
// API: GET /v1/tenants/{tenantID}/break-glass-members
func (s *DomainsService) ListBreakGlassMembers(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[BreakGlassMember], error) {
	return call[Page[BreakGlassMember]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/break-glass-members",
		pathArgs: []string{tenantID},
	}, opts)
}

// AddBreakGlassMember makes a member a break-glass member.
//
// API: PUT /v1/tenants/{tenantID}/break-glass-members/{userID}
func (s *DomainsService) AddBreakGlassMember(ctx context.Context, tenantID, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/break-glass-members/{userID}",
		pathArgs: []string{tenantID, userID},
	}, nil, opts)
}

// RemoveBreakGlassMember removes a break-glass member.
//
// API: DELETE /v1/tenants/{tenantID}/break-glass-members/{userID}
func (s *DomainsService) RemoveBreakGlassMember(ctx context.Context, tenantID, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/break-glass-members/{userID}",
		pathArgs: []string{tenantID, userID},
	}, nil, opts)
}

// Domain is an email domain a tenant claims. Once verified it can require single sign-on and provision
// accounts just in time.
type Domain struct {
	CreatedAt      time.Time `json:"created_at,omitzero"`
	DNSRecordName  string    `json:"dns_record_name,omitempty"`
	DNSRecordType  string    `json:"dns_record_type,omitempty"`
	DNSRecordValue string    `json:"dns_record_value,omitempty"`
	Domain         string    `json:"domain,omitempty"`
	ID             string    `json:"id,omitempty"`
	// Single sign-on creates an account at sign-in for an address on the domain that has none.
	JitEnabled    bool      `json:"jit_enabled"`
	LastCheckedAt time.Time `json:"last_checked_at,omitzero"`
	LastError     string    `json:"last_error,omitempty"`
	// Addresses on the domain sign in only through the organisation's active SAML, LDAP or OpenID Connect
	// connection (break-glass members and administrators single sign-on can't sign in excepted).
	SSOEnabled        bool      `json:"sso_enabled"`
	Status            string    `json:"status,omitempty"`
	VerificationToken string    `json:"verification_token,omitempty"`
	VerifiedAt        time.Time `json:"verified_at,omitzero"`
}

// CreateDomainParams is the request body of DomainsService.Create.
type CreateDomainParams struct {
	Domain string `json:"domain,omitempty"`
}

// UpdateDomainParams is the request body of DomainsService.Update.
type UpdateDomainParams struct {
	SSOEnabled *bool `json:"sso_enabled,omitempty"`
	JitEnabled *bool `json:"jit_enabled,omitempty"`
}

// LoginDomainStatus is the result of DomainsService.LoginDomain.
type LoginDomainStatus struct {
	// When this result was found; null without a hostname.
	CheckedAt     *time.Time  `json:"checked_at,omitempty"`
	DNSCheckedAt  time.Time   `json:"dns_checked_at,omitzero"`
	DNSDetail     string      `json:"dns_detail,omitempty"`
	DNSVerified   bool        `json:"dns_verified"`
	EndpointState string      `json:"endpoint_state,omitempty"`
	EndpointURL   string      `json:"endpoint_url,omitempty"`
	LoginDomain   string      `json:"login_domain,omitempty"`
	Records       []DNSRecord `json:"records,omitempty"`
	Status        string      `json:"status,omitempty"`
	TLSDetail     string      `json:"tls_detail,omitempty"`
	TLSState      string      `json:"tls_state,omitempty"`
}

// DNSRecord is one element of LoginDomainStatus.Records.
type DNSRecord struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

// BreakGlassMember is one item of the list DomainsService.ListBreakGlassMembers returns.
type BreakGlassMember struct {
	UserID      string `json:"user_id,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	// The owner who named them, while that account exists.
	CreatedBy string    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}
