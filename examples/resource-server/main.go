// Command resource-server protects an HTTP API with Qeet ID access tokens,
// verified locally against Qeet ID's published signing keys.
//
//	QEETID_AUDIENCE=https://api.example.com go run ./examples/resource-server
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	qeetid "github.com/qeetgroup/qeet-id-go"
	"github.com/qeetgroup/qeet-id-go/jwt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// The issuer is in Qeet ID's discovery document; read it once at start-up.
	client, err := qeetid.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	meta, err := client.OAuth.Discover(ctx)
	if err != nil {
		log.Fatal(err)
	}

	v, err := jwt.NewVerifier(meta.Issuer, os.Getenv("QEETID_AUDIENCE"), jwt.WithJWKSURL(meta.JWKSURI))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /invoices", v.Middleware(jwt.RequireScope("invoices.read")(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			claims, _ := jwt.ClaimsFromContext(r.Context())
			fmt.Fprintf(w, "invoices of tenant %s for user %s\n", claims.TenantID, claims.UserID)
		}))))

	srv := &http.Server{Addr: ":8081", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
