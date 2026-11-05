package qeetid

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// OAuthService runs the OAuth 2.0 and OpenID Connect requests in which your
// application is the client of Qeet ID: obtaining, exchanging, refreshing,
// introspecting and revoking tokens, the device and backchannel (CIBA) flows,
// and reading provider metadata.
//
// These calls authenticate the application with its client ID and secret,
// sent with HTTP Basic, not with the client's own credential. For a public
// client, leave ClientSecret empty and the client ID is sent in the form.
//
// Token endpoint failures are *OAuthError values carrying the RFC 6749 error
// code. A request refused by rate limiting is an *APIError matching
// ErrRateLimited, which the client retries after the Retry-After the server
// gives.
type OAuthService struct{ c *Client }

// Token types for token exchange (RFC 8693).
const (
	TokenTypeAccessToken  = "urn:ietf:params:oauth:token-type:access_token"
	TokenTypeRefreshToken = "urn:ietf:params:oauth:token-type:refresh_token"
)

const (
	grantTypeDeviceCode    = "urn:ietf:params:oauth:grant-type:device_code"
	grantTypeTokenExchange = "urn:ietf:params:oauth:grant-type:token-exchange"
	grantTypeCIBA          = "urn:openid:params:grant-type:ciba"
)

// ClientCredentialsParams are the parameters of OAuthService.ClientCredentials.
type ClientCredentialsParams struct {
	// ClientID is the service principal's client ID.
	ClientID string
	// ClientSecret is the service principal's secret.
	ClientSecret string
	// Scope narrows the token to these space-separated scopes. Empty asks for
	// every scope the service principal holds.
	Scope string
}

// ClientCredentials obtains an access token for a service principal with the
// client credentials grant. Token.ExpiresIn and Token.ExpiresAt give its
// lifetime. To have the client fetch and renew such tokens for its own
// requests, use WithClientCredentials instead.
//
// API: POST /v1/oauth/token
func (s *OAuthService) ClientCredentials(ctx context.Context, params ClientCredentialsParams, opts ...RequestOption) (*Token, error) {
	if params.ClientID == "" || params.ClientSecret == "" {
		return nil, errors.New("qeetid: ClientCredentials: ClientID and ClientSecret are required")
	}
	return s.c.clientCredentialsToken(ctx, params.ClientID, params.ClientSecret, params.Scope, opts...)
}

// ExchangeCodeParams are the parameters of OAuthService.ExchangeCode.
type ExchangeCodeParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret. Leave it empty for a public
	// client, which proves itself with CodeVerifier.
	ClientSecret string
	// Code is the authorization code from the redirect.
	Code string
	// RedirectURI is the redirect URI the authorization request used.
	RedirectURI string
	// CodeVerifier is the PKCE verifier whose challenge the authorization
	// request sent.
	CodeVerifier string
	// Resource binds the access token to a resource server (RFC 8707).
	Resource string
}

// ExchangeCode redeems an authorization code for tokens (the
// authorization_code grant).
//
// API: POST /v1/oauth/token-code
func (s *OAuthService) ExchangeCode(ctx context.Context, params ExchangeCodeParams, opts ...RequestOption) (*Token, error) {
	if params.Code == "" {
		return nil, errors.New("qeetid: ExchangeCode: Code is required")
	}
	form := url.Values{"grant_type": {"authorization_code"}, "code": {params.Code}}
	setIf(form, "redirect_uri", params.RedirectURI)
	setIf(form, "code_verifier", params.CodeVerifier)
	setIf(form, "resource", params.Resource)
	return s.token(ctx, params.ClientID, params.ClientSecret, form, opts)
}

// RefreshTokenParams are the parameters of OAuthService.RefreshToken.
type RefreshTokenParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret; empty for a public client.
	ClientSecret string
	// RefreshToken is the refresh token to redeem. The server rotates it: use
	// the refresh token in the answer next time.
	RefreshToken string
	// Resource switches the resource server the new access token is bound to
	// (RFC 8707). Empty keeps the one the grant is bound to.
	Resource string
}

