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
| [AGENTS.md](AGENTS.md) | **Rules, commands, what CI runs** |
| [docs/llm/context.md](docs/llm/context.md) | How this repository works |
| [docs/llm/boundaries.md](docs/llm/boundaries.md) | What it owns, and what it must not touch |
| [docs/llm/workflows.md](docs/llm/workflows.md) | Step-by-step for common changes |
| [docs/llm/architecture-map.md](docs/llm/architecture-map.md) | "Where is X?" — fastest path to a file |

Parent context: **L0** `qeetgroup/qeet-context` · **L1** `qeetgroup/qeet-id-context`.

## What this repository is

The server-side **Go SDK for Qeet ID**. Module `github.com/qeetgroup/qeet-id-go`, package `qeetid`,
with leaf packages `webhook` and `jwt`. Go 1.26+, standard library only. It wraps the contract
owned by `qeet-id-server` and never adds behaviour the API lacks.

## Non-negotiables

1. Standard library only; a dependency is the owner's decision.
2. Match the server's OpenAPI documents exactly; the handlers decide ambiguities.
3. Retries only where safe; never retry a 401 or `auth.session_revoked`.
4. Never weaken `webhook/` or `jwt/` verification.
5. `context.Context` first, `opts ...RequestOption` last.

## Commands

`make build` · `make test` · `make vet` · `make fmt` · `make lint` · `make vuln` · `make tidy` ·
`make ci`

## Before finishing

```bash
make ci
git diff
```
