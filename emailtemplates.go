package qeetid

import (
	"context"
	"net/http"
)

// EmailTemplatesService customizes the emails Qeet ID sends for a tenant.
type EmailTemplatesService struct{ c *Client }

// List lists every template with the tenant's customization, if any.
//
// API: GET /v1/tenants/{tenantID}/email-templates
func (s *EmailTemplatesService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*EmailTemplates, error) {
	return call[EmailTemplates](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/email-templates",
		pathArgs: []string{tenantID},
	}, opts)
}

// Get returns a template.
//
// API: GET /v1/tenants/{tenantID}/email-templates/{key}
func (s *EmailTemplatesService) Get(ctx context.Context, tenantID, key string, opts ...RequestOption) (*EmailTemplate, error) {
	return call[EmailTemplate](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/email-templates/{key}",
		pathArgs: []string{tenantID, key},
	}, opts)
}

// Update saves a tenant's customization of a template.
//
// API: PUT /v1/tenants/{tenantID}/email-templates/{key}
func (s *EmailTemplatesService) Update(ctx context.Context, tenantID, key string, params UpdateEmailTemplateParams, opts ...RequestOption) (*EmailTemplate, error) {
	return call[EmailTemplate](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/email-templates/{key}",
		pathArgs: []string{tenantID, key},
		body:     params,
	}, opts)
}

// Reset drops a tenant's customization and returns the default template.
//
// API: DELETE /v1/tenants/{tenantID}/email-templates/{key}
func (s *EmailTemplatesService) Reset(ctx context.Context, tenantID, key string, opts ...RequestOption) (*EmailTemplate, error) {
	return call[EmailTemplate](ctx, s.c, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/email-templates/{key}",
		pathArgs: []string{tenantID, key},
	}, opts)
}

// Preview renders a template with sample values. Nothing is sent.
//
// API: POST /v1/tenants/{tenantID}/email-templates/{key}/preview
func (s *EmailTemplatesService) Preview(ctx context.Context, tenantID, key string, params PreviewEmailTemplateParams, opts ...RequestOption) (*EmailPreview, error) {
	return call[EmailPreview](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/email-templates/{key}/preview",
		pathArgs: []string{tenantID, key},
		body:     params,
	}, opts)
}

// SendTest sends a template to the calling person's verified email address.
//
// API: POST /v1/tenants/{tenantID}/email-templates/{key}/test
func (s *EmailTemplatesService) SendTest(ctx context.Context, tenantID, key string, params SendTestEmailParams, opts ...RequestOption) (*TestEmailResult, error) {
	return call[TestEmailResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/email-templates/{key}/test",
		pathArgs: []string{tenantID, key},
		body:     params,
	}, opts)
}

// EmailTemplates is the result of EmailTemplatesService.List.
type EmailTemplates struct {
	Items        []EmailTemplate            `json:"items,omitempty"`
	Capabilities EmailTemplatesCapabilities `json:"capabilities,omitzero"`
}

// EmailTemplate is an email Qeet ID sends, with the tenant's customization if it has one.
type EmailTemplate struct {
	Subject       string                     `json:"subject,omitempty"`
	Body          string                     `json:"body,omitempty"`
	Options       EmailTemplateOptions       `json:"options,omitzero"`
	Key           string                     `json:"key,omitempty"`
	Name          string                     `json:"name,omitempty"`
	Description   string                     `json:"description,omitempty"`
	Variables     []string                   `json:"variables,omitempty"`
	Custom        bool                       `json:"custom"`
	DeliveryScope EmailTemplateDeliveryScope `json:"delivery_scope,omitempty"`
}

// EmailTemplateOptions is the type of EmailTemplate.Options.
type EmailTemplateOptions struct {
	Preheader     string                 `json:"preheader,omitempty"`
	BodyHTML      string                 `json:"body_html,omitempty"`
	DefaultLocale string                 `json:"default_locale,omitempty"`
	ShowLogo      bool                   `json:"show_logo"`
	Alignment     EmailTemplateAlignment `json:"alignment,omitempty"`
	// Non-English BCP-47 language variants; English remains the base content.
	Translations map[string]EmailLocalizedContent `json:"translations,omitempty"`
}

// EmailTemplateAlignment is a string enumeration used by EmailTemplateOptions.Alignment.
type EmailTemplateAlignment string

// Values of EmailTemplateAlignment. The server may add values; handle an unknown one gracefully.
const (
	EmailTemplateAlignmentLeft   EmailTemplateAlignment = "left"
	EmailTemplateAlignmentCenter EmailTemplateAlignment = "center"
)

// EmailLocalizedContent is a value of the EmailTemplateOptions.Translations map. Subject and preheader
// limits count Unicode characters. Text and HTML limits count UTF-8 bytes. Action URL variables must
// remain in both text and HTML when provided. Variables in URLs are allowed only as standalone
// action-link href values, never in image sources or within external URLs.
type EmailLocalizedContent struct {
	Subject   string `json:"subject,omitempty"`
	Body      string `json:"body,omitempty"`
	Preheader string `json:"preheader,omitempty"`
	BodyHTML  string `json:"body_html,omitempty"`
}

// EmailTemplateDeliveryScope is a string enumeration used by EmailTemplate.DeliveryScope.
type EmailTemplateDeliveryScope string

// Values of EmailTemplateDeliveryScope. The server may add values; handle an unknown one gracefully.
const (
	EmailTemplateDeliveryScopeTenant   EmailTemplateDeliveryScope = "tenant"
	EmailTemplateDeliveryScopePlatform EmailTemplateDeliveryScope = "platform"
)

// EmailTemplatesCapabilities is the type of EmailTemplates.Capabilities.
type EmailTemplatesCapabilities struct {
	RichText     bool `json:"rich_text"`
	Localization bool `json:"localization"`
	TestEmail    bool `json:"test_email"`
}

// UpdateEmailTemplateParams is the request body of EmailTemplatesService.Update. Subject and preheader
// limits count Unicode characters. Text and HTML limits count UTF-8 bytes. Action URL variables must
// remain in both text and HTML when provided. Variables in URLs are allowed only as standalone
// action-link href values, never in image sources or within external URLs.
type UpdateEmailTemplateParams struct {
	Subject string                `json:"subject,omitempty"`
	Body    string                `json:"body,omitempty"`
	Options *EmailTemplateOptions `json:"options,omitempty"`
}

// PreviewEmailTemplateParams is the request body of EmailTemplatesService.Preview.
type PreviewEmailTemplateParams struct {
	// Each value is limited to 4096 UTF-8 bytes.
	Vars   map[string]string          `json:"vars,omitempty"`
	Locale string                     `json:"locale,omitempty"`
	Draft  *UpdateEmailTemplateParams `json:"draft,omitempty"`
}

// EmailPreview is the result of EmailTemplatesService.Preview.
type EmailPreview struct {
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
	HTML    string `json:"html,omitempty"`
	Locale  string `json:"locale,omitempty"`
}

// SendTestEmailParams is the request body of EmailTemplatesService.SendTest.
type SendTestEmailParams struct {
	Draft  UpdateEmailTemplateParams `json:"draft"`
	Locale string                    `json:"locale,omitempty"`
}

// TestEmailResult is the result of EmailTemplatesService.SendTest.
type TestEmailResult struct {
	Status TestEmailStatus `json:"status,omitempty"`
}

// TestEmailStatus is a string enumeration used by TestEmailResult.Status.
type TestEmailStatus string

// Values of TestEmailStatus. The server may add values; handle an unknown one gracefully.
const (
	TestEmailStatusAccepted TestEmailStatus = "accepted"
)
