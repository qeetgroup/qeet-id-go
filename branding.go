package qeetid

import (
	"context"
	"net/http"
)

// BrandingService manages an organization's hosted sign-in branding.
type BrandingService struct{ c *Client }

// Get returns an organization's branding.
//
// API: GET /v1/organizations/{organizationID}/branding
func (s *BrandingService) Get(ctx context.Context, organizationID string, opts ...RequestOption) (*Branding, error) {
	return call[Branding](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/branding",
		pathArgs: []string{organizationID},
	}, opts)
}

// Update replaces an organization's branding. Read it with Get, change it, and send it back.
//
// API: PUT /v1/organizations/{organizationID}/branding
func (s *BrandingService) Update(ctx context.Context, organizationID string, params Branding, opts ...RequestOption) (*Branding, error) {
	return call[Branding](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/organizations/{organizationID}/branding",
		pathArgs: []string{organizationID},
		body:     params,
	}, opts)
}

// Branding is the look of an organization's hosted sign-in pages and emails.
type Branding struct {
	OrganizationID string `json:"organization_id,omitempty"`
	LogoURL        string `json:"logo_url,omitempty"`
	PrimaryColor   string `json:"primary_color,omitempty"`
	// Brand-panel accent, independent of page background.
	SecondaryColor string `json:"secondary_color,omitempty"`
	CustomDomain   string `json:"custom_domain,omitempty"`
	EmailFromName  string `json:"email_from_name,omitempty"`
	// Accepted and ignored, and never returned. Email is sent from the platform's own address (an
	// organization's domain would need verifying with the email provider); replies go to
	// settings.email_reply_to.
	EmailFromAddress string           `json:"email_from_address,omitempty"`
	Settings         BrandingSettings `json:"settings,omitzero"`
}

// BrandingSettings is the type of Branding.Settings.
type BrandingSettings struct {
	// Hosted-login page background; empty keeps the default theme.
	BackgroundColor string `json:"background_color,omitempty"`
	// Single reply mailbox for organization-scoped messages; empty omits Reply-To.
	EmailReplyTo string `json:"email_reply_to,omitempty"`
}
