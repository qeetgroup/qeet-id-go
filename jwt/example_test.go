package jwt_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/qeetgroup/qeet-id-go/jwt"
)

func ExampleVerifier_Middleware() {
	// In production the verifier fetches Qeet ID's published keys; here a
	// local server stands in for them.
	jwks, sign := fakeQeetID()
	defer jwks.Close()

	v, err := jwt.NewVerifier("qeet-id", "https://api.example.com",
		jwt.WithJWKSURL(jwks.URL+"/.well-known/jwks.json"))
	if err != nil {
		log.Fatal(err)
	}
	api := v.Middleware(jwt.RequireScope("invoices.read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := jwt.ClaimsFromContext(r.Context())
		fmt.Fprintf(w, "invoices of tenant %s for %s", claims.TenantID, claims.UserID)
	})))

	req := httptest.NewRequest(http.MethodGet, "/invoices", nil)
	req.Header.Set("Authorization", "Bearer "+sign(map[string]any{
		"iss": "qeet-id", "aud": "https://api.example.com", "sub": "u-1", "user_id": "u-1",
		"tenant_id": "t-1", "sid": "s-1", "scope": "openid invoices.read",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	}))
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())
	// Output: 200 invoices of tenant t-1 for u-1
}

func ExampleVerifier_Verify() {
	jwks, sign := fakeQeetID()
	defer jwks.Close()

	v, err := jwt.NewVerifier("qeet-id", "https://api.example.com",
		jwt.WithJWKSURL(jwks.URL+"/.well-known/jwks.json"))
	if err != nil {
		log.Fatal(err)
	}
	expired := sign(map[string]any{"iss": "qeet-id", "aud": "https://api.example.com", "sub": "u-1",
		"exp": time.Now().Add(-time.Hour).Unix()})
	_, err = v.Verify(context.Background(), expired)
	fmt.Println(err)
	// Output: jwt: invalid token: expired
}

// fakeQeetID publishes one ES256 key as a JWKS and returns a function that
// signs claims with it, as Qeet ID does.
func fakeQeetID() (*httptest.Server, func(map[string]any) string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	pub, err := key.PublicKey.Bytes()
	if err != nil {
		log.Fatal(err)
	}
	enc := base64.RawURLEncoding
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kty": "EC", "crv": "P-256", "alg": "ES256", "use": "sig", "kid": "key-1",
			"x": enc.EncodeToString(pub[1:33]), "y": enc.EncodeToString(pub[33:]),
		}}})
	}))
	sign := func(claims map[string]any) string {
		h, _ := json.Marshal(map[string]string{"alg": "ES256", "typ": "JWT", "kid": "key-1"})
		p, _ := json.Marshal(claims)
		input := enc.EncodeToString(h) + "." + enc.EncodeToString(p)
		digest := sha256.Sum256([]byte(input))
		der, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
		if err != nil {
			log.Fatal(err)
		}
		var rs struct{ R, S *big.Int }
		if _, err := asn1.Unmarshal(der, &rs); err != nil {
			log.Fatal(err)
		}
		sig := make([]byte, 64)
		rs.R.FillBytes(sig[:32])
		rs.S.FillBytes(sig[32:])
		return input + "." + enc.EncodeToString(sig)
	}
	return srv, sign
}
