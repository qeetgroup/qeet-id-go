package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type Secret struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Scope     string `json:"scope"`
	Last4     string `json:"last4"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateSecretInput struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
	Value string `json:"value"`
}

type UpdateSecretInput struct {
	Scope *string `json:"scope,omitempty"`
	Value *string `json:"value,omitempty"`
}

type VaultGetResult struct {
	Value string `json:"value"`
}

// VaultService is the encrypted secrets store used by agent/developer
// credentials (3rd-party OAuth tokens, API keys agents hold on a user's
// behalf).
type VaultService struct{ t *transport.Transport }

// Get fetches the value of a vault secret by name (agent-scoped endpoint).
func (s *VaultService) Get(ctx context.Context, name string) (*VaultGetResult, error) {
	var out VaultGetResult
	err := s.t.Get(ctx, "/v1/vault/"+url.PathEscape(name), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *VaultService) ListSecrets(ctx context.Context, tenantID string) ([]Secret, error) {
	var env marshal.Envelope[Secret]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/secrets", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *VaultService) CreateSecret(ctx context.Context, tenantID string, in CreateSecretInput) (*Secret, error) {
	var out Secret
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/secrets", in, &out)
	return &out, err
}

func (s *VaultService) UpdateSecret(ctx context.Context, tenantID, id string, in UpdateSecretInput) (*Secret, error) {
	var out Secret
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/secrets/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *VaultService) RevealSecret(ctx context.Context, tenantID, id string) (*VaultGetResult, error) {
	var out VaultGetResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/secrets/"+url.PathEscape(id)+"/reveal", struct{}{}, &out)
	return &out, err
}

func (s *VaultService) DeleteSecret(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/secrets/"+url.PathEscape(id), nil)
}
