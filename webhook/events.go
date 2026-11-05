package webhook

import (
	"encoding/json"
	"errors"
	"time"
)

// Event is a webhook delivery's envelope.
type Event struct {
	// ID identifies the event. It is the same in every delivery of the event,
	// retries included, so deduplicate on it.
	ID string `json:"id"`
	// Type names the event, such as UserCreated. An event's name is its
	// version: an incompatible change to its data comes as a new name, while
	// new fields may be added to any event's data.
	Type EventType `json:"type"`
	// TenantID is the tenant the event belongs to.
	TenantID string `json:"tenant_id"`
	// OccurredAt is when the event happened.
	OccurredAt time.Time `json:"occurred_at"`
	// Data is the event's own payload. Decode it with DecodeData.
	Data json.RawMessage `json:"data"`
	// DeliveryID is the X-Qeet-Delivery-ID header of the delivery that
	// carried the event. It is not part of the signed body.
	DeliveryID string `json:"-"`
}

// DecodeData decodes the event's data into v, typically one of the data types
// of this package that matches e.Type.
func (e *Event) DecodeData(v any) error {
	if len(e.Data) == 0 {
		return errors.New("webhook: the event has no data")
	}
	return json.Unmarshal(e.Data, v)
}

// EventType names a webhook event. The server may add event types; ignore
// the ones you do not handle.
type EventType string

// The event catalogue: every event a subscription can name. List it at run
// time with WebhooksService.ListEventTypes in the root package.
const (
	UserCreated         EventType = "user.created"          // a person was added, by an administrator or over SCIM
	UserUpdated         EventType = "user.updated"          // SCIM changed a person's profile or status
	UserSuspended       EventType = "user.suspended"        // a person's account was suspended
	UserReactivated     EventType = "user.reactivated"      // a suspended account was reactivated
	UserDeleted         EventType = "user.deleted"          // a person was deleted
	UserPasswordSet     EventType = "user.password_set"     // an administrator set a password, ending the person's sessions
	UserBulkImported    EventType = "user.bulk_imported"    // a bulk import finished
	AuthPasswordChanged EventType = "auth.password_changed" // a person changed their own password
	SessionRevoked      EventType = "session.revoked"       // a session ended; its tokens stop working
	TokenClaimsChange   EventType = "token.claims_change"   // a person's roles changed, so their next token differs

	GroupCreated         EventType = "group.created"
	GroupUpdated         EventType = "group.updated"
	GroupDeleted         EventType = "group.deleted"
	GroupMemberAdded     EventType = "group.member_added"
	GroupMemberRemoved   EventType = "group.member_removed"
	GroupTemplateApplied EventType = "group_template.applied"

	DomainCreated  EventType = "domain.created"
	DomainVerified EventType = "domain.verified"
	DomainDeleted  EventType = "domain.deleted"

	AgentCreated            EventType = "agent.created"
	AgentSuspended          EventType = "agent.suspended"
	AgentResumed            EventType = "agent.resumed"
	AgentDecommissioned     EventType = "agent.decommissioned"
	AgentDeleted            EventType = "agent.deleted"
	AgentSponsorTransferred EventType = "agent.sponsor_transferred"
	AgentKillSwitch         EventType = "agent.kill_switch"
	AgentResumeAll          EventType = "agent.resume_all"

	RelationTupleWritten EventType = "rebac.tuple_written"
	RelationTupleDeleted EventType = "rebac.tuple_deleted"

	ABACPolicyCreated EventType = "abac_policy.created"
	ABACPolicyUpdated EventType = "abac_policy.updated"
	ABACPolicyDeleted EventType = "abac_policy.deleted"

	AuthHookCreated EventType = "auth_hook.created"
	AuthHookUpdated EventType = "auth_hook.updated"
	AuthHookDeleted EventType = "auth_hook.deleted"

	LogSinkCreated EventType = "siem.sink_created"
	LogSinkUpdated EventType = "siem.sink_updated"
	LogSinkDeleted EventType = "siem.sink_deleted"

	RateLimitsUpdated EventType = "rate_limits.updated"
	RateLimitsDeleted EventType = "rate_limits.deleted"

	CredentialIssued  EventType = "credential.issued"
	CredentialRevoked EventType = "credential.revoked"

	OIDCClientLockedOut EventType = "oidc.client_locked_out"

	TenantSuspended   EventType = "tenant.suspended"
	TenantReactivated EventType = "tenant.reactivated"
	TenantDeleted     EventType = "tenant.deleted"

	// TestPing is the event WebhooksService.Test sends. No subscription names
	// it, and every subscription can receive it.
	TestPing EventType = "test.ping"
)

