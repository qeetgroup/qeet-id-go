package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Permission is a platform-defined permission key (e.g. "billing:write").
// The catalog is fixed by the platform, not user-creatable — there is no
// Create/Update/Delete for permissions, only List and the checks below.
type Permission struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
}

// PermissionCheck is a single RBAC authorization query (maps to GET /v1/check).
type PermissionCheck struct {
	User       string
	Tenant     string
	Permission string
}

// AuthzExplanation is the response shape for GET /v1/check?explain=true.
type AuthzExplanation struct {
	Allowed bool             `json:"allowed"`
	Paths   []AuthzGrantStep `json:"paths,omitempty"`
	Reason  string           `json:"reason,omitempty"` // set on denial only
}

// AuthzGrantStep is one grant in an authorization explanation's path.
type AuthzGrantStep struct {
	Permission string `json:"permission"`
	GrantedBy  string `json:"granted_by"` // "role:<name>"
	Via        string `json:"via"`        // "direct" | "group:<name>"
	GroupID    string `json:"group_id,omitempty"`
	RoleID     string `json:"role_id"`
}

// PermissionsService lists the platform's permission catalog and provides
// the Check/CheckAll/Explain authorization queries — kept here rather than a
// bespoke root-level method set, since a permission check is naturally a
// Permissions operation.
type PermissionsService struct{ t *transport.Transport }

// List returns the full permission catalog. This is a platform-wide list,
// not tenant-scoped, and not paginated.
func (s *PermissionsService) List(ctx context.Context) ([]Permission, error) {
	var env marshal.Envelope[Permission]
	err := s.t.Get(ctx, "/v1/permissions", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// Effective returns every permission key a user currently holds within a
// tenant (the union of direct role grants and group-derived ones) — the
// resolved result Check/Explain reason about internally.
func (s *PermissionsService) Effective(ctx context.Context, userID, tenantID string) ([]string, error) {
	path := "/v1/users/" + url.PathEscape(userID) + "/tenants/" + url.PathEscape(tenantID) + "/permissions"
	var out struct {
		Permissions []string `json:"permissions"`
	}
	err := s.t.Get(ctx, path, transport.RequestOptions{}, &out)
	return out.Permissions, err
}

// Check resolves a single permission check — the hot-path call made on
// nearly every request.
func (s *PermissionsService) Check(ctx context.Context, check PermissionCheck) (bool, error) {
	q := url.Values{}
	q.Set("user_id", check.User)
	q.Set("tenant_id", check.Tenant)
	q.Set("permission", check.Permission)
	var res struct {
		Allowed bool `json:"allowed"`
	}
	if err := s.t.Get(ctx, "/v1/check", transport.RequestOptions{Query: q}, &res); err != nil {
		return false, err
	}
	return res.Allowed, nil
}

// CheckAll returns true only if every permission passes.
func (s *PermissionsService) CheckAll(ctx context.Context, user, tenant string, permissions []string) (bool, error) {
	for _, p := range permissions {
		ok, err := s.Check(ctx, PermissionCheck{User: user, Tenant: tenant, Permission: p})
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// Explain resolves a permission check and returns the grant path that
// decided it — which role, and whether it came from a direct assignment or
// a group. Denials carry a Reason instead of Paths.
func (s *PermissionsService) Explain(ctx context.Context, check PermissionCheck) (*AuthzExplanation, error) {
	q := url.Values{}
	q.Set("user_id", check.User)
	q.Set("tenant_id", check.Tenant)
	q.Set("permission", check.Permission)
	q.Set("explain", "true")
	var res AuthzExplanation
	if err := s.t.Get(ctx, "/v1/check", transport.RequestOptions{Query: q}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
