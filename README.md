<div align="center">

# 🔐 Qeet ID for Go

**The Go SDK for Qeet ID: manage tenants, people, access and federation from your servers, receive webhooks safely, and verify access tokens in your own APIs.**

[![Go](https://img.shields.io/badge/Go-1.26%20%7C%201.27-00ADD8?logo=go&logoColor=white)](#-quick-start)
[![context.Context](https://img.shields.io/badge/context-aware-00ADD8)](#-quick-start)
[![Errors](https://img.shields.io/badge/errors-typed-7C3AED)](#-errors)
[![Pagination](https://img.shields.io/badge/pagination-iter.Seq2-0EA5E9)](#-pagination)
[![Dependencies](https://img.shields.io/badge/dependencies-standard%20library%20only-16A34A)](go.mod)
[![Tests](https://img.shields.io/badge/tests-passing-16A34A)](#-building)

[📖 Docs](https://docs.qeet.in) · [🚀 Quick start](#-quick-start) · [🧩 API](https://apis.qeet.in) · [📝 Changelog](CHANGELOG.md)

</div>

## 📦 What's inside

| | |
|:--|:--|
| 🧭 **48 services** | One client with a service per area of the API — `client.Users`, `client.Roles`, `client.Webhooks`, `client.SAMLConnections` and more — about 330 operations in all |
| 🔐 **Three ways to authenticate** | An API key, a service principal with client credentials (tokens cached and renewed for you), or any bearer token source |
| 🚨 **Typed errors** | `*APIError` with the server's code, message and request ID, `*OAuthError` for token endpoints, and sentinels for `errors.Is` |
| 🔁 **Safe retries** | Exponential backoff with jitter, `Retry-After`, idempotency keys; never a retry that could apply a write twice |
| 📄 **Pagination** | `Page[T]` for one page, `All` iterators (`iter.Seq2`) for every item, `Truncated` on capped lists |
| 🪝 **Webhooks** | Package `webhook`: constant-time signature checks, replay protection, secret rotation, typed events, a `net/http` handler |
| 🔑 **Token verification** | Package `jwt`: ES256 access tokens checked locally against the published keys, with middleware |
| 🤝 **OAuth** | Client credentials, code exchange, refresh, token exchange, device and backchannel flows, introspection, revocation |
| 🪶 **No dependencies** | The Go standard library only |

## 🚀 Quick start

**Requirements:** Go 1.26 or 1.27 — the SDK supports the two most recent Go releases, following the Go team's
[release policy](https://go.dev/doc/devel/release#policy) — and a server-side credential for your Qeet ID tenant: an
API key, or a service principal's client ID and secret.

**Install:**

```bash
go get github.com/qeetgroup/qeet-id-go
```

The package name is `qeetid`; `goimports` writes the import as `qeetid "github.com/qeetgroup/qeet-id-go"`.

**Use.** Create one client and share it — it is safe for concurrent use:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	qeetid "github.com/qeetgroup/qeet-id-go"
)

func main() {
	client, err := qeetid.NewClient(qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	user, err := client.Users.Create(ctx, qeetid.CreateUserParams{
		TenantID:    os.Getenv("QEETID_TENANT_ID"),
		Email:       "ada@brightleaffoods.in",
		DisplayName: "Ada Lovelace",
	})
	if errors.Is(err, qeetid.ErrConflict) {
		log.Fatal("a user with this email already exists")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created", user.ID)
}
```

Every method takes a `context.Context` first, then the path parameters in the order the API path names them, then a
parameter struct when the call has one, and optional `RequestOption` values last. It returns an `error` last.

> [!WARNING]
> API keys and client secrets belong on servers. Never ship one to a browser or a mobile app.

## ⚙️ Options

```go
client, err := qeetid.NewClient(
	qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")),
	qeetid.WithTimeout(10*time.Second),
	qeetid.WithMaxRetries(3),
	qeetid.WithLogger(slog.Default()),
	qeetid.WithUserAgentSuffix("billing-worker/2.3.0"),
)
```

| Option | Default | What it does |
|:--|:--|:--|
| `WithAPIKey(key)` | — | Sends `Authorization: ApiKey <key>` |
| `WithClientCredentials(id, secret, scopes...)` | — | Fetches, caches and renews service principal tokens |
| `WithTokenSource(ts)` | — | Sends `Authorization: Bearer <token>` from your own source |
| `WithBaseURL(url)` | `https://api.id.qeet.in` | https only; http is allowed for a loopback host such as a local development server |
| `WithHTTPClient(c)` | a new `http.Client` | Your own transport, proxy or TLS settings |
| `WithTimeout(d)` | 30 s | Bounds each attempt; the context bounds the whole call |
| `WithMaxRetries(n)` | 2 | `0` turns retries off |
| `WithRetryDelay(min, max)` | 500 ms, 8 s | Bounds of the backoff between retries |
| `WithLogger(l)` | no logging | A `*slog.Logger`: each attempt at debug, each retry at warn — never credentials, headers or bodies |
| `WithUserAgentSuffix(s)` | — | Appended to `qeet-id-go/<version> (<go version>; <os>/<arch>)` |
| `WithHeader(k, v)` | — | An extra header on every request; the client's own headers cannot be overridden |

Use exactly one credential option. A client without one can still call the public endpoints — `client.OAuth`,
`client.Credentials.Verify`, `client.LDAP.Authenticate` — and returns `qeetid.ErrNoCredentials` for anything else,
without sending the request.

## 🔐 Authentication

| Credential | Use it when | Option |
|:--|:--|:--|
| 🗝️ **API key** | A backend calls the management API | `WithAPIKey` |
| 🤖 **Service principal** | You want short-lived tokens instead of a long-lived key | `WithClientCredentials` |
| 🎫 **Token source** | You already hold a token, such as an AI agent's | `WithTokenSource` |

```go
// Service principal: one token shared across goroutines, renewed a minute before it expires.
client, err := qeetid.NewClient(qeetid.WithClientCredentials(
	os.Getenv("QEETID_CLIENT_ID"), os.Getenv("QEETID_CLIENT_SECRET"),
	"user.read", "user.write", // optional: narrow the token's scopes
))
```

A [`golang.org/x/oauth2`](https://pkg.go.dev/golang.org/x/oauth2) token source fits in a few lines, so the SDK does not
depend on that module — `clientcredentials.Config` works against Qeet ID unchanged with `TokenURL` set to
`https://api.id.qeet.in/v1/oauth/token`:

```go
qeetid.WithTokenSource(qeetid.TokenSourceFunc(func(context.Context) (string, error) {
	t, err := ts.Token() // ts is an oauth2.TokenSource
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}))
```

## 🧭 Services

| Area | Services |
|:--|:--|
| 👥 **Tenants and people** | `Tenants` · `Users` · `Groups` · `Invitations` · `Domains` · `Branding` |
| 🛡️ **Access control** | `Roles` · `Permissions` · `ABACPolicies` · `Relationships` · `AuthZEN` |
| 🔏 **Sign-in and security** | `AuthPolicy` · `SecurityPolicy` · `IPRules` · `Threats` |
| 🌐 **Applications and federation** | `OAuth` · `OIDCClients` · `OAuthResources` · `OAuthGrants` · `DeviceAuthorizations` · `SigningKeys` · `ServicePrincipals` · `SocialProviders` · `SAMLConnections` · `SAMLServiceProviders` · `OIDCConnections` · `LDAP` · `Directory` · `SCIM` · `AdminPortalLinks` |
| 🛠️ **Developer** | `APIKeys` · `Agents` · `Credentials` · `AuthHooks` · `Webhooks` · `Secrets` · `TokenVault` |
| 📊 **Operations** | `Audit` · `Analytics` · `Compliance` · `Retention` · `GDPR` · `LogSinks` · `RateLimits` · `EmailTemplates` · `EmailDelivery` · `Billing` · `Entitlements` |

Each method's documentation names the endpoint it calls. A few of them:

```go
// Give a user a role, then ask whether they hold a permission — and why.
err = client.Roles.AssignToUser(ctx, userID, tenantID, roleID)
check, err := client.Permissions.Check(ctx, &qeetid.CheckPermissionParams{
	UserID: userID, TenantID: tenantID, Permission: "billing.write", Explain: true,
})

// Relationship-based access: may this user view the document?
rel, err := client.Relationships.Check(ctx, tenantID, qeetid.CheckRelationParams{
	Object: "document:readme", Relation: "viewer", UserID: userID,
}, nil)

// Directory sign-in answers with a token pair or a second-factor challenge — never both.
res, err := client.LDAP.Authenticate(ctx, connectionID, qeetid.LDAPAuthenticateParams{Username: "ada", Password: pw})
switch r := res.(type) {
case *qeetid.TokenPair:
	// signed in: r.AccessToken, r.RefreshToken
case *qeetid.MFAChallenge:
	// complete r.MFAToken at POST /v1/auth/mfa
}
```

## 🚨 Errors

```go
user, err := client.Users.Get(ctx, userID)
switch {
case errors.Is(err, qeetid.ErrNotFound):
	// no such user, or not visible to this credential
case err != nil:
	if apiErr, ok := errors.AsType[*qeetid.APIError](err); ok {
		log.Printf("qeet id: %d %s: %s (request ID %s)", apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.RequestID)
	}
	return err
}
```

| Check | Matches |
|:--|:--|
| `errors.Is(err, qeetid.ErrInvalidRequest)` | 400 and 422, and OAuth errors other than `invalid_client` |
| `errors.Is(err, qeetid.ErrUnauthorized)` | 401, and OAuth `invalid_client` |
| `errors.Is(err, qeetid.ErrSessionRevoked)` | 401 `auth.session_revoked`: the session ended — sign the person in again |
| `errors.Is(err, qeetid.ErrPaymentRequired)` | 402: the tenant's plan does not cover the request |
| `errors.Is(err, qeetid.ErrForbidden)` · `ErrNotFound` · `ErrConflict` | 403 · 404 · 409 |
| `errors.Is(err, qeetid.ErrRateLimited)` · `ErrUnavailable` | 429 · 503; `APIError.RetryAfter` says how long to wait |
| `errors.As(err, &apiErr)` with `*qeetid.APIError` | Any API error: `StatusCode`, stable `Code`, `Message`, `Detail`, `Fields`, `Metadata`, `RequestID`, `RetryAfter` |
| `errors.As(err, &oauthErr)` with `*qeetid.OAuthError` | Token endpoint errors in the RFC 6749 shape, with a typed `Code` such as `qeetid.OAuthInvalidGrant` |
| `errors.Is(err, context.DeadlineExceeded)` | Network failures wrap their cause |

Branch on `Code`, never on `Message`. Error text never contains a credential.

## 🔁 Retries & timeouts

The client retries only when a retry cannot apply a request twice:

| Failure | Retried for |
|:--|:--|
| 🌐 Network error, 408, 500, 502, 503, 504 | GET and DELETE, and create calls carrying an idempotency key |
| ⏳ 429 or 503 **with** `Retry-After` — the server refused before acting | Every request |
| 🔄 409 `idempotency.in_progress` | Requests carrying an idempotency key |
| ⛔ 401, `auth.session_revoked` included, and any other 4xx | Never |

Backoff is exponential with jitter: between half and all of 500 ms × 2ⁿ, capped at 8 s. A `Retry-After` is honoured
exactly, up to `qeetid.MaxRetryAfter` (60 s); a longer one comes back to you in `APIError.RetryAfter`. The client never
sleeps past the context's deadline, and every attempt gets its own `WithTimeout`.

**Idempotency keys.** Nine calls accept one: `Users.Create`, `Groups.Create`, `Invitations.Create`,
`ServicePrincipals.Create`, `APIKeys.Create`, `Agents.Create`, `Webhooks.Create`, `Billing.Checkout` and
`Billing.PayInvoice`. The client sends a random key with them, which makes its own retries safe. Derive your own to
make retries across restarts or queue redeliveries safe too — any other call refuses the option:

```go
user, err := client.Users.Create(ctx, params, qeetid.WithIdempotencyKey("hr-import-2026-11-05-row-17"))
```

## 📄 Pagination

`List` methods return a `*qeetid.Page[T]` with `Items`, `NextCursor`, `Total` and `Truncated`. Cursor lists also have an
`All` method — a Go iterator that fetches the next page only when the loop needs it:

```go
for user, err := range client.Users.All(ctx, &qeetid.ListUsersParams{Limit: 100, Status: qeetid.UserStatusActive}) {
	if err != nil {
		return err
	}
	fmt.Println(user.Email) // break whenever you like; paging stops with the loop
}
```

Some lists come back whole but capped by the server. When such a list is longer than the cap, `Page.Truncated` is true:
narrow the request with a filter to see the rest.

## 🪝 Webhooks

Package [`webhook`](webhook) verifies the deliveries Qeet ID sends and parses them into typed events. Each delivery is
a POST whose body is `{id, type, tenant_id, occurred_at, data}`, signed in `X-Qeet-Signature` with `sha256=` and the
HMAC-SHA256 of `X-Qeet-Timestamp`, a `.`, and the raw body.

```go
import "github.com/qeetgroup/qeet-id-go/webhook"

v, err := webhook.NewVerifier(os.Getenv("QEETID_WEBHOOK_SECRET"))
if err != nil {
	log.Fatal(err)
}
http.Handle("POST /webhooks/qeet-id", v.Handler(func(ctx context.Context, e *webhook.Event) error {
	switch e.Type {
	case webhook.UserCreated:
		var d webhook.UserData
		if err := e.DecodeData(&d); err != nil {
			return err
		}
		// provision d.UserID
	case webhook.SessionRevoked:
		var d webhook.SessionRevokedData
		if err := e.DecodeData(&d); err != nil {
			return err
		}
		// stop accepting tokens whose sid is d.SessionID
	}
	return nil // 204 — return an error to answer 500, and Qeet ID retries later
}))
```

| | |
|:--|:--|
| ⚖️ **Constant time** | Signatures are compared with `hmac.Equal` |
| ⏱️ **Replay protection** | A timestamp more than five minutes from the verifier's clock is refused; `webhook.WithTolerance` changes the window |
| 🔄 **Secret rotation** | For 24 hours after a rotation Qeet ID sends two signatures; either one is accepted, and `webhook.WithSecret` lets you hold both secrets too |
| 🧾 **Typed events** | The whole event catalogue as `webhook.EventType` constants, with data types for user, group, session and claims events |
| 🧪 **Testing** | `webhook.Sign` builds signed deliveries for tests of your own handlers |

Handling the request yourself? Pass the raw body unchanged — re-encoding the JSON breaks the signature:

```go
body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
event, err := v.ConstructEvent(body, r.Header)
if errors.Is(err, webhook.ErrInvalidDelivery) {
	http.Error(w, "invalid signature", http.StatusUnauthorized)
	return
}
```

Deliveries are at least once and in no guaranteed order: deduplicate on `Event.ID`.

## 🔑 Verifying access tokens

Package [`jwt`](jwt) verifies Qeet ID access tokens in your own API, with no call to Qeet ID per request:

```go
import "github.com/qeetgroup/qeet-id-go/jwt"

meta, err := client.OAuth.Discover(ctx) // the issuer and the key set's address
v, err := jwt.NewVerifier(meta.Issuer, "https://api.example.com", jwt.WithJWKSURL(meta.JWKSURI))

mux.Handle("GET /invoices", v.Middleware(jwt.RequireScope("invoices.read")(http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		claims, _ := jwt.ClaimsFromContext(r.Context())
		fmt.Fprintln(w, claims.UserID, claims.TenantID, claims.SessionID)
	}))))
```

| | |
|:--|:--|
| 🔒 **Strict** | ES256 only, a `kid` that names a published key, `exp` and `nbf` with 30 s leeway, issuer and audience both required |
| 🔄 **Key rotation** | Keys cached for ten minutes and fetched again when a token names a new `kid` — no restart needed |
| 🧱 **Resilient** | Unknown key IDs trigger at most one fetch a minute; cached keys keep working through an outage of the key endpoint |
| 🌐 **Middleware** | 401 for a missing or invalid token, 403 for a missing scope, 503 when the keys cannot be fetched |

Local verification cannot see a session that ended before its token expired. Where revocation must take effect at once,
handle the `session.revoked` webhook or call `client.OAuth.Introspect`.

## 🤝 OAuth

`client.OAuth` runs the requests in which your application is the OAuth client of Qeet ID. They authenticate the
application with its client ID and secret, not with the client's own credential:

```go
tok, err := client.OAuth.ExchangeCode(ctx, qeetid.ExchangeCodeParams{
	ClientID: appID, ClientSecret: appSecret,
	Code: code, RedirectURI: "https://app.example.com/callback", CodeVerifier: verifier,
})
if oauthErr, ok := errors.AsType[*qeetid.OAuthError](err); ok && oauthErr.Code == qeetid.OAuthInvalidGrant {
	// the code was used, has expired, or does not match
}
fmt.Println(tok.ExpiresIn, tok.Expiry())
```

| Method | Endpoint |
|:--|:--|
| `ClientCredentials` | `POST /v1/oauth/token` |
| `ExchangeCode` · `RefreshToken` · `ExchangeToken` · `PollDeviceToken` · `PollBackchannelToken` | `POST /v1/oauth/token-code` |
| `AuthorizeDevice` · `AuthorizeBackchannel` | `POST /v1/oauth/device_authorization` · `POST /v1/oauth/bc-authorize` |
| `Introspect` · `Revoke` | `POST /oauth/introspect` · `POST /oauth/revoke` |
| `UserInfo` · `Discover` | `GET /v1/oauth/userinfo` · `GET /.well-known/openid-configuration` |

## 🧪 Examples

Runnable examples sit beside the code (`example_test.go` in each package), and complete programs are in
[examples/](examples):

| Program | Shows |
|:--|:--|
| 🚀 [quickstart](examples/quickstart) | An API key, creating a user, assigning a role, paging through users |
| 🤖 [service-principal](examples/service-principal) | Client credentials and debug logging |
| 🪝 [webhook-receiver](examples/webhook-receiver) | Verifying and handling webhook deliveries |
| 🔑 [resource-server](examples/resource-server) | Protecting an API with access-token middleware |

**Left out on purpose:** calls that act for the signed-in person (`/v1/me`, `/v1/auth/*`, passkey and MFA enrolment),
browser redirect flows, the SCIM protocol endpoints identity providers call, and console features such as the activity
stream, search, notifications and the Qeet AI assistant.

## 🔨 Building

```sh
make ci     # gofmt, go vet, go build, go test -race, golangci-lint, govulncheck — what CI runs
```

Every command is in [CONTRIBUTING.md](CONTRIBUTING.md); what changed between versions is in the
[changelog](CHANGELOG.md).

## 📚 Learn more

| | |
|:--|:--|
| 📖 [docs.qeet.in](https://docs.qeet.in) | Qeet ID documentation |
| 🧩 [apis.qeet.in](https://apis.qeet.in) | The Qeet ID API reference |
| 🏗️ [Architecture](docs/architecture/ARCHITECTURE.md) | How the client, the `webhook` and `jwt` packages are built |
| 🧠 [Design decision](docs/design-decisions/ADR-0004-sdk-design.md) | The practices this SDK follows, with their sources |
| 🤝 [Contributing](CONTRIBUTING.md) | Setting up, commands, adding a method |
| 📝 [Changelog](CHANGELOG.md) | What's in each version |
| 🔒 [Security](SECURITY.md) | How to report a vulnerability |
| 📄 [Licence](LICENSE) | The terms of use |

---

<div align="center">

Made with 🧡 by Qeet Group · [⬆️ Back to top](#-qeet-id-for-go)

</div>
