package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// UsersService manages the people in a tenant, their sessions, devices and verification state.
type UsersService struct{ c *Client }

// Create creates a user in a tenant.
//
// It accepts WithIdempotencyKey. Without one, the client sends a random key so that a retry cannot
// apply the request twice.
//
// API: POST /v1/users
func (s *UsersService) Create(ctx context.Context, params CreateUserParams, opts ...RequestOption) (*User, error) {
	return call[User](ctx, s.c, &request{method: http.MethodPost, path: "/v1/users", body: params, idempotent: true}, opts)
}

// Get returns a user.
//
// API: GET /v1/users/{id}
func (s *UsersService) Get(ctx context.Context, userID string, opts ...RequestOption) (*User, error) {
	return call[User](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users/{userID}", pathArgs: []string{userID}}, opts)
}

// Update changes a user's profile or status. Fields left nil are not changed.
//
// API: PATCH /v1/users/{id}
func (s *UsersService) Update(ctx context.Context, userID string, params UpdateUserParams, opts ...RequestOption) (*User, error) {
	return call[User](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/users/{userID}",
		pathArgs: []string{userID},
		body:     params,
	}, opts)
}

// Delete soft-deletes a user. The user can be brought back with Restore until it is purged.
//
// API: DELETE /v1/users/{id}
func (s *UsersService) Delete(ctx context.Context, userID string, opts ...RequestOption) (*UserDeletion, error) {
	return call[UserDeletion](ctx, s.c, &request{method: http.MethodDelete, path: "/v1/users/{userID}", pathArgs: []string{userID}}, opts)
}

// List lists the users of the caller's tenant, one page at a time.
//
// API: GET /v1/users
func (s *UsersService) List(ctx context.Context, params *ListUsersParams, opts ...RequestOption) (*Page[User], error) {
	return call[Page[User]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users", query: params.values()}, opts)
}

// All returns an iterator over every item List returns, across all pages. It fetches a page only when
// the loop needs it, and breaking out of the loop stops paging. A failed request ends the iteration
// with a non-nil error.
func (s *UsersService) All(ctx context.Context, params *ListUsersParams, opts ...RequestOption) iter.Seq2[*User, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[User], error) {
		return s.List(ctx, params.withCursor(cursor), opts...)
	})
}

// SetPassword sets a user's password. It ends the user's sessions.
//
// API: POST /v1/users/{id}/password
func (s *UsersService) SetPassword(ctx context.Context, userID string, params SetPasswordParams, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/password",
		pathArgs: []string{userID},
		body:     params,
	}, nil, opts)
}

// BulkCreate creates up to 1,000 users in one request and reports the outcome of each row.
//
// API: POST /v1/users/bulk
func (s *UsersService) BulkCreate(ctx context.Context, params BulkCreateUsersParams, opts ...RequestOption) (*BulkImportResult, error) {
	return call[BulkImportResult](ctx, s.c, &request{method: http.MethodPost, path: "/v1/users/bulk", body: params}, opts)
}

// ListDeleted lists soft-deleted users (the recycle bin).
//
// API: GET /v1/users/deleted
func (s *UsersService) ListDeleted(ctx context.Context, params *ListDeletedUsersParams, opts ...RequestOption) (*Page[DeletedUser], error) {
	return call[Page[DeletedUser]](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users/deleted", query: params.values()}, opts)
}

// Restore undoes a soft delete.
//
// API: POST /v1/users/{id}/restore
func (s *UsersService) Restore(ctx context.Context, userID string, opts ...RequestOption) (*User, error) {
	return call[User](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/restore",
		pathArgs: []string{userID},
	}, opts)
}

// Purge permanently deletes a soft-deleted user. It cannot be undone.
//
// API: DELETE /v1/users/{id}/purge
func (s *UsersService) Purge(ctx context.Context, userID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/users/{userID}/purge",
		pathArgs: []string{userID},
	}, nil, opts)
}

// RemoveMembership removes a person from the caller's tenant without deleting their account.
//
// API: DELETE /v1/users/{id}/membership
func (s *UsersService) RemoveMembership(ctx context.Context, userID string, opts ...RequestOption) (*MembershipRemoval, error) {
	return call[MembershipRemoval](ctx, s.c, &request{
		method:   http.MethodDelete,
		path:     "/v1/users/{userID}/membership",
		pathArgs: []string{userID},
	}, opts)
}

// ResetMFA removes every second factor of a user who has lost access to them.
//
// API: DELETE /v1/users/{id}/mfa
func (s *UsersService) ResetMFA(ctx context.Context, userID string, opts ...RequestOption) (*Message, error) {
	return call[Message](ctx, s.c, &request{method: http.MethodDelete, path: "/v1/users/{userID}/mfa", pathArgs: []string{userID}}, opts)
}

// SetMFARequired sets whether a user must use a second factor, whatever the tenant's policy says.
//
// API: PUT /v1/users/{id}/mfa-required
func (s *UsersService) SetMFARequired(ctx context.Context, userID string, params SetMFARequiredParams, opts ...RequestOption) (*MFARequirement, error) {
	return call[MFARequirement](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/users/{userID}/mfa-required",
		pathArgs: []string{userID},
		body:     params,
	}, opts)
}

