package qeetid

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

type TokenExchangeInput struct {
	ClientID       string
	ClientSecret   string
	SubjectToken   string
	Scope          string // optional downscoped permissions
	ActorToken     string // optional — for RFC 8693 delegation
	ActorTokenType string
}

type TokenExchangeResult struct {
	AccessToken     string `json:"access_token"`
	TokenType       string `json:"token_type"`
	ExpiresIn       int    `json:"expires_in"`
	Scope           string `json:"scope,omitempty"`
	IssuedTokenType string `json:"issued_token_type,omitempty"`
}

type IntrospectResult struct {
	Active    bool   `json:"active"`
	Sub       string `json:"sub,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
	ActorType string `json:"actor_type,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
}

// DeviceAuthorizationResult is the RFC 8628 device-flow starting point.
type DeviceAuthorizationResult struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval,omitempty"`
}

// DeviceContext is the verification-page context for a pending device code
// (GET /v1/oauth/device?user_code=).
type DeviceContext struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope,omitempty"`
}

// BackchannelAuthorizeInput starts a CIBA (client-initiated backchannel
// auth) request.
type BackchannelAuthorizeInput struct {
	ClientID       string
	ClientSecret   string
	LoginHint      string
	Scope          string
	BindingMessage string
}

// BackchannelAuthorizeResult is the pending-request handle from BackchannelAuthorize.
type BackchannelAuthorizeResult struct {
	AuthReqID string `json:"auth_req_id"`
	ExpiresIn int    `json:"expires_in"`
	Interval  int    `json:"interval"`
}

// PendingBackchannelRequest is one of the caller's own pending CIBA requests.
type PendingBackchannelRequest struct {
	ID        string `json:"id"`
	ClientID  string `json:"client_id"`
	LoginHint string `json:"login_hint"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

// SigningKey is one ES256 key in the platform's JWKS (published + retired).
type SigningKey struct {
	Kid    string `json:"kid"`
	Alg    string `json:"alg"`
	Use    string `json:"use"`
	Status string `json:"status"` // active | retired
}

// RotateSigningKeyResult carries the newly-minted key. PrivateKeyPEM is
// shown once — the platform never exposes it again after this response.
type RotateSigningKeyResult struct {
	Kid           string `json:"kid"`
	Alg           string `json:"alg"`
	PrivateKeyPEM string `json:"private_key_pem"`
	Warning       string `json:"warning,omitempty"`
}

// OAuthGrant is a client_credentials/authorization_code grant a tenant's
// user or service has active against an OIDC client.
type OAuthGrant struct {
	ID        string   `json:"id"`
	ClientID  string   `json:"client_id"`
	UserID    string   `json:"user_id,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	CreatedAt string   `json:"created_at"`
}

