# AGENTS.md — qeet-id-go

**The model-neutral instruction file for coding agents.** [CLAUDE.md](CLAUDE.md),
[GEMINI.md](GEMINI.md) and [.github/copilot-instructions.md](.github/copilot-instructions.md) point
here and add nothing architectural.

## What this repository is

The **official server-side Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package
`qeetid`, v0.1.0, Go 1.23, **zero third-party dependencies** (there is no `go.sum`).

It is a **client of a contract it does not own.** The contract lives in `qeet-id-server`
(`qeet-id-server/api/openapi/`). This SDK wraps it — it never adds behaviour the API does not have.

## Context hierarchy

```text
qeet-context (L0)  →  qeet-id-context (L1)  →  qeet-id-server (the contract)
                                            →  qeet-id-go (L2 — this repository)
```

Read L0/L1 when a task needs organization or product understanding. This repository does not
restate them.

## Read before changing code

| File | For |
|---|---|
| [qeet-repo.yml](qeet-repo.yml) | Machine-readable identity |
| [docs/llm/architecture-map.md](docs/llm/architecture-map.md) | "Where is X?" |
| [docs/llm/context.md](docs/llm/context.md) | How this SDK works |
| [docs/llm/boundaries.md](docs/llm/boundaries.md) | What it owns; cross-SDK differences |
| [docs/llm/workflows.md](docs/llm/workflows.md) | Adding a resource, a method, a release |

Repository ADRs: `docs/design-decisions/` (three, all binding).

## Rules

### 1. Zero third-party dependencies
`CONTRIBUTING.md` mandates stdlib only. **Adding a dependency changes the SDK's contract with every
consumer** and needs a decision, not a commit. There is deliberately no `go.sum`.

### 2. One package, file per resource — ADR-0001
All public API is `package qeetid` at the repository root, one `.go` file per resource
(`users.go`, `agents.go`, …). Large resources split into `_models.go` / `_requests.go` /
`_responses.go`. **Do not create subpackages for public API.**

### 3. Flat client fields — ADR-0002
Services are flat fields on `Client` (`c.Users`, `c.Agents`), grouped only by comment banner.
**Do not introduce nested namespace groups.**

### 4. Shared machinery lives in `internal/` — ADR-0003
Transport, JWKS, pagination, validation and constants are internal. Public files must not
re-implement them.

### 5. `context.Context` is always the first argument
Every method that performs I/O takes `ctx context.Context` first.

### 6. Never invent an endpoint
`CONTRIBUTING.md`: verify against the **real backend handler**, not just the OpenAPI document.
A method wrapping a route that does not exist is worse than a missing method.

### 7. Never weaken security machinery
Do not change JWKS verification (`internal/auth/jwks.go`), webhook HMAC verification
(`webhooks_verify.go`), or auth-header injection (`internal/transport/request.go`) without review.
**Retry idempotency is a safety rule, not a tuning knob**: 429 always retries; 5xx retries only on
GET/DELETE. Never make POST/PATCH/PUT retryable.

### 8. Errors are values, not panics
Return `(T, error)`. Public errors are `*qeetid.Error` (a type alias for `transport.Error`).
**A transport failure is a wrapped `fmt.Errorf`, not an `*Error`** — Go has no `NetworkError` type,
unlike the Node SDK.

### 9. Tests
`internal/testutil` provides a real `httptest.Server`. Use it; do not hand-roll a mock.

## Commands

```bash
make build   # go build ./...
make test    # go test -race -count=1 ./...
make vet     # go vet ./...
make fmt     # FAILS if gofmt -l is non-empty
make lint    # golangci-lint run
make tidy    # go mod tidy
```

That is the complete target set. **Do not invent targets.**

## What CI enforces

`ci.yml` (matrix Go 1.23 + 1.24): `go vet`, `go build`, `go test -race` + coverage, and an inline
gofmt check. `lint.yml`: golangci-lint. `security.yml`: govulncheck + gitleaks, weekly.
`codeql.yml`: CodeQL with autobuild. `release.yml`: on tag `v*` — **the tag is the artifact**;
Go modules need no build step.

## Before you finish

```bash
make fmt && make vet && make lint && make test
git diff
```

## Distribution status — know this

**This SDK has never been fetched by the Go module proxy** and carries no release tag, while the
organization advertises "drop-in SDKs". See `qeet-id-context/DRIFT-REGISTER.md` QID-009. Do not
describe it as available to users.

## Escalate rather than proceed

Stop and report if a task would add a dependency, change the auth scheme, weaken signature
verification, make a non-idempotent method retryable, or require changing another repository.
