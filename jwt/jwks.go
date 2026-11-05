package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/qeetgroup/qeet-id-go/internal/version"
)

const maxJWKSBytes = 1 << 20

// keySet caches the published signing keys.
type keySet struct {
	url        string
	hc         *http.Client
	now        func() time.Time
	ttl        time.Duration
	minRefresh time.Duration

	mu          sync.Mutex
	keys        map[string]*ecdsa.PublicKey
	fetchedAt   time.Time
	lastAttempt time.Time
	inFlight    bool   // a fetch is running; callers wait for it
	generation  uint64 // completed fetches, successful or not
	lastErr     error  // the last fetch's error, nil on success

	fetchMu sync.Mutex // one fetch at a time; waiters reuse its result
}

// get returns the key named kid. It fetches the key set when there is none,
// when the cached one is older than the TTL, or when kid is unknown; fetches
// are at least minRefresh apart. A failed fetch falls back to the cached keys
// when they hold kid.
func (k *keySet) get(ctx context.Context, kid string) (*ecdsa.PublicKey, error) {
	key, fresh, empty, mayFetch := k.lookup(kid)
	if key != nil && fresh {
		return key, nil
	}
	if !mayFetch {
		switch {
		case key != nil:
			return key, nil
		case empty:
			return nil, fmt.Errorf("%w: waiting before fetching the keys again", ErrKeysUnavailable)
		}
		return nil, ErrUnknownKey
	}
	if err := k.refresh(ctx); err != nil {
		if key != nil {
			return key, nil // serve the stale key rather than fail every request
		}
		return nil, err
	}
	if key, _, _, _ = k.lookup(kid); key != nil {
		return key, nil
	}
	return nil, ErrUnknownKey
}

// lookup returns the cached key for kid, whether the cache is within its TTL,
// whether it is empty, and whether a fetch may be made now.
func (k *keySet) lookup(kid string) (key *ecdsa.PublicKey, fresh, empty, mayFetch bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	now := k.now()
	key = k.keys[kid]
	empty = k.keys == nil
	fresh = !empty && now.Sub(k.fetchedAt) < k.ttl
	cooledDown := k.inFlight || k.lastAttempt.IsZero() || now.Sub(k.lastAttempt) >= k.minRefresh
	mayFetch = cooledDown && (empty || key == nil || !fresh)
	return key, fresh, empty, mayFetch
}

// refresh fetches the key set. Concurrent callers wait for one fetch and share
// its outcome.
func (k *keySet) refresh(ctx context.Context) error {
	k.mu.Lock()
	generation := k.generation
	k.mu.Unlock()

	k.fetchMu.Lock()
	defer k.fetchMu.Unlock()

	k.mu.Lock()
	if k.generation != generation {
		err := k.lastErr // a fetch completed while this caller waited
		k.mu.Unlock()
		return err
	}
	k.lastAttempt, k.inFlight = k.now(), true
	k.mu.Unlock()

	keys, err := k.fetch(ctx)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrKeysUnavailable, err)
	}
	k.mu.Lock()
	k.inFlight, k.lastErr = false, err
	k.generation++
	if err == nil {
		k.keys, k.fetchedAt = keys, k.now()
	}
	k.mu.Unlock()
	return err
}

type jwk struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
}

func (k *keySet) fetch(ctx context.Context) (map[string]*ecdsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", version.UserAgent())
	resp, err := k.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s answered %d", k.url, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxJWKSBytes))
	if err != nil {
		return nil, err
	}
	var doc struct {
		Keys []jwk `json:"keys"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}
	keys := make(map[string]*ecdsa.PublicKey, len(doc.Keys))
	for _, j := range doc.Keys {
		if pub, err := j.publicKey(); err == nil && j.Kid != "" {
			keys[j.Kid] = pub
		}
	}
	if len(keys) == 0 {
		return nil, errors.New("the JWKS holds no usable ES256 key")
	}
	return keys, nil
}

// publicKey converts an EC P-256 signing JWK to a public key. Keys for other
// curves, algorithms or uses are rejected.
func (j jwk) publicKey() (*ecdsa.PublicKey, error) {
	if j.Kty != "EC" || j.Crv != "P-256" || (j.Alg != "" && j.Alg != algES256) || (j.Use != "" && j.Use != "sig") {
		return nil, errors.New("not an ES256 signing key")
	}
	x, err1 := base64.RawURLEncoding.DecodeString(j.X)
	y, err2 := base64.RawURLEncoding.DecodeString(j.Y)
	if err1 != nil || err2 != nil || len(x) > 32 || len(y) > 32 {
		return nil, errors.New("malformed coordinates")
	}
	point := make([]byte, 65)
	point[0] = 4 // uncompressed
	copy(point[1+32-len(x):33], x)
	copy(point[33+32-len(y):], y)
	return ecdsa.ParseUncompressedPublicKey(elliptic.P256(), point)
}
