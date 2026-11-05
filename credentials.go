package qeetid

import (
	"context"
	"net/http"
	"time"
)

// CredentialsService issues, lists, revokes and verifies W3C verifiable credentials.
type CredentialsService struct{ c *Client }

// Issue issues a credential as a signed JWT.
//
// API: POST /v1/tenants/{tenantID}/credentials
func (s *CredentialsService) Issue(ctx context.Context, tenantID string, params IssueCredentialParams, opts ...RequestOption) (*IssuedCredential, error) {
	return call[IssuedCredential](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/credentials",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// List lists a tenant's credentials.
//
// API: GET /v1/tenants/{tenantID}/credentials
func (s *CredentialsService) List(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[Credential], error) {
	return call[Page[Credential]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/credentials",
		pathArgs: []string{tenantID},
	}, opts)
}

// Revoke revokes a credential.
//
// API: POST /v1/tenants/{tenantID}/credentials/{id}/revoke
func (s *CredentialsService) Revoke(ctx context.Context, tenantID string, credentialID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/credentials/{credentialID}/revoke",
		pathArgs: []string{tenantID, credentialID},
	}, nil, opts)
}

// Verify checks a presented credential's signature, expiry and revocation.
//
// The endpoint is public: the client sends no credentials with this request.
//
// API: POST /v1/credentials/verify
func (s *CredentialsService) Verify(ctx context.Context, params VerifyCredentialParams, opts ...RequestOption) (*CredentialVerification, error) {
	return call[CredentialVerification](ctx, s.c, &request{method: http.MethodPost, path: "/v1/credentials/verify", body: params, public: true}, opts)
}

// IssueCredentialParams is the request body of CredentialsService.Issue.
type IssueCredentialParams struct {
	Subject    string         `json:"subject,omitempty"`
	Type       string         `json:"type,omitempty"`
	Claims     map[string]any `json:"claims,omitempty"`
	TTLSeconds *int           `json:"ttl_seconds,omitempty"`
	// The person the credential is for, a member of the organisation (else 422 vc.user_not_member).
	// Recorded with the credential, not in the signed JWT-VC.
	UserID string `json:"user_id,omitempty"`
}

// IssuedCredential carries a newly issued credential as a signed JWT.
type IssuedCredential struct {
	CredentialID string    `json:"credential_id,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitzero"`
	JWT          string    `json:"jwt,omitempty"`
}

// Credential is a W3C verifiable credential a tenant issued.
type Credential struct {
	ExpiresAt time.Time `json:"expires_at,omitzero"`
	ID        string    `json:"id,omitempty"`
	IssuedAt  time.Time `json:"issued_at,omitzero"`
	Revoked   bool      `json:"revoked"`
	Subject   string    `json:"subject,omitempty"`
	Type      string    `json:"type,omitempty"`
	// The person the credential was issued to; null when the issuer named no one, or the account was
	// deleted.
	User *CredentialUser `json:"user,omitempty"`
}

// CredentialUser is the person the credential was issued to; null when the issuer named no one, or the
// account was deleted.
type CredentialUser struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// VerifyCredentialParams is the request body of CredentialsService.Verify.
type VerifyCredentialParams struct {
	// The JWT-VC string.
	Credential string `json:"credential,omitempty"`
}

// CredentialVerification is the outcome of checking a presented credential.
type CredentialVerification struct {
	Issuer string `json:"issuer,omitempty"`
	// Why an invalid credential is refused, in English. Show reason_code's meaning instead.
	Reason string `json:"reason,omitempty"`
	// Why an invalid credential is refused: not a JWT-VC, a bad signature or expired, not one this server
	// issued, or revoked.
	ReasonCode CredentialVerificationReasonCode `json:"reason_code,omitempty"`
	Subject    string                           `json:"subject,omitempty"`
	Valid      bool                             `json:"valid"`
	VC         map[string]any                   `json:"vc,omitempty"`
}

// CredentialVerificationReasonCode is a string enumeration used by CredentialVerification.ReasonCode.
// Why an invalid credential is refused: not a JWT-VC, a bad signature or expired, not one this server
// issued, or revoked.
type CredentialVerificationReasonCode string

// Values of CredentialVerificationReasonCode. The server may add values; handle an unknown one gracefully.
const (
	CredentialVerificationReasonCodeVCNotCredential    CredentialVerificationReasonCode = "vc.not_credential"
	CredentialVerificationReasonCodeVCSignatureInvalid CredentialVerificationReasonCode = "vc.signature_invalid"
	CredentialVerificationReasonCodeVCUnknown          CredentialVerificationReasonCode = "vc.unknown"
	CredentialVerificationReasonCodeVCRevoked          CredentialVerificationReasonCode = "vc.revoked"
)
