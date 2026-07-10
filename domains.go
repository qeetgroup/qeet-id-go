package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Domain is a custom domain pending or completed DNS verification.
// VerifiedAt is nil until the DNS record is confirmed — there is no separate
// boolean flag. DNSRecordName/Type/Value are what the caller needs to
// actually create in their DNS provider. There is no per-domain Get in the
// backend — only Create, Delete, Verify, and List.
type Domain struct {
	ID                string `json:"id"`
	Domain            string `json:"domain"`
	VerificationToken string `json:"verification_token,omitempty"`
	DNSRecordName     string `json:"dns_record_name,omitempty"`
	DNSRecordType     string `json:"dns_record_type,omitempty"`
	DNSRecordValue    string `json:"dns_record_value,omitempty"`
	VerifiedAt        string `json:"verified_at,omitempty"`
	CreatedAt         string `json:"created_at"`
}

// Verified reports whether DNS verification has completed.
func (d Domain) Verified() bool { return d.VerifiedAt != "" }

type CreateDomainInput struct {
	Domain string `json:"domain"`
}

// DomainsService manages custom domains.
type DomainsService struct{ t *transport.Transport }

func (s *DomainsService) Create(ctx context.Context, tenantID string, in CreateDomainInput) (*Domain, error) {
	var out Domain
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/domains", in, &out)
	return &out, err
}

func (s *DomainsService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/domains/"+url.PathEscape(id), nil)
}

func (s *DomainsService) Verify(ctx context.Context, tenantID, id string) (*Domain, error) {
	var out Domain
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/domains/"+url.PathEscape(id)+"/verify", struct{}{}, &out)
	return &out, err
}

func (s *DomainsService) List(ctx context.Context, tenantID string) ([]Domain, error) {
	var env marshal.Envelope[Domain]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/domains", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
