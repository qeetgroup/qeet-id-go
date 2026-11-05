# ADR-0004: The SDK design for the 1.0 API contract

**Status:** Accepted (2026-11-05). Amends ADR-0001, keeps ADR-0002, supersedes ADR-0003.

## Context

The first draft of this SDK was written against an earlier API. By the server's 1.0 release the
contract had moved: new fields and endpoints, a different webhook signature, a token endpoint split
in two, an RFC 6749 error shape on the OAuth endpoints, idempotency keys, capped lists and a
revoked-session error code. The draft also had gaps a production SDK cannot have: no typed error
sentinels, retries that ignored idempotency keys, a JWT verifier that accepted tokens without a
`kid`, and a webhook verifier that could not verify a real delivery. No version had been tagged, so
the API could be redesigned without breaking anyone.

## Decision

### Practices followed, with sources

| Practice | Source |
|---|---|
| One client type, `NewClient(opts ...Option)`; required values as arguments, optional ones as variadic options that can fail at construction | [Google Go style, best practices: option structures and variadic options](https://google.github.io/styleguide/go/best-practices#option-structure); [Azure SDK Go guidelines: client constructors](https://azure.github.io/azure-sdk/golang_introduction.html) |
| `context.Context` first on every call that does I/O; `error` last | [Go Code Review Comments: Contexts](https://go.dev/wiki/CodeReviewComments#contexts); [Google Go style decisions](https://google.github.io/styleguide/go/decisions#contexts) |
| Every method accepts an options parameter so options can be added without breaking callers | [Azure SDK Go guidelines: method options](https://azure.github.io/azure-sdk/golang_introduction.html) |
| No stutter (`qeetid.Client`, `client.Users.Get`), initialisms in one case (`ID`, `URL`, `HTTP`, `JWKS`), no `Get` prefix on accessors, short receivers | [Google Go style decisions: naming](https://google.github.io/styleguide/go/decisions#naming); [Go Code Review Comments: Initialisms](https://go.dev/wiki/CodeReviewComments#initialisms) |
| Doc comments on every exported name, starting with the name; runnable examples | [Effective Go: commentary](https://go.dev/doc/effective_go#commentary); [Go Code Review Comments: Doc Comments, Examples](https://go.dev/wiki/CodeReviewComments#doc-comments) |
| Structured error type plus sentinels matched with `errors.Is`/`errors.As`; lower-case error strings; `%w` to keep causes | [Google Go style, best practices: errors](https://google.github.io/styleguide/go/best-practices#error-handling); [Azure SDK Go guidelines: errors](https://azure.github.io/azure-sdk/golang_introduction.html) |
| Retries with exponential backoff and jitter, only for safe requests, honouring `Retry-After`, never retrying a 4xx without `Retry-After` | [Azure core guidelines: retry policy](https://azure.github.io/azure-sdk/general_azurecore.html); [AWS SDK retry behaviour](https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html) |
| Idempotency keys sent automatically so client retries are safe | [stripe-go README: automatic retries](https://github.com/stripe/stripe-go) |
| Pagination as `iter.Seq2` returned by an `All` method; one page with `List` | [Go blog: range over function types](https://go.dev/blog/range-functions); [stripe-go list iteration](https://github.com/stripe/stripe-go) |
| No secrets in logs; structured logging with `log/slog`; telemetry in the User-Agent with an application suffix | [Azure core guidelines: logging and telemetry](https://azure.github.io/azure-sdk/general_azurecore.html) |
| Webhook verification in its own package: constant-time HMAC, timestamp tolerance, multiple signatures during rotation | [stripe-go `webhook` package](https://pkg.go.dev/github.com/stripe/stripe-go/v82/webhook); [WorkOS Go webhook verifier](https://github.com/workos/workos-go) |
| Token verification that pins the algorithm, requires `kid`, issuer and audience, and caches keys with a refresh cooldown | [RFC 8725, JWT best current practices](https://www.rfc-editor.org/rfc/rfc8725); [RFC 6750, bearer token errors](https://www.rfc-editor.org/rfc/rfc6750) |
| OAuth errors in the RFC 6749 §5.2 shape; interoperate with `golang.org/x/oauth2` through a small interface instead of a dependency | [RFC 6749 §5.2](https://www.rfc-editor.org/rfc/rfc6749#section-5.2); [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) |
| Few packages; split only what is conceptually distinct and used separately | [Organizing a Go module](https://go.dev/doc/modules/layout); [Google Go style, best practices: package size](https://google.github.io/styleguide/go/best-practices#package-size) |
| Semantic versions as tags; v0 while the API may change; `/v2` only for a new major; `retract` for a bad release | [Module version numbering](https://go.dev/doc/modules/version-numbers); [go.mod `retract`](https://go.dev/ref/mod#go-mod-file-retract) |
| Support the two most recent Go releases | [Go release policy](https://go.dev/doc/devel/release#policy) |

Also consulted for client shape: [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2),
[Google Cloud Go](https://github.com/googleapis/google-cloud-go) and
[go-auth0](https://github.com/auth0/go-auth0).

### Layout

- **Root package `qeetid`**: the client, its options, credentials, transport, errors and
  pagination as unexported machinery, and one file per service holding the service, its methods
  and its types. Keeping the machinery unexported in the root package removes the type aliases the
  `internal/transport` layer needed (ADR-0003).
- **`webhook`**: receivers of webhook deliveries need no API client and usually run in a different
  service. A separate package gives them a small surface with natural names (`webhook.Verifier`,
  `webhook.Event`), as stripe-go does.
- **`jwt`**: resource servers that verify access tokens need no API client either. The package is
  named for what it verifies; it does not collide with anything in the module.
- **`internal/version`**: the version string and User-Agent, shared by all three.

Service fields stay flat on `Client` (ADR-0002).

### Contract fidelity

Models are written from the five OpenAPI documents of `qeet-id-server` and checked against the
handlers where a document is loose. Timestamps are `time.Time`, enumerations typed string
constants, nullable values pointers (strings excepted, where empty means absent). In request
bodies, empty strings are omitted, because the server reads an absent field as empty while `""`
fails to decode a UUID field. A route test sends every method's request and checks its method,
path, credentials and idempotency header.

### Retries

Safe-to-replay requests (GET, HEAD, DELETE, or any request carrying an idempotency key) are retried
after network errors and 408/500/502/503/504. A 429 or 503 **with** `Retry-After` is retried for
any method: the server's rate limiter and overload guard refuse before the handler runs, and the
contract says such a request "changed nothing". A 401, `auth.session_revoked` above all, is never
retried, even though the server marks that error `retryable` — the flag tells a browser client to
sign in again, not to resend.

### Dependencies

None outside the standard library. `golang.org/x/oauth2` users adapt their token source with a
three-line `TokenSourceFunc`.

## Consequences

- Every exported name changed shape; the CHANGELOG lists the breaking changes.
- `goimports` writes `qeetid "github.com/qeetgroup/qeet-id-go"`, because the last element of the
  module path is not the package name. Changing either is a decision for the owner.
- Go 1.26 is the minimum, following the support policy. The examples use `errors.AsType` and
  `new(expr)`, both new in Go 1.26, and the tests use `testing/synctest`.
