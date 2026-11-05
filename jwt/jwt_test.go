package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testIssuer   = "qeet-id"
	testAudience = "https://api.brightleaffoods.in"
)

// signer holds an ES256 key published under kid, as the server publishes it.
type signer struct {
	kid string
	key *ecdsa.PrivateKey
}

func newSigner(t *testing.T, kid string) *signer {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return &signer{kid: kid, key: k}
}

// jwk is the key as /.well-known/jwks.json publishes it.
func (s *signer) jwk() map[string]string {
	raw, err := s.key.PublicKey.Bytes() // 0x04 || X || Y
	if err != nil {
		panic(err)
	}
	return map[string]string{
		"kty": "EC", "crv": "P-256", "use": "sig", "alg": "ES256", "kid": s.kid,
		"x": base64.RawURLEncoding.EncodeToString(raw[1:33]),
		"y": base64.RawURLEncoding.EncodeToString(raw[33:]),
	}
}

func b64(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// sign builds a compact JWS the way the server's issuer does: ES256 with the
// raw r||s signature and the kid header.
func (s *signer) sign(t *testing.T, header map[string]any, claims map[string]any) string {
	t.Helper()
	if header == nil {
		header = map[string]any{"alg": "ES256", "typ": "JWT", "kid": s.kid}
	}
	input := b64(header) + "." + b64(claims)
	digest := sha256.Sum256([]byte(input))
	der, err := ecdsa.SignASN1(rand.Reader, s.key, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	var rs struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(der, &rs); err != nil {
		t.Fatal(err)
	}
	sig := make([]byte, 64)
	rs.R.FillBytes(sig[:32])
	rs.S.FillBytes(sig[32:])
	return input + "." + base64.RawURLEncoding.EncodeToString(sig)
}

// accessClaims are the claims of a person's access token as the server issues
// them (platform/crypto/tokens Claims).
func accessClaims(now time.Time) map[string]any {
	return map[string]any{
		"iss": testIssuer, "aud": []string{testAudience}, "sub": "u-1",
		"user_id": "u-1", "tenant_id": "t-1", "sid": "s-1", "scope": "openid invoices.read",
		"iat": now.Unix(), "nbf": now.Unix(), "exp": now.Add(10 * time.Minute).Unix(), "jti": "j-1",
		"claims": map[string]any{"plan": "pro"},
	}
}

// jwksServer publishes the signers' keys and counts the fetches.
type jwksServer struct {
	*httptest.Server
	mu      sync.Mutex
	signers []*signer
	fetches atomic.Int32
	down    atomic.Bool
}

func newJWKSServer(t *testing.T, signers ...*signer) *jwksServer {
	t.Helper()
	js := &jwksServer{signers: signers}
	js.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		js.fetches.Add(1)
		if js.down.Load() {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		js.mu.Lock()
		keys := make([]map[string]string, 0, len(js.signers))
		for _, s := range js.signers {
			keys = append(keys, s.jwk())
		}
		js.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": keys})
	}))
	t.Cleanup(js.Close)
	return js
}

func (js *jwksServer) publish(signers ...*signer) {
	js.mu.Lock()
	js.signers = signers
	js.mu.Unlock()
}

type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) Now() time.Time      { c.mu.Lock(); defer c.mu.Unlock(); return c.now }
func (c *clock) add(d time.Duration) { c.mu.Lock(); c.now = c.now.Add(d); c.mu.Unlock() }

