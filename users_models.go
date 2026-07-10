package qeetid

type User struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenant_id,omitempty"`
	Email       string         `json:"email"`
	DisplayName string         `json:"display_name,omitempty"`
	Status      string         `json:"status"`
	Phone       string         `json:"phone,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
}

// ListParams narrows a List call. All fields optional.
type ListParams struct {
	Tenant string
	Limit  int
	Cursor string
}

type UserPage struct {
	Data       []User
	NextCursor string
}

// BulkImportSource selects the vendor export format for BulkImport — it
// dictates both the expected Content-Type and the raw body shape.
type BulkImportSource string

const (
	// BulkImportAuth0 expects NDJSON (Auth0's user-export format).
	BulkImportAuth0 BulkImportSource = "auth0"
	// BulkImportCognito expects the Cognito User Import Job CSV template.
	BulkImportCognito BulkImportSource = "cognito"
	// BulkImportAzureB2C expects a Microsoft Graph /users list response
	// ({"value": [...]}).
	BulkImportAzureB2C BulkImportSource = "azure_b2c"
)
