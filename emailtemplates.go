package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// EmailTemplate is a resolved template — the tenant's custom override if one
// exists (Custom is true), otherwise the platform default. There's a fixed
// catalog of Keys (Name/Description/Variables are metadata for building an
// editor UI); templates aren't independently created, only customized via
// Update and reverted via Reset.
type EmailTemplate struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Variables   []string `json:"variables,omitempty"`
	Custom      bool     `json:"custom"`
}

type UpdateEmailTemplateInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// PreviewEmailTemplateResult is the rendered output with Vars substituted —
// {{name}} placeholders left intact are unfilled.
type PreviewEmailTemplateResult struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// EmailTemplatesService manages transactional email template overrides.
type EmailTemplatesService struct{ t *transport.Transport }

func (s *EmailTemplatesService) List(ctx context.Context, tenantID string) ([]EmailTemplate, error) {
	var env marshal.Envelope[EmailTemplate]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/email-templates", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *EmailTemplatesService) Get(ctx context.Context, tenantID, key string) (*EmailTemplate, error) {
	var out EmailTemplate
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/email-templates/"+url.PathEscape(key), transport.RequestOptions{}, &out)
	return &out, err
}

// Update sets a custom override for a template (full replace of subject+body).
func (s *EmailTemplatesService) Update(ctx context.Context, tenantID, key string, in UpdateEmailTemplateInput) (*EmailTemplate, error) {
	var out EmailTemplate
	err := s.t.Put(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/email-templates/"+url.PathEscape(key), in, &out)
	return &out, err
}

// Reset discards a custom override, reverting to the platform default.
func (s *EmailTemplatesService) Reset(ctx context.Context, tenantID, key string) (*EmailTemplate, error) {
	var out EmailTemplate
	err := s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/email-templates/"+url.PathEscape(key), &out)
	return &out, err
}

// Preview renders the template's current subject/body with vars substituted
// for {{placeholder}} tokens — nothing is sent.
func (s *EmailTemplatesService) Preview(ctx context.Context, tenantID, key string, vars map[string]string) (*PreviewEmailTemplateResult, error) {
	var out PreviewEmailTemplateResult
	body := map[string]any{"vars": vars}
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/email-templates/"+url.PathEscape(key)+"/preview", body, &out)
	return &out, err
}