// RefreshToken renews a grant with its refresh token (the refresh_token
// grant).
//
// API: POST /v1/oauth/token-code
func (s *OAuthService) RefreshToken(ctx context.Context, params RefreshTokenParams, opts ...RequestOption) (*Token, error) {
	if params.RefreshToken == "" {
		return nil, errors.New("qeetid: RefreshToken: RefreshToken is required")
	}
	form := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {params.RefreshToken}}
	setIf(form, "resource", params.Resource)
	return s.token(ctx, params.ClientID, params.ClientSecret, form, opts)
}

// TokenExchangeParams are the parameters of OAuthService.ExchangeToken.
type TokenExchangeParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// SubjectToken is the token of the party whose authority is exercised.
	SubjectToken string
	// SubjectTokenType is the type of SubjectToken. Empty means
	// TokenTypeAccessToken, the only type the server accepts.
	SubjectTokenType string
	// ActorToken names the delegate acting for the subject, such as an AI
	// agent. The issued token carries an "act" claim for it.
	ActorToken string
	// ActorTokenType is the type of ActorToken. Empty means
	// TokenTypeAccessToken when ActorToken is set.
	ActorTokenType string
	// RequestedTokenType is the type of token wanted. Empty lets the server
	// choose.
	RequestedTokenType string
	// Scope narrows the issued token to these space-separated scopes.
	Scope string
	// Resource binds the issued token to a resource server (RFC 8707).
	Resource string
}

// ExchangeToken exchanges a token for another, narrower or delegated one
// (RFC 8693 token exchange).
//
// API: POST /v1/oauth/token-code
func (s *OAuthService) ExchangeToken(ctx context.Context, params TokenExchangeParams, opts ...RequestOption) (*Token, error) {
	if params.SubjectToken == "" {
		return nil, errors.New("qeetid: ExchangeToken: SubjectToken is required")
	}
	subjectType := params.SubjectTokenType
	if subjectType == "" {
		subjectType = TokenTypeAccessToken
	}
	form := url.Values{
		"grant_type":         {grantTypeTokenExchange},
		"subject_token":      {params.SubjectToken},
		"subject_token_type": {subjectType},
	}
	if params.ActorToken != "" {
		actorType := params.ActorTokenType
		if actorType == "" {
			actorType = TokenTypeAccessToken
		}
		form.Set("actor_token", params.ActorToken)
		form.Set("actor_token_type", actorType)
	}
	setIf(form, "requested_token_type", params.RequestedTokenType)
	setIf(form, "scope", params.Scope)
	setIf(form, "resource", params.Resource)
	return s.token(ctx, params.ClientID, params.ClientSecret, form, opts)
}

// DeviceAuthorizationParams are the parameters of OAuthService.AuthorizeDevice.
type DeviceAuthorizationParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret; empty for a public client.
	ClientSecret string
	// Scope lists the space-separated scopes to ask for.
	Scope string
}

// DeviceAuthorizationResult is the device authorization response (RFC 8628
// section 3.2).
type DeviceAuthorizationResult struct {
	// DeviceCode is the code the device polls with. Keep it secret.
	DeviceCode string `json:"device_code"`
	// UserCode is the short code the person enters, such as "WDJB-MJHT".
	UserCode string `json:"user_code"`
	// VerificationURI is where the person enters UserCode.
	VerificationURI string `json:"verification_uri"`
	// VerificationURIComplete is VerificationURI with the code filled in,
	// suitable for a QR code.
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	// ExpiresIn is how long, in seconds, the codes stay valid.
	ExpiresIn int `json:"expires_in"`
	// Interval is the minimum number of seconds between polls.
	Interval int `json:"interval,omitempty"`
}

