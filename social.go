package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// SocialProvider is a configured social-login provider (Google, GitHub,
// etc.). ClientSecret is write-only — never returned.
type SocialProvider struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	Provider     string `json:"provider"`
	ClientID     string `json:"client_id"`
	DiscoveryURL string `json:"discovery_url,omitempty"`
	Enabled      bool   `json:"enabled"`
	CreatedAt    string `json:"created_at"`
}

type UpsertSocialProviderInput struct {
	TenantID     string `json:"tenant_id"`
	Provider     string `json:"provider"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	DiscoveryURL string `json:"discovery_url"`
}

// ExternalIdentity links a user to an upstream social/OIDC provider account.
type ExternalIdentity struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Provider string `json:"provider"`
	Subject  string `json:"subject"`
	Email    string `json:"email,omitempty"`
	LinkedAt string `json:"linked_at"`
}

// SocialService manages social-login provider configuration and the
// identities users have linked. The browser-redirect endpoints
// (start/callback/exchange) belong to the client-side SDK, not here.
type SocialService struct{ t *transport.Transport }

func (s *SocialService) ListProviders(ctx context.Context, tenantID string) ([]SocialProvider, error) {
	var env marshal.Envelope[SocialProvider]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/social/providers", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// UpsertProvider creates or updates a tenant's provider config.
func (s *SocialService) UpsertProvider(ctx context.Context, in UpsertSocialProviderInput) (*SocialProvider, error) {
	var out SocialProvider
	err := s.t.Post(ctx, "/v1/social/providers", in, &out)
	return &out, err
}

func (s *SocialService) ListUserIdentities(ctx context.Context, userID string) ([]ExternalIdentity, error) {
	var env marshal.Envelope[ExternalIdentity]
	err := s.t.Get(ctx, "/v1/users/"+url.PathEscape(userID)+"/social/identities", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *SocialService) UnlinkIdentity(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/social/identities/"+url.PathEscape(id), nil)
}
