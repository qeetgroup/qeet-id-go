package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// LogSink is a configured SIEM forwarding destination. The auth token is
// write-only — never returned on any response.
type LogSink struct {
	ID              string `json:"id"`
	Type            string `json:"type"` // splunk_hec | datadog | http
	Endpoint        string `json:"endpoint"`
	Enabled         bool   `json:"enabled"`
	LastForwardedAt string `json:"last_forwarded_at,omitempty"`
	LastError       string `json:"last_error,omitempty"`
	CreatedAt       string `json:"created_at"`
}

type CreateLogSinkInput struct {
	Type     string `json:"type"` // splunk_hec | datadog | http
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}

// LogSinksService manages SIEM log-forwarding destinations. The OpenAPI tag
// is "Log Sinks" and every path is /log-sinks — there is no /siem path
// despite the backend Go package being named "siem".
type LogSinksService struct{ t *transport.Transport }

func (s *LogSinksService) Create(ctx context.Context, tenantID string, in CreateLogSinkInput) (*LogSink, error) {
	var out LogSink
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/log-sinks", in, &out)
	return &out, err
}

func (s *LogSinksService) List(ctx context.Context, tenantID string) ([]LogSink, error) {
	var env marshal.Envelope[LogSink]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/log-sinks", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *LogSinksService) SetEnabled(ctx context.Context, tenantID, id string, enabled bool) (*LogSink, error) {
	var out LogSink
	body := map[string]bool{"enabled": enabled}
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/log-sinks/"+url.PathEscape(id), body, &out)
	return &out, err
}

func (s *LogSinksService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/log-sinks/"+url.PathEscape(id), nil)
}
