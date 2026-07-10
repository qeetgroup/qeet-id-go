package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Role is a tenant-scoped RBAC role. Permissions aren't embedded here — grant
// them individually via GrantPermission (a role's permission set is a
// separate collection, not a field on the role itself).
type Role struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsSystem    bool   `json:"is_system"`
	CreatedAt   string `json:"created_at"`
}

type CreateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// RolesService manages RBAC roles. There is no per-role Get/Update/Delete in
// the backend — only List, Create, and the assign/permission-grant
// operations below.
type RolesService struct{ t *transport.Transport }

// List returns every role defined for a tenant. Not paginated — the backend
// returns the full set in one response.
func (s *RolesService) List(ctx context.Context, tenantID string) ([]Role, error) {
	var env marshal.Envelope[Role]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/roles", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *RolesService) Create(ctx context.Context, tenantID string, in CreateRoleInput) (*Role, error) {
	var out Role
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/roles", in, &out)
	return &out, err
}

// GrantPermission adds a permission to a role — every user holding the role
// immediately gains it.
func (s *RolesService) GrantPermission(ctx context.Context, roleID, permissionID string) error {
	path := "/v1/roles/" + url.PathEscape(roleID) + "/permissions/" + url.PathEscape(permissionID)
	return s.t.Post(ctx, path, nil, nil)
}

// RevokePermission removes a permission from a role.
func (s *RolesService) RevokePermission(ctx context.Context, roleID, permissionID string) error {
	path := "/v1/roles/" + url.PathEscape(roleID) + "/permissions/" + url.PathEscape(permissionID)
	return s.t.Delete(ctx, path, nil)
}

// AssignToUser grants a role to a user within a tenant.
func (s *RolesService) AssignToUser(ctx context.Context, userID, tenantID, roleID string) error {
	path := "/v1/users/" + url.PathEscape(userID) + "/tenants/" + url.PathEscape(tenantID) + "/roles/" + url.PathEscape(roleID)
	return s.t.Post(ctx, path, nil, nil)
}

// RemoveFromUser revokes a role previously assigned to a user.
func (s *RolesService) RemoveFromUser(ctx context.Context, userID, tenantID, roleID string) error {
	path := "/v1/users/" + url.PathEscape(userID) + "/tenants/" + url.PathEscape(tenantID) + "/roles/" + url.PathEscape(roleID)
	return s.t.Delete(ctx, path, nil)
}
