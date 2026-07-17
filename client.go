// Package qeetid is the server-side Go SDK for Qeet ID: manage users,
// organizations, roles, and every other management resource; run
// authorization checks; and verify sessions/JWTs.
//
// Authenticate with a secret API key (`qk_…`); never embed it in client code.
// The package has no third-party dependencies — only the standard library.
//
//	client := qeetid.New(qeetid.Config{APIKey: os.Getenv("QEETID_API_KEY")})
//	claims, err := client.Sessions.Verify(ctx, token)
//	ok, err := client.Permissions.Check(ctx, qeetid.PermissionCheck{
//		User: claims.UserID, Tenant: claims.TenantID, Permission: "billing:write",
//	})
//
// Every resource sits directly on Client as a *XService field — comment
// banners below group them (Identity / Authentication / Authorization /
// Administration) for documentation only; there is no nesting, so every
// resource is one field access away: client.Users, client.Sessions,
// client.Webhooks.
package qeetid

import "github.com/qeetgroup/qeet-id-go/internal/transport"

// Client is the Qeet ID API client. Construct once with New and reuse it; it
// is safe for concurrent use.
type Client struct {
	// Identity — who exists: human users, organizations, and machine identities.
	Users             *UsersService
	Organizations     *OrganizationsService
	ServicePrincipals *ServicePrincipalsService
	Agents            *AgentsService
	Domains           *DomainsService

	// Authentication — proving who's calling: sessions, federation protocols,
	// and their supporting policy.
	Sessions      *SessionsService
	OAuth         *OAuthService
	OIDC          *OIDCService
	SAML          *SAMLService
	SAMLProviders *SAMLServiceProvidersService
	SCIM          *SCIMService
	LDAP          *LDAPService
	Social        *SocialService
	MFA           *MFAService
	Credentials   *CredentialsService
	AuthHooks     *AuthHooksService
	AuthPolicy    *AuthPolicyService
	Policy        *PolicyService
	IPRules       *IPRulesService
	BotDetection  *BotDetectionService
	RiskSettings  *RiskSettingsService

	// Authorization — what an authenticated caller is allowed to do.
	Roles         *RolesService
	Permissions   *PermissionsService
	Groups        *GroupsService
	Relationships *RelationTuplesService
	Decisions     *AuthZENService

	// Administration — tenant operations: branding, developer tooling,
	// compliance, and billing for the Qeet ID platform itself.
	Branding       *BrandingService
	Invitations    *InvitationsService
	EmailTemplates *EmailTemplatesService
	APIKeys        *APIKeysService
	Vault          *VaultService
	TokenVault     *TokenVaultService
	Webhooks       *WebhooksService
	AuditLogs      *AuditLogsService
	Analytics      *AnalyticsService
	GDPR           *GDPRService
	Billing        *BillingService
	Retention      *RetentionService
	RateLimits     *RateLimitsService
	LogSinks       *LogSinksService
	AdminLinks     *AdminLinksService

	transport *transport.Transport
}

// New builds a client. APIKey is required.
func New(cfg Config) *Client {
	t := transport.New(transport.Options{
		APIKey:     cfg.APIKey,
		BaseURL:    cfg.BaseURL,
		HTTPClient: cfg.HTTPClient,
		Timeout:    cfg.Timeout,
		MaxRetries: cfg.MaxRetries,
		Headers:    cfg.Headers,
		UserAgent:  cfg.UserAgent,
		Logger:     cfg.Logger,
	})
	return &Client{
		transport: t,

		Users:             &UsersService{t: t},
		Organizations:     &OrganizationsService{t: t},
		ServicePrincipals: &ServicePrincipalsService{t: t},
		Agents:            &AgentsService{t: t},
		Domains:           &DomainsService{t: t},

		Sessions:      newSessionsService(t.BaseURL(), t.HTTPClient()),
		OAuth:         newOAuthService(t),
		OIDC:          &OIDCService{t: t},
		SAML:          &SAMLService{t: t},
		SAMLProviders: &SAMLServiceProvidersService{t: t},
		SCIM:          &SCIMService{t: t},
		LDAP:          &LDAPService{t: t},
		Social:        &SocialService{t: t},
		MFA:           &MFAService{t: t},
		Credentials:   &CredentialsService{t: t},
		AuthHooks:     &AuthHooksService{t: t},
		AuthPolicy:    &AuthPolicyService{t: t},
		Policy:        &PolicyService{t: t},
		IPRules:       &IPRulesService{t: t},
		BotDetection:  &BotDetectionService{t: t},
		RiskSettings:  &RiskSettingsService{t: t},

		Roles:         &RolesService{t: t},
		Permissions:   &PermissionsService{t: t},
		Groups:        &GroupsService{t: t},
		Relationships: &RelationTuplesService{t: t},
		Decisions:     &AuthZENService{t: t},

		Branding:       &BrandingService{t: t},
		Invitations:    &InvitationsService{t: t},
		EmailTemplates: &EmailTemplatesService{t: t},
		APIKeys:        &APIKeysService{t: t},
		Vault:          &VaultService{t: t},
		TokenVault:     &TokenVaultService{t: t},
		Webhooks:       &WebhooksService{t: t},
		AuditLogs:      &AuditLogsService{t: t},
		Analytics:      &AnalyticsService{t: t},
		GDPR:           &GDPRService{t: t},
		Billing:        &BillingService{t: t},
		Retention:      &RetentionService{t: t},
		RateLimits:     &RateLimitsService{t: t},
		LogSinks:       &LogSinksService{t: t},
		AdminLinks:     &AdminLinksService{t: t},
	}
}
