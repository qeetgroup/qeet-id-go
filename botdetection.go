package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// BotEvent is a single recorded bot-scoring verdict (offline User-Agent
// heuristic).
type BotEvent struct {
	ID        string  `json:"id"`
	IP        string  `json:"ip,omitempty"`
	UserAgent string  `json:"user_agent"`
	Score     float64 `json:"score"`   // offline UA heuristic, 0.0–1.0 (discrete buckets)
	Verdict   string  `json:"verdict"` // allowed | challenged | blocked
	CreatedAt string  `json:"created_at"`
}

// BotDetectionStats reports blocked/challenged counts over the last 24h plus
// the tenant's current score threshold.
type BotDetectionStats struct {
	Blocked24h    int     `json:"blocked_24h"`
	Challenged24h int     `json:"challenged_24h"`
	Threshold     float64 `json:"threshold"`
}

// BotDetectionOverview is the bot-detection overview: the most recent recorded
// verdicts plus aggregate stats.
type BotDetectionOverview struct {
	Recent []BotEvent        `json:"recent"`
	Stats  BotDetectionStats `json:"stats"`
}

type BotDetectionSettings struct {
	UACheck        bool    `json:"ua_check"`
	Honeypot       bool    `json:"honeypot"`
	Captcha        bool    `json:"captcha"`
	Signature      bool    `json:"signature"`
	ScoreThreshold float64 `json:"score_threshold"`
}

// UpdateBotDetectionSettingsInput is a full replace, not a partial patch: the
// backend decodes it directly into its settings record (no pointer/optional
// fields), so any omitted field is written back as its zero value rather than
// left unchanged.
type UpdateBotDetectionSettingsInput struct {
	UACheck        bool    `json:"ua_check"`
	Honeypot       bool    `json:"honeypot"`
	Captcha        bool    `json:"captcha"`
	Signature      bool    `json:"signature"`
	ScoreThreshold float64 `json:"score_threshold"`
}

// BotDetectionService exposes tenant-scoped bot-detection settings and stats —
// one of the threat-detection surfaces alongside RiskSettingsService.
// Detection is detect-only: a "blocked" verdict means "would block", so this
// never affects the auth path itself, only what's surfaced here. There is no
// pagination or list — Overview is a single aggregate read and settings are a
// singleton get/put record, the same shape as AuthPolicyService/PolicyService.
type BotDetectionService struct{ t *transport.Transport }

// Overview returns the most recent recorded verdicts plus aggregate stats.
func (s *BotDetectionService) Overview(ctx context.Context, tenantID string) (*BotDetectionOverview, error) {
	var out BotDetectionOverview
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/security/bots", transport.RequestOptions{}, &out)
	return &out, err
}

// GetSettings returns the tenant's current bot-detection settings.
func (s *BotDetectionService) GetSettings(ctx context.Context, tenantID string) (*BotDetectionSettings, error) {
	var out BotDetectionSettings
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/security/bots/settings", transport.RequestOptions{}, &out)
	return &out, err
}

// PutSettings replaces the tenant's bot-detection settings (full replace).
func (s *BotDetectionService) PutSettings(ctx context.Context, tenantID string, in UpdateBotDetectionSettingsInput) (*BotDetectionSettings, error) {
	var out BotDetectionSettings
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/security/bots/settings", in, &out)
	return &out, err
}
