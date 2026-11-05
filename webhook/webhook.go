// Package webhook verifies and parses the webhook deliveries Qeet ID sends to
// your endpoints.
//
// Each delivery is a POST whose JSON body is an envelope, decoded into Event:
//
//	{"id": "...", "type": "user.created", "tenant_id": "...", "occurred_at": "...", "data": {...}}
//
// Qeet ID signs every delivery. X-Qeet-Signature is "sha256=" and the hex
// HMAC-SHA256, keyed by the subscription's secret, of the X-Qeet-Timestamp
// header (Unix seconds), a ".", and the raw body. For 24 hours after the
// secret is rotated the header holds two comma-separated values, one per
// secret. A Verifier accepts a delivery when any value matches one of its
// secrets, compares in constant time, and refuses a timestamp more than five
// minutes away from its clock, so a captured delivery cannot be replayed
// later.
//
// The simplest use is the handler:
//
//	v, err := webhook.NewVerifier(os.Getenv("QEETID_WEBHOOK_SECRET"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	http.Handle("POST /webhooks/qeet-id", v.Handler(func(ctx context.Context, e *webhook.Event) error {
//		switch e.Type {
//		case webhook.UserCreated:
//			var d webhook.UserData
//			if err := e.DecodeData(&d); err != nil {
//				return err
//			}
//			// provision d.UserID ...
//		}
//		return nil // answers 204; an error answers 500 and Qeet ID retries
//	}))
//
// Deliveries are at least once and in no guaranteed order: deduplicate on
// Event.ID, which is the same in every delivery of one event.
package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Delivery headers set by Qeet ID.
const (
	// SignatureHeader carries one or two "sha256=<hex>" values, comma-separated.
	SignatureHeader = "X-Qeet-Signature"
	// TimestampHeader carries the time of sending in Unix seconds. It is part
	// of what is signed.
	TimestampHeader = "X-Qeet-Timestamp"
	// EventHeader carries the event type, the same as Event.Type.
	EventHeader = "X-Qeet-Event"
	// DeliveryIDHeader identifies one delivery attempt series of one event to
	// one subscription.
	DeliveryIDHeader = "X-Qeet-Delivery-ID"
)

const (
	// DefaultTolerance is how far a delivery's timestamp may be from the
	// verifier's clock, in either direction.
	DefaultTolerance = 5 * time.Minute
	// DefaultMaxBodyBytes caps the body Handler reads.
	DefaultMaxBodyBytes = 1 << 20

	signaturePrefix = "sha256="
	maxSignatures   = 8
)

// Errors returned by Verify and ConstructEvent. Every verification failure
// also matches ErrInvalidDelivery, so a handler can refuse them all with one
// errors.Is check.
var (
	// ErrInvalidDelivery matches every reason a delivery is refused.
	ErrInvalidDelivery = errors.New("webhook: invalid delivery")
	// ErrMissingSignature reports a delivery without an X-Qeet-Signature
	// header, or without a sha256 value in it.
	ErrMissingSignature = fmt.Errorf("%w: missing signature", ErrInvalidDelivery)
	// ErrMissingTimestamp reports a delivery without a valid X-Qeet-Timestamp
	// header.
	ErrMissingTimestamp = fmt.Errorf("%w: missing or malformed timestamp", ErrInvalidDelivery)
	// ErrTimestampOutOfTolerance reports a delivery whose timestamp is further
	// from the verifier's clock than its tolerance: a replay, or a clock that
	// is badly off.
	ErrTimestampOutOfTolerance = fmt.Errorf("%w: timestamp outside the tolerance", ErrInvalidDelivery)
	// ErrNoValidSignature reports a delivery none of whose signatures matches
	// any of the verifier's secrets: a forged or altered delivery, or the
	// wrong secret.
	ErrNoValidSignature = fmt.Errorf("%w: no signature matches", ErrInvalidDelivery)
	// ErrMalformedPayload reports a verified body that is not an event
	// envelope.
	ErrMalformedPayload = fmt.Errorf("%w: malformed payload", ErrInvalidDelivery)
)

// Verifier checks deliveries against one or more subscription secrets. It is
// safe for concurrent use.
type Verifier struct {
	secrets   [][]byte
	tolerance time.Duration
	maxBody   int64
	now       func() time.Time
}

// Option configures a Verifier.
type Option interface {
	apply(*Verifier)
}

type optionFunc func(*Verifier)

func (f optionFunc) apply(v *Verifier) { f(v) }

// WithTolerance sets how far a delivery's timestamp may be from the
// verifier's clock (default DefaultTolerance). A non-positive d leaves the
// default.
func WithTolerance(d time.Duration) Option {
	return optionFunc(func(v *Verifier) {
		if d > 0 {
			v.tolerance = d
		}
	})
}

// WithSecret adds a secret the verifier also accepts, for example while you
// roll a newly rotated secret out to every instance of your service.
func WithSecret(secret string) Option {
	return optionFunc(func(v *Verifier) {
		if secret != "" {
			v.secrets = append(v.secrets, []byte(secret))
		}
	})
}

