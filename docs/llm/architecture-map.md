# Architecture Map — qeet-id-go

**Level:** L2 · **Last verified:** 2026-08-28
**Verification scope:** every path below was confirmed to exist.

| Need | Path |
|---|---|
| Repository identity | [`qeet-repo.yml`](../../qeet-repo.yml) |
| Agent instructions | [`AGENTS.md`](../../AGENTS.md) |
| **Client constructor + all 44 services** | `client.go` |
| Client options | `config.go` |
| Public error type (alias) | `errors.go` |
| Version constant | `version.go` |
| Discovery helper | `discovery.go` |
| **Webhook HMAC verification** | `webhooks_verify.go` |
| Pagination docs (no code) | `pagination.go` |

## Internal machinery

| Concern | Path |
|---|---|
| **HTTP transport, auth header injection** | `internal/transport/request.go` |
| Transport construction, options | `internal/transport/transport.go` |
| **Error type, `parseErrorBody`, `readResponse`** | `internal/transport/response.go` |
| Retry policy, backoff, `Retry-After` | `internal/transport/retry.go` |
| OAuth form posts (Basic auth) | `internal/transport/form.go` |
| Logger interface | `internal/transport/middleware.go` |
| User-Agent | `internal/transport/useragent.go` |
| **JWKS cache + ES256 verification** | `internal/auth/jwks.go` |
| **All defaults, header names, TTLs** | `internal/constants/constants.go` |
| Generic `Paginate[T] -> iter.Seq2` | `internal/pagination/` |
| `Envelope[T]` items/data resolution | `internal/marshal/` |
| Fail-fast argument validation | `internal/validation/` |
| Version (isolated to avoid a cycle) | `internal/version/version.go` |
| Test server + recorder | `internal/testutil/testutil.go` |

## Resources by category

Flat files in the root package, grouped by the comment banners in `client.go`:

```text
Identity        users.go  organizations.go  serviceprincipals.go  agents.go  domains.go
Authentication  sessions.go  oauth.go  oidc.go  saml.go  saml_providers.go  scim.go
                ldap.go  social.go  mfa.go  credentials.go  authhooks.go
                authpolicy.go  policy.go  iprules.go  botdetection.go  risksettings.go
Authorization   roles.go  permissions.go  groups.go  relationtuples.go  authzen.go
Administration  branding.go  invitations.go  emailtemplates.go  apikeys.go  vault.go
                tokenvault.go  webhooks.go  auditlogs.go  analytics.go  gdpr.go
                billing.go  retention.go  ratelimits.go  logsinks.go  adminlinks.go
```

Large resources add `_models.go` / `_requests.go` / `_responses.go` — see `users.go`,
`organizations.go`, `sessions.go`.

## Decisions

| ADR | Decision |
|---|---|
| `docs/design-decisions/ADR-0001-single-package.md` | Single package, file per resource |
| `docs/design-decisions/ADR-0002-flat-vs-grouped-api.md` | Flat `Client` fields, not nested groups |
| `docs/design-decisions/ADR-0003-internal-transport-layering.md` | Shared machinery in `internal/` |

Also `docs/architecture/ARCHITECTURE.md`, `docs/FUTURE.md`, `docs/release-notes/v0.1.0.md`.

## Build, test, CI

| | |
|---|---|
| Commands | `Makefile` |
| Tests | 6 `_test.go` files, stdlib `testing` + `httptest` |
| CI | `.github/workflows/{ci,lint,security,codeql,release}.yml` |
| Contribution rules | `CONTRIBUTING.md` |
| Disclosure policy | `SECURITY.md` |

## Upstream

The API contract this SDK implements lives in **`qeet-id-server`**, under its OpenAPI documents — five
OpenAPI 3.1 documents. It is not vendored here.