// Security returns a user's security posture: factors, passkeys, password state and sessions.
//
// API: GET /v1/users/{id}/security
func (s *UsersService) Security(ctx context.Context, userID string, opts ...RequestOption) (*UserSecurity, error) {
	return call[UserSecurity](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/users/{userID}/security",
		pathArgs: []string{userID},
	}, opts)
}

// Access summarizes what a user can reach: roles, groups, applications and permission counts.
//
// API: GET /v1/users/{id}/access
func (s *UsersService) Access(ctx context.Context, userID string, opts ...RequestOption) (*UserAccess, error) {
	return call[UserAccess](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users/{userID}/access", pathArgs: []string{userID}}, opts)
}

// Stats returns summary counts for the caller's user directory.
//
// API: GET /v1/users/stats
func (s *UsersService) Stats(ctx context.Context, opts ...RequestOption) (*UserStats, error) {
	return call[UserStats](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users/stats"}, opts)
}

// Trends returns 30-day trends for the caller's user directory.
//
// API: GET /v1/users/trends
func (s *UsersService) Trends(ctx context.Context, opts ...RequestOption) (*UserTrends, error) {
	return call[UserTrends](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users/trends"}, opts)
}

// ListSuspended lists suspended users with review metrics.
//
// API: GET /v1/users/suspended
func (s *UsersService) ListSuspended(ctx context.Context, params *ListSuspendedUsersParams, opts ...RequestOption) (*SuspendedUsers, error) {
	return call[SuspendedUsers](ctx, s.c, &request{method: http.MethodGet, path: "/v1/users/suspended", query: params.values()}, opts)
}

// Reactivate reactivates a suspended user whose home is the caller's tenant.
//
// API: POST /v1/users/{id}/reactivate
func (s *UsersService) Reactivate(ctx context.Context, userID string, opts ...RequestOption) (*User, error) {
	return call[User](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/reactivate",
		pathArgs: []string{userID},
	}, opts)
}

// ReviewSuspension records a review of a suspended user: the reason, risk and review state.
//
// API: PUT /v1/users/{id}/suspension
func (s *UsersService) ReviewSuspension(ctx context.Context, userID string, params ReviewSuspensionParams, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPut,
		path:     "/v1/users/{userID}/suspension",
		pathArgs: []string{userID},
		body:     params,
	}, nil, opts)
}

// ListSessions lists a user's live sessions.
//
// API: GET /v1/users/{id}/sessions
func (s *UsersService) ListSessions(ctx context.Context, userID string, opts ...RequestOption) (*Page[UserSession], error) {
	return call[Page[UserSession]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/users/{userID}/sessions",
		pathArgs: []string{userID},
	}, opts)
}

// RevokeSessions ends every session of a user and revokes their application grants, remembered devices
// and pending authorizations.
//
// API: POST /v1/users/{id}/sessions/revoke-all
func (s *UsersService) RevokeSessions(ctx context.Context, userID string, opts ...RequestOption) (*SessionRevocation, error) {
	return call[SessionRevocation](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/sessions/revoke-all",
		pathArgs: []string{userID},
	}, opts)
}

// EndSession ends one of a user's sessions.
//
// API: DELETE /v1/tenants/{tenantID}/users/{id}/sessions/{sessionID}
func (s *UsersService) EndSession(ctx context.Context, tenantID string, userID string, sessionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/users/{userID}/sessions/{sessionID}",
		pathArgs: []string{tenantID, userID, sessionID},
	}, nil, opts)
}

// ListTrustedDevices lists the devices a user asked to be remembered.
//
// API: GET /v1/users/{id}/trusted-devices
func (s *UsersService) ListTrustedDevices(ctx context.Context, userID string, opts ...RequestOption) (*Page[TrustedDevice], error) {
	return call[Page[TrustedDevice]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/users/{userID}/trusted-devices",
		pathArgs: []string{userID},
	}, opts)
}

// ForgetTrustedDevices forgets every remembered device of a user.
//
// API: DELETE /v1/users/{id}/trusted-devices
func (s *UsersService) ForgetTrustedDevices(ctx context.Context, userID string, opts ...RequestOption) (*ForgottenDevices, error) {
	return call[ForgottenDevices](ctx, s.c, &request{
		method:   http.MethodDelete,
		path:     "/v1/users/{userID}/trusted-devices",
		pathArgs: []string{userID},
	}, opts)
}

// ForgetTrustedDevice forgets one remembered device of a user.
//
// API: DELETE /v1/users/{id}/trusted-devices/{deviceID}
func (s *UsersService) ForgetTrustedDevice(ctx context.Context, userID string, deviceID string, opts ...RequestOption) (*ForgottenDevices, error) {
	return call[ForgottenDevices](ctx, s.c, &request{
		method:   http.MethodDelete,
		path:     "/v1/users/{userID}/trusted-devices/{deviceID}",
		pathArgs: []string{userID, deviceID},
	}, opts)
}

// ListPushDevices lists the phones a user linked for push approval.
//
// API: GET /v1/users/{id}/push-devices
func (s *UsersService) ListPushDevices(ctx context.Context, userID string, opts ...RequestOption) (*Page[PushDevice], error) {
	return call[Page[PushDevice]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/users/{userID}/push-devices",
		pathArgs: []string{userID},
	}, opts)
}

