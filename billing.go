package qeetid

import (
	"context"
	"net/http"
	"time"
)

// BillingService manages the tenant's Qeet ID plan, subscription, invoices and billing profile.
type BillingService struct{ c *Client }

// ListPlans lists the plans a tenant can subscribe to.
//
// API: GET /v1/billing/plans
func (s *BillingService) ListPlans(ctx context.Context, opts ...RequestOption) (*Page[Plan], error) {
	return call[Page[Plan]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/billing/plans"}, opts)
}

// Subscription returns the tenant's subscription. Status is none when there is no subscription.
//
// API: GET /v1/tenants/{tenantID}/billing/subscription
func (s *BillingService) Subscription(ctx context.Context, tenantID string, opts ...RequestOption) (*Subscription, error) {
	return call[Subscription](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/billing/subscription",
		pathArgs: []string{tenantID},
	}, opts)
}

// Subscribe subscribes the tenant to a plan or changes its plan.
//
// API: PUT /v1/tenants/{tenantID}/billing/subscription
func (s *BillingService) Subscribe(ctx context.Context, tenantID string, params SubscribeParams, opts ...RequestOption) (*Subscription, error) {
	return call[Subscription](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/billing/subscription",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// CancelSubscription cancels the subscription at the end of the current period.
//
// API: POST /v1/tenants/{tenantID}/billing/subscription/cancel
func (s *BillingService) CancelSubscription(ctx context.Context, tenantID string, opts ...RequestOption) (*SubscriptionCancellation, error) {
	return call[SubscriptionCancellation](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/billing/subscription/cancel",
		pathArgs: []string{tenantID},
	}, opts)
}

// StartTrial starts a trial of a plan for the tenant.
//
// API: POST /v1/tenants/{tenantID}/billing/trial
func (s *BillingService) StartTrial(ctx context.Context, tenantID string, params StartTrialParams, opts ...RequestOption) (*Subscription, error) {
	return call[Subscription](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/billing/trial",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Checkout starts a checkout for a paid plan. Send the person to CheckoutURL when Status is checkout.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/tenants/{tenantID}/billing/checkout
func (s *BillingService) Checkout(ctx context.Context, tenantID string, params CheckoutParams, opts ...RequestOption) (*CheckoutResult, error) {
	return call[CheckoutResult](ctx, s.c, &request{
		method:     http.MethodPost,
		path:       "/v1/tenants/{tenantID}/billing/checkout",
		pathArgs:   []string{tenantID},
		body:       params,
		idempotent: true,
	}, opts)
}

// ListInvoices lists the tenant's invoices.
//
// API: GET /v1/tenants/{tenantID}/billing/invoices
func (s *BillingService) ListInvoices(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Invoice], error) {
	return call[Page[Invoice]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/billing/invoices",
		pathArgs: []string{tenantID},
	}, opts)
}

// PayInvoice starts a payment of an open invoice.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/tenants/{tenantID}/billing/invoices/{invoiceID}/pay
func (s *BillingService) PayInvoice(ctx context.Context, tenantID, invoiceID string, params PayInvoiceParams, opts ...RequestOption) (*CheckoutResult, error) {
	return call[CheckoutResult](ctx, s.c, &request{
		method:     http.MethodPost,
		path:       "/v1/tenants/{tenantID}/billing/invoices/{invoiceID}/pay",
		pathArgs:   []string{tenantID, invoiceID},
		body:       params,
		idempotent: true,
	}, opts)
}

// PaymentMethod returns the payment method on file, if any.
//
// API: GET /v1/tenants/{tenantID}/billing/payment-method
func (s *BillingService) PaymentMethod(ctx context.Context, tenantID string, opts ...RequestOption) (*PaymentMethodResult, error) {
	return call[PaymentMethodResult](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/billing/payment-method",
		pathArgs: []string{tenantID},
	}, opts)
}

// Profile returns the tenant's billing profile.
//
// API: GET /v1/tenants/{tenantID}/billing/profile
func (s *BillingService) Profile(ctx context.Context, tenantID string, opts ...RequestOption) (*BillingProfile, error) {
	return call[BillingProfile](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/billing/profile",
		pathArgs: []string{tenantID},
	}, opts)
}

