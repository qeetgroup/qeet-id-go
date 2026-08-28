# GitHub Copilot — qeet-id-go

**Canonical instructions: [`AGENTS.md`](../AGENTS.md).** This file is a summary; it adds no
architecture.

## Repository

The server-side **Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`, Go 1.23, **zero third-party dependencies**. It wraps a contract owned by `qeet-id-server` — it never adds behaviour the API lacks.

Context: **L0** `qeet-context` (organization) → **L1** `qeet-id-context` (product) → **L2** this
repository → source.

## Structure

Public API is a **single root package**, one file per resource (`users.go`, `agents.go`, …), with 44 flat `*Service` fields on `Client`. Shared machinery is in `internal/{transport,auth,constants,pagination,marshal,validation,version,testutil}`. **Never create a subpackage for public API** (ADR-0001) and **never introduce nested namespace groups** (ADR-0002).

## Rules

1. **Zero third-party dependencies.** Stdlib only — there is deliberately no `go.sum`. Adding one is a decision, not a commit.
2. **Never invent an endpoint.** Verify against the real handler in `qeet-id-server`, not just the OpenAPI document.
3. **Retry idempotency is a safety rule:** 429 always retries; 5xx only on GET/DELETE. Making POST/PATCH/PUT retryable duplicates server-side writes.
4. **Never weaken JWKS or webhook verification** (`internal/auth/jwks.go`, `webhooks_verify.go`), including to make a test pass.
5. **`context.Context` is always the first argument** on any method that performs I/O.
6. A transport failure is a wrapped `fmt.Errorf`, **not** an `*Error` — there is no `NetworkError` type here, unlike the Node SDK.

## Commands

`make build` · `make test` · `make vet` · `make fmt` · `make lint` · `make tidy`

## Do not

- suggest adding a third-party dependency
- suggest a subpackage for public API, or nested client namespaces
- make a non-idempotent HTTP method retryable
- log, embed, or echo an API key — not in an error, not in a debug line
- invent Make targets — `build test vet fmt lint tidy` is the complete set