func newTestVerifier(t *testing.T, js *jwksServer, clk *clock, opts ...Option) *Verifier {
	t.Helper()
	base := []Option{WithJWKSURL(js.URL + "/.well-known/jwks.json"), WithClock(clk.Now)}
	v, err := NewVerifier(testIssuer, testAudience, append(base, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestVerifyAccessToken(t *testing.T) {
	s := newSigner(t, "kid-1")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)

	c, err := v.Verify(context.Background(), s.sign(t, nil, accessClaims(clk.Now())))
	if err != nil {
		t.Fatal(err)
	}
	if c.UserID != "u-1" || c.TenantID != "t-1" || c.SessionID != "s-1" || c.Issuer != testIssuer ||
		!c.HasScope("invoices.read") || c.HasScope("invoices.write") || c.Custom["plan"] != "pro" || c.ID != "j-1" {
		t.Errorf("claims = %+v", c)
	}
	if c.ExpiresAt.Unix() != clk.Now().Add(10*time.Minute).Unix() {
		t.Errorf("ExpiresAt = %v", c.ExpiresAt)
	}
	// The second verification uses the cached keys.
	if _, err := v.Verify(context.Background(), s.sign(t, nil, accessClaims(clk.Now()))); err != nil {
		t.Fatal(err)
	}
	if n := js.fetches.Load(); n != 1 {
		t.Errorf("fetched the JWKS %d times, want 1", n)
	}
}

func TestServiceAndDelegatedTokens(t *testing.T) {
	s := newSigner(t, "kid-1")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)

	svc := map[string]any{"iss": testIssuer, "aud": testAudience, "sub": "sp-1", "tenant_id": "t-1",
		"scope": "user.read", "actor_type": "service", "exp": clk.Now().Add(time.Minute).Unix()}
	c, err := v.Verify(context.Background(), s.sign(t, nil, svc))
	if err != nil {
		t.Fatal(err)
	}
	if c.ActorType != "service" || c.Subject != "sp-1" || c.UserID != "" || len(c.Audience) != 1 {
		t.Errorf("service claims = %+v", c)
	}

	delegated := accessClaims(clk.Now())
	delegated["act"] = map[string]any{"sub": "agent-7"}
	delegated["exchange_depth"] = 1
	delegated["client_id"] = "app-1"
	c, err = v.Verify(context.Background(), s.sign(t, nil, delegated))
	if err != nil {
		t.Fatal(err)
	}
	if c.Actor == nil || c.Actor.Subject != "agent-7" || c.ExchangeDepth != 1 || c.ClientID != "app-1" {
		t.Errorf("delegated claims = %+v", c)
	}
}

func TestVerifyRefusals(t *testing.T) {
	s := newSigner(t, "kid-1")
	other := newSigner(t, "kid-1") // same kid, different key: a forgery
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)
	now := clk.Now()

	with := func(mut func(map[string]any)) map[string]any {
		c := accessClaims(now)
		mut(c)
		return c
	}
	good := s.sign(t, nil, accessClaims(now))
	parts := strings.Split(good, ".")
	hmacToken := func() string {
		h := b64(map[string]any{"alg": "HS256", "kid": "kid-1"}) + "." + parts[1]
		mac := hmac.New(sha256.New, []byte("public key bytes"))
		mac.Write([]byte(h))
		return h + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	}

	cases := []struct {
		name  string
		token string
		want  error
	}{
		{"expired", s.sign(t, nil, with(func(c map[string]any) { c["exp"] = now.Add(-time.Minute).Unix() })), ErrExpired},
		{"no exp", s.sign(t, nil, with(func(c map[string]any) { delete(c, "exp") })), ErrExpired},
		{"not yet valid", s.sign(t, nil, with(func(c map[string]any) { c["nbf"] = now.Add(time.Minute).Unix() })), ErrNotYetValid},
		{"wrong issuer", s.sign(t, nil, with(func(c map[string]any) { c["iss"] = "https://evil.example" })), ErrIssuerMismatch},
		{"wrong audience", s.sign(t, nil, with(func(c map[string]any) { c["aud"] = "qeet-id" })), ErrAudienceMismatch},
		{"no audience", s.sign(t, nil, with(func(c map[string]any) { delete(c, "aud") })), ErrAudienceMismatch},
		{"alg none", b64(map[string]any{"alg": "none", "kid": "kid-1"}) + "." + parts[1] + ".", ErrUnsupportedAlgorithm},
		{"alg HS256", hmacToken(), ErrUnsupportedAlgorithm},
		{"alg RS256", s.sign(t, map[string]any{"alg": "RS256", "kid": "kid-1"}, accessClaims(now)), ErrUnsupportedAlgorithm},
		{"missing kid", s.sign(t, map[string]any{"alg": "ES256"}, accessClaims(now)), ErrUnknownKey},
		{"critical header", s.sign(t, map[string]any{"alg": "ES256", "kid": "kid-1", "crit": []string{"b64"}}, accessClaims(now)), ErrMalformed},
		{"forged signature", other.sign(t, nil, accessClaims(now)), ErrInvalidSignature},
		{"tampered payload", parts[0] + "." + b64(with(func(c map[string]any) { c["user_id"] = "admin" })) + "." + parts[2], ErrInvalidSignature},
		{"short signature", parts[0] + "." + parts[1] + "." + parts[2][:20], ErrInvalidSignature},
		{"two parts", parts[0] + "." + parts[1], ErrMalformed},
		{"four parts", good + ".x", ErrMalformed},
		{"garbage", "not-a-token", ErrMalformed},
		{"huge", strings.Repeat("a", maxTokenBytes+1), ErrMalformed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := v.Verify(context.Background(), tc.token)
			if !errors.Is(err, tc.want) || !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("Verify = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestLeeway(t *testing.T) {
	s := newSigner(t, "kid-1")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk, WithLeeway(time.Minute))
	c := accessClaims(clk.Now())
	c["exp"] = clk.Now().Add(-50 * time.Second).Unix()
	if _, err := v.Verify(context.Background(), s.sign(t, nil, c)); err != nil {
		t.Fatalf("token expired within the leeway refused: %v", err)
	}
	c["exp"] = clk.Now().Add(-70 * time.Second).Unix()
	if _, err := v.Verify(context.Background(), s.sign(t, nil, c)); !errors.Is(err, ErrExpired) {
		t.Fatalf("token expired beyond the leeway: %v", err)
	}
}

func TestKeyRotation(t *testing.T) {
	oldKey, newKey := newSigner(t, "kid-old"), newSigner(t, "kid-new")
	js := newJWKSServer(t, oldKey)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)
	ctx := context.Background()

	if _, err := v.Verify(ctx, oldKey.sign(t, nil, accessClaims(clk.Now()))); err != nil {
		t.Fatal(err)
	}
	// The server activates a new key and keeps the old one for its grace
	// period. A token naming the new kid makes the verifier fetch again.
	js.publish(newKey, oldKey)
	clk.add(2 * time.Minute)
	if _, err := v.Verify(ctx, newKey.sign(t, nil, accessClaims(clk.Now()))); err != nil {
		t.Fatalf("rotated-in key refused: %v", err)
	}
	if n := js.fetches.Load(); n != 2 {
		t.Fatalf("fetches = %d, want 2", n)
	}
	// Once the old key is retired and the cache expires, its tokens fail.
	js.publish(newKey)
	clk.add(DefaultKeyCacheTTL)
	if _, err := v.Verify(ctx, oldKey.sign(t, nil, accessClaims(clk.Now()))); !errors.Is(err, ErrUnknownKey) {
		t.Fatalf("retired key: %v, want ErrUnknownKey", err)
	}
}

func TestUnknownKidCannotHammerTheJWKS(t *testing.T) {
	s, stranger := newSigner(t, "kid-1"), newSigner(t, "kid-random")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)
	ctx := context.Background()
	if _, err := v.Verify(ctx, s.sign(t, nil, accessClaims(clk.Now()))); err != nil {
		t.Fatal(err)
	}
	for range 50 {
		if _, err := v.Verify(ctx, stranger.sign(t, nil, accessClaims(clk.Now()))); !errors.Is(err, ErrUnknownKey) {
			t.Fatalf("err = %v", err)
		}
	}
	if n := js.fetches.Load(); n > 2 {
		t.Fatalf("50 unknown kids caused %d fetches, want at most 2", n)
	}
}

