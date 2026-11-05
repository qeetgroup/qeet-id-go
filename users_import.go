package qeetid

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// ImportSource names the identity provider a user export file comes from. It
// decides the format UsersService.Import expects.
type ImportSource string

// Supported import sources.
const (
	// ImportSourceAuth0 expects Auth0's user export as NDJSON, one user per
	// line (Content-Type application/x-ndjson).
	ImportSourceAuth0 ImportSource = "auth0"
	// ImportSourceCognito expects Amazon Cognito's user import CSV template
	// (Content-Type text/csv).
	ImportSourceCognito ImportSource = "cognito"
	// ImportSourceAzureB2C expects a Microsoft Graph /users list response,
	// {"value": [...]} (Content-Type application/json).
	ImportSourceAzureB2C ImportSource = "azure_b2c"
)

// contentType is the media type the API expects for the source's file.
func (s ImportSource) contentType() (string, bool) {
	switch s {
	case ImportSourceAuth0:
		return "application/x-ndjson", true
	case ImportSourceCognito:
		return "text/csv", true
	case ImportSourceAzureB2C:
		return "application/json", true
	}
	return "", false
}

// Import creates users from another identity provider's user export file,
// sent as is. Rows are imported one by one: the result reports how many
// succeeded and why each failed row was refused. No provider exports a
// portable password hash, so imported users have no password; they sign in
// with a passkey, a magic link, a one-time code or a reset.
//
// API: POST /v1/users/bulk/import
func (s *UsersService) Import(ctx context.Context, source ImportSource, file []byte, opts ...RequestOption) (*BulkImportResult, error) {
	contentType, ok := source.contentType()
	if !ok {
		return nil, errors.New("qeetid: Import: unknown source " + string(source))
	}
	if len(file) == 0 {
		return nil, errors.New("qeetid: Import: the file is empty")
	}
	return call[BulkImportResult](ctx, s.c, &request{
		method:      http.MethodPost,
		path:        "/v1/users/bulk/import",
		query:       url.Values{"source": {string(source)}},
		rawBody:     file,
		contentType: contentType,
	}, opts)
}
