package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// WebhooksService manages webhook subscriptions and their deliveries. To receive deliveries, see
// package github.com/qeetgroup/qeet-id-go/webhook.
type WebhooksService struct{ c *Client }

// Create creates a subscription. The result carries its signing secret, shown only once.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/webhooks
func (s *WebhooksService) Create(ctx context.Context, params CreateWebhookParams, opts ...RequestOption) (*Webhook, error) {
	return call[Webhook](ctx, s.c, &request{method: http.MethodPost, path: "/v1/webhooks", body: params, idempotent: true}, opts)
}

// Get returns a subscription.
//
// API: GET /v1/webhooks/{id}
func (s *WebhooksService) Get(ctx context.Context, webhookID string, opts ...RequestOption) (*Webhook, error) {
	return call[Webhook](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/webhooks/{webhookID}",
		pathArgs: []string{webhookID},
	}, opts)
}

// List lists a tenant's subscriptions.
//
// API: GET /v1/tenants/{tenantID}/webhooks
func (s *WebhooksService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Webhook], error) {
	return call[Page[Webhook]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/webhooks",
		pathArgs: []string{tenantID},
	}, opts)
}

// Update changes a subscription's URL or events. Empty fields are not changed.
//
// API: PATCH /v1/tenants/{tenantID}/webhooks/{id}
func (s *WebhooksService) Update(ctx context.Context, tenantID string, webhookID string, params UpdateWebhookParams, opts ...RequestOption) (*Webhook, error) {
	return call[Webhook](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/webhooks/{webhookID}",
		pathArgs: []string{tenantID, webhookID},
		body:     params,
	}, opts)
}

// Disable disables a subscription. Deliveries queue until Enable.
//
// API: DELETE /v1/webhooks/{id}
func (s *WebhooksService) Disable(ctx context.Context, webhookID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/webhooks/{webhookID}",
		pathArgs: []string{webhookID},
	}, nil, opts)
}

// Enable enables a disabled subscription and sends the deliveries that queued meanwhile.
//
// API: POST /v1/tenants/{tenantID}/webhooks/{id}/enable
func (s *WebhooksService) Enable(ctx context.Context, tenantID string, webhookID string, opts ...RequestOption) (*Webhook, error) {
	return call[Webhook](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/webhooks/{webhookID}/enable",
		pathArgs: []string{tenantID, webhookID},
	}, opts)
}

// RotateSecret replaces a subscription's signing secret. For 24 hours deliveries carry a signature for
// each secret.
//
// API: POST /v1/tenants/{tenantID}/webhooks/{id}/rotate-secret
func (s *WebhooksService) RotateSecret(ctx context.Context, tenantID string, webhookID string, opts ...RequestOption) (*RotatedWebhookSecret, error) {
	return call[RotatedWebhookSecret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/webhooks/{webhookID}/rotate-secret",
		pathArgs: []string{tenantID, webhookID},
	}, opts)
}

// Test queues a test.ping delivery to a subscription.
//
// API: POST /v1/webhooks/{id}/test
func (s *WebhooksService) Test(ctx context.Context, webhookID string, opts ...RequestOption) (*WebhookTest, error) {
	return call[WebhookTest](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/webhooks/{webhookID}/test",
		pathArgs: []string{webhookID},
	}, opts)
}

// ListEventTypes lists the events a subscription can name.
//
// API: GET /v1/tenants/{tenantID}/webhooks/events
func (s *WebhooksService) ListEventTypes(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[WebhookEventType], error) {
	return call[Page[WebhookEventType]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/webhooks/events",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListDeliveries lists a subscription's deliveries, newest first, one page at a time.
//
// API: GET /v1/webhooks/{id}/deliveries
func (s *WebhooksService) ListDeliveries(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) (*Page[WebhookDelivery], error) {
	return call[Page[WebhookDelivery]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/webhooks/{webhookID}/deliveries",
		pathArgs: []string{webhookID},
		query:    params.values(),
	}, opts)
}

// AllDeliveries returns an iterator over every item ListDeliveries returns, across all pages. It
// fetches a page only when the loop needs it, and breaking out of the loop stops paging. A failed
// request ends the iteration with a non-nil error.
func (s *WebhooksService) AllDeliveries(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) iter.Seq2[*WebhookDelivery, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[WebhookDelivery], error) {
		return s.ListDeliveries(ctx, webhookID, params.withCursor(cursor), opts...)
	})
}

// ListAttempts lists the attempts made to send one delivery.
//
// API: GET /v1/tenants/{tenantID}/webhooks/{id}/deliveries/{deliveryID}/attempts
func (s *WebhooksService) ListAttempts(ctx context.Context, tenantID string, webhookID string, deliveryID string, opts ...RequestOption) (*Page[WebhookAttempt], error) {
	return call[Page[WebhookAttempt]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/webhooks/{webhookID}/deliveries/{deliveryID}/attempts",
		pathArgs: []string{tenantID, webhookID, deliveryID},
	}, opts)
}

// RetryDelivery sends a delivery again.
//
// API: POST /v1/webhooks/{id}/deliveries/{deliveryID}/retry
func (s *WebhooksService) RetryDelivery(ctx context.Context, webhookID string, deliveryID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/webhooks/{webhookID}/deliveries/{deliveryID}/retry",
		pathArgs: []string{webhookID, deliveryID},
	}, nil, opts)
}

// RetryGivenUp sends every given-up delivery of a subscription again.
//
// API: POST /v1/tenants/{tenantID}/webhooks/{id}/deliveries/retry
func (s *WebhooksService) RetryGivenUp(ctx context.Context, tenantID string, webhookID string, opts ...RequestOption) (*RequeuedDeliveries, error) {
	return call[RequeuedDeliveries](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/webhooks/{webhookID}/deliveries/retry",
		pathArgs: []string{tenantID, webhookID},
	}, opts)
}

