package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Webhook is a subscription. DisabledAt is nil while active; there is no
// separate "enabled" flag on the wire — a subscription is enabled exactly
// when DisabledAt is empty. There is also no rename/update endpoint —
// recreate the subscription to change its URL or events.
type Webhook struct {
	ID         string   `json:"id"`
	TenantID   string   `json:"tenant_id"`
	URL        string   `json:"url"`
	Events     []string `json:"events"`
	DisabledAt string   `json:"disabled_at,omitempty"`
	CreatedAt  string   `json:"created_at"`
	Secret     string   `json:"secret,omitempty"` // only present in the Create response
}

// Enabled reports whether the subscription is currently active.
func (w Webhook) Enabled() bool { return w.DisabledAt == "" }

type CreateWebhookInput struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

type WebhookDelivery struct {
	ID             string `json:"id"`
	WebhookID      string `json:"webhook_id"`
	Event          string `json:"event"`
	Status         string `json:"status"`
	ResponseStatus int    `json:"response_status,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// WebhooksService manages webhook subscriptions. Create/Get/Delete/Test/
// ListDeliveries/RetryDelivery are NOT tenant-path-scoped — the tenant is
// resolved from the caller's own API key, so the key must be tenant-scoped
// for these calls to work. Only List takes an explicit tenantID (the
// backend requires it in the path there, even though the same implicit
// scoping is available). See webhooks_verify.go for inbound-delivery HMAC
// verification helpers.
type WebhooksService struct{ t *transport.Transport }

func (s *WebhooksService) Create(ctx context.Context, in CreateWebhookInput) (*Webhook, error) {
	var out Webhook
	err := s.t.Post(ctx, "/v1/webhooks", in, &out)
	return &out, err
}

func (s *WebhooksService) Get(ctx context.Context, id string) (*Webhook, error) {
	var out Webhook
	err := s.t.Get(ctx, "/v1/webhooks/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

// Delete disables the subscription (the backend calls this "disable"
// internally; DELETE is the verb, matching every other resource's naming in
// this SDK).
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/webhooks/"+url.PathEscape(id), nil)
}

func (s *WebhooksService) Test(ctx context.Context, id string) error {
	return s.t.Post(ctx, "/v1/webhooks/"+url.PathEscape(id)+"/test", struct{}{}, nil)
}

func (s *WebhooksService) List(ctx context.Context, tenantID string) ([]Webhook, error) {
	var env marshal.Envelope[Webhook]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/webhooks", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *WebhooksService) ListDeliveries(ctx context.Context, webhookID string) ([]WebhookDelivery, error) {
	var env marshal.Envelope[WebhookDelivery]
	err := s.t.Get(ctx, "/v1/webhooks/"+url.PathEscape(webhookID)+"/deliveries", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *WebhooksService) RetryDelivery(ctx context.Context, webhookID, deliveryID string) error {
	path := "/v1/webhooks/" + url.PathEscape(webhookID) + "/deliveries/" + url.PathEscape(deliveryID) + "/retry"
	return s.t.Post(ctx, path, struct{}{}, nil)
}
