package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TenantsService manages tenants (organizations). Every other resource belongs to exactly one tenant.
type TenantsService struct{ c *Client }

// List lists the tenants the caller can see, one page at a time.
//
// API: GET /v1/tenants
func (s *TenantsService) List(ctx context.Context, params *ListTenantsParams, opts ...RequestOption) (*Page[Tenant], error) {
	return call[Page[Tenant]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/tenants", query: params.values()}, opts)
}

// All returns an iterator over every item List returns, across all pages. It fetches a page only when
// the loop needs it, and breaking out of the loop stops paging. A failed request ends the iteration
// with a non-nil error.
func (s *TenantsService) All(ctx context.Context, params *ListTenantsParams, opts ...RequestOption) iter.Seq2[*Tenant, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[Tenant], error) {
		return s.List(ctx, params.withCursor(cursor), opts...)
	})
}

// Create creates a tenant. The calling person becomes its owner; the result carries the new tenant and
// a token pair scoped to it.
//
// API: POST /v1/tenants
func (s *TenantsService) Create(ctx context.Context, params CreateTenantParams, opts ...RequestOption) (*CreateTenantResult, error) {
	return call[CreateTenantResult](ctx, s.c, &request{method: http.MethodPost, path: "/v1/tenants", body: params}, opts)
}

// Get returns a tenant.
//
// API: GET /v1/tenants/{tenantID}
func (s *TenantsService) Get(ctx context.Context, tenantID string, opts ...RequestOption) (*Tenant, error) {
	return call[Tenant](ctx, s.c, &request{method: http.MethodGet, path: "/v1/tenants/{tenantID}", pathArgs: []string{tenantID}}, opts)
}

// Update changes a tenant's name, status, plan, region or metadata.
//
// API: PATCH /v1/tenants/{tenantID}
func (s *TenantsService) Update(ctx context.Context, tenantID string, params UpdateTenantParams, opts ...RequestOption) (*Tenant, error) {
	return call[Tenant](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Delete schedules a tenant for deletion. The tenant is deleted when DeletionDueAt passes unless
// CancelDeletion is called first.
//
// API: DELETE /v1/tenants/{tenantID}
func (s *TenantsService) Delete(ctx context.Context, tenantID string, opts ...RequestOption) (*TenantDeletion, error) {
	return call[TenantDeletion](ctx, s.c, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}",
		pathArgs: []string{tenantID},
	}, opts)
}

// CancelDeletion cancels a scheduled tenant deletion.
//
// API: POST /v1/tenants/{tenantID}/deletion/cancel
func (s *TenantsService) CancelDeletion(ctx context.Context, tenantID string, opts ...RequestOption) (*Tenant, error) {
	return call[Tenant](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/deletion/cancel",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListTenantsParams holds the optional query parameters of TenantsService.List. A nil
// *ListTenantsParams sends none.
type ListTenantsParams struct {
	// Maximum items to return.
	Limit int
	// Opaque pagination cursor from a prior response's next_cursor.
	Cursor string
}

func (p *ListTenantsParams) values() url.Values {
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

func (p *ListTenantsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListTenantsParams) withCursor(cursor string) *ListTenantsParams {
	var c ListTenantsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// Tenant is an organization on Qeet ID. Every user, group, role and setting belongs to exactly one
// tenant.
type Tenant struct {
	ID     string       `json:"id,omitempty"`
	Slug   string       `json:"slug,omitempty"`
	Name   string       `json:"name,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
	Plan   TenantPlan   `json:"plan,omitempty"`
	// The region this service is hosted in, where the organisation's data is held. It is the same for
	// every organisation of a deployment and isn't a per-organisation data-residency choice.
	Region    string         `json:"region,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitzero"`
	UpdatedAt time.Time      `json:"updated_at,omitzero"`
	// When a scheduled deletion takes effect; null when none is.
	DeletionDueAt *time.Time `json:"deletion_due_at,omitempty"`
	// List-only. Distinct members (rbac.user_roles holders) in this org.
	MemberCount int `json:"member_count"`
	// List-only. Members with at least one MFA factor enrolled.
	MFAEnabledCount int `json:"mfa_enabled_count"`
}

// TenantStatus is a string enumeration used by Tenant.Status.
type TenantStatus string

// Values of TenantStatus. The server may add values; handle an unknown one gracefully.
const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
)

// TenantPlan is a string enumeration used by Tenant.Plan.
type TenantPlan string

// Values of TenantPlan. The server may add values; handle an unknown one gracefully.
const (
	TenantPlanFree       TenantPlan = "free"
	TenantPlanPro        TenantPlan = "pro"
	TenantPlanEnterprise TenantPlan = "enterprise"
)

// CreateTenantParams is the request body of TenantsService.Create.
type CreateTenantParams struct {
	Slug string     `json:"slug,omitempty"`
	Name string     `json:"name,omitempty"`
	Plan TenantPlan `json:"plan,omitempty"`
	// Optional. Only the region this service is hosted in is accepted (422 tenant.region_unavailable,
	// with metadata.available, otherwise); the organisation is created there either way.
	Region   string         `json:"region,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// CreateTenantResult is the result of TenantsService.Create: the new tenant and a token pair for its
// owner, scoped to it.
type CreateTenantResult struct {
	Tenant       Tenant    `json:"tenant,omitzero"`
	TenantID     string    `json:"tenant_id,omitempty"`
	Roles        []string  `json:"roles,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitzero"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	SessionID    string    `json:"session_id,omitempty"`
	UserID       string    `json:"user_id,omitempty"`
}

// UpdateTenantParams is the request body of TenantsService.Update.
type UpdateTenantParams struct {
	Name   *string       `json:"name,omitempty"`
	Status *TenantStatus `json:"status,omitempty"`
	Plan   *TenantPlan   `json:"plan,omitempty"`
	// Read-only: the same value is accepted, as a form sends every field back; another is refused with
	// 422 tenant.region_read_only.
	Region   *string        `json:"region,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// TenantDeletion is the result of TenantsService.Delete.
type TenantDeletion struct {
	TenantID      string    `json:"tenant_id,omitempty"`
	DeletionDueAt time.Time `json:"deletion_due_at,omitzero"`
}
