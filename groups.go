package qeetid

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// GroupsService manages groups, their members and the roles granted to them.
type GroupsService struct{ c *Client }

// Create creates a group.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/groups
func (s *GroupsService) Create(ctx context.Context, params CreateGroupParams, opts ...RequestOption) (*Group, error) {
	return call[Group](ctx, s.c, &request{method: http.MethodPost, path: "/v1/groups", body: params, idempotent: true}, opts)
}

// Update changes a group's name, description or parent.
//
// API: PATCH /v1/groups/{id}
func (s *GroupsService) Update(ctx context.Context, groupID string, params UpdateGroupParams, opts ...RequestOption) (*Group, error) {
	return call[Group](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/groups/{groupID}",
		pathArgs: []string{groupID},
		body:     params,
	}, opts)
}

// Delete deletes a group.
//
// API: DELETE /v1/groups/{id}
func (s *GroupsService) Delete(ctx context.Context, groupID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{method: http.MethodDelete, path: "/v1/groups/{groupID}", pathArgs: []string{groupID}}, nil, opts)
}

// List lists a tenant's groups.
//
// API: GET /v1/tenants/{tenantID}/groups
func (s *GroupsService) List(ctx context.Context, tenantID string, params *ListGroupsParams, opts ...RequestOption) (*Page[Group], error) {
	return call[Page[Group]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/groups",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// ListMembers lists a group's members. Page.Truncated reports a list cut short by the server's cap.
//
// API: GET /v1/groups/{id}/members
func (s *GroupsService) ListMembers(ctx context.Context, groupID string, opts ...RequestOption) (*Page[GroupMember], error) {
	return call[Page[GroupMember]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/groups/{groupID}/members",
		pathArgs: []string{groupID},
	}, opts)
}

// AddMember adds a user to a group.
//
// API: POST /v1/groups/{id}/members/{userID}
func (s *GroupsService) AddMember(ctx context.Context, groupID string, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/groups/{groupID}/members/{userID}",
		pathArgs: []string{groupID, userID},
	}, nil, opts)
}

// RemoveMember removes a user from a group.
//
// API: DELETE /v1/groups/{id}/members/{userID}
func (s *GroupsService) RemoveMember(ctx context.Context, groupID string, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/groups/{groupID}/members/{userID}",
		pathArgs: []string{groupID, userID},
	}, nil, opts)
}

// ListRoles lists the roles granted to a group. Every member inherits them.
//
// API: GET /v1/tenants/{tenantID}/groups/{groupID}/roles
func (s *GroupsService) ListRoles(ctx context.Context, tenantID string, groupID string, opts ...RequestOption) (*Page[GroupRole], error) {
	return call[Page[GroupRole]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/groups/{groupID}/roles",
		pathArgs: []string{tenantID, groupID},
	}, opts)
}

// GrantRole grants a role to a group.
//
// API: POST /v1/tenants/{tenantID}/groups/{groupID}/roles/{roleID}
func (s *GroupsService) GrantRole(ctx context.Context, tenantID string, groupID string, roleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/groups/{groupID}/roles/{roleID}",
		pathArgs: []string{tenantID, groupID, roleID},
	}, nil, opts)
}

// RevokeRole revokes a role from a group.
//
// API: DELETE /v1/tenants/{tenantID}/groups/{groupID}/roles/{roleID}
func (s *GroupsService) RevokeRole(ctx context.Context, tenantID string, groupID string, roleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/groups/{groupID}/roles/{roleID}",
		pathArgs: []string{tenantID, groupID, roleID},
	}, nil, opts)
}

// Import creates many groups in one request and reports the outcome of each row.
//
// API: POST /v1/groups/import
func (s *GroupsService) Import(ctx context.Context, params ImportGroupsParams, opts ...RequestOption) (*GroupImportResult, error) {
	return call[GroupImportResult](ctx, s.c, &request{method: http.MethodPost, path: "/v1/groups/import", body: params}, opts)
}

