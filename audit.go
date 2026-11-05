package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AuditService reads a tenant's tamper-evident audit log and manages its anomaly detection.
type AuditService struct{ c *Client }

// List lists audit events, newest first, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/audit
func (s *AuditService) List(ctx context.Context, tenantID string, params *ListAuditEventsParams, opts ...RequestOption) (*Page[AuditEvent], error) {
	return call[Page[AuditEvent]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/audit",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// All returns an iterator over every item List returns, across all pages. It fetches a page only when
// the loop needs it, and breaking out of the loop stops paging. A failed request ends the iteration
// with a non-nil error.
func (s *AuditService) All(ctx context.Context, tenantID string, params *ListAuditEventsParams, opts ...RequestOption) iter.Seq2[*AuditEvent, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[AuditEvent], error) {
		return s.List(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// Verify walks the tenant's hash chain and reports the first broken link, if any.
//
// API: GET /v1/tenants/{tenantID}/audit/verify
func (s *AuditService) Verify(ctx context.Context, tenantID string, opts ...RequestOption) (*AuditChainVerification, error) {
	return call[AuditChainVerification](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/audit/verify",
		pathArgs: []string{tenantID},
	}, opts)
}

// LastVerification returns the outcome of the most recent chain verification.
//
// API: GET /v1/tenants/{tenantID}/audit/verification
func (s *AuditService) LastVerification(ctx context.Context, tenantID string, opts ...RequestOption) (*AuditVerificationStatus, error) {
	return call[AuditVerificationStatus](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/audit/verification",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListAnomalies lists behavioral anomalies found in the audit log.
//
// API: GET /v1/tenants/{tenantID}/audit/anomalies
func (s *AuditService) ListAnomalies(ctx context.Context, tenantID string, params *ListAuditAnomaliesParams, opts ...RequestOption) (*Page[AuditAnomaly], error) {
	return call[Page[AuditAnomaly]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/audit/anomalies",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AnomalySummary returns counts of open and resolved anomalies.
//
// API: GET /v1/tenants/{tenantID}/audit/anomalies/summary
func (s *AuditService) AnomalySummary(ctx context.Context, tenantID string, opts ...RequestOption) (*AuditAnomalySummary, error) {
	return call[AuditAnomalySummary](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/audit/anomalies/summary",
		pathArgs: []string{tenantID},
	}, opts)
}

// ResolveAnomaly marks an anomaly resolved.
//
// API: POST /v1/tenants/{tenantID}/audit/anomalies/{id}/resolve
func (s *AuditService) ResolveAnomaly(ctx context.Context, tenantID string, anomalyID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/audit/anomalies/{anomalyID}/resolve",
		pathArgs: []string{tenantID, anomalyID},
	}, nil, opts)
}

// AnomalySettings returns the anomaly-detection settings.
//
// API: GET /v1/tenants/{tenantID}/audit/anomaly-settings
func (s *AuditService) AnomalySettings(ctx context.Context, tenantID string, opts ...RequestOption) (*AuditAnomalySettings, error) {
	return call[AuditAnomalySettings](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/audit/anomaly-settings",
		pathArgs: []string{tenantID},
	}, opts)
}

// UpdateAnomalySettings replaces the anomaly-detection settings.
//
// API: PUT /v1/tenants/{tenantID}/audit/anomaly-settings
func (s *AuditService) UpdateAnomalySettings(ctx context.Context, tenantID string, params AuditAnomalySettings, opts ...RequestOption) (*AuditAnomalySettings, error) {
	return call[AuditAnomalySettings](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/audit/anomaly-settings",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// ListAuditEventsParams holds the optional query parameters of AuditService.List. A nil
// *ListAuditEventsParams sends none.
type ListAuditEventsParams struct {
	// Finds words of the action, resource type, actor type, user agent, address and metadata, each split
	// into words (auth.login_succeeded is auth, login and succeeded), every plain word as a prefix; or
	// the events of an actor whose name contains it or whose email begins with it. "Quoted phrases",
	// -exclusions and OR are read as websearch_to_tsquery reads them.
	Q string
	// Comma-separated actions; events of any of them.
	Action string
	// Comma-separated resource types; events about any of them.
	ResourceType string
	ResourceID   string
	// Events from this address.
	IP string
	// Inclusive lower bound (RFC3339).
	From time.Time
	// Inclusive upper bound (RFC3339).
	To          time.Time
	ActorUserID string
	// Page size. Defaults to 50; the server clamps values above 200.
	Limit int
	// The previous page's next_cursor (the id of its last event), for keyset pagination.
	Cursor string
}

func (p *ListAuditEventsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Q != "" {
		q.Set("q", p.Q)
	}
	if p.Action != "" {
		q.Set("action", p.Action)
	}
	if p.ResourceType != "" {
		q.Set("resource_type", p.ResourceType)
	}
	if p.ResourceID != "" {
		q.Set("resource_id", p.ResourceID)
	}
	if p.IP != "" {
		q.Set("ip", p.IP)
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.UTC().Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.UTC().Format(time.RFC3339))
	}
	if p.ActorUserID != "" {
		q.Set("actor_user_id", p.ActorUserID)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ListAuditEventsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListAuditEventsParams) withCursor(cursor string) *ListAuditEventsParams {
	var c ListAuditEventsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// AuditEvent is one entry of a tenant's tamper-evident audit log.
type AuditEvent struct {
	Action       string    `json:"action,omitempty"`
	ActorType    string    `json:"actor_type,omitempty"`
	ActorUserID  string    `json:"actor_user_id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitzero"`
	ID           string    `json:"id,omitempty"`
	IP           string    `json:"ip,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	ResourceID   string    `json:"resource_id,omitempty"`
	ResourceType string    `json:"resource_type,omitempty"`
	TenantID     string    `json:"tenant_id,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	// What the event recorded beside its fields; geo is where the request came from.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// AuditChainVerification is the outcome of walking a tenant's audit hash chain.
type AuditChainVerification struct {
	BrokenAtID string `json:"broken_at_id,omitempty"`
	// Why the chain is broken at broken_at_id, in English. Show reason_code's meaning instead.
	BrokenReason   string `json:"broken_reason,omitempty"`
	LastVerifiedID string `json:"last_verified_id,omitempty"`
	Ok             bool   `json:"ok"`
	// Set when the chain is broken: the row's link to the one before it is wrong, or its contents no
	// longer match its hash.
	ReasonCode  AuditChainVerificationReasonCode `json:"reason_code,omitempty"`
	RowsChecked int                              `json:"rows_checked"`
}

// AuditChainVerificationReasonCode is a string enumeration used by AuditChainVerification.ReasonCode.
// Set when the chain is broken: the row's link to the one before it is wrong, or its contents no
// longer match its hash.
type AuditChainVerificationReasonCode string

// Values of AuditChainVerificationReasonCode. The server may add values; handle an unknown one gracefully.
const (
	AuditChainVerificationReasonCodeAuditChainPrevHashMismatch AuditChainVerificationReasonCode = "audit.chain_prev_hash_mismatch"
	AuditChainVerificationReasonCodeAuditChainRowHashMismatch  AuditChainVerificationReasonCode = "audit.chain_row_hash_mismatch"
)

// AuditVerificationStatus is the result of AuditService.LastVerification.
type AuditVerificationStatus struct {
	// Whether the scheduled verification has reached the chain.
	Verified bool `json:"verified"`
	// No break found; false also while not yet verified.
	Intact         bool       `json:"intact"`
	LastVerifiedAt *time.Time `json:"last_verified_at,omitempty"`
	// When a break was found; it stays until resolved.
	BrokenAt            *time.Time `json:"broken_at,omitempty"`
	LastVerifiedEventID string     `json:"last_verified_event_id,omitempty"`
	// When the latest scheduled run finished.
	LastRunAt *time.Time `json:"last_run_at,omitempty"`
	// all (a full walk) or incremental; empty before any run.
	LastRunScope string `json:"last_run_scope,omitempty"`
}

// AuditAnomalyStatus is a string enumeration used by ListAuditAnomaliesParams.Status.
type AuditAnomalyStatus string

// Values of AuditAnomalyStatus. The server may add values; handle an unknown one gracefully.
const (
	AuditAnomalyStatusOpen     AuditAnomalyStatus = "open"
	AuditAnomalyStatusResolved AuditAnomalyStatus = "resolved"
)

// ListAuditAnomaliesParams holds the optional query parameters of AuditService.ListAnomalies. A nil
// *ListAuditAnomaliesParams sends none.
type ListAuditAnomaliesParams struct {
	Status AuditAnomalyStatus
	Limit  int
}

func (p *ListAuditAnomaliesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Status != "" {
		q.Set("status", string(p.Status))
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	return q
}

// AuditAnomaly is an audit event that departs from its actor's usual behavior.
type AuditAnomaly struct {
	ID          string `json:"id,omitempty"`
	TenantID    string `json:"tenant_id,omitempty"`
	EventID     string `json:"event_id,omitempty"`
	ActorUserID string `json:"actor_user_id,omitempty"`
	ActorEmail  string `json:"actor_email,omitempty"`
	// 0..1 weighted-novelty anomaly score.
	Score        float64              `json:"score"`
	Reasons      []AuditAnomalyReason `json:"reasons,omitempty"`
	Status       AuditAnomalyStatus   `json:"status,omitempty"`
	ResolvedAt   *time.Time           `json:"resolved_at,omitempty"`
	ResolvedBy   string               `json:"resolved_by,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitzero"`
	Action       string               `json:"action,omitempty"`
	ResourceType string               `json:"resource_type,omitempty"`
	IP           string               `json:"ip,omitempty"`
	EventAt      time.Time            `json:"event_at,omitzero"`
}

// AuditAnomalyReason is a string enumeration used by AuditAnomaly.Reasons.
type AuditAnomalyReason string

// Values of AuditAnomalyReason. The server may add values; handle an unknown one gracefully.
const (
	AuditAnomalyReasonNewActionType AuditAnomalyReason = "new_action_type"
	AuditAnomalyReasonUnusualHour   AuditAnomalyReason = "unusual_hour"
	AuditAnomalyReasonNewIP         AuditAnomalyReason = "new_ip"
)

// AuditAnomalySummary is the result of AuditService.AnomalySummary.
type AuditAnomalySummary struct {
	HighScoreOpen int `json:"high_score_open"`
	Open          int `json:"open"`
	Resolved7d    int `json:"resolved_7d"`
}

// AuditAnomalySettings is the result of AuditService.AnomalySettings.
type AuditAnomalySettings struct {
	TenantID string `json:"tenant_id,omitempty"`
	Enabled  bool   `json:"enabled"`
	// 0..1; an event scoring at or above this is flagged.
	ScoreThreshold float64 `json:"score_threshold"`
	// Cold-start guard — an actor needs this many historical events before anything is flagged for them.
	MinBaselineEvents int `json:"min_baseline_events"`
}
