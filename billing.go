package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Plan is a subscribable Qeet ID plan. Prices map currency (ISO 4217) to
// amount in minor units (e.g. cents) — never a float.
type Plan struct {
	ID          string           `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Interval    string           `json:"interval"`
	Features    []string         `json:"features,omitempty"`
	Prices      map[string]int64 `json:"prices"` // currency -> minor units
}

// Subscription is a tenant's current Qeet ID billing state. Status "none"
// (all other fields zero-valued) means the tenant has no subscription.
type Subscription struct {
	PlanCode           string `json:"plan_code,omitempty"`
	PlanName           string `json:"plan_name,omitempty"`
	Currency           string `json:"currency,omitempty"`
	AmountMinor        int64  `json:"amount_minor,omitempty"`
	Interval           string `json:"interval,omitempty"`
	Status             string `json:"status"`
	CurrentPeriodStart string `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   string `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd  bool   `json:"cancel_at_period_end"`
}

type PutSubscriptionInput struct {
	PlanCode string `json:"plan_code"`
	Currency string `json:"currency"`
}

// Invoice is one billed period for a tenant's subscription.
type Invoice struct {
	ID          string `json:"id"`
	PlanCode    string `json:"plan_code"`
	Currency    string `json:"currency"`
	AmountMinor int64  `json:"amount_minor"`
	Status      string `json:"status"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
	IssuedAt    string `json:"issued_at"`
}

type CheckoutInput struct {
	PlanCode   string `json:"plan_code"`
	Currency   string `json:"currency"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

// CheckoutResult is either an immediate activation ("active", no
// CheckoutURL — e.g. a free plan) or a redirect to the payment provider
// ("checkout").
type CheckoutResult struct {
	Status      string `json:"status"` // active | checkout
	CheckoutURL string `json:"checkout_url,omitempty"`
	Provider    string `json:"provider,omitempty"`
}

// CancelSubscriptionResult confirms the subscription will lapse at the end
// of the current billing period (no proration/immediate cancellation today).
type CancelSubscriptionResult struct {
	CancelAtPeriodEnd bool `json:"cancel_at_period_end"`
}

// BillingService manages billing for the Qeet ID platform itself — what a
// tenant pays for using Qeet ID — not a general-purpose payments surface.
type BillingService struct{ t *transport.Transport }

func (s *BillingService) ListPlans(ctx context.Context) ([]Plan, error) {
	var env marshal.Envelope[Plan]
	err := s.t.Get(ctx, "/v1/billing/plans", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *BillingService) GetSubscription(ctx context.Context, tenantID string) (*Subscription, error) {
	var out Subscription
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/billing/subscription", transport.RequestOptions{}, &out)
	return &out, err
}

func (s *BillingService) PutSubscription(ctx context.Context, tenantID string, in PutSubscriptionInput) (*Subscription, error) {
	var out Subscription
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/billing/subscription", in, &out)
	return &out, err
}

func (s *BillingService) CancelSubscription(ctx context.Context, tenantID string) (*CancelSubscriptionResult, error) {
	var out CancelSubscriptionResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/billing/subscription/cancel", nil, &out)
	return &out, err
}

func (s *BillingService) ListInvoices(ctx context.Context, tenantID string) ([]Invoice, error) {
	var env marshal.Envelope[Invoice]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/billing/invoices", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *BillingService) Checkout(ctx context.Context, tenantID string, in CheckoutInput) (*CheckoutResult, error) {
	var out CheckoutResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/billing/checkout", in, &out)
	return &out, err
}
