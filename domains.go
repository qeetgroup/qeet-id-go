package qeetid

import (
	"context"
	"net/http"
	"time"
)

// DomainsService proves control of email domains and manages the organization's custom sign-in domain and
// break-glass members.
type DomainsService struct{ c *Client }

// List lists an organization's domains.
//
// API: GET /v1/organizations/{organizationID}/domains
func (s *DomainsService) List(ctx context.Context, organizationID string, opts ...RequestOption) (*Page[Domain], error) {
	return call[Page[Domain]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/domains",
		pathArgs: []string{organizationID},
	}, opts)
}

// Create claims a domain. Publish the DNS record the result names, then call Verify.
//
// API: POST /v1/organizations/{organizationID}/domains
func (s *DomainsService) Create(ctx context.Context, organizationID string, params CreateDomainParams, opts ...RequestOption) (*Domain, error) {
	return call[Domain](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/domains",
		pathArgs: []string{organizationID},
		body:     params,
	}, opts)
}

// Verify checks a claimed domain's DNS record and marks the domain verified when it matches.
//
// API: POST /v1/organizations/{organizationID}/domains/{id}/verify
func (s *DomainsService) Verify(ctx context.Context, organizationID, domainID string, opts ...RequestOption) (*Domain, error) {
	return call[Domain](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/domains/{domainID}/verify",
		pathArgs: []string{organizationID, domainID},
	}, opts)
}

// Update changes a verified domain's single sign-on and just-in-time provisioning settings.
//
// API: PATCH /v1/organizations/{organizationID}/domains/{id}
func (s *DomainsService) Update(ctx context.Context, organizationID, domainID string, params UpdateDomainParams, opts ...RequestOption) (*Domain, error) {
	return call[Domain](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/organizations/{organizationID}/domains/{domainID}",
		pathArgs: []string{organizationID, domainID},
		body:     params,
	}, opts)
}

// Delete removes a domain.
//
// API: DELETE /v1/organizations/{organizationID}/domains/{id}
func (s *DomainsService) Delete(ctx context.Context, organizationID, domainID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}/domains/{domainID}",
		pathArgs: []string{organizationID, domainID},
	}, nil, opts)
}

// LoginDomain returns the state of the organization's custom sign-in domain.
//
// API: GET /v1/organizations/{organizationID}/login-domain
func (s *DomainsService) LoginDomain(ctx context.Context, organizationID string, opts ...RequestOption) (*LoginDomainStatus, error) {
	return call[LoginDomainStatus](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/login-domain",
		pathArgs: []string{organizationID},
	}, opts)
}

// CheckLoginDomain checks the organization's custom sign-in domain now and returns its state.
//
// API: POST /v1/organizations/{organizationID}/login-domain/check
func (s *DomainsService) CheckLoginDomain(ctx context.Context, organizationID string, opts ...RequestOption) (*LoginDomainStatus, error) {
	return call[LoginDomainStatus](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/login-domain/check",
		pathArgs: []string{organizationID},
	}, opts)
}

// ListBreakGlassMembers lists the members who may sign in with a password even when single sign-on is
// required.
//
// API: GET /v1/organizations/{organizationID}/break-glass-members
func (s *DomainsService) ListBreakGlassMembers(ctx context.Context, organizationID string, opts ...RequestOption) (*Page[BreakGlassMember], error) {
	return call[Page[BreakGlassMember]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/break-glass-members",
		pathArgs: []string{organizationID},
	}, opts)
}

// AddBreakGlassMember makes a member a break-glass member.
//
// API: PUT /v1/organizations/{organizationID}/break-glass-members/{userID}
func (s *DomainsService) AddBreakGlassMember(ctx context.Context, organizationID, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPut,
		path:     "/v1/organizations/{organizationID}/break-glass-members/{userID}",
		pathArgs: []string{organizationID, userID},
	}, nil, opts)
}

// RemoveBreakGlassMember removes a break-glass member.
//
// API: DELETE /v1/organizations/{organizationID}/break-glass-members/{userID}
func (s *DomainsService) RemoveBreakGlassMember(ctx context.Context, organizationID, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}/break-glass-members/{userID}",
		pathArgs: []string{organizationID, userID},
	}, nil, opts)
}

// Domain is an email domain an organization claims. Once verified it can require single sign-on and provision
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
	// Addresses on the domain sign in only through the organization's active SAML, LDAP or OpenID Connect
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
