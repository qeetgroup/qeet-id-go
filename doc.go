// Package qeetid is the Go client for the Qeet ID API: manage tenants, users,
// groups, roles, applications, federation and every other Qeet ID resource
// from your servers.
//
// # Getting started
//
// Create one Client and share it; it is safe for concurrent use:
//
//	client, err := qeetid.NewClient(qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")))
//	if err != nil {
//		return err
//	}
//	user, err := client.Users.Get(ctx, userID)
//
// Every method takes a context.Context first and returns an error last.
// Resources are grouped by service: client.Users, client.Roles,
// client.Webhooks and so on, one per area of the API.
//
// # Authentication
//
// Give NewClient one credential:
//
//   - WithAPIKey: an API key, sent as "Authorization: ApiKey <key>".
//   - WithClientCredentials: a service principal's client ID and secret. The
//     client obtains short-lived tokens with the OAuth client credentials grant
//     and renews them before they expire.
//   - WithTokenSource: any other bearer token, such as an agent's.
//
// Credentials belong on servers. Never ship an API key or a client secret to a
// browser or a mobile app.
//
// # Errors
//
// When the API answers with an error the method returns an *APIError, which
// carries the HTTP status, the stable error code, a message, the request ID
// and, for rate limiting and overload, how long to wait. The OAuth token
// endpoints answer in the RFC 6749 shape instead, returned as an *OAuthError.
// Both match the sentinel errors with errors.Is:
//
//	_, err := client.Users.Get(ctx, userID)
//	switch {
//	case errors.Is(err, qeetid.ErrNotFound):
//		// no such user
//	case errors.Is(err, qeetid.ErrRateLimited):
//		// slow down
//	}
//	if apiErr, ok := errors.AsType[*qeetid.APIError](err); ok {
//		log.Printf("qeet id error %s (request %s)", apiErr.Code, apiErr.RequestID)
//	}
//
// A failure to reach the API is returned as an error that wraps the cause, so
// errors.Is(err, context.DeadlineExceeded) and errors.As with *net.OpError
// work.
//
// # Retries and timeouts
//
// The client retries a failed request when that is safe, with exponential
// backoff and jitter: a GET or DELETE, or a create call carrying an
// idempotency key, after a network error or a 408, 500, 502, 503 or 504; and
// any request the server refused before acting on it, which it signals with
// 429 or 503 and a Retry-After header. It waits as long as Retry-After asks,
// up to MaxRetryAfter. It never retries a 401: a revoked session
// (ErrSessionRevoked) needs the person to sign in again. Each attempt has its
// own timeout (WithTimeout); the context bounds the whole call.
//
// # Idempotency
//
// The create calls that support it (Users.Create, Groups.Create,
// Invitations.Create, ServicePrincipals.Create, APIKeys.Create, Agents.Create,
// Webhooks.Create and Billing.Checkout, and Billing.PayInvoice) send an
// Idempotency-Key header. The client generates one per call so its own retries
// are safe; pass WithIdempotencyKey to make retries across process restarts
// safe as well.
//
// # Pagination
//
// List methods return one *Page. Cursor-paginated lists also have an All
// method that walks every page as a Go iterator:
//
//	for user, err := range client.Users.All(ctx, &qeetid.ListUsersParams{Limit: 100}) {
//		if err != nil {
//			return err
//		}
//		fmt.Println(user.Email)
//	}
//
// Lists the server returns whole report Page.Truncated when they hit the
// server's cap.
//
// # Related packages
//
// Package github.com/qeetgroup/qeet-id-go/webhook verifies and parses the
// webhook deliveries Qeet ID sends to your endpoints. Package
// github.com/qeetgroup/qeet-id-go/jwt verifies Qeet ID access tokens locally
// against the published signing keys, with HTTP middleware.
package qeetid
