# Boundaries — qeet-id-go

**Level:** L2 · **Last verified:** 2026-08-28
**Verification scope:** ownership derived from the repository tree and the L1 map in
`qeet-id-context/REPOSITORIES.md`. Cross-SDK differences verified by reading all three SDKs.

## Owns

Typed Go bindings for the Qeet ID REST API · HTTP transport with retry and backoff · cursor
pagination as `iter.Seq2` · error typing · JWKS/ES256 token verification · webhook HMAC verification.

## Does not own

| Not owned | Owner |
|---|---|
| **The API contract** | `qeet-id-server` — see its OpenAPI documents |
| Any authorization or tenancy decision | `qeet-id-server` |
| Token issuance, session semantics | `qeet-id-server` |
| The browser/React story | `qeet-id-react` |
| The Node story | `qeet-id-node` |
| Product documentation | `qeet-id-docs` |
| Organization standards | `qeet-context` (L0) |
| Product architecture | `qeet-id-context` (L1) |

**The SDK never adds behaviour the API lacks.** If a method needs an endpoint that does not exist,
the endpoint is the work — not a client-side simulation of it.

## Consumes

`api.id.qeet.in` over HTTPS, authenticated with `Authorization: ApiKey <qk_...>`; the JWKS document
at `/.well-known/jwks.json`; nothing else. No database, no filesystem, no third-party package.

## Provides

A Go module. **Consumers are external and unsurveyable once published** — every exported name is a
contract under Go's compatibility expectations.

## Security boundaries

| Boundary | Enforcement |
|---|---|
| API key → wire | `internal/transport/request.go` — never logged, never in an error |
| Token → claims | `internal/auth/jwks.go` — ES256, unknown `kid` rejected |
| Webhook body → event | `webhooks_verify.go` — constant-time HMAC compare |
| Retry → side effects | `internal/transport/retry.go` — non-idempotent methods never retried |

**This SDK is not a trust boundary.** It runs in the caller's process with the caller's API key. It
cannot enforce anything; it can only avoid weakening what the server enforces.

## Cross-SDK differences — read before "harmonising" anything

The three Qeet ID SDKs deliberately and accidentally differ. **Do not change one to match another
without a decision.**

| Concern | Go | Node | React |
|---|---|---|---|
| Error model | `(T, error)`, `*Error` | **throws**, 4 classes | throws `AuthenticationError` |
| Network failure | wrapped `fmt.Errorf` | `NetworkError` class | — |
| Webhook header const | `X-Qeet-Signature` | **lowercase** `x-qeet-signature` | n/a |
| `ConstructEventFromRequest` | **yes** | deliberately omitted | n/a |
| Webhook payload | deferred (`json.RawMessage`) | eagerly parsed | n/a |
| Pagination | `iter.Seq2` | async generator | **none** |
| Auth | `ApiKey` header | `ApiKey` header | **cookie + CSRF** |
| Base URL default | `https://api.id.qeet.in` | same | **none — required** |
| Client type name | `Client` | `QeetID` | `QeetIDClient` |
| Wire field casing | as-is | snake_case kept (ADR-0004) | **mapped to camelCase** |
| Config escape hatch | `HTTPClient` | `fetch` | neither |

The last row of the casing comparison is a genuine contradiction: Node's ADR-0004 keeps wire fields
snake_case; React maps them to camelCase. Same organization, opposite decisions. Both are recorded;
neither is a bug to fix here.

## Safe to change without coordination

Internal refactoring behind an unchanged public API · adding a test · a bug fix that does not alter a
signature · doc comments · adding a **new** method or resource file.

## Requires coordination

| Change | Why |
|---|---|
| **Any exported name** | Go compatibility — a rename breaks every consumer |
| Auth scheme | Matches the server and the Node SDK |
| Retry policy | A wrong change duplicates server-side writes |
| Webhook header or algorithm | Must match `qeet-id-server` and the Node SDK |
| Error type shape | Consumers switch on it |
| Go version floor | Affects who can compile |

Product-level fan-out: `qeet-id-context/CHANGE-MATRIX.md`.

## Hard limits

1. **Never add a third-party dependency.** Stdlib only; there is no `go.sum`.
2. **Never make POST/PATCH/PUT retryable.**
3. **Never log or embed an API key** — not in an error, not in a debug line.
4. **Never skip signature or `kid` verification**, including "just for tests".
5. **Never wrap an endpoint you have not verified in the backend handler.**
6. **Never create a subpackage for public API** — ADR-0001.
7. **Never change another repository** from a task scoped to this one.