// AuthorizeDevice starts the device flow (RFC 8628): show the person UserCode
// and VerificationURI, then poll with PollDeviceToken.
//
// API: POST /v1/oauth/device_authorization
func (s *OAuthService) AuthorizeDevice(ctx context.Context, params DeviceAuthorizationParams, opts ...RequestOption) (*DeviceAuthorizationResult, error) {
	form := url.Values{}
	setIf(form, "scope", params.Scope)
	req, err := oauthRequest("/v1/oauth/device_authorization", params.ClientID, params.ClientSecret, form)
	if err != nil {
		return nil, err
	}
	return call[DeviceAuthorizationResult](ctx, s.c, req, opts)
}

// DeviceTokenParams are the parameters of OAuthService.PollDeviceToken.
type DeviceTokenParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret; empty for a public client.
	ClientSecret string
	// DeviceCode is the code AuthorizeDevice returned.
	DeviceCode string
}

// PollDeviceToken asks once whether the person approved a device
// authorization. Until they act it fails with an *OAuthError whose Code is
// OAuthAuthorizationPending; poll again after the interval, and add five
// seconds to the interval on OAuthSlowDown. OAuthAccessDenied and
// OAuthExpiredToken end the flow.
//
// API: POST /v1/oauth/token-code
func (s *OAuthService) PollDeviceToken(ctx context.Context, params DeviceTokenParams, opts ...RequestOption) (*Token, error) {
	if params.DeviceCode == "" {
		return nil, errors.New("qeetid: PollDeviceToken: DeviceCode is required")
	}
	form := url.Values{"grant_type": {grantTypeDeviceCode}, "device_code": {params.DeviceCode}}
	return s.token(ctx, params.ClientID, params.ClientSecret, form, opts)
}

// BackchannelAuthorizationParams are the parameters of
// OAuthService.AuthorizeBackchannel.
type BackchannelAuthorizationParams struct {
	// ClientID is the application's client ID. Only confidential clients may
	// use the backchannel flow.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// LoginHint is the email address of the person to authenticate, a member
	// of the application's tenant.
	LoginHint string
	// Scope lists the space-separated scopes to ask for.
	Scope string
	// BindingMessage is shown to the person with the request, such as a
	// transaction reference.
	BindingMessage string
}

// BackchannelAuthorizationResult is the backchannel authentication response
// (OpenID CIBA section 7.3).
type BackchannelAuthorizationResult struct {
	// AuthReqID identifies the request when polling.
	AuthReqID string `json:"auth_req_id"`
	// ExpiresIn is how long, in seconds, the request stays valid.
	ExpiresIn int `json:"expires_in"`
	// Interval is the minimum number of seconds between polls.
	Interval int `json:"interval"`
}

// AuthorizeBackchannel asks a person to approve a sign-in on a device they
// already use (OpenID CIBA, poll mode). Poll for the outcome with
// PollBackchannelToken.
//
// API: POST /v1/oauth/bc-authorize
func (s *OAuthService) AuthorizeBackchannel(ctx context.Context, params BackchannelAuthorizationParams, opts ...RequestOption) (*BackchannelAuthorizationResult, error) {
	if params.LoginHint == "" {
		return nil, errors.New("qeetid: AuthorizeBackchannel: LoginHint is required")
	}
	if params.ClientSecret == "" {
		return nil, errors.New("qeetid: AuthorizeBackchannel: only a confidential client with a ClientSecret may use the backchannel flow")
	}
	form := url.Values{"login_hint": {params.LoginHint}}
	setIf(form, "scope", params.Scope)
	setIf(form, "binding_message", params.BindingMessage)
	req, err := oauthRequest("/v1/oauth/bc-authorize", params.ClientID, params.ClientSecret, form)
	if err != nil {
		return nil, err
	}
	return call[BackchannelAuthorizationResult](ctx, s.c, req, opts)
}

// BackchannelTokenParams are the parameters of
// OAuthService.PollBackchannelToken.
type BackchannelTokenParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// AuthReqID is the identifier AuthorizeBackchannel returned.
	AuthReqID string
}