// UserData is the data of the user events. Every one carries UserID and
// TenantID; Email and Source are set where the event's origin has them
// (Source is "admin" or "scim").
type UserData struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email,omitempty"`
	Source   string `json:"source,omitempty"`
}

// GroupMemberData is the data of GroupMemberAdded and GroupMemberRemoved.
type GroupMemberData struct {
	GroupID  string `json:"group_id"`
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	// Source is "scim", "saml" or "ldap" when a directory or sign-in
	// connection changed the membership; empty for an administrator.
	Source string `json:"source,omitempty"`
}

// SessionRevokedData is the data of SessionRevoked: a session ended before its
// tokens expired. An application that verifies access tokens offline should
// stop accepting the session's tokens on receipt.
type SessionRevokedData struct {
	// EventID identifies this revocation signal.
	EventID string `json:"event_id"`
	// UserID is the person whose session ended.
	UserID string `json:"user_id"`
	// SessionID is the session that ended; it is the sid claim of its tokens.
	SessionID string `json:"session_id"`
	// TenantID is the session's tenant; empty for a session without one.
	TenantID string `json:"tenant_id,omitempty"`
	// Reason says why the session ended.
	Reason RevocationReason `json:"reason"`
	// RevokedAt is when it ended.
	RevokedAt time.Time `json:"revoked_at"`
}

// RevocationReason says why a session ended.
type RevocationReason string

// Reasons a session ends. The server may add reasons.
const (
	ReasonLogout                  RevocationReason = "logout"
	ReasonReuseDetected           RevocationReason = "reuse_detected"
	ReasonOtherSessionsRevoked    RevocationReason = "other_sessions_revoked"
	ReasonAllSessionsRevoked      RevocationReason = "all_sessions_revoked"
	ReasonAdminRevokedAll         RevocationReason = "admin_revoked_all"
	ReasonAdminRevoked            RevocationReason = "admin_revoked"
	ReasonPasswordChanged         RevocationReason = "password_changed"
	ReasonPasswordReset           RevocationReason = "password_reset"
	ReasonAdminPasswordSet        RevocationReason = "admin_password_set"
	ReasonTenantSwitch            RevocationReason = "tenant_switch"
	ReasonSuspended               RevocationReason = "suspended"
	ReasonDeprovisioned           RevocationReason = "deprovisioned"
	ReasonReactivated             RevocationReason = "reactivated"
	ReasonDeleted                 RevocationReason = "deleted"
	ReasonRemovedFromOrganization RevocationReason = "removed_from_organization"
	ReasonTenantSuspended         RevocationReason = "tenant_suspended"
	ReasonTenantDeleted           RevocationReason = "tenant_deleted"
	ReasonExpired                 RevocationReason = "expired"
	ReasonSessionLimit            RevocationReason = "session_limit"
)

// ClaimsChangeData is the data of TokenClaimsChange: a person's roles changed,
// so the access token they hold no longer matches what their next one will
// carry.
type ClaimsChangeData struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	// Change is "role_assigned" or "role_unassigned".
	Change string `json:"change"`
}

// TestPingData is the data of TestPing.
type TestPingData struct {
	SubscriptionID string `json:"subscription_id"`
}
