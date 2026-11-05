package webhook_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/qeetgroup/qeet-id-go/webhook"
)

// A delivery as Qeet ID sends it, signed with the subscription's secret.
func signedRequest(secret, body string) *http.Request {
	now := time.Now()
	r := httptest.NewRequest(http.MethodPost, "/webhooks/qeet-id", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set(webhook.TimestampHeader, strconv.FormatInt(now.Unix(), 10))
	r.Header.Set(webhook.SignatureHeader, webhook.Sign(secret, now, []byte(body)))
	r.Header.Set(webhook.EventHeader, "user.created")
	return r
}

func ExampleVerifier_Handler() {
	secret := "whsec_example" // in production: os.Getenv("QEETID_WEBHOOK_SECRET")
	v, err := webhook.NewVerifier(secret)
	if err != nil {
		log.Fatal(err)
	}
	handler := v.Handler(func(_ context.Context, e *webhook.Event) error {
		switch e.Type {
		case webhook.UserCreated:
			var d webhook.UserData
			if err := e.DecodeData(&d); err != nil {
				return err
			}
			fmt.Println("provision", d.UserID, d.Email)
		case webhook.SessionRevoked:
			var d webhook.SessionRevokedData
			if err := e.DecodeData(&d); err != nil {
				return err
			}
			fmt.Println("forget session", d.SessionID)
		}
		return nil
	})

	body := `{"id":"evt-1","type":"user.created","tenant_id":"t1","occurred_at":"2026-11-05T10:00:00Z",` +
		`"data":{"user_id":"u1","tenant_id":"t1","email":"ada@example.com","source":"admin"}}`
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, signedRequest(secret, body))
	fmt.Println(w.Code)

	// A forged delivery is refused before the handler runs.
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, signedRequest("someone-else", body))
	fmt.Println(w.Code)
	// Output:
	// provision u1 ada@example.com
	// 204
	// 401
}

func ExampleVerifier_ConstructEvent() {
	secret := "whsec_example"
	body := `{"id":"evt-2","type":"test.ping","tenant_id":"t1","occurred_at":"2026-11-05T10:00:00Z","data":{"subscription_id":"w1"}}`
	r := signedRequest(secret, body)

	v, err := webhook.NewVerifier(secret)
	if err != nil {
		log.Fatal(err)
	}
	// Read the raw body yourself (capped) and pass it unchanged.
	payload := []byte(body)
	e, err := v.ConstructEvent(payload, r.Header)
	if err != nil {
		fmt.Println("refused:", err)
		return
	}
	fmt.Println(e.Type, e.ID)
	// Output: test.ping evt-2
}
