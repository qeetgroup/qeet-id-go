package qeetid

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// RolesService manages a tenant's roles, the permissions they grant and who holds them.
type RolesService struct{ c *Client }

// List lists a tenant's roles.
//
// API: GET /v1/tenants/{tenantID}/roles
func (s *RolesService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Role], error) {
	return call[Page[Role]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/roles",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create creates a role.
//
// API: POST /v1/tenants/{tenantID}/roles
func (s *RolesService) Create(ctx context.Context, tenantID string, params CreateRoleParams, opts ...RequestOption) (*Role, error) {
	return call[Role](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/roles",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Update renames a role or changes its description.
//
// API: PATCH /v1/tenants/{tenantID}/roles/{roleID}
func (s *RolesService) Update(ctx context.Context, tenantID, roleID string, params UpdateRoleParams, opts ...RequestOption) (*Role, error) {
	return call[Role](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/roles/{roleID}",
		pathArgs: []string{tenantID, roleID},
		body:     params,
	}, opts)
}

// Delete deletes a role.
//
// API: DELETE /v1/tenants/{tenantID}/roles/{roleID}
func (s *RolesService) Delete(ctx context.Context, tenantID, roleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/roles/{roleID}",
		pathArgs: []string{tenantID, roleID},
	}, nil, opts)
}

// ListPermissions lists the permissions a role grants.
//
// API: GET /v1/roles/{roleID}/permissions
func (s *RolesService) ListPermissions(ctx context.Context, roleID string, opts ...RequestOption) (*Page[Permission], error) {
	return call[Page[Permission]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/roles/{roleID}/permissions",
		pathArgs: []string{roleID},
	}, opts)
}

// SetPermissions replaces the permissions a role grants and returns the new set.
//
// API: PUT /v1/tenants/{tenantID}/roles/{roleID}/permissions
func (s *RolesService) SetPermissions(ctx context.Context, tenantID, roleID string, params SetRolePermissionsParams, opts ...RequestOption) (*Page[Permission], error) {
	return call[Page[Permission]](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/roles/{roleID}/permissions",
		pathArgs: []string{tenantID, roleID},
		body:     params,
	}, opts)
}

// GrantPermission adds a permission to a role.
//
// API: POST /v1/roles/{roleID}/permissions/{permID}
func (s *RolesService) GrantPermission(ctx context.Context, roleID, permissionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/roles/{roleID}/permissions/{permissionID}",
		pathArgs: []string{roleID, permissionID},
	}, nil, opts)
}

// RevokePermission removes a permission from a role.
//
// API: DELETE /v1/roles/{roleID}/permissions/{permID}
func (s *RolesService) RevokePermission(ctx context.Context, roleID, permissionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/roles/{roleID}/permissions/{permissionID}",
		pathArgs: []string{roleID, permissionID},
	}, nil, opts)
}

// Holders lists the people and groups that hold a role, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/roles/{roleID}/holders
func (s *RolesService) Holders(ctx context.Context, tenantID, roleID string, params *RoleHoldersParams, opts ...RequestOption) (*RoleHolders, error) {
	return call[RoleHolders](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/roles/{roleID}/holders",
		pathArgs: []string{tenantID, roleID},
		query:    params.values(),
	}, opts)
}

// AssignToUser gives a user a role in a tenant.
//
// API: POST /v1/users/{userID}/tenants/{tenantID}/roles/{roleID}
func (s *RolesService) AssignToUser(ctx context.Context, userID, tenantID, roleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/tenants/{tenantID}/roles/{roleID}",
		pathArgs: []string{userID, tenantID, roleID},
	}, nil, opts)
}

// UnassignFromUser takes a role away from a user.
//
// API: DELETE /v1/users/{userID}/tenants/{tenantID}/roles/{roleID}
func (s *RolesService) UnassignFromUser(ctx context.Context, userID, tenantID, roleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/users/{userID}/tenants/{tenantID}/roles/{roleID}",
		pathArgs: []string{userID, tenantID, roleID},
	}, nil, opts)
}

// Role is a named set of permissions in a tenant.
type Role struct {
	ID          string `json:"id,omitempty"`
	TenantID    string `json:"tenant_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// A built-in role (owner, admin, member or viewer), whose permissions can't be changed and which
	// doesn't count against the plan's custom roles.
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	// Who holds the role, in the roles list only.
	Holders RoleHolderCounts `json:"holders,omitzero"`
}

// RoleHolderCounts is the type of Role.Holders. Who holds the role, in the roles list only.
type RoleHolderCounts struct {
	// People granted the role themselves.
	PeopleDirect int `json:"people_direct"`
	// People holding it only through a group.
	PeopleThroughGroups int `json:"people_through_groups"`
	// Groups granted the role.
	Groups int `json:"groups"`
}

// CreateRoleParams is the request body of RolesService.Create.
type CreateRoleParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateRoleParams is the request body of RolesService.Update.
type UpdateRoleParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Permission is an action a role can grant, such as "user.write". The catalog is fixed by Qeet ID.
type Permission struct {
	ID          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
}

// SetRolePermissionsParams is the request body of RolesService.SetPermissions.
type SetRolePermissionsParams struct {
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

// RoleHoldersParams holds the optional query parameters of RolesService.Holders. A nil
// *RoleHoldersParams sends none.
type RoleHoldersParams struct {
	Limit int
	// The next_cursor of the previous page.
	Cursor string
}

func (p *RoleHoldersParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

// RoleHolders is the result of RolesService.Holders.
type RoleHolders struct {
	People     []RoleHolderPerson `json:"people,omitempty"`
	Groups     []RoleHolderGroup  `json:"groups,omitempty"`
	NextCursor string             `json:"next_cursor,omitempty"`
}

// RoleHolderPerson is one element of RoleHolders.People.
type RoleHolderPerson struct {
	UserID      string               `json:"user_id,omitempty"`
	Email       string               `json:"email,omitempty"`
	DisplayName string               `json:"display_name,omitempty"`
	Direct      bool                 `json:"direct"`
	ViaGroups   []RoleHolderViaGroup `json:"via_groups,omitempty"`
}

// RoleHolderViaGroup is one element of RoleHolderPerson.ViaGroups.
type RoleHolderViaGroup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RoleHolderGroup is one element of RoleHolders.Groups.
type RoleHolderGroup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
