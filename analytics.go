package qeetid

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AnalyticsService reads a tenant's sign-in and user analytics.
type AnalyticsService struct{ c *Client }

// Overview returns the tenant's key figures and trends.
//
// API: GET /v1/tenants/{tenantID}/analytics/overview
func (s *AnalyticsService) Overview(ctx context.Context, tenantID string, params *AnalyticsOverviewParams, opts ...RequestOption) (*AnalyticsOverview, error) {
	return call[AnalyticsOverview](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/analytics/overview",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AnalyticsOverviewParams holds the optional query parameters of AnalyticsService.Overview. A nil
// *AnalyticsOverviewParams sends none.
type AnalyticsOverviewParams struct {
	// Start of a chosen range (RFC 3339), returned as range.
	From time.Time
	// End of the range (RFC 3339); now by default.
	To time.Time
	// The range as days back from to; 30 by default, at most 366.
	Days int
	// Count only sign-ins of this first method (password, passkey, magic_link, social, saml, ...).
	Method string
}

func (p *AnalyticsOverviewParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.UTC().Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.UTC().Format(time.RFC3339))
	}
	if p.Days != 0 {
		q.Set("days", strconv.FormatInt(int64(p.Days), 10))
	}
	if p.Method != "" {
		q.Set("method", p.Method)
	}
	return q
}

// AnalyticsOverview holds a tenant's key sign-in and user figures and their trends.
type AnalyticsOverview struct {
	FailedLoginsHourly24h []AnalyticsHourlyPoint         `json:"failed_logins_hourly_24h,omitempty"`
	FailedTrend14d        []AnalyticsTrendPoint          `json:"failed_trend_14d,omitempty"`
	GeneratedAt           time.Time                      `json:"generated_at,omitzero"`
	KPIs                  AnalyticsOverviewKPIs          `json:"kpis,omitzero"`
	LoginActivity14d      []AnalyticsActivityPoint       `json:"login_activity_14d,omitempty"`
	LoginMethodsMix       []AnalyticsMethodShare         `json:"login_methods_mix,omitempty"`
	LoginTrend14d         []AnalyticsTrendPoint          `json:"login_trend_14d,omitempty"`
	MFAMethodsAdoption    []AnalyticsMethodCount         `json:"mfa_methods_adoption,omitempty"`
	MFATrend14d           []AnalyticsTrendPoint          `json:"mfa_trend_14d,omitempty"`
	UserTrend14d          []AnalyticsTrendPoint          `json:"user_trend_14d,omitempty"`
	WeeklyActivity8w      []AnalyticsWeeklyActivityPoint `json:"weekly_activity_8w,omitempty"`
	// The sign-ins day by day (UTC) over the chosen range, without organisation switches; present when
	// from, to or days was given.
	Range AnalyticsOverviewRange `json:"range,omitzero"`
}

// AnalyticsHourlyPoint is one element of AnalyticsOverview.FailedLoginsHourly24h.
type AnalyticsHourlyPoint struct {
	Attempts int    `json:"attempts"`
	Hour     string `json:"hour,omitempty"`
}

// AnalyticsTrendPoint is one element of AnalyticsOverview.FailedTrend14d.
type AnalyticsTrendPoint struct {
	Date  string  `json:"date,omitempty"`
	Value float64 `json:"value"`
}

// AnalyticsOverviewKPIs is the type of AnalyticsOverview.KPIs.
type AnalyticsOverviewKPIs struct {
	AvgSessionsPerUser AnalyticsMetric `json:"avg_sessions_per_user,omitzero"`
	DAU                AnalyticsMetric `json:"dau,omitzero"`
	FailedLogins24h    AnalyticsMetric `json:"failed_logins_24h,omitzero"`
	LoginsToday        AnalyticsMetric `json:"logins_today,omitzero"`
	MAU                AnalyticsMetric `json:"mau,omitzero"`
	MFAAdoptionPct     AnalyticsMetric `json:"mfa_adoption_pct,omitzero"`
	StickinessPct      AnalyticsMetric `json:"stickiness_pct,omitzero"`
	TotalUsers         AnalyticsMetric `json:"total_users,omitzero"`
}

// AnalyticsMetric is the type of AnalyticsOverviewKPIs.AvgSessionsPerUser.
type AnalyticsMetric struct {
	DeltaPct float64 `json:"delta_pct"`
	Value    float64 `json:"value"`
}

// AnalyticsActivityPoint is one element of AnalyticsOverview.LoginActivity14d.
type AnalyticsActivityPoint struct {
	Date      string `json:"date,omitempty"`
	LDAP      int    `json:"ldap"`
	MagicLink int    `json:"magic_link"`
	OIDC      int    `json:"oidc"`
	Passkey   int    `json:"passkey"`
	Password  int    `json:"password"`
	SAML      int    `json:"saml"`
	Social    int    `json:"social"`
}

// AnalyticsMethodShare is one element of AnalyticsOverview.LoginMethodsMix.
type AnalyticsMethodShare struct {
	Method string  `json:"method,omitempty"`
	Value  float64 `json:"value"`
}

// AnalyticsMethodCount is one element of AnalyticsOverview.MFAMethodsAdoption.
type AnalyticsMethodCount struct {
	Method string `json:"method,omitempty"`
	Users  int    `json:"users"`
}

// AnalyticsWeeklyActivityPoint is one element of AnalyticsOverview.WeeklyActivity8w.
type AnalyticsWeeklyActivityPoint struct {
	DAU  int    `json:"dau"`
	WAU  int    `json:"wau"`
	Week string `json:"week,omitempty"`
}

// AnalyticsOverviewRange is the sign-ins day by day (UTC) over the chosen range, without organisation
// switches; present when from, to or days was given.
type AnalyticsOverviewRange struct {
	From   time.Time                   `json:"from,omitzero"`
	To     time.Time                   `json:"to,omitzero"`
	Method string                      `json:"method,omitempty"`
	Days   []AnalyticsOverviewRangeDay `json:"days,omitempty"`
}

// AnalyticsOverviewRangeDay is one element of AnalyticsOverviewRange.Days.
type AnalyticsOverviewRangeDay struct {
	Date    string `json:"date,omitempty"`
	SignIns int    `json:"sign_ins"`
	People  int    `json:"people"`
	Failed  int    `json:"failed"`
}
