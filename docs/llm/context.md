# Repository Context — qeet-id-go

**Level:** L2 · **Status:** active · **Evidence state:** verified · **Last verified:** 2026-08-28
**Verification scope:** package layout, constants, transport behaviour, CI jobs and test setup read
from source. Distribution status checked against the Go module proxy.

## Identity

The server-side **Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`,
v0.1.0, Go 1.23 (needs `iter.Seq2`), **stdlib only — no third-party dependencies and no `go.sum`**.

**Status: `development`.** Published as a Go module in the sense that a tag would resolve — but the
module has **never been fetched by the Go module proxy** and the repository carries no release tag.
See `qeet-id-context/DRIFT-REGISTER.md` QID-009.

## Context inheritance

```text
qeet-context (L0)      organization standards
      ↓
qeet-id-context (L1)   product architecture, cross-repo contracts
      ↓
qeet-id-server         the API contract this SDK implements
      ↓
qeet-id-go (L2)        THIS REPOSITORY — how the Go client is built
```

## Responsibilities

Wrapping the Qeet ID REST API for Go callers: typed request/response models, HTTP transport with
retry, cursor pagination, error typing, JWKS-based token verification, and webhook signature
verification.

## Non-responsibilities

**It owns no behaviour.** Not authorization decisions, not tenant scoping, not token issuance, not
validation semantics — all of that is `qeet-id-server`. An SDK that decided something the API did
not would be a security bug.

It also does not own: the OpenAPI contract, the browser story (`qeet-id-react`), or docs
(`qeet-id-docs`).

## Architecture

```text
caller
  ↓  qeetid.New(Config) *Client
Client            44 flat *Service fields, grouped by comment banner only
  ↓
internal/transport   URL build → auth header → retry → response parse → typed error
  ↓
api.id.qeet.in
```

Three binding ADRs shape this: **single package, file per resource** (ADR-0001), **flat client
fields** (ADR-0002), **shared machinery in `internal/`** (ADR-0003).

Two services are constructed differently from the rest —
`Sessions: newSessionsService(t.BaseURL(), t.HTTPClient())` and `OAuth: newOAuthService(t)`;
everything else is `&XService{t: t}`.

## Transport

**Evidence:** `internal/transport/`, `internal/constants/constants.go`

| Property | Value |
|---|---|
| Base URL default | `https://api.id.qeet.in` |
| Auth header | `Authorization: ApiKey <qk_...>` |
| OAuth form endpoints | `Basic <base64>` — **never retried** |
| Timeout | 10s (ignored entirely if `Config.HTTPClient` is supplied) |
| Retries | 2 |
| Retry policy | **429 always; 5xx only when idempotent (GET/DELETE)** |
| Backoff | `250ms * 2^attempt + rand(0..100ms)`, honours `Retry-After` |
| Response cap | 1 MiB |
| Correlation | `X-Request-Id` |

**The retry policy is a correctness rule.** Making POST/PATCH/PUT retryable would duplicate
side effects on the server.

## Errors

`(T, error)` returns. `qeetid.Error` is a **type alias** for `transport.Error`, carrying
`Status`, `Code`, `Message`, `RequestID`, `RetryAfterSeconds`, with `IsUnauthorized()`-style
predicates.

> **A transport failure is a wrapped `fmt.Errorf`, not an `*Error`.** There is no `NetworkError`
> type. The Node SDK splits this into four classes; Go does not. Callers must not assume every
> failure is an `*Error`.

## Pagination

`internal/pagination.Paginate[T](ctx, startCursor, fetchPage) iter.Seq2[T, error]` — a Go 1.23
iterator with a `next == cursor` loop guard. Per-resource `All(ctx, ...)` methods wrap it. Errors
are yielded as `(zero, err)`.

## Cryptography

`internal/auth/jwks.go` — JWKS cache (**5 min TTL, 1 min refresh cooldown**), ES256 verification,
30s default clock skew. Deliberately independent of transport, because verification needs no API
key. A `sync.Mutex` guards refresh; unlike Node there is **no in-flight de-duplication**.

`webhooks_verify.go` — HMAC-SHA256 over the raw body. Headers `X-Qeet-Signature` and `X-Qeet-Event`
(canonical casing). `ConstructEventFromRequest` exists **only in Go**.

## Testing

6 `_test.go` files, ~43 test functions, stdlib `testing` with a real `httptest.Server` from
`internal/testutil`. Coverage is uneven: `CONTRIBUTING.md` claims "every resource file has a
companion `_test.go`" — **it does not** (6 test files against ~44 resources). Treat that claim as
aspirational.

## CI/CD

| Workflow | Runs |
|---|---|
| `ci.yml` | matrix Go 1.23/1.24 — `go vet`, `go build`, `go test -race` + coverage, inline gofmt check |
| `lint.yml` | `golangci-lint` (version `latest`) |
| `security.yml` | `govulncheck` + `gitleaks`, weekly cron |
| `codeql.yml` | CodeQL, autobuild, weekly cron |
| `release.yml` | on tag `v*` — vet, build, test, then a GitHub release. **The tag is the artifact** |

## Security-critical areas

| Area | Path | Risk | Review |
|---|---|---|---|
| JWKS / ES256 verification | `internal/auth/jwks.go` | **Critical** | Security review |
| Webhook HMAC | `webhooks_verify.go` | **Critical** | Security review |
| Auth header injection | `internal/transport/request.go` | High | Security review |
| Retry idempotency | `internal/transport/retry.go` | High | Review — a wrong change duplicates writes |

`SECURITY.md`: disclose privately to **security@qeet.in**, never a public issue; 2-business-day ack.

## Known constraints

- **Zero third-party dependencies is a hard rule**, not a preference.
- Go 1.23 minimum, because the pagination API returns `iter.Seq2`.
- `Config.HTTPClient` **overrides `Timeout` entirely** — a caller supplying a client owns its timeout.
- Test coverage is thin relative to the resource count.
- Not published; no consumers to break yet, which is the one thing that makes contract changes cheap
  right now.

## Documentation authority

Source > tests > `CONTRIBUTING.md` > README. The **API contract** is `qeet-id-server/api/openapi/`,
not anything in this repository.
