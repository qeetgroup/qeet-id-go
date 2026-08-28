# Workflows — qeet-id-go

**Level:** L2 · **Last verified:** 2026-08-28
**Verification scope:** every command checked against the `Makefile` and `.github/workflows/`.

## Set up

```bash
git clone && cd qeet-id-go
make build
make test
```

No services required — tests use an in-process `httptest.Server`. There is no `.env`.

## Add a method to an existing resource

```text
1  verify the endpoint          read the handler in qeet-id-server — NOT just the OpenAPI doc
2  open <resource>.go           match the neighbouring method exactly
3  ctx first                    func (s *XService) Do(ctx context.Context, ...) (T, error)
4  models                       big resources use <resource>_models.go / _requests.go / _responses.go
5  validation                   internal/validation for fail-fast argument checks
6  test                         internal/testutil.NewServer — never hand-roll a mock
7  README table row             the resource table is part of the contract
8  make fmt vet lint test
```

**Never wrap an endpoint you have not confirmed exists.** `CONTRIBUTING.md` is explicit: verify
against the real backend handler.

## Add a whole resource

```text
1  new file <resource>.go in the ROOT package        (never a subpackage — ADR-0001)
2  type XService struct { t *transport.Transport }
3  add the field to Client in client.go, under the right comment banner   (flat — ADR-0002)
4  wire it in New():  X: &XService{t: t}
5  <resource>_test.go
6  README table row
7  make fmt vet lint test
```

## Add pagination to a list method

```text
1  internal/pagination.Paginate[T](ctx, startCursor, fetchPage)
2  expose an All(ctx, ...) iter.Seq2[T, error] on the service
3  yield errors as (zero, err) — do not swallow
4  test both a single page and a multi-page walk
```

The `next == cursor` loop guard is already in `Paginate`. Do not reimplement it per resource.

## Change transport behaviour

**Review required.**

```text
1  internal/transport/{request,response,retry}.go
2  constants stay in internal/constants/constants.go — never inline a default
3  RETRY IDEMPOTENCY IS A SAFETY RULE: 429 always; 5xx only for GET/DELETE
4  internal/transport tests
5  make test    (-race matters here)
```

**Never make a non-idempotent method retryable.** A retried POST duplicates a server-side write.

## Change cryptography

**Security review required.**

```text
1  internal/auth/jwks.go (ES256, kid, cache TTLs) or webhooks_verify.go (HMAC)
2  cite the RFC — do not improvise
3  constant-time comparison for signatures
4  tests incl. NEGATIVE cases: bad signature, unknown kid, expired token
5  make test
```

Never weaken verification to make a test pass. Never add a "skip verification" flag.

## Release

```text
1  bump internal/version/version.go
2  CHANGELOG.md
3  docs/release-notes/vX.Y.Z.md
4  make fmt vet lint test
5  tag vX.Y.Z and push the tag
```

`release.yml` runs vet/build/test then creates a GitHub release. **A Go module has nothing to
build — the tag is the artifact.**

> This SDK has never been fetched by the module proxy. A first real release is a product decision,
> not a routine step — see `qeet-id-context/DRIFT-REGISTER.md` QID-009.

## Finish any task

```bash
make fmt && make vet && make lint && make test
git diff
```

`make fmt` **fails** on unformatted files rather than rewriting them, matching CI.

### Escalate rather than proceed

Adding a dependency · changing an exported name · changing the auth scheme · weakening verification ·
making a non-idempotent method retryable · anything requiring a change in another repository.
