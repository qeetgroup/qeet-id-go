package qeetid

import (
	"context"
	"net/http"
	"time"
)

// SAMLServiceProvidersService manages applications that sign people in with Qeet ID over SAML (Qeet ID
// as the identity provider).
type SAMLServiceProvidersService struct{ c *Client }

// List lists a tenant's SAML service providers.
//
// API: GET /v1/tenants/{tenantID}/saml-providers
func (s *SAMLServiceProvidersService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[SAMLServiceProvider], error) {
	return call[Page[SAMLServiceProvider]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/saml-providers",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create registers a service provider.
//
// API: POST /v1/tenants/{tenantID}/saml-providers
func (s *SAMLServiceProvidersService) Create(ctx context.Context, tenantID string, params SAMLServiceProviderParams, opts ...RequestOption) (*SAMLServiceProvider, error) {
	return call[SAMLServiceProvider](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml-providers",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Get returns a service provider.
//
// API: GET /v1/tenants/{tenantID}/saml-providers/{id}
func (s *SAMLServiceProvidersService) Get(ctx context.Context, tenantID, providerID string, opts ...RequestOption) (*SAMLServiceProvider, error) {
	return call[SAMLServiceProvider](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}",
		pathArgs: []string{tenantID, providerID},
	}, opts)
}

// Update changes a service provider. Empty fields are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/saml-providers/{id}
func (s *SAMLServiceProvidersService) Update(ctx context.Context, tenantID, providerID string, params SAMLServiceProviderParams, opts ...RequestOption) (*SAMLServiceProvider, error) {
	return call[SAMLServiceProvider](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}",
		pathArgs: []string{tenantID, providerID},
		body:     params,
	}, opts)
}

// Delete deletes a service provider.
//
// API: DELETE /v1/tenants/{tenantID}/saml-providers/{id}
func (s *SAMLServiceProvidersService) Delete(ctx context.Context, tenantID, providerID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}",
		pathArgs: []string{tenantID, providerID},
	}, nil, opts)
}

// ListAssignments lists the people and groups assigned to a service provider.
//
// API: GET /v1/tenants/{tenantID}/saml-providers/{id}/assignments
func (s *SAMLServiceProvidersService) ListAssignments(ctx context.Context, tenantID, providerID string, opts ...RequestOption) (*SAMLServiceProviderAssignments, error) {
	return call[SAMLServiceProviderAssignments](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}/assignments",
		pathArgs: []string{tenantID, providerID},
	}, opts)
}

// AssignUser assigns a person to a service provider.
//
// API: POST /v1/tenants/{tenantID}/saml-providers/{id}/assignments/users/{userID}
func (s *SAMLServiceProvidersService) AssignUser(ctx context.Context, tenantID, providerID, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}/assignments/users/{userID}",
		pathArgs: []string{tenantID, providerID, userID},
	}, nil, opts)
}

// UnassignUser unassigns a person from a service provider.
//
// API: DELETE /v1/tenants/{tenantID}/saml-providers/{id}/assignments/users/{userID}
func (s *SAMLServiceProvidersService) UnassignUser(ctx context.Context, tenantID, providerID, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}/assignments/users/{userID}",
		pathArgs: []string{tenantID, providerID, userID},
	}, nil, opts)
}

// AssignGroup assigns a group to a service provider.
//
// API: POST /v1/tenants/{tenantID}/saml-providers/{id}/assignments/groups/{groupID}
func (s *SAMLServiceProvidersService) AssignGroup(ctx context.Context, tenantID, providerID, groupID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}/assignments/groups/{groupID}",
		pathArgs: []string{tenantID, providerID, groupID},
	}, nil, opts)
}

// UnassignGroup unassigns a group from a service provider.
//
// API: DELETE /v1/tenants/{tenantID}/saml-providers/{id}/assignments/groups/{groupID}
func (s *SAMLServiceProvidersService) UnassignGroup(ctx context.Context, tenantID, providerID, groupID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/saml-providers/{providerID}/assignments/groups/{groupID}",
		pathArgs: []string{tenantID, providerID, groupID},
	}, nil, opts)
}

// SAMLServiceProvider is an application that signs people in with Qeet ID over SAML, with Qeet ID as
// the identity provider.
type SAMLServiceProvider struct {
	ID       string `json:"id,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Name     string `json:"name,omitempty"`
	EntityID string `json:"entity_id,omitempty"`
	ACSURL   string `json:"acs_url,omitempty"`
	// urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress,
	// urn:oasis:names:tc:SAML:2.0:nameid-format:persistent or
	// urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified
	NameIDFormat string `json:"name_id_format,omitempty"`
	// Where the NameID comes from — email, user_id, persistent_id or metadata.<key>
	NameIDAttribute string `json:"name_id_attribute,omitempty"`
	// The app's own certificate (PEM), kept for require_signed_requests, encrypt_assertions or both.
	Certificate string `json:"certificate,omitempty"`
	// The app's sign-in requests must be signed with its certificate; one that isn't is refused with 403
	// saml.request_signature_invalid. Sign-in the IdP starts has no request and is unaffected.
	RequireSignedRequests bool `json:"require_signed_requests"`
	// Each assertion, signed as always, is encrypted to the app's certificate (AES-256-GCM, its key by
	// RSA-OAEP).
	EncryptAssertions bool `json:"encrypt_assertions"`
	// What the certificate says about itself; null without one.
	CertificateDetails *SAMLServiceProviderCertificateDetails `json:"certificate_details,omitempty"`
	Status             string                                 `json:"status,omitempty"`
	// The attributes each assertion carries beside the NameID, in this order; one the person has no value
	// for is left out
	Attributes []SAMLAttributeMapping `json:"attributes,omitempty"`
	// Only the people assigned to the SP, themselves or through a group, may sign in to it
	AssignmentRequired bool       `json:"assignment_required"`
	CreatedAt          time.Time  `json:"created_at,omitzero"`
	UpdatedAt          time.Time  `json:"updated_at,omitzero"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
}

