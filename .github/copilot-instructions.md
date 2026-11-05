# GitHub Copilot — qeet-id-go

**Canonical instructions: [`AGENTS.md`](../AGENTS.md).** This file is a summary; it adds no
architecture.

## Repository

The server-side **Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`,
Go 1.26+, standard library only. It wraps the contract owned by `qeet-id-server` (its five OpenAPI
documents) and never adds behaviour the API lacks.

Context: **L0** `qeet-context` (organization) → **L1** `qeet-id-context` (product) → **L2** this
repository → source.

## Structure

The root package holds the client (`client.go`, `options.go`, `auth.go`, `transport.go`,
`errors.go`, `pagination.go`) and one file per service (`users.go`, `webhooks.go`, …) with its
methods and types. Leaf packages: `webhook/` (verify deliveries) and `jwt/` (verify access tokens).
`internal/version` holds the version. See ADR-0004 in `docs/design-decisions/`.

## Rules

1. Standard library only.
2. Match the OpenAPI documents exactly; the server's handlers decide ambiguities.
3. `context.Context` first, `opts ...RequestOption` last; lists return `*Page[T]`.
4. Retries only where safe (`retryPlan`); never retry a 401 or `auth.session_revoked`.
5. Never weaken `webhook/` or `jwt/` verification; never log or embed a credential.
6. Every method has a route test case.

## Commands

`make build` · `make test` · `make vet` · `make fmt` · `make lint` · `make vuln` · `make tidy` ·
`make ci`

## Do not

- suggest adding a third-party dependency
- add a public package other than `webhook` and `jwt` without a decision
- widen the retry rules or retry a 401
- log, embed or echo a credential, in an error or a log line
- reference tracker tickets in code, comments or docs
