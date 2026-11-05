package qeetid

import (
	"context"
	"net/http"
)

// AuthPolicyService manages how a tenant's people may sign in.
type AuthPolicyService struct{ c *Client }

// Get returns a tenant's sign-in policy.
//
// API: GET /v1/tenants/{tenantID}/auth-policy
func (s *AuthPolicyService) Get(ctx context.Context, tenantID string, opts ...RequestOption) (*AuthPolicy, error) {
	return call[AuthPolicy](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/auth-policy",
		pathArgs: []string{tenantID},
	}, opts)
}

// Update replaces a tenant's sign-in policy. Read it with Get, change it, and send it back.
//
// API: PUT /v1/tenants/{tenantID}/auth-policy
func (s *AuthPolicyService) Update(ctx context.Context, tenantID string, params AuthPolicy, opts ...RequestOption) (*AuthPolicy, error) {
	return call[AuthPolicy](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/auth-policy",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Usage reports how the tenant's people sign in today, to show the effect of a policy change.
//
// API: GET /v1/tenants/{tenantID}/auth-policy/usage
func (s *AuthPolicyService) Usage(ctx context.Context, tenantID string, opts ...RequestOption) (*AuthPolicyUsage, error) {
	return call[AuthPolicyUsage](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/auth-policy/usage",
		pathArgs: []string{tenantID},
	}, opts)
}

// AuthPolicy is a tenant's sign-in policy: the methods people may use, password rules, lockout and
// remembered devices.
type AuthPolicy struct {
	PasswordEnabled          bool `json:"password_enabled"`
	PasswordMinLength        int  `json:"password_min_length"`
	PasswordRequireUppercase bool `json:"password_require_uppercase"`
	PasswordRequireNumber    bool `json:"password_require_number"`
	PasswordRequireSymbol    bool `json:"password_require_symbol"`
	MagicLinkEnabled         bool `json:"magic_link_enabled"`
	MagicLinkTTLMinutes      int  `json:"magic_link_ttl_minutes"`
	PasskeyEnabled           bool `json:"passkey_enabled"`
	OTPEmailEnabled          bool `json:"otp_email_enabled"`
	OTPSMSEnabled            bool `json:"otp_sms_enabled"`
	SelfRegistrationEnabled  bool `json:"self_registration_enabled"`
	RememberDeviceEnabled    bool `json:"remember_device_enabled"`
	// How long a device stays remembered after the sign-in that remembered it, in days (30 by default; 0
	// on update keeps the default). Applies to devices remembered after it is set.
	RememberDeviceDays int `json:"remember_device_days"`
	// Failed password sign-ins to one account from one network, each within the window of the last, that
	// lock the account for that network. null = the platform default. Must lie within lockout_limits
	// (default bounds 3..20).
	LockoutThreshold *int `json:"lockout_threshold,omitempty"`
	// A failure is forgotten once this many minutes pass with no new failure. null = the platform
	// default. Must lie within lockout_limits (default bounds 1..1440).
	LockoutWindowMinutes *int `json:"lockout_window_minutes,omitempty"`
	// How long a lock lasts. null = the platform default. Must lie within lockout_limits (default bounds
	// 1..1440).
	LockoutDurationMinutes *int          `json:"lockout_duration_minutes,omitempty"`
	LockoutLimits          LockoutLimits `json:"lockout_limits,omitzero"`
}

// LockoutLimits is the platform defaults for the lockout settings (LOGIN_LOCKOUT_THRESHOLD / _WINDOW /
// _DURATION, built-in 5 / 15 / 15) and the inclusive bounds a tenant may choose within
// (LOGIN_LOCKOUT_*_MIN / _MAX). Ignored on update.
type LockoutLimits struct {
	Defaults LockoutSettings `json:"defaults,omitzero"`
	Min      LockoutSettings `json:"min,omitzero"`
	Max      LockoutSettings `json:"max,omitzero"`
}

// LockoutSettings is the type of LockoutLimits.Defaults.
type LockoutSettings struct {
	Threshold       int `json:"threshold"`
	WindowMinutes   int `json:"window_minutes"`
	DurationMinutes int `json:"duration_minutes"`
}

// AuthPolicyUsage is the result of AuthPolicyService.Usage.
type AuthPolicyUsage struct {
	People           int                     `json:"people"`
	Password         AuthPolicyUsagePassword `json:"password,omitzero"`
	Passkey          AuthPolicyUsagePasskey  `json:"passkey,omitzero"`
	Social           []SocialSignInUsage     `json:"social,omitempty"`
	MagicLinkEnabled bool                    `json:"magic_link_enabled"`
}

// AuthPolicyUsagePassword is the type of AuthPolicyUsage.Password.
type AuthPolicyUsagePassword struct {
	People int `json:"people"`
	Only   int `json:"only"`
}

// AuthPolicyUsagePasskey is the type of AuthPolicyUsage.Passkey.
type AuthPolicyUsagePasskey struct {
	People int `json:"people"`
	Only   int `json:"only"`
}

// SocialSignInUsage is one element of AuthPolicyUsage.Social.
type SocialSignInUsage struct {
	Provider string `json:"provider,omitempty"`
	Linked   int    `json:"linked"`
	Only     int    `json:"only"`
}