// DisconnectPushDevice disconnects one of a user's linked phones.
//
// API: DELETE /v1/users/{id}/push-devices/{deviceID}
func (s *UsersService) DisconnectPushDevice(ctx context.Context, userID string, deviceID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/users/{userID}/push-devices/{deviceID}",
		pathArgs: []string{userID, deviceID},
	}, nil, opts)
}

// ListPasskeys lists a user's passkeys.
//
// API: GET /v1/tenants/{tenantID}/users/{id}/passkeys
func (s *UsersService) ListPasskeys(ctx context.Context, tenantID string, userID string, opts ...RequestOption) (*Page[UserPasskey], error) {
	return call[Page[UserPasskey]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/users/{userID}/passkeys",
		pathArgs: []string{tenantID, userID},
	}, opts)
}

// DeletePasskey removes one of a user's passkeys.
//
// API: DELETE /v1/tenants/{tenantID}/users/{id}/passkeys/{passkeyID}
func (s *UsersService) DeletePasskey(ctx context.Context, tenantID string, userID string, passkeyID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/users/{userID}/passkeys/{passkeyID}",
		pathArgs: []string{tenantID, userID, passkeyID},
	}, nil, opts)
}

// SignIns summarizes a user's sign-ins over a period: methods, addresses, browsers and failures.
//
// API: GET /v1/tenants/{tenantID}/users/{id}/sign-ins
func (s *UsersService) SignIns(ctx context.Context, tenantID string, userID string, params *SignInSummaryParams, opts ...RequestOption) (*SignInSummary, error) {
	return call[SignInSummary](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/users/{userID}/sign-ins",
		pathArgs: []string{tenantID, userID},
		query:    params.values(),
	}, opts)
}

// StartEmailVerification sends a verification code to the user's email address, or to the address in
// params.
//
// API: POST /v1/users/{id}/verify/email/start
func (s *UsersService) StartEmailVerification(ctx context.Context, userID string, params StartEmailVerificationParams, opts ...RequestOption) (*Message, error) {
	return call[Message](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/verify/email/start",
		pathArgs: []string{userID},
		body:     params,
	}, opts)
}

// ConfirmEmailVerification confirms a user's email address with the code StartEmailVerification sent.
//
// API: POST /v1/users/{id}/verify/email/confirm
func (s *UsersService) ConfirmEmailVerification(ctx context.Context, userID string, params ConfirmVerificationParams, opts ...RequestOption) (*Message, error) {
	return call[Message](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/verify/email/confirm",
		pathArgs: []string{userID},
		body:     params,
	}, opts)
}

// StartPhoneVerification sends a verification code to the user's phone, or to the number in params.
//
// API: POST /v1/users/{id}/verify/phone/start
func (s *UsersService) StartPhoneVerification(ctx context.Context, userID string, params StartPhoneVerificationParams, opts ...RequestOption) (*Message, error) {
	return call[Message](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/verify/phone/start",
		pathArgs: []string{userID},
		body:     params,
	}, opts)
}

// ConfirmPhoneVerification confirms a user's phone number with the code StartPhoneVerification sent.
//
// API: POST /v1/users/{id}/verify/phone/confirm
func (s *UsersService) ConfirmPhoneVerification(ctx context.Context, userID string, params ConfirmVerificationParams, opts ...RequestOption) (*Message, error) {
	return call[Message](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/users/{userID}/verify/phone/confirm",
		pathArgs: []string{userID},
		body:     params,
	}, opts)
}

// ListSocialIdentities lists the social accounts linked to a user.
//
// API: GET /v1/users/{userID}/social/identities
func (s *UsersService) ListSocialIdentities(ctx context.Context, userID string, opts ...RequestOption) (*Page[ExternalIdentity], error) {
	return call[Page[ExternalIdentity]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/users/{userID}/social/identities",
		pathArgs: []string{userID},
	}, opts)
}

