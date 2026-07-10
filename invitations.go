package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Invitation is a pending (or resolved) org invitation. There is no
// per-invitation Get or Resend in the backend — only Create, Delete, and
// List. RoleID is the role the invitee will hold once accepted.
type Invitation struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	Email      string `json:"email"`
	RoleID     string `json:"role_id,omitempty"`
	Status     string `json:"status"`
	ExpiresAt  string `json:"expires_at"`
	AcceptedAt string `json:"accepted_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type CreateInvitationInput struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id,omitempty"`
}

// CreateInvitationResult carries the raw invite token alongside the
// invitation record — admins frequently need to copy the link directly when
// email delivery isn't trustworthy yet.
type CreateInvitationResult struct {
	Invite Invitation `json:"invite"`
	Token  string     `json:"token"`
}

// InvitationsService manages org invitations. Accepting an invitation is
// deliberately not wrapped here — like login/signup, it's an end-user auth
// action (it returns a fresh session token pair for the invitee), not an
// admin management operation.
type InvitationsService struct{ t *transport.Transport }

func (s *InvitationsService) Create(ctx context.Context, in CreateInvitationInput) (*CreateInvitationResult, error) {
	var out CreateInvitationResult
	err := s.t.Post(ctx, "/v1/invites", in, &out)
	return &out, err
}

func (s *InvitationsService) Delete(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/invites/"+url.PathEscape(id), nil)
}

func (s *InvitationsService) List(ctx context.Context, tenantID string) ([]Invitation, error) {
	var env marshal.Envelope[Invitation]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/invites", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}
