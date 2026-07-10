# Architecture

## One package, many files

`qeet-id-go` is a single Go package (`qeetid`) with ~40 resource files at
the repo root, each following the same shape:

```go
type UsersService struct{ t *transport.Transport }

func (s *UsersService) Create(ctx context.Context, in CreateUserInput) (*User, error) {
	var out User
	err := s.t.Post(ctx, "/v1/users", in, &out)
	return &out, err
}
```

Larger resources split into companion files (`users.go` / `users_models.go`
/ `users_requests.go` / `users_responses.go`); smaller ones stay in one
file. Every `*XService` is a field directly on `Client` (`client.go`),
grouped only by comment banners (Identity / Authentication / Authorization
/ Administration) — there's no nested namespace type. `client.Users`,
`client.Sessions`, `client.Webhooks` are all one field access away.

## Why flat, not package-per-resource

The alternative — `qeetid/users`, `qeetid/oauth`, ... as separate importable
packages (Clerk v2's model) — was tried and reverted. It solves nothing a
flat package doesn't already solve, and it introduces a real problem: the
root package would need to import every resource package to build the
`Client` aggregator, so no resource package could import the root package
back without an import cycle. That forces shared types (`Error`, `Config`)
into an `internal/` layer with public re-export aliases just to dodge the
cycle — real complexity paid for a namespacing benefit a comment banner
gets for free. See [ADR-0002](../design-decisions/ADR-0002-flat-vs-grouped-api.md).

## `internal/` layering

Even though the *public* API is one package, the *implementation* is split
into focused internal packages — each independently testable, none
importable from outside this module:

| Package | Owns |
|---|---|
| `internal/transport` | HTTP execution: auth header, retry/backoff, JSON (de)serialization, the concrete `Error` type, form-encoded requests (`DoForm`, for OAuth grants) |
| `internal/auth` | JWKS caching + ES256 signature verification — the crypto behind `Sessions.Verify` |
| `internal/pagination` | The generic `Paginate[T]` cursor-walking iterator every `All()` method calls |
| `internal/marshal` | The `{items:[]}` / `{data:[]}` envelope-unwrap helper used by every `List` method |
| `internal/constants` | Default base URL/timeout/retries, header names |
| `internal/validation` | Fail-fast required-field/UUID checks |
| `internal/version` | The SDK version string |
| `internal/testutil` | Shared `httptest` server + request recorder for every resource's tests |

Root package files import these freely (one direction only — root → internal,
never the reverse), so there's no cycle risk despite the layering.

## The `Error` type alias

`errors.go` at the root is one line:

```go
type Error = transport.Error
```

This is a type *alias*, not a new type — `qeetid.Error` and
`transport.Error` are the same type at compile time. Every resource method
returns whatever `internal/transport` constructs, and callers write
`errors.As(err, &qeetid.Error{})` without ever knowing an `internal/`
package exists.

## Authorization checks live on `Permissions`

`Permissions.Check` / `CheckAll` / `Explain` are the RBAC hot-path calls —
made on nearly every authenticated request. They're methods on
`PermissionsService`, not a bespoke root-level `Client.Can()` — a permission
check is a `Permissions` operation, and putting it there means there's
exactly one place to look for every permission-related capability instead
of a special case at the top level.

## Sessions doesn't use the API-key transport

`SessionsService` (JWKS verification) never sends an API key — it hits a
public JWKS endpoint. It's constructed separately in `New()`, wrapping
`internal/auth.JWKSVerifier` directly rather than going through
`internal/transport`.

## OAuth's two transports

`OAuthService` legitimately needs two different request shapes:
- RFC-standard grant/introspection endpoints (`TokenExchange`, `Introspect`,
  `Revoke`, device flow, CIBA) authenticate via OIDC client credentials
  (HTTP Basic, optional) over form-encoded bodies — `transport.DoForm`.
- Its `SigningKeys`/`Grants`/`Devices` sub-resources are ordinary
  ApiKey-authed JSON admin endpoints — the normal `transport.Get`/`Post`
  path.

Both live on the same `*transport.Transport`, so `OAuthService` doesn't need
two different HTTP clients — just two different methods on one.
