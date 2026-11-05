package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type claimsKey struct{}

// ContextWithClaims returns a copy of ctx carrying claims. Middleware calls it;
// tests of handlers that read claims can call it too.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// ClaimsFromContext returns the claims Middleware stored in ctx.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey{}).(*Claims)
	return c, ok && c != nil
}

// Middleware returns a handler that verifies the bearer token in each
// request's Authorization header and calls next with the claims in the
// request's context (see ClaimsFromContext). It answers, as RFC 6750
// describes:
//
//   - 401 with WWW-Authenticate: Bearer when the request carries no token;
//   - 401 with error="invalid_token" when the token is refused;
//   - 503 when the signing keys cannot be fetched, so clients retry instead
//     of discarding a token that may be valid.
func (v *Verifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			w.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		claims, err := v.Verify(r.Context(), token)
		switch {
		case errors.Is(err, ErrKeysUnavailable):
			http.Error(w, "token verification unavailable", http.StatusServiceUnavailable)
			return
		case err != nil:
			w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ContextWithClaims(r.Context(), claims)))
	})
}

// RequireScope returns middleware that lets a request through only when the
// claims Middleware stored carry every one of scopes. Otherwise it answers 403
// with error="insufficient_scope" (RFC 6750), or 401 when there are no claims.
// Wrap it inside Middleware:
//
//	v.Middleware(jwt.RequireScope("invoices.read")(handler))
func RequireScope(scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				w.Header().Set("WWW-Authenticate", `Bearer`)
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			for _, s := range scopes {
				if !claims.HasScope(s) {
					w.Header().Set("WWW-Authenticate", `Bearer error="insufficient_scope", scope="`+strings.Join(scopes, " ")+`"`)
					http.Error(w, "insufficient scope", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// bearerToken extracts the token of an "Authorization: Bearer <token>" header.
// The scheme is case-insensitive.
func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(header), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	return token, token != ""
}
