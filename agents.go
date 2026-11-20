package qeetid

import (
	"context"
	"errors"
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
// API: POST /v1/organizations/{organizationID}/agents
func (s *AgentsService) Create(ctx context.Context, organizationID string, params CreateAgentParams, opts ...RequestOption) (*Agent, error) {
	return call[Agent](ctx, s.c, &request{
		method:     http.MethodPost,
		path:       "/v1/organizations/{organizationID}/agents",
		pathArgs:   []string{organizationID},
		body:       params,
		idempotent: true,
	}, opts)
}

// List lists an organization's agents. Page.Truncated reports a list cut short by the server's cap.
//
// API: GET /v1/organizations/{organizationID}/agents
func (s *AgentsService) List(ctx context.Context, organizationID string, opts ...RequestOption) (*Page[Agent], error) {
	return call[Page[Agent]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/agents",
		pathArgs: []string{organizationID},
	}, opts)
}

// Update enables or disables an agent.
//
// API: PATCH /v1/organizations/{organizationID}/agents/{id}
func (s *AgentsService) Update(ctx context.Context, organizationID, agentID string, params UpdateAgentParams, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPatch,
		path:     "/v1/organizations/{organizationID}/agents/{agentID}",
		pathArgs: []string{organizationID, agentID},
		body:     params,
	}, nil, opts)
}

// Delete deletes an agent.
//
// API: DELETE /v1/organizations/{organizationID}/agents/{id}
func (s *AgentsService) Delete(ctx context.Context, organizationID, agentID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}/agents/{agentID}",
		pathArgs: []string{organizationID, agentID},
	}, nil, opts)
}

// Suspend suspends an agent. Its tokens stop working until Resume.
//
// API: POST /v1/organizations/{organizationID}/agents/{id}/suspend
func (s *AgentsService) Suspend(ctx context.Context, organizationID, agentID string, opts ...RequestOption) (*AgentStatusResult, error) {
	return call[AgentStatusResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/{agentID}/suspend",
		pathArgs: []string{organizationID, agentID},
	}, opts)
}

// Resume makes a suspended agent active again.
//
// API: POST /v1/organizations/{organizationID}/agents/{id}/resume
func (s *AgentsService) Resume(ctx context.Context, organizationID, agentID string, opts ...RequestOption) (*AgentStatusResult, error) {
	return call[AgentStatusResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/{agentID}/resume",
		pathArgs: []string{organizationID, agentID},
	}, opts)
}

// Decommission decommissions an agent for good. It cannot be resumed.
//
// API: POST /v1/organizations/{organizationID}/agents/{id}/decommission
func (s *AgentsService) Decommission(ctx context.Context, organizationID, agentID string, opts ...RequestOption) (*AgentStatusResult, error) {
	return call[AgentStatusResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/{agentID}/decommission",
		pathArgs: []string{organizationID, agentID},
	}, opts)
}

// SuspendAll suspends every active agent of the organization at once (the kill switch).
//
// API: POST /v1/organizations/{organizationID}/agents/kill-all
func (s *AgentsService) SuspendAll(ctx context.Context, organizationID string, opts ...RequestOption) (*SuspendedAgents, error) {
	return call[SuspendedAgents](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/kill-all",
		pathArgs: []string{organizationID},
	}, opts)
}

// ResumeAll makes every suspended agent of the organization active again.
//
// API: POST /v1/organizations/{organizationID}/agents/resume-all
func (s *AgentsService) ResumeAll(ctx context.Context, organizationID string, opts ...RequestOption) (*ResumedAgents, error) {
	return call[ResumedAgents](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/resume-all",
		pathArgs: []string{organizationID},
	}, opts)
}

// RotateSecret replaces an agent's secret. The old secret keeps working for the overlap given.
//
// API: POST /v1/organizations/{organizationID}/agents/{id}/rotate-secret
func (s *AgentsService) RotateSecret(ctx context.Context, organizationID, agentID string, params RotateSecretParams, opts ...RequestOption) (*RotatedAgentSecret, error) {
	return call[RotatedAgentSecret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/{agentID}/rotate-secret",
		pathArgs: []string{organizationID, agentID},
		body:     params,
	}, opts)
}

// ListSponsoredBy lists the agents a person sponsors. Page.Truncated reports a list cut short by the
// server's cap.
//
// API: GET /v1/organizations/{organizationID}/agents/sponsored-by/{userID}
func (s *AgentsService) ListSponsoredBy(ctx context.Context, organizationID, userID string, opts ...RequestOption) (*Page[Agent], error) {
	return call[Page[Agent]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/agents/sponsored-by/{userID}",
		pathArgs: []string{organizationID, userID},
	}, opts)
}

// TransferSponsorship gives every agent a person sponsors to another sponsor, as when the person
// leaves.
//
// API: POST /v1/organizations/{organizationID}/agents/sponsored-by/{userID}/transfer
func (s *AgentsService) TransferSponsorship(ctx context.Context, organizationID, userID string, params TransferSponsorshipParams, opts ...RequestOption) (*SponsorshipTransfer, error) {
	return call[SponsorshipTransfer](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/agents/sponsored-by/{userID}/transfer",
		pathArgs: []string{organizationID, userID},
		body:     params,
	}, opts)
}

// Token obtains a short-lived access token for an agent from its ID and secret, with the client
// credentials grant of the OAuth token endpoint. It is OAuthService.ClientCredentials with the agent's
// ID as the client ID. Agent tokens are obtained again, never refreshed; to have the client do that for
// its own requests, use WithClientCredentials with the agent's ID and secret.
//
// The agent authenticates with HTTP Basic; the client's own credential is not sent. A wrong ID or
// secret fails with an *OAuthError whose Code is OAuthInvalidClient, and an agent that isn't active
// with OAuthUnauthorizedClient.
//
// API: POST /v1/oauth/token
func (s *AgentsService) Token(ctx context.Context, params AgentTokenParams, opts ...RequestOption) (*Token, error) {
	if params.AgentID == "" || params.Secret == "" {
		return nil, errors.New("qeetid: Agents.Token: AgentID and Secret are required")
	}
	return s.c.clientCredentialsToken(ctx, params.AgentID, params.Secret, "", opts...)
}

// CreateAgentParams is the request body of AgentsService.Create.
type CreateAgentParams struct {
	Name            string   `json:"name,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
	TokenTTLSeconds *int     `json:"token_ttl_seconds,omitempty"`
	// The named human owner accountable for this agent. Must be a member of this organization.
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
	// Must be a member of this organization.
	ToUserID string `json:"to_user_id,omitempty"`
}

// SponsorshipTransfer is the result of AgentsService.TransferSponsorship.
type SponsorshipTransfer struct {
	Transferred int `json:"transferred"`
}

// AgentTokenParams are the parameters of AgentsService.Token.
type AgentTokenParams struct {
	// AgentID is the agent's ID, sent as the client ID.
	AgentID string
	// Secret is the agent's secret, sent as the client secret.
	Secret string
}