// PollBackchannelToken asks once whether the person approved a backchannel
// request. It fails like PollDeviceToken while the request is pending.
//
// API: POST /v1/oauth/token-code
func (s *OAuthService) PollBackchannelToken(ctx context.Context, params BackchannelTokenParams, opts ...RequestOption) (*Token, error) {
	if params.AuthReqID == "" {
		return nil, errors.New("qeetid: PollBackchannelToken: AuthReqID is required")
	}
	form := url.Values{"grant_type": {grantTypeCIBA}, "auth_req_id": {params.AuthReqID}}
	return s.token(ctx, params.ClientID, params.ClientSecret, form, opts)
}

// IntrospectParams are the parameters of OAuthService.Introspect.
type IntrospectParams struct {
	// ClientID is the application's client ID. Introspection needs a
	// confidential client.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// Token is the access or refresh token to examine.
	Token string
	// TokenTypeHint is "access_token" or "refresh_token", when known.
	TokenTypeHint string
}

// Introspection is a token introspection response (RFC 7662). An application
// learns only about tokens it owns; any other token, and any token that is no
// longer valid, comes back with Active false and nothing else.
type Introspection struct {
	// Active reports whether the token is valid now: signed, unexpired, and
	// its session, account, membership and agent all still active.
	Active bool
	// TokenType is "Bearer" for an access token or "refresh_token".
	TokenType string
	// Subject is the token's subject (the sub claim).
	Subject string
	// Scope lists the token's space-separated scopes.
	Scope string
	// Issuer is the token's issuer.
	Issuer string
	// Audience lists the token's audiences.
	Audience []string
	// ClientID is the application the token was issued to.
	ClientID string
	// Resource lists the resource servers the token was issued for.
	Resource []string
	// ExpiresAt is when the token expires.
	ExpiresAt time.Time
	// IssuedAt is when the token was issued.
	IssuedAt time.Time
	// TenantID is the tenant the token is scoped to.
	TenantID string
	// SessionID is the session the token belongs to.
	SessionID string
	// ActorType is "service" or "agent" for a machine token; empty for a
	// person's token.
	ActorType string
	// AgentID is the agent an agent token belongs to.
	AgentID string
	// Actor is the delegate named by a token exchange's "act" claim.
	Actor *TokenActor
}

// TokenActor is the acting party of a delegated token (RFC 8693 section 4.1).
type TokenActor struct {
	// Subject is the actor's subject.
	Subject string `json:"sub"`
}

