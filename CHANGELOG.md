# Changelog

All notable changes to the Qeet ID Go SDK are documented here. The format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] — 2026-07-10

Initial release.

### Added

**Client**
- `qeetid.New(qeetid.Config{...})` — one client, every resource a
  `*XService` field directly on `Client` (comment-banner grouped in
  `client.go`: Identity / Authentication / Authorization / Administration —
  no nesting).
- `Config`: `APIKey`, `BaseURL`, `HTTPClient`, `Timeout`, `MaxRetries`,
  `Headers` (`http.Header`), `UserAgent`, `Logger`.
- Zero third-party dependencies — standard library only.

**Identity** — `Users` (CRUD, bulk create/import, recycle-bin, MFA reset,
email/phone verification, auto-pagination), `Organizations`,
`ServicePrincipals`, `Agents` (ephemeral tokens + full
suspend/resume/decommission/kill-all/sponsor-transfer lifecycle), `Domains`.

**Authentication** — `Sessions` (local ES256 JWT verification against
cached JWKS), `OAuth` (RFC 8693 token exchange, RFC 7662 introspection, RFC
7009 revocation, an MCP token guard, RFC 8628 device flow, CIBA, signing-key
rotation, grant/device-session admin views), `OIDC` (tenant-scoped client
CRUD + shadow-AI discovery/review), `SAML` (Qeet ID as SP), `SAMLProviders`
(Qeet ID as IdP — the mirror image of `SAML`, a distinct resource despite
the shared URL prefix), `SCIM` (tenant-admin provisioning config), `LDAP`
(connections, bind test, public authenticate passthrough), `Social`
(provider config, linked identities), `MFA` (admin-initiated factor reset —
the backend has no admin endpoint to list a user's factors), `Credentials`
(W3C Verifiable Credentials), `AuthHooks` (a multi-record collection:
list/create/update-by-id/delete-by-id, not a singleton), `AuthPolicy`,
`Policy` (an older, broader combined security-policy record alongside
`AuthPolicy`/`IPRules`), `IPRules` (rules, a dry-run `Check`, and
enforcement on/off).

**Authorization** — `Roles` (tenant-scoped create/list, user assignment,
permission grants — there is no per-role Get/Update/Delete in the backend),
`Permissions` (the platform's permission catalog, a user's effective
permissions, plus `Check`/`CheckAll`/`Explain` — RBAC with grant-path
explanation), `Groups` (membership + role bindings), `Relationships`
(Zanzibar-style ReBAC: tuple CRUD, recursive `Check` with `explain`, `Graph`
identity-graph expansion), `Decisions` (AuthZEN unified evaluation fronting
both RBAC and ReBAC).

**Administration** — `Branding`, `Invitations` (accepting an invite is
deliberately not wrapped — like login/signup, it's an end-user auth action),
`EmailTemplates` (a fixed catalog of resolved templates — list, get,
override, reset-to-default, preview-render), `APIKeys` (create, revoke,
list — no per-key Get or Rotate in the backend), `Vault` (encrypted secrets
store), `Webhooks` (management CRUD + HMAC-SHA256 inbound-delivery
verification via `ConstructEvent`/`ConstructEventFromRequest` — most
operations are scoped by the caller's own API key, not a tenant path
segment, and there is no Update), `AuditLogs` (hash-chain read +
whole-chain `Verify` + free-text search + behavioral-baseline anomaly
detection), `Analytics` (dashboard overview), `GDPR` (erasure + export
requests), `Billing`, `Retention`, `RateLimits`, `LogSinks` (SIEM
forwarding), `AdminLinks` (delegated admin-portal links).

**Cross-cutting**
- Auto-pagination: every listable resource has an `All(...)` method
  returning an `iter.Seq2[T, error]` (Go 1.23+ range-over-func).
- Automatic retry with exponential backoff + jitter on `429`/`5xx`
  (idempotent requests only for `5xx`).
- Typed `*qeetid.Error` with `Status`/`Code`/`Message`/`RequestID`/
  `RetryAfterSeconds` and `Is*` predicate helpers.
- Zero-config OIDC discovery (`Discover`, `NewFromDiscovery`) from
  `/.well-known/openid-configuration`.
- Optional per-request `Logger` observability hook — no logging dependency
  baked into the core.

### Notes

Every resource, path, and field shape in this release was verified directly
against the live Qeet ID backend's OpenAPI specs and (where the spec body
was still generic) the actual Go handler/domain structs — not assumed from
documentation. See [docs/architecture/ARCHITECTURE.md](docs/architecture/ARCHITECTURE.md)
for the internal package layout and [docs/design-decisions/](docs/design-decisions/)
for the API-shape trade-offs (single package vs. per-resource packages, flat
vs. grouped `Client` fields).

[0.1.0]: https://github.com/qeetgroup/qeet-id-go/releases/tag/v0.1.0
