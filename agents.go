package qeetid

import (
	"context"
	"net/http"
	"time"
)

// AgentsService manages AI agent identities, their sponsors and their lifecycle.
type AgentsService struct{ c *Client }

// Create registers an agent. The result carries its secret, shown only once.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/tenants/{tenantID}/agents
func (s *AgentsService) Create(ctx context.Context, tenantID string, params CreateAgentParams, opts ...RequestOption) (*Agent, error) {
	return call[Agent](ctx, s.c, &request{
		method:     http.MethodPost,
		path:       "/v1/tenants/{tenantID}/agents",
		pathArgs:   []string{tenantID},
		body:       params,
		idempotent: true,
	}, opts)
}

// List lists a tenant's agents. Page.Truncated reports a list cut short by the server's cap.
//
// API: GET /v1/tenants/{tenantID}/agents
func (s *AgentsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Agent], error) {
	return call[Page[Agent]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/agents",
		pathArgs: []string{tenantID},
	}, opts)
}

// Update enables or disables an agent.
//
// API: PATCH /v1/tenants/{tenantID}/agents/{id}
func (s *AgentsService) Update(ctx context.Context, tenantID string, agentID string, params UpdateAgentParams, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/agents/{agentID}",
		pathArgs: []string{tenantID, agentID},
		body:     params,
	}, nil, opts)
}

// Delete deletes an agent.
//
// API: DELETE /v1/tenants/{tenantID}/agents/{id}
func (s *AgentsService) Delete(ctx context.Context, tenantID string, agentID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/agents/{agentID}",
		pathArgs: []string{tenantID, agentID},
	}, nil, opts)
}

// Suspend suspends an agent. Its tokens stop working until Resume.
//
// API: POST /v1/tenants/{tenantID}/agents/{id}/suspend
func (s *AgentsService) Suspend(ctx context.Context, tenantID string, agentID string, opts ...RequestOption) (*AgentStatusResult, error) {
	return call[AgentStatusResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/{agentID}/suspend",
		pathArgs: []string{tenantID, agentID},
	}, opts)
}

// Resume makes a suspended agent active again.
//
// API: POST /v1/tenants/{tenantID}/agents/{id}/resume
func (s *AgentsService) Resume(ctx context.Context, tenantID string, agentID string, opts ...RequestOption) (*AgentStatusResult, error) {
	return call[AgentStatusResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/{agentID}/resume",
		pathArgs: []string{tenantID, agentID},
	}, opts)
}

// Decommission decommissions an agent for good. It cannot be resumed.
//
// API: POST /v1/tenants/{tenantID}/agents/{id}/decommission
func (s *AgentsService) Decommission(ctx context.Context, tenantID string, agentID string, opts ...RequestOption) (*AgentStatusResult, error) {
	return call[AgentStatusResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/{agentID}/decommission",
		pathArgs: []string{tenantID, agentID},
	}, opts)
}

// SuspendAll suspends every active agent of the tenant at once (the kill switch).
//
// API: POST /v1/tenants/{tenantID}/agents/kill-all
func (s *AgentsService) SuspendAll(ctx context.Context, tenantID string, opts ...RequestOption) (*SuspendedAgents, error) {
	return call[SuspendedAgents](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/kill-all",
		pathArgs: []string{tenantID},
	}, opts)
}

// ResumeAll makes every suspended agent of the tenant active again.
//
// API: POST /v1/tenants/{tenantID}/agents/resume-all
func (s *AgentsService) ResumeAll(ctx context.Context, tenantID string, opts ...RequestOption) (*ResumedAgents, error) {
	return call[ResumedAgents](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/resume-all",
		pathArgs: []string{tenantID},
	}, opts)
}

// RotateSecret replaces an agent's secret. The old secret keeps working for the overlap given.
//
// API: POST /v1/tenants/{tenantID}/agents/{id}/rotate-secret
func (s *AgentsService) RotateSecret(ctx context.Context, tenantID string, agentID string, params RotateSecretParams, opts ...RequestOption) (*RotatedAgentSecret, error) {
	return call[RotatedAgentSecret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/{agentID}/rotate-secret",
		pathArgs: []string{tenantID, agentID},
		body:     params,
	}, opts)
}

