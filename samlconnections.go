package qeetid

import (
	"context"
	"net/http"
	"time"
)

// SAMLConnectionsService manages SAML connections to a tenant's identity providers (Qeet ID as the
// service provider).
type SAMLConnectionsService struct{ c *Client }

// List lists a tenant's SAML connections.
//
// API: GET /v1/tenants/{tenantID}/saml
func (s *SAMLConnectionsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[SAMLConnection], error) {
	return call[Page[SAMLConnection]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/saml",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create creates a connection from its settings.
//
// API: POST /v1/tenants/{tenantID}/saml
func (s *SAMLConnectionsService) Create(ctx context.Context, tenantID string, params SAMLConnectionParams, opts ...RequestOption) (*SAMLConnection, error) {
	return call[SAMLConnection](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Import creates a connection from the identity provider's metadata.
//
// API: POST /v1/tenants/{tenantID}/saml/import
func (s *SAMLConnectionsService) Import(ctx context.Context, tenantID string, params ImportSAMLConnectionParams, opts ...RequestOption) (*SAMLConnection, error) {
	return call[SAMLConnection](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml/import",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Get returns a connection.
//
// API: GET /v1/tenants/{tenantID}/saml/{id}
func (s *SAMLConnectionsService) Get(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*SAMLConnection, error) {
	return call[SAMLConnection](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// Update changes a connection. Empty fields are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/saml/{id}
func (s *SAMLConnectionsService) Update(ctx context.Context, tenantID string, connectionID string, params SAMLConnectionParams, opts ...RequestOption) (*SAMLConnection, error) {
	return call[SAMLConnection](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// Delete deletes a connection.
//
// API: DELETE /v1/tenants/{tenantID}/saml/{id}
func (s *SAMLConnectionsService) Delete(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}",
		pathArgs: []string{tenantID, connectionID},
	}, nil, opts)
}

// ApplyMetadata updates a connection from the identity provider's metadata.
//
// API: POST /v1/tenants/{tenantID}/saml/{id}/metadata
func (s *SAMLConnectionsService) ApplyMetadata(ctx context.Context, tenantID string, connectionID string, params ApplySAMLMetadataParams, opts ...RequestOption) (*SAMLConnection, error) {
	return call[SAMLConnection](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/metadata",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// AddCertificate adds an identity provider signing certificate, for a certificate rollover.
//
// API: POST /v1/tenants/{tenantID}/saml/{id}/certificates
func (s *SAMLConnectionsService) AddCertificate(ctx context.Context, tenantID string, connectionID string, params AddSAMLCertificateParams, opts ...RequestOption) (*SAMLCertificate, error) {
	return call[SAMLCertificate](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/certificates",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// RemoveCertificate removes an identity provider signing certificate.
//
// API: DELETE /v1/tenants/{tenantID}/saml/{id}/certificates/{certificateID}
func (s *SAMLConnectionsService) RemoveCertificate(ctx context.Context, tenantID string, connectionID string, certificateID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/certificates/{certificateID}",
		pathArgs: []string{tenantID, connectionID, certificateID},
	}, nil, opts)
}

// Test checks a connection's settings.
//
// API: POST /v1/tenants/{tenantID}/saml/{id}/test
func (s *SAMLConnectionsService) Test(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*SAMLTestResult, error) {
	return call[SAMLTestResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/test",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// ListGroupMappings lists how identity provider groups map to tenant groups.
//
// API: GET /v1/tenants/{tenantID}/saml/{id}/group-mappings
func (s *SAMLConnectionsService) ListGroupMappings(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*Page[GroupMapping], error) {
	return call[Page[GroupMapping]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/group-mappings",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// CreateGroupMapping maps an identity provider group to a tenant group.
//
// API: POST /v1/tenants/{tenantID}/saml/{id}/group-mappings
func (s *SAMLConnectionsService) CreateGroupMapping(ctx context.Context, tenantID string, connectionID string, params CreateGroupMappingParams, opts ...RequestOption) (*GroupMapping, error) {
	return call[GroupMapping](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/group-mappings",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// DeleteGroupMapping removes a group mapping.
//
// API: DELETE /v1/tenants/{tenantID}/saml/{id}/group-mappings/{mappingID}
func (s *SAMLConnectionsService) DeleteGroupMapping(ctx context.Context, tenantID string, connectionID string, mappingID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/saml/{connectionID}/group-mappings/{mappingID}",
		pathArgs: []string{tenantID, connectionID, mappingID},
	}, nil, opts)
}

// SAMLConnection connects a tenant to a SAML identity provider, with Qeet ID as the service provider.
type SAMLConnection struct {
	// The connection's own entity ID, for the IdP.
	SPEntityID string `json:"sp_entity_id,omitempty"`
	// Where the IdP posts assertions.
	ACSURL string `json:"acs_url,omitempty"`
	// The connection's own metadata, for the IdP to import.
	SPMetadataURL string `json:"sp_metadata_url,omitempty"`
	ID            string `json:"id,omitempty"`
	TenantID      string `json:"tenant_id,omitempty"`
	Name          string `json:"name,omitempty"`
	IDPEntityID   string `json:"idp_entity_id,omitempty"`
	IDPSSOURL     string `json:"idp_sso_url,omitempty"`
	// The oldest of the signing certificates the connection trusts.
	IDPCertificate string `json:"idp_certificate,omitempty"`
	// Every IdP signing certificate the connection trusts, oldest first. An assertion signed with any of
	// them is accepted, so an IdP can roll its key over without a sign-in outage.
	Certificates []SAMLCertificate `json:"certificates,omitempty"`
	// The IdP metadata URL the connection refreshes its SSO URL and certificates from; null when it
	// doesn't refresh.
	MetadataURL string `json:"metadata_url,omitempty"`
	// When the connection last took its settings from the IdP's metadata.
	MetadataRefreshedAt *time.Time `json:"metadata_refreshed_at,omitempty"`
	// Why the last refresh from metadata_url failed; null once one succeeds.
	MetadataError  string `json:"metadata_error,omitempty"`
	EmailAttribute string `json:"email_attribute,omitempty"`
	// The attribute carrying the person's display name. When it is empty or the assertion lacks it, the
	// common display-name attributes are read, else the given and family names joined. Each sign-in sets
	// the account's display name to what the assertion carries.
	NameAttribute string `json:"name_attribute,omitempty"`
	// The attribute carrying the person's groups, matched against the connection's group mappings and the
	// external ids of groups SCIM provisioned. Empty (the default) reads groups,
	// http://schemas.microsoft.com/ws/2008/06/identity/claims/groups and memberOf. When Microsoft Entra
	// ID leaves the groups out for being too many, memberships stay as they are.
	GroupsAttribute string `json:"groups_attribute,omitempty"`
	Status          string `json:"status,omitempty"`
	// Accept unsolicited assertions that carry no InResponseTo, i.e. IdP-initiated SSO. Off by default.
	// Enabling it relaxes only the AuthnRequest correlation — signature, issuer, audience, destination
	// and time-condition validation are unchanged, and response and assertion IDs are still consumed for
	// replay protection.
	AllowIDPInitiated bool       `json:"allow_idp_initiated"`
	CreatedAt         time.Time  `json:"created_at,omitzero"`
	UpdatedAt         time.Time  `json:"updated_at,omitzero"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
}

// SAMLCertificate is one element of SAMLConnection.Certificates.
type SAMLCertificate struct {
	ID string `json:"id,omitempty"`
	// PEM or base64 DER X.509 certificate
	Certificate string `json:"certificate,omitempty"`
	// The certificate's subject; empty if it doesn't parse.
	Subject string `json:"subject,omitempty"`
	// SHA-256 fingerprint as colon-separated hex; empty if it doesn't parse.
	FingerprintSHA256 string     `json:"fingerprint_sha256,omitempty"`
	NotBefore         *time.Time `json:"not_before,omitempty"`
	// When the certificate expires; an assertion it signed is refused after then.
	NotAfter *time.Time `json:"not_after,omitempty"`
	// Whether an administrator added the certificate or the IdP's metadata listed it.
	Source    SAMLCertificateSource `json:"source,omitempty"`
	CreatedAt time.Time             `json:"created_at,omitzero"`
}

// SAMLCertificateSource is a string enumeration used by SAMLCertificate.Source. Whether an
// administrator added the certificate or the IdP's metadata listed it.
type SAMLCertificateSource string

// Values of SAMLCertificateSource. The server may add values; handle an unknown one gracefully.
const (
	SAMLCertificateSourceManual   SAMLCertificateSource = "manual"
	SAMLCertificateSourceMetadata SAMLCertificateSource = "metadata"
)

// SAMLConnectionParams is the name, with all of the IdP's entity ID, SSO URL and certificate, or none
// of them for a draft.
type SAMLConnectionParams struct {
	Name        string `json:"name,omitempty"`
	IDPEntityID string `json:"idp_entity_id,omitempty"`
	IDPSSOURL   string `json:"idp_sso_url,omitempty"`
	// PEM or base64 DER X.509 signing certificate. On update it becomes the only certificate the
	// connection trusts; add and remove certificates one at a time with the certificates endpoints.
	IDPCertificate string `json:"idp_certificate,omitempty"`
	EmailAttribute string `json:"email_attribute,omitempty"`
	NameAttribute  string `json:"name_attribute,omitempty"`
	// The attribute carrying the person's groups; empty reads the common ones (see SAMLConnection). Over
	// 256 characters is refused with 422 validation_failed.
	GroupsAttribute string `json:"groups_attribute,omitempty"`
	Status          string `json:"status,omitempty"`
}

// ImportSAMLConnectionParams is the request body of SAMLConnectionsService.Import.
type ImportSAMLConnectionParams struct {
	Name string `json:"name,omitempty"`
	// Where to fetch the IdP's metadata from, and refresh it from later. Give this or metadata_xml.
	MetadataURL string `json:"metadata_url,omitempty"`
	// The IdP's SAML 2.0 metadata document. Give this or metadata_url.
	MetadataXML string `json:"metadata_xml,omitempty"`
	// The entity id of the IdP to use, when the metadata describes several.
	IDPEntityID    string `json:"idp_entity_id,omitempty"`
	EmailAttribute string `json:"email_attribute,omitempty"`
	NameAttribute  string `json:"name_attribute,omitempty"`
	// The assertion attribute carrying a person's groups; empty reads the common ones.
	GroupsAttribute string `json:"groups_attribute,omitempty"`
	Status          string `json:"status,omitempty"`
}

// ApplySAMLMetadataParams is the request body of SAMLConnectionsService.ApplyMetadata.
type ApplySAMLMetadataParams struct {
	MetadataURL string `json:"metadata_url,omitempty"`
	MetadataXML string `json:"metadata_xml,omitempty"`
}

// AddSAMLCertificateParams is the request body of SAMLConnectionsService.AddCertificate.
type AddSAMLCertificateParams struct {
	// PEM or base64 DER X.509 signing certificate
	Certificate string `json:"certificate,omitempty"`
}

// SAMLTestResult is the result of SAMLConnectionsService.Test.
type SAMLTestResult struct {
	Checks []SAMLTestCheck `json:"checks,omitempty"`
	Ok     bool            `json:"ok"`
}

// SAMLTestCheck is one element of SAMLTestResult.Checks.
type SAMLTestCheck struct {
	Detail string `json:"detail,omitempty"`
	Name   string `json:"name,omitempty"`
	Ok     bool   `json:"ok"`
}

// GroupMapping maps a group of an identity provider or directory to a tenant group.
type GroupMapping struct {
	ID string `json:"id,omitempty"`
	// The group as the identity provider asserts it (a name, an object id or an LDAP DN), compared
	// ignoring case.
	IDPGroup  string    `json:"idp_group,omitempty"`
	GroupID   string    `json:"group_id,omitempty"`
	GroupName string    `json:"group_name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}

// CreateGroupMappingParams is the request body of SAMLConnectionsService.CreateGroupMapping.
type CreateGroupMappingParams struct {
	IDPGroup string `json:"idp_group,omitempty"`
	GroupID  string `json:"group_id,omitempty"`
}
