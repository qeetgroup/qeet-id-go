# CLAUDE.md — qeet-id-go

**Read [AGENTS.md](AGENTS.md) first.** It is the model-neutral instruction file and the source of
truth for this repository. This file adds only Claude-specific guidance.

```text
CLAUDE.md  →  AGENTS.md  →  docs/llm/*
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

## The things most likely to trip you up here

1. **Standard library only.** No `go.sum`; a new dependency is the owner's decision.
2. **The contract is the server's OpenAPI documents**, with the handlers as tiebreaker. Never wrap
   an endpoint that does not exist; copy field names and types exactly.
3. **Retry rules are safety rules** (`retryPlan` in `transport.go`): never retry a 401 or
   `auth.session_revoked`; retry writes only with an idempotency key or a 429/503 with
   `Retry-After`.
4. **Never weaken `webhook/` or `jwt/` verification**, including to make a test pass.
5. **`context.Context` first, `opts ...RequestOption` last** on every method.
6. **Every method needs a route case**, or `TestRoutesCoverEveryService` fails.

## Working style

- Read before editing. Match the neighbouring service file's shape.
- Use the architecture map instead of guessing a path.
- Do not read `.env*` or secret files into anything you write.

## Finishing a change

```bash
make ci
git diff
```

## Escalate rather than proceed

Stop and report if a task would weaken a security control, change the module path or package name,
add a dependency, or require modifying another repository.
