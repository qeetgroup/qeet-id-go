# ADR-0003: Shared HTTP/crypto machinery lives in `internal/`

**Status:** Accepted

## Context

Every resource needs the same request lifecycle: build a request, attach
auth, retry on 429/5xx, decode JSON or map an error. `Sessions.Verify` needs
JWKS caching and ES256 verification. Duplicating either across ~40 resource
files, or inlining it all in the root package, both make the actual
resource logic harder to find.

## Decision

- `internal/transport` owns HTTP execution (`Transport.Do`/`Get`/`Post`/
  `Patch`/`Put`/`Delete`/`DoForm`) and the concrete `Error` type.
- `internal/auth` owns JWKS caching + ES256 verification
  (`JWKSVerifier.Verify`).
- `internal/pagination` owns the generic cursor-walking iterator
  (`Paginate[T]`).
- `internal/marshal` owns the `{items:[]}`/`{data:[]}` envelope-unwrap
  helper.
- Root-package files (`sessions.go`, every resource file) are thin: they
  hold a `*transport.Transport` and call it; `errors.go` re-exports
  `transport.Error` as `qeetid.Error` via a type alias.

## Why

- **Testable in isolation.** `internal/transport`'s retry/backoff logic,
  `internal/auth`'s JWKS/ES256 verification, and `internal/pagination`'s
  iterator each have their own unit tests, independent of any specific
  resource — a bug in retry logic gets caught by one focused test file, not
  by chasing it through 40 resources that each reimplement it slightly
  differently.
- **`internal/` enforces the boundary for free.** Go's compiler refuses
  imports of `internal/` packages from outside this module — there's no way
  for a consumer to depend on transport internals even by accident, so the
  public surface stays exactly `Client` + its resource types.
- **One retry policy, one JSON-error mapping, one place to change either.**

## Consequence: the `Error` type alias

Because `internal/transport` is the package that actually parses HTTP
error responses, it has to own the concrete `Error` struct. The root
package can't redefine its own `Error` type and convert at the boundary
without either losing information or adding a conversion step to every
single resource method. A type alias (`type Error = transport.Error`)
avoids both — see [ARCHITECTURE.md](../architecture/ARCHITECTURE.md#the-error-type-alias).
