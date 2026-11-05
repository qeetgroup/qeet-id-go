package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// SCIMService manages a tenant's SCIM provisioning tokens and lists what the identity provider
// provisioned.
type SCIMService struct{ c *Client }

// Config returns the tenant's SCIM provisioning state.
//
// API: GET /v1/tenants/{tenantID}/scim
func (s *SCIMService) Config(ctx context.Context, tenantID string, opts ...RequestOption) (*SCIMConfig, error) {
	return call[SCIMConfig](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/scim",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListTokens lists the tenant's SCIM tokens.
//
// API: GET /v1/tenants/{tenantID}/scim/tokens
func (s *SCIMService) ListTokens(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[SCIMToken], error) {
	return call[Page[SCIMToken]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/scim/tokens",
		pathArgs: []string{tenantID},
	}, opts)
}

// CreateToken creates a SCIM token. The result carries the token, shown only once.
//
// API: POST /v1/tenants/{tenantID}/scim/tokens
func (s *SCIMService) CreateToken(ctx context.Context, tenantID string, params CreateSCIMTokenParams, opts ...RequestOption) (*CreatedSCIMToken, error) {
	return call[CreatedSCIMToken](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/scim/tokens",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// RotateToken replaces a SCIM token. The old token keeps working for the overlap given.
//
// API: POST /v1/tenants/{tenantID}/scim/tokens/{id}/rotate
func (s *SCIMService) RotateToken(ctx context.Context, tenantID string, tokenID string, params RotateSCIMTokenParams, opts ...RequestOption) (*RotatedSCIMToken, error) {
	return call[RotatedSCIMToken](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/scim/tokens/{tokenID}/rotate",
		pathArgs: []string{tenantID, tokenID},
		body:     params,
	}, opts)
}

// RevokeToken revokes a SCIM token.
//
// API: DELETE /v1/tenants/{tenantID}/scim/tokens/{id}
func (s *SCIMService) RevokeToken(ctx context.Context, tenantID string, tokenID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/scim/tokens/{tokenID}",
		pathArgs: []string{tenantID, tokenID},
	}, nil, opts)
}

// ReplaceAllTokens issues a new SCIM token and revokes every earlier one at once. The result carries
// the token, shown only once.
//
// API: POST /v1/tenants/{tenantID}/scim/token
func (s *SCIMService) ReplaceAllTokens(ctx context.Context, tenantID string, opts ...RequestOption) (*ReplacedSCIMTokens, error) {
	return call[ReplacedSCIMTokens](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/scim/token",
		pathArgs: []string{tenantID},
	}, opts)
}

// RevokeAllTokens revokes every SCIM token of the tenant, which stops provisioning.
//
// API: DELETE /v1/tenants/{tenantID}/scim/token
func (s *SCIMService) RevokeAllTokens(ctx context.Context, tenantID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/scim/token",
		pathArgs: []string{tenantID},
	}, nil, opts)
}

