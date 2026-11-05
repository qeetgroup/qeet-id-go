# Workflows — qeet-id-go

**Level:** L2 · **Last verified:** 2026-11-05
**Verification scope:** every command checked against the `Makefile` and `.github/workflows/`.

## Set up

```bash
git clone https://github.com/qeetgroup/qeet-id-go && cd qeet-id-go
make build
make test
```

No services are needed: tests answer requests in process. There is no `.env`.

## Add a method to an existing service

```text
1  find the operation           qeet-id-server/api/openapi/<doc>.yaml; read its handler when the
                                document is loose (untyped body, oneOf, "required" from struct tags)
2  open <service>.go            match the neighbouring methods exactly
3  signature                    ctx, path args in path order, params XParams (body),
                                params/options *XParams (query), opts ...RequestOption
4  body                         return call[T](ctx, s.c, &request{method, path, pathArgs, query,
                                body, idempotent, public}, opts)
5  types                        field names and JSON tags exactly as the contract; time.Time for
                                date-time; typed string constants for enums; omitempty on request
                                strings
6  doc comment                  "<Name> <does what>." … ends with "API: VERB /path"
7  route case                   routes_cases_test.go (or handRouteCases in routes_test.go)
8  make ci
```

**Never wrap an endpoint you have not confirmed exists.**

## Add a service

```text
1  new file <service>.go in the root package
2  type XService struct{ c *Client } with a doc comment
3  add the field to Client in client.go, in the right group, and wire it in NewClient
4  route cases for every method (TestRoutesCoverEveryService enforces it)
5  README "Services" table
6  make ci
```

## Add a cursor-paginated list

```text
1  List returns call[Page[T]](...) with a query params struct that has Cursor and Limit
2  give the params struct cursor() and withCursor(string) helpers
3  add All<Suffix>(ctx, ..., params, opts...) iter.Seq2[*T, error] returning paginate(...)
4  test a multi-page walk, a break and an error
```

## Change transport or retry behaviour

**Review required.** `transport.go` (`do`, `send`, `retryPlan`), `errors.go`. Keep: no retry of a
401 or `auth.session_revoked`; writes retried only with an idempotency key or a 429/503 with
`Retry-After`; bodies encoded once; no credentials in logs or errors. Extend `retry_test.go`.

## Change webhook or token verification

**Security review required.** `webhook/webhook.go`, `jwt/`. Never relax: constant-time comparison,
timestamp tolerance, ES256 only, `kid` required, issuer and audience required. When the server's
signing changes, record new fixtures from its signing code into `webhook/testdata/`.

## Update dependencies and tools

```text
1  go get -u ./... && go mod tidy          (there are no module dependencies today)
2  go directive = the older of the two supported Go releases; CI matrix = both
3  tool versions: Makefile (GOLANGCI_LINT_VERSION, GOVULNCHECK_VERSION),
   lint.yml, security.yml — keep them equal
4  make ci
```

## Release

```text
1  update CHANGELOG.md: move [Unreleased] under the new version
2  set the version in internal/version/version.go to match
3  merge to main — release.yml runs the gate and pushes the next patch tag (bump with the
   workflow's manual dispatch for a minor or major)
```

A tag publishes a module version that cannot be changed, only retracted with a newer version. The
first tag is the owner's decision.
