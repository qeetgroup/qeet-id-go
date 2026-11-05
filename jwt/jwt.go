// Package jwt verifies Qeet ID access tokens locally, against the signing keys
// Qeet ID publishes, so a resource server can authenticate requests without a
// call to Qeet ID per request.
//
// Qeet ID signs access tokens with ES256 and names the key in the kid header.
// A Verifier fetches the public keys from the JWKS endpoint, caches them,
// fetches them again when a token names a key it has not seen (a key
// rotation), and checks the signature, the expiry, the issuer and the
// audience:
//
//	v, err := jwt.NewVerifier(issuer, audience)
//	if err != nil {
//		log.Fatal(err)
//	}
//	http.Handle("/api/", v.Middleware(api))
//
//	func api(w http.ResponseWriter, r *http.Request) {
//		claims, _ := jwt.ClaimsFromContext(r.Context())
//		fmt.Fprintln(w, "hello", claims.UserID)
//	}
//
// The issuer is the iss value of Qeet ID's tokens; read it from the discovery
// document (OAuthService.Discover in the root package). The audience is the
// API's own audience, or the resource indicator the token was requested for.
//
// Local verification cannot see a session that ended before its token
// expired. Access tokens are short-lived; where revocation must take effect at
// once, subscribe to the session.revoked webhook event or ask
// OAuthService.Introspect.
package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Defaults of a Verifier. Each can be changed with an Option.
const (
	// DefaultJWKSURL is where Qeet ID publishes its public signing keys.
	DefaultJWKSURL = "https://api.id.qeet.in/.well-known/jwks.json"
	// DefaultLeeway is the clock skew allowed when checking exp and nbf.
	DefaultLeeway = 30 * time.Second
	// DefaultKeyCacheTTL is how long fetched keys are used before they are
	// fetched again.
	DefaultKeyCacheTTL = 10 * time.Minute
	// DefaultMinRefreshInterval is the shortest time between two fetches
	// triggered by unknown key IDs, so that tokens naming random keys cannot
	// make the verifier hammer the JWKS endpoint.
	DefaultMinRefreshInterval = time.Minute

	maxTokenBytes = 16 << 10
	algES256      = "ES256"
)

// Errors returned by Verify. Every problem with the token itself matches
// ErrInvalidToken; failing to obtain the signing keys matches
// ErrKeysUnavailable instead, because the token may well be valid.
var (
	// ErrInvalidToken matches every reason a token is refused.
	ErrInvalidToken = errors.New("jwt: invalid token")
	// ErrMalformed reports a token that is not a well-formed JWS compact JWT.
	ErrMalformed = fmt.Errorf("%w: malformed", ErrInvalidToken)
	// ErrUnsupportedAlgorithm reports an alg other than ES256, including
	// "none" and HMAC algorithms.
	ErrUnsupportedAlgorithm = fmt.Errorf("%w: unsupported algorithm", ErrInvalidToken)
	// ErrUnknownKey reports a kid that is missing or not among the published
	// keys, even after fetching them again.
	ErrUnknownKey = fmt.Errorf("%w: unknown signing key", ErrInvalidToken)
	// ErrInvalidSignature reports a signature that does not verify.
	ErrInvalidSignature = fmt.Errorf("%w: signature does not verify", ErrInvalidToken)
	// ErrExpired reports a token past its exp, beyond the leeway.
	ErrExpired = fmt.Errorf("%w: expired", ErrInvalidToken)
	// ErrNotYetValid reports a token before its nbf, beyond the leeway.
	ErrNotYetValid = fmt.Errorf("%w: not yet valid", ErrInvalidToken)
	// ErrIssuerMismatch reports an iss other than the verifier's issuer.
	ErrIssuerMismatch = fmt.Errorf("%w: issuer mismatch", ErrInvalidToken)
	// ErrAudienceMismatch reports an aud that does not include the verifier's
	// audience.
	ErrAudienceMismatch = fmt.Errorf("%w: audience mismatch", ErrInvalidToken)
	// ErrKeysUnavailable reports that the signing keys could not be fetched.
	ErrKeysUnavailable = errors.New("jwt: signing keys unavailable")
)

// Verifier verifies Qeet ID access tokens. It is safe for concurrent use;
// create one and share it, so that its key cache is shared too.
type Verifier struct {
	issuer   string
	audience string
	leeway   time.Duration
	now      func() time.Time
	keys     *keySet
}

// Option configures a Verifier.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

type config struct {
	jwksURL    string
	hc         *http.Client
	leeway     time.Duration
	now        func() time.Time
	cacheTTL   time.Duration
	minRefresh time.Duration
}

// WithJWKSURL sets where the verifier fetches the signing keys (default
// DefaultJWKSURL), such as the jwks_uri of the discovery document. It must
// use https, or http on a loopback host.
func WithJWKSURL(u string) Option {
	return optionFunc(func(c *config) { c.jwksURL = u })
}

// WithHTTPClient sets the HTTP client used to fetch the signing keys. The
// default has a 10-second timeout.
func WithHTTPClient(hc *http.Client) Option {
	return optionFunc(func(c *config) {
		if hc != nil {
			c.hc = hc
		}
	})
}

