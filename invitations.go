package qeetid

import (
	"context"
	"net/http"
	"time"
)

// InvitationsService invites people to join a tenant.
type InvitationsService struct{ c *Client }

// Create invites a person by email. The accept link goes only to the invited address; the result
// reports whether the email went out.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/invites
func (s *InvitationsService) Create(ctx context.Context, params CreateInvitationParams, opts ...RequestOption) (*InvitationResult, error) {
	return call[InvitationResult](ctx, s.c, &request{method: http.MethodPost, path: "/v1/invites", body: params, idempotent: true}, opts)
}

// List lists a tenant's invitations.
//
// API: GET /v1/tenants/{tenantID}/invites
func (s *InvitationsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Invitation], error) {
	return call[Page[Invitation]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/invites",
		pathArgs: []string{tenantID},
	}, opts)
}

// Revoke revokes a pending invitation.
//
// API: DELETE /v1/invites/{id}
func (s *InvitationsService) Revoke(ctx context.Context, invitationID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/invites/{invitationID}",
		pathArgs: []string{invitationID},
	}, nil, opts)
}

// Resend sends a pending invitation's email again.
//
// API: POST /v1/invites/{id}/resend
func (s *InvitationsService) Resend(ctx context.Context, invitationID string, opts ...RequestOption) (*InvitationResult, error) {
	return call[InvitationResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/invites/{invitationID}/resend",
		pathArgs: []string{invitationID},
	}, opts)
}

// CreateInvitationParams is the request body of InvitationsService.Create.
type CreateInvitationParams struct {
	Email    string  `json:"email,omitempty"`
	RoleID   *string `json:"role_id,omitempty"`
	TenantID string  `json:"tenant_id,omitempty"`
}

// InvitationResult carries the invitation an Invitations call created or sent again, with whether its
// email went out.
type InvitationResult struct {
	Invite Invitation `json:"invite,omitzero"`
}

// Invitation is an invitation for a person to join a tenant.
type Invitation struct {
	AcceptedAt *time.Time         `json:"accepted_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitzero"`
	Delivery   InvitationDelivery `json:"delivery,omitzero"`
	Email      string             `json:"email,omitempty"`
	ExpiresAt  time.Time          `json:"expires_at,omitzero"`
	ID         string             `json:"id,omitempty"`
	RoleID     string             `json:"role_id,omitempty"`
	Status     string             `json:"status,omitempty"`
	TenantID   string             `json:"tenant_id,omitempty"`
}

// InvitationDelivery reports whether an invitation's email was sent. The invitation exists either way.
type InvitationDelivery struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}
