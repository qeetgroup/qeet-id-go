package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// PurgeRequest is a GDPR erasure request against a user, with a grace
// window before it actually runs.
type PurgeRequest struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	UserID      string `json:"user_id"`
	RequestedBy string `json:"requested_by,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Status      string `json:"status"`
	GraceUntil  string `json:"grace_until"`
	CompletedAt string `json:"completed_at,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// CreatePurgeInput — TenantID/UserID identify who to erase; the requester
// is inferred server-side from the caller's principal, not sent in the body.
type CreatePurgeInput struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	Reason   string `json:"reason,omitempty"`
}

// ExportRequest is a GDPR data-export request. Payload is populated inline
// once Status is "ready" — there is no separate download-URL field; fetch
// GetExport again once ready and read Payload directly.
type ExportRequest struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenant_id"`
	UserID      string         `json:"user_id"`
	RequestedBy string         `json:"requested_by,omitempty"`
	Status      string         `json:"status"` // pending | ready | failed
	Payload     map[string]any `json:"payload,omitempty"`
	Error       string         `json:"error,omitempty"`
	CompletedAt string         `json:"completed_at,omitempty"`
	CreatedAt   string         `json:"created_at"`
}

type CreateExportInput struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

// GDPRService manages erasure (purge) and data-export requests.
type GDPRService struct{ t *transport.Transport }

// CreatePurge is not tenant-path-scoped — tenant/user are identified in the body.
func (s *GDPRService) CreatePurge(ctx context.Context, in CreatePurgeInput) (*PurgeRequest, error) {
	var out PurgeRequest
	err := s.t.Post(ctx, "/v1/gdpr/purge", in, &out)
	return &out, err
}

func (s *GDPRService) CancelPurge(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/gdpr/purge/"+url.PathEscape(id), nil)
}

func (s *GDPRService) ListPurge(ctx context.Context, tenantID string) ([]PurgeRequest, error) {
	var env marshal.Envelope[PurgeRequest]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/gdpr/purge", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// CreateExport is async (202 Accepted) — poll GetExport until Status is
// "ready" or "failed".
func (s *GDPRService) CreateExport(ctx context.Context, in CreateExportInput) (*ExportRequest, error) {
	var out ExportRequest
	err := s.t.Post(ctx, "/v1/gdpr/export", in, &out)
	return &out, err
}

func (s *GDPRService) ListExports(ctx context.Context, tenantID string) ([]ExportRequest, error) {
	var env marshal.Envelope[ExportRequest]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/gdpr/export", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *GDPRService) GetExport(ctx context.Context, tenantID, id string) (*ExportRequest, error) {
	var out ExportRequest
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/gdpr/export/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}
