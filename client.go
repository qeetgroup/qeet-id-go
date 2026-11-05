package qeetid

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/qeetgroup/qeet-id-go/internal/version"
)

// Client is a client for the Qeet ID API. Create one with NewClient and share
// it: a Client is safe for concurrent use, and reusing it reuses connections
// and cached tokens.
//
// Each field groups the operations on one kind of resource, as in
// client.Users.Get(ctx, userID). The fields are set by NewClient and must not
// be changed.
type Client struct {
	// Tenants and the people, groups and invitations in them.
	Tenants     *TenantsService
	Users       *UsersService
	Groups      *GroupsService
	Invitations *InvitationsService
	Domains     *DomainsService
	Branding    *BrandingService

	// Access control: roles, permissions, attribute- and relationship-based
	// rules, and AuthZEN decisions.
	Roles         *RolesService
	Permissions   *PermissionsService
	ABACPolicies  *ABACPoliciesService
	Relationships *RelationshipsService
	AuthZEN       *AuthZENService

	// Sign-in policy and threat protection.
	AuthPolicy     *AuthPolicyService
	SecurityPolicy *SecurityPolicyService
	IPRules        *IPRulesService
	Threats        *ThreatsService

	// Applications, machine identities and federation.
	OIDCClients          *OIDCClientsService
	OAuthResources       *OAuthResourcesService
	OAuthGrants          *OAuthGrantsService
	DeviceAuthorizations *DeviceAuthorizationsService
	SigningKeys          *SigningKeysService
	ServicePrincipals    *ServicePrincipalsService
	SocialProviders      *SocialProvidersService
	SAMLConnections      *SAMLConnectionsService
	SAMLServiceProviders *SAMLServiceProvidersService
	OIDCConnections      *OIDCConnectionsService
	LDAP                 *LDAPService
	Directory            *DirectoryService
	SCIM                 *SCIMService
	AdminPortalLinks     *AdminPortalLinksService

	// Developer resources.
	APIKeys     *APIKeysService
	Agents      *AgentsService
	Credentials *CredentialsService
	AuthHooks   *AuthHooksService
	Webhooks    *WebhooksService
	Secrets     *SecretsService
	TokenVault  *TokenVaultService

	// Operations: audit, compliance, data lifecycle, email and billing.
	Audit          *AuditService
	Analytics      *AnalyticsService
	Compliance     *ComplianceService
	Retention      *RetentionService
	GDPR           *GDPRService
	LogSinks       *LogSinksService
	RateLimits     *RateLimitsService
	EmailTemplates *EmailTemplatesService
	EmailDelivery  *EmailDeliveryService
	Billing        *BillingService
	Entitlements   *EntitlementsService

	baseURL       string
	hc            *http.Client
	cred          credential
	timeout       time.Duration
	maxRetries    int
	minRetryDelay time.Duration
	maxRetryDelay time.Duration
	logger        *slog.Logger
	userAgent     string
	headers       http.Header
}

