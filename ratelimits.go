package qeetid

import (
	"context"
	"net/http"
)

// RateLimitsService manages a tenant's rate limits.
type RateLimitsService struct{ c *Client }

// Get returns the tenant's effective rate limits, the defaults and the allowed ceiling.
//
// API: GET /v1/tenants/{tenantID}/rate-limits
func (s *RateLimitsService) Get(ctx context.Context, tenantID string, opts ...RequestOption) (*RateLimits, error) {
	return call[RateLimits](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/rate-limits",
		pathArgs: []string{tenantID},
	}, opts)
}

// Update sets the tenant's rate limits for the buckets given.
//
// API: PUT /v1/tenants/{tenantID}/rate-limits
func (s *RateLimitsService) Update(ctx context.Context, tenantID string, params UpdateRateLimitsParams, opts ...RequestOption) (*RateLimits, error) {
	return call[RateLimits](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/rate-limits",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Reset returns the tenant to the default rate limits.
//
// API: DELETE /v1/tenants/{tenantID}/rate-limits
func (s *RateLimitsService) Reset(ctx context.Context, tenantID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/rate-limits",
		pathArgs: []string{tenantID},
	}, nil, opts)
}

// RateLimits holds a tenant's effective rate limits, with the defaults and bounds they must respect.
type RateLimits struct {
	Tenant   RateLimitBucket  `json:"tenant,omitzero"`
	User     RateLimitBucket  `json:"user,omitzero"`
	APIKey   RateLimitBucket  `json:"api_key,omitzero"`
	Defaults RateLimitBuckets `json:"defaults,omitzero"`
	Ceiling  RateLimitBuckets `json:"ceiling,omitzero"`
	// The lowest rate, in requests a second, an override may set.
	MinRate float64 `json:"min_rate"`
	// The platform's fixed limits, the same for every organisation.
	BuiltIn []BuiltInRateLimit `json:"built_in,omitempty"`
}

// RateLimitBucket is a token-bucket limit: Rate requests a second, with bursts of up to Capacity.
type RateLimitBucket struct {
	// Sustained tokens per second.
	Rate float64 `json:"rate"`
	// Maximum burst capacity.
	Capacity int `json:"capacity"`
}

// RateLimitBuckets is the type of RateLimits.Defaults.
type RateLimitBuckets struct {
	Tenant RateLimitBucket `json:"tenant,omitzero"`
	User   RateLimitBucket `json:"user,omitzero"`
	APIKey RateLimitBucket `json:"api_key,omitzero"`
}

// BuiltInRateLimit is one element of RateLimits.BuiltIn.
type BuiltInRateLimit struct {
	Name string `json:"name,omitempty"`
	// Whose requests share a bucket.
	Per RateLimitScope `json:"per,omitempty"`
	// What the limit covers.
	Applies string `json:"applies,omitempty"`
	// Requests a second.
	Rate float64 `json:"rate"`
	// The burst.
	Capacity int `json:"capacity"`
}

// RateLimitScope is a string enumeration used by BuiltInRateLimit.Per. Whose requests share a bucket.
type RateLimitScope string

// Values of RateLimitScope. The server may add values; handle an unknown one gracefully.
const (
	RateLimitScopeIP      RateLimitScope = "ip"
	RateLimitScopeUser    RateLimitScope = "user"
	RateLimitScopeAddress RateLimitScope = "address"
)

// UpdateRateLimitsParams is the request body of RateLimitsService.Update. Per-bucket rate limits for a
// tenant. On GET these are the effective values (defaults merged with overrides); on PUT, any bucket
// supplied is upserted as an override. The answer also carries the defaults, the ceiling an override
// may reach (the defaults themselves: an organisation may tighten its limits, never lift them), the
// lowest rate an override may set, and the built-in limits every organisation shares; those are
// read-only and ignored on PUT.
type UpdateRateLimitsParams struct {
	Tenant *RateLimitBucket `json:"tenant,omitempty"`
	User   *RateLimitBucket `json:"user,omitempty"`
	APIKey *RateLimitBucket `json:"api_key,omitempty"`
}
