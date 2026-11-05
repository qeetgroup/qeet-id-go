# Boundaries — qeet-id-go

**Level:** L2 · **Last verified:** 2026-11-05
**Verification scope:** ownership derived from the repository tree and the server's OpenAPI
documents.

## Owns

Typed Go bindings for the Qeet ID REST API · the HTTP pipeline (credentials, retries, timeouts,
logging, idempotency keys) · error typing · pagination · webhook delivery verification (`webhook`)
· access token verification (`jwt`).

## Does not own

| Not owned | Owner |
|---|---|
| **The API contract** | `qeet-id-server` — its five OpenAPI documents and handlers |
| Any authorization or tenancy decision | `qeet-id-server` |
| Token issuance, session semantics, the webhook signature | `qeet-id-server` |
| End-user flows in browsers and apps | the browser and mobile SDKs |
| The Node SDK | `qeet-id-node` |
| Product documentation | `qeet-id-docs` |

**The SDK never adds behaviour the API lacks.** Client-side helpers are limited to protocol
mechanics: retries, pagination, verification.

## Consumes

The Qeet ID API over HTTPS (`https://api.id.qeet.in` by default), its JWKS document, and nothing
else. No database, no filesystem, no third-party module.

## Provides

A Go module with three public packages. Once tagged, **every exported name is a contract** under
Go's compatibility rules.

## Leaves out on purpose

End-user self-service endpoints (`/v1/me`, `/v1/auth/*`, passkey registration, MFA enrolment,
account export and deletion), browser redirect and form-post flows, SCIM protocol endpoints, and
console features (activity stream, search, notifications, the Qeet AI assistant, sales leads,
pre-signup checkout). The README lists them.

## Security boundaries

| Boundary | Enforcement |
|---|---|
| Credential → wire | `auth.go`, `transport.go` — https only outside loopback; never logged or in errors |
| Retry → side effects | `retryPlan` in `transport.go` |
| Webhook body → event | `webhook/webhook.go` — HMAC over timestamp and body, constant time, tolerance |
| Token → claims | `jwt/` — ES256 only, `kid`, issuer and audience required |

**This SDK is not a trust boundary.** It runs in the caller's process with the caller's
credentials. It can only avoid weakening what the server enforces.

## Cross-SDK differences

The Node and Go SDKs are reworked separately. Do not change one to match the other without a
decision; record differences here.

| Concern | Go |
|---|---|
| Errors | values: `*APIError`, `*OAuthError`, sentinels with `errors.Is` |
| Pagination | `Page[T]` and `iter.Seq2` iterators |
| Webhooks | package `webhook`, `Verifier` with `Handler` |
| Token verification | package `jwt`, `Verifier` with `Middleware` |
| Wire field names | Go names with the contract's JSON tags |

## Requires coordination

| Change | Why |
|---|---|
| Any exported name, once tagged | Go compatibility |
| Module path or package name | the import line every customer writes |
| Retry policy | a wrong change duplicates server-side writes |
| Webhook signature or JWT verification | must match `qeet-id-server` |
| Go version floor | decides who can compile |