// WithLeeway sets the clock skew allowed when checking exp and nbf (default
// DefaultLeeway). A negative value is treated as zero.
func WithLeeway(d time.Duration) Option {
	return optionFunc(func(c *config) { c.leeway = max(d, 0) })
}

// WithClock sets the function the verifier reads the current time from.
func WithClock(now func() time.Time) Option {
	return optionFunc(func(c *config) {
		if now != nil {
			c.now = now
		}
	})
}

// WithKeyCacheTTL sets how long fetched keys are used before they are fetched
// again (default DefaultKeyCacheTTL).
func WithKeyCacheTTL(d time.Duration) Option {
	return optionFunc(func(c *config) {
		if d > 0 {
			c.cacheTTL = d
		}
	})
}

// WithMinRefreshInterval sets the shortest time between two key fetches
// triggered by unknown key IDs (default DefaultMinRefreshInterval).
func WithMinRefreshInterval(d time.Duration) Option {
	return optionFunc(func(c *config) {
		if d >= 0 {
			c.minRefresh = d
		}
	})
}

// NewVerifier returns a Verifier that accepts tokens whose iss equals issuer
// and whose aud includes audience. Both are required: checking them is what
// stops a token issued for another purpose from being accepted here.
// NewVerifier makes no network calls; keys are fetched on first use.
func NewVerifier(issuer, audience string, opts ...Option) (*Verifier, error) {
	if issuer == "" || audience == "" {
		return nil, errors.New("jwt: NewVerifier: the issuer and audience are required")
	}
	cfg := config{
		jwksURL:    DefaultJWKSURL,
		hc:         &http.Client{Timeout: 10 * time.Second},
		leeway:     DefaultLeeway,
		now:        time.Now,
		cacheTTL:   DefaultKeyCacheTTL,
		minRefresh: DefaultMinRefreshInterval,
	}
	for _, o := range opts {
		if o != nil {
			o.apply(&cfg)
		}
	}
	if err := checkJWKSURL(cfg.jwksURL); err != nil {
		return nil, err
	}
	return &Verifier{
		issuer:   issuer,
		audience: audience,
		leeway:   cfg.leeway,
		now:      cfg.now,
		keys: &keySet{
			url:        cfg.jwksURL,
			hc:         cfg.hc,
			now:        cfg.now,
			ttl:        cfg.cacheTTL,
			minRefresh: cfg.minRefresh,
		},
	}, nil
}

// Verify checks token and returns its claims. A token with any problem is
// refused with an error matching ErrInvalidToken; failing to fetch the keys
// returns an error matching ErrKeysUnavailable.
func (v *Verifier) Verify(ctx context.Context, token string) (*Claims, error) {
	if len(token) > maxTokenBytes {
		return nil, ErrMalformed
	}
	headerB64, rest, ok1 := strings.Cut(token, ".")
	payloadB64, sigB64, ok2 := strings.Cut(rest, ".")
	if !ok1 || !ok2 || strings.Contains(sigB64, ".") || headerB64 == "" || payloadB64 == "" {
		return nil, ErrMalformed
	}

	var header struct {
		Alg  string          `json:"alg"`
		Kid  string          `json:"kid"`
		Crit json.RawMessage `json:"crit"`
	}
	if err := decodeSegment(headerB64, &header); err != nil {
		return nil, err
	}
	if header.Alg != algES256 {
		return nil, ErrUnsupportedAlgorithm
	}
	if len(header.Crit) > 0 {
		return nil, fmt.Errorf("%w: unsupported critical header", ErrMalformed)
	}
	if header.Kid == "" {
		return nil, ErrUnknownKey
	}

	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil || len(sig) != 64 {
		return nil, ErrInvalidSignature
	}
	key, err := v.keys.get(ctx, header.Kid)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(headerB64 + "." + payloadB64))
	der, err := asn1.Marshal(struct{ R, S *big.Int }{
		R: new(big.Int).SetBytes(sig[:32]),
		S: new(big.Int).SetBytes(sig[32:]),
	})
	if err != nil || !ecdsa.VerifyASN1(key, digest[:], der) {
		return nil, ErrInvalidSignature
	}

	claims, err := parseClaims(payloadB64)
	if err != nil {
		return nil, err
	}
	now := v.now()
	switch {
	case claims.ExpiresAt.IsZero() || now.After(claims.ExpiresAt.Add(v.leeway)):
		return nil, ErrExpired
	case !claims.NotBefore.IsZero() && now.Add(v.leeway).Before(claims.NotBefore):
		return nil, ErrNotYetValid
	case claims.Issuer != v.issuer:
		return nil, ErrIssuerMismatch
	case !claims.HasAudience(v.audience):
		return nil, ErrAudienceMismatch
	}
	return claims, nil
}

func decodeSegment(seg string, v any) error {
	b, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return ErrMalformed
	}
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.UseNumber()
	if err := dec.Decode(v); err != nil {
		return ErrMalformed
	}
	return nil
}

func checkJWKSURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return fmt.Errorf("jwt: %q is not an absolute URL", raw)
	}
	switch u.Scheme {
	case "https":
		return nil
	case "http":
		host := u.Hostname()
		if host == "localhost" || strings.HasSuffix(host, ".localhost") {
			return nil
		}
		if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
			return nil
		}
	}
	return fmt.Errorf("jwt: the JWKS URL %q must use https; http is allowed only for a loopback host", raw)
}
