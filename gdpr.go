package qeetid

import (
	"context"
	"net/http"
	"time"
)

// GDPRService handles erasure and data-export requests.
type GDPRService struct{ c *Client }

// CreatePurge requests the erasure of a person's data. It runs when the grace period ends unless
// cancelled.
//
// API: POST /v1/gdpr/purge
func (s *GDPRService) CreatePurge(ctx context.Context, params CreatePurgeRequestParams, opts ...RequestOption) (*PurgeRequest, error) {
	return call[PurgeRequest](ctx, s.c, &request{method: http.MethodPost, path: "/v1/gdpr/purge", body: params}, opts)
}

// CancelPurge cancels a pending erasure request.
//
// API: DELETE /v1/gdpr/purge/{id}
func (s *GDPRService) CancelPurge(ctx context.Context, requestID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/gdpr/purge/{requestID}",
		pathArgs: []string{requestID},
	}, nil, opts)
}

// ListPurges lists a tenant's erasure requests.
//
// API: GET /v1/tenants/{tenantID}/gdpr/purge
func (s *GDPRService) ListPurges(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[PurgeRequest], error) {
	return call[Page[PurgeRequest]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/gdpr/purge",
		pathArgs: []string{tenantID},
	}, opts)
}

// CreateExport requests an export of a person's data. Poll GetExport until it is ready.
//
// API: POST /v1/gdpr/export
func (s *GDPRService) CreateExport(ctx context.Context, params CreateExportRequestParams, opts ...RequestOption) (*ExportRequest, error) {
	return call[ExportRequest](ctx, s.c, &request{method: http.MethodPost, path: "/v1/gdpr/export", body: params}, opts)
}

// ListExports lists a tenant's export requests.
//
// API: GET /v1/tenants/{tenantID}/gdpr/export
func (s *GDPRService) ListExports(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[ExportRequest], error) {
	return call[Page[ExportRequest]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/gdpr/export",
		pathArgs: []string{tenantID},
	}, opts)
}

// GetExport returns an export request, with its payload once ready.
//
// API: GET /v1/tenants/{tenantID}/gdpr/export/{id}
func (s *GDPRService) GetExport(ctx context.Context, tenantID string, requestID string, opts ...RequestOption) (*ExportRequest, error) {
	return call[ExportRequest](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/gdpr/export/{requestID}",
		pathArgs: []string{tenantID, requestID},
	}, opts)
}

// CreatePurgeRequestParams is the request body of GDPRService.CreatePurge.
type CreatePurgeRequestParams struct {
	Reason   string `json:"reason,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// PurgeRequest is a request to erase a person's data. It runs when its grace period ends unless
// cancelled.
type PurgeRequest struct {
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitzero"`
	GraceUntil  time.Time  `json:"grace_until,omitzero"`
	ID          string     `json:"id,omitempty"`
	Reason      string     `json:"reason,omitempty"`
	RequestedBy string     `json:"requested_by,omitempty"`
	Status      string     `json:"status,omitempty"`
	TenantID    string     `json:"tenant_id,omitempty"`
	// The account to erase; kept after the account itself is deleted.
	UserID string `json:"user_id,omitempty"`
	// The person's email as the request was filed; null once they are erased.
	SubjectEmail string `json:"subject_email,omitempty"`
	// Their display name as the request was filed; null once they are erased.
	SubjectName string `json:"subject_name,omitempty"`
	// The erasure ran, so who it was filed for is no longer kept.
	SubjectRedacted bool `json:"subject_redacted"`
	// How many days a request waits before the erasure runs.
	GracePeriodDays int `json:"grace_period_days"`
}

// CreateExportRequestParams is the request body of GDPRService.CreateExport.
type CreateExportRequestParams struct {
	TenantID string `json:"tenant_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// ExportRequest is a request to export a person's data. Payload is set once Status is ready.
type ExportRequest struct {
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitzero"`
	Error       string     `json:"error,omitempty"`
	ID          string     `json:"id,omitempty"`
	// The person's data in the organisation, while status is ready: profile (with metadata), sessions,
	// passkeys, roles, mfa_totp_enabled, identities (linked social and single sign-on accounts),
	// otp_factors (destinations masked), app_consents, app_grants, groups, sign_ins (the latest 1000) and
	// generated_at.
	Payload     map[string]any `json:"payload,omitempty"`
	RequestedBy string         `json:"requested_by,omitempty"`
	// pending, ready, failed, or expired once the payload has been deleted (seven days after it was
	// built, or when the person was erased or their account deleted).
	Status   string `json:"status,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	// The account exported; kept after the account itself is deleted.
	UserID string `json:"user_id,omitempty"`
}
