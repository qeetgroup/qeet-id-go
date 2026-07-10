package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Group is a tenant-scoped group, optionally nested under a parent. There is
// no per-group Get or Update in the backend — only Create, Delete, List, and
// membership/role operations.
type Group struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	ParentID    string `json:"parent_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// CreateGroupInput — TenantID is deliberately not a field here: the backend
// always derives it from the caller's own API key and ignores any value
// sent in the body.
type CreateGroupInput struct {
	ParentID    string `json:"parent_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type GroupMember struct {
	UserID  string `json:"user_id"`
	GroupID string `json:"group_id"`
	AddedAt string `json:"added_at"`
}

// GroupRole is a role granted to a group — inherited by every member.
type GroupRole struct {
	RoleID    string `json:"role_id"`
	Name      string `json:"name"`
	GrantedAt string `json:"granted_at"`
}

// GroupsService manages groups and group membership. Create/Delete/member
// operations are NOT tenant-path-scoped (Create derives tenant from the
// caller's API key; Delete/member ops operate on an already-tenant-owned
// group ID) — only List and the role-binding methods take an explicit
// tenantID, matching the backend's own path shapes.
type GroupsService struct{ t *transport.Transport }

func (s *GroupsService) Create(ctx context.Context, in CreateGroupInput) (*Group, error) {
	var out Group
	err := s.t.Post(ctx, "/v1/groups", in, &out)
	return &out, err
}

func (s *GroupsService) Delete(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/groups/"+url.PathEscape(id), nil)
}

func (s *GroupsService) List(ctx context.Context, tenantID string) ([]Group, error) {
	var env marshal.Envelope[Group]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/groups", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// AddMember adds a user to a group — no request body, both IDs are path
// parameters.
func (s *GroupsService) AddMember(ctx context.Context, groupID, userID string) error {
	path := "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(userID)
	return s.t.Post(ctx, path, nil, nil)
}

func (s *GroupsService) RemoveMember(ctx context.Context, groupID, userID string) error {
	return s.t.Delete(ctx, "/v1/groups/"+url.PathEscape(groupID)+"/members/"+url.PathEscape(userID), nil)
}

func (s *GroupsService) ListMembers(ctx context.Context, groupID string) ([]GroupMember, error) {
	var env marshal.Envelope[GroupMember]
	err := s.t.Get(ctx, "/v1/groups/"+url.PathEscape(groupID)+"/members", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// ListRoles lists the roles granted directly to a group (tenant-scoped,
// distinct from the /v1/groups/{id} management path above).
func (s *GroupsService) ListRoles(ctx context.Context, tenantID, groupID string) ([]GroupRole, error) {
	path := "/v1/tenants/" + url.PathEscape(tenantID) + "/groups/" + url.PathEscape(groupID) + "/roles"
	var env marshal.Envelope[GroupRole]
	err := s.t.Get(ctx, path, transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// GrantRole grants a role to a group — every member inherits it.
func (s *GroupsService) GrantRole(ctx context.Context, tenantID, groupID, roleID string) error {
	path := "/v1/tenants/" + url.PathEscape(tenantID) + "/groups/" + url.PathEscape(groupID) + "/roles/" + url.PathEscape(roleID)
	return s.t.Post(ctx, path, nil, nil)
}

// RevokeRole revokes a role previously granted to a group.
func (s *GroupsService) RevokeRole(ctx context.Context, tenantID, groupID, roleID string) error {
	path := "/v1/tenants/" + url.PathEscape(tenantID) + "/groups/" + url.PathEscape(groupID) + "/roles/" + url.PathEscape(roleID)
	return s.t.Delete(ctx, path, nil)
}
