package qeetid

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// ABACPoliciesService manages attribute-based access policies and evaluates them.
type ABACPoliciesService struct{ c *Client }

// List lists a tenant's policies.
//
// API: GET /v1/tenants/{tenantID}/abac/policies
func (s *ABACPoliciesService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[ABACPolicy], error) {
	return call[Page[ABACPolicy]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/abac/policies",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create creates a policy.
//
// API: POST /v1/tenants/{tenantID}/abac/policies
func (s *ABACPoliciesService) Create(ctx context.Context, tenantID string, params ABACPolicyParams, opts ...RequestOption) (*ABACPolicy, error) {
	return call[ABACPolicy](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/abac/policies",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Get returns a policy.
//
// API: GET /v1/tenants/{tenantID}/abac/policies/{id}
func (s *ABACPoliciesService) Get(ctx context.Context, tenantID, policyID string, opts ...RequestOption) (*ABACPolicy, error) {
	return call[ABACPolicy](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/abac/policies/{policyID}",
		pathArgs: []string{tenantID, policyID},
	}, opts)
}

// Update changes a policy and records a new version.
//
// API: PATCH /v1/tenants/{tenantID}/abac/policies/{id}
func (s *ABACPoliciesService) Update(ctx context.Context, tenantID, policyID string, params ABACPolicyParams, opts ...RequestOption) (*ABACPolicy, error) {
	return call[ABACPolicy](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/abac/policies/{policyID}",
		pathArgs: []string{tenantID, policyID},
		body:     params,
	}, opts)
}

// Delete deletes a policy.
//
// API: DELETE /v1/tenants/{tenantID}/abac/policies/{id}
func (s *ABACPoliciesService) Delete(ctx context.Context, tenantID, policyID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/abac/policies/{policyID}",
		pathArgs: []string{tenantID, policyID},
	}, nil, opts)
}

// ListVersions lists the saved versions of a policy.
//
// API: GET /v1/tenants/{tenantID}/abac/policies/{id}/versions
func (s *ABACPoliciesService) ListVersions(ctx context.Context, tenantID, policyID string, opts ...RequestOption) (*Page[ABACPolicyVersion], error) {
	return call[Page[ABACPolicyVersion]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/abac/policies/{policyID}/versions",
		pathArgs: []string{tenantID, policyID},
	}, opts)
}

// RestoreVersion makes a saved version of a policy current again.
//
// API: POST /v1/tenants/{tenantID}/abac/policies/{id}/versions/{version}/restore
func (s *ABACPoliciesService) RestoreVersion(ctx context.Context, tenantID, policyID, version string, opts ...RequestOption) (*ABACPolicy, error) {
	return call[ABACPolicy](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/abac/policies/{policyID}/versions/{version}/restore",
		pathArgs: []string{tenantID, policyID, version},
	}, opts)
}

// Evaluate evaluates the tenant's policies for a subject, resource and action.
//
// API: POST /v1/tenants/{tenantID}/abac/evaluate
func (s *ABACPoliciesService) Evaluate(ctx context.Context, tenantID string, params ABACEvaluateParams, options *ABACEvaluateOptions, opts ...RequestOption) (*ABACDecision, error) {
	return call[ABACDecision](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/abac/evaluate",
		pathArgs: []string{tenantID},
		query:    options.values(),
		body:     params,
	}, opts)
}

