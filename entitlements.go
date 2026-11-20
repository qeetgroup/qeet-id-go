package qeetid

import (
	"context"
	"net/http"
)

// EntitlementsService reports what an organization's plan includes and how much of it is used.
type EntitlementsService struct{ c *Client }

// Get returns the features and limits of the organization's plan.
//
// API: GET /v1/organizations/{organizationID}/entitlements
func (s *EntitlementsService) Get(ctx context.Context, organizationID string, opts ...RequestOption) (*Entitlements, error) {
	return call[Entitlements](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/entitlements",
		pathArgs: []string{organizationID},
	}, opts)
}

// Usage returns the organization's resource usage against its plan limits.
//
// API: GET /v1/organizations/{organizationID}/entitlements/usage
func (s *EntitlementsService) Usage(ctx context.Context, organizationID string, opts ...RequestOption) (*EntitlementUsage, error) {
	return call[EntitlementUsage](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/entitlements/usage",
		pathArgs: []string{organizationID},
	}, opts)
}

// Entitlements lists the features and limits of an organization's plan.
type Entitlements struct {
	// Effective plan tier (annual "_year" variants folded to the base tier).
	Plan PlanCode `json:"plan,omitempty"`
	// Boolean feature flags keyed by feature name.
	Features map[string]bool `json:"features,omitempty"`
	// Numeric resource caps keyed by resource name; -1 means unlimited.
	Limits map[string]int `json:"limits,omitempty"`
}

// PlanCode is a string enumeration used by Entitlements.Plan. Effective plan tier (annual "_year"
// variants folded to the base tier).
type PlanCode string

// Values of PlanCode. The server may add values; handle an unknown one gracefully.
const (
	PlanCodeFree       PlanCode = "free"
	PlanCodeStarter    PlanCode = "starter"
	PlanCodePro        PlanCode = "pro"
	PlanCodeEnterprise PlanCode = "enterprise"
)

// EntitlementUsage is the result of EntitlementsService.Usage.
type EntitlementUsage struct {
	Usage map[string]int `json:"usage,omitempty"`
}
