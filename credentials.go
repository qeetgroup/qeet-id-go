package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type Credential struct {
	ID        string `json:"id"`
	Subject   string `json:"subject"`
	Type      string `json:"type"`
	IssuedAt  string `json:"issued_at"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Revoked   bool   `json:"revoked"`
}

type IssueCredentialInput struct {
	Subject    string         `json:"subject"`
	Type       string         `json:"type"`
	Claims     map[string]any `json:"claims,omitempty"`
	TTLSeconds int            `json:"ttl_seconds,omitempty"`
}

type IssueCredentialResult struct {
	CredentialID string `json:"credential_id"`
	JWT          string `json:"jwt"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}

type VerifyCredentialResult struct {
	Valid   bool           `json:"valid"`
	Reason  string         `json:"reason,omitempty"`
	Subject string         `json:"subject,omitempty"`
	Issuer  string         `json:"issuer,omitempty"`
	VC      map[string]any `json:"vc,omitempty"`
}

// CredentialsService issues, lists, revokes, and verifies W3C Verifiable
// Credentials.
type CredentialsService struct{ t *transport.Transport }

func (s *CredentialsService) Issue(ctx context.Context, tenantID string, in IssueCredentialInput) (*IssueCredentialResult, error) {
	var out IssueCredentialResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/credentials", in, &out)
	return &out, err
}

func (s *CredentialsService) List(ctx context.Context, tenantID string) ([]Credential, error) {
	var env marshal.Envelope[Credential]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/credentials", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *CredentialsService) Revoke(ctx context.Context, tenantID, id string) error {
	return s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/credentials/"+url.PathEscape(id)+"/revoke", struct{}{}, nil)
}

// Verify is the public endpoint — no API key required. Relying parties call
// this to confirm a presented JWT-VC is authentic and not revoked.
func (s *CredentialsService) Verify(ctx context.Context, jwt string) (*VerifyCredentialResult, error) {
	var out VerifyCredentialResult
	body := map[string]string{"credential": jwt}
	err := s.t.Post(ctx, "/v1/credentials/verify", body, &out)
	return &out, err
}
