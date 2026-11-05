package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ThreatsService manages bot detection, risk-based sign-in settings and security anomalies.
type ThreatsService struct{ c *Client }

// Bots returns recent bot-scoring verdicts and their totals, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/security/bots
func (s *ThreatsService) Bots(ctx context.Context, tenantID string, params *BotTelemetryParams, opts ...RequestOption) (*BotTelemetry, error) {
	return call[BotTelemetry](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/security/bots",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// BotSettings returns a tenant's bot-detection settings.
//
// API: GET /v1/tenants/{tenantID}/security/bots/settings
func (s *ThreatsService) BotSettings(ctx context.Context, tenantID string, opts ...RequestOption) (*BotSettings, error) {
	return call[BotSettings](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/security/bots/settings",
		pathArgs: []string{tenantID},
	}, opts)
}

// UpdateBotSettings replaces a tenant's bot-detection settings.
//
// API: PUT /v1/tenants/{tenantID}/security/bots/settings
func (s *ThreatsService) UpdateBotSettings(ctx context.Context, tenantID string, params BotSettings, opts ...RequestOption) (*BotSettings, error) {
	return call[BotSettings](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/security/bots/settings",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// RiskSettings returns a tenant's risk-based sign-in settings.
//
// API: GET /v1/tenants/{tenantID}/security/risk-settings
func (s *ThreatsService) RiskSettings(ctx context.Context, tenantID string, opts ...RequestOption) (*RiskSettings, error) {
	return call[RiskSettings](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/security/risk-settings",
		pathArgs: []string{tenantID},
	}, opts)
}

// UpdateRiskSettings replaces a tenant's risk-based sign-in settings.
//
// API: PUT /v1/tenants/{tenantID}/security/risk-settings
func (s *ThreatsService) UpdateRiskSettings(ctx context.Context, tenantID string, params RiskSettings, opts ...RequestOption) (*RiskSettings, error) {
	return call[RiskSettings](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/security/risk-settings",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// ListAnomalies lists detected sign-in anomalies, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/security/anomalies
func (s *ThreatsService) ListAnomalies(ctx context.Context, tenantID string, params *ListSecurityAnomaliesParams, opts ...RequestOption) (*Page[SecurityAnomaly], error) {
	return call[Page[SecurityAnomaly]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/security/anomalies",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AllAnomalies returns an iterator over every item ListAnomalies returns, across all pages. It fetches
// a page only when the loop needs it, and breaking out of the loop stops paging. A failed request ends
// the iteration with a non-nil error.
func (s *ThreatsService) AllAnomalies(ctx context.Context, tenantID string, params *ListSecurityAnomaliesParams, opts ...RequestOption) iter.Seq2[*SecurityAnomaly, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[SecurityAnomaly], error) {
		return s.ListAnomalies(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// AnomalySummary returns counts of open and resolved anomalies.
//
// API: GET /v1/tenants/{tenantID}/security/anomalies/summary
func (s *ThreatsService) AnomalySummary(ctx context.Context, tenantID string, opts ...RequestOption) (*SecurityAnomalySummary, error) {
	return call[SecurityAnomalySummary](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/security/anomalies/summary",
		pathArgs: []string{tenantID},
	}, opts)
}

// ResolveAnomaly marks an anomaly resolved.
//
// API: POST /v1/tenants/{tenantID}/security/anomalies/{id}/resolve
func (s *ThreatsService) ResolveAnomaly(ctx context.Context, tenantID string, anomalyID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/security/anomalies/{anomalyID}/resolve",
		pathArgs: []string{tenantID, anomalyID},
	}, nil, opts)
}

// BotVerdict is a string enumeration used by BotTelemetryParams.Verdict.
type BotVerdict string

// Values of BotVerdict. The server may add values; handle an unknown one gracefully.
const (
	BotVerdictAllowed    BotVerdict = "allowed"
	BotVerdictChallenged BotVerdict = "challenged"
	BotVerdictBlocked    BotVerdict = "blocked"
)

// BotTelemetryParams holds the optional query parameters of ThreatsService.Bots. A nil
// *BotTelemetryParams sends none.
type BotTelemetryParams struct {
	Limit int
	// The previous page's next_cursor.
	Cursor  string
	Verdict BotVerdict
}

func (p *BotTelemetryParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	if p.Verdict != "" {
		q.Set("verdict", string(p.Verdict))
	}
	return q
}

// BotTelemetry is the result of ThreatsService.Bots.
type BotTelemetry struct {
	// A page of detections, newest first.
	Recent []BotEvent `json:"recent,omitempty"`
	// Present when there is a next page.
	NextCursor string `json:"next_cursor,omitempty"`
	// How many detections match.
	Total int      `json:"total"`
	Stats BotStats `json:"stats,omitzero"`
}

// BotEvent is one element of BotTelemetry.Recent.
type BotEvent struct {
	// The account the sign-in attempt was for, when it exists.
	UserID    string `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	// How that account's next sign-in from the same address, within five minutes, ended.
	Outcome   BotEventOutcome `json:"outcome,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitzero"`
	ID        string          `json:"id,omitempty"`
	IP        string          `json:"ip,omitempty"`
	Score     float64         `json:"score"`
	UserAgent string          `json:"user_agent,omitempty"`
	Verdict   string          `json:"verdict,omitempty"`
}

// BotEventOutcome is a string enumeration used by BotEvent.Outcome. How that account's next sign-in
// from the same address, within five minutes, ended.
type BotEventOutcome string

// Values of BotEventOutcome. The server may add values; handle an unknown one gracefully.
const (
	BotEventOutcomeSignedIn BotEventOutcome = "signed_in"
	BotEventOutcomeFailed   BotEventOutcome = "failed"
	BotEventOutcomeUnknown  BotEventOutcome = "unknown"
)

// BotStats is the type of BotTelemetry.Stats.
type BotStats struct {
	Blocked24h    int     `json:"blocked_24h"`
	Challenged24h int     `json:"challenged_24h"`
	Threshold     float64 `json:"threshold"`
}

// BotSettings configures how a tenant scores sign-in traffic for bots.
type BotSettings struct {
	Captcha        bool    `json:"captcha"`
	Honeypot       bool    `json:"honeypot"`
	ScoreThreshold float64 `json:"score_threshold"`
	Signature      bool    `json:"signature"`
	UACheck        bool    `json:"ua_check"`
}

// RiskSettings configures when a sign-in counts as risky and must pass a second factor.
type RiskSettings struct {
	MediumThreshold float64                     `json:"medium_threshold"`
	HighThreshold   float64                     `json:"high_threshold"`
	ForceMFAAtLevel RiskSettingsForceMFAAtLevel `json:"force_mfa_at_level,omitempty"`
	// Off by default. Also needs a geo signal to do anything — see GEO_COUNTRY_HEADER.
	ImpossibleTravelEnabled bool `json:"impossible_travel_enabled"`
	// A login from a different country than the last one, sooner than this many hours later, is flagged.
	MinTravelHours float64 `json:"min_travel_hours"`
	// Off by default. Flags a login from a browser+OS combination never seen before for this user.
	DeviceReputationEnabled bool `json:"device_reputation_enabled"`
	// Whether the deployment names a trusted proxy country header (GEO_COUNTRY_HEADER). Without one
	// impossible travel never fires. Ignored on update.
	CountryHeaderConfigured bool `json:"country_header_configured"`
	// When the instance answering last read a country from that header; null if none since it started (no
	// header, or the proxy isn't sending it). Ignored on update.
	CountryLastSeenAt *time.Time `json:"country_last_seen_at,omitempty"`
}

// RiskSettingsForceMFAAtLevel is a string enumeration used by RiskSettings.ForceMFAAtLevel.
type RiskSettingsForceMFAAtLevel string

// Values of RiskSettingsForceMFAAtLevel. The server may add values; handle an unknown one gracefully.
const (
	RiskSettingsForceMFAAtLevelMedium RiskSettingsForceMFAAtLevel = "medium"
	RiskSettingsForceMFAAtLevelHigh   RiskSettingsForceMFAAtLevel = "high"
)

// SecurityAnomalyStatus is a string enumeration used by ListSecurityAnomaliesParams.Status.
type SecurityAnomalyStatus string

// Values of SecurityAnomalyStatus. The server may add values; handle an unknown one gracefully.
const (
	SecurityAnomalyStatusOpen     SecurityAnomalyStatus = "open"
	SecurityAnomalyStatusResolved SecurityAnomalyStatus = "resolved"
)

// ListSecurityAnomaliesParams holds the optional query parameters of ThreatsService.ListAnomalies. A
// nil *ListSecurityAnomaliesParams sends none.
type ListSecurityAnomaliesParams struct {
	Limit int
	// The previous page's next_cursor.
	Cursor string
	Status SecurityAnomalyStatus
}

func (p *ListSecurityAnomaliesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	if p.Status != "" {
		q.Set("status", string(p.Status))
	}
	return q
}

func (p *ListSecurityAnomaliesParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListSecurityAnomaliesParams) withCursor(cursor string) *ListSecurityAnomaliesParams {
	var c ListSecurityAnomaliesParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// SecurityAnomaly is a suspicious sign-in pattern Qeet ID detected.
type SecurityAnomaly struct {
	// The person who resolved it; absent while open.
	ResolvedBy      string `json:"resolved_by,omitempty"`
	ResolvedByEmail string `json:"resolved_by_email,omitempty"`
	ID              string `json:"id,omitempty"`
	// e.g. credential_stuffing (an account locked after repeated failed sign-ins) or password_spraying (a
	// source network throttled for failing against many accounts).
	Type      string `json:"type,omitempty"`
	Severity  string `json:"severity,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Status    string `json:"status,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	IP        string `json:"ip,omitempty"`
	// The detection's structured context. For a sign-in lock: reason, account, network, counter
	// (account_network or account), failures, threshold, window_minutes, locked_until. For a throttled
	// source: reason, network, counter (source), distinct_accounts, window_minutes, throttled_until.
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitzero"`
	ResolvedAt time.Time      `json:"resolved_at,omitzero"`
}

// SecurityAnomalySummary is the result of ThreatsService.AnomalySummary.
type SecurityAnomalySummary struct {
	// Distinct accounts of the organisation with an open event, counting the account an event names and
	// every account a password spray targeted (a password_spraying event lists them in
	// metadata.targeted_user_ids, only the organisation's own).
	AffectedAccounts int `json:"affected_accounts"`
	HighSeverity24h  int `json:"high_severity_24h"`
	Open             int `json:"open"`
	Resolved24h      int `json:"resolved_24h"`
}