// CreateWebhookParams is the request body of WebhooksService.Create.
type CreateWebhookParams struct {
	// Event names from the event catalogue.
	Events []string `json:"events,omitempty"`
	// Ignored; the organisation is the caller's.
	TenantID string `json:"tenant_id,omitempty"`
	URL      string `json:"url,omitempty"`
}

// Webhook is a webhook subscription: an endpoint and the events it receives. Secret is set only when
// the subscription is created.
type Webhook struct {
	CreatedAt  time.Time  `json:"created_at,omitzero"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
	// Who disabled it, an administrator (disabled) or its failures (failing).
	DisabledReason WebhookDisabledReason `json:"disabled_reason,omitempty"`
	Events         []string              `json:"events,omitempty"`
	// When deliveries started failing, with no success since.
	FailingSince *time.Time `json:"failing_since,omitempty"`
	ID           string     `json:"id,omitempty"`
	// Until when the secret replaced by the last rotation also signs deliveries.
	PreviousSecretExpiresAt *time.Time `json:"previous_secret_expires_at,omitempty"`
	Secret                  string     `json:"secret,omitempty"`
	TenantID                string     `json:"tenant_id,omitempty"`
	URL                     string     `json:"url,omitempty"`
	// The administrator who disabled it.
	DisabledBy string `json:"disabled_by,omitempty"`
	// Its latest delivery.
	LastDeliveredAt *time.Time `json:"last_delivered_at,omitempty"`
	// The status its endpoint answered last.
	LastStatusCode *int `json:"last_status_code,omitempty"`
	// Its deliveries not yet tried (held while it is disabled), retrying, and given up.
	Counts WebhookCounts `json:"counts,omitzero"`
}

// WebhookDisabledReason is a string enumeration used by Webhook.DisabledReason. Who disabled it, an
// administrator (disabled) or its failures (failing).
type WebhookDisabledReason string

// Values of WebhookDisabledReason. The server may add values; handle an unknown one gracefully.
const (
	WebhookDisabledReasonDisabled WebhookDisabledReason = "disabled"
	WebhookDisabledReasonFailing  WebhookDisabledReason = "failing"
	WebhookDisabledReasonNone     WebhookDisabledReason = "None"
)

// WebhookCounts is the type of Webhook.Counts. Its deliveries not yet tried (held while it is
// disabled), retrying, and given up.
type WebhookCounts struct {
	Pending  int `json:"pending"`
	Retrying int `json:"retrying"`
	GaveUp   int `json:"gave_up"`
}

// UpdateWebhookParams is the request body of WebhooksService.Update.
type UpdateWebhookParams struct {
	URL    string   `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
}

// RotatedWebhookSecret is the result of WebhooksService.RotateSecret.
type RotatedWebhookSecret struct {
	PreviousSecretExpiresAt time.Time `json:"previous_secret_expires_at,omitzero"`
	Secret                  string    `json:"secret,omitempty"`
}

// WebhookTest is the result of WebhooksService.Test.
type WebhookTest struct {
	Queued       bool   `json:"queued"`
	Subscription string `json:"subscription,omitempty"`
}

// WebhookEventType names an event a subscription can receive.
type WebhookEventType struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

// WebhookDeliveryStatus is a string enumeration used by ListWebhookDeliveriesParams.Status.
type WebhookDeliveryStatus string

// Values of WebhookDeliveryStatus. The server may add values; handle an unknown one gracefully.
const (
	WebhookDeliveryStatusDelivered WebhookDeliveryStatus = "delivered"
	WebhookDeliveryStatusPending   WebhookDeliveryStatus = "pending"
	WebhookDeliveryStatusRetrying  WebhookDeliveryStatus = "retrying"
	WebhookDeliveryStatusDead      WebhookDeliveryStatus = "dead"
)

// ListWebhookDeliveriesParams holds the optional query parameters of WebhooksService.ListDeliveries. A
// nil *ListWebhookDeliveriesParams sends none.
type ListWebhookDeliveriesParams struct {
	Limit int
	// The previous page's next_cursor.
	Cursor string
	// delivered, pending (not tried yet), retrying or dead (given up).
	Status WebhookDeliveryStatus
}

func (p *ListWebhookDeliveriesParams) values() url.Values {
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

func (p *ListWebhookDeliveriesParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListWebhookDeliveriesParams) withCursor(cursor string) *ListWebhookDeliveriesParams {
	var c ListWebhookDeliveriesParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// WebhookDelivery is one event queued for a subscription, and the state of sending it.
type WebhookDelivery struct {
	Attempt       int       `json:"attempt"`
	CreatedAt     time.Time `json:"created_at,omitzero"`
	DeadAt        time.Time `json:"dead_at,omitzero"`
	DeliveredAt   time.Time `json:"delivered_at,omitzero"`
	Error         string    `json:"error,omitempty"`
	EventType     string    `json:"event_type,omitempty"`
	ID            string    `json:"id,omitempty"`
	NextAttemptAt time.Time `json:"next_attempt_at,omitzero"`
	Payload       string    `json:"payload,omitempty"`
	ResponseBody  string    `json:"response_body,omitempty"`
	StatusCode    int       `json:"status_code"`
}

// WebhookAttempt is one try at sending a delivery.
type WebhookAttempt struct {
	AttemptedAt time.Time `json:"attempted_at,omitzero"`
	// Null when the endpoint gave no answer.
	StatusCode *int   `json:"status_code,omitempty"`
	DurationMs int    `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}

// RequeuedDeliveries is the result of WebhooksService.RetryGivenUp.
type RequeuedDeliveries struct {
	Requeued int `json:"requeued"`
}
