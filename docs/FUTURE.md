# Future

Capabilities common competitor CIAM SDKs (Auth0, Clerk, WorkOS) expose that
qeet-id-go doesn't, because the Qeet ID *backend* doesn't have them yet.
These aren't SDK gaps to fill with stub methods — a stub that compiles but
throws `ErrNotImplemented` at runtime is worse than not having the method at
all. They're named here, each tied to a real tracked item, so this SDK adds
real coverage the moment the backend ships them.

## Nested org hierarchy / tenancy primitives

Tracked: qeet-id ROADMAP.md item 3.5. Today every tenant is flat — no
parent/child organization relationships. Once the backend introduces
`platform/tenancy` and a `parent_id` on organizations, this SDK adds it to
`OrganizationsService`.

## SPIFFE / workload identity (JWT-SVID)

Tracked: qeet-id ROADMAP.md item 3.6. CNCF-standard workload identity for
service-to-service auth — genuinely industry-wide whitespace, not just a
Qeet ID gap, but not yet built.

## Tenant-admin notification/broadcast management

Confirmed absent from the current API surface entirely: `GET /v1/notifications`
and `POST /v1/notifications/mark-all-read` exist, but only for the
*authenticated end-user's own inbox* — there is no tenant-admin path for
configuring or broadcasting notifications (no
`/v1/tenants/{tenantID}/notifications` or equivalent). If this becomes a
requirement, it needs a backend API first; there's nothing here for an SDK
to wrap yet.
