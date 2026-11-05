# qeet-id-go

[![Go Reference](https://pkg.go.dev/badge/github.com/qeetgroup/qeet-id-go.svg)](https://pkg.go.dev/github.com/qeetgroup/qeet-id-go)
[![Go](https://img.shields.io/badge/go-1.26%20%7C%201.27-00ADD8?logo=go)](go.mod)

The Go SDK for [Qeet ID](https://qeet.in), the passkeys-first identity platform. Use it on your
servers to manage tenants, users, groups, roles, applications, single sign-on, webhooks and the
rest of Qeet ID, to receive webhook deliveries safely, and to verify Qeet ID access tokens in your
own APIs.

- One client, `qeetid.NewClient`, with a service per area of the API: `client.Users`,
  `client.Roles`, `client.Webhooks` and 45 more.
- Typed errors that work with `errors.Is` and `errors.As`, carrying the server's error code and
  request ID.
- Safe automatic retries with backoff, `Retry-After` support and idempotency keys.
- Pagination with Go iterators.
- Package [`webhook`](#receiving-webhooks) to verify and parse webhook deliveries.
- Package [`jwt`](#verifying-access-tokens) to verify access tokens locally, with HTTP middleware.
- No dependencies outside the Go standard library.

> **Status:** the module has not had its first tagged release yet. Until v1.0.0 its API may change
> between minor versions; see [CHANGELOG.md](CHANGELOG.md).

## Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [Authentication](#authentication)
- [Errors](#errors)
- [Retries and timeouts](#retries-and-timeouts)
- [Idempotency](#idempotency)
- [Pagination](#pagination)
- [Receiving webhooks](#receiving-webhooks)
- [Verifying access tokens](#verifying-access-tokens)
- [OAuth flows](#oauth-flows)
- [Services](#services)
- [What the SDK does not cover](#what-the-sdk-does-not-cover)
- [Logging](#logging)
- [Examples](#examples)
- [Versioning and Go support](#versioning-and-go-support)
- [Contributing, security and license](#contributing-security-and-license)

## Requirements

Go 1.26 or 1.27, the two Go releases currently supported by the Go team. The SDK follows Go's
support policy: when a new Go release comes out, the oldest one is dropped.

You need a Qeet ID tenant and a server-side credential: an API key, or a service principal's client
ID and secret.

## Installation

```bash
go get github.com/qeetgroup/qeet-id-go
```

The package name is `qeetid`. `goimports` writes the import with that name:

```go
import qeetid "github.com/qeetgroup/qeet-id-go"
```

## Quick start

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
	// Create one client and share it; it is safe for concurrent use.
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

Every method takes a `context.Context` first and returns an `error` last. Path parameters come
next, in the order the API path names them, then a parameter struct when the call has one, and
optional `RequestOption` values last.

> [!WARNING]
> API keys and client secrets belong on servers. Never ship one to a browser or a mobile app.

## Configuration

Pass options to `NewClient`:

```go
client, err := qeetid.NewClient(
	qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")),
	qeetid.WithTimeout(10*time.Second),             // per attempt; default 30s
	qeetid.WithMaxRetries(3),                       // default 2; 0 turns retries off
	qeetid.WithRetryDelay(time.Second, 15*time.Second),
	qeetid.WithLogger(slog.Default()),              // default: no logging
	qeetid.WithUserAgentSuffix("billing-worker/2.3.0"),
	qeetid.WithHeader("X-Correlation-Id", corrID),  // sent with every request
	qeetid.WithHTTPClient(&http.Client{Transport: myTransport}),
)
```

| Option | Default | Notes |
|---|---|---|
| `WithAPIKey(key)` | none | Sends `Authorization: ApiKey <key>`. |
| `WithClientCredentials(id, secret, scopes...)` | none | Fetches and renews service principal tokens. |
| `WithTokenSource(ts)` | none | Sends `Authorization: Bearer <token>` from your own source. |
| `WithBaseURL(url)` | `https://api.id.qeet.in` | Must be https; http is allowed only for a loopback host such as a local development server. |
| `WithHTTPClient(c)` | a new `http.Client` | Custom transport, proxy or TLS settings. |
| `WithTimeout(d)` | 30 seconds | Bounds each attempt. The context bounds the whole call. |
| `WithMaxRetries(n)` | 2 | Retries happen only when they are safe; see below. |
| `WithRetryDelay(min, max)` | 500 ms, 8 s | Bounds of the exponential backoff. |
| `WithLogger(l)` | none | A `*slog.Logger`; see [Logging](#logging). |
| `WithUserAgentSuffix(s)` | none | Appended to `qeet-id-go/<version> (<go version>; <os>/<arch>)`. |
| `WithHeader(k, v)` | none | Cannot override the headers the client manages. |

Use exactly one credential option. A client without one can still call the public endpoints, such
as `client.OAuth`, `client.Credentials.Verify` and `client.LDAP.Authenticate`; any other call
returns `qeetid.ErrNoCredentials` without sending a request.

## Authentication

**API key.** The simplest choice for a backend:

```go
client, err := qeetid.NewClient(qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")))
```

**Service principal.** Short-lived tokens instead of a long-lived key. The client obtains a token
with the OAuth client credentials grant, shares it across goroutines, and fetches a new one a
minute before it expires:

```go
client, err := qeetid.NewClient(qeetid.WithClientCredentials(
	os.Getenv("QEETID_CLIENT_ID"), os.Getenv("QEETID_CLIENT_SECRET"),
	"user.read", "user.write", // optional: narrow the token's scopes
))
```

**Any other token.** `WithTokenSource` takes anything with a
`Token(ctx context.Context) (string, error)` method, such as an AI agent's token. A
[`golang.org/x/oauth2`](https://pkg.go.dev/golang.org/x/oauth2) token source fits in a few lines,
so the SDK does not depend on that module:

```go
qeetid.WithTokenSource(qeetid.TokenSourceFunc(func(context.Context) (string, error) {
	t, err := ts.Token() // ts is an oauth2.TokenSource
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}))
```

If you already use [`golang.org/x/oauth2/clientcredentials`](https://pkg.go.dev/golang.org/x/oauth2/clientcredentials),
it works against Qeet ID unchanged: set `TokenURL` to `https://api.id.qeet.in/v1/oauth/token` and
adapt its token source as above.

## Errors

When the API answers with an error status, methods return an `*qeetid.APIError`:

```go
user, err := client.Users.Get(ctx, userID)
switch {
case errors.Is(err, qeetid.ErrNotFound):
	// no such user, or not visible to this credential
case errors.Is(err, qeetid.ErrRateLimited):
	// slow down; the client already retried
case err != nil:
	if apiErr, ok := errors.AsType[*qeetid.APIError](err); ok {
		log.Printf("qeet id: %d %s: %s (request ID %s)",
			apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.RequestID)
	}
	return err
}
```

`APIError` carries the whole error envelope: `StatusCode`, the stable machine `Code` (branch on
it, never on `Message`), `Message`, `Detail`, `Fields` (per-field validation errors),
`Metadata`, `RequestID` (quote it to support) and `RetryAfter`.

| Sentinel | Matches |
|---|---|
| `ErrInvalidRequest` | 400 and 422, and OAuth errors other than `invalid_client` |
| `ErrUnauthorized` | 401, and OAuth `invalid_client` |
| `ErrSessionRevoked` | 401 with code `auth.session_revoked`: the session has ended; sign the person in again |
| `ErrPaymentRequired` | 402: the tenant's plan does not cover the request |
| `ErrForbidden` | 403 |
| `ErrNotFound` | 404 |
| `ErrConflict` | 409 |
| `ErrRateLimited` | 429 |
| `ErrUnavailable` | 503 |
| `ErrNoCredentials` | an authenticated call on a client without a credential |

The OAuth token endpoints answer in the RFC 6749 shape `{"error", "error_description"}`. Those
failures are `*qeetid.OAuthError` values with a typed `Code` such as `qeetid.OAuthInvalidGrant`.

A failure to reach the API wraps its cause, so `errors.Is(err, context.DeadlineExceeded)` and
`errors.As(err, &netErr)` work. Error text never contains credentials.

## Retries and timeouts

The client retries a failed request only when retrying cannot apply it twice:

| Failure | Retried for |
|---|---|
| Network error, 408, 500, 502, 503, 504 | GET and DELETE, and create calls that carry an idempotency key |
| 429 or 503 **with** `Retry-After` (the server refused before acting) | every request |
| 409 `idempotency.in_progress` | requests that carry an idempotency key |
| 401, including `auth.session_revoked` | never |
| any other 4xx | never |

Retries use exponential backoff with jitter (between half and all of 500 ms × 2ⁿ, capped at 8 s by
default). When the server sends `Retry-After`, the client waits exactly that long, up to
`qeetid.MaxRetryAfter` (60 seconds); a longer wait is returned to you in `APIError.RetryAfter`. The
client never sleeps past the context's deadline. Each attempt gets its own `WithTimeout`; the
context you pass bounds the call as a whole.

## Idempotency

Nine calls accept an `Idempotency-Key`: `Users.Create`, `Groups.Create`, `Invitations.Create`,
`ServicePrincipals.Create`, `APIKeys.Create`, `Agents.Create`, `Webhooks.Create`,
`Billing.Checkout` and `Billing.PayInvoice`. A request sent again with the same key within 24 hours
gets the first answer instead of acting twice.

The client sends a random key with these calls by itself, so its own retries are safe. To make
retries across process restarts or message redeliveries safe too, derive the key from your own
data:

```go
user, err := client.Users.Create(ctx, params,
	qeetid.WithIdempotencyKey("hr-import-2026-11-05-row-17"))
```

The same key with a different body is refused with `idempotency.key_reused`. Any other call refuses
`WithIdempotencyKey` with an error, without sending the request, because the server would ignore the
key.

## Pagination

List methods return a `*qeetid.Page[T]`:

```go
page, err := client.Users.List(ctx, &qeetid.ListUsersParams{Limit: 50, Status: qeetid.UserStatusActive})
for _, u := range page.Items { /* ... */ }
if page.HasMore() {
	next, err := client.Users.List(ctx, &qeetid.ListUsersParams{Limit: 50, Cursor: page.NextCursor})
}
```

Lists the server pages with cursors also have an `All` method that returns a Go iterator over every
item. It fetches a page only when the loop needs it; breaking out of the loop stops paging:

```go
for user, err := range client.Users.All(ctx, &qeetid.ListUsersParams{Limit: 100}) {
	if err != nil {
		return err
	}
	fmt.Println(user.Email)
}
```

Some lists come back whole but are capped by the server. When such a list is longer than the cap,
`Page.Truncated` is true: narrow the request with a filter to see the rest. `Page.Total` reports
the number of matches for lists that count them.

## Receiving webhooks

Package [`webhook`](https://pkg.go.dev/github.com/qeetgroup/qeet-id-go/webhook) verifies the
deliveries Qeet ID sends and parses them into typed events. Each delivery is a POST whose body is
the envelope `{id, type, tenant_id, occurred_at, data}`. `X-Qeet-Signature` holds `sha256=` and the
hex HMAC-SHA256, keyed by the subscription's secret, of `X-Qeet-Timestamp`, a `.`, and the raw
body.

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
		// provision d.UserID in your system
	case webhook.SessionRevoked:
		var d webhook.SessionRevokedData
		if err := e.DecodeData(&d); err != nil {
			return err
		}
		// stop accepting tokens with sid == d.SessionID
	}
	return nil // 204; returning an error answers 500 and Qeet ID retries later
}))
```

The verifier:

- compares signatures in constant time;
- refuses a timestamp more than five minutes from its clock, so a captured delivery cannot be
  replayed (`webhook.WithTolerance` changes the window);
- accepts a delivery when any of the comma-separated signatures matches. For 24 hours after a secret
  rotation Qeet ID signs with both secrets, so endpoints keep working while you deploy the new one;
  `webhook.WithSecret` lets one verifier accept both on your side too.

If you handle the request yourself, read the raw body and pass it unchanged; re-encoding the JSON
breaks the signature:

```go
body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
event, err := v.ConstructEvent(body, r.Header)
if errors.Is(err, webhook.ErrInvalidDelivery) {
	http.Error(w, "invalid signature", http.StatusUnauthorized)
	return
}
```

Deliveries are at least once and in no guaranteed order: deduplicate on `Event.ID`. Use
`webhook.Sign` to build signed deliveries in tests of your own handlers.

## Verifying access tokens

Package [`jwt`](https://pkg.go.dev/github.com/qeetgroup/qeet-id-go/jwt) verifies Qeet ID access
tokens in your own API without a call to Qeet ID per request. It accepts only ES256, requires a
`kid` that names a published key, and checks the signature, `exp` and `nbf` (with 30 seconds of
leeway), the issuer and the audience.

```go
import "github.com/qeetgroup/qeet-id-go/jwt"

// The issuer is in Qeet ID's discovery document.
meta, err := client.OAuth.Discover(ctx)
v, err := jwt.NewVerifier(meta.Issuer, "https://api.example.com", jwt.WithJWKSURL(meta.JWKSURI))

mux.Handle("GET /invoices", v.Middleware(jwt.RequireScope("invoices.read")(invoicesHandler)))

func invoicesHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := jwt.ClaimsFromContext(r.Context())
	// claims.UserID, claims.TenantID, claims.SessionID, claims.Scope ...
}
```

The verifier caches the signing keys for ten minutes and fetches them again when a token names a key
it has not seen, so key rotations need no restart. Unknown key IDs trigger at most one fetch a
minute, and cached keys keep working through an outage of the key endpoint. The middleware answers
401 for a missing or invalid token, 403 when `RequireScope` fails, and 503 when the keys cannot be
fetched.

Local verification cannot see a session that ended before its token expired. Access tokens are
short-lived; where revocation must take effect at once, handle the `session.revoked` webhook event
or call `client.OAuth.Introspect`.

## OAuth flows

`client.OAuth` runs the requests in which your application is the OAuth client of Qeet ID. These
calls authenticate the application with its client ID and secret, not with the client's own
credential:

```go
tok, err := client.OAuth.ExchangeCode(ctx, qeetid.ExchangeCodeParams{
	ClientID: appID, ClientSecret: appSecret,
	Code: code, RedirectURI: "https://app.example.com/callback", CodeVerifier: verifier,
})
if oauthErr, ok := errors.AsType[*qeetid.OAuthError](err); ok && oauthErr.Code == qeetid.OAuthInvalidGrant {
	// the code was used, expired or does not match
}
fmt.Println(tok.ExpiresIn, tok.Expiry())
```

| Method | Endpoint |
|---|---|
| `ClientCredentials` | `POST /v1/oauth/token` |
| `ExchangeCode`, `RefreshToken`, `ExchangeToken` (RFC 8693), `PollDeviceToken`, `PollBackchannelToken` | `POST /v1/oauth/token-code` |
| `AuthorizeDevice` (RFC 8628) | `POST /v1/oauth/device_authorization` |
| `AuthorizeBackchannel` (CIBA) | `POST /v1/oauth/bc-authorize` |
| `Introspect` (RFC 7662), `Revoke` (RFC 7009) | `POST /oauth/introspect`, `POST /oauth/revoke` |
| `UserInfo`, `Discover` | `GET /v1/oauth/userinfo`, `GET /.well-known/openid-configuration` |

`Token.ExpiresIn` is the lifetime in seconds; `Token.Expiry()` turns it into a time.

## Services

| Area | Services |
|---|---|
| Tenants and people | `Tenants`, `Users`, `Groups`, `Invitations`, `Domains`, `Branding` |
| Access control | `Roles`, `Permissions`, `ABACPolicies`, `Relationships`, `AuthZEN` |
| Sign-in and security | `AuthPolicy`, `SecurityPolicy`, `IPRules`, `Threats` |
| Applications and federation | `OAuth`, `OIDCClients`, `OAuthResources`, `OAuthGrants`, `DeviceAuthorizations`, `SigningKeys`, `ServicePrincipals`, `SocialProviders`, `SAMLConnections`, `SAMLServiceProviders`, `OIDCConnections`, `LDAP`, `Directory`, `SCIM`, `AdminPortalLinks` |
| Developer | `APIKeys`, `Agents`, `Credentials`, `AuthHooks`, `Webhooks`, `Secrets`, `TokenVault` |
| Operations | `Audit`, `Analytics`, `Compliance`, `Retention`, `GDPR`, `LogSinks`, `RateLimits`, `EmailTemplates`, `EmailDelivery`, `Billing`, `Entitlements` |

Each method's documentation names the API endpoint it calls. A few examples:

```go
// Give a user a role in a tenant.
err := client.Roles.AssignToUser(ctx, userID, tenantID, roleID)

// Ask whether the user holds a permission, and why.
check, err := client.Permissions.Check(ctx, &qeetid.CheckPermissionParams{
	UserID: userID, TenantID: tenantID, Permission: "billing.write", Explain: true,
})

// Relationship-based access: may Ada view the document?
rel, err := client.Relationships.Check(ctx, tenantID, qeetid.CheckRelationParams{
	Object: "document:readme", Relation: "viewer", UserID: userID,
}, nil)

// Sign a person in through an LDAP connection: a token pair or a second-factor challenge.
res, err := client.LDAP.Authenticate(ctx, connectionID, qeetid.LDAPAuthenticateParams{
	Username: "ada", Password: password,
})
switch r := res.(type) {
case *qeetid.TokenPair:
	// signed in
case *qeetid.MFAChallenge:
	// complete r.MFAToken at POST /v1/auth/mfa
}
```

## What the SDK does not cover

The SDK is for servers acting for a tenant. It leaves out, on purpose:

- end-user self-service calls that act on the signed-in person (`/v1/me`, `/v1/auth/*`, passkey
  registration, MFA enrolment, account export and deletion), which belong in the browser and
  mobile SDKs;
- browser redirect and form-post flows (OAuth authorize and logout, social and SAML sign-in, admin
  portal pages) and the SCIM protocol endpoints that identity providers call;
- console features: activity history and its live stream, universal search, notifications, the
  Qeet AI assistant, contact-sales leads and pre-signup checkout.

## Logging

`WithLogger` takes a `*slog.Logger`. The client logs each HTTP attempt at debug level (method, path,
status, duration, attempt number and the server's request ID) and each retry at warn level (reason
and wait). It never logs credentials, headers, query strings or bodies.

## Examples

Runnable examples are in the package documentation on
[pkg.go.dev](https://pkg.go.dev/github.com/qeetgroup/qeet-id-go#pkg-examples), and complete
programs are in [examples/](examples/):

| Program | Shows |
|---|---|
| [quickstart](examples/quickstart) | API key, creating a user, assigning a role, paging through users |
| [service-principal](examples/service-principal) | Client credentials and debug logging |
| [webhook-receiver](examples/webhook-receiver) | Verifying and handling webhook deliveries |
| [resource-server](examples/resource-server) | Protecting an API with access-token middleware |

## Versioning and Go support

The module follows [semantic versioning](https://semver.org). Releases are Git tags; until v1.0.0
a minor version may change the API, and [CHANGELOG.md](CHANGELOG.md) lists every breaking change.
The SDK supports the two most recent Go releases, matching the Go team's own
[release policy](https://go.dev/doc/devel/release#policy).

## Contributing, security and license

See [CONTRIBUTING.md](CONTRIBUTING.md) for development commands. Report vulnerabilities privately as
[SECURITY.md](SECURITY.md) describes. The SDK is released under the terms in [LICENSE](LICENSE).
