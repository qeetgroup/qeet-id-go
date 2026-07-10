package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type AuthPolicySettings struct {
	TenantID                 string   `json:"tenant_id"`
	PasswordMinLength        int      `json:"password_min_length,omitempty"`
	PasswordRequireUppercase bool     `json:"password_require_uppercase,omitempty"`
	PasswordRequireNumbers   bool     `json:"password_require_numbers,omitempty"`
	PasswordRequireSymbols   bool     `json:"password_require_symbols,omitempty"`
	AllowedLoginMethods      []string `json:"allowed_login_methods,omitempty"`
	MfaRequired              bool     `json:"mfa_required,omitempty"`
	SessionDurationSeconds   int      `json:"session_duration_seconds,omitempty"`
	UpdatedAt                string   `json:"updated_at,omitempty"`
}

type UpdateAuthPolicyInput struct {
	PasswordMinLength        *int     `json:"password_min_length,omitempty"`
	PasswordRequireUppercase *bool    `json:"password_require_uppercase,omitempty"`
	PasswordRequireNumbers   *bool    `json:"password_require_numbers,omitempty"`
	PasswordRequireSymbols   *bool    `json:"password_require_symbols,omitempty"`
	AllowedLoginMethods      []string `json:"allowed_login_methods,omitempty"`
	MfaRequired              *bool    `json:"mfa_required,omitempty"`
	SessionDurationSeconds   *int     `json:"session_duration_seconds,omitempty"`
}

// AuthPolicyService manages tenant-wide login policy (password rules, MFA
// requirement, session duration).
type AuthPolicyService struct{ t *transport.Transport }

func (s *AuthPolicyService) Get(ctx context.Context, tenantID string) (*AuthPolicySettings, error) {
	var out AuthPolicySettings
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/auth-policy", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *AuthPolicyService) Update(ctx context.Context, tenantID string, in UpdateAuthPolicyInput) (*AuthPolicySettings, error) {
	var out AuthPolicySettings
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/auth-policy", in, &out)
	return &out, err
}
