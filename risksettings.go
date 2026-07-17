package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// RiskSettings holds a tenant's adaptive-MFA risk-assessment thresholds.
type RiskSettings struct {
	MediumThreshold float64 `json:"medium_threshold"`
	HighThreshold   float64 `json:"high_threshold"`
	ForceMFAAtLevel string  `json:"force_mfa_at_level"` // medium | high
	// ImpossibleTravelEnabled flags a login from a different country than the
	// user's last-seen one, sooner than MinTravelHours could plausibly allow.
	// Off by default: it needs a country signal from the caller (an upstream
	// proxy geo header) and, like any new heuristic, shouldn't affect logins
	// until a tenant opts in.
	ImpossibleTravelEnabled bool    `json:"impossible_travel_enabled"`
	MinTravelHours          float64 `json:"min_travel_hours"`
	// DeviceReputationEnabled flags a login from a browser+OS combination never
	// seen before for this user.
	DeviceReputationEnabled bool `json:"device_reputation_enabled"`
}

// UpdateRiskSettingsInput is a full replace, not a partial patch: the backend
// decodes it directly into its settings record (no pointer/optional fields),
// so any omitted field is written back as its zero value — which the backend
// then clamps or defaults (e.g. an omitted MinTravelHours becomes the server's
// default, not "leave unchanged").
type UpdateRiskSettingsInput struct {
	MediumThreshold         float64 `json:"medium_threshold"`
	HighThreshold           float64 `json:"high_threshold"`
	ForceMFAAtLevel         string  `json:"force_mfa_at_level"`
	ImpossibleTravelEnabled bool    `json:"impossible_travel_enabled"`
	MinTravelHours          float64 `json:"min_travel_hours"`
	DeviceReputationEnabled bool    `json:"device_reputation_enabled"`
}

// RiskSettingsService manages tenant-scoped risk-assessment settings — one of
// the threat-detection surfaces alongside BotDetectionService. These
// thresholds drive adaptive MFA: a request scored at or above ForceMFAAtLevel
// is forced through a second factor even if the device is otherwise trusted.
// No pagination, no list — this is a singleton get/put record, the same shape
// as AuthPolicyService/PolicyService.
type RiskSettingsService struct{ t *transport.Transport }

// Get returns the tenant's current risk-assessment settings.
func (s *RiskSettingsService) Get(ctx context.Context, tenantID string) (*RiskSettings, error) {
	var out RiskSettings
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/security/risk-settings", transport.RequestOptions{}, &out)
	return &out, err
}

// Put replaces the tenant's risk-assessment settings (full replace).
func (s *RiskSettingsService) Put(ctx context.Context, tenantID string, in UpdateRiskSettingsInput) (*RiskSettings, error) {
	var out RiskSettings
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/security/risk-settings", in, &out)
	return &out, err
}
