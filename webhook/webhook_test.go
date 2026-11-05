package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// delivery is a webhook delivery recorded from the Qeet ID server's own
// signing code (internal/developer/webhooks: sign and deliver), with the
// secrets it was signed with. See testdata/deliveries.json.
type delivery struct {
	Name    string            `json:"name"`
	Secrets []string          `json:"secrets"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func loadDeliveries(t *testing.T) map[string]delivery {
	t.Helper()
	b, err := os.ReadFile("testdata/deliveries.json")
	if err != nil {
		t.Fatal(err)
	}
	var list []delivery
	if err := json.Unmarshal(b, &list); err != nil {
		t.Fatal(err)
	}
	out := map[string]delivery{}
	for _, d := range list {
		out[d.Name] = d
	}
	return out
}

func (d delivery) header() http.Header {
	h := http.Header{}
	for k, v := range d.Headers {
		h.Set(k, v)
	}
	return h
}

// sentAt is the delivery's timestamp, used as the verifier's clock.
func (d delivery) sentAt(t *testing.T) time.Time {
	t.Helper()
	secs, err := strconv.ParseInt(d.Headers[TimestampHeader], 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return time.Unix(secs, 0)
}

func verifierAt(t *testing.T, secret string, now time.Time, opts ...Option) *Verifier {
	t.Helper()
	v, err := NewVerifier(secret, append([]Option{WithClock(func() time.Time { return now })}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestServerDeliveriesVerify(t *testing.T) {
	for name, d := range loadDeliveries(t) {
		t.Run(name, func(t *testing.T) {
			for _, secret := range d.Secrets {
				v := verifierAt(t, secret, d.sentAt(t).Add(30*time.Second))
				e, err := v.ConstructEvent([]byte(d.Body), d.header())
				if err != nil {
					t.Fatalf("secret %q: %v", secret[:4], err)
				}
				if string(e.Type) != d.Headers[EventHeader] || e.DeliveryID != d.Headers[DeliveryIDHeader] || e.ID == "" || e.TenantID == "" || e.OccurredAt.IsZero() {
					t.Errorf("event = %+v", e)
				}
			}
		})
	}
}

func TestServerDeliveryData(t *testing.T) {
	ds := loadDeliveries(t)

	d := ds["user_created"]
	e, err := verifierAt(t, d.Secrets[0], d.sentAt(t)).ConstructEvent([]byte(d.Body), d.header())
	if err != nil {
		t.Fatal(err)
	}
	var user UserData
	if err := e.DecodeData(&user); err != nil {
		t.Fatal(err)
	}
	if e.Type != UserCreated || user.Email != "ada@brightleaffoods.in" || user.Source != "admin" || user.UserID == "" || user.TenantID != e.TenantID {
		t.Errorf("user.created = %+v / %+v", e, user)
	}

	d = ds["session_revoked_rotation"]
	e, err = verifierAt(t, d.Secrets[1], d.sentAt(t)).ConstructEvent([]byte(d.Body), d.header())
	if err != nil {
		t.Fatal(err)
	}
	var revoked SessionRevokedData
	if err := e.DecodeData(&revoked); err != nil {
		t.Fatal(err)
	}
	if e.Type != SessionRevoked || revoked.Reason != ReasonAdminRevoked || revoked.SessionID == "" || revoked.RevokedAt.IsZero() {
		t.Errorf("session.revoked = %+v", revoked)
	}

	d = ds["test_ping"]
	e, err = verifierAt(t, d.Secrets[0], d.sentAt(t)).ConstructEvent([]byte(d.Body), d.header())
	if err != nil {
		t.Fatal(err)
	}
	var ping TestPingData
	if err := e.DecodeData(&ping); err != nil || e.Type != TestPing || ping.SubscriptionID == "" {
		t.Errorf("test.ping = %+v, %v", ping, err)
	}
}

// TestServerSignVector checks Sign against the vector of the server's own
// unit test for its signing function (TestSign_MatchesHMACSHA256Hex).
func TestServerSignVector(t *testing.T) {
	secret := "whsec_abc123"
	body := []byte(`{"event":"user.created","id":"01J"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("1700000000." + string(body)))
	want := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if got := Sign(secret, time.Unix(1700000000, 0), body); got != want {
		t.Fatalf("Sign = %s, want %s", got, want)
	}
}

func TestRotationAcceptsEitherSecret(t *testing.T) {
	d := loadDeliveries(t)["session_revoked_rotation"]
	if !strings.Contains(d.Headers[SignatureHeader], ",") {
		t.Fatal("fixture should carry two signatures")
	}
	for _, secret := range d.Secrets {
		if err := verifierAt(t, secret, d.sentAt(t)).Verify([]byte(d.Body), d.header()); err != nil {
			t.Errorf("secret %q refused: %v", secret[:4], err)
		}
	}
	// A receiver holding both secrets while it rolls the new one out.
	v := verifierAt(t, "unrelated", d.sentAt(t), WithSecret(d.Secrets[1]))
	if err := v.Verify([]byte(d.Body), d.header()); err != nil {
		t.Errorf("WithSecret not honoured: %v", err)
	}
}

