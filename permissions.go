package qeetid

import (
	"context"
	"net/http"
	"net/url"
)

// PermissionsService lists the permission catalog and answers role-based access checks.
type PermissionsService struct{ c *Client }

// List lists every permission a role can grant.
//
// API: GET /v1/permissions
func (s *PermissionsService) List(ctx context.Context, opts ...RequestOption) (*Page[Permission], error) {
	return call[Page[Permission]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/permissions"}, opts)
}

// Check reports whether a user holds a permission in a tenant. Set Explain to get the grant path that
// decided it.
//
// API: GET /v1/check
func (s *PermissionsService) Check(ctx context.Context, params *CheckPermissionParams, opts ...RequestOption) (*PermissionCheck, error) {
	return call[PermissionCheck](ctx, s.c, &request{method: http.MethodGet, path: "/v1/check", query: params.values()}, opts)
}

// Effective returns every permission a user holds in a tenant, directly or through groups, with the
// grants behind them.
//
// API: GET /v1/users/{userID}/tenants/{tenantID}/permissions
func (s *PermissionsService) Effective(ctx context.Context, userID string, tenantID string, opts ...RequestOption) (*EffectivePermissions, error) {
	return call[EffectivePermissions](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/users/{userID}/tenants/{tenantID}/permissions",
		pathArgs: []string{userID, tenantID},
	}, opts)
}

// CheckPermissionParams holds the optional query parameters of PermissionsService.Check. A nil
// *CheckPermissionParams sends none.
type CheckPermissionParams struct {
	UserID     string
	TenantID   string
	Permission string
	// When "true", also return the role grant path(s), a reason on denial and the decision's trace.
	Explain bool
}

func (p *CheckPermissionParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.UserID != "" {
		q.Set("user_id", p.UserID)
	}
	if p.TenantID != "" {
		q.Set("tenant_id", p.TenantID)
	}
	if p.Permission != "" {
		q.Set("permission", p.Permission)
	}
	if p.Explain {
		q.Set("explain", "true")
	}
	return q
}

// PermissionCheck is the answer to a role-based access check.
type PermissionCheck struct {
	Allowed   bool           `json:"allowed"`
	DecidedBy AuthzDecidedBy `json:"decided_by,omitzero"`
	Paths     []GrantStep    `json:"paths,omitempty"`
	// Set only on a denial, with explain=true, in English. Show reason_code's meaning instead.
	Reason string `json:"reason,omitempty"`
	// Set with reason. authz.no_role_grant when no role grants the permission; an ABAC code when a policy
	// refused it.
	ReasonCode PermissionCheckReasonCode `json:"reason_code,omitempty"`
	// With an ABAC code, the policy (policy_id, policy_name, priority), to a caller holding policy.read.
	ReasonParams map[string]any `json:"reason_params,omitempty"`
	// With explain=true. Each step in order, prefixed with its model (tenant, abac, rbac, rebac, none).
	Trace []string `json:"trace,omitempty"`
}

// AuthzDecidedBy is the type of PermissionCheck.DecidedBy. What decided a composed authorization
// decision. The ABAC policy (rule, policy_id) is named only to a caller holding policy.read; anyone
// else sees rule "an attribute policy".
type AuthzDecidedBy struct {
	Model AuthzDecidedByModel `json:"model,omitempty"`
	// What decided it, in English. Show reason_code's meaning instead.
	Rule string `json:"rule,omitempty"`
	// What decided it: no organisation, an engine failure (refused), a role or a relationship tuple
	// granting it, nothing allowing it, or an ABAC policy (allowing, denying, denying as an attribute it
	// compares is absent, or refusing as its condition failed to evaluate).
	ReasonCode AuthzDecidedByReasonCode `json:"reason_code,omitempty"`
	// What reason_code names: permission for authz.role_grant; relation and object for
	// authz.relation_grant; policy_id, policy_name and priority for an ABAC policy, shown only to a
	// caller holding policy.read.
	ReasonParams map[string]any `json:"reason_params,omitempty"`
	PolicyID     string         `json:"policy_id,omitempty"`
}

