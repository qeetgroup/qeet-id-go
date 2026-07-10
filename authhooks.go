package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// AuthHook is an HMAC-signed custom login-flow hook. Secret is write-only —
// never returned. A tenant may have multiple hooks; this is a genuine
// collection (List/Create/Update-by-id/Delete-by-id), not a singleton.
type AuthHook struct {
	ID        string `json:"id"`
	Trigger   string `json:"trigger"`
	URL       string `json:"url"`
	Enabled   bool   `json:"enabled"`
	FailOpen  bool   `json:"fail_open"`
	CreatedAt string `json:"created_at"`
}

// CreateAuthHookInput — Secret signs the HMAC payload sent to URL on every
// invocation; store it to verify inbound calls. FailOpen defaults to true
// server-side if omitted (a hook outage doesn't block login).
type CreateAuthHookInput struct {
	URL      string `json:"url"`
	Secret   string `json:"secret"`
	FailOpen *bool  `json:"fail_open,omitempty"`
}

// UpdateAuthHookInput — both fields are always sent (a full replace, not a
// partial patch, despite the PATCH verb): the backend has no notion of
// "leave unset fields alone" here.
type UpdateAuthHookInput struct {
	Enabled  bool `json:"enabled"`
	FailOpen bool `json:"fail_open"`
}

// AuthHooksService manages HMAC-signed pre/post-login webhooks.
type AuthHooksService struct{ t *transport.Transport }

func (s *AuthHooksService) List(ctx context.Context, tenantID string) ([]AuthHook, error) {
	var env marshal.Envelope[AuthHook]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/auth-hooks", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *AuthHooksService) Create(ctx context.Context, tenantID string, in CreateAuthHookInput) (*AuthHook, error) {
	var out AuthHook
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/auth-hooks", in, &out)
	return &out, err
}

func (s *AuthHooksService) Update(ctx context.Context, tenantID, id string, in UpdateAuthHookInput) error {
	path := "/v1/tenants/" + url.PathEscape(tenantID) + "/auth-hooks/" + url.PathEscape(id)
	return s.t.Patch(ctx, path, in, nil)
}

func (s *AuthHooksService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/auth-hooks/"+url.PathEscape(id), nil)
}
