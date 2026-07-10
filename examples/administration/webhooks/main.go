// Command webhooks verifies inbound Qeet ID webhooks and prints them.
//
//	QEETID_WEBHOOK_SECRET=whsec_… go run ./examples/administration/webhooks
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qeetgroup/qeet-id-go"
)

func main() {
	secret := os.Getenv("QEETID_WEBHOOK_SECRET")

	http.HandleFunc("/webhooks/qeet", func(w http.ResponseWriter, r *http.Request) {
		event, err := qeetid.ConstructEventFromRequest(r, secret)
		if err != nil {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}
		switch event.Type {
		case "user.created", "user.deleted":
			log.Printf("%s: %s", event.Type, event.Payload)
		default:
			log.Printf("unhandled event %s", event.Type)
		}
		w.WriteHeader(http.StatusOK)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
