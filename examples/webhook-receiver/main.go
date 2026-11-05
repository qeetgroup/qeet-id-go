// Command webhook-receiver receives Qeet ID webhook deliveries, verifies their
// signatures and handles user and session events.
//
//	QEETID_WEBHOOK_SECRET=... go run ./examples/webhook-receiver
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/qeetgroup/qeet-id-go/webhook"
)

func main() {
	v, err := webhook.NewVerifier(os.Getenv("QEETID_WEBHOOK_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	// Deliveries are at least once: remember the events already handled.
	var mu sync.Mutex
	seen := map[string]bool{}

	mux := http.NewServeMux()
	mux.Handle("POST /webhooks/qeet-id", v.Handler(func(_ context.Context, e *webhook.Event) error {
		mu.Lock()
		dup := seen[e.ID]
		seen[e.ID] = true
		mu.Unlock()
		if dup {
			return nil
		}
		switch e.Type {
		case webhook.UserCreated, webhook.UserDeleted, webhook.UserSuspended:
			var d webhook.UserData
			if err := e.DecodeData(&d); err != nil {
				return err
			}
			log.Printf("%s: user %s in tenant %s", e.Type, d.UserID, d.TenantID)
		case webhook.SessionRevoked:
			var d webhook.SessionRevokedData
			if err := e.DecodeData(&d); err != nil {
				return err
			}
			log.Printf("session %s of user %s ended (%s)", d.SessionID, d.UserID, d.Reason)
		case webhook.TestPing:
			log.Print("test delivery received")
		}
		return nil
	}))

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