// UnmarshalJSON decodes an introspection response, in which aud and resource
// may each be a string or a list and exp and iat are Unix seconds.
func (in *Introspection) UnmarshalJSON(b []byte) error {
	var raw struct {
		Active    bool        `json:"active"`
		TokenType string      `json:"token_type"`
		Subject   string      `json:"sub"`
		Scope     string      `json:"scope"`
		Issuer    string      `json:"iss"`
		Audience  stringList  `json:"aud"`
		ClientID  string      `json:"client_id"`
		Resource  stringList  `json:"resource"`
		Exp       json.Number `json:"exp"`
		Iat       json.Number `json:"iat"`
		TenantID  string      `json:"tenant_id"`
		SessionID string      `json:"sid"`
		ActorType string      `json:"actor_type"`
		AgentID   string      `json:"agent_id"`
		Actor     *TokenActor `json:"act"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	*in = Introspection{
		Active: raw.Active, TokenType: raw.TokenType, Subject: raw.Subject, Scope: raw.Scope,
		Issuer: raw.Issuer, Audience: raw.Audience, ClientID: raw.ClientID, Resource: raw.Resource,
		ExpiresAt: unixTime(raw.Exp), IssuedAt: unixTime(raw.Iat), TenantID: raw.TenantID,
		SessionID: raw.SessionID, ActorType: raw.ActorType, AgentID: raw.AgentID, Actor: raw.Actor,
	}
	return nil
}

// Introspect asks whether a token is active and what it carries (RFC 7662).
// Introspection is the authoritative answer on revocation: a token whose
// session ended or whose account was suspended is reported inactive at once,
// while a signature check alone accepts it until it expires.
//
// API: POST /oauth/introspect
func (s *OAuthService) Introspect(ctx context.Context, params IntrospectParams, opts ...RequestOption) (*Introspection, error) {
	if params.Token == "" {
		return nil, errors.New("qeetid: Introspect: Token is required")
	}
	if params.ClientSecret == "" {
		return nil, errors.New("qeetid: Introspect: a public client cannot introspect; ClientSecret is required")
	}
	form := url.Values{"token": {params.Token}}
	setIf(form, "token_type_hint", params.TokenTypeHint)
	req, err := oauthRequest("/oauth/introspect", params.ClientID, params.ClientSecret, form)
	if err != nil {
		return nil, err
	}
	return call[Introspection](ctx, s.c, req, opts)
}

// RevokeParams are the parameters of OAuthService.Revoke.
type RevokeParams struct {
	// ClientID is the application's client ID.
	ClientID string
	// ClientSecret is the application's secret; empty for a public client.
	ClientSecret string
	// Token is the refresh token to revoke. Access tokens cannot be revoked;
	// they expire.
	Token string
	// TokenTypeHint is "refresh_token" or "access_token", when known.
	TokenTypeHint string
}

// Revoke revokes a refresh token and the whole grant it belongs to (RFC 7009).
// It succeeds for unknown and already invalid tokens, as the RFC requires. An
// access token is refused with OAuthUnsupportedTokenType.
//
// API: POST /oauth/revoke
func (s *OAuthService) Revoke(ctx context.Context, params RevokeParams, opts ...RequestOption) error {
	if params.Token == "" {
		return errors.New("qeetid: Revoke: Token is required")
	}
	form := url.Values{"token": {params.Token}}
	setIf(form, "token_type_hint", params.TokenTypeHint)
	req, err := oauthRequest("/oauth/revoke", params.ClientID, params.ClientSecret, form)
	if err != nil {
		return err
	}
	return s.c.do(ctx, req, nil, opts)
}

// UserInfo holds the OpenID Connect claims about the person an access token
// was issued for. Each claim is present only when the token's scopes grant it.
type UserInfo struct {
	// Subject is the person's user ID.
	Subject string `json:"sub"`
	// TenantID is the person's tenant.
	TenantID string `json:"tenant_id,omitempty"`
	// Email and EmailVerified come with the email scope.
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	// Name, Picture, Locale and UpdatedAt come with the profile scope.
	// UpdatedAt is in Unix seconds.
	Name      string `json:"name,omitempty"`
	Picture   string `json:"picture,omitempty"`
	Locale    string `json:"locale,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
	// PhoneNumber and PhoneNumberVerified come with the phone scope.
	PhoneNumber         string `json:"phone_number,omitempty"`
	PhoneNumberVerified bool   `json:"phone_number_verified,omitempty"`
}

// UserInfo returns the claims about the person accessToken was issued for. It
// sends accessToken, not the client's credential. The token must carry the
// openid scope.
//
// API: GET /v1/oauth/userinfo
func (s *OAuthService) UserInfo(ctx context.Context, accessToken string, opts ...RequestOption) (*UserInfo, error) {
	if accessToken == "" {
		return nil, errors.New("qeetid: UserInfo: accessToken is required")
	}
	return call[UserInfo](ctx, s.c, &request{method: http.MethodGet, path: "/v1/oauth/userinfo", bearer: accessToken}, opts)
}

// ProviderMetadata is Qeet ID's OpenID Connect discovery document. Issuer is
// the value the iss claim of every token carries; pass it to
// github.com/qeetgroup/qeet-id-go/jwt.NewVerifier.
type ProviderMetadata struct {
	Issuer                                    string   `json:"issuer"`
	AuthorizationEndpoint                     string   `json:"authorization_endpoint"`
	TokenEndpoint                             string   `json:"token_endpoint"`
	DeviceAuthorizationEndpoint               string   `json:"device_authorization_endpoint,omitempty"`
	BackchannelAuthenticationEndpoint         string   `json:"backchannel_authentication_endpoint,omitempty"`
	UserinfoEndpoint                          string   `json:"userinfo_endpoint,omitempty"`
	JWKSURI                                   string   `json:"jwks_uri"`
	RevocationEndpoint                        string   `json:"revocation_endpoint,omitempty"`
	IntrospectionEndpoint                     string   `json:"introspection_endpoint,omitempty"`
	EndSessionEndpoint                        string   `json:"end_session_endpoint,omitempty"`
	ResponseTypesSupported                    []string `json:"response_types_supported,omitempty"`
	ResponseModesSupported                    []string `json:"response_modes_supported,omitempty"`
	GrantTypesSupported                       []string `json:"grant_types_supported,omitempty"`
	SubjectTypesSupported                     []string `json:"subject_types_supported,omitempty"`
	IDTokenSigningAlgValuesSupported          []string `json:"id_token_signing_alg_values_supported,omitempty"`
	ScopesSupported                           []string `json:"scopes_supported,omitempty"`
	ClaimsSupported                           []string `json:"claims_supported,omitempty"`
	CodeChallengeMethodsSupported             []string `json:"code_challenge_methods_supported,omitempty"`
	TokenEndpointAuthMethodsSupported         []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	IntrospectionEndpointAuthMethodsSupported []string `json:"introspection_endpoint_auth_methods_supported,omitempty"`
	RevocationEndpointAuthMethodsSupported    []string `json:"revocation_endpoint_auth_methods_supported,omitempty"`
	BackchannelTokenDeliveryModesSupported    []string `json:"backchannel_token_delivery_modes_supported,omitempty"`
	ActorTypesSupported                       []string `json:"actor_types_supported,omitempty"`
	ResourceIndicatorsSupported               bool     `json:"resource_indicators_supported,omitempty"`
}

// Discover fetches Qeet ID's OpenID Connect discovery document.
//
// The endpoint is public: the client sends no credentials with this request.
//
// API: GET /.well-known/openid-configuration
func (s *OAuthService) Discover(ctx context.Context, opts ...RequestOption) (*ProviderMetadata, error) {
	return call[ProviderMetadata](ctx, s.c, &request{method: http.MethodGet, path: "/.well-known/openid-configuration", public: true}, opts)
}

// token posts a grant to the interactive token endpoint.
func (s *OAuthService) token(ctx context.Context, clientID, clientSecret string, form url.Values, opts []RequestOption) (*Token, error) {
	req, err := oauthRequest("/v1/oauth/token-code", clientID, clientSecret, form)
	if err != nil {
		return nil, err
	}
	tok, err := call[Token](ctx, s.c, req, opts)
	if err != nil {
		return nil, err
	}
	tok.received = time.Now()
	return tok, nil
}

// oauthRequest builds a form POST that authenticates the application: with
// HTTP Basic when it has a secret, or by client_id in the form when it is a
// public client.
func oauthRequest(path, clientID, clientSecret string, form url.Values) (*request, error) {
	if clientID == "" {
		return nil, errors.New("qeetid: OAuth: ClientID is required")
	}
	req := &request{method: http.MethodPost, path: path, form: form, public: true}
	if clientSecret != "" {
		req.clientID, req.clientSecret = clientID, clientSecret
	} else {
		form.Set("client_id", clientID)
	}
	return req, nil
}

func setIf(form url.Values, key, value string) {
	if value != "" {
		form.Set(key, value)
	}
}

// stringList decodes a JSON string or array of strings.
type stringList []string

// UnmarshalJSON accepts a string or an array of strings.
func (l *stringList) UnmarshalJSON(b []byte) error {
	var one string
	if err := json.Unmarshal(b, &one); err == nil {
		if one != "" {
			*l = []string{one}
		}
		return nil
	}
	var many []string
	if err := json.Unmarshal(b, &many); err != nil {
		return err
	}
	*l = many
	return nil
}

func unixTime(n json.Number) time.Time {
	if n == "" {
		return time.Time{}
	}
	if secs, err := n.Int64(); err == nil && secs > 0 {
		return time.Unix(secs, 0)
	}
	if f, err := n.Float64(); err == nil && f > 0 {
		return time.Unix(int64(f), 0)
	}
	return time.Time{}
}