// WithMaxBodyBytes caps the body Handler reads (default DefaultMaxBodyBytes).
func WithMaxBodyBytes(n int64) Option {
	return optionFunc(func(v *Verifier) {
		if n > 0 {
			v.maxBody = n
		}
	})
}

// WithClock sets the function the verifier reads the current time from, for
// tests that replay recorded deliveries.
func WithClock(now func() time.Time) Option {
	return optionFunc(func(v *Verifier) {
		if now != nil {
			v.now = now
		}
	})
}

// NewVerifier returns a Verifier for the subscription secret Qeet ID returned
// when the subscription was created or its secret was rotated.
func NewVerifier(secret string, opts ...Option) (*Verifier, error) {
	if secret == "" {
		return nil, errors.New("webhook: the secret is empty")
	}
	v := &Verifier{
		secrets:   [][]byte{[]byte(secret)},
		tolerance: DefaultTolerance,
		maxBody:   DefaultMaxBodyBytes,
		now:       time.Now,
	}
	for _, o := range opts {
		if o != nil {
			o.apply(v)
		}
	}
	return v, nil
}

// Verify checks that payload, the raw request body, was signed by Qeet ID
// with one of the verifier's secrets at a time within the tolerance. It reads
// the signature and timestamp from header. Verify against the bytes received;
// re-encoding the JSON changes them and breaks the signature.
func (v *Verifier) Verify(payload []byte, header http.Header) error {
	ts := strings.TrimSpace(header.Get(TimestampHeader))
	secs, err := strconv.ParseInt(ts, 10, 64)
	if err != nil || secs <= 0 {
		return ErrMissingTimestamp
	}
	signatures := parseSignatures(header.Get(SignatureHeader))
	if len(signatures) == 0 {
		return ErrMissingSignature
	}
	if age := v.now().Sub(time.Unix(secs, 0)); age > v.tolerance || age < -v.tolerance {
		return ErrTimestampOutOfTolerance
	}
	for _, secret := range v.secrets {
		expected := computeSignature(secret, ts, payload)
		for _, sig := range signatures {
			if hmac.Equal(sig, expected) {
				return nil
			}
		}
	}
	return ErrNoValidSignature
}

// ConstructEvent verifies a delivery and decodes its envelope.
func (v *Verifier) ConstructEvent(payload []byte, header http.Header) (*Event, error) {
	if err := v.Verify(payload, header); err != nil {
		return nil, err
	}
	var e Event
	if err := json.Unmarshal(payload, &e); err != nil || e.ID == "" || e.Type == "" {
		return nil, ErrMalformedPayload
	}
	e.DeliveryID = header.Get(DeliveryIDHeader)
	return &e, nil
}

// HandlerFunc processes a verified event. Returning an error makes Handler
// answer 500, so Qeet ID delivers the event again later.
type HandlerFunc func(ctx context.Context, e *Event) error

// Handler returns an http.Handler that reads a delivery, verifies it, decodes
// it and passes it to fn. It answers:
//
//   - 204 No Content when fn returns nil;
//   - 401 Unauthorized when the delivery fails verification;
//   - 400 Bad Request when the body is unreadable, too large or not an event;
//   - 405 Method Not Allowed for anything but POST;
//   - 500 Internal Server Error when fn returns an error.
//
// Qeet ID retries a 5xx, 408 or 429 with backoff and gives up on any other
// 4xx at once, so return an error from fn only for failures worth retrying.
// Answer quickly: hand slow work to a queue.
func (v *Verifier) Handler(fn HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, v.maxBody))
		if err != nil {
			http.Error(w, "unreadable body", http.StatusBadRequest)
			return
		}
		e, err := v.ConstructEvent(body, r.Header)
		switch {
		case errors.Is(err, ErrMalformedPayload):
			http.Error(w, "malformed event", http.StatusBadRequest)
			return
		case err != nil:
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
		if err := fn(r.Context(), e); err != nil {
			http.Error(w, "event not processed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// Sign returns the X-Qeet-Signature value Qeet ID would send for payload
// signed with secret at timestamp. Use it to build signed deliveries in tests
// of your own handlers; send the timestamp in TimestampHeader.
func Sign(secret string, timestamp time.Time, payload []byte) string {
	ts := strconv.FormatInt(timestamp.Unix(), 10)
	return signaturePrefix + hex.EncodeToString(computeSignature([]byte(secret), ts, payload))
}

// computeSignature is the server's signing function: HMAC-SHA256 over the
// timestamp, ".", and the body.
func computeSignature(secret []byte, timestamp string, payload []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	return mac.Sum(nil)
}

// parseSignatures decodes the sha256 values of a signature header, ignoring
// values with other prefixes and anything malformed.
func parseSignatures(header string) [][]byte {
	var out [][]byte
	for part := range strings.SplitSeq(header, ",") {
		if len(out) == maxSignatures {
			break
		}
		part = strings.TrimSpace(part)
		hexSig, ok := strings.CutPrefix(part, signaturePrefix)
		if !ok || len(hexSig) != 2*sha256.Size {
			continue
		}
		sig, err := hex.DecodeString(hexSig)
		if err != nil {
			continue
		}
		out = append(out, sig)
	}
	return out
}
