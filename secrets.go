package qeetid

import (
	"context"
	"net/http"
	"time"
)

// SecretsService manages an organization's encrypted secrets vault.
type SecretsService struct{ c *Client }

// List lists an organization's secrets without their values. Page.Truncated reports a list cut short by the
// server's cap.
//
// API: GET /v1/organizations/{organizationID}/secrets
func (s *SecretsService) List(ctx context.Context, organizationID string, opts ...RequestOption) (*Page[Secret], error) {
	return call[Page[Secret]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/secrets",
		pathArgs: []string{organizationID},
	}, opts)
}

// Create stores a secret.
//
// API: POST /v1/organizations/{organizationID}/secrets
func (s *SecretsService) Create(ctx context.Context, organizationID string, params CreateSecretParams, opts ...RequestOption) (*Secret, error) {
	return call[Secret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/secrets",
		pathArgs: []string{organizationID},
		body:     params,
	}, opts)
}

// Update changes a secret's value, scope or rotation interval.
//
// API: PATCH /v1/organizations/{organizationID}/secrets/{id}
func (s *SecretsService) Update(ctx context.Context, organizationID, secretID string, params UpdateSecretParams, opts ...RequestOption) (*Secret, error) {
	return call[Secret](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/organizations/{organizationID}/secrets/{secretID}",
		pathArgs: []string{organizationID, secretID},
		body:     params,
	}, opts)
}

// Delete deletes a secret.
//
// API: DELETE /v1/organizations/{organizationID}/secrets/{id}
func (s *SecretsService) Delete(ctx context.Context, organizationID, secretID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}/secrets/{secretID}",
		pathArgs: []string{organizationID, secretID},
	}, nil, opts)
}

// Reveal returns a secret's value. Every reveal is audited.
//
// API: POST /v1/organizations/{organizationID}/secrets/{id}/reveal
func (s *SecretsService) Reveal(ctx context.Context, organizationID, secretID string, opts ...RequestOption) (*RevealedSecret, error) {
	return call[RevealedSecret](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/organizations/{organizationID}/secrets/{secretID}/reveal",
		pathArgs: []string{organizationID, secretID},
	}, opts)
}

// GetByName returns the value of a secret by name, for the caller's organization.
//
// API: GET /v1/vault/{name}
func (s *SecretsService) GetByName(ctx context.Context, name string, opts ...RequestOption) (*NamedSecret, error) {
	return call[NamedSecret](ctx, s.c, &request{method: http.MethodGet, path: "/v1/vault/{name}", pathArgs: []string{name}}, opts)
}

// Secret is a value stored in an organization's encrypted vault, without the value itself.
type Secret struct {
	CreatedAt time.Time `json:"created_at,omitzero"`
	ID        string    `json:"id,omitempty"`
	Last4     string    `json:"last4,omitempty"`
	Name      string    `json:"name,omitempty"`
	Scope     string    `json:"scope,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	// When the value was last set; a scope or interval change leaves it.
	ValueChangedAt time.Time  `json:"value_changed_at,omitzero"`
	LastRevealedAt *time.Time `json:"last_revealed_at,omitempty"`
	// The person who last revealed the secret; null when no one has or a machine credential did. Email is
	// empty once the person no longer exists.
	LastRevealedBy *SecretLastRevealedBy `json:"last_revealed_by,omitempty"`
	// The id of the machine credential, such as an API key, that last revealed the secret; left out when
	// a person did.
	LastRevealedCredential string `json:"last_revealed_credential,omitempty"`
	RotationIntervalDays   *int   `json:"rotation_interval_days,omitempty"`
	// The rotation interval after the value last changed; null without an interval.
	RotationDueAt   *time.Time `json:"rotation_due_at,omitempty"`
	RotationOverdue bool       `json:"rotation_overdue"`
}

// SecretLastRevealedBy is the person who last revealed the secret; null when no one has or a machine
// credential did. Email is empty once the person no longer exists.
type SecretLastRevealedBy struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// CreateSecretParams is the request body of SecretsService.Create.
type CreateSecretParams struct {
	Name  string `json:"name,omitempty"`
	Scope string `json:"scope,omitempty"`
	Value string `json:"value,omitempty"`
	// Days after the value last changed that the secret is due to be rotated; null or left out for no
	// reminder.
	RotationIntervalDays *int `json:"rotation_interval_days,omitempty"`
}

// UpdateSecretParams is the request body of SecretsService.Update.
type UpdateSecretParams struct {
	Scope *string `json:"scope,omitempty"`
	Value *string `json:"value,omitempty"`
	// Left out keeps the interval; null removes it.
	RotationIntervalDays *int `json:"rotation_interval_days,omitempty"`
}

// RevealedSecret is the result of SecretsService.Reveal.
type RevealedSecret struct {
	Value string `json:"value,omitempty"`
}

// NamedSecret is the result of SecretsService.GetByName.
type NamedSecret struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
