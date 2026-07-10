package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// RetentionPolicy controls automatic purging of soft-deleted users.
type RetentionPolicy struct {
	DeletedUsersEnabled bool `json:"deleted_users_enabled"`
	DeletedUsersDays    int  `json:"deleted_users_days"` // clamped server-side to [1, 3650]
}

// RetentionPreviewResult is a dry-run count — nothing is deleted.
type RetentionPreviewResult struct {
	RipeDeletedUsers int `json:"ripe_deleted_users"`
	DeletedUsersDays int `json:"deleted_users_days"`
}

// RetentionRunResult reports how many records Run actually purged.
type RetentionRunResult struct {
	Purged int `json:"purged"`
}

// RetentionService manages the tenant's data-retention policy.
type RetentionService struct{ t *transport.Transport }

func (s *RetentionService) Get(ctx context.Context, tenantID string) (*RetentionPolicy, error) {
	var out RetentionPolicy
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/retention", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *RetentionService) Put(ctx context.Context, tenantID string, in RetentionPolicy) (*RetentionPolicy, error) {
	var out RetentionPolicy
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/retention", in, &out)
	return &out, err
}

// Preview reports how many records are currently ripe for purge, without
// deleting anything.
func (s *RetentionService) Preview(ctx context.Context, tenantID string) (*RetentionPreviewResult, error) {
	var out RetentionPreviewResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/retention/preview", nil, &out)
	return &out, err
}

// Run purges ripe records immediately, ahead of the scheduled sweep.
func (s *RetentionService) Run(ctx context.Context, tenantID string) (*RetentionRunResult, error) {
	var out RetentionRunResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/retention/run", nil, &out)
	return &out, err
}