func TestRefusals(t *testing.T) {
	d := loadDeliveries(t)["user_created"]
	secret, sent := d.Secrets[0], d.sentAt(t)
	cases := []struct {
		name   string
		mutate func(h http.Header, body []byte) []byte
		now    time.Time
		opts   []Option
		want   error
	}{
		{"wrong secret", nil, sent, nil, ErrNoValidSignature},
		{"tampered body", func(_ http.Header, b []byte) []byte {
			return bytes.Replace(b, []byte("ada@"), []byte("eve@"), 1)
		}, sent, nil, ErrNoValidSignature},
		{"re-encoded body", func(_ http.Header, b []byte) []byte {
			var v map[string]any
			_ = json.Unmarshal(b, &v)
			out, _ := json.MarshalIndent(v, "", " ")
			return out
		}, sent, nil, ErrNoValidSignature},
		{"timestamp changed", func(h http.Header, b []byte) []byte {
			h.Set(TimestampHeader, strconv.FormatInt(sent.Unix()+1, 10))
			return b
		}, sent, nil, ErrNoValidSignature},
		{"too old", nil, sent.Add(5*time.Minute + time.Second), nil, ErrTimestampOutOfTolerance},
		{"from the future", nil, sent.Add(-5*time.Minute - time.Second), nil, ErrTimestampOutOfTolerance},
		{"outside a tighter tolerance", nil, sent.Add(31 * time.Second), []Option{WithTolerance(30 * time.Second)}, ErrTimestampOutOfTolerance},
		{"no signature", func(h http.Header, b []byte) []byte { h.Del(SignatureHeader); return b }, sent, nil, ErrMissingSignature},
		{"signature without prefix", func(h http.Header, b []byte) []byte {
			h.Set(SignatureHeader, strings.TrimPrefix(h.Get(SignatureHeader), "sha256="))
			return b
		}, sent, nil, ErrMissingSignature},
		{"other algorithm", func(h http.Header, b []byte) []byte {
			h.Set(SignatureHeader, "sha1="+strings.Repeat("a", 40))
			return b
		}, sent, nil, ErrMissingSignature},
		{"no timestamp", func(h http.Header, b []byte) []byte { h.Del(TimestampHeader); return b }, sent, nil, ErrMissingTimestamp},
		{"malformed timestamp", func(h http.Header, b []byte) []byte { h.Set(TimestampHeader, "yesterday"); return b }, sent, nil, ErrMissingTimestamp},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := d.header()
			body := []byte(d.Body)
			if tc.mutate != nil {
				body = tc.mutate(h, body)
			}
			s := secret
			if tc.name == "wrong secret" {
				s = "not-the-secret"
			}
			err := verifierAt(t, s, tc.now, tc.opts...).Verify(body, h)
			if !errors.Is(err, tc.want) || !errors.Is(err, ErrInvalidDelivery) {
				t.Fatalf("Verify = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestWithinTolerance(t *testing.T) {
	d := loadDeliveries(t)["user_created"]
	for _, skew := range []time.Duration{-5 * time.Minute, 0, 5 * time.Minute} {
		if err := verifierAt(t, d.Secrets[0], d.sentAt(t).Add(skew)).Verify([]byte(d.Body), d.header()); err != nil {
			t.Errorf("skew %v refused: %v", skew, err)
		}
	}
}

func TestNewVerifierRequiresASecret(t *testing.T) {
	if _, err := NewVerifier(""); err == nil {
		t.Fatal("empty secret accepted")
	}
}

func TestMalformedPayload(t *testing.T) {
	now := time.Unix(1793858052, 0)
	body := []byte(`{"not":"an event"}`)
	h := http.Header{}
	h.Set(TimestampHeader, strconv.FormatInt(now.Unix(), 10))
	h.Set(SignatureHeader, Sign("s3cret", now, body))
	_, err := verifierAt(t, "s3cret", now).ConstructEvent(body, h)
	if !errors.Is(err, ErrMalformedPayload) {
		t.Fatalf("err = %v, want ErrMalformedPayload", err)
	}
}

func TestHandler(t *testing.T) {
	d := loadDeliveries(t)["user_created"]
	v := verifierAt(t, d.Secrets[0], d.sentAt(t))
	var got *Event
	failNext := false
	h := v.Handler(func(_ context.Context, e *Event) error {
		if failNext {
			return errors.New("database unavailable")
		}
		got = e
		return nil
	})
	send := func(method, body string, header http.Header) int {
		r := httptest.NewRequest(method, "/webhooks/qeet-id", strings.NewReader(body))
		for k, vs := range header {
			r.Header[k] = vs
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}

	if code := send(http.MethodPost, d.Body, d.header()); code != http.StatusNoContent || got == nil || got.Type != UserCreated {
		t.Fatalf("valid delivery: %d, event %+v", code, got)
	}
	if code := send(http.MethodPost, d.Body+" ", d.header()); code != http.StatusUnauthorized {
		t.Errorf("altered delivery: %d, want 401", code)
	}
	if code := send(http.MethodGet, "", d.header()); code != http.StatusMethodNotAllowed {
		t.Errorf("GET: %d, want 405", code)
	}
	failNext = true
	if code := send(http.MethodPost, d.Body, d.header()); code != http.StatusInternalServerError {
		t.Errorf("handler error: %d, want 500 so Qeet ID retries", code)
	}

	small := verifierAt(t, d.Secrets[0], d.sentAt(t), WithMaxBodyBytes(16)).Handler(func(context.Context, *Event) error { return nil })
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(d.Body))
	r.Header = d.header()
	w := httptest.NewRecorder()
	small.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("oversized body: %d, want 400", w.Code)
	}
}

func TestParseSignaturesIgnoresJunk(t *testing.T) {
	good := strings.Repeat("ab", 32)
	header := " sha256=" + good + " ,sha256=zz,sha256=" + good[:10] + ",v2=" + good + ",," + "sha256=" + good
	if n := len(parseSignatures(header)); n != 2 {
		t.Fatalf("parsed %d signatures, want 2", n)
	}
	many := strings.Repeat("sha256="+good+",", 50)
	if n := len(parseSignatures(many)); n != maxSignatures {
		t.Fatalf("parsed %d signatures, want the cap %d", n, maxSignatures)
	}
}
