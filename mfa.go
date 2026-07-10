package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// MFAService is the admin-side MFA surface. The backend has no admin
// endpoint to list a user's enrolled factors or force MFA on — only an
// admin-initiated reset (clearing all factors, e.g. after a lost device).
// This mirrors Users.ResetMFA; both reach the same endpoint.
type MFAService struct{ t *transport.Transport }

func (s *MFAService) Reset(ctx context.Context, userID string) error {
	return s.t.Delete(ctx, "/v1/users/"+url.PathEscape(userID)+"/mfa", nil)
}
