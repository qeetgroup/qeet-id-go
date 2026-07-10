# Contributing

## Setup

```bash
git clone https://github.com/qeetgroup/qeet-id-go
cd qeet-id-go
make build test
```

Requires Go 1.23+ (the SDK uses `iter.Seq2` range-over-func iterators).

## Conventions

- **One package.** Every resource is a file at the repo root (`users.go`,
  `oauth.go`, ...); larger resources split into `_models.go`/`_requests.go`/
  `_responses.go` companions. See
  [docs/architecture/ARCHITECTURE.md](docs/architecture/ARCHITECTURE.md) for
  why this beats a package-per-resource split at this SDK's size.
- **Zero third-party dependencies.** Standard library only, in the public
  module. This is a deliberate constraint, not an oversight — see
  [docs/design-decisions/](docs/design-decisions/) before proposing a new
  dependency.
- **Every resource method takes `context.Context` first.**
- **Shared HTTP/crypto machinery lives in `internal/`** — never duplicate
  retry/backoff, JSON envelope-unwrapping, or JWKS logic in a resource file;
  add to `internal/transport`, `internal/marshal`, or `internal/auth`
  instead and call from there.
- **Match the backend, not the OpenAPI spec's `additionalProperties: true`
  placeholders.** Several backend routes have unenriched spec bodies —
  verify the real Go handler/domain struct in the `qeet-id` backend repo
  before adding a typed response.

## Tests

Every resource file has a companion `_test.go` asserting method/path/body
shape against an `httptest.Server`, using the shared helper in
`internal/testutil`. Run:

```bash
make test          # go test -race ./...
```

## Adding a new resource

1. Add the file(s) at the repo root, following an existing resource (e.g.
   `domains.go`) as a template: a `*XService struct{ t *transport.Transport }`,
   methods on it, request/response types.
2. Add the `*XService` field to `Client` in `client.go`, under the right
   comment banner (Identity / Authentication / Authorization /
   Administration), and construct it in `New()`.
3. Add a `_test.go` with method/path/body-shape assertions.
4. Update the resource table in `README.md`.

## Reporting bugs

Open an issue with the exact method call, expected vs. actual behavior, and
(if applicable) the `RequestID` from the returned `*qeetid.Error`.