// AuthzDecidedByModel is a string enumeration used by AuthzDecidedBy.Model.
type AuthzDecidedByModel string

// Values of AuthzDecidedByModel. The server may add values; handle an unknown one gracefully.
const (
	AuthzDecidedByModelTenant AuthzDecidedByModel = "tenant"
	AuthzDecidedByModelABAC   AuthzDecidedByModel = "abac"
	AuthzDecidedByModelRBAC   AuthzDecidedByModel = "rbac"
	AuthzDecidedByModelReBAC  AuthzDecidedByModel = "rebac"
	AuthzDecidedByModelNone   AuthzDecidedByModel = "none"
)

// AuthzDecidedByReasonCode is a string enumeration used by AuthzDecidedBy.ReasonCode. What decided it:
// no organisation, an engine failure (refused), a role or a relationship tuple granting it, nothing
// allowing it, or an ABAC policy (allowing, denying, denying as an attribute it compares is absent, or
// refusing as its condition failed to evaluate).
type AuthzDecidedByReasonCode string

// Values of AuthzDecidedByReasonCode. The server may add values; handle an unknown one gracefully.
const (
	AuthzDecidedByReasonCodeAuthzNoTenant        AuthzDecidedByReasonCode = "authz.no_tenant"
	AuthzDecidedByReasonCodeAuthzEngineError     AuthzDecidedByReasonCode = "authz.engine_error"
	AuthzDecidedByReasonCodeAuthzRoleGrant       AuthzDecidedByReasonCode = "authz.role_grant"
	AuthzDecidedByReasonCodeAuthzRelationGrant   AuthzDecidedByReasonCode = "authz.relation_grant"
	AuthzDecidedByReasonCodeAuthzNothingAllows   AuthzDecidedByReasonCode = "authz.nothing_allows"
	AuthzDecidedByReasonCodeAuthzAllowedByPolicy AuthzDecidedByReasonCode = "authz.allowed_by_policy"
	AuthzDecidedByReasonCodeAuthzDeniedByPolicy  AuthzDecidedByReasonCode = "authz.denied_by_policy"
	AuthzDecidedByReasonCodeAuthzAttributeAbsent AuthzDecidedByReasonCode = "authz.attribute_absent"
	AuthzDecidedByReasonCodeAuthzConditionError  AuthzDecidedByReasonCode = "authz.condition_error"
)

// GrantStep is one element of PermissionCheck.Paths. One link in an authorization decision's grant
// path: a role that confers the requested permission and how that role reaches the user.
type GrantStep struct {
	Permission string `json:"permission,omitempty"`
	// The conferring role, e.g. "role:admin".
	GrantedBy string `json:"granted_by,omitempty"`
	// "direct" for a user grant, or "group:<name>" for a group-derived grant.
	Via string `json:"via,omitempty"`
	// Set only when via is a group grant.
	GroupID string `json:"group_id,omitempty"`
	// The group's name, set only when via is a group grant.
	GroupName string `json:"group_name,omitempty"`
	RoleID    string `json:"role_id,omitempty"`
	RoleName  string `json:"role_name,omitempty"`
}

// PermissionCheckReasonCode is a string enumeration used by PermissionCheck.ReasonCode. Set with
// reason. authz.no_role_grant when no role grants the permission; an ABAC code when a policy refused
// it.
type PermissionCheckReasonCode string

// Values of PermissionCheckReasonCode. The server may add values; handle an unknown one gracefully.
const (
	PermissionCheckReasonCodeAuthzNoRoleGrant     PermissionCheckReasonCode = "authz.no_role_grant"
	PermissionCheckReasonCodeAuthzDeniedByPolicy  PermissionCheckReasonCode = "authz.denied_by_policy"
	PermissionCheckReasonCodeAuthzAttributeAbsent PermissionCheckReasonCode = "authz.attribute_absent"
	PermissionCheckReasonCodeAuthzConditionError  PermissionCheckReasonCode = "authz.condition_error"
)

// EffectivePermissions is the result of PermissionsService.Effective.
type EffectivePermissions struct {
	Permissions []string    `json:"permissions,omitempty"`
	Grants      []GrantStep `json:"grants,omitempty"`
}
