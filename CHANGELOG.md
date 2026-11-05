# Changelog

All notable changes to the Qeet ID Go SDK are documented here. The format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

A rebuild of the SDK against the Qeet ID 1.0 API contract. No version of the
module has been tagged yet, so the changes below are measured against the
0.1.0 draft described further down.

### Breaking changes

- **Client construction.** `qeetid.New(qeetid.Config{...})` is replaced by
  `qeetid.NewClient(opts ...Option) (*Client, error)` with functional options:
  `WithAPIKey`, `WithClientCredentials`, `WithTokenSource`, `WithBaseURL`,
  `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`, `WithRetryDelay`,
  `WithLogger` (`log/slog`), `WithUserAgentSuffix` and `WithHeader`. The
  `Logger` interface is gone. A base URL must use https (http only on a
  loopback host).
- **Every method** takes `context.Context` first and accepts
  `...RequestOption` last; list and filter parameters are pointers to
  parameter structs (`nil` sends none); request bodies are `...Params`
  structs.
- **Errors.** `qeetid.Error` is replaced by `*APIError`, which carries the
  whole error envelope (`Code`, `Message`, `Detail`, `Fields`, `Metadata`,
  `RequestID`, `RetryAfter`), and `*OAuthError` for the RFC 6749 errors of
  the OAuth endpoints. Sentinels `ErrInvalidRequest`, `ErrUnauthorized`,
  `ErrSessionRevoked`, `ErrPaymentRequired`, `ErrForbidden`, `ErrNotFound`,
  `ErrConflict`, `ErrRateLimited`, `ErrUnavailable` and `ErrNoCredentials`
  work with `errors.Is`. The `Is*` predicate methods are gone.
- **Lists** return `*Page[T]` (`Items`, `NextCursor`, `Truncated`, `Total`)
  instead of slices or per-resource page types; cursor lists keep an `All`
  iterator, which now yields `*T`.
- **Renamed services:** `Organizations` → `Tenants` (and `Organization` →
  `Tenant`), `Vault` → `Secrets`, `Policy` → `SecurityPolicy`,
  `Decisions` → `AuthZEN`, `Social` → `SocialProviders`, `SAML` → `SAMLConnections`, `OIDC`
  → `OIDCClients`, `AuditLogs` → `Audit`, `AdminLinks` → `AdminPortalLinks`;
  `BotDetection` and `RiskSettings` are folded into `Threats`, and `MFA.Reset`
  into `Users.ResetMFA`. `OAuth.Grants`, `OAuth.Devices` and
  `OAuth.SigningKeys` become `OAuthGrants`, `DeviceAuthorizations` and
  `SigningKeys`.
- **Token verification** moves from `Client.Sessions` to package `jwt`. The
  new verifier requires an issuer and an audience and refuses a token
  without a `kid`.
- **Webhook verification** moves from the root package to package `webhook`.
  The old functions signed the body alone and could not verify a real
  delivery.
- **Invitations:** `CreateInvitationResult.Token` is gone; the server sends
  the accept link only to the invited address. `Invitations.Create` returns
  an `InvitationResult` whose invitation reports whether its email went out.
- **LDAP sign-in:** `LDAP.Authenticate` returns a `SignInResult`, either a
  `*TokenPair` or an `*MFAChallenge`.
- **OAuth:** the token-exchange, introspection, revocation, device and
  backchannel calls take parameter structs and authenticate the calling
  application; token exchange now goes to `/v1/oauth/token-code`. The
  cookie-only `DeviceContext`, `DeviceDecision`, `BackchannelPending` and
  `BackchannelDecision`, the client-side `OAuth.Verify` helper,
  `Permissions.CheckAll` (use `AuthZEN.EvaluateBatch`) and
  `OAuth.SigningKeys.Rotate` (the endpoint no longer exists) are removed.
- **Models** follow the contract's field names and types: timestamps are
  `time.Time`, enumerations are typed string constants, and many fields
  were added, renamed or removed (see "Fixed").
- Go 1.26 or later is required.

### Added

- Service principal authentication with `WithClientCredentials`: tokens are
  cached, shared and renewed before they expire.
- Retries with exponential backoff and jitter that run only when safe, honour
  `Retry-After` (up to `MaxRetryAfter`) and never retry a 401 or
  `auth.session_revoked`.
- `Idempotency-Key` support (`WithIdempotencyKey`) on the nine calls that
  accept it; the client sends a random key with them by itself.
- `Page.Truncated` on the lists the server caps.
- `OAuthService`: `ClientCredentials`, `ExchangeCode`, `RefreshToken`,
  `ExchangeToken`, `AuthorizeDevice`, `PollDeviceToken`,
  `AuthorizeBackchannel`, `PollBackchannelToken`, `Introspect`, `Revoke`,
  `UserInfo` and `Discover`; `Token` exposes `expires_in` and `Expiry()`.
- Package `webhook`: `Verifier` with timestamp tolerance, secret rotation,
  constant-time comparison, typed events and a `net/http` handler.
- Package `jwt`: ES256 verification against the JWKS with key caching,
  rotation handling and `net/http` middleware.
- About 140 operations the previous draft did not cover, among them ABAC
  policies, AuthZEN batch evaluation and decision logs, role holders and
  permissions, user sessions, devices, passkeys and security posture,
  suspended-user review, group templates and import, domain settings and
  break-glass members, OIDC connections, OAuth resources, SAML metadata,
  certificates, assignments and group mappings, directory sync, SCIM tokens,
  webhook updates, secret rotation, delivery attempts and retries, API key,
  agent and service principal secret rotation, compliance evidence, email
  delivery reports and billing profile, trial and invoice payment.
- Runnable examples and four example programs.

### Fixed

- Requests and models now match the API contract, among them: `Users.Create`
  sends the required `tenant_id`; `Tenants.Create` returns the new tenant
  with its owner's token pair; `Tenants.Delete` reports the scheduled
  deletion; the security policy's durations are strings such as "720h", not
  nanoseconds; OIDC client creation and secret rotation return the secret
  the way the server sends it; SAML connections, auth policy, branding,
  webhooks, agents, secrets and billing types carry the server's fields.
- OAuth token endpoint errors are parsed from the RFC 6749 shape instead of
  being reported as `http_400`.
- Introspection authenticates the calling application, which the server
  requires.

### Removed

- `New`, `Config`, `Error`, `Logger`, `ListParams` and the per-resource page
  types, replaced as described above.
- `Discover` and `NewFromDiscovery`; use `Client.OAuth.Discover`, and pass
  its `JWKSURI` to `jwt.WithJWKSURL`.
- `VerifyWebhookSignature`, `ConstructEvent`, `ConstructEventFromRequest`,
  `WebhookEvent` and the webhook header constants; use package `webhook`.
- `Client.Sessions`, `Claims` and `VerifyOptions`; use package `jwt`.

The module still has no dependencies outside the standard library.

## [0.1.0] — 2026-07-10 (draft, never tagged)

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

[Unreleased]: https://github.com/qeetgroup/qeet-id-go/commits/main
