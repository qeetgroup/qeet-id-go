package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testKey holds a generated ES256 keypair plus its JWKS-document encoding,
// for building signed test tokens without touching a real server.
type testKey struct {
	priv *ecdsa.PrivateKey
	kid  string
}

func newTestKey(t *testing.T, kid string) *testKey {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return &testKey{priv: priv, kid: kid}
}

func (k *testKey) jwk() map[string]string {
	return map[string]string{
		"kty": "EC",
		"crv": "P-256",
		"kid": k.kid,
		"x":   base64.RawURLEncoding.EncodeToString(k.priv.X.Bytes()),
		"y":   base64.RawURLEncoding.EncodeToString(k.priv.Y.Bytes()),
	}
}

func b64(v any) string {
	b, _ := json.Marshal(v)
	return base64.RawURLEncoding.EncodeToString(b)
}

// sign builds a compact ES256 JWT for the given claims, signed by k.
func (k *testKey) sign(t *testing.T, claims map[string]any) string {
	t.Helper()
	header := map[string]any{"alg": "ES256", "kid": k.kid}
	signingInput := b64(header) + "." + b64(claims)
	digest := sha256.Sum256([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, k.priv, digest[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func jwksServer(t *testing.T, keys ...*testKey) *httptest.Server {
	t.Helper()
	var jwks []map[string]string
	for _, k := range keys {
		jwks = append(jwks, k.jwk())
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": jwks})
	}))
}

func validClaims() map[string]any {
	now := time.Now().Unix()
	return map[string]any{
		"user_id":   "u1",
		"tenant_id": "t1",
		"sid":       "s1",
		"scope":     "read write",
		"sub":       "u1",
		"iss":       "https://id.example.com",
		"aud":       "https://api.example.com",
		"iat":       now,
		"exp":       now + 300,
	}
}

func TestVerify_ValidToken(t *testing.T) {
	k := newTestKey(t, "key1")
	srv := jwksServer(t, k)
	defer srv.Close()

	v := NewJWKSVerifier(srv.URL, http.DefaultClient)
	token := k.sign(t, validClaims())

	claims, err := v.Verify(context.Background(), token, VerifyOptions{})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.UserID != "u1" || claims.TenantID != "t1" || claims.SessionID != "s1" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestVerify_IssuerAudienceEnforced(t *testing.T) {
	k := newTestKey(t, "key1")
	srv := jwksServer(t, k)
	defer srv.Close()
	v := NewJWKSVerifier(srv.URL, http.DefaultClient)
	token := k.sign(t, validClaims())

	if _, err := v.Verify(context.Background(), token, VerifyOptions{Issuer: "https://wrong.example.com"}); err == nil {
		t.Fatal("wrong issuer accepted")
	}
	if _, err := v.Verify(context.Background(), token, VerifyOptions{Audience: "https://wrong.example.com"}); err == nil {
		t.Fatal("wrong audience accepted")
	}
	if _, err := v.Verify(context.Background(), token, VerifyOptions{
		Issuer: "https://id.example.com", Audience: "https://api.example.com",
	}); err != nil {
		t.Fatalf("correct issuer+audience rejected: %v", err)
	}
}

func TestVerify_Expired(t *testing.T) {
	k := newTestKey(t, "key1")
	srv := jwksServer(t, k)
	defer srv.Close()
	v := NewJWKSVerifier(srv.URL, http.DefaultClient)

	claims := validClaims()
	claims["exp"] = time.Now().Add(-time.Hour).Unix()
	token := k.sign(t, claims)

	if _, err := v.Verify(context.Background(), token, VerifyOptions{}); err == nil {
		t.Fatal("expired token accepted")
	}
}

func TestVerify_TamperedSignature(t *testing.T) {
	k := newTestKey(t, "key1")
	srv := jwksServer(t, k)
	defer srv.Close()
	v := NewJWKSVerifier(srv.URL, http.DefaultClient)

	token := k.sign(t, validClaims())
	tampered := token[:len(token)-4] + "AAAA"
	if _, err := v.Verify(context.Background(), tampered, VerifyOptions{}); err == nil {
		t.Fatal("tampered signature accepted")
	}
}

func TestVerify_DifferentKeyRejected(t *testing.T) {
	signer := newTestKey(t, "key1")
	other := newTestKey(t, "key1") // same kid, different keypair
	srv := jwksServer(t, other)    // JWKS only publishes "other"'s key
	defer srv.Close()
	v := NewJWKSVerifier(srv.URL, http.DefaultClient)

	token := signer.sign(t, validClaims())
	if _, err := v.Verify(context.Background(), token, VerifyOptions{}); err == nil {
		t.Fatal("token signed by an unpublished key was accepted")
	}
}

func TestVerify_UnknownKidForcesRefresh(t *testing.T) {
	k1 := newTestKey(t, "key1")
	k2 := newTestKey(t, "key2")
	srv := jwksServer(t, k1, k2)
	defer srv.Close()
	v := NewJWKSVerifier(srv.URL, http.DefaultClient)

	// Prime the cache with a request for key1, then verify a key2 token —
	// this must trigger a refresh rather than failing on a stale cache.
	if _, err := v.Verify(context.Background(), k1.sign(t, validClaims()), VerifyOptions{}); err != nil {
		t.Fatalf("priming Verify: %v", err)
	}
	if _, err := v.Verify(context.Background(), k2.sign(t, validClaims()), VerifyOptions{}); err != nil {
		t.Fatalf("key2 token rejected: %v", err)
	}
}

func TestVerify_MalformedToken(t *testing.T) {
	v := NewJWKSVerifier("http://unused.invalid", http.DefaultClient)
	if _, err := v.Verify(context.Background(), "not-a-jwt", VerifyOptions{}); err == nil {
		t.Fatal("malformed token accepted")
	}
}

func TestVerify_UnsupportedAlgRejected(t *testing.T) {
	v := NewJWKSVerifier("http://unused.invalid", http.DefaultClient)
	header := b64(map[string]any{"alg": "HS256", "kid": "x"})
	payload := b64(validClaims())
	token := header + "." + payload + ".sig"
	if _, err := v.Verify(context.Background(), token, VerifyOptions{}); err == nil {
		t.Fatal("HS256 token accepted — alg-confusion risk")
	}
}

func TestSetJWKSURL_DropsCache(t *testing.T) {
	k1 := newTestKey(t, "key1")
	srv1 := jwksServer(t, k1)
	defer srv1.Close()
	k2 := newTestKey(t, "key2")
	srv2 := jwksServer(t, k2)
	defer srv2.Close()

	v := NewJWKSVerifier(srv1.URL, http.DefaultClient)
	if _, err := v.Verify(context.Background(), k1.sign(t, validClaims()), VerifyOptions{}); err != nil {
		t.Fatalf("priming against srv1: %v", err)
	}

	v.SetJWKSURL(srv2.URL)
	if _, err := v.Verify(context.Background(), k2.sign(t, validClaims()), VerifyOptions{}); err != nil {
		t.Fatalf("key2 token against srv2 rejected after SetJWKSURL: %v", err)
	}
}
