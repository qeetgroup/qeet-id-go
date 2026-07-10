package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// RateLimitBucket is a token-bucket config: rate is tokens/sec, capacity is
// the burst size.
type RateLimitBucket struct {
	Rate     float64 `json:"rate"`
	Capacity int     `json:"capacity"`
}

// TenantRateLimits is the effective (defaults merged with overrides) rate
// limit config for a tenant, per bucket scope.
type TenantRateLimits struct {
	Tenant RateLimitBucket `json:"tenant"`
	User   RateLimitBucket `json:"user"`
	APIKey RateLimitBucket `json:"api_key"`
}

// PutRateLimitsInput upserts whichever buckets are supplied — omitted
// buckets keep their current value.
type PutRateLimitsInput struct {
	Tenant *RateLimitBucket `json:"tenant,omitempty"`
	User   *RateLimitBucket `json:"user,omitempty"`
	APIKey *RateLimitBucket `json:"api_key,omitempty"`
}

// RateLimitsService manages per-tenant rate-limit overrides.
type RateLimitsService struct{ t *transport.Transport }

func (s *RateLimitsService) Get(ctx context.Context, tenantID string) (*TenantRateLimits, error) {
	var out TenantRateLimits
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/rate-limits", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *RateLimitsService) Put(ctx context.Context, tenantID string, in PutRateLimitsInput) (*TenantRateLimits, error) {
	var out TenantRateLimits
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/rate-limits", in, &out)
	return &out, err
}

// Reset clears all overrides, reverting the tenant to platform defaults.
func (s *RateLimitsService) Reset(ctx context.Context, tenantID string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/rate-limits", nil)
}
