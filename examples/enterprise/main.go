// Command enterprise demonstrates the patterns a production deployment
// actually needs together: a session-verifying + permission-checking HTTP
// middleware, structured request logging via the Logger hook, and a
// configured timeout/retry budget — the "how do I wire this into my real
// service" example, not a single-call demo.
//
//	QEETID_API_KEY=qk_… go run ./examples/enterprise
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/qeetgroup/qeet-id-go"
)

// requestLogger implements qeetid.Logger with the standard library's log
// package — swap in your structured logger of choice (slog, zap, zerolog)
// without qeet-id-go ever depending on one.
type requestLogger struct{}

func (requestLogger) LogRequest(method, path string, status int, duration time.Duration, requestID string) {
	log.Printf("qeetid %s %s -> %d in %s (request %s)", method, path, status, duration, requestID)
}

type contextKey string

const claimsKey contextKey = "qeetid-claims"

// requireAuth verifies the bearer token and checks a permission before
// calling next — the shape of a typical enterprise auth middleware.
func requireAuth(client *qeetid.Client, permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := client.Sessions.Verify(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ok, err := client.Permissions.Check(r.Context(), qeetid.PermissionCheck{
			User: claims.UserID, Tenant: claims.TenantID, Permission: permission,
		})
		if err != nil || !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	client := qeetid.New(qeetid.Config{
		APIKey:     os.Getenv("QEETID_API_KEY"),
		Timeout:    5 * time.Second,
		MaxRetries: 3,
		UserAgent:  "acme-api/1.0.0",
		Logger:     requestLogger{},
	})

	billingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(claimsKey).(*qeetid.Claims)
		w.Write([]byte("billing data for tenant " + claims.TenantID))
	})

	http.Handle("/billing", requireAuth(client, "billing:read", billingHandler))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