// ABACPolicy is an attribute-based access policy: an effect that applies to a resource type and action
// when its condition holds.
type ABACPolicy struct {
	ID          string     `json:"id,omitempty"`
	TenantID    string     `json:"tenant_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Effect      ABACEffect `json:"effect,omitempty"`
	// Resource type matched by this policy; '*' matches any.
	ResourceType string `json:"resource_type,omitempty"`
	// Action matched by this policy; '*' matches any.
	Action string `json:"action,omitempty"`
	// Recursive JSON condition tree. An empty object {} means unconditional. Logical nodes:
	// {"all":[...]}, {"any":[...]}, {"not": node}. Comparison leaf:
	// {"attr":"subject.dept","op":"eq","value":"eng"}. A leaf's value may instead reference another
	// attribute, {"attr":"subject.region","op":"eq","value":{"attr":"resource.region"}}, for eq, ne, gt,
	// gte, lt, lte, in, nin and contains (not regex, prefix, suffix or exists). The reference is exactly
	// {"attr":"<subject|resource|action|context>.<path>"}, names an attribute other than the leaf's own,
	// and is the whole value, never an array element. A string value is always a literal, even
	// "resource.region". A comparison with a missing attribute (absent or null), on either side and with
	// any operator, ne and nin included, is unknown; not of unknown is unknown, all is unknown unless a
	// child is false, and any unknown unless a child is true. A condition left unknown fails closed: a
	// deny policy applies (reason authz.attribute_absent) and an allow policy grants nothing. exists
	// alone tests presence and is never unknown, so a deny meant only for requests that carry an
	// attribute guards it with exists in an all (otherwise a deny for "*" on an attribute the platform's
	// own route check doesn't supply refuses every request there). prefix and suffix need a non-empty
	// string value and match only string attributes; any other value is refused with 422
	// abac.condition_invalid, as is contains or regex with an empty value. Two more leaves combine roles
	// and relationships with attributes: {"role":"<name>"} holds when the subject has the role, by name,
	// in the organisation (directly or through a group), as the platform supplies subject.roles when it
	// decides for a person itself; it is unknown when no roles are supplied.
	// {"relation":{"relation":"viewer", "object":"document:1"}} holds when the subject (subject.id, a
	// user) has the relationship on the object in the organisation's ReBAC tuples, the object given as
	// type:id or as an attribute reference {"attr":"resource.object"}; it is unknown with no subject id
	// or no object. Like a missing attribute, an unknown role or relation applies a deny and grants
	// nothing.
	Condition map[string]any `json:"condition,omitempty"`
	// Higher priority is evaluated first; deny wins at any priority.
	Priority  int       `json:"priority"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	// The policy's saves counted from 1; each create, change or restore adds one.
	Version int `json:"version"`
	// The person who made the last save; null for a machine caller, or one from before it was kept.
	UpdatedBy string `json:"updated_by,omitempty"`
}

// ABACEffect is a string enumeration used by ABACPolicy.Effect.
type ABACEffect string

// Values of ABACEffect. The server may add values; handle an unknown one gracefully.
const (
	ABACEffectAllow ABACEffect = "allow"
	ABACEffectDeny  ABACEffect = "deny"
)

// ABACPolicyParams is the request body of ABACPoliciesService.Create.
type ABACPolicyParams struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Effect      ABACEffect `json:"effect,omitempty"`
	// Resource type to match; '*' matches any resource type.
	ResourceType string `json:"resource_type,omitempty"`
	// Action to match; '*' matches any action.
	Action string `json:"action,omitempty"`
	// Optional condition tree. Omit or send {} for unconditional. Same grammar as ABACPolicy.condition,
	// including attribute references; a malformed condition or reference is refused with 422
	// abac.condition_invalid.
	Condition map[string]any `json:"condition,omitempty"`
	// Higher priority policies are evaluated first.
	Priority *int `json:"priority,omitempty"`
	// Whether the policy is evaluated; left out, it is saved disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// ABACPolicyVersion is one item of the list ABACPoliciesService.ListVersions returns.
type ABACPolicyVersion struct {
	Version      int            `json:"version"`
	Name         string         `json:"name,omitempty"`
	Description  string         `json:"description,omitempty"`
	Effect       ABACEffect     `json:"effect,omitempty"`
	ResourceType string         `json:"resource_type,omitempty"`
	Action       string         `json:"action,omitempty"`
	Condition    map[string]any `json:"condition,omitempty"`
	Priority     int            `json:"priority"`
	Enabled      bool           `json:"enabled"`
	SavedBy      string         `json:"saved_by,omitempty"`
	SavedAt      time.Time      `json:"saved_at,omitzero"`
}

// ABACEvaluateParams is the request body of ABACPoliciesService.Evaluate. Input for an ABAC policy
// evaluation. Subject/resource/context carry the attribute bags.
type ABACEvaluateParams struct {
	// Subject (caller) attributes, e.g. {"id":"user:123","department":"eng"}.
	Subject map[string]any `json:"subject,omitempty"`
	// Resource attributes.
	Resource *ABACEvaluateParamsResource `json:"resource,omitempty"`
	// The action being attempted, e.g. "read"; policies see it as action.name.
	Action string `json:"action,omitempty"`
	// Optional attributes of the action, the action namespace beside action.name (which wins over an
	// attribute of that name), e.g. {"channel":"api"}.
	ActionAttrs map[string]any `json:"action_attrs,omitempty"`
	// Request-context attributes, e.g. {"ip":"1.2.3.4","hour":14}.
	Context map[string]any `json:"context,omitempty"`
	// An unsaved policy to evaluate with the stored ones, writing nothing. It is checked as a save is
	// (name may be left out); with replaces, it stands in for that stored policy. It applies only when
	// enabled and its resource type and action match the request, so a disabled draft only takes the
	// replaced policy out.
	Draft *ABACEvaluateParamsDraft `json:"draft,omitempty"`
}

