package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// AnalyticsOverview is the single payload behind the admin dashboard's KPI
// cards and charts. Where the underlying data isn't recorded yet, fields
// come back as empty buckets rather than being omitted.
type AnalyticsOverview struct {
	GeneratedAt string `json:"generated_at"`

	KPIs struct {
		MAU                AnalyticsMetric `json:"mau"`
		LoginsToday        AnalyticsMetric `json:"logins_today"`
		MFAAdoptionPct     AnalyticsMetric `json:"mfa_adoption_pct"`
		FailedLogins24h    AnalyticsMetric `json:"failed_logins_24h"`
		DAU                AnalyticsMetric `json:"dau"`
		TotalUsers         AnalyticsMetric `json:"total_users"`
		AvgSessionsPerUser AnalyticsMetric `json:"avg_sessions_per_user"`
		StickinessPct      AnalyticsMetric `json:"stickiness_pct"`
	} `json:"kpis"`

	WeeklyActivity8w []AnalyticsWeeklyActivityPoint `json:"weekly_activity_8w"`

	UserTrend14d   []AnalyticsTrendPoint `json:"user_trend_14d"`
	LoginTrend14d  []AnalyticsTrendPoint `json:"login_trend_14d"`
	MFATrend14d    []AnalyticsTrendPoint `json:"mfa_trend_14d"`
	FailedTrend14d []AnalyticsTrendPoint `json:"failed_trend_14d"`

	LoginActivity14d      []AnalyticsActivityPoint `json:"login_activity_14d"`
	LoginMethodsMix       []AnalyticsMethodSlice   `json:"login_methods_mix"`
	MFAMethodsAdoption    []AnalyticsMethodCount   `json:"mfa_methods_adoption"`
	FailedLoginsHourly24h []AnalyticsHourlyPoint   `json:"failed_logins_hourly_24h"`
}

type AnalyticsMetric struct {
	Value    float64 `json:"value"`
	DeltaPct float64 `json:"delta_pct"`
}

type AnalyticsTrendPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// AnalyticsActivityPoint is a daily bucket grouped by login method. Missing
// methods come back as 0, not omitted, so a stacked-area chart never gaps.
type AnalyticsActivityPoint struct {
	Date     string `json:"date"`
	Password int64  `json:"password"`
	Passkey  int64  `json:"passkey"`
	Social   int64  `json:"social"`
	SAML     int64  `json:"saml"`
	OIDC     int64  `json:"oidc"`
}

type AnalyticsMethodSlice struct {
	Method string  `json:"method"`
	Value  float64 `json:"value"`
}

type AnalyticsMethodCount struct {
	Method string `json:"method"`
	Users  int64  `json:"users"`
}

type AnalyticsHourlyPoint struct {
	Hour     string `json:"hour"`
	Attempts int64  `json:"attempts"`
}

// AnalyticsWeeklyActivityPoint is one ISO week's WAU/average-DAU bucket.
type AnalyticsWeeklyActivityPoint struct {
	Week string `json:"week"` // "Wnn"
	WAU  int64  `json:"wau"`
	DAU  int64  `json:"dau"`
}

// AnalyticsService reads the admin dashboard's aggregated KPIs.
type AnalyticsService struct{ t *transport.Transport }

func (s *AnalyticsService) Overview(ctx context.Context, tenantID string) (*AnalyticsOverview, error) {
	var out AnalyticsOverview
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/analytics/overview", transport.RequestOptions{}, &out)
	return &out, err
}
