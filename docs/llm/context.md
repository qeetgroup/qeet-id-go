# Repository Context — qeet-id-go

**Level:** L2 · **Status:** active · **Evidence state:** verified · **Last verified:** 2026-11-05
**Verification scope:** package layout, defaults, transport behaviour, CI jobs and tests read from
source; endpoints checked against `qeet-id-server` develop.

## Identity

The server-side **Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`
(imported as `qeetid "github.com/qeetgroup/qeet-id-go"`), with leaf packages `webhook` and `jwt`.
Version string 0.1.0, **no tag pushed yet**. Go 1.26+, standard library only.

## Context inheritance

```text
qeet-context (L0)      organization standards
      ↓
qeet-id-context (L1)   product architecture, cross-repo contracts
      ↓
qeet-id-server         the API contract this SDK implements
      ↓
qeet-id-go (L2)        THIS REPOSITORY
```

## Responsibilities

Typed bindings for 48 services (about 330 operations) of the Qeet ID API, credentials (API key,
service principal client credentials, any bearer token source), safe retries, idempotency keys,
pagination, typed errors, webhook verification and access-token verification.

## Non-responsibilities

**It owns no behaviour.** Authorization, tenancy, token issuance and validation rules are the
server's. It does not own the OpenAPI contract, end-user flows, or documentation sites.

## Architecture

```text
qeetid.NewClient(opts...) *Client
Client.<Service>.<Method>(ctx, path args, params, opts...)
  ↓  call[T] → Client.do(ctx, *request, out, opts)          transport.go
     options → path → idempotency key → body encoded once
     → attempt (timeout, credential, User-Agent) → decode or *APIError/*OAuthError
     → retryPlan → backoff / Retry-After
  ↓
api.id.qeet.in
```

Design record: [ADR-0004](../design-decisions/ADR-0004-sdk-design.md).

## Transport

| Property | Value |
|---|---|
| Base URL default | `https://api.id.qeet.in` (https required outside loopback) |
| Credentials | `ApiKey <key>`, or `Bearer <token>` (token source or client credentials) |
| OAuth endpoints | HTTP Basic with the application's client ID and secret; public clients send `client_id` |
| Timeout | 30 s per attempt |
| Retries | 2; backoff 500 ms × 2ⁿ with jitter, capped at 8 s |
| Retry rule | see `retryPlan`: never 401 or `auth.session_revoked` |
| `Retry-After` | honoured up to 60 s; longer is returned in `APIError.RetryAfter` |
| Idempotency | random UUID on the nine supporting create calls; `WithIdempotencyKey` to choose |
| Response cap | 32 MiB |
| Correlation | `X-Request-Id`, or the envelope's `request_id` |
| Logging | `log/slog`, debug per attempt, warn per retry; no credentials, headers or bodies |

## Errors

`*APIError` carries the envelope `{error: {code, message, detail, translation_key, retryable,
fields, metadata, request_id}}` plus status and `Retry-After`. `*OAuthError` carries the RFC 6749
`error`/`error_description`. Both match `ErrInvalidRequest`, `ErrUnauthorized`,
`ErrSessionRevoked`, `ErrPaymentRequired`, `ErrForbidden`, `ErrNotFound`, `ErrConflict`,
`ErrRateLimited` or `ErrUnavailable` as their status and code fit. Network failures wrap the cause.

## Pagination

`Page[T]{Items, NextCursor, Truncated, Total, TotalCapped}`. Cursor lists have `All…` methods
returning `iter.Seq2[*T, error]`; capped lists report `Truncated`.

## Cryptography

`webhook`: HMAC-SHA256 of `timestamp + "." + body`, `hmac.Equal`, five-minute tolerance, any of
several signatures and secrets. `jwt`: ES256 with raw r‖s signatures converted to ASN.1, keys from
the JWKS cached ten minutes, unknown `kid` refresh at most once a minute, shared fetches, stale keys
during outages.

## Testing

In-process fakes (`helpers_test.go`), a route case per method, retry tests under
`testing/synctest`, webhook fixtures recorded from the server's signing code, generated ES256 keys
for `jwt`, runnable examples.

## CI/CD

| Workflow | Runs |
|---|---|
| `ci.yml` | Go 1.26.x and 1.27.x: gofmt, vet, build, `go test -race` with coverage |
| `lint.yml` | golangci-lint v2.15.0 |
| `security.yml` | govulncheck v1.9.0, gitleaks; weekly |
| `codeql.yml` | CodeQL |
| `release.yml` | on push to `main`: gate, then the next patch tag and a GitHub release |

## Documentation authority

Source > tests > ADR-0004 > README. The **API contract** is `qeet-id-server/api/openapi/`.