// ListTemplates lists the tenant's group templates.
//
// API: GET /v1/groups/templates
func (s *GroupsService) ListTemplates(ctx context.Context, opts ...RequestOption) (*Page[GroupTemplate], error) {
	return call[Page[GroupTemplate]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/groups/templates"}, opts)
}

// CreateTemplate creates a group template: a tree of groups that can be created in one step.
//
// API: POST /v1/groups/templates
func (s *GroupsService) CreateTemplate(ctx context.Context, params CreateGroupTemplateParams, opts ...RequestOption) (*GroupTemplate, error) {
	return call[GroupTemplate](ctx, s.c, &request{method: http.MethodPost, path: "/v1/groups/templates", body: params}, opts)
}

// DeleteTemplate deletes a group template.
//
// API: DELETE /v1/groups/templates/{id}
func (s *GroupsService) DeleteTemplate(ctx context.Context, templateID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/groups/templates/{templateID}",
		pathArgs: []string{templateID},
	}, nil, opts)
}

// ApplyTemplate creates the groups a template describes and returns them.
//
// API: POST /v1/groups/templates/{id}/apply
func (s *GroupsService) ApplyTemplate(ctx context.Context, templateID string, params ApplyGroupTemplateParams, opts ...RequestOption) ([]Group, error) {
	var out struct {
		V []Group `json:"items"`
	}
	if err := s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/groups/templates/{templateID}/apply",
		pathArgs: []string{templateID},
		body:     params,
	}, &out, opts); err != nil {
		return nil, err
	}
	return out.V, nil
}

