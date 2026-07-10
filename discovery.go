package qeetid

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

const discoveryPath = "/.well-known/openid-configuration"

// DiscoveryDocument is the OIDC/OAuth provider metadata (OpenID Connect
// Discovery / RFC 8414) published at {issuer}/.well-known/openid-configuration.
// Raw holds every field, including the Qeet-specific extensions
// (actor_types_supported, resource_indicators_supported).
type DiscoveryDocument struct {
	Issuer                           string         `json:"issuer"`
	AuthorizationEndpoint            string         `json:"authorization_endpoint"`
	TokenEndpoint                    string         `json:"token_endpoint"`
	UserinfoEndpoint                 string         `json:"userinfo_endpoint"`
	JWKSURI                          string         `json:"jwks_uri"`
	RevocationEndpoint               string         `json:"revocation_endpoint"`
	IntrospectionEndpoint            string         `json:"introspection_endpoint"`
	EndSessionEndpoint               string         `json:"end_session_endpoint"`
	DeviceAuthorizationEndpoint      string         `json:"device_authorization_endpoint"`
	BackchannelAuthEndpoint          string         `json:"backchannel_authentication_endpoint"`
	GrantTypesSupported              []string       `json:"grant_types_supported"`
	IDTokenSigningAlgValuesSupported []string       `json:"id_token_signing_alg_values_supported"`
	CodeChallengeMethodsSupported    []string       `json:"code_challenge_methods_supported"`
	ActorTypesSupported              []string       `json:"actor_types_supported"`
	ResourceIndicatorsSupported      bool           `json:"resource_indicators_supported"`
	Raw                              map[string]any `json:"-"`
}

// Discover fetches provider metadata from issuer + /.well-known/openid-configuration.
// issuer may be the bare base URL or already include the well-known path. hc is
// optional (defaults to a client with a 10s timeout).
func Discover(ctx context.Context, issuer string, hc *http.Client) (*DiscoveryDocument, error) {
	if hc == nil {
		hc = &http.Client{Timeout: 10 * time.Second}
	}
	u := strings.TrimRight(issuer, "/")
	if !strings.HasSuffix(u, discoveryPath) {
		u += discoveryPath
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, &Error{Code: "discovery_error", Message: "build request: " + err.Error()}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", transport.UserAgent())

	res, err := hc.Do(req)
	if err != nil {
		return nil, &Error{Code: "network_error", Message: "discovery fetch: " + err.Error()}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, &Error{Status: res.StatusCode, Code: "discovery_error", Message: "discovery fetch failed"}
	}

	var raw map[string]any
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, &Error{Code: "discovery_error", Message: "decode metadata: " + err.Error()}
	}
	// Re-marshal into the typed struct so callers get both typed and raw views.
	b, _ := json.Marshal(raw)
	var doc DiscoveryDocument
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, &Error{Code: "discovery_error", Message: "parse metadata: " + err.Error()}
	}
	doc.Raw = raw
	return &doc, nil
}

// Discover fetches this client's provider metadata (from its configured base URL).
func (c *Client) Discover(ctx context.Context) (*DiscoveryDocument, error) {
	return Discover(ctx, c.transport.BaseURL(), c.transport.HTTPClient())
}

// NewFromDiscovery builds a client and self-configures it from the provider's
// discovery document. Unlike New, it makes one network call (to fetch metadata)
// and wires session verification to the published jwks_uri — so a self-hosted
// instance serving JWKS at a non-default path works with no extra config. The
// resolved metadata is returned alongside the client.
func NewFromDiscovery(ctx context.Context, cfg Config) (*Client, *DiscoveryDocument, error) {
	c := New(cfg)
	doc, err := c.Discover(ctx)
	if err != nil {
		return nil, nil, err
	}
	if doc.JWKSURI != "" {
		c.Sessions.setJWKSURL(doc.JWKSURI)
	}
	return c, doc, nil
}
