package qeetid

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"id":"evt_1","type":"user.created"}`)
	good := sign(secret, body)

	if err := VerifyWebhookSignature(body, good, secret); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}
	if err := VerifyWebhookSignature(body, good, "wrong-secret"); err == nil {
		t.Fatal("wrong secret accepted")
	}
	if err := VerifyWebhookSignature(body, good, ""); err == nil {
		t.Fatal("empty secret accepted")
	}
	noPrefix := good[len("sha256="):]
	if err := VerifyWebhookSignature(body, noPrefix, secret); err == nil {
		t.Fatal("signature without sha256= prefix accepted")
	}
	if err := VerifyWebhookSignature([]byte(`{"id":"evt_1","type":"user.deleted"}`), good, secret); err == nil {
		t.Fatal("tampered body accepted")
	}
}

func TestConstructEvent(t *testing.T) {
	secret := "whsec_construct"
	body := []byte(`{"id":"evt_9","data":{"user_id":"u1"}}`)
	ev, err := ConstructEvent(body, sign(secret, body), "user.created", secret)
	if err != nil {
		t.Fatalf("ConstructEvent: %v", err)
	}
	if ev.Type != "user.created" {
		t.Fatalf("Type = %q", ev.Type)
	}
	var parsed struct {
		ID   string `json:"id"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	if err := ev.Data(&parsed); err != nil {
		t.Fatalf("Data: %v", err)
	}
	if parsed.ID != "evt_9" || parsed.Data.UserID != "u1" {
		t.Fatalf("unexpected payload: %+v", parsed)
	}
}

func TestConstructEventFromRequest(t *testing.T) {
	secret := "whsec_req"
	body := []byte(`{"hello":"world"}`)
	req := httptest.NewRequest(http.MethodPost, "/hook", bytes.NewReader(body))
	req.Header.Set(WebhookSignatureHeader, sign(secret, body))
	req.Header.Set(WebhookEventHeader, "test.ping")

	ev, err := ConstructEventFromRequest(req, secret)
	if err != nil {
		t.Fatalf("ConstructEventFromRequest: %v", err)
	}
	if ev.Type != "test.ping" {
		t.Fatalf("Type = %q", ev.Type)
	}
	if string(ev.Payload) != string(body) {
		t.Fatalf("Payload = %s", ev.Payload)
	}
}