// ABACEvaluateParamsResource is the type of ABACEvaluateParams.Resource. Resource attributes.
type ABACEvaluateParamsResource struct {
	// Resource type, e.g. "document".
	Type string `json:"type,omitempty"`
	// Resource identifier.
	ID string `json:"id,omitempty"`
	// Additional resource attributes; type and id above win over attributes of those names.
	Attrs map[string]any `json:"attrs,omitempty"`
}

// ABACEvaluateParamsDraft is an unsaved policy to evaluate with the stored ones, writing nothing. It
// is checked as a save is (name may be left out); with replaces, it stands in for that stored policy.
// It applies only when enabled and its resource type and action match the request, so a disabled draft
// only takes the replaced policy out.
type ABACEvaluateParamsDraft struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Effect      ABACEffect `json:"effect,omitempty"`
	// Resource type to match; '*' matches any resource type.
	ResourceType string `json:"resource_type,omitempty"`
	// Action to match; '*' matches any action.
	Action string `json:"action,omitempty"`
	// Optional condition tree. Omit or send {} for unconditional. Same grammar as ABACPolicy.condition,
	// including attribute references; a malformed condition or reference is refused with 422
	// abac.condition_invalid.
	Condition map[string]any `json:"condition,omitempty"`
	// Higher priority policies are evaluated first.
	Priority *int `json:"priority,omitempty"`
	// Whether the policy is evaluated; left out, it is saved disabled.
	Enabled *bool `json:"enabled,omitempty"`
	// The stored policy the draft stands in for; one the organisation doesn't have is 404.
	Replaces string `json:"replaces,omitempty"`
}

// ABACEvaluateOptions holds the optional query parameters of ABACPoliciesService.Evaluate. A nil
// *ABACEvaluateOptions sends none.
type ABACEvaluateOptions struct {
	// When true, the response includes a human-readable trace of each evaluated policy.
	Explain bool
}

func (p *ABACEvaluateOptions) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Explain {
		q.Set("explain", "true")
	}
	return q
}

// ABACDecision is the outcome of evaluating a tenant's attribute-based policies.
type ABACDecision struct {
	// True when an allow policy matched and no deny policy matched.
	Allow  bool       `json:"allow"`
	Effect ABACEffect `json:"effect,omitempty"`
	// ID of the policy that determined the outcome, if any.
	MatchedPolicyID string `json:"matched_policy_id,omitempty"`
	// Human-readable explanation of the decision. Show reason_code's meaning instead.
	Reason string `json:"reason,omitempty"`
	// Why: an allow policy matched; a deny policy matched; a deny policy applied as an attribute it
	// compares is absent; a policy's condition failed to evaluate, so it refused (the error is logged,
	// not returned); or no allow policy matched.
	ReasonCode ABACDecisionReasonCode `json:"reason_code,omitempty"`
	// The policy that decided (policy_id, policy_name, priority); absent for authz.no_matching_allow.
	ReasonParams map[string]any `json:"reason_params,omitempty"`
	// Grant-path trace (present only when ?explain=true).
	Trace []string `json:"trace,omitempty"`
	// Present, true, when the request's draft policy determined the outcome. A new draft has no
	// matched_policy_id; one standing in for a stored policy has that policy's.
	MatchedDraft bool `json:"matched_draft"`
}

// ABACDecisionReasonCode is a string enumeration used by ABACDecision.ReasonCode. Why: an allow policy
// matched; a deny policy matched; a deny policy applied as an attribute it compares is absent; a
// policy's condition failed to evaluate, so it refused (the error is logged, not returned); or no
// allow policy matched.
type ABACDecisionReasonCode string

// Values of ABACDecisionReasonCode. The server may add values; handle an unknown one gracefully.
const (
	ABACDecisionReasonCodeAuthzAllowedByPolicy ABACDecisionReasonCode = "authz.allowed_by_policy"
	ABACDecisionReasonCodeAuthzDeniedByPolicy  ABACDecisionReasonCode = "authz.denied_by_policy"
	ABACDecisionReasonCodeAuthzAttributeAbsent ABACDecisionReasonCode = "authz.attribute_absent"
	ABACDecisionReasonCodeAuthzConditionError  ABACDecisionReasonCode = "authz.condition_error"
	ABACDecisionReasonCodeAuthzNoMatchingAllow ABACDecisionReasonCode = "authz.no_matching_allow"
)