// CreateUserParams is the request body of UsersService.Create.
type CreateUserParams struct {
	TenantID    string `json:"tenant_id,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	// E.164
	Phone    string         `json:"phone,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// User is a person with an account in a tenant.
type User struct {
	ID              string         `json:"id,omitempty"`
	TenantID        string         `json:"tenant_id,omitempty"`
	Email           string         `json:"email,omitempty"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	Phone           string         `json:"phone,omitempty"`
	PhoneVerifiedAt *time.Time     `json:"phone_verified_at,omitempty"`
	DisplayName     string         `json:"display_name,omitempty"`
	Status          UserStatus     `json:"status,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitzero"`
	UpdatedAt       time.Time      `json:"updated_at,omitzero"`
	// Role names in the tenant reading the user, on the list and GET /v1/users/{id}.
	Roles []string `json:"roles,omitempty"`
	// The same roles with their IDs.
	AssignedRoles []UserAssignedRole `json:"assigned_roles,omitempty"`
	// True if the user has a confirmed TOTP, a verified OTP factor, or a push device; on the list and GET
	// /v1/users/{id}.
	MFAEnabled bool `json:"mfa_enabled"`
	// Number of groups the user belongs to in this tenant; on the list and GET /v1/users/{id}.
	GroupsCount int `json:"groups_count"`
	// Their most recent session activity or sign-in in this tenant, or null; on the list and GET
	// /v1/users/{id}.
	LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
	// The address of that session or sign-in, or null.
	LastSeenIP string `json:"last_seen_ip,omitempty"`
	// The name of the organisation that manages the account (tenant_id).
	HomeTenantName string `json:"home_tenant_name,omitempty"`
	// On GET /v1/users/{id}, for a suspended person whose home is this tenant; absent otherwise.
	Suspension *UserSuspension `json:"suspension,omitempty"`
	// AvatarURL is the address of the person's picture, if they have one.
	AvatarURL string `json:"avatar_url,omitempty"`
}

// UserStatus is a string enumeration used by User.Status.
type UserStatus string

// Values of UserStatus. The server may add values; handle an unknown one gracefully.
const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusInvited   UserStatus = "invited"
)

// UserAssignedRole is one element of User.AssignedRoles.
type UserAssignedRole struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UserSuspension is the type of User.Suspension. On GET /v1/users/{id}, for a suspended person whose
// home is this tenant; absent otherwise.
type UserSuspension struct {
	Reason      string     `json:"reason,omitempty"`
	Source      string     `json:"source,omitempty"`
	SuspendedAt *time.Time `json:"suspended_at,omitempty"`
	// The person who last suspended them by hand; null for a suspension by policy or provisioning.
	SuspendedBy      string `json:"suspended_by,omitempty"`
	SuspendedByEmail string `json:"suspended_by_email,omitempty"`
}

// UpdateUserParams is the request body of UsersService.Update.
type UpdateUserParams struct {
	DisplayName *string        `json:"display_name,omitempty"`
	Phone       *string        `json:"phone,omitempty"`
	Status      *UserStatus    `json:"status,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// UserDeletion is the result of UsersService.Delete.
type UserDeletion struct {
	ID        string    `json:"id,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitzero"`
	DeletedBy string    `json:"deleted_by,omitempty"`
}

// MFAFilter is a string enumeration used by ListUsersParams.MFA.
type MFAFilter string

// Values of MFAFilter. The server may add values; handle an unknown one gracefully.
const (
	MFAFilterTrue  MFAFilter = "true"
	MFAFilterFalse MFAFilter = "false"
)

// LastSeenFilter is a string enumeration used by ListUsersParams.LastSeen.
type LastSeenFilter string

// Values of LastSeenFilter. The server may add values; handle an unknown one gracefully.
const (
	LastSeenFilterToday LastSeenFilter = "today"
	LastSeenFilterWeek  LastSeenFilter = "week"
	LastSeenFilterMonth LastSeenFilter = "month"
	LastSeenFilterOlder LastSeenFilter = "older"
	LastSeenFilterNever LastSeenFilter = "never"
)

// UserSort is a string enumeration used by ListUsersParams.Sort.
type UserSort string

// Values of UserSort. The server may add values; handle an unknown one gracefully.
const (
	UserSortCreated  UserSort = "created"
	UserSortName     UserSort = "name"
	UserSortLastSeen UserSort = "last_seen"
)

// SortOrder is a string enumeration used by ListUsersParams.Order.
type SortOrder string

// Values of SortOrder. The server may add values; handle an unknown one gracefully.
const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// ListUsersParams holds the optional query parameters of UsersService.List. A nil *ListUsersParams
// sends none.
type ListUsersParams struct {
	// Maximum items to return.
	Limit int
	// Opaque pagination cursor from a prior response's next_cursor.
	Cursor string
	// Offset for numbered paging (jump to page N). Alternative to cursor.
	Offset        int
	Q             string
	Role          string
	Department    string
	MFA           MFAFilter
	EmailVerified *bool
	Status        UserStatus
	LastSeen      LastSeenFilter
	Sort          UserSort
	Order         SortOrder
}

func (p *ListUsersParams) values() url.Values {
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
	if p.Offset != 0 {
		q.Set("offset", strconv.FormatInt(int64(p.Offset), 10))
	}
	if p.Q != "" {
		q.Set("q", p.Q)
	}
	if p.Role != "" {
		q.Set("role", p.Role)
	}
	if p.Department != "" {
		q.Set("department", p.Department)
	}
	if p.MFA != "" {
		q.Set("mfa", string(p.MFA))
	}
	if p.EmailVerified != nil {
		q.Set("email_verified", strconv.FormatBool(*p.EmailVerified))
	}
	if p.Status != "" {
		q.Set("status", string(p.Status))
	}
	if p.LastSeen != "" {
		q.Set("last_seen", string(p.LastSeen))
	}
	if p.Sort != "" {
		q.Set("sort", string(p.Sort))
	}
	if p.Order != "" {
		q.Set("order", string(p.Order))
	}
	return q
}

func (p *ListUsersParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListUsersParams) withCursor(cursor string) *ListUsersParams {
	var c ListUsersParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// SetPasswordParams is the request body of UsersService.SetPassword.
type SetPasswordParams struct {
	Password string `json:"password,omitempty"`
}

// BulkCreateUsersParams is the request body of UsersService.BulkCreate.
type BulkCreateUsersParams struct {
	// Optional; when present must match the caller's tenant.
	TenantID string `json:"tenant_id,omitempty"`
	// A role of the organisation to give every user the import makes, with the user in one transaction,
	// audited as role.assigned. The whole import is refused, making nobody, when the caller lacks
	// role.write or the authority the role carries (403 rbac.permission_denied,
	// rbac.grant_exceeds_caller, rbac.owner_only), or the role isn't the organisation's (404). Each
	// user's role takes one of the plan's seats; once they are all taken, the rest of the rows are
	// refused one by one with billing.plan_limit, making no account.
	RoleID string     `json:"role_id,omitempty"`
	Users  []BulkUser `json:"users,omitempty"`
}

// BulkUser is one element of BulkCreateUsersParams.Users.
type BulkUser struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	// E.164
	Phone string `json:"phone,omitempty"`
}

// BulkImportResult is the result of UsersService.BulkCreate.
type BulkImportResult struct {
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	// The users the import made, by their row.
	Created []BulkImportCreated `json:"created,omitempty"`
	Errors  []BulkImportError   `json:"errors,omitempty"`
}

// BulkImportCreated is one element of BulkImportResult.Created.
type BulkImportCreated struct {
	// The row's place in the request
	Line  int    `json:"line"`
	Email string `json:"email,omitempty"`
	ID    string `json:"id,omitempty"`
}

// BulkImportError is one element of BulkImportResult.Errors.
type BulkImportError struct {
	// POST /v1/users/bulk: the 1-based index in the users array. POST /v1/users/bulk/import: the line of
	// the export file (azure_b2c: the 1-based position in value). Absent when the whole file failed.
	Line  int    `json:"line"`
	Email string `json:"email,omitempty"`
	// Stable reason the row failed; word the error from this. Password, duplicate and domain refusals
	// reuse the API error codes.
	Code BulkImportErrorCode `json:"code,omitempty"`
	// The row's field the error concerns; absent when none.
	Field BulkImportErrorField `json:"field,omitempty"`
	// English text kept for compatibility; it may change. Show wording keyed on code.
	Message string `json:"message,omitempty"`
}

// BulkImportErrorCode is a string enumeration used by BulkImportError.Code. Stable reason the row
// failed; word the error from this. Password, duplicate and domain refusals reuse the API error codes.
type BulkImportErrorCode string

// Values of BulkImportErrorCode. The server may add values; handle an unknown one gracefully.
const (
	BulkImportErrorCodeUserImportFieldMissing        BulkImportErrorCode = "user.import_field_missing"
	BulkImportErrorCodeUserImportFieldInvalid        BulkImportErrorCode = "user.import_field_invalid"
	BulkImportErrorCodeUserImportFieldTooLong        BulkImportErrorCode = "user.import_field_too_long"
	BulkImportErrorCodeUserImportRowMalformed        BulkImportErrorCode = "user.import_row_malformed"
	BulkImportErrorCodeUserImportSourceInvalid       BulkImportErrorCode = "user.import_source_invalid"
	BulkImportErrorCodeAuthEmailExists               BulkImportErrorCode = "auth.email_exists"
	BulkImportErrorCodeUserDomainManagedElsewhere    BulkImportErrorCode = "user.domain_managed_elsewhere"
	BulkImportErrorCodeAuthpolicyPasswordTooShort    BulkImportErrorCode = "authpolicy.password_too_short"
	BulkImportErrorCodeAuthpolicyPasswordNoUppercase BulkImportErrorCode = "authpolicy.password_no_uppercase"
	BulkImportErrorCodeAuthpolicyPasswordNoNumber    BulkImportErrorCode = "authpolicy.password_no_number"
	BulkImportErrorCodeAuthpolicyPasswordNoSymbol    BulkImportErrorCode = "authpolicy.password_no_symbol"
	BulkImportErrorCodeAuthpolicyPasswordBreached    BulkImportErrorCode = "authpolicy.password_breached"
	BulkImportErrorCodeAuthPasswordWeak              BulkImportErrorCode = "auth.password_weak"
	BulkImportErrorCodeInternal                      BulkImportErrorCode = "internal"
)

// BulkImportErrorField is a string enumeration used by BulkImportError.Field. The row's field the
// error concerns; absent when none.
type BulkImportErrorField string

// Values of BulkImportErrorField. The server may add values; handle an unknown one gracefully.
const (
	BulkImportErrorFieldEmail       BulkImportErrorField = "email"
	BulkImportErrorFieldPassword    BulkImportErrorField = "password"
	BulkImportErrorFieldDisplayName BulkImportErrorField = "display_name"
	BulkImportErrorFieldPhone       BulkImportErrorField = "phone"
)

// ListDeletedUsersParams holds the optional query parameters of UsersService.ListDeleted. A nil
// *ListDeletedUsersParams sends none.
type ListDeletedUsersParams struct {
	// Maximum items to return.
	Limit int
}

func (p *ListDeletedUsersParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	return q
}

// DeletedUser is a soft-deleted user, as listed by UsersService.ListDeleted.
type DeletedUser struct {
	ID          string    `json:"id,omitempty"`
	Email       string    `json:"email,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitzero"`
	// The admin who deleted the user; null when not recorded.
	DeletedBy string    `json:"deleted_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}

// MembershipRemoval is the result of UsersService.RemoveMembership.
type MembershipRemoval struct {
	ID            string `json:"id,omitempty"`
	RolesRemoved  int    `json:"roles_removed"`
	SessionsEnded int    `json:"sessions_ended"`
}

// Message is the confirmation some actions answer with.
type Message struct {
	Message string `json:"message,omitempty"`
}

// SetMFARequiredParams is the request body of UsersService.SetMFARequired.
type SetMFARequiredParams struct {
	Required bool `json:"required"`
}

// MFARequirement is the result of UsersService.SetMFARequired.
type MFARequirement struct {
	MFARequired bool `json:"mfa_required"`
}

// UserSecurity is the result of UsersService.Security.
type UserSecurity struct {
	MFARequired            bool       `json:"mfa_required"`
	MFAEnabled             bool       `json:"mfa_enabled"`
	TOTPEnabled            bool       `json:"totp_enabled"`
	OTPFactors             int        `json:"otp_factors"`
	PushDevices            int        `json:"push_devices"`
	Passkeys               int        `json:"passkeys"`
	PasskeyLastUsedAt      *time.Time `json:"passkey_last_used_at,omitempty"`
	RecoveryCodesRemaining int        `json:"recovery_codes_remaining"`
	PasswordSet            bool       `json:"password_set"`
	PasswordChangedAt      *time.Time `json:"password_changed_at,omitempty"`
	ActiveSessions         int        `json:"active_sessions"`
	DistinctDevices        int        `json:"distinct_devices"`
}

// UserAccess is the result of UsersService.Access.
type UserAccess struct {
	Organization         *UserAccessOrganization `json:"organization,omitempty"`
	Roles                []UserAccessRole        `json:"roles,omitempty"`
	Groups               []UserAccessGroup       `json:"groups,omitempty"`
	ApplicationsCount    int                     `json:"applications_count"`
	PoliciesCount        int                     `json:"policies_count"`
	PermissionsDirect    int                     `json:"permissions_direct"`
	PermissionsInherited int                     `json:"permissions_inherited"`
	PermissionsTotal     int                     `json:"permissions_total"`
}

// UserAccessOrganization is the type of UserAccess.Organization.
type UserAccessOrganization struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// UserAccessRole is one element of UserAccess.Roles.
type UserAccessRole struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UserAccessGroup is one element of UserAccess.Groups.
type UserAccessGroup struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UserStats is the result of UsersService.Stats.
type UserStats struct {
	Total      int               `json:"total"`
	Active     int               `json:"active"`
	Suspended  int               `json:"suspended"`
	Invited    int               `json:"invited"`
	MFAEnabled int               `json:"mfa_enabled"`
	MFAMissing int               `json:"mfa_missing"`
	NewLast30d int               `json:"new_last_30d"`
	LastSeen   UserStatsLastSeen `json:"last_seen,omitzero"`
	// Members holding the built-in owner or admin role with no confirmed second factor.
	AdminsWithoutMFA int             `json:"admins_without_mfa"`
	EmailUnverified  int             `json:"email_unverified"`
	ByRole           []UserRoleCount `json:"by_role,omitempty"`
	// Up to 200 departments, largest first.
	ByDepartment []UserDepartmentCount `json:"by_department,omitempty"`
}

// UserStatsLastSeen is the type of UserStats.LastSeen.
type UserStatsLastSeen struct {
	Today int `json:"today"`
	Week  int `json:"week"`
	Month int `json:"month"`
	Older int `json:"older"`
	Never int `json:"never"`
}

// UserRoleCount is one element of UserStats.ByRole.
type UserRoleCount struct {
	Role  string `json:"role,omitempty"`
	Count int    `json:"count"`
}

// UserDepartmentCount is one element of UserStats.ByDepartment.
type UserDepartmentCount struct {
	Department string `json:"department,omitempty"`
	Count      int    `json:"count"`
}

// UserTrends is the result of UsersService.Trends.
type UserTrends struct {
	// 30-day daily cumulative member count (oldest first).
	Total []int `json:"total,omitempty"`
	// 30-day daily cumulative MFA-enrolled count (oldest first).
	MFAEnabled []int `json:"mfa_enabled,omitempty"`
}

// SuspensionRisk is a string enumeration used by ListSuspendedUsersParams.Risk.
type SuspensionRisk string

// Values of SuspensionRisk. The server may add values; handle an unknown one gracefully.
const (
	SuspensionRiskUnknown SuspensionRisk = "unknown"
	SuspensionRiskLow     SuspensionRisk = "low"
	SuspensionRiskMedium  SuspensionRisk = "medium"
	SuspensionRiskHigh    SuspensionRisk = "high"
)

// SuspensionReview is a string enumeration used by ListSuspendedUsersParams.Review.
type SuspensionReview string

// Values of SuspensionReview. The server may add values; handle an unknown one gracefully.
const (
	SuspensionReviewPending  SuspensionReview = "pending"
	SuspensionReviewReviewed SuspensionReview = "reviewed"
	SuspensionReviewExternal SuspensionReview = "external"
)

// ListSuspendedUsersParams holds the optional query parameters of UsersService.ListSuspended. A nil
// *ListSuspendedUsersParams sends none.
type ListSuspendedUsersParams struct {
	Q      string
	Reason string
	Source string
	Team   string
	Risk   SuspensionRisk
	Review SuspensionReview
	Limit  int
	Offset int
}

func (p *ListSuspendedUsersParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Q != "" {
		q.Set("q", p.Q)
	}
	if p.Reason != "" {
		q.Set("reason", p.Reason)
	}
	if p.Source != "" {
		q.Set("source", p.Source)
	}
	if p.Team != "" {
		q.Set("team", p.Team)
	}
	if p.Risk != "" {
		q.Set("risk", string(p.Risk))
	}
	if p.Review != "" {
		q.Set("review", string(p.Review))
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Offset != 0 {
		q.Set("offset", strconv.FormatInt(int64(p.Offset), 10))
	}
	return q
}

// SuspendedUsers is the result of UsersService.ListSuspended.
type SuspendedUsers struct {
	Items   []SuspendedUser        `json:"items,omitempty"`
	Total   int                    `json:"total"`
	Limit   int                    `json:"limit"`
	Offset  int                    `json:"offset"`
	Summary SuspendedUsersSummary  `json:"summary,omitzero"`
	Reasons []SuspendedUsersReason `json:"reasons,omitempty"`
	Risks   []SuspendedUsersRisk   `json:"risks,omitempty"`
	Queue   []SuspendedUser        `json:"queue,omitempty"`
	Teams   []string               `json:"teams,omitempty"`
}

// SuspendedUser is a suspended person with the details of their suspension and its review.
type SuspendedUser struct {
	ID          string           `json:"id,omitempty"`
	Email       string           `json:"email,omitempty"`
	DisplayName string           `json:"display_name,omitempty"`
	AvatarURL   string           `json:"avatar_url,omitempty"`
	Team        string           `json:"team,omitempty"`
	Reason      string           `json:"reason,omitempty"`
	Source      string           `json:"source,omitempty"`
	Risk        SuspensionRisk   `json:"risk,omitempty"`
	Review      SuspensionReview `json:"review,omitempty"`
	SuspendedAt *time.Time       `json:"suspended_at,omitempty"`
	LastSignIn  *time.Time       `json:"last_sign_in,omitempty"`
	CanManage   bool             `json:"can_manage"`
}

// SuspendedUsersSummary is the type of SuspendedUsers.Summary.
type SuspendedUsersSummary struct {
	Total    int `json:"total"`
	HighRisk int `json:"high_risk"`
	Policy   int `json:"policy"`
	Manual   int `json:"manual"`
	Pending  int `json:"pending"`
}

// SuspendedUsersReason is one element of SuspendedUsers.Reasons.
type SuspendedUsersReason struct {
	Key   string `json:"key,omitempty"`
	Count int    `json:"count"`
}

// SuspendedUsersRisk is one element of SuspendedUsers.Risks.
type SuspendedUsersRisk struct {
	Key   string `json:"key,omitempty"`
	Count int    `json:"count"`
}

// ReviewSuspensionParams is the request body of UsersService.ReviewSuspension.
type ReviewSuspensionParams struct {
	Reason SuspensionReason `json:"reason,omitempty"`
	Risk   SuspensionRisk   `json:"risk,omitempty"`
	Review SuspensionReview `json:"review,omitempty"`
}

// SuspensionReason is a string enumeration used by ReviewSuspensionParams.Reason.
type SuspensionReason string

// Values of SuspensionReason. The server may add values; handle an unknown one gracefully.
const (
	SuspensionReasonUnspecified        SuspensionReason = "unspecified"
	SuspensionReasonProvider           SuspensionReason = "provider"
	SuspensionReasonSuspiciousActivity SuspensionReason = "suspicious_activity"
	SuspensionReasonPolicyViolation    SuspensionReason = "policy_violation"
	SuspensionReasonAccountCompromise  SuspensionReason = "account_compromise"
	SuspensionReasonOffboarding        SuspensionReason = "offboarding"
	SuspensionReasonInactive           SuspensionReason = "inactive"
	SuspensionReasonOther              SuspensionReason = "other"
)

// UserSession is one of a user's live sign-in sessions.
type UserSession struct {
	ID         string    `json:"id,omitempty"`
	IP         string    `json:"ip,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitzero"`
	LastSeenAt time.Time `json:"last_seen_at,omitzero"`
	// The device and browser in words ("Chrome on macOS"), as the person's own session list says.
	Device     string     `json:"device,omitempty"`
	DeviceKind DeviceKind `json:"device_kind,omitempty"`
	// Where the session was opened
	Country string `json:"country,omitempty"`
	// The city it was opened from
	City string `json:"city,omitempty"`
}

// DeviceKind is a string enumeration used by UserSession.DeviceKind.
type DeviceKind string

// Values of DeviceKind. The server may add values; handle an unknown one gracefully.
const (
	DeviceKindDesktop DeviceKind = "desktop"
	DeviceKindMobile  DeviceKind = "mobile"
	DeviceKindTablet  DeviceKind = "tablet"
	DeviceKindUnknown DeviceKind = "unknown"
)

// SessionRevocation is the result of UsersService.RevokeSessions.
type SessionRevocation struct {
	// Number of sign-ins ended, all places together.
	Revoked int `json:"revoked"`
	// Sessions ended.
	Sessions int `json:"sessions"`
	// Hosted-login browser sessions ended.
	HostedLoginSessions int `json:"hosted_login_sessions"`
	// Application refresh grants ended.
	AppGrants int `json:"app_grants"`
	// Devices remembered to skip the second factor that were forgotten; each is asked for the second
	// factor at its next sign-in.
	TrustedDevices int `json:"trusted_devices"`
	// Unredeemed authorisation codes spent, and device or back-channel sign-in requests (pending, or
	// approved but not collected) denied.
	PendingAuthorizations int `json:"pending_authorizations"`
}

// TrustedDevice is a browser a user asked Qeet ID to remember, which skips the second factor until it
// expires.
type TrustedDevice struct {
	ID string `json:"id,omitempty"`
	// The organisation the device was remembered in.
	TenantID string `json:"tenant_id,omitempty"`
	// The device and browser in words, such as "Chrome on macOS".
	Device     string     `json:"device,omitempty"`
	DeviceKind DeviceKind `json:"device_kind,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitzero"`
	// When the device last skipped the second factor.
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at,omitzero"`
	// The request came from this device.
	Current bool `json:"current"`
}

// ForgottenDevices is the result of UsersService.ForgetTrustedDevices.
type ForgottenDevices struct {
	// How many devices were forgotten.
	Forgotten int `json:"forgotten"`
}

// PushDevice is a phone a user linked to approve sign-ins with push notifications.
type PushDevice struct {
	ID         string             `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Platform   PushDevicePlatform `json:"platform,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitzero"`
	LastSeenAt time.Time          `json:"last_seen_at,omitzero"`
}

// PushDevicePlatform is a string enumeration used by PushDevice.Platform.
type PushDevicePlatform string

// Values of PushDevicePlatform. The server may add values; handle an unknown one gracefully.
const (
	PushDevicePlatformIos     PushDevicePlatform = "ios"
	PushDevicePlatformAndroid PushDevicePlatform = "android"
)

// UserPasskey is a passkey registered to a user.
type UserPasskey struct {
	ID         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitzero"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// SignInSummaryParams holds the optional query parameters of UsersService.SignIns. A nil
// *SignInSummaryParams sends none.
type SignInSummaryParams struct {
	From time.Time
	To   time.Time
	Days int
}

func (p *SignInSummaryParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.UTC().Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.UTC().Format(time.RFC3339))
	}
	if p.Days != 0 {
		q.Set("days", strconv.FormatInt(int64(p.Days), 10))
	}
	return q
}

// SignInSummary is the result of UsersService.SignIns.
type SignInSummary struct {
	From             time.Time            `json:"from,omitzero"`
	To               time.Time            `json:"to,omitzero"`
	SignIns          int                  `json:"sign_ins"`
	WithSecondStep   int                  `json:"with_second_step"`
	ByMethod         []SignInMethodCount  `json:"by_method,omitempty"`
	ByAddress        []SignInAddressCount `json:"by_address,omitempty"`
	ByBrowser        []SignInBrowserCount `json:"by_browser,omitempty"`
	Failures         int                  `json:"failures"`
	FailuresByReason []SignInFailureCount `json:"failures_by_reason,omitempty"`
}

// SignInMethodCount is one element of SignInSummary.ByMethod.
type SignInMethodCount struct {
	Method         string `json:"method,omitempty"`
	Count          int    `json:"count"`
	WithSecondStep int    `json:"with_second_step"`
}

// SignInAddressCount is one element of SignInSummary.ByAddress.
type SignInAddressCount struct {
	IP     string    `json:"ip,omitempty"`
	Count  int       `json:"count"`
	LastAt time.Time `json:"last_at,omitzero"`
}

// SignInBrowserCount is one element of SignInSummary.ByBrowser.
type SignInBrowserCount struct {
	Device string `json:"device,omitempty"`
	Kind   string `json:"kind,omitempty"`
	Count  int    `json:"count"`
}

// SignInFailureCount is one element of SignInSummary.FailuresByReason.
type SignInFailureCount struct {
	Reason string `json:"reason,omitempty"`
	Count  int    `json:"count"`
}

// StartEmailVerificationParams is the request body of UsersService.StartEmailVerification.
type StartEmailVerificationParams struct {
	// Optional; defaults to the user's email on file.
	Email string `json:"email,omitempty"`
}

// ConfirmVerificationParams is the request body of UsersService.ConfirmEmailVerification.
type ConfirmVerificationParams struct {
	Code string `json:"code,omitempty"`
}

// StartPhoneVerificationParams is the request body of UsersService.StartPhoneVerification.
type StartPhoneVerificationParams struct {
	// Optional; defaults to the user's phone on file.
	Phone string `json:"phone,omitempty"`
}

// ExternalIdentity is a social account linked to a user.
type ExternalIdentity struct {
	Email    string    `json:"email,omitempty"`
	ID       string    `json:"id,omitempty"`
	LinkedAt time.Time `json:"linked_at,omitzero"`
	Provider string    `json:"provider,omitempty"`
	Subject  string    `json:"subject,omitempty"`
	TenantID string    `json:"tenant_id,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
}
