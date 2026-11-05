# Architecture

## Packages

```text
github.com/qeetgroup/qeet-id-go          package qeetid — the API client
github.com/qeetgroup/qeet-id-go/webhook  package webhook — verify and parse deliveries
github.com/qeetgroup/qeet-id-go/jwt      package jwt — verify access tokens, middleware
github.com/qeetgroup/qeet-id-go/internal/version
```

Why two leaf packages and no more: [ADR-0004](../design-decisions/ADR-0004-sdk-design.md).

## The client

```text
caller
  ↓  qeetid.NewClient(opts...) → *Client          options.go, client.go
Client.Users, Client.Roles, ... (48 services)      one file per service
  ↓  call[T] / do(ctx, *request, out, opts)        transport.go
     expandPath → credential (auth.go) → send → errorFromResponse (errors.go) → retryPlan
  ↓
api.id.qeet.in
```

Every service method builds a `request` (method, path template, path arguments, query, body,
whether the endpoint accepts an `Idempotency-Key`, whether it is public) and hands it to
`Client.do`:

1. Request options are applied; an idempotency key on an endpoint that ignores it is refused.
2. Path arguments are checked (none may be empty) and escaped.
3. A create call that accepts `Idempotency-Key` gets a random UUID when the caller gave none and
   retries are on.
4. The body is encoded once, so every retry sends identical bytes.
5. Each attempt runs with its own timeout and the client's credential: `ApiKey`, a bearer token
   from a `TokenSource`, or the client credentials token source. OAuth calls use HTTP Basic with
   the application's credentials instead; public calls send none.
6. A 2xx body is decoded into the result. Anything else becomes an `*APIError` (error envelope) or
   `*OAuthError` (RFC 6749 shape).
7. `retryPlan` decides whether to retry and how long to wait; see ADR-0004.

Lists decode into `Page[T]`. `paginate` turns cursor lists into `iter.Seq2[*T, error]`.

## Credentials

| Option | Header | Source |
|---|---|---|
| `WithAPIKey` | `Authorization: ApiKey <key>` | `apiKeyCredential` |
| `WithTokenSource` | `Authorization: Bearer <token>` | caller's `TokenSource` |
| `WithClientCredentials` | `Authorization: Bearer <token>` | `clientCredentialsSource`: one cached token, renewed a minute before expiry, one fetch at a time |

## Webhook package

`Verifier.Verify` reads `X-Qeet-Timestamp` and `X-Qeet-Signature`, checks the timestamp against the
tolerance, and compares `HMAC-SHA256(secret, timestamp + "." + body)` with each `sha256=` value in
constant time, for each configured secret. `ConstructEvent` decodes the envelope;
`Handler` wraps both for `net/http`.

## JWT package

`Verifier.Verify` checks the header (`alg` ES256, `kid` present, no `crit`), the signature with the
key named by `kid`, then `exp`/`nbf` with leeway, `iss` and `aud`. `keySet` caches the JWKS for
`DefaultKeyCacheTTL`, refreshes on an unknown `kid` at most once per `DefaultMinRefreshInterval`,
shares one fetch between concurrent callers, and keeps serving cached keys when a refresh fails.

## Tests

| Test | Covers |
|---|---|
| `routes_test.go`, `routes_cases_test.go` | every method's method, path, credentials and idempotency header |
| `retry_test.go` | the retry rules, `Retry-After`, deadlines, timeouts (in `testing/synctest`) |
| `errors_test.go` | envelope and OAuth error parsing, sentinels |
| `oauth_test.go` | every OAuth error code, `expires_in`, client authentication, the token source |
| `resources_test.go` | pagination, truncated lists, invitations, LDAP sign-in, imports, query encoding |
| `webhook/webhook_test.go` | deliveries recorded from the server's signing code, refusals, rotation, handler |
| `jwt/jwt_test.go` | valid and refused tokens, key rotation, refresh cooldown, outages, middleware |
| `example_test.go` and each package's examples | runnable documentation |
