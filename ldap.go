package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// LDAPConnection is a tenant's LDAP/AD bind config. BindPassword is
// write-only — never returned.
type LDAPConnection struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	Name           string `json:"name"`
	ServerURL      string `json:"server_url"`
	StartTLS       bool   `json:"start_tls"`
	SkipTLSVerify  bool   `json:"skip_tls_verify"`
	BindDN         string `json:"bind_dn"`
	BaseDN         string `json:"base_dn"`
	UserFilter     string `json:"user_filter"`
	EmailAttribute string `json:"email_attribute"`
	NameAttribute  string `json:"name_attribute"`
	Status         string `json:"status"` // draft | active | disabled
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	LastLoginAt    string `json:"last_login_at,omitempty"`
}

// CreateLDAPConnectionInput — Name, ServerURL, BindDN, BindPassword, BaseDN
// are required. ServerURL must start ldap:// or ldaps://. Unset optional
// fields get server-side defaults: UserFilter "(uid=%s)", EmailAttribute
// "mail", NameAttribute "cn", Status "draft".
type CreateLDAPConnectionInput struct {
	Name           string `json:"name"`
	ServerURL      string `json:"server_url"`
	StartTLS       bool   `json:"start_tls,omitempty"`
	SkipTLSVerify  bool   `json:"skip_tls_verify,omitempty"`
	BindDN         string `json:"bind_dn"`
	BindPassword   string `json:"bind_password"`
	BaseDN         string `json:"base_dn"`
	UserFilter     string `json:"user_filter,omitempty"`
	EmailAttribute string `json:"email_attribute,omitempty"`
	NameAttribute  string `json:"name_attribute,omitempty"`
	Status         string `json:"status,omitempty"`
}

type UpdateLDAPConnectionInput struct {
	Name           *string `json:"name,omitempty"`
	ServerURL      *string `json:"server_url,omitempty"`
	StartTLS       *bool   `json:"start_tls,omitempty"`
	SkipTLSVerify  *bool   `json:"skip_tls_verify,omitempty"`
	BindDN         *string `json:"bind_dn,omitempty"`
	BindPassword   *string `json:"bind_password,omitempty"`
	BaseDN         *string `json:"base_dn,omitempty"`
	UserFilter     *string `json:"user_filter,omitempty"`
	EmailAttribute *string `json:"email_attribute,omitempty"`
	NameAttribute  *string `json:"name_attribute,omitempty"`
	Status         *string `json:"status,omitempty"` // draft | active | disabled
}

// LDAPTestResult is the outcome of a bind test — {"ok":true} on success; a
// failure surfaces as an error instead (dial failed, bind failed, etc.).
type LDAPTestResult struct {
	OK bool `json:"ok"`
}

// LDAPService manages LDAP/AD connections.
type LDAPService struct{ t *transport.Transport }

func (s *LDAPService) Create(ctx context.Context, tenantID string, in CreateLDAPConnectionInput) (*LDAPConnection, error) {
	var out LDAPConnection
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ldap", in, &out)
	return &out, err
}

func (s *LDAPService) Get(ctx context.Context, tenantID, id string) (*LDAPConnection, error) {
	var out LDAPConnection
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ldap/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *LDAPService) Update(ctx context.Context, tenantID, id string, in UpdateLDAPConnectionInput) (*LDAPConnection, error) {
	var out LDAPConnection
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ldap/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *LDAPService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ldap/"+url.PathEscape(id), nil)
}

func (s *LDAPService) List(ctx context.Context, tenantID string) ([]LDAPConnection, error) {
	var env marshal.Envelope[LDAPConnection]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ldap", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *LDAPService) Test(ctx context.Context, tenantID, id string) (*LDAPTestResult, error) {
	var out LDAPTestResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ldap/"+url.PathEscape(id)+"/test", struct{}{}, &out)
	return &out, err
}

// Authenticate is a public, unversioned passthrough (no /v1 prefix, no auth)
// for legacy apps doing direct LDAP-bind authentication against a configured
// connection.
func (s *LDAPService) Authenticate(ctx context.Context, connectionID, username, password string) error {
	body := map[string]string{"username": username, "password": password}
	return s.t.Post(ctx, "/ldap/"+url.PathEscape(connectionID)+"/authenticate", body, nil)
}
