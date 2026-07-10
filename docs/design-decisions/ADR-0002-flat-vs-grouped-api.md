# ADR-0002: Flat `Client` fields, not nested namespace groups

**Status:** Accepted (superseded an earlier grouped design in the same
implementation pass)

## Context

At ~35 resources, `Client` has a lot of fields. One proposal: nest them
under four namespace structs — `Client.Authentication.Sessions`,
`Client.Authorization.Roles`, `Client.Identity.Users`,
`Client.Administration.Webhooks` — so autocomplete on `client.` shows four
options instead of thirty-five.

This was actually implemented once during this SDK's restructure, then
reverted in favor of the flat design below.

## Decision

Every `*XService` is a direct field on `Client`. The four categories
(Identity / Authentication / Authorization / Administration) exist only as
comment banners grouping the field declarations in `client.go` — not as Go
types.

```go
type Client struct {
	// Identity
	Users *UsersService
	...
	// Authentication
	Sessions *SessionsService
	...
}
```

## Why

- **Fewer keystrokes on the common path.** `client.Users.Create(...)` vs.
  `client.Identity.Users.Create(...)` — the extra segment buys
  autocomplete-list brevity at the cost of every single call site being
  longer.
- **Matches Stripe/WorkOS/AWS SDK precedent.** These are widely regarded as
  well-designed server SDKs, and none group service clients into
  intermediate namespace objects at comparable or larger surface sizes.
- **Comment banners solve the actual problem.** The concern was
  discoverability/mental model, not runtime behavior — a well-organized
  `client.go` (and a grouped table in the README) gives a reader the same
  mental map without changing how every call site reads.
- **No import-cycle costs.** Unlike ADR-0001's rejected alternative, nested
  *structs* (not packages) don't cause import cycles — but they still don't
  earn their keep given the point above.

## Rejected alternative: four namespace groups

Implemented as `AuthenticationGroup`, `AuthorizationGroup`, `IdentityGroup`,
`AdministrationGroup` structs, each holding a subset of services, with
`Client` holding one field per group. Reverted once weighed against the
ergonomic cost on every call site in the SDK's own examples and tests.

## Related: `Permissions.Check`/`CheckAll`/`Explain`

The RBAC check methods (`Can`/`CanAll`/`Explain` in an earlier draft) moved
onto `PermissionsService` rather than staying as bespoke `Client`-level (or,
in the reverted design, `AuthorizationGroup`-level) methods. A permission
check is a `Permissions` operation; giving it a home there means there's one
place — not a special case — to look for anything permission-related.
