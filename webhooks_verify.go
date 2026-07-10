package qeetid

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Webhook delivery headers set by the Qeet ID dispatcher.
const (
	// WebhookSignatureHeader carries "sha256=<hex>" — HMAC-SHA256 of the raw
	// body keyed by the subscription's signing secret.
	WebhookSignatureHeader = "X-Qeet-Signature"
	// WebhookEventHeader carries the event type (e.g. "user.created").
	WebhookEventHeader = "X-Qeet-Event"

	signaturePrefix = "sha256="
	maxWebhookBytes = 1 << 20 // 1 MiB cap on inbound webhook bodies
)

// WebhookEvent is a verified inbound webhook delivery. Payload is the exact
// bytes whose signature was verified — decode it with Data.
type WebhookEvent struct {
	Type    string          // the X-Qeet-Event header value
	Payload json.RawMessage // the raw, verified request body
}

// Data unmarshals the verified payload into v.
func (e *WebhookEvent) Data(v any) error { return json.Unmarshal(e.Payload, v) }

// VerifyWebhookSignature recomputes HMAC-SHA256(secret, payload) and compares
// it, in constant time, against the "sha256=<hex>" value from the
// X-Qeet-Signature header. It returns nil when the signature is valid. Always
// verify against the RAW request bytes — never a re-serialized body.
func VerifyWebhookSignature(payload []byte, signatureHeader, secret string) error {
	if secret == "" {
		return &Error{Code: "invalid_signature", Message: "webhook secret is empty"}
	}
	if !strings.HasPrefix(signatureHeader, signaturePrefix) {
		return &Error{Code: "invalid_signature", Message: "signature header missing sha256= prefix"}
	}
	got := strings.TrimPrefix(signatureHeader, signaturePrefix)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	want := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(got), []byte(want)) {
		return &Error{Code: "invalid_signature", Message: "webhook signature mismatch"}
	}
	return nil
}

// ConstructEvent verifies an inbound webhook and returns the parsed event.
// Pass the raw request body, the X-Qeet-Signature and X-Qeet-Event header
// values, and the subscription's signing secret (shown once at create time).
func ConstructEvent(payload []byte, signatureHeader, eventHeader, secret string) (*WebhookEvent, error) {
	if err := VerifyWebhookSignature(payload, signatureHeader, secret); err != nil {
		return nil, err
	}
	// Copy so callers can reuse the input buffer without mutating the event.
	buf := make(json.RawMessage, len(payload))
	copy(buf, payload)
	return &WebhookEvent{Type: eventHeader, Payload: buf}, nil
}

// ConstructEventFromRequest reads, verifies, and parses a webhook straight from
// an *http.Request — the ergonomic entry point for an HTTP handler. It consumes
// r.Body (capped at 1 MiB) and reads the signature/event headers for you.
func ConstructEventFromRequest(r *http.Request, secret string) (*WebhookEvent, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxWebhookBytes))
	if err != nil {
		return nil, &Error{Code: "read_error", Message: "reading webhook body: " + err.Error()}
	}
	return ConstructEvent(body, r.Header.Get(WebhookSignatureHeader), r.Header.Get(WebhookEventHeader), secret)
}
