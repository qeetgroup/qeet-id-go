package qeetid

import (
	"context"
	"net/http"
	"time"
)

// AuthHooksService manages hooks: HTTPS endpoints Qeet ID calls during sign-in to allow, deny or
// enrich it.
//
// Qeet ID signs each call to a hook as it signs a webhook delivery: X-Qeet-Timestamp and an
// X-Qeet-Signature over the timestamp, a "." and the body, with two signatures for 24 hours after the
// hook's secret changes. Verify a call in the hook's endpoint with webhook.Verifier.Verify, which
// refuses a stale timestamp and accepts either signature.
type AuthHooksService struct{ c *Client }

// List lists an organization's hooks.
//
// API: GET /v1/organizations/{organizationID}/auth-hooks
func (s *AuthHooksService) List(ctx context.Context, organizationID string, opts ...RequestOption) (*Page[AuthHook], error) {
	return call[Page[AuthHook]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/auth-hooks",
		pathArgs: []string{organizationID},
	}, opts)
}

// Create creates a hook.
//
// API: POST /v1/organizations/{organizationID}/auth-hooks
func (s *AuthHooksService) Create(ctx context.Context, organizationID string, params CreateAuthHookParams, opts ...RequestOption) (*AuthHook, error) {
	return call[AuthHook](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/auth-hooks",
		pathArgs: []string{organizationID},
		body:     params,
	}, opts)
}

// Update changes a hook. Fields left nil are not changed.
//
// API: PATCH /v1/organizations/{organizationID}/auth-hooks/{id}
func (s *AuthHooksService) Update(ctx context.Context, organizationID, hookID string, params UpdateAuthHookParams, opts ...RequestOption) (*AuthHookState, error) {
	return call[AuthHookState](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/organizations/{organizationID}/auth-hooks/{hookID}",
		pathArgs: []string{organizationID, hookID},
		body:     params,
	}, opts)
}

// Delete deletes a hook.
//
// API: DELETE /v1/organizations/{organizationID}/auth-hooks/{id}
func (s *AuthHooksService) Delete(ctx context.Context, organizationID, hookID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}/auth-hooks/{hookID}",
		pathArgs: []string{organizationID, hookID},
	}, nil, opts)
}

// Test calls a hook with a sample request and reports its answer.
//
// API: POST /v1/organizations/{organizationID}/auth-hooks/{id}/test
func (s *AuthHooksService) Test(ctx context.Context, organizationID, hookID string, opts ...RequestOption) (*AuthHookTest, error) {
	return call[AuthHookTest](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/auth-hooks/{hookID}/test",
		pathArgs: []string{organizationID, hookID},
	}, opts)
}

// AuthHook is an HTTPS endpoint Qeet ID calls during sign-in to allow, deny or add claims to it.
type AuthHook struct {
	CreatedAt time.Time `json:"created_at,omitzero"`
	Enabled   bool      `json:"enabled"`
	FailOpen  bool      `json:"fail_open"`
	ID        string    `json:"id,omitempty"`
	Trigger   string    `json:"trigger,omitempty"`
	URL       string    `json:"url,omitempty"`
	// Its last call by a sign-in.
	LastCalledAt *time.Time `json:"last_called_at,omitempty"`
	// How that call ended; skipped when the hook gave no answer.
	LastOutcome    AuthHookLastOutcome `json:"last_outcome,omitempty"`
	LastStatusCode *int                `json:"last_status_code,omitempty"`
	// Sign-in calls since this release; test calls aren't counted.
	Calls  int `json:"calls"`
	Skips  int `json:"skips"`
	Denies int `json:"denies"`
	// Until when the secret the last change of secret replaced also signs the hook's calls; nil when
	// only one secret signs them.
	PreviousSecretExpiresAt *time.Time `json:"previous_secret_expires_at,omitempty"`
}

// AuthHookLastOutcome is a string enumeration used by AuthHook.LastOutcome. How that call ended;
// skipped when the hook gave no answer.
type AuthHookLastOutcome string

// Values of AuthHookLastOutcome. The server may add values; handle an unknown one gracefully.
const (
	AuthHookLastOutcomeAllowed AuthHookLastOutcome = "allowed"
	AuthHookLastOutcomeDenied  AuthHookLastOutcome = "denied"
	AuthHookLastOutcomeSkipped AuthHookLastOutcome = "skipped"
)

// CreateAuthHookParams is the request body of AuthHooksService.Create.
type CreateAuthHookParams struct {
	URL    string `json:"url,omitempty"`
	Secret string `json:"secret,omitempty"`
	// What a sign-in does when the hook gives no answer: it can't be read, its signing secret can't be
	// opened, the call fails or times out, or its answer can't be read. True (the default on create) lets
	// the sign-in proceed; false refuses it as temporarily unavailable. Where the hook can't be read, the
	// setting last read for the organization is followed. Every such skip is recorded (audit
	// auth_hook.skipped with the reason and outcome; metric auth_hook_skipped_total).
	FailOpen *bool `json:"fail_open,omitempty"`
}

// UpdateAuthHookParams is the request body of AuthHooksService.Update.
type UpdateAuthHookParams struct {
	URL *string `json:"url,omitempty"`
	// A new signing secret; write-only.
	Secret  *string `json:"secret,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
	// What a sign-in does when the hook gives no answer: it can't be read, its signing secret can't be
	// opened, the call fails or times out, or its answer can't be read. True (the default on create) lets
	// the sign-in proceed; false refuses it as temporarily unavailable. Where the hook can't be read, the
	// setting last read for the organization is followed. Every such skip is recorded (audit
	// auth_hook.skipped with the reason and outcome; metric auth_hook_skipped_total).
	FailOpen *bool `json:"fail_open,omitempty"`
}

// AuthHookState is the result of AuthHooksService.Update.
type AuthHookState struct {
	Enabled  bool `json:"enabled"`
	FailOpen bool `json:"fail_open"`
}

// AuthHookTest is the result of AuthHooksService.Test.
type AuthHookTest struct {
	StatusCode int                  `json:"status_code"`
	DurationMs int                  `json:"duration_ms"`
	Decision   AuthHookTestDecision `json:"decision,omitempty"`
	Message    string               `json:"message,omitempty"`
	Claims     map[string]any       `json:"claims,omitempty"`
	Error      string               `json:"error,omitempty"`
}

// AuthHookTestDecision is a string enumeration used by AuthHookTest.Decision.
type AuthHookTestDecision string

// Values of AuthHookTestDecision. The server may add values; handle an unknown one gracefully.
const (
	AuthHookTestDecisionAllow AuthHookTestDecision = "allow"
	AuthHookTestDecisionDeny  AuthHookTestDecision = "deny"
	AuthHookTestDecisionNone  AuthHookTestDecision = "none"
)