func TestKeysUnavailable(t *testing.T) {
	s := newSigner(t, "kid-1")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)
	ctx := context.Background()

	js.down.Store(true)
	_, err := v.Verify(ctx, s.sign(t, nil, accessClaims(clk.Now())))
	if !errors.Is(err, ErrKeysUnavailable) || errors.Is(err, ErrInvalidToken) {
		t.Fatalf("err = %v, want ErrKeysUnavailable only", err)
	}

	// Once keys are cached, an outage at refresh time serves the stale keys.
	js.down.Store(false)
	clk.add(2 * time.Minute)
	if _, err := v.Verify(ctx, s.sign(t, nil, accessClaims(clk.Now()))); err != nil {
		t.Fatal(err)
	}
	js.down.Store(true)
	clk.add(DefaultKeyCacheTTL + time.Minute)
	if _, err := v.Verify(ctx, s.sign(t, nil, accessClaims(clk.Now()))); err != nil {
		t.Fatalf("stale keys not used during an outage: %v", err)
	}
}

func TestConcurrentFirstUseFetchesOnce(t *testing.T) {
	s := newSigner(t, "kid-1")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)
	token := s.sign(t, nil, accessClaims(clk.Now()))
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			if _, err := v.Verify(context.Background(), token); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if n := js.fetches.Load(); n != 1 {
		t.Fatalf("fetches = %d, want 1", n)
	}
}

