package qeetid

import (
	"context"
	"net/http"
	"time"
)

// SecurityPolicyService manages a tenant's session, password and second-factor policy.
type SecurityPolicyService struct{ c *Client }

// Get returns a tenant's security policy.
//
// API: GET /v1/tenants/{tenantID}/policy
func (s *SecurityPolicyService) Get(ctx context.Context, tenantID string, opts ...RequestOption) (*SecurityPolicy, error) {
	return call[SecurityPolicy](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/policy",
		pathArgs: []string{tenantID},
	}, opts)
}

// Replace replaces a tenant's whole security policy.
//
// API: PUT /v1/tenants/{tenantID}/policy
func (s *SecurityPolicyService) Replace(ctx context.Context, tenantID string, params SecurityPolicy, opts ...RequestOption) (*SecurityPolicy, error) {
	return call[SecurityPolicy](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/policy",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Update changes part of a tenant's security policy. Fields left nil are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/policy
func (s *SecurityPolicyService) Update(ctx context.Context, tenantID string, params UpdateSecurityPolicyParams, opts ...RequestOption) (*SecurityPolicy, error) {
	return call[SecurityPolicy](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/policy",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// SecurityPolicy is a tenant's session, password and second-factor policy.
type SecurityPolicy struct {
	TenantID string `json:"tenant_id,omitempty"`
	// Accepted and ignored, and never returned. Password rules are the auth policy's (AuthPolicy, PUT
	// /v1/tenants/{tenantID}/auth-policy), which sign-in and the compliance checks apply.
	PasswordMinLength int `json:"password_min_length"`
	// Accepted and ignored, and never returned; see password_min_length.
	PasswordComplexity string `json:"password_complexity,omitempty"`
	// How long a session lasts after sign-in however active it is, and an application's grant however
	// often it is renewed: a duration such as "720h" or "90m", between 1h and 8760h. Applies to sessions
	// and grants started after it is saved. Omitted or zero on update: the platform default.
	SessionMaxAge string `json:"session_max_age,omitempty"`
	// How long a session may go without being renewed before it ends: a duration of at least 15m and no
	// longer than session_max_age. Applies to sessions started after it is saved. Omitted or zero on
	// update: the platform default (no longer than session_max_age).
	SessionIdleTimeout string `json:"session_idle_timeout,omitempty"`
	// How many live sessions one person may hold in the organisation at once; null for no limit (the
	// default). Counted and enforced as a session opens, so simultaneous sign-ins cannot exceed it.
	MaxSessions *int `json:"max_sessions,omitempty"`
	// What a sign-in at max_sessions does: evict_oldest ends the person's least recently used session
	// (renewal stops at once; an access token it already issued lasts out its own lifetime,
	// ACCESS_TOKEN_TTL), announced to session.revoked subscribers as session_limit; refuse refuses the
	// sign-in with 403 auth.session_limit. Both are audited.
	SessionLimitAction SecurityPolicySessionLimitAction `json:"session_limit_action,omitempty"`
	// Who must use a second factor to sign in, enforced at password sign-in: optional (only people who
	// set one up are challenged), required (everyone), admin_only (anyone holding a role, directly or
	// through a group, that grants a write permission). Separately, an account an administrator marks as
	// requiring MFA must, whatever is set here. An account that must and holds no factor sets up an
	// authenticator app before any session is issued.
	MFAEnforcement SecurityPolicyMFAEnforcement `json:"mfa_enforcement,omitempty"`
	// For this many days after the requirement is turned on, the accounts that already existed may put
	// setting up a second factor off (they are prompted at each sign-in). 0, the default: no grace
	// period.
	MFAGraceDays int `json:"mfa_grace_days"`
	// When the requirement was turned on (from optional); null while optional. Moving between required
	// and admin_only keeps it.
	MFAEnforcedSince *time.Time `json:"mfa_enforced_since,omitempty"`
	// How long a step-up keeps a session's sensitive actions unlocked: a duration between 1m and 60m. A
	// step-up belongs to the session that earned it; it ends with that session, and a changed password or
	// an administrative factor reset clears it. Omitted or zero on update: 5m.
	StepUpWindow string `json:"step_up_window,omitempty"`
	// The factor types that count as a step-up (otp: an emailed or text-message code; webauthn: a
	// passkey). Null or empty: every one. A session stepped up by a factor that is not listed is asked
	// again.
	StepUpMethods []StepUpMethod `json:"step_up_methods,omitempty"`
	Settings      map[string]any `json:"settings,omitempty"`
}

// SecurityPolicySessionLimitAction is a string enumeration used by SecurityPolicy.SessionLimitAction.
// What a sign-in at max_sessions does: evict_oldest ends the person's least recently used session
// (renewal stops at once; an access token it already issued lasts out its own lifetime,
// ACCESS_TOKEN_TTL), announced to session.revoked subscribers as session_limit; refuse refuses the
// sign-in with 403 auth.session_limit. Both are audited.
type SecurityPolicySessionLimitAction string

// Values of SecurityPolicySessionLimitAction. The server may add values; handle an unknown one gracefully.
const (
	SecurityPolicySessionLimitActionEvictOldest SecurityPolicySessionLimitAction = "evict_oldest"
	SecurityPolicySessionLimitActionRefuse      SecurityPolicySessionLimitAction = "refuse"
)

// SecurityPolicyMFAEnforcement is a string enumeration used by SecurityPolicy.MFAEnforcement. Who must
// use a second factor to sign in, enforced at password sign-in: optional (only people who set one up
// are challenged), required (everyone), admin_only (anyone holding a role, directly or through a
// group, that grants a write permission). Separately, an account an administrator marks as requiring
// MFA must, whatever is set here. An account that must and holds no factor sets up an authenticator
// app before any session is issued.
type SecurityPolicyMFAEnforcement string

// Values of SecurityPolicyMFAEnforcement. The server may add values; handle an unknown one gracefully.
const (
	SecurityPolicyMFAEnforcementOptional  SecurityPolicyMFAEnforcement = "optional"
	SecurityPolicyMFAEnforcementRequired  SecurityPolicyMFAEnforcement = "required"
	SecurityPolicyMFAEnforcementAdminOnly SecurityPolicyMFAEnforcement = "admin_only"
)

// StepUpMethod is a string enumeration used by SecurityPolicy.StepUpMethods.
type StepUpMethod string

// Values of StepUpMethod. The server may add values; handle an unknown one gracefully.
const (
	StepUpMethodTOTP         StepUpMethod = "totp"
	StepUpMethodRecoveryCode StepUpMethod = "recovery_code"
	StepUpMethodOTP          StepUpMethod = "otp"
	StepUpMethodWebAuthn     StepUpMethod = "webauthn"
	StepUpMethodPush         StepUpMethod = "push"
)

// UpdateSecurityPolicyParams is the request body of SecurityPolicyService.Update.
type UpdateSecurityPolicyParams struct {
	// How long a session lasts after sign-in however active it is, and an application's grant however
	// often it is renewed: a duration such as "720h" or "90m", between 1h and 8760h. Applies to sessions
	// and grants started after it is saved. Omitted or zero on update: the platform default.
	SessionMaxAge *string `json:"session_max_age,omitempty"`
	// How long a session may go without being renewed before it ends: a duration of at least 15m and no
	// longer than session_max_age. Applies to sessions started after it is saved. Omitted or zero on
	// update: the platform default (no longer than session_max_age).
	SessionIdleTimeout *string `json:"session_idle_timeout,omitempty"`
	// How many live sessions one person may hold in the organisation at once; null for no limit (the
	// default). Counted and enforced as a session opens, so simultaneous sign-ins cannot exceed it.
	MaxSessions *int `json:"max_sessions,omitempty"`
	// What a sign-in at max_sessions does: evict_oldest ends the person's least recently used session
	// (renewal stops at once; an access token it already issued lasts out its own lifetime,
	// ACCESS_TOKEN_TTL), announced to session.revoked subscribers as session_limit; refuse refuses the
	// sign-in with 403 auth.session_limit. Both are audited.
	SessionLimitAction *SecurityPolicySessionLimitAction `json:"session_limit_action,omitempty"`
	// Who must use a second factor to sign in, enforced at password sign-in: optional (only people who
	// set one up are challenged), required (everyone), admin_only (anyone holding a role, directly or
	// through a group, that grants a write permission). Separately, an account an administrator marks as
	// requiring MFA must, whatever is set here. An account that must and holds no factor sets up an
	// authenticator app before any session is issued.
	MFAEnforcement *SecurityPolicyMFAEnforcement `json:"mfa_enforcement,omitempty"`
	// For this many days after the requirement is turned on, the accounts that already existed may put
	// setting up a second factor off (they are prompted at each sign-in). 0, the default: no grace
	// period.
	MFAGraceDays *int `json:"mfa_grace_days,omitempty"`
	// How long a step-up keeps a session's sensitive actions unlocked: a duration between 1m and 60m. A
	// step-up belongs to the session that earned it; it ends with that session, and a changed password or
	// an administrative factor reset clears it. Omitted or zero on update: 5m.
	StepUpWindow *string `json:"step_up_window,omitempty"`
	// The factor types that count as a step-up (otp: an emailed or text-message code; webauthn: a
	// passkey). Null or empty: every one. A session stepped up by a factor that is not listed is asked
	// again.
	StepUpMethods []StepUpMethod `json:"step_up_methods,omitempty"`
	Settings      map[string]any `json:"settings,omitempty"`
}