// ListSponsoredBy lists the agents a person sponsors. Page.Truncated reports a list cut short by the
// server's cap.
//
// API: GET /v1/tenants/{tenantID}/agents/sponsored-by/{userID}
func (s *AgentsService) ListSponsoredBy(ctx context.Context, tenantID string, userID string, opts ...RequestOption) (*Page[Agent], error) {
	return call[Page[Agent]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/agents/sponsored-by/{userID}",
		pathArgs: []string{tenantID, userID},
	}, opts)
}

// TransferSponsorship gives every agent a person sponsors to another sponsor, as when the person
// leaves.
//
// API: POST /v1/tenants/{tenantID}/agents/sponsored-by/{userID}/transfer
func (s *AgentsService) TransferSponsorship(ctx context.Context, tenantID string, userID string, params TransferSponsorshipParams, opts ...RequestOption) (*SponsorshipTransfer, error) {
	return call[SponsorshipTransfer](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/agents/sponsored-by/{userID}/transfer",
		pathArgs: []string{tenantID, userID},
		body:     params,
	}, opts)
}

// Token mints a short-lived access token for an agent from its ID and secret. Agent tokens are minted
// again, never refreshed.
//
// The endpoint is public: the client sends no credentials with this request.
//
// API: POST /v1/agents/token
func (s *AgentsService) Token(ctx context.Context, params AgentTokenParams, opts ...RequestOption) (*AgentToken, error) {
	return call[AgentToken](ctx, s.c, &request{method: http.MethodPost, path: "/v1/agents/token", body: params, public: true}, opts)
}

// CreateAgentParams is the request body of AgentsService.Create.
type CreateAgentParams struct {
	Name            string   `json:"name,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
	TokenTTLSeconds *int     `json:"token_ttl_seconds,omitempty"`
	// The named human owner accountable for this agent. Must be a member of this tenant.
	SponsorUserID string `json:"sponsor_user_id,omitempty"`
}

// Agent is an AI agent identity. Its secret is returned only when it is created.
type Agent struct {
	CreatedAt       time.Time `json:"created_at,omitzero"`
	Disabled        bool      `json:"disabled"`
	ID              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Scopes          []string  `json:"scopes,omitempty"`
	Secret          string    `json:"secret,omitempty"`
	SponsorUserID   string    `json:"sponsor_user_id,omitempty"`
	Status          string    `json:"status,omitempty"`
	TokenTTLSeconds int       `json:"token_ttl_seconds"`
	// When the agent was last issued a token, to the minute.
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	// Present while the secret replaced at its last rotation still works, until then.
	PreviousSecretExpiresAt time.Time `json:"previous_secret_expires_at,omitzero"`
}

// UpdateAgentParams is the request body of AgentsService.Update.
type UpdateAgentParams struct {
	Disabled *bool `json:"disabled,omitempty"`
}

// AgentStatusResult is the result of AgentsService.Suspend.
type AgentStatusResult struct {
	Status AgentStatus `json:"status,omitempty"`
}

// AgentStatus is a string enumeration used by AgentStatusResult.Status.
type AgentStatus string

// Values of AgentStatus. The server may add values; handle an unknown one gracefully.
const (
	AgentStatusActive         AgentStatus = "active"
	AgentStatusSuspended      AgentStatus = "suspended"
	AgentStatusDecommissioned AgentStatus = "decommissioned"
)

// SuspendedAgents is the result of AgentsService.SuspendAll.
type SuspendedAgents struct {
	Suspended int `json:"suspended"`
}

// ResumedAgents is the result of AgentsService.ResumeAll.
type ResumedAgents struct {
	Resumed int `json:"resumed"`
}

// RotatedAgentSecret is the result of AgentsService.RotateSecret.
type RotatedAgentSecret struct {
	ID     string `json:"id,omitempty"`
	Secret string `json:"secret,omitempty"`
	// When the old secret stops working.
	PreviousSecretExpiresAt time.Time `json:"previous_secret_expires_at,omitzero"`
}

// TransferSponsorshipParams is the request body of AgentsService.TransferSponsorship.
type TransferSponsorshipParams struct {
	// Must be a member of this tenant.
	ToUserID string `json:"to_user_id,omitempty"`
}

// SponsorshipTransfer is the result of AgentsService.TransferSponsorship.
type SponsorshipTransfer struct {
	Transferred int `json:"transferred"`
}

// AgentTokenParams is the request body of AgentsService.Token.
type AgentTokenParams struct {
	AgentID string `json:"agent_id,omitempty"`
	Secret  string `json:"secret,omitempty"`
}

// AgentToken is a short-lived access token minted for an agent.
type AgentToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}
