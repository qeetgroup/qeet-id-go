package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// SCIMConfig describes a tenant's SCIM provisioning setup. TokenSet/
// TokenPrefix let the console show "a token exists, starts with X" without
// ever re-exposing the secret.
type SCIMConfig struct {
	TokenSet         bool   `json:"token_set"`
	TokenPrefix      string `json:"token_prefix,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	LastUsedAt       string `json:"last_used_at,omitempty"`
	ProvisionedCount int    `json:"provisioned_count"`
}

// RotateSCIMTokenResult carries the new plaintext bearer token, shown once.
type RotateSCIMTokenResult struct {
	Token  string     `json:"token"`
	Config SCIMConfig `json:"config"`
}

// ProvisionedUser is a user provisioned into Qeet ID by the tenant's IdP via
// SCIM — distinct from the full User type; this is the admin-facing summary.
type ProvisionedUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Status      string `json:"status"`
	ExternalID  string `json:"external_id,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// SCIMService is the tenant-admin config surface for SCIM provisioning —
// not the /scim/v2/* protocol endpoints themselves, which the customer's
// IdP (Okta, Azure AD, etc.) calls directly and this SDK has no reason to
// wrap.
type SCIMService struct{ t *transport.Transport }

func (s *SCIMService) GetConfig(ctx context.Context, tenantID string) (*SCIMConfig, error) {
	var out SCIMConfig
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/scim", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *SCIMService) RotateToken(ctx context.Context, tenantID string) (*RotateSCIMTokenResult, error) {
	var out RotateSCIMTokenResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/scim/token", struct{}{}, &out)
	return &out, err
}

func (s *SCIMService) RevokeToken(ctx context.Context, tenantID string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/scim/token", nil)
}

func (s *SCIMService) ListProvisionedUsers(ctx context.Context, tenantID string) ([]ProvisionedUser, error) {
	var env marshal.Envelope[ProvisionedUser]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/scim/users", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
