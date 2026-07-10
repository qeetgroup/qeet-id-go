package qeetid

// Organization is a tenant in Qeet ID terms — the wire path is still
// /v1/tenants; "Organization" is the SDK-facing name, matching the term the
// rest of the CIAM industry (Auth0, WorkOS, Clerk) uses for the same concept.
type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Region    string `json:"region,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateOrganizationInput struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Region string `json:"region,omitempty"`
}

type UpdateOrganizationInput struct {
	Name   *string `json:"name,omitempty"`
	Region *string `json:"region,omitempty"`
}

type OrganizationPage struct {
	Data       []Organization
	NextCursor string
}
