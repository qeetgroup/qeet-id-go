package qeetid

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

const discoveryJSON = `{
	"issuer":"https://id.example.com",
	"authorization_endpoint":"https://id.example.com/v1/oauth/authorize",
	"token_endpoint":"https://id.example.com/v1/oauth/token-code",
	"userinfo_endpoint":"https://id.example.com/v1/oauth/userinfo",
	"jwks_uri":"https://id.example.com/keys/jwks.json",
	"introspection_endpoint":"https://id.example.com/oauth/introspect",
	"id_token_signing_alg_values_supported":["ES256"],
	"code_challenge_methods_supported":["S256"],
	"actor_types_supported":["user","service","agent"],
	"resource_indicators_supported":true
}`

func TestDiscover(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != discoveryPath {
			t.Errorf("path = %q, want %q", r.URL.Path, discoveryPath)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(discoveryJSON))
	}))
	defer srv.Close()

	doc, err := Discover(context.Background(), srv.URL, nil)
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if doc.Issuer != "https://id.example.com" {
		t.Fatalf("Issuer = %q", doc.Issuer)
	}
	if doc.JWKSURI != "https://id.example.com/keys/jwks.json" {
		t.Fatalf("JWKSURI = %q", doc.JWKSURI)
	}
	if len(doc.ActorTypesSupported) != 3 {
		t.Fatalf("ActorTypesSupported = %v", doc.ActorTypesSupported)
	}
	if !doc.ResourceIndicatorsSupported {
		t.Fatal("ResourceIndicatorsSupported should be true")
	}
	if doc.Raw["issuer"] != "https://id.example.com" {
		t.Fatal("Raw should retain the full document")
	}
}

// TestNewFromDiscovery_WiresJWKSURL proves the custom jwks_uri from
// discovery is actually used — not just that it's non-nil — by serving a
// real JWKS at a path distinct from the default /.well-known/jwks.json and
// verifying a token signed by that key round-trips through
// Sessions.Verify.
func TestNewFromDiscovery_WiresJWKSURL(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	kid := "custom-key"

	mux := http.NewServeMux()
	var baseURL string
	mux.HandleFunc(discoveryPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"issuer":"x","jwks_uri":"` + baseURL + `/custom/jwks.json"}`))
	})
	mux.HandleFunc("/custom/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kty": "EC", "crv": "P-256", "kid": kid,
			"x": base64.RawURLEncoding.EncodeToString(priv.X.Bytes()),
			"y": base64.RawURLEncoding.EncodeToString(priv.Y.Bytes()),
		}}})
	})
	// If a Verify call hits the DEFAULT /.well-known/jwks.json instead, that
	// means the rewiring silently failed — fail loudly instead of 404ing.
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		t.Error("Sessions.Verify hit the default JWKS path — discovery rewiring did not take effect")
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	baseURL = srv.URL

	c, doc, err := NewFromDiscovery(context.Background(), Config{APIKey: "qk_x.y", BaseURL: srv.URL})
	if err != nil {
		t.Fatalf("NewFromDiscovery: %v", err)
	}
	if doc.JWKSURI != srv.URL+"/custom/jwks.json" {
		t.Fatalf("doc.JWKSURI = %q", doc.JWKSURI)
	}

	now := time.Now().Unix()
	claims := map[string]any{"user_id": "u1", "iat": now, "exp": now + 300}
	token := signES256(t, priv, kid, claims)

	got, err := c.Sessions.Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("Verify after discovery rewiring: %v", err)
	}
	if got.UserID != "u1" {
		t.Fatalf("UserID = %q", got.UserID)
	}
}

func signES256(t *testing.T, priv *ecdsa.PrivateKey, kid string, claims map[string]any) string {
	t.Helper()
	b64 := func(v any) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	signingInput := b64(map[string]any{"alg": "ES256", "kid": kid}) + "." + b64(claims)
	digest := sha256.Sum256([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
}
