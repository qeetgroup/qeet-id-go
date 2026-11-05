package qeetid

import (
	"context"
	"net/http"
	"time"
)

// DeviceAuthorizationsService lists and revokes device-flow (RFC 8628) authorizations.
type DeviceAuthorizationsService struct{ c *Client }

// List lists pending and approved device authorizations. Page.Truncated reports a list cut short by
// the server's cap.
//
// API: GET /v1/tenants/{tenantID}/oauth/devices
func (s *DeviceAuthorizationsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[DeviceAuthorization], error) {
	return call[Page[DeviceAuthorization]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oauth/devices",
		pathArgs: []string{tenantID},
	}, opts)
}

// Revoke revokes a device authorization.
//
// API: DELETE /v1/tenants/{tenantID}/oauth/devices/{id}
func (s *DeviceAuthorizationsService) Revoke(ctx context.Context, tenantID string, authorizationID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/oauth/devices/{authorizationID}",
		pathArgs: []string{tenantID, authorizationID},
	}, nil, opts)
}

// DeviceAuthorization is a device-flow (RFC 8628) authorization. It never carries the device code.
type DeviceAuthorization struct {
	ID       string `json:"id,omitempty"`
	ClientID string `json:"client_id,omitempty"`
	// The application's name; empty once the application is deleted.
	ClientName   string                    `json:"client_name,omitempty"`
	UserCode     string                    `json:"user_code,omitempty"`
	Status       DeviceAuthorizationStatus `json:"status,omitempty"`
	UserID       string                    `json:"user_id,omitempty"`
	UserEmail    string                    `json:"user_email,omitempty"`
	Scopes       []string                  `json:"scopes,omitempty"`
	CreatedAt    time.Time                 `json:"created_at,omitzero"`
	ExpiresAt    time.Time                 `json:"expires_at,omitzero"`
	LastPolledAt *time.Time                `json:"last_polled_at,omitempty"`
	// When the device redeemed the approval for tokens; null until then.
	ConsumedAt *time.Time `json:"consumed_at,omitempty"`
}

// DeviceAuthorizationStatus is a string enumeration used by DeviceAuthorization.Status.
type DeviceAuthorizationStatus string

// Values of DeviceAuthorizationStatus. The server may add values; handle an unknown one gracefully.
const (
	DeviceAuthorizationStatusPending    DeviceAuthorizationStatus = "pending"
	DeviceAuthorizationStatusAuthorized DeviceAuthorizationStatus = "authorized"
	DeviceAuthorizationStatusDenied     DeviceAuthorizationStatus = "denied"
)
