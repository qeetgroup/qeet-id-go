package qeetid

import (
	"context"
	"iter"
	"net/url"
	"strconv"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/pagination"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// AuditLog is one row from the hash-chained audit log. Metadata and the
// hash-chain fields (prev_hash/row_hash) aren't part of the list response —
// only Verify walks the chain.
type AuditLog struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	ActorUserID  string `json:"actor_user_id,omitempty"`
	ActorType    string `json:"actor_type,omitempty"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
	IP           string `json:"ip,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// AuditLogListParams narrows a List call. Search applies a free-text
// websearch_to_tsquery filter (?q=) over action/resource_type/actor_type/
// user_agent/metadata — supports "quoted phrases", -exclusions, and OR.
type AuditLogListParams struct {
	Action       string
	ResourceType string
	ActorUserID  string
	Search       string
	Limit        int
	Cursor       string
}

type AuditLogPage struct {
	Data       []AuditLog
	NextCursor string
}

// AuditChainVerification is the result of walking a tenant's hash chain from
// the seed (GET /audit/verify) — a broken chain names the first bad row.
type AuditChainVerification struct {
	OK             bool   `json:"ok"`
	RowsChecked    int    `json:"rows_checked"`
	LastVerifiedID string `json:"last_verified_id,omitempty"`
	BrokenAtID     string `json:"broken_at_id,omitempty"`
	BrokenReason   string `json:"broken_reason,omitempty"`
}

// AuditAnomaly is a flagged deviation from an actor's behavioral baseline.
type AuditAnomaly struct {
	ID           string   `json:"id"`
	TenantID     string   `json:"tenant_id"`
	EventID      string   `json:"event_id"`
	ActorUserID  string   `json:"actor_user_id,omitempty"`
	ActorEmail   string   `json:"actor_email,omitempty"`
	Score        float64  `json:"score"`
	Reasons      []string `json:"reasons"` // new_action_type | unusual_hour | new_ip
	Status       string   `json:"status"`  // open | resolved
	ResolvedAt   string   `json:"resolved_at,omitempty"`
	ResolvedBy   string   `json:"resolved_by,omitempty"`
	Action       string   `json:"action"`
	ResourceType string   `json:"resource_type"`
	IP           string   `json:"ip,omitempty"`
	EventAt      string   `json:"event_at"`
	CreatedAt    string   `json:"created_at"`
}

// AuditAnomalySummary is the counts view for a tenant's anomaly queue.
type AuditAnomalySummary struct {
	Open          int `json:"open"`
	Resolved7d    int `json:"resolved_7d"`
	HighScoreOpen int `json:"high_score_open"`
}

// AuditAnomalySettings tunes the behavioral-baseline scorer.
type AuditAnomalySettings struct {
	TenantID          string  `json:"tenant_id"`
	Enabled           bool    `json:"enabled"`
	ScoreThreshold    float64 `json:"score_threshold"`
	MinBaselineEvents int     `json:"min_baseline_events"`
}

// UpdateAuditAnomalySettingsInput is the PUT body for AuditAnomalies.PutSettings.
type UpdateAuditAnomalySettingsInput struct {
	Enabled           *bool    `json:"enabled,omitempty"`
	ScoreThreshold    *float64 `json:"score_threshold,omitempty"`
	MinBaselineEvents *int     `json:"min_baseline_events,omitempty"`
}

// AuditLogsService reads the hash-chained audit log and its anomaly-detection
// sub-resource.
type AuditLogsService struct{ t *transport.Transport }

func (s *AuditLogsService) List(ctx context.Context, tenantID string, params AuditLogListParams) (*AuditLogPage, error) {
	q := url.Values{}
	if params.Action != "" {
		q.Set("action", params.Action)
	}
	if params.ResourceType != "" {
		q.Set("resource_type", params.ResourceType)
	}
	if params.ActorUserID != "" {
		q.Set("actor_user_id", params.ActorUserID)
	}
	if params.Search != "" {
		q.Set("q", params.Search)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	var env struct {
		marshal.Envelope[AuditLog]
		NextCursor string `json:"next_cursor"`
	}
	if err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit", transport.RequestOptions{Query: q}, &env); err != nil {
		return nil, err
	}
	return &AuditLogPage{Data: env.Resolve(), NextCursor: env.NextCursor}, nil
}

// All iterates every audit-log entry across pages.
func (s *AuditLogsService) All(ctx context.Context, tenantID string, params AuditLogListParams) iter.Seq2[AuditLog, error] {
	return pagination.Paginate(ctx, params.Cursor, func(ctx context.Context, cursor string) ([]AuditLog, string, error) {
		p := params
		p.Cursor = cursor
		page, err := s.List(ctx, tenantID, p)
		if err != nil {
			return nil, "", err
		}
		return page.Data, page.NextCursor, nil
	})
}

// Verify walks the tenant's hash chain from the seed and recomputes each
// row's hash, returning the first broken link (if any). No side effects.
// (Not per-entry — the backend has no such endpoint; it verifies the whole
// chain in one pass.)
func (s *AuditLogsService) Verify(ctx context.Context, tenantID string) (*AuditChainVerification, error) {
	var out AuditChainVerification
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit/verify", transport.RequestOptions{}, &out)
	return &out, err
}

// Anomalies lists flagged behavioral-baseline deviations. status filters to
// "open" or "resolved"; empty returns both.
func (s *AuditLogsService) Anomalies(ctx context.Context, tenantID, status string, limit int) ([]AuditAnomaly, error) {
	q := url.Values{}
	if status != "" {
		q.Set("status", status)
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var env marshal.Envelope[AuditAnomaly]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit/anomalies", transport.RequestOptions{Query: q}, &env)
	return env.Resolve(), err
}

func (s *AuditLogsService) AnomalySummary(ctx context.Context, tenantID string) (*AuditAnomalySummary, error) {
	var out AuditAnomalySummary
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit/anomalies/summary", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *AuditLogsService) ResolveAnomaly(ctx context.Context, tenantID, id string) error {
	return s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit/anomalies/"+url.PathEscape(id)+"/resolve", struct{}{}, nil)
}

func (s *AuditLogsService) GetAnomalySettings(ctx context.Context, tenantID string) (*AuditAnomalySettings, error) {
	var out AuditAnomalySettings
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit/anomaly-settings", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *AuditLogsService) PutAnomalySettings(ctx context.Context, tenantID string, in UpdateAuditAnomalySettingsInput) (*AuditAnomalySettings, error) {
	var out AuditAnomalySettings
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/audit/anomaly-settings", in, &out)
	return &out, err
}
