# GEMINI.md — qeet-id-go

**Read [AGENTS.md](AGENTS.md) first.** It is the model-neutral instruction file and the source of
truth for this repository. This file is a pointer and adds no architecture.

```text
GEMINI.md  →  AGENTS.md  →  docs/llm/*
```

## Canonical context

| File | For |
|---|---|
| [qeet-repo.yml](qeet-repo.yml) | Machine-readable repository identity |
| [AGENTS.md](AGENTS.md) | **Rules, commands, what CI enforces** |
| [docs/llm/context.md](docs/llm/context.md) | How this repository actually works |
| [docs/llm/boundaries.md](docs/llm/boundaries.md) | What it owns, and what it must not touch |
| [docs/llm/workflows.md](docs/llm/workflows.md) | Step-by-step for common changes |
| [docs/llm/architecture-map.md](docs/llm/architecture-map.md) | "Where is X?" — fastest path to a file |

Parent context: **L0** `qeetgroup/qeet-context` · **L1** `qeetgroup/qeet-id-context`.

## What this repository is

The server-side **Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`, Go 1.23, **zero third-party dependencies**. It wraps a contract owned by `qeet-id-server` — it never adds behaviour the API lacks.

## Non-negotiables

1. **Zero third-party dependencies.** Stdlib only — there is deliberately no `go.sum`. Adding one is a decision, not a commit.
2. **Never invent an endpoint.** Verify against the real handler in `qeet-id-server`, not just the OpenAPI document.
3. **Retry idempotency is a safety rule:** 429 always retries; 5xx only on GET/DELETE. Making POST/PATCH/PUT retryable duplicates server-side writes.
4. **Never weaken JWKS or webhook verification** (`internal/auth/jwks.go`, `webhooks_verify.go`), including to make a test pass.
5. **`context.Context` is always the first argument** on any method that performs I/O.
6. A transport failure is a wrapped `fmt.Errorf`, **not** an `*Error` — there is no `NetworkError` type here, unlike the Node SDK.

## Commands

`make build` · `make test` · `make vet` · `make fmt` · `make lint` · `make tidy`

That list is complete — **do not invent commands.**

## Before finishing

```bash
make fmt && make vet && make lint && make test
git diff
```
