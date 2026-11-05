package qeetid

import (
	"context"
	"net/http"
	"time"
)

// LogSinksService manages the streams that forward a tenant's audit events to a SIEM.
type LogSinksService struct{ c *Client }

// List lists a tenant's log sinks.
//
// API: GET /v1/tenants/{tenantID}/log-sinks
func (s *LogSinksService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[LogSink], error) {
	return call[Page[LogSink]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/log-sinks",
		pathArgs: []string{tenantID},
	}, opts)
}

// Create creates a log sink.
//
// API: POST /v1/tenants/{tenantID}/log-sinks
func (s *LogSinksService) Create(ctx context.Context, tenantID string, params LogSinkParams, opts ...RequestOption) (*LogSink, error) {
	return call[LogSink](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/log-sinks",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// SetEnabled enables or disables a log sink.
//
// API: PATCH /v1/tenants/{tenantID}/log-sinks/{id}
func (s *LogSinksService) SetEnabled(ctx context.Context, tenantID, sinkID string, params SetLogSinkEnabledParams, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/log-sinks/{sinkID}",
		pathArgs: []string{tenantID, sinkID},
		body:     params,
	}, nil, opts)
}

// Delete deletes a log sink.
//
// API: DELETE /v1/tenants/{tenantID}/log-sinks/{id}
func (s *LogSinksService) Delete(ctx context.Context, tenantID, sinkID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/log-sinks/{sinkID}",
		pathArgs: []string{tenantID, sinkID},
	}, nil, opts)
}

// Test sends a test event to a saved log sink.
//
// API: POST /v1/tenants/{tenantID}/log-sinks/{id}/test
func (s *LogSinksService) Test(ctx context.Context, tenantID, sinkID string, opts ...RequestOption) (*LogSinkTest, error) {
	return call[LogSinkTest](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/log-sinks/{sinkID}/test",
		pathArgs: []string{tenantID, sinkID},
	}, opts)
}

// TestSettings sends a test event to a log sink before it is saved.
//
// API: POST /v1/tenants/{tenantID}/log-sinks/test
func (s *LogSinksService) TestSettings(ctx context.Context, tenantID string, params LogSinkParams, opts ...RequestOption) (*LogSinkTest, error) {
	return call[LogSinkTest](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/log-sinks/test",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// LogSink forwards a tenant's audit events to a SIEM. Its token is never returned.
type LogSink struct {
	CreatedAt       time.Time `json:"created_at,omitzero"`
	Enabled         bool      `json:"enabled"`
	Endpoint        string    `json:"endpoint,omitempty"`
	ID              string    `json:"id,omitempty"`
	LastError       string    `json:"last_error,omitempty"`
	LastForwardedAt time.Time `json:"last_forwarded_at,omitzero"`
	Type            string    `json:"type,omitempty"`
}

// LogSinkParams is the request body of LogSinksService.Create.
type LogSinkParams struct {
	Type     LogSinkType `json:"type,omitempty"`
	Endpoint string      `json:"endpoint,omitempty"`
	Token    string      `json:"token,omitempty"`
}

// LogSinkType is a string enumeration used by LogSinkParams.Type.
type LogSinkType string

// Values of LogSinkType. The server may add values; handle an unknown one gracefully.
const (
	LogSinkTypeSplunkHEC LogSinkType = "splunk_hec"
	LogSinkTypeDatadog   LogSinkType = "datadog"
	LogSinkTypeHTTP      LogSinkType = "http"
)

// SetLogSinkEnabledParams is the request body of LogSinksService.SetEnabled.
type SetLogSinkEnabledParams struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// LogSinkTest is the result of LogSinksService.Test.
type LogSinkTest struct {
	// Whether the destination accepted the event (a 2xx answer).
	Delivered bool `json:"delivered"`
	// The destination's HTTP status; absent when it gave no answer.
	Status int `json:"status"`
	// Why it wasn't delivered.
	Error      string `json:"error,omitempty"`
	DurationMs int    `json:"duration_ms"`
}
