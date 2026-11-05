package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// OAuthGrantsService lists and revokes the grants people gave applications.
type OAuthGrantsService struct{ c *Client }

// List lists active grants, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/oauth/grants
func (s *OAuthGrantsService) List(ctx context.Context, tenantID string, params *ListOAuthGrantsParams, opts ...RequestOption) (*Page[OAuthGrant], error) {
	return call[Page[OAuthGrant]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/oauth/grants",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// All returns an iterator over every item List returns, across all pages. It fetches a page only when
// the loop needs it, and breaking out of the loop stops paging. A failed request ends the iteration
// with a non-nil error.
func (s *OAuthGrantsService) All(ctx context.Context, tenantID string, params *ListOAuthGrantsParams, opts ...RequestOption) iter.Seq2[*OAuthGrant, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[OAuthGrant], error) {
		return s.List(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// Revoke revokes a grant and the tokens issued under it.
//
// API: DELETE /v1/tenants/{tenantID}/oauth/grants/{id}
func (s *OAuthGrantsService) Revoke(ctx context.Context, tenantID, grantID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/oauth/grants/{grantID}",
		pathArgs: []string{tenantID, grantID},
	}, nil, opts)
}

// RevokeMany revokes several grants, or every grant of one application.
//
// API: POST /v1/tenants/{tenantID}/oauth/grants/revoke
func (s *OAuthGrantsService) RevokeMany(ctx context.Context, tenantID string, params RevokeOAuthGrantsParams, opts ...RequestOption) (*RevokedGrants, error) {
	return call[RevokedGrants](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/oauth/grants/revoke",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// ListOAuthGrantsParams holds the optional query parameters of OAuthGrantsService.List. A nil
// *ListOAuthGrantsParams sends none.
type ListOAuthGrantsParams struct {
	// Only the grants given to this application.
	ClientID string
	// Page size, 1 to 500; 100 by default.
	Limit int
	// The next_cursor of the previous page.
	Cursor string
}

func (p *ListOAuthGrantsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.ClientID != "" {
		q.Set("client_id", p.ClientID)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ListOAuthGrantsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListOAuthGrantsParams) withCursor(cursor string) *ListOAuthGrantsParams {
	var c ListOAuthGrantsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// OAuthGrant is a person's consent for an application, with the refresh tokens issued under it.
type OAuthGrant struct {
	// When the grant began.
	AuthorizedAt time.Time `json:"authorized_at,omitzero"`
	ClientName   string    `json:"client_name,omitempty"`
	// When the person last consented to the application; null for an application that asks for no
	// consent.
	ConsentedAt *time.Time `json:"consented_at,omitempty"`
	// When the application was last issued tokens under the grant (the same as issued_at).
	LastUsedAt time.Time `json:"last_used_at,omitzero"`
	UserName   string    `json:"user_name,omitempty"`
	ClientID   string    `json:"client_id,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitzero"`
	ID         string    `json:"id,omitempty"`
	IssuedAt   time.Time `json:"issued_at,omitzero"`
	Scopes     []string  `json:"scopes,omitempty"`
	UserEmail  string    `json:"user_email,omitempty"`
	UserID     string    `json:"user_id,omitempty"`
}

// RevokeOAuthGrantsParams is the request body of OAuthGrantsService.RevokeMany.
type RevokeOAuthGrantsParams struct {
	IDs      []string `json:"ids,omitempty"`
	ClientID string   `json:"client_id,omitempty"`
}

// RevokedGrants is the result of OAuthGrantsService.RevokeMany.
type RevokedGrants struct {
	Revoked int `json:"revoked"`
}
