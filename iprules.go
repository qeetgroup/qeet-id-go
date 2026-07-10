package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type IPRule struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	CIDR      string `json:"cidr"`
	Label     string `json:"label,omitempty"`
	Action    string `json:"action"` // allow | deny
	CreatedAt string `json:"created_at"`
}

type CreateIPRuleInput struct {
	CIDR   string `json:"cidr"`
	Label  string `json:"label,omitempty"`
	Action string `json:"action"`
}

// IPRuleCheckResult is the outcome of testing an address against a tenant's
// rule set. Enabled reports whether enforcement is on at all — Allowed is
// only meaningful when it is.
type IPRuleCheckResult struct {
	Enabled bool   `json:"enabled"`
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// IPRulesService manages tenant IP allow/deny rules.
type IPRulesService struct{ t *transport.Transport }

func (s *IPRulesService) Create(ctx context.Context, tenantID string, in CreateIPRuleInput) (*IPRule, error) {
	var out IPRule
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ip-rules", in, &out)
	return &out, err
}

func (s *IPRulesService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ip-rules/"+url.PathEscape(id), nil)
}

func (s *IPRulesService) List(ctx context.Context, tenantID string) ([]IPRule, error) {
	var env marshal.Envelope[IPRule]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ip-rules", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// Check tests whether ip would be allowed under the tenant's current rule
// set and enforcement setting — useful to dry-run a rule change before
// enabling enforcement.
func (s *IPRulesService) Check(ctx context.Context, tenantID, ip string) (*IPRuleCheckResult, error) {
	var out IPRuleCheckResult
	body := map[string]string{"ip": ip}
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ip-rules/check", body, &out)
	return &out, err
}

// SetEnforcement turns IP-rule enforcement on or off for the tenant.
// Existing rules are unaffected either way — this only toggles whether
// they're evaluated on login.
func (s *IPRulesService) SetEnforcement(ctx context.Context, tenantID string, enabled bool) error {
	body := map[string]bool{"enabled": enabled}
	return s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/ip-rules/config", body, nil)
}
