package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type Agent struct {
	ID              string   `json:"id"`
	TenantID        string   `json:"tenant_id"`
	Name            string   `json:"name"`
	Scopes          []string `json:"scopes"`
	TokenTTLSeconds int      `json:"token_ttl_seconds"`
	Disabled        bool     `json:"disabled"`
	CreatedAt       string   `json:"created_at"`
	Secret          string   `json:"secret,omitempty"` // only on create
}

type CreateAgentInput struct {
	Name            string   `json:"name"`
	Scopes          []string `json:"scopes,omitempty"`
	TokenTTLSeconds int      `json:"token_ttl_seconds,omitempty"`
}

// UpdateAgentInput's only mutable field, per the backend, is Disabled.
type UpdateAgentInput struct {
	Disabled bool `json:"disabled"`
}

type AgentTokenResult struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// AgentStatus is the lifecycle state returned by Suspend/Resume/Decommission.
type AgentStatus struct {
	Status string `json:"status"` // active | suspended | decommissioned
}

// KillAllResult reports how many agents an incident-response kill-switch
// suspended.
type KillAllResult struct {
	Suspended int `json:"suspended"`
}

// TransferSponsorshipInput names the new sponsor for an offboarded user's
// agents.
type TransferSponsorshipInput struct {
	ToUserID string `json:"to_user_id"`
}

// TransferSponsorshipResult reports how many agents were transferred.
type TransferSponsorshipResult struct {
	Transferred int `json:"transferred"`
}

// AgentsService manages AI-agent identities: ephemeral tokens and the
// suspend/resume/decommission lifecycle.
type AgentsService struct{ t *transport.Transport }

func (s *AgentsService) Create(ctx context.Context, tenantID string, in CreateAgentInput) (*Agent, error) {
	var out Agent
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents", in, &out)
	return &out, err
}

func (s *AgentsService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents/"+url.PathEscape(id), nil)
}

func (s *AgentsService) List(ctx context.Context, tenantID string) ([]Agent, error) {
	var env marshal.Envelope[Agent]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// Update toggles Disabled — the only field the backend allows mutating here.
func (s *AgentsService) Update(ctx context.Context, tenantID, id string, in UpdateAgentInput) (*Agent, error) {
	var out Agent
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents/"+url.PathEscape(id), in, &out)
	return &out, err
}

// Token mints a short-lived access token for an AI agent.
func (s *AgentsService) Token(ctx context.Context, tenantID, agentID, secret, scope string) (*AgentTokenResult, error) {
	body := map[string]string{
		"tenant_id": tenantID,
		"agent_id":  agentID,
		"secret":    secret,
	}
	if scope != "" {
		body["scope"] = scope
	}
	var out AgentTokenResult
	err := s.t.Post(ctx, "/v1/agents/token", body, &out)
	return &out, err
}

// Suspend reversibly disables an agent — it can be Resumed.
func (s *AgentsService) Suspend(ctx context.Context, tenantID, id string) (*AgentStatus, error) {
	var out AgentStatus
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents/"+url.PathEscape(id)+"/suspend", nil, &out)
	return &out, err
}

// Resume reverses a Suspend.
func (s *AgentsService) Resume(ctx context.Context, tenantID, id string) (*AgentStatus, error) {
	var out AgentStatus
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents/"+url.PathEscape(id)+"/resume", nil, &out)
	return &out, err
}

// Decommission is terminal — a decommissioned agent cannot be Resumed.
func (s *AgentsService) Decommission(ctx context.Context, tenantID, id string) (*AgentStatus, error) {
	var out AgentStatus
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents/"+url.PathEscape(id)+"/decommission", nil, &out)
	return &out, err
}

// KillAll is an incident-response switch: suspends every agent in the tenant.
func (s *AgentsService) KillAll(ctx context.Context, tenantID string) (*KillAllResult, error) {
	var out KillAllResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/agents/kill-all", nil, &out)
	return &out, err
}

// ListSponsoredBy lists every agent a given human user sponsors.
func (s *AgentsService) ListSponsoredBy(ctx context.Context, tenantID, userID string) ([]Agent, error) {
	path := "/v1/tenants/" + url.PathEscape(tenantID) + "/agents/sponsored-by/" + url.PathEscape(userID)
	var env marshal.Envelope[Agent]
	err := s.t.Get(ctx, path, transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// TransferSponsorship reassigns every agent sponsored by userID to a new
// sponsor — the standard step when offboarding a human sponsor.
func (s *AgentsService) TransferSponsorship(ctx context.Context, tenantID, userID string, in TransferSponsorshipInput) (*TransferSponsorshipResult, error) {
	path := "/v1/tenants/" + url.PathEscape(tenantID) + "/agents/sponsored-by/" + url.PathEscape(userID) + "/transfer"
	var out TransferSponsorshipResult
	err := s.t.Post(ctx, path, in, &out)
	return &out, err
}
