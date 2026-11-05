package qeetid

import (
	"context"
	"net/http"
	"time"
)

// RetentionService manages how long a tenant keeps soft-deleted users.
type RetentionService struct{ c *Client }

// Get returns the tenant's retention policy.
//
// API: GET /v1/tenants/{tenantID}/retention
func (s *RetentionService) Get(ctx context.Context, tenantID string, opts ...RequestOption) (*RetentionPolicy, error) {
	return call[RetentionPolicy](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/retention",
		pathArgs: []string{tenantID},
	}, opts)
}

// Update replaces the tenant's retention policy.
//
// API: PUT /v1/tenants/{tenantID}/retention
func (s *RetentionService) Update(ctx context.Context, tenantID string, params RetentionPolicy, opts ...RequestOption) (*RetentionPolicy, error) {
	return call[RetentionPolicy](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/retention",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Preview counts the users the policy would purge now, without deleting anything.
//
// API: POST /v1/tenants/{tenantID}/retention/preview
func (s *RetentionService) Preview(ctx context.Context, tenantID string, opts ...RequestOption) (*RetentionPreview, error) {
	return call[RetentionPreview](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/retention/preview",
		pathArgs: []string{tenantID},
	}, opts)
}

// Run purges the users the policy makes ripe now, ahead of the scheduled sweep.
//
// API: POST /v1/tenants/{tenantID}/retention/run
func (s *RetentionService) Run(ctx context.Context, tenantID string, opts ...RequestOption) (*RetentionRun, error) {
	return call[RetentionRun](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/retention/run",
		pathArgs: []string{tenantID},
	}, opts)
}

// RetentionPolicy sets how long a tenant keeps soft-deleted users before purging them.
type RetentionPolicy struct {
	DeletedUsersDays    int  `json:"deleted_users_days"`
	DeletedUsersEnabled bool `json:"deleted_users_enabled"`
	// When the policy was last changed; absent until it is set.
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	// Who last changed it.
	UpdatedBy      string `json:"updated_by,omitempty"`
	UpdatedByEmail string `json:"updated_by_email,omitempty"`
}

// RetentionPreview is the result of RetentionService.Preview.
type RetentionPreview struct {
	DeletedUsersDays int `json:"deleted_users_days"`
	RipeDeletedUsers int `json:"ripe_deleted_users"`
}

// RetentionRun is the result of RetentionService.Run.
type RetentionRun struct {
	Purged int `json:"purged"`
}