// OAuthDeviceAuthorization is an admin-visible device-flow session (distinct
// from the RFC 8628 device-flow endpoints above — this is the management
// view over active/pending device authorizations).
type OAuthDeviceAuthorization struct {
	ID           string   `json:"id"`
	ClientID     string   `json:"client_id"`
	UserCode     string   `json:"user_code"`
	Status       string   `json:"status"` // pending | authorized | denied
	UserID       string   `json:"user_id,omitempty"`
	UserEmail    string   `json:"user_email,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	CreatedAt    string   `json:"created_at"`
	ExpiresAt    string   `json:"expires_at"`
	LastPolledAt string   `json:"last_polled_at,omitempty"`
}

// SigningKeysService manages the platform's ES256 signing-key set.
type SigningKeysService struct{ t *transport.Transport }

func (s *SigningKeysService) List(ctx context.Context) ([]SigningKey, error) {
	var out struct {
		Keys []SigningKey `json:"keys"`
	}
	err := s.t.Get(ctx, "/v1/oidc/signing-keys", transport.RequestOptions{}, &out)
	return out.Keys, err
}

func (s *SigningKeysService) Rotate(ctx context.Context) (*RotateSigningKeyResult, error) {
	var out RotateSigningKeyResult
	err := s.t.Post(ctx, "/v1/oidc/signing-keys/rotate", nil, &out)
	return &out, err
}

// OAuthGrantsService is the admin view over active client grants.
type OAuthGrantsService struct{ t *transport.Transport }

func (s *OAuthGrantsService) List(ctx context.Context, tenantID string) ([]OAuthGrant, error) {
	var env marshal.Envelope[OAuthGrant]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oauth/grants", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *OAuthGrantsService) Revoke(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oauth/grants/"+url.PathEscape(id), nil)
}

// OAuthDevicesService is the admin view over device-flow sessions.
type OAuthDevicesService struct{ t *transport.Transport }

func (s *OAuthDevicesService) List(ctx context.Context, tenantID string) ([]OAuthDeviceAuthorization, error) {
	var env marshal.Envelope[OAuthDeviceAuthorization]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oauth/devices", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

func (s *OAuthDevicesService) Revoke(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/oauth/devices/"+url.PathEscape(id), nil)
}

// OAuthService provides RFC 8693 token exchange, RFC 7662 introspection, an
// MCP token guard, RFC 8628 device flow, and CIBA. These use form-encoded
// requests with OIDC client credentials (transport.DoForm) rather than the
// management API's ApiKey header — an OAuth client_id/secret pair, not a
// qk_… key, authenticates them. SigningKeys/Grants/Devices are the admin
// JSON sub-resources and use the normal ApiKey-authed path.
type OAuthService struct {
	t *transport.Transport

	SigningKeys *SigningKeysService
	Grants      *OAuthGrantsService
	Devices     *OAuthDevicesService
}

func newOAuthService(t *transport.Transport) *OAuthService {
	return &OAuthService{
		t:           t,
		SigningKeys: &SigningKeysService{t: t},
		Grants:      &OAuthGrantsService{t: t},
		Devices:     &OAuthDevicesService{t: t},
	}
}

// TokenExchange implements RFC 8693 downscoping and delegation.
func (s *OAuthService) TokenExchange(ctx context.Context, in TokenExchangeInput) (*TokenExchangeResult, error) {
	form := url.Values{
		"grant_type":           {"urn:ietf:params:oauth:grant-type:token-exchange"},
		"subject_token":        {in.SubjectToken},
		"subject_token_type":   {"urn:ietf:params:oauth:token-type:access_token"},
		"requested_token_type": {"urn:ietf:params:oauth:token-type:access_token"},
	}
	if in.Scope != "" {
		form.Set("scope", in.Scope)
	}
	if in.ActorToken != "" {
		form.Set("actor_token", in.ActorToken)
		typ := in.ActorTokenType
		if typ == "" {
			typ = "urn:ietf:params:oauth:token-type:access_token"
		}
		form.Set("actor_token_type", typ)
	}
	var out TokenExchangeResult
	err := s.t.DoForm(ctx, http.MethodPost, "/v1/oauth/token", form, in.ClientID, in.ClientSecret, &out)
	return &out, err
}

// Introspect implements RFC 7662 token introspection. Note the unprefixed
// path — /oauth/introspect, not /v1/oauth/introspect (confirmed against the
// spec and the router's CSRF-exempt path list; the previous SDK version had
// this wrong and would 404).
func (s *OAuthService) Introspect(ctx context.Context, token string) (*IntrospectResult, error) {
	form := url.Values{"token": {token}}
	var out IntrospectResult
	err := s.t.DoForm(ctx, http.MethodPost, "/oauth/introspect", form, "", "", &out)
	return &out, err
}

// Revoke implements RFC 7009 token revocation. Same unprefixed-path fix as
// Introspect.
func (s *OAuthService) Revoke(ctx context.Context, token string) error {
	form := url.Values{"token": {token}}
	return s.t.DoForm(ctx, http.MethodPost, "/oauth/revoke", form, "", "", nil)
}

// Verify is an MCP token guard: introspects the token and returns an error if
// it is inactive or does not carry requiredScope (empty string skips scope check).
func (s *OAuthService) Verify(ctx context.Context, token, requiredScope string) (*IntrospectResult, error) {
	result, err := s.Introspect(ctx, token)
	if err != nil {
		return nil, err
	}
	if !result.Active {
		return nil, &Error{Status: 401, Code: "token_inactive", Message: "token is not active"}
	}
	if requiredScope != "" {
		scopes := strings.Fields(result.Scope)
		found := false
		for _, sc := range scopes {
			if sc == requiredScope {
				found = true
				break
			}
		}
		if !found {
			return nil, &Error{Status: 403, Code: "insufficient_scope", Message: "required scope: " + requiredScope}
		}
	}
	return result, nil
}

// DeviceAuthorize starts an RFC 8628 device-flow grant.
func (s *OAuthService) DeviceAuthorize(ctx context.Context, clientID, clientSecret, scope string) (*DeviceAuthorizationResult, error) {
	form := url.Values{"client_id": {clientID}}
	if scope != "" {
		form.Set("scope", scope)
	}
	var out DeviceAuthorizationResult
	err := s.t.DoForm(ctx, http.MethodPost, "/v1/oauth/device_authorization", form, "", "", &out)
	_ = clientSecret // reserved: confidential clients may need Basic auth here too
	return &out, err
}

// DeviceContext fetches the verification-page context for a pending device
// code — the SSO-cookie-authenticated user is about to approve/deny it.
func (s *OAuthService) DeviceContext(ctx context.Context, userCode string) (*DeviceContext, error) {
	q := url.Values{"user_code": {userCode}}
	var out DeviceContext
	err := s.t.Get(ctx, "/v1/oauth/device", transport.RequestOptions{Query: q}, &out)
	return &out, err
}

// DeviceDecision approves or denies a pending device-flow user code.
func (s *OAuthService) DeviceDecision(ctx context.Context, userCode string, approve bool) error {
	body := map[string]any{"user_code": userCode, "approve": approve}
	return s.t.Post(ctx, "/v1/oauth/device/decision", body, nil)
}

// BackchannelAuthorize starts a CIBA request.
func (s *OAuthService) BackchannelAuthorize(ctx context.Context, in BackchannelAuthorizeInput) (*BackchannelAuthorizeResult, error) {
	form := url.Values{"login_hint": {in.LoginHint}}
	if in.Scope != "" {
		form.Set("scope", in.Scope)
	}
	if in.BindingMessage != "" {
		form.Set("binding_message", in.BindingMessage)
	}
	if in.ClientID != "" {
		form.Set("client_id", in.ClientID)
	}
	var out BackchannelAuthorizeResult
	err := s.t.DoForm(ctx, http.MethodPost, "/v1/oauth/bc-authorize", form, in.ClientID, in.ClientSecret, &out)
	return &out, err
}

// BackchannelPending lists the caller's own pending CIBA requests.
func (s *OAuthService) BackchannelPending(ctx context.Context) ([]PendingBackchannelRequest, error) {
	var env marshal.Envelope[PendingBackchannelRequest]
	err := s.t.Get(ctx, "/v1/oauth/bc-authorize/pending", transport.RequestOptions{}, &env)
	return env.Resolve(), err
}

// BackchannelDecision approves or denies a pending CIBA request.
func (s *OAuthService) BackchannelDecision(ctx context.Context, id string, approve bool) error {
	body := map[string]any{"id": id, "approve": approve}
	return s.t.Post(ctx, "/v1/oauth/bc-authorize/decision", body, nil)
}