// NewClient returns a Client configured by opts. Give it one credential:
// WithAPIKey, WithClientCredentials or WithTokenSource. A client without one
// can still call the public endpoints, such as token verification for verifiable
// credentials and the OAuth endpoints that authenticate the calling
// application.
//
//	client, err := qeetid.NewClient(qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")))
//
// NewClient makes no network calls. It returns an error when an option is
// invalid.
func NewClient(opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		if o != nil {
			o.applyOption(&cfg)
		}
	}
	if len(cfg.credentials) > 1 {
		names := make([]string, len(cfg.credentials))
		for i, cr := range cfg.credentials {
			names[i] = cr.name
		}
		cfg.errs = append(cfg.errs, fmt.Errorf("qeetid: use one credential option, got %s", strings.Join(names, ", ")))
	}
	base, err := normalizeBaseURL(cfg.baseURL)
	if err != nil {
		cfg.errs = append(cfg.errs, err)
	}
	if err := errors.Join(cfg.errs...); err != nil {
		return nil, err
	}

	c := &Client{
		baseURL:       base,
		hc:            cfg.httpClient,
		timeout:       cfg.timeout,
		maxRetries:    cfg.maxRetries,
		minRetryDelay: cfg.minRetryDelay,
		maxRetryDelay: cfg.maxRetryDelay,
		logger:        cfg.logger,
		userAgent:     version.UserAgent(),
		headers:       cfg.headers,
	}
	if c.hc == nil {
		c.hc = &http.Client{}
	}
	if cfg.userAgentSuffix != "" {
		c.userAgent += " " + cfg.userAgentSuffix
	}
	if len(cfg.credentials) == 1 {
		switch cr := cfg.credentials[0]; {
		case cr.apiKey != "":
			c.cred = apiKeyCredential{key: cr.apiKey}
		case cr.ts != nil:
			c.cred = bearerCredential{ts: cr.ts}
		case cr.cc != nil:
			c.cred = bearerCredential{ts: newClientCredentialsSource(c, cr.cc)}
		}
	}

	c.Tenants = &TenantsService{c: c}
	c.Users = &UsersService{c: c}
	c.Groups = &GroupsService{c: c}
	c.Invitations = &InvitationsService{c: c}
	c.Domains = &DomainsService{c: c}
	c.Branding = &BrandingService{c: c}

	c.Roles = &RolesService{c: c}
	c.Permissions = &PermissionsService{c: c}
	c.ABACPolicies = &ABACPoliciesService{c: c}
	c.Relationships = &RelationshipsService{c: c}
	c.AuthZEN = &AuthZENService{c: c}

	c.AuthPolicy = &AuthPolicyService{c: c}
	c.SecurityPolicy = &SecurityPolicyService{c: c}
	c.IPRules = &IPRulesService{c: c}
	c.Threats = &ThreatsService{c: c}

	c.OIDCClients = &OIDCClientsService{c: c}
	c.OAuthResources = &OAuthResourcesService{c: c}
	c.OAuthGrants = &OAuthGrantsService{c: c}
	c.DeviceAuthorizations = &DeviceAuthorizationsService{c: c}
	c.SigningKeys = &SigningKeysService{c: c}
	c.ServicePrincipals = &ServicePrincipalsService{c: c}
	c.SocialProviders = &SocialProvidersService{c: c}
	c.SAMLConnections = &SAMLConnectionsService{c: c}
	c.SAMLServiceProviders = &SAMLServiceProvidersService{c: c}
	c.OIDCConnections = &OIDCConnectionsService{c: c}
	c.LDAP = &LDAPService{c: c}
	c.Directory = &DirectoryService{c: c}
	c.SCIM = &SCIMService{c: c}
	c.AdminPortalLinks = &AdminPortalLinksService{c: c}

	c.APIKeys = &APIKeysService{c: c}
	c.Agents = &AgentsService{c: c}
	c.Credentials = &CredentialsService{c: c}
	c.AuthHooks = &AuthHooksService{c: c}
	c.Webhooks = &WebhooksService{c: c}
	c.Secrets = &SecretsService{c: c}
	c.TokenVault = &TokenVaultService{c: c}

	c.Audit = &AuditService{c: c}
	c.Analytics = &AnalyticsService{c: c}
	c.Compliance = &ComplianceService{c: c}
	c.Retention = &RetentionService{c: c}
	c.GDPR = &GDPRService{c: c}
	c.LogSinks = &LogSinksService{c: c}
	c.RateLimits = &RateLimitsService{c: c}
	c.EmailTemplates = &EmailTemplatesService{c: c}
	c.EmailDelivery = &EmailDeliveryService{c: c}
	c.Billing = &BillingService{c: c}
	c.Entitlements = &EntitlementsService{c: c}
	return c, nil
}

// BaseURL returns the API base URL the client sends requests to.
func (c *Client) BaseURL() string { return c.baseURL }

// normalizeBaseURL checks that raw is an absolute https URL (or http on a
// loopback host) without a query or fragment, and trims a trailing slash.
func normalizeBaseURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("qeetid: WithBaseURL: %w", err)
	}
	if u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return "", fmt.Errorf("qeetid: WithBaseURL: %q is not an absolute http(s) URL", raw)
	}
	if u.Scheme == "http" && !isLoopback(u.Hostname()) {
		return "", fmt.Errorf("qeetid: WithBaseURL: %q must use https; http is allowed only for a loopback host", raw)
	}
	if u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return "", fmt.Errorf("qeetid: WithBaseURL: %q must not carry credentials, a query or a fragment", raw)
	}
	return strings.TrimRight(u.String(), "/"), nil
}

func isLoopback(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
