# AGENTS.md — qeet-id-go

**The model-neutral instruction file for coding agents.** [CLAUDE.md](CLAUDE.md),
[GEMINI.md](GEMINI.md) and [.github/copilot-instructions.md](.github/copilot-instructions.md) point
here and add nothing architectural.

## What this repository is

The **server-side Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`
(imported as `qeetid "github.com/qeetgroup/qeet-id-go"`), plus two leaf packages: `webhook`
(receiving deliveries) and `jwt` (verifying access tokens). Version string 0.1.0, **not yet
tagged**. Go 1.26+ (the two supported Go releases). **No dependencies outside the standard
library** (there is no `go.sum`).

It is a **client of a contract it does not own.** The contract is the five OpenAPI documents in
`qeet-id-server/api/openapi/` (`auth`, `management`, `federation`, `developer`, `operations`); the
server's handlers under `internal/` decide where the documents are ambiguous. The SDK never adds
behaviour the API does not have.

## Context hierarchy

```text
qeet-context (L0)  →  qeet-id-context (L1)  →  qeet-id-server (the contract)
                                            →  qeet-id-go (L2 — this repository)
```

## Read before changing code

| File | For |
|---|---|
| [qeet-repo.yml](qeet-repo.yml) | Machine-readable identity |
| [docs/llm/architecture-map.md](docs/llm/architecture-map.md) | "Where is X?" |
| [docs/llm/context.md](docs/llm/context.md) | How this SDK works |
| [docs/llm/boundaries.md](docs/llm/boundaries.md) | What it owns and must not touch |
| [docs/llm/workflows.md](docs/llm/workflows.md) | Adding a method, a resource, a release |
| [docs/design-decisions/](docs/design-decisions/) | ADR-0004 is the current design; 0001–0003 are history |

## Rules

### 1. Standard library only
Adding a dependency changes the SDK's contract with every consumer and needs the owner's decision.
`golang.org/x/oauth2` interoperates through `TokenSource` without a dependency.

### 2. Layout — ADR-0004
All API resources live in the root package `qeetid`, one file per service (`users.go`,
`webhooks.go`, …) holding the service, its methods and its types. Transport, retries, errors and
pagination are unexported code in the root package (`transport.go`, `errors.go`,
`pagination.go`, `auth.go`, `options.go`). `webhook/` and `jwt/` are the only sub-packages: they
serve programs that need no API client. `internal/version` holds the version string. Do not add
other public packages without a decision.

### 3. Shape of a method
`func (s *XService) Verb(ctx context.Context, <path args>, params XParams, opts ...RequestOption)`.
Context first; path arguments in the order the path names them; a request body as a value
`...Params` struct; optional query parameters as a pointer struct (`nil` sends none); `opts` last.
Lists return `*Page[T]`; cursor lists also get `All…` returning `iter.Seq2[*T, error]`. Every
method's doc comment ends with `API: VERB /path`.

### 4. Match the contract exactly
Field names, types and enums come from the OpenAPI documents. When a document is ambiguous, read
the handler in `qeet-id-server/internal/` and follow the code. **Never wrap an endpoint that does
not exist.** In request bodies an empty string is left out (`omitempty`): the server reads an
absent field as empty, while `""` can fail to decode a UUID field.

### 5. Retries are a safety rule, not a tuning knob
`retryPlan` in `transport.go`: network errors and 408/500/502/503/504 are retried only for
GET/HEAD/DELETE or a request carrying an `Idempotency-Key`; 429 or 503 **with** `Retry-After` is
retried for any method (the server refused before acting); 409 `idempotency.in_progress` only with
a key; **401 and `auth.session_revoked` never**. Do not widen this.

### 6. Never weaken verification or leak secrets
`webhook/webhook.go` (HMAC over `timestamp.body`, constant time, tolerance), `jwt/` (ES256 only,
`kid` required, issuer and audience required) and `auth.go` (credentials) are security code.
Credentials never appear in logs or error text. Keep the base URL https-only outside loopback.

### 7. Errors are values
`*APIError` (error envelope) and `*OAuthError` (RFC 6749 shape) match the sentinels in
`errors.go` through their `Is` methods. Network failures are wrapped with `%w`.

### 8. Tests
Every method has a route case (`routes_cases_test.go` for generated methods, `handRouteCases` in
`routes_test.go` for hand-written ones); `TestRoutesCoverEveryService` fails otherwise. Use the
in-process fakes in `helpers_test.go`; tests that wait on `Retry-After` run inside
`testing/synctest`. Webhook tests replay deliveries recorded from the server's signing code
(`webhook/testdata/deliveries.json`).

### 9. No ticket references
Name the behaviour, never a tracker key, in code, comments, tests and docs.

## Commands

```bash
make build   # go build ./...
make test    # go test -race -count=1 ./...
make vet     # go vet ./...
make fmt     # fails if gofmt -l is not empty
make lint    # golangci-lint run ./...  (v2.15.0, config in .golangci.yml)
make vuln    # govulncheck v1.9.0
make tidy    # go mod tidy
make ci      # fmt, vet, build, test, lint, vuln — what the workflows run
```

## What CI runs

`ci.yml` (Go 1.26.x and 1.27.x): gofmt, `go vet`, `go build`, `go test -race` with coverage.
`lint.yml`: golangci-lint v2.15.0. `security.yml`: govulncheck v1.9.0 and gitleaks, weekly too.
`codeql.yml`: CodeQL. `release.yml`: run by hand on `main` (patch, minor or major), runs the gate
and **pushes the next tag**, which publishes a module version. Merging into `main` tags nothing.

## Branches

Two long-lived branches: `develop`, where work lands, and `main`, what is released. Changes reach
`main` only by a pull request from `develop`. There are no release branches: a version is a
`vX.Y.Z` tag, pushed when `release.yml` is run on `main`.

## Before you finish

```bash
make ci
git diff
```

## Escalate rather than proceed

Stop and report if a task would add a dependency, change the module path or package name, weaken
signature or token verification, widen the retry rules, or require changing another repository.