func TestNewVerifierChecks(t *testing.T) {
	if _, err := NewVerifier("", testAudience); err == nil {
		t.Error("empty issuer accepted")
	}
	if _, err := NewVerifier(testIssuer, ""); err == nil {
		t.Error("empty audience accepted")
	}
	if _, err := NewVerifier(testIssuer, testAudience, WithJWKSURL("http://keys.example.com/jwks.json")); err == nil {
		t.Error("plain http JWKS URL accepted")
	}
	if _, err := NewVerifier(testIssuer, testAudience); err != nil {
		t.Errorf("defaults refused: %v", err)
	}
}

func TestJWKRejectsOtherKeys(t *testing.T) {
	s := newSigner(t, "k")
	j := s.jwk()
	for _, mut := range []func(map[string]string){
		func(m map[string]string) { m["crv"] = "P-384" },
		func(m map[string]string) { m["kty"] = "RSA" },
		func(m map[string]string) { m["alg"] = "ES384" },
		func(m map[string]string) { m["use"] = "enc" },
		func(m map[string]string) { m["x"] = "!!" },
	} {
		m := map[string]string{}
		for k, v := range j {
			m[k] = v
		}
		mut(m)
		k := jwk{Kty: m["kty"], Crv: m["crv"], X: m["x"], Y: m["y"], Kid: m["kid"], Use: m["use"], Alg: m["alg"]}
		if _, err := k.publicKey(); err == nil {
			t.Errorf("accepted %v", m)
		}
	}
}

func TestMiddleware(t *testing.T) {
	s := newSigner(t, "kid-1")
	js := newJWKSServer(t, s)
	clk := &clock{now: time.Now()}
	v := newTestVerifier(t, js, clk)
	h := v.Middleware(RequireScope("invoices.read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := ClaimsFromContext(r.Context())
		_, _ = w.Write([]byte(c.UserID))
	})))
	call := func(auth string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, "/invoices", nil)
		if auth != "" {
			r.Header.Set("Authorization", auth)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}

	if w := call("Bearer " + s.sign(t, nil, accessClaims(clk.Now()))); w.Code != http.StatusOK || w.Body.String() != "u-1" {
		t.Errorf("valid: %d %q", w.Code, w.Body.String())
	}
	if w := call("bearer " + s.sign(t, nil, accessClaims(clk.Now()))); w.Code != http.StatusOK {
		t.Errorf("lower-case scheme: %d", w.Code)
	}
	if w := call(""); w.Code != http.StatusUnauthorized || w.Header().Get("WWW-Authenticate") != "Bearer" {
		t.Errorf("missing: %d %q", w.Code, w.Header().Get("WWW-Authenticate"))
	}
	if w := call("Basic dXNlcjpwYXNz"); w.Code != http.StatusUnauthorized {
		t.Errorf("basic: %d", w.Code)
	}
	if w := call("Bearer garbage"); w.Code != http.StatusUnauthorized || !strings.Contains(w.Header().Get("WWW-Authenticate"), "invalid_token") {
		t.Errorf("invalid: %d %q", w.Code, w.Header().Get("WWW-Authenticate"))
	}
	noScope := accessClaims(clk.Now())
	noScope["scope"] = "openid"
	if w := call("Bearer " + s.sign(t, nil, noScope)); w.Code != http.StatusForbidden || !strings.Contains(w.Header().Get("WWW-Authenticate"), "insufficient_scope") {
		t.Errorf("insufficient scope: %d %q", w.Code, w.Header().Get("WWW-Authenticate"))
	}
	js.down.Store(true)
	fresh := newTestVerifier(t, js, clk)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+s.sign(t, nil, accessClaims(clk.Now())))
	fresh.Middleware(http.NotFoundHandler()).ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("keys down: %d, want 503", w.Code)
	}
}