// UpdateProfile replaces the tenant's billing profile.
//
// API: PUT /v1/tenants/{tenantID}/billing/profile
func (s *BillingService) UpdateProfile(ctx context.Context, tenantID string, params BillingProfile, opts ...RequestOption) (*BillingProfile, error) {
	return call[BillingProfile](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/billing/profile",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Plan is a Qeet ID plan a tenant can subscribe to. Prices are in minor units, such as paise or cents.
type Plan struct {
	Code        string         `json:"code,omitempty"`
	Description string         `json:"description,omitempty"`
	Features    []string       `json:"features,omitempty"`
	ID          string         `json:"id,omitempty"`
	Interval    string         `json:"interval,omitempty"`
	Name        string         `json:"name,omitempty"`
	Prices      map[string]int `json:"prices,omitempty"`
}

// Subscription is a tenant's Qeet ID subscription. Status is none when the tenant has none.
type Subscription struct {
	AmountMinor        int                `json:"amount_minor"`
	CancelAtPeriodEnd  bool               `json:"cancel_at_period_end"`
	Currency           string             `json:"currency,omitempty"`
	CurrentPeriodEnd   *time.Time         `json:"current_period_end,omitempty"`
	CurrentPeriodStart *time.Time         `json:"current_period_start,omitempty"`
	Interval           string             `json:"interval,omitempty"`
	PlanCode           string             `json:"plan_code,omitempty"`
	PlanName           string             `json:"plan_name,omitempty"`
	Status             SubscriptionStatus `json:"status,omitempty"`
	TrialEnd           *time.Time         `json:"trial_end,omitempty"`
	// When a canceled subscription ended.
	EndedAt *time.Time `json:"ended_at,omitempty"`
}

// SubscriptionStatus is a string enumeration used by Subscription.Status.
type SubscriptionStatus string

// Values of SubscriptionStatus. The server may add values; handle an unknown one gracefully.
const (
	SubscriptionStatusNone     SubscriptionStatus = "none"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

// SubscribeParams is the request body of BillingService.Subscribe.
type SubscribeParams struct {
	Currency string `json:"currency,omitempty"`
	PlanCode string `json:"plan_code,omitempty"`
}

// SubscriptionCancellation is the result of BillingService.CancelSubscription.
type SubscriptionCancellation struct {
	CancelAtPeriodEnd bool `json:"cancel_at_period_end"`
}

// StartTrialParams is the request body of BillingService.StartTrial.
type StartTrialParams struct {
	// The plan to try.
	PlanCode string `json:"plan_code,omitempty"`
	// ISO 4217 currency code, such as INR or USD.
	Currency string `json:"currency,omitempty"`
}

// CheckoutParams is the request body of BillingService.Checkout.
type CheckoutParams struct {
	PlanCode string `json:"plan_code,omitempty"`
	Currency string `json:"currency,omitempty"`
	// ISO-3166-1 alpha-2 billing country; selects the payment provider. Omitted → currency-based routing.
	Country    string `json:"country,omitempty"`
	SuccessURL string `json:"success_url,omitempty"`
	CancelURL  string `json:"cancel_url,omitempty"`
}

// CheckoutResult is the outcome of starting a checkout. Send the person to CheckoutURL when Status is
// checkout.
type CheckoutResult struct {
	CheckoutURL string `json:"checkout_url,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Status      string `json:"status,omitempty"`
}

// Invoice is a bill for a tenant's subscription.
type Invoice struct {
	// The total charged, taxable_amount_minor plus tax_amount_minor; for a checkout, what the payment
	// provider charged.
	AmountMinor int    `json:"amount_minor"`
	Currency    string `json:"currency,omitempty"`
	ID          string `json:"id,omitempty"`
	// The GST invoice number, one series per Indian financial year (April to March), such as
	// 2026-27/000042. Null on an invoice that charges nothing.
	InvoiceNumber string    `json:"invoice_number,omitempty"`
	IssuedAt      time.Time `json:"issued_at,omitzero"`
	PeriodEnd     time.Time `json:"period_end,omitzero"`
	PeriodStart   time.Time `json:"period_start,omitzero"`
	PlaceOfSupply string    `json:"place_of_supply,omitempty"`
	PlanCode      string    `json:"plan_code,omitempty"`
	// paid when a payment was captured or nothing is owed; open until it is paid, due at due_at, for a
	// plan granted before payment (invoice-only billing) or a renewal; void for a renewal no longer owed,
	// because the subscription ended or a new plan was bought.
	Status InvoiceStatus `json:"status,omitempty"`
	// purchase for a plan bought or changed; renewal for the next period of a plan, billed when its
	// period ends.
	Kind InvoiceKind `json:"kind,omitempty"`
	// When an open invoice falls due; null once it is paid or when nothing was owed.
	DueAt  *time.Time `json:"due_at,omitempty"`
	PaidAt *time.Time `json:"paid_at,omitempty"`
	// The billing profile as it stood when the invoice was issued; null on invoices issued before it was
	// kept.
	BilledTo *BillingProfile `json:"billed_to,omitempty"`
	// The seller's GSTIN; empty when the operator has set none.
	SupplierGSTIN string `json:"supplier_gstin,omitempty"`
	// The services accounting code of what was sold; empty when the operator has set none.
	Sac string `json:"sac,omitempty"`
	// When the invoice was emailed to its billing email.
	EmailedAt          *time.Time `json:"emailed_at,omitempty"`
	TaxAmountMinor     int        `json:"tax_amount_minor"`
	TaxRateBps         int        `json:"tax_rate_bps"`
	TaxType            string     `json:"tax_type,omitempty"`
	TaxableAmountMinor int        `json:"taxable_amount_minor"`
}

// InvoiceStatus is a string enumeration used by Invoice.Status. paid when a payment was captured or
// nothing is owed; open until it is paid, due at due_at, for a plan granted before payment
// (invoice-only billing) or a renewal; void for a renewal no longer owed, because the subscription
// ended or a new plan was bought.
type InvoiceStatus string

// Values of InvoiceStatus. The server may add values; handle an unknown one gracefully.
const (
	InvoiceStatusPaid InvoiceStatus = "paid"
	InvoiceStatusOpen InvoiceStatus = "open"
	InvoiceStatusVoid InvoiceStatus = "void"
)

// InvoiceKind is a string enumeration used by Invoice.Kind. purchase for a plan bought or changed;
// renewal for the next period of a plan, billed when its period ends.
type InvoiceKind string

// Values of InvoiceKind. The server may add values; handle an unknown one gracefully.
const (
	InvoiceKindPurchase InvoiceKind = "purchase"
	InvoiceKindRenewal  InvoiceKind = "renewal"
)

// BillingProfile holds the details printed on a tenant's invoices.
type BillingProfile struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	BillingEmail string `json:"billing_email,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	LegalName    string `json:"legal_name,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	State        string `json:"state,omitempty"`
	TaxID        string `json:"tax_id,omitempty"`
	TaxIDType    string `json:"tax_id_type,omitempty"`
}

// PayInvoiceParams is the request body of BillingService.PayInvoice.
type PayInvoiceParams struct {
	SuccessURL string `json:"success_url,omitempty"`
	CancelURL  string `json:"cancel_url,omitempty"`
}

// PaymentMethodResult is the result of BillingService.PaymentMethod.
type PaymentMethodResult struct {
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
}

// PaymentMethod is the type of PaymentMethodResult.PaymentMethod.
type PaymentMethod struct {
	// stripe, razorpay or sandbox.
	Provider string `json:"provider,omitempty"`
	// card, upi, netbanking, wallet, or the provider's own word.
	Type string `json:"type,omitempty"`
	// The card network, bank or wallet.
	Brand      string    `json:"brand,omitempty"`
	Last4      string    `json:"last4,omitempty"`
	ExpMonth   int       `json:"exp_month"`
	ExpYear    int       `json:"exp_year"`
	LastUsedAt time.Time `json:"last_used_at,omitzero"`
}
