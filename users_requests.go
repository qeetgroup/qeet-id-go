package qeetid

type CreateUserInput struct {
	Email       string         `json:"email"`
	DisplayName string         `json:"display_name,omitempty"`
	Phone       string         `json:"phone,omitempty"`
	Password    string         `json:"password,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type UpdateUserInput struct {
	DisplayName *string        `json:"display_name,omitempty"`
	Phone       *string        `json:"phone,omitempty"`
	Status      *string        `json:"status,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// BulkUserInput is one row in a BulkCreate call.
type BulkUserInput struct {
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

// BulkCreateInput is the request body for POST /v1/users/bulk.
type BulkCreateInput struct {
	TenantID string          `json:"tenant_id,omitempty"`
	Users    []BulkUserInput `json:"users"` // 1-1000 items
}

// VerifyEmailStartInput optionally overrides the email a verification code
// is sent to (defaults to the user's current email).
type VerifyEmailStartInput struct {
	Email string `json:"email,omitempty"`
}

// VerifyEmailConfirmInput is the code from VerifyEmailStart.
type VerifyEmailConfirmInput struct {
	Code string `json:"code"`
}

// VerifyPhoneStartInput optionally overrides the phone a verification code
// is sent to (defaults to the user's current phone).
type VerifyPhoneStartInput struct {
	Phone string `json:"phone,omitempty"`
}

// VerifyPhoneConfirmInput is the code from VerifyPhoneStart.
type VerifyPhoneConfirmInput struct {
	Code string `json:"code"`
}
