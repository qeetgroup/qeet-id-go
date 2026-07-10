package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// SecurityPolicy is a tenant's combined security policy: IP allow/deny
// lists, password rules, session lifetime, and MFA enforcement in one
// record. It's a distinct, older resource from AuthPolicy/IPRules — those
// are more granular views over overlapping concerns; this is the original
// combined record, still live at its own path.
type SecurityPolicy struct {
	TenantID           string         `json:"tenant_id"`
	IPAllowlist        []string       `json:"ip_allowlist,omitempty"`
	IPDenylist         []string       `json:"ip_denylist,omitempty"`
	PasswordMinLength  int            `json:"password_min_length,omitempty"`
	PasswordComplexity string         `json:"password_complexity,omitempty"`
	SessionMaxAge      int64          `json:"session_max_age,omitempty"` // nanoseconds
	MFAEnforcement     string         `json:"mfa_enforcement,omitempty"`
	Settings           map[string]any `json:"settings,omitempty"`
}

// PolicyService manages the combined per-tenant security policy record.
type PolicyService struct{ t *transport.Transport }

func (s *PolicyService) Get(ctx context.Context, tenantID string) (*SecurityPolicy, error) {
	var out SecurityPolicy
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/policy", transport.RequestOptions{}, &out)
	return &out, err
}

// Put replaces the entire policy record — this is a full overwrite, not a
// partial patch.
func (s *PolicyService) Put(ctx context.Context, tenantID string, in SecurityPolicy) (*SecurityPolicy, error) {
	var out SecurityPolicy
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/policy", in, &out)
	return &out, err
}
