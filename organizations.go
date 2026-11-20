package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// OrganizationsService manages organizations. Every other resource belongs to exactly one organization.
type OrganizationsService struct{ c *Client }

// List lists the organizations the caller can see, one page at a time.
//
// API: GET /v1/organizations
func (s *OrganizationsService) List(ctx context.Context, params *ListOrganizationsParams, opts ...RequestOption) (*Page[Organization], error) {
	return call[Page[Organization]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/organizations", query: params.values()}, opts)
}

// All returns an iterator over every item List returns, across all pages. It fetches a page only when
// the loop needs it, and breaking out of the loop stops paging. A failed request ends the iteration
// with a non-nil error.
func (s *OrganizationsService) All(ctx context.Context, params *ListOrganizationsParams, opts ...RequestOption) iter.Seq2[*Organization, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[Organization], error) {
		return s.List(ctx, params.withCursor(cursor), opts...)
	})
}

// Create creates an organization. The calling person becomes its owner; the result carries the new organization and
// a token pair scoped to it.
//
// API: POST /v1/organizations
func (s *OrganizationsService) Create(ctx context.Context, params CreateOrganizationParams, opts ...RequestOption) (*CreateOrganizationResult, error) {
	return call[CreateOrganizationResult](ctx, s.c, &request{method: http.MethodPost, path: "/v1/organizations", body: params}, opts)
}

// Get returns an organization.
//
// API: GET /v1/organizations/{organizationID}
func (s *OrganizationsService) Get(ctx context.Context, organizationID string, opts ...RequestOption) (*Organization, error) {
	return call[Organization](ctx, s.c, &request{method: http.MethodGet, path: "/v1/organizations/{organizationID}", pathArgs: []string{organizationID}}, opts)
}

// Update changes an organization's name, status, logo or metadata. The plan changes through Billing.
//
// API: PATCH /v1/organizations/{organizationID}
func (s *OrganizationsService) Update(ctx context.Context, organizationID string, params UpdateOrganizationParams, opts ...RequestOption) (*Organization, error) {
	return call[Organization](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/organizations/{organizationID}",
		pathArgs: []string{organizationID},
		body:     params,
	}, opts)
}

// Delete schedules an organization for deletion. The organization is deleted when DeletionDueAt passes unless
// CancelDeletion is called first.
//
// API: DELETE /v1/organizations/{organizationID}
func (s *OrganizationsService) Delete(ctx context.Context, organizationID string, opts ...RequestOption) (*OrganizationDeletion, error) {
	return call[OrganizationDeletion](ctx, s.c, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}",
		pathArgs: []string{organizationID},
	}, opts)
}

// CancelDeletion cancels a scheduled organization deletion.
//
// API: POST /v1/organizations/{organizationID}/deletion/cancel
func (s *OrganizationsService) CancelDeletion(ctx context.Context, organizationID string, opts ...RequestOption) (*Organization, error) {
	return call[Organization](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/deletion/cancel",
		pathArgs: []string{organizationID},
	}, opts)
}

// ListOrganizationsParams holds the optional query parameters of OrganizationsService.List. A nil
// *ListOrganizationsParams sends none.
type ListOrganizationsParams struct {
	// Maximum items to return.
	Limit int
	// Opaque pagination cursor from a prior response's next_cursor.
	Cursor string
}

func (p *ListOrganizationsParams) values() url.Values {
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

func (p *ListOrganizationsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListOrganizationsParams) withCursor(cursor string) *ListOrganizationsParams {
	var c ListOrganizationsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// Organization is an organization on Qeet ID. Every user, group, role and setting belongs to exactly one
// organization.
type Organization struct {
	ID     string             `json:"id,omitempty"`
	Slug   string             `json:"slug,omitempty"`
	Name   string             `json:"name,omitempty"`
	Status OrganizationStatus `json:"status,omitempty"`
	Plan   OrganizationPlan   `json:"plan,omitempty"`
	// The region this service is hosted in, where the organization's data is held. It is the same for
	// every organization of a deployment and isn't a per-organization data-residency choice.
	Region string `json:"region,omitempty"`
	// The organization's logo, a hosted URL or a small inline data URL; empty when none is set.
	LogoURL   string         `json:"logo_url,omitempty"`
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

// OrganizationStatus is a string enumeration used by Organization.Status.
type OrganizationStatus string

// Values of OrganizationStatus. The server may add values; handle an unknown one gracefully.
const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
)

// OrganizationPlan is a string enumeration used by Organization.Plan.
type OrganizationPlan string

// Values of OrganizationPlan. The server may add values; handle an unknown one gracefully.
const (
	OrganizationPlanFree       OrganizationPlan = "free"
	OrganizationPlanStarter    OrganizationPlan = "starter"
	OrganizationPlanPro        OrganizationPlan = "pro"
	OrganizationPlanEnterprise OrganizationPlan = "enterprise"
)

// CreateOrganizationParams is the request body of OrganizationsService.Create.
type CreateOrganizationParams struct {
	Slug string `json:"slug,omitempty"`
	Name string `json:"name,omitempty"`
	// Only the free plan is created here (it may be left out); a paid plan answers 402
	// billing.checkout_required, since it starts with checkout.
	Plan OrganizationPlan `json:"plan,omitempty"`
	// Optional. Only the region this service is hosted in is accepted (422 organization.region_unavailable,
	// with metadata.available, otherwise); the organization is created there either way.
	Region string `json:"region,omitempty"`
	// Optional logo, a hosted URL or a small inline data URL, up to 3,000,000 characters.
	LogoURL  string         `json:"logo_url,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// CreateOrganizationResult is the result of OrganizationsService.Create: the new organization and a token pair for its
// owner, scoped to it.
type CreateOrganizationResult struct {
	Organization   Organization `json:"organization,omitzero"`
	OrganizationID string       `json:"organization_id,omitempty"`
	Roles          []string     `json:"roles,omitempty"`
	AccessToken    string       `json:"access_token,omitempty"`
	TokenType      string       `json:"token_type,omitempty"`
	ExpiresAt      time.Time    `json:"expires_at,omitzero"`
	RefreshToken   string       `json:"refresh_token,omitempty"`
	SessionID      string       `json:"session_id,omitempty"`
	UserID         string       `json:"user_id,omitempty"`
}

// UpdateOrganizationParams is the request body of OrganizationsService.Update.
type UpdateOrganizationParams struct {
	Name   *string             `json:"name,omitempty"`
	Status *OrganizationStatus `json:"status,omitempty"`
	// Read-only: the same value is accepted, as a form sends every field back; another is refused with
	// 422 organization.region_read_only.
	Region *string `json:"region,omitempty"`
	// The logo, a hosted URL or a small inline data URL, up to 3,000,000 characters; empty removes it.
	LogoURL  *string        `json:"logo_url,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// OrganizationDeletion is the result of OrganizationsService.Delete.
type OrganizationDeletion struct {
	OrganizationID string    `json:"organization_id,omitempty"`
	DeletionDueAt  time.Time `json:"deletion_due_at,omitzero"`
}