// CreateGroupParams is the request body of GroupsService.Create.
type CreateGroupParams struct {
	Description string  `json:"description,omitempty"`
	Name        string  `json:"name,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	TenantID    string  `json:"tenant_id,omitempty"`
}

// Group is a set of users in a tenant. Groups can nest, and the roles granted to a group pass to its
// members.
type Group struct {
	CreatedAt   time.Time `json:"created_at,omitzero"`
	Description string    `json:"description,omitempty"`
	ID          string    `json:"id,omitempty"`
	MemberCount int       `json:"member_count"`
	Name        string    `json:"name,omitempty"`
	ParentID    string    `json:"parent_id,omitempty"`
	TenantID    string    `json:"tenant_id,omitempty"`
	// The roles granted to the group itself, in the list only and only for a caller with role.read.
	Roles []GroupRole `json:"roles,omitempty"`
}

// GroupRole is a role granted to a group.
type GroupRole struct {
	RoleID string `json:"role_id,omitempty"`
	Name   string `json:"name,omitempty"`
	// GrantedAt is when the role was granted to the group. It is set by GroupsService.ListRoles.
	GrantedAt time.Time `json:"granted_at,omitzero"`
}

// GroupGrantsFilter is a string enumeration used by ListGroupsParams.Grants.
type GroupGrantsFilter string

// Values of GroupGrantsFilter. The server may add values; handle an unknown one gracefully.
const (
	GroupGrantsFilterSome GroupGrantsFilter = "some"
	GroupGrantsFilterNone GroupGrantsFilter = "none"
)

// ListGroupsParams holds the optional query parameters of GroupsService.List. A nil *ListGroupsParams
// sends none.
type ListGroupsParams struct {
	// Only the groups granted this role. Needs role.read.
	Role string
	// Only the groups that grant some role, or none. Needs role.read.
	Grants GroupGrantsFilter
}

func (p *ListGroupsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Role != "" {
		q.Set("role", p.Role)
	}
	if p.Grants != "" {
		q.Set("grants", string(p.Grants))
	}
	return q
}

// GroupMember is a user's membership of a group.
type GroupMember struct {
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

// ImportGroupsParams is the request body of GroupsService.Import.
type ImportGroupsParams struct {
	Groups []ImportGroup `json:"groups,omitempty"`
}

// ImportGroup is one element of ImportGroupsParams.Groups.
type ImportGroup struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// Parent group name, resolved against existing and imported rows.
	Parent string `json:"parent,omitempty"`
}

// GroupImportResult is the result of GroupsService.Import.
type GroupImportResult struct {
	Created int              `json:"created"`
	Failed  int              `json:"failed"`
	Rows    []GroupImportRow `json:"rows,omitempty"`
	Skipped int              `json:"skipped"`
}

// GroupImportRow is one element of GroupImportResult.Rows.
type GroupImportRow struct {
	// Set when the row isn't created: it has no name; a group of that name exists, or an earlier row made
	// one; its parent is no existing group and no other row (or the rows form a cycle); or it failed on
	// the server. Word the message from this.
	Code GroupImportRowCode `json:"code,omitempty"`
	// The row field the code concerns, when one does.
	Field GroupImportRowField `json:"field,omitempty"`
	Name  string              `json:"name,omitempty"`
	// The code's meaning in English, kept for compatibility.
	Reason string               `json:"reason,omitempty"`
	Status GroupImportRowStatus `json:"status,omitempty"`
}

// GroupImportRowCode is a string enumeration used by GroupImportRow.Code. Set when the row isn't
// created: it has no name; a group of that name exists, or an earlier row made one; its parent is no
// existing group and no other row (or the rows form a cycle); or it failed on the server. Word the
// message from this.
type GroupImportRowCode string

// Values of GroupImportRowCode. The server may add values; handle an unknown one gracefully.
const (
	GroupImportRowCodeGroupImportNameRequired   GroupImportRowCode = "group.import_name_required"
	GroupImportRowCodeGroupImportNameExists     GroupImportRowCode = "group.import_name_exists"
	GroupImportRowCodeGroupImportParentNotFound GroupImportRowCode = "group.import_parent_not_found"
	GroupImportRowCodeInternal                  GroupImportRowCode = "internal"
)

// GroupImportRowField is a string enumeration used by GroupImportRow.Field. The row field the code
// concerns, when one does.
type GroupImportRowField string

// Values of GroupImportRowField. The server may add values; handle an unknown one gracefully.
const (
	GroupImportRowFieldName   GroupImportRowField = "name"
	GroupImportRowFieldParent GroupImportRowField = "parent"
)

// GroupImportRowStatus is a string enumeration used by GroupImportRow.Status.
type GroupImportRowStatus string

// Values of GroupImportRowStatus. The server may add values; handle an unknown one gracefully.
const (
	GroupImportRowStatusCreated GroupImportRowStatus = "created"
	GroupImportRowStatusSkipped GroupImportRowStatus = "skipped"
	GroupImportRowStatusFailed  GroupImportRowStatus = "failed"
)

// GroupTemplate is a saved tree of groups that GroupsService.ApplyTemplate creates in one step.
type GroupTemplate struct {
	CreatedAt   time.Time           `json:"created_at,omitzero"`
	Description string              `json:"description,omitempty"`
	GroupCount  int                 `json:"group_count"`
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Nodes       []GroupTemplateNode `json:"nodes,omitempty"`
	TenantID    string              `json:"tenant_id,omitempty"`
	UpdatedAt   time.Time           `json:"updated_at,omitzero"`
}

// GroupTemplateNode is one element of GroupTemplate.Nodes.
type GroupTemplateNode struct {
	Children    []GroupTemplateNode `json:"children,omitempty"`
	Description string              `json:"description,omitempty"`
	Name        string              `json:"name,omitempty"`
}

// CreateGroupTemplateParams is the request body of GroupsService.CreateTemplate.
type CreateGroupTemplateParams struct {
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Nodes       []GroupTemplateNode `json:"nodes,omitempty"`
}

// ApplyGroupTemplateParams is the request body of GroupsService.ApplyTemplate.
type ApplyGroupTemplateParams struct {
	ParentID string `json:"parent_id,omitempty"`
	// Put before every group's name, such as "Pune · ".
	Prefix string `json:"prefix,omitempty"`
	// A new name for the template's one top group.
	RootName string `json:"root_name,omitempty"`
}