// SAMLServiceProviderCertificateDetails is the type of SAMLServiceProvider.CertificateDetails. What
// the certificate says about itself; null without one.
type SAMLServiceProviderCertificateDetails struct {
	Subject           string    `json:"subject,omitempty"`
	FingerprintSHA256 string    `json:"fingerprint_sha256,omitempty"`
	NotBefore         time.Time `json:"not_before,omitzero"`
	NotAfter          time.Time `json:"not_after,omitzero"`
}

// SAMLAttributeMapping is one element of SAMLServiceProvider.Attributes.
type SAMLAttributeMapping struct {
	// The attribute's SAML Name, unique within the SP
	Name string `json:"name,omitempty"`
	// email; user_id (the account ID, the same at every app); persistent_id (an opaque value made for
	// this person at this SP, which never changes); display_name; given_name and family_name (the
	// user-metadata key of that name, else the person's directory profile first_name or last_name);
	// groups (the names of their groups in the organisation, one value each); roles (the names of their
	// roles there, directly or through a group, one value each); or metadata.<key> (that user-metadata
	// key: a string, number or boolean, or an array of them as one value each)
	Source string `json:"source,omitempty"`
	// The attribute's NameFormat; basic when absent
	NameFormat SAMLAttributeMappingNameFormat `json:"name_format,omitempty"`
}

// SAMLAttributeMappingNameFormat is a string enumeration used by SAMLAttributeMapping.NameFormat. The
// attribute's NameFormat; basic when absent
type SAMLAttributeMappingNameFormat string

// Values of SAMLAttributeMappingNameFormat. The server may add values; handle an unknown one gracefully.
const (
	SAMLAttributeMappingNameFormatBasic       SAMLAttributeMappingNameFormat = "basic"
	SAMLAttributeMappingNameFormatURI         SAMLAttributeMappingNameFormat = "uri"
	SAMLAttributeMappingNameFormatUnspecified SAMLAttributeMappingNameFormat = "unspecified"
)

// SAMLServiceProviderParams is the request body of SAMLServiceProvidersService.Create.
type SAMLServiceProviderParams struct {
	Name     string `json:"name,omitempty"`
	EntityID string `json:"entity_id,omitempty"`
	ACSURL   string `json:"acs_url,omitempty"`
	// urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress (the default),
	// urn:oasis:names:tc:SAML:2.0:nameid-format:persistent or
	// urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified
	NameIDFormat string `json:"name_id_format,omitempty"`
	// Where the NameID comes from — email, user_id, persistent_id or metadata.<key>; by default
	// persistent_id for a persistent NameID and email otherwise
	NameIDAttribute string `json:"name_id_attribute,omitempty"`
	// The app's own certificate (PEM). It must be kept for require_signed_requests, encrypt_assertions or
	// both (422 saml.certificate_unused otherwise), and either needs one (422 saml.certificate_required);
	// encryption needs an RSA key (422 saml.encryption_key_unsupported).
	Certificate string `json:"certificate,omitempty"`
	// Refuse sign-in requests not signed with the certificate; off by default
	RequireSignedRequests *bool `json:"require_signed_requests,omitempty"`
	// Encrypt each assertion to the certificate; off by default
	EncryptAssertions *bool `json:"encrypt_assertions,omitempty"`
	// The attributes to send; absent sends the default email (from email) and name (from display_name),
	// an empty list none
	Attributes []SAMLAttributeMapping `json:"attributes,omitempty"`
	// Sign in only the people assigned to the SP; off by default
	AssignmentRequired *bool `json:"assignment_required,omitempty"`
}

// SAMLServiceProviderAssignments is the result of SAMLServiceProvidersService.ListAssignments.
type SAMLServiceProviderAssignments struct {
	Users  []SAMLServiceProviderAssignmentsUser  `json:"users,omitempty"`
	Groups []SAMLServiceProviderAssignmentsGroup `json:"groups,omitempty"`
}

// SAMLServiceProviderAssignmentsUser is one element of SAMLServiceProviderAssignments.Users.
type SAMLServiceProviderAssignmentsUser struct {
	UserID      string    `json:"user_id,omitempty"`
	Email       string    `json:"email,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	AssignedAt  time.Time `json:"assigned_at,omitzero"`
	AssignedBy  string    `json:"assigned_by,omitempty"`
}

// SAMLServiceProviderAssignmentsGroup is one element of SAMLServiceProviderAssignments.Groups.
type SAMLServiceProviderAssignmentsGroup struct {
	GroupID    string    `json:"group_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	AssignedAt time.Time `json:"assigned_at,omitzero"`
	AssignedBy string    `json:"assigned_by,omitempty"`
}