// ListUsers lists the users the identity provider provisioned, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/scim/users
func (s *SCIMService) ListUsers(ctx context.Context, tenantID string, params *ListSCIMResourcesParams, opts ...RequestOption) (*Page[SCIMUser], error) {
	return call[Page[SCIMUser]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/scim/users",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AllUsers returns an iterator over every item ListUsers returns, across all pages. It fetches a page
// only when the loop needs it, and breaking out of the loop stops paging. A failed request ends the
// iteration with a non-nil error.
func (s *SCIMService) AllUsers(ctx context.Context, tenantID string, params *ListSCIMResourcesParams, opts ...RequestOption) iter.Seq2[*SCIMUser, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[SCIMUser], error) {
		return s.ListUsers(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// ListGroups lists the groups the identity provider provisioned, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/scim/groups
func (s *SCIMService) ListGroups(ctx context.Context, tenantID string, params *ListSCIMResourcesParams, opts ...RequestOption) (*Page[SCIMGroup], error) {
	return call[Page[SCIMGroup]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/scim/groups",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AllGroups returns an iterator over every item ListGroups returns, across all pages. It fetches a
// page only when the loop needs it, and breaking out of the loop stops paging. A failed request ends
// the iteration with a non-nil error.
func (s *SCIMService) AllGroups(ctx context.Context, tenantID string, params *ListSCIMResourcesParams, opts ...RequestOption) iter.Seq2[*SCIMGroup, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[SCIMGroup], error) {
		return s.ListGroups(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// SCIMConfig is the state of a tenant's SCIM provisioning.
type SCIMConfig struct {
	// When the newest working token was issued.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// The latest use of any working token, recorded at most once a minute.
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`
	ProvisionedCount int        `json:"provisioned_count"`
	// How many tokens work, rotated ones still in their overlap included.
	TokenCount int `json:"token_count"`
	// The newest working token's prefix.
	TokenPrefix string `json:"token_prefix,omitempty"`
	// Whether any token works.
	TokenSet bool `json:"token_set"`
	// How the last change the identity provider sent was answered; null before the first. Reads are not
	// recorded.
	LastRequest *SCIMConfigLastRequest `json:"last_request,omitempty"`
}

// SCIMConfigLastRequest is the type of SCIMConfig.LastRequest. How the last change the identity
// provider sent was answered; null before the first. Reads are not recorded.
type SCIMConfigLastRequest struct {
	At         time.Time                      `json:"at,omitzero"`
	Operation  SCIMConfigLastRequestOperation `json:"operation,omitempty"`
	Resource   SCIMConfigLastRequestResource  `json:"resource,omitempty"`
	HTTPStatus int                            `json:"http_status"`
	Status     SCIMConfigLastRequestStatus    `json:"status,omitempty"`
	// The directory error code, such as DIR-409.
	ErrorCode string `json:"error_code,omitempty"`
	// The SCIM error type the answer gave, such as scim.uniqueness or scim.invalidValue.
	Reason string `json:"reason,omitempty"`
}

// SCIMConfigLastRequestOperation is a string enumeration used by SCIMConfigLastRequest.Operation.
type SCIMConfigLastRequestOperation string

// Values of SCIMConfigLastRequestOperation. The server may add values; handle an unknown one gracefully.
const (
	SCIMConfigLastRequestOperationPost   SCIMConfigLastRequestOperation = "post"
	SCIMConfigLastRequestOperationPut    SCIMConfigLastRequestOperation = "put"
	SCIMConfigLastRequestOperationPatch  SCIMConfigLastRequestOperation = "patch"
	SCIMConfigLastRequestOperationDelete SCIMConfigLastRequestOperation = "delete"
)

// SCIMConfigLastRequestResource is a string enumeration used by SCIMConfigLastRequest.Resource.
type SCIMConfigLastRequestResource string

// Values of SCIMConfigLastRequestResource. The server may add values; handle an unknown one gracefully.
const (
	SCIMConfigLastRequestResourceUser  SCIMConfigLastRequestResource = "user"
	SCIMConfigLastRequestResourceGroup SCIMConfigLastRequestResource = "group"
)

// SCIMConfigLastRequestStatus is a string enumeration used by SCIMConfigLastRequest.Status.
type SCIMConfigLastRequestStatus string

// Values of SCIMConfigLastRequestStatus. The server may add values; handle an unknown one gracefully.
const (
	SCIMConfigLastRequestStatusSuccess SCIMConfigLastRequestStatus = "success"
	SCIMConfigLastRequestStatusWarning SCIMConfigLastRequestStatus = "warning"
	SCIMConfigLastRequestStatusFailed  SCIMConfigLastRequestStatus = "failed"
)

// SCIMToken is a bearer token an identity provider uses to provision over SCIM. The token itself is
// shown only when created.
type SCIMToken struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// The token's first characters, to recognise it by.
	TokenPrefix string `json:"token_prefix,omitempty"`
	// active: works. rotated: replaced, and works until expires_at, when its overlap ends. expired: past
	// expires_at. revoked: revoked.
	Status    SCIMTokenStatus `json:"status,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitzero"`
	ExpiresAt *time.Time      `json:"expires_at,omitempty"`
	// Recorded at most once a minute.
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	// The token that replaced this one when it was rotated.
	ReplacedBy string `json:"replaced_by,omitempty"`
}

// SCIMTokenStatus is a string enumeration used by SCIMToken.Status. active: works. rotated: replaced,
// and works until expires_at, when its overlap ends. expired: past expires_at. revoked: revoked.
type SCIMTokenStatus string

// Values of SCIMTokenStatus. The server may add values; handle an unknown one gracefully.
const (
	SCIMTokenStatusActive  SCIMTokenStatus = "active"
	SCIMTokenStatusRotated SCIMTokenStatus = "rotated"
	SCIMTokenStatusExpired SCIMTokenStatus = "expired"
	SCIMTokenStatusRevoked SCIMTokenStatus = "revoked"
)

// CreateSCIMTokenParams is the request body of SCIMService.CreateToken.
type CreateSCIMTokenParams struct {
	Name string `json:"name,omitempty"`
	// When the token stops working; in the future. Omit for none.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CreatedSCIMToken is the result of SCIMService.CreateToken.
type CreatedSCIMToken struct {
	SCIMToken SCIMToken `json:"scim_token,omitzero"`
	// The token to give the identity provider, shown only here.
	Token string `json:"token,omitempty"`
}

// RotateSCIMTokenParams is the request body of SCIMService.RotateToken.
type RotateSCIMTokenParams struct {
	// How long the old token keeps working, in hours. Default 168.
	OverlapHours *int `json:"overlap_hours,omitempty"`
	// When the new token stops working; in the future. Omit for none.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// RotatedSCIMToken is the result of SCIMService.RotateToken.
type RotatedSCIMToken struct {
	Replaced  SCIMToken `json:"replaced,omitzero"`
	SCIMToken SCIMToken `json:"scim_token,omitzero"`
	// The new token to give the identity provider, shown only here.
	Token string `json:"token,omitempty"`
}

// ReplacedSCIMTokens is the result of SCIMService.ReplaceAllTokens.
type ReplacedSCIMTokens struct {
	Config SCIMConfig `json:"config,omitzero"`
	Token  string     `json:"token,omitempty"`
}

// ListSCIMResourcesParams holds the optional query parameters of SCIMService.ListUsers. A nil
// *ListSCIMResourcesParams sends none.
type ListSCIMResourcesParams struct {
	// Page size, 1 to 200; 50 by default.
	Limit int
	// The next_cursor of the previous page.
	Cursor string
}

func (p *ListSCIMResourcesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ListSCIMResourcesParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListSCIMResourcesParams) withCursor(cursor string) *ListSCIMResourcesParams {
	var c ListSCIMResourcesParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// SCIMUser is a user an identity provider provisioned over SCIM.
type SCIMUser struct {
	CreatedAt   time.Time `json:"created_at,omitzero"`
	DisplayName string    `json:"display_name,omitempty"`
	Email       string    `json:"email,omitempty"`
	ExternalID  string    `json:"external_id,omitempty"`
	ID          string    `json:"id,omitempty"`
	Status      string    `json:"status,omitempty"`
}

// SCIMGroup is a group an identity provider provisioned over SCIM.
type SCIMGroup struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Members    int       `json:"members"`
	CreatedAt  time.Time `json:"created_at,omitzero"`
	UpdatedAt  time.Time `json:"updated_at,omitzero"`
}
