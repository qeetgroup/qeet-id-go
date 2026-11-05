# Architecture Map — qeet-id-go

**Level:** L2 · **Last verified:** 2026-11-05
**Verification scope:** every path below was confirmed to exist.

| Need | Path |
|---|---|
| Repository identity | [`qeet-repo.yml`](../../qeet-repo.yml) |
| Agent instructions | [`AGENTS.md`](../../AGENTS.md) |
| Package documentation | `doc.go` |
| **Client, `NewClient`, all 48 services** | `client.go` |
| Client and request options, defaults | `options.go` |
| Credentials, `TokenSource`, client credentials cache | `auth.go` |
| **Request pipeline, retries, logging, idempotency keys** | `transport.go` |
| **`APIError`, `OAuthError`, sentinels, error codes** | `errors.go` |
| `Page[T]`, the `All` iterator helper | `pagination.go` |
| `Token` (OAuth token response) | `token.go` |
| OAuth flows, introspection, discovery | `oauth.go` |
| LDAP sign-in and its `TokenPair`/`MFAChallenge` result | `ldap_authenticate.go` |
| User import from another provider's export | `users_import.go` |
| Group update with "move to top" | `groups_update.go` |
| Version constant | `version.go`, `internal/version/version.go` |
| A service and its types | `<service>.go`, e.g. `users.go`, `webhooks.go`, `samlconnections.go` |

## Sub-packages

| Need | Path |
|---|---|
| **Webhook signature verification, handler** | `webhook/webhook.go` |
| Webhook event envelope, catalogue, data types | `webhook/events.go` |
| Recorded server deliveries for tests | `webhook/testdata/deliveries.json` |
| **Access token verification** | `jwt/jwt.go` |
| Claims | `jwt/claims.go` |
| **JWKS cache and refresh** | `jwt/jwks.go` |
| Middleware, `RequireScope` | `jwt/middleware.go` |

## Tests

| Need | Path |
|---|---|
| Route case for every generated method | `routes_cases_test.go` |
| Route cases for hand-written methods, coverage check | `routes_test.go` |
| Fakes (in-process transport) | `helpers_test.go` |
| Retry rules | `retry_test.go` |
| Errors | `errors_test.go` |
| OAuth | `oauth_test.go` |
| Pagination and individual resources | `resources_test.go` |
| Runnable examples | `example_test.go`, `webhook/example_test.go`, `jwt/example_test.go` |

## Docs and tooling

| Need | Path |
|---|---|
| Design decision in force | `docs/design-decisions/ADR-0004-sdk-design.md` |
| Architecture overview | `docs/architecture/ARCHITECTURE.md` |
| Example programs | `examples/` |
| Lint configuration | `.golangci.yml` |
| CI | `.github/workflows/` |
