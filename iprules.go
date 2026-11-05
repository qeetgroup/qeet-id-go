package qeetid

import (
	"context"
	"net/http"
	"time"
)

// IPRulesService manages the network addresses a tenant's people may sign in from.
type IPRulesService struct{ c *Client }

// List lists a tenant's rules and whether they are enforced.
//
// API: GET /v1/tenants/{tenantID}/ip-rules
func (s *IPRulesService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*IPRules, error) {
	return call[IPRules](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/ip-rules",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create adds a rule.
//
// API: POST /v1/tenants/{tenantID}/ip-rules
func (s *IPRulesService) Create(ctx context.Context, tenantID string, params CreateIPRuleParams, opts ...RequestOption) (*IPRule, error) {
	return call[IPRule](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/ip-rules",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Delete deletes a rule.
//
// API: DELETE /v1/tenants/{tenantID}/ip-rules/{ruleID}
func (s *IPRulesService) Delete(ctx context.Context, tenantID string, ruleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/ip-rules/{ruleID}",
		pathArgs: []string{tenantID, ruleID},
	}, nil, opts)
}

// Check reports whether an address would be allowed under the tenant's rules.
//
// API: POST /v1/tenants/{tenantID}/ip-rules/check
func (s *IPRulesService) Check(ctx context.Context, tenantID string, params CheckIPParams, opts ...RequestOption) (*IPCheck, error) {
	return call[IPCheck](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/ip-rules/check",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// SetEnforcement turns rule enforcement on or off. The rules themselves are kept either way.
//
// API: PUT /v1/tenants/{tenantID}/ip-rules/config
func (s *IPRulesService) SetEnforcement(ctx context.Context, tenantID string, params SetIPEnforcementParams, opts ...RequestOption) (*IPEnforcement, error) {
	return call[IPEnforcement](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/ip-rules/config",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// IPRules is the result of IPRulesService.List.
type IPRules struct {
	Enabled bool     `json:"enabled"`
	Items   []IPRule `json:"items,omitempty"`
}

// IPRule allows or denies sign-in from a range of network addresses.
type IPRule struct {
	Action    string    `json:"action,omitempty"`
	CIDR      string    `json:"cidr,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	ID        string    `json:"id,omitempty"`
	Label     string    `json:"label,omitempty"`
	TenantID  string    `json:"tenant_id,omitempty"`
}

// CreateIPRuleParams is the request body of IPRulesService.Create.
type CreateIPRuleParams struct {
	Action IPRuleAction `json:"action,omitempty"`
	CIDR   string       `json:"cidr,omitempty"`
	Label  string       `json:"label,omitempty"`
}

// IPRuleAction is a string enumeration used by CreateIPRuleParams.Action.
type IPRuleAction string

// Values of IPRuleAction. The server may add values; handle an unknown one gracefully.
const (
	IPRuleActionAllow IPRuleAction = "allow"
	IPRuleActionDeny  IPRuleAction = "deny"
)

// CheckIPParams is the request body of IPRulesService.Check.
type CheckIPParams struct {
	IP string `json:"ip,omitempty"`
}

// IPCheck is the result of IPRulesService.Check.
type IPCheck struct {
	Allowed bool `json:"allowed"`
	Enabled bool `json:"enabled"`
	// The reason in English. Show reason_code's meaning instead.
	Reason string `json:"reason,omitempty"`
	// Allowed; refused by a deny rule; outside every allow rule; refused as the address can't be read; or
	// allowed because the rules are off.
	ReasonCode IPCheckReasonCode `json:"reason_code,omitempty"`
	// For ip.deny_rule, the rule's cidr.
	ReasonParams map[string]any `json:"reason_params,omitempty"`
	// What the organisation's rules decide whatever the enforcement flag, so they can be tried before
	// they are turned on; with enforcement on it matches the fields above.
	Rules IPCheckRules `json:"rules,omitzero"`
}

// IPCheckReasonCode is a string enumeration used by IPCheck.ReasonCode. Allowed; refused by a deny
// rule; outside every allow rule; refused as the address can't be read; or allowed because the rules
// are off.
type IPCheckReasonCode string

// Values of IPCheckReasonCode. The server may add values; handle an unknown one gracefully.
const (
	IPCheckReasonCodeIPAllowed             IPCheckReasonCode = "ip.allowed"
	IPCheckReasonCodeIPDenyRule            IPCheckReasonCode = "ip.deny_rule"
	IPCheckReasonCodeIPNotInAllowRules     IPCheckReasonCode = "ip.not_in_allow_rules"
	IPCheckReasonCodeIPAddressUnknown      IPCheckReasonCode = "ip.address_unknown"
	IPCheckReasonCodeIPEnforcementDisabled IPCheckReasonCode = "ip.enforcement_disabled"
)

// IPCheckRules is the type of IPCheck.Rules. What the organisation's rules decide whatever the
// enforcement flag, so they can be tried before they are turned on; with enforcement on it matches the
// fields above.
type IPCheckRules struct {
	Allowed bool `json:"allowed"`
	// The reason in English.
	Reason     string                 `json:"reason,omitempty"`
	ReasonCode IPCheckRulesReasonCode `json:"reason_code,omitempty"`
	// For ip.deny_rule, the rule's cidr.
	ReasonParams map[string]any `json:"reason_params,omitempty"`
}

// IPCheckRulesReasonCode is a string enumeration used by IPCheckRules.ReasonCode.
type IPCheckRulesReasonCode string

// Values of IPCheckRulesReasonCode. The server may add values; handle an unknown one gracefully.
const (
	IPCheckRulesReasonCodeIPAllowed         IPCheckRulesReasonCode = "ip.allowed"
	IPCheckRulesReasonCodeIPDenyRule        IPCheckRulesReasonCode = "ip.deny_rule"
	IPCheckRulesReasonCodeIPNotInAllowRules IPCheckRulesReasonCode = "ip.not_in_allow_rules"
	IPCheckRulesReasonCodeIPAddressUnknown  IPCheckRulesReasonCode = "ip.address_unknown"
)

// SetIPEnforcementParams is the request body of IPRulesService.SetEnforcement.
type SetIPEnforcementParams struct {
	Enabled bool `json:"enabled"`
}

// IPEnforcement is the result of IPRulesService.SetEnforcement.
type IPEnforcement struct {
	Enabled bool `json:"enabled"`
}
