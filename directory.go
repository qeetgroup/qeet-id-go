package qeetid

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DirectoryService manages directory sync: profile mappings, errors, activity and failure alerts.
type DirectoryService struct{ c *Client }

// ListConnections lists the tenant's directory connections.
//
// API: GET /v1/tenants/{tenantID}/directory/connections
func (s *DirectoryService) ListConnections(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[DirectoryConnection], error) {
	return call[Page[DirectoryConnection]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/directory/connections",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListMappings lists a connection's effective profile mapping rules and the fields they can target.
//
// API: GET /v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings
func (s *DirectoryService) ListMappings(ctx context.Context, tenantID string, connectionID string, opts ...RequestOption) (*DirectoryMappings, error) {
	return call[DirectoryMappings](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings",
		pathArgs: []string{tenantID, connectionID},
	}, opts)
}

// CreateMapping adds a profile mapping rule.
//
// API: POST /v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings
func (s *DirectoryService) CreateMapping(ctx context.Context, tenantID string, connectionID string, params DirectoryMapping, opts ...RequestOption) (*DirectoryMapping, error) {
	return call[DirectoryMapping](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// UpdateMapping replaces a profile mapping rule. Version must be the version last read.
//
// API: PUT /v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings/{mappingID}
func (s *DirectoryService) UpdateMapping(ctx context.Context, tenantID string, connectionID string, mappingID string, params DirectoryMapping, opts ...RequestOption) (*DirectoryMapping, error) {
	return call[DirectoryMapping](ctx, s.c, &request{
		method:   http.MethodPut,
		path:     "/v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings/{mappingID}",
		pathArgs: []string{tenantID, connectionID, mappingID},
		body:     params,
	}, opts)
}

// DeleteMapping deletes a profile mapping rule.
//
// API: DELETE /v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings/{mappingID}
func (s *DirectoryService) DeleteMapping(ctx context.Context, tenantID string, connectionID string, mappingID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings/{mappingID}",
		pathArgs: []string{tenantID, connectionID, mappingID},
	}, nil, opts)
}

// PreviewMappings applies the saved rules to sample directory data without changing anything.
//
// API: POST /v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings/preview
func (s *DirectoryService) PreviewMappings(ctx context.Context, tenantID string, connectionID string, params PreviewMappingsParams, opts ...RequestOption) (*MappingPreview, error) {
	return call[MappingPreview](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/directory/connections/{connectionID}/mappings/preview",
		pathArgs: []string{tenantID, connectionID},
		body:     params,
	}, opts)
}

// ListErrors lists directory sync errors with a remediation summary.
//
// API: GET /v1/tenants/{tenantID}/directory/errors
func (s *DirectoryService) ListErrors(ctx context.Context, tenantID string, params *ListDirectoryErrorsParams, opts ...RequestOption) (*DirectoryErrors, error) {
	return call[DirectoryErrors](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/directory/errors",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// ResolveError marks a directory error resolved.
//
// API: POST /v1/tenants/{tenantID}/directory/errors/{eventID}/resolve
func (s *DirectoryService) ResolveError(ctx context.Context, tenantID string, eventID string, opts ...RequestOption) (*DirectoryActionResult, error) {
	return call[DirectoryActionResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/directory/errors/{eventID}/resolve",
		pathArgs: []string{tenantID, eventID},
	}, opts)
}

// RetryError retries the failed connection check behind a directory error.
//
// API: POST /v1/tenants/{tenantID}/directory/errors/{eventID}/retry
func (s *DirectoryService) RetryError(ctx context.Context, tenantID string, eventID string, opts ...RequestOption) (*DirectoryActionResult, error) {
	return call[DirectoryActionResult](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/directory/errors/{eventID}/retry",
		pathArgs: []string{tenantID, eventID},
	}, opts)
}

// ListActivity lists recorded directory operations.
//
// API: GET /v1/tenants/{tenantID}/directory/activity
func (s *DirectoryService) ListActivity(ctx context.Context, tenantID string, params *ListDirectoryActivityParams, opts ...RequestOption) (*DirectoryActivity, error) {
	return call[DirectoryActivity](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/directory/activity",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// ListAlerts lists directory failure alerts.
//
// API: GET /v1/tenants/{tenantID}/directory/alerts
func (s *DirectoryService) ListAlerts(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[DirectoryAlert], error) {
	return call[Page[DirectoryAlert]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/directory/alerts",
		pathArgs: []string{tenantID},
	}, opts)
}

// CreateAlert creates a directory failure alert.
//
// API: POST /v1/tenants/{tenantID}/directory/alerts
func (s *DirectoryService) CreateAlert(ctx context.Context, tenantID string, params CreateDirectoryAlertParams, opts ...RequestOption) (*DirectoryAlert, error) {
	return call[DirectoryAlert](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/directory/alerts",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// UpdateAlert changes or pauses a directory failure alert.
//
// API: PATCH /v1/tenants/{tenantID}/directory/alerts/{alertID}
func (s *DirectoryService) UpdateAlert(ctx context.Context, tenantID string, alertID string, params UpdateDirectoryAlertParams, opts ...RequestOption) (*DirectoryAlert, error) {
	return call[DirectoryAlert](ctx, s.c, &request{
		method:   http.MethodPatch,
		path:     "/v1/tenants/{tenantID}/directory/alerts/{alertID}",
		pathArgs: []string{tenantID, alertID},
		body:     params,
	}, opts)
}

// DeleteAlert removes a directory failure alert.
//
// API: DELETE /v1/tenants/{tenantID}/directory/alerts/{alertID}
func (s *DirectoryService) DeleteAlert(ctx context.Context, tenantID string, alertID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/directory/alerts/{alertID}",
		pathArgs: []string{tenantID, alertID},
	}, nil, opts)
}

// Profile returns the directory attributes behind a person's profile.
//
// API: GET /v1/tenants/{tenantID}/directory/profiles/{userID}
func (s *DirectoryService) Profile(ctx context.Context, tenantID string, userID string, opts ...RequestOption) (*Page[DirectoryProfileValue], error) {
	return call[Page[DirectoryProfileValue]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/directory/profiles/{userID}",
		pathArgs: []string{tenantID, userID},
	}, opts)
}

// DirectoryConnection is one item of the list DirectoryService.ListConnections returns.
type DirectoryConnection struct {
	ID             string                    `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Reference      string                    `json:"reference,omitempty"`
	Type           DirectoryConnectionType   `json:"type,omitempty"`
	Status         DirectoryConnectionStatus `json:"status,omitempty"`
	Health         DirectoryConnectionHealth `json:"health,omitempty"`
	Environment    string                    `json:"environment,omitempty"`
	Scope          DirectoryConnectionScope  `json:"scope,omitempty"`
	CreatedAt      *time.Time                `json:"created_at,omitempty"`
	LastActivityAt *time.Time                `json:"last_activity_at,omitempty"`
	LastSyncAt     *time.Time                `json:"last_sync_at,omitempty"`
	// SCIM, the people it provisioned; LDAP, the people with an identity this connection signed in (not
	// deleted ones).
	UserCount *int `json:"user_count,omitempty"`
	// SCIM, the groups it provisioned; null for LDAP, which provisions none.
	GroupCount *int   `json:"group_count,omitempty"`
	ServerURL  string `json:"server_url,omitempty"`
}

// DirectoryConnectionType is a string enumeration used by DirectoryConnection.Type.
type DirectoryConnectionType string

// Values of DirectoryConnectionType. The server may add values; handle an unknown one gracefully.
const (
	DirectoryConnectionTypeSCIM DirectoryConnectionType = "scim"
	DirectoryConnectionTypeLDAP DirectoryConnectionType = "ldap"
)

// DirectoryConnectionStatus is a string enumeration used by DirectoryConnection.Status.
type DirectoryConnectionStatus string

// Values of DirectoryConnectionStatus. The server may add values; handle an unknown one gracefully.
const (
	DirectoryConnectionStatusConfigured DirectoryConnectionStatus = "configured"
	DirectoryConnectionStatusActive     DirectoryConnectionStatus = "active"
	DirectoryConnectionStatusDraft      DirectoryConnectionStatus = "draft"
	DirectoryConnectionStatusDisabled   DirectoryConnectionStatus = "disabled"
)

// DirectoryConnectionHealth is a string enumeration used by DirectoryConnection.Health.
type DirectoryConnectionHealth string

// Values of DirectoryConnectionHealth. The server may add values; handle an unknown one gracefully.
const (
	DirectoryConnectionHealthHealthy DirectoryConnectionHealth = "healthy"
	DirectoryConnectionHealthWarning DirectoryConnectionHealth = "warning"
	DirectoryConnectionHealthFailed  DirectoryConnectionHealth = "failed"
	DirectoryConnectionHealthUnknown DirectoryConnectionHealth = "unknown"
)

// DirectoryConnectionScope is a string enumeration used by DirectoryConnection.Scope.
type DirectoryConnectionScope string

// Values of DirectoryConnectionScope. The server may add values; handle an unknown one gracefully.
const (
	DirectoryConnectionScopeUsersAndGroups DirectoryConnectionScope = "users-and-groups"
	DirectoryConnectionScopeUsersOnSignin  DirectoryConnectionScope = "users-on-signin"
)

// DirectoryMappings is the result of DirectoryService.ListMappings.
type DirectoryMappings struct {
	Items   []DirectoryMapping `json:"items,omitempty"`
	Targets map[string]string  `json:"targets,omitempty"`
}

// DirectoryMapping is a rule that maps a directory attribute to a profile field.
type DirectoryMapping struct {
	ID              string                      `json:"id,omitempty"`
	SourceAttribute string                      `json:"source_attribute,omitempty"`
	TargetField     DirectoryMappingTargetField `json:"target_field,omitempty"`
	ValueType       DirectoryMappingValueType   `json:"value_type,omitempty"`
	Transform       DirectoryMappingTransform   `json:"transform,omitempty"`
	Pattern         string                      `json:"pattern,omitempty"`
	Enabled         bool                        `json:"enabled"`
	Required        bool                        `json:"required"`
	Version         int                         `json:"version"`
}

// DirectoryMappingTargetField is a string enumeration used by DirectoryMapping.TargetField.
type DirectoryMappingTargetField string

// Values of DirectoryMappingTargetField. The server may add values; handle an unknown one gracefully.
const (
	DirectoryMappingTargetFieldEmail       DirectoryMappingTargetField = "email"
	DirectoryMappingTargetFieldDisplayName DirectoryMappingTargetField = "display_name"
	DirectoryMappingTargetFieldFirstName   DirectoryMappingTargetField = "first_name"
	DirectoryMappingTargetFieldLastName    DirectoryMappingTargetField = "last_name"
	DirectoryMappingTargetFieldDepartment  DirectoryMappingTargetField = "department"
	DirectoryMappingTargetFieldTitle       DirectoryMappingTargetField = "title"
	DirectoryMappingTargetFieldManagerID   DirectoryMappingTargetField = "manager_id"
	DirectoryMappingTargetFieldGroups      DirectoryMappingTargetField = "groups"
)

// DirectoryMappingValueType is a string enumeration used by DirectoryMapping.ValueType.
type DirectoryMappingValueType string

// Values of DirectoryMappingValueType. The server may add values; handle an unknown one gracefully.
const (
	DirectoryMappingValueTypeString DirectoryMappingValueType = "string"
	DirectoryMappingValueTypeArray  DirectoryMappingValueType = "array"
)

// DirectoryMappingTransform is a string enumeration used by DirectoryMapping.Transform.
type DirectoryMappingTransform string

// Values of DirectoryMappingTransform. The server may add values; handle an unknown one gracefully.
const (
	DirectoryMappingTransformNone      DirectoryMappingTransform = "none"
	DirectoryMappingTransformTrim      DirectoryMappingTransform = "trim"
	DirectoryMappingTransformLowercase DirectoryMappingTransform = "lowercase"
	DirectoryMappingTransformUppercase DirectoryMappingTransform = "uppercase"
	DirectoryMappingTransformExtract   DirectoryMappingTransform = "extract"
)

// PreviewMappingsParams is the request body of DirectoryService.PreviewMappings.
type PreviewMappingsParams struct {
	Source map[string]any `json:"source,omitempty"`
}

// MappingPreview is the result of DirectoryService.PreviewMappings.
type MappingPreview struct {
	Valid   bool           `json:"valid"`
	Profile map[string]any `json:"profile,omitempty"`
	Values  map[string]any `json:"values,omitempty"`
	// Every field the person would get, with its value and where it came from: a mapping, or the
	// connection's own default where no mapping sets it (SCIM userName and displayName; an LDAP
	// connection's email and name attributes, mail and cn unless set).
	Effective map[string]MappingPreviewEffectiveEntry `json:"effective,omitempty"`
	Issues    []MappingPreviewIssue                   `json:"issues,omitempty"`
}

// MappingPreviewEffectiveEntry is a value of the MappingPreview.Effective map.
type MappingPreviewEffectiveEntry struct {
	Value json.RawMessage                  `json:"value,omitempty"`
	From  MappingPreviewEffectiveEntryFrom `json:"from,omitempty"`
	// The source attribute it was read from.
	Attribute string `json:"attribute,omitempty"`
}

// MappingPreviewEffectiveEntryFrom is a string enumeration used by MappingPreviewEffectiveEntry.From.
type MappingPreviewEffectiveEntryFrom string

// Values of MappingPreviewEffectiveEntryFrom. The server may add values; handle an unknown one gracefully.
const (
	MappingPreviewEffectiveEntryFromMapping MappingPreviewEffectiveEntryFrom = "mapping"
	MappingPreviewEffectiveEntryFromDefault MappingPreviewEffectiveEntryFrom = "default"
)

// MappingPreviewIssue is one element of MappingPreview.Issues.
type MappingPreviewIssue struct {
	// The mapping at fault; the nil UUID for directory.mapping_preview_too_large.
	MappingID string `json:"mapping_id,omitempty"`
	// Target profile field of the mapping at fault; empty for directory.mapping_preview_too_large.
	Field string `json:"field,omitempty"`
	// Stable reason for the issue; word it from this.
	Code MappingPreviewIssueCode `json:"code,omitempty"`
	// English text kept for compatibility; it may change.
	Message string `json:"message,omitempty"`
}

// MappingPreviewIssueCode is a string enumeration used by MappingPreviewIssue.Code. Stable reason for
// the issue; word it from this.
type MappingPreviewIssueCode string

// Values of MappingPreviewIssueCode. The server may add values; handle an unknown one gracefully.
const (
	MappingPreviewIssueCodeDirectoryMappingPreviewTooLarge      MappingPreviewIssueCode = "directory.mapping_preview_too_large"
	MappingPreviewIssueCodeDirectoryMappingSourceInvalid        MappingPreviewIssueCode = "directory.mapping_source_invalid"
	MappingPreviewIssueCodeDirectoryMappingSourceCredential     MappingPreviewIssueCode = "directory.mapping_source_credential"
	MappingPreviewIssueCodeDirectoryMappingTargetUnsupported    MappingPreviewIssueCode = "directory.mapping_target_unsupported"
	MappingPreviewIssueCodeDirectoryMappingTargetDuplicate      MappingPreviewIssueCode = "directory.mapping_target_duplicate"
	MappingPreviewIssueCodeDirectoryMappingValueTypeMismatch    MappingPreviewIssueCode = "directory.mapping_value_type_mismatch"
	MappingPreviewIssueCodeDirectoryMappingTransformUnsupported MappingPreviewIssueCode = "directory.mapping_transform_unsupported"
	MappingPreviewIssueCodeDirectoryMappingPatternUnexpected    MappingPreviewIssueCode = "directory.mapping_pattern_unexpected"
	MappingPreviewIssueCodeDirectoryMappingPatternInvalid       MappingPreviewIssueCode = "directory.mapping_pattern_invalid"
	MappingPreviewIssueCodeDirectoryMappingSourceMissing        MappingPreviewIssueCode = "directory.mapping_source_missing"
	MappingPreviewIssueCodeDirectoryMappingSourceEmpty          MappingPreviewIssueCode = "directory.mapping_source_empty"
	MappingPreviewIssueCodeDirectoryMappingValueNotString       MappingPreviewIssueCode = "directory.mapping_value_not_string"
	MappingPreviewIssueCodeDirectoryMappingValueNotArray        MappingPreviewIssueCode = "directory.mapping_value_not_array"
	MappingPreviewIssueCodeDirectoryMappingValueTooLong         MappingPreviewIssueCode = "directory.mapping_value_too_long"
	MappingPreviewIssueCodeDirectoryMappingValueTooMany         MappingPreviewIssueCode = "directory.mapping_value_too_many"
	MappingPreviewIssueCodeDirectoryMappingValueUnmatched       MappingPreviewIssueCode = "directory.mapping_value_unmatched"
	MappingPreviewIssueCodeDirectoryMappingEmailInvalid         MappingPreviewIssueCode = "directory.mapping_email_invalid"
	MappingPreviewIssueCodeDirectoryMappingDisplayNameTooLong   MappingPreviewIssueCode = "directory.mapping_display_name_too_long"
)

// DirectoryEventSeverity is a string enumeration used by ListDirectoryErrorsParams.Severity.
type DirectoryEventSeverity string

// Values of DirectoryEventSeverity. The server may add values; handle an unknown one gracefully.
const (
	DirectoryEventSeverityWarning  DirectoryEventSeverity = "warning"
	DirectoryEventSeverityCritical DirectoryEventSeverity = "critical"
	DirectoryEventSeverityEmpty    DirectoryEventSeverity = ""
)

// DirectoryEventResolution is a string enumeration used by ListDirectoryErrorsParams.Resolution.
type DirectoryEventResolution string

// Values of DirectoryEventResolution. The server may add values; handle an unknown one gracefully.
const (
	DirectoryEventResolutionOpen     DirectoryEventResolution = "open"
	DirectoryEventResolutionResolved DirectoryEventResolution = "resolved"
	DirectoryEventResolutionRetried  DirectoryEventResolution = "retried"
)

// ListDirectoryErrorsParams holds the optional query parameters of DirectoryService.ListErrors. A nil
// *ListDirectoryErrorsParams sends none.
type ListDirectoryErrorsParams struct {
	Days         int
	ConnectionID string
	Severity     DirectoryEventSeverity
	Resolution   DirectoryEventResolution
	Q            string
	Limit        int
	Offset       int
}

func (p *ListDirectoryErrorsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Days != 0 {
		q.Set("days", strconv.FormatInt(int64(p.Days), 10))
	}
	if p.ConnectionID != "" {
		q.Set("connection_id", p.ConnectionID)
	}
	if p.Severity != "" {
		q.Set("severity", string(p.Severity))
	}
	if p.Resolution != "" {
		q.Set("resolution", string(p.Resolution))
	}
	if p.Q != "" {
		q.Set("q", p.Q)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Offset != 0 {
		q.Set("offset", strconv.FormatInt(int64(p.Offset), 10))
	}
	return q
}

// DirectoryErrors is the result of DirectoryService.ListErrors.
type DirectoryErrors struct {
	Items        []DirectoryEvent            `json:"items,omitempty"`
	Total        int                         `json:"total"`
	Limit        int                         `json:"limit"`
	Offset       int                         `json:"offset"`
	Summary      DirectoryErrorsSummary      `json:"summary,omitzero"`
	Categories   []DirectoryErrorsCategory   `json:"categories,omitempty"`
	Connections  []DirectoryErrorsConnection `json:"connections,omitempty"`
	Remediations []DirectoryEvent            `json:"remediations,omitempty"`
}

// DirectoryEvent is a recorded directory operation or error.
type DirectoryEvent struct {
	ID           string                     `json:"id,omitempty"`
	ConnectionID string                     `json:"connection_id,omitempty"`
	Operation    string                     `json:"operation,omitempty"`
	Trigger      DirectoryEventTrigger      `json:"trigger,omitempty"`
	ResourceType DirectoryEventResourceType `json:"resource_type,omitempty"`
	ResourceID   string                     `json:"resource_id,omitempty"`
	StartedAt    time.Time                  `json:"started_at,omitzero"`
	DurationMs   int                        `json:"duration_ms"`
	HTTPStatus   int                        `json:"http_status"`
	Status       DirectoryEventStatus       `json:"status,omitempty"`
	UsersSynced  int                        `json:"users_synced"`
	GroupsSynced int                        `json:"groups_synced"`
	ErrorCode    string                     `json:"error_code,omitempty"`
	Message      string                     `json:"message,omitempty"`
	Category     string                     `json:"category,omitempty"`
	Severity     DirectoryEventSeverity     `json:"severity,omitempty"`
	Resolution   DirectoryEventResolution   `json:"resolution,omitempty"`
	ResolvedAt   *time.Time                 `json:"resolved_at,omitempty"`
	RetryOf      string                     `json:"retry_of,omitempty"`
	Retryable    bool                       `json:"retryable"`
	// The failure's own code, such as ldap.dial_failed, ldap.service_bind_failed,
	// ldap.tls_verify_required or ldap.invalid_credentials, or the SCIM error type an answer to the
	// identity provider gave (scim.uniqueness, say). Absent on a success, and on events from before it
	// was kept.
	Reason string `json:"reason,omitempty"`
	// Safe details of the failure: host, the directory server it used.
	Params map[string]string `json:"params,omitempty"`
	// The login name an LDAP sign-in tried.
	Subject string `json:"subject,omitempty"`
	// The person an LDAP sign-in was for, once the directory knew them.
	UserID string `json:"user_id,omitempty"`
}

// DirectoryEventTrigger is a string enumeration used by DirectoryEvent.Trigger.
type DirectoryEventTrigger string

// Values of DirectoryEventTrigger. The server may add values; handle an unknown one gracefully.
const (
	DirectoryEventTriggerPush   DirectoryEventTrigger = "push"
	DirectoryEventTriggerManual DirectoryEventTrigger = "manual"
	DirectoryEventTriggerSignIn DirectoryEventTrigger = "sign-in"
)

// DirectoryEventResourceType is a string enumeration used by DirectoryEvent.ResourceType.
type DirectoryEventResourceType string

// Values of DirectoryEventResourceType. The server may add values; handle an unknown one gracefully.
const (
	DirectoryEventResourceTypeUser       DirectoryEventResourceType = "user"
	DirectoryEventResourceTypeGroup      DirectoryEventResourceType = "group"
	DirectoryEventResourceTypeConnection DirectoryEventResourceType = "connection"
)

// DirectoryEventStatus is a string enumeration used by DirectoryEvent.Status.
type DirectoryEventStatus string

// Values of DirectoryEventStatus. The server may add values; handle an unknown one gracefully.
const (
	DirectoryEventStatusSuccess DirectoryEventStatus = "success"
	DirectoryEventStatusWarning DirectoryEventStatus = "warning"
	DirectoryEventStatusFailed  DirectoryEventStatus = "failed"
)

// DirectoryErrorsSummary is the type of DirectoryErrors.Summary.
type DirectoryErrorsSummary struct {
	Open          int  `json:"open"`
	Critical      int  `json:"critical"`
	Warning       int  `json:"warning"`
	Retried       int  `json:"retried"`
	ImpactedUsers *int `json:"impacted_users,omitempty"`
}

// DirectoryErrorsCategory is one element of DirectoryErrors.Categories.
type DirectoryErrorsCategory struct {
	Key   string `json:"key,omitempty"`
	Count int    `json:"count"`
}

// DirectoryErrorsConnection is one element of DirectoryErrors.Connections.
type DirectoryErrorsConnection struct {
	Key   string `json:"key,omitempty"`
	Count int    `json:"count"`
}

// DirectoryActionResult is the result of DirectoryService.ResolveError.
type DirectoryActionResult struct {
	Succeeded bool `json:"succeeded"`
}

// ListDirectoryActivityParams holds the optional query parameters of DirectoryService.ListActivity. A
// nil *ListDirectoryActivityParams sends none.
type ListDirectoryActivityParams struct {
	Days         int
	ConnectionID string
	Status       DirectoryEventStatus
	Q            string
	Limit        int
	Offset       int
}

func (p *ListDirectoryActivityParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Days != 0 {
		q.Set("days", strconv.FormatInt(int64(p.Days), 10))
	}
	if p.ConnectionID != "" {
		q.Set("connection_id", p.ConnectionID)
	}
	if p.Status != "" {
		q.Set("status", string(p.Status))
	}
	if p.Q != "" {
		q.Set("q", p.Q)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Offset != 0 {
		q.Set("offset", strconv.FormatInt(int64(p.Offset), 10))
	}
	return q
}

// DirectoryActivity is the result of DirectoryService.ListActivity.
type DirectoryActivity struct {
	Items   []DirectoryEvent          `json:"items,omitempty"`
	Total   int                       `json:"total"`
	Limit   int                       `json:"limit"`
	Offset  int                       `json:"offset"`
	Summary DirectoryActivitySummary  `json:"summary,omitzero"`
	Volume  []DirectoryActivityVolume `json:"volume,omitempty"`
}

// DirectoryActivitySummary is the type of DirectoryActivity.Summary.
type DirectoryActivitySummary struct {
	Total             int  `json:"total"`
	Successful        int  `json:"successful"`
	Warnings          int  `json:"warnings"`
	Failed            int  `json:"failed"`
	AverageDurationMs *int `json:"average_duration_ms,omitempty"`
}

// DirectoryActivityVolume is one element of DirectoryActivity.Volume.
type DirectoryActivityVolume struct {
	Date   string `json:"date,omitempty"`
	Users  int    `json:"users"`
	Groups int    `json:"groups"`
}

// DirectoryAlert notifies administrators when directory operations keep failing.
type DirectoryAlert struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// The connection it watches; null for every connection.
	ConnectionID string `json:"connection_id,omitempty"`
	// Failures within 15 minutes that notify.
	Threshold int `json:"threshold"`
	// False while paused.
	Enabled bool `json:"enabled"`
	// The person it notifies (who made it).
	UserID         string     `json:"user_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitzero"`
	LastNotifiedAt *time.Time `json:"last_notified_at,omitempty"`
}

// CreateDirectoryAlertParams is the request body of DirectoryService.CreateAlert.
type CreateDirectoryAlertParams struct {
	Name         string `json:"name,omitempty"`
	ConnectionID string `json:"connection_id,omitempty"`
	Threshold    int    `json:"threshold"`
}

// UpdateDirectoryAlertParams is the request body of DirectoryService.UpdateAlert.
type UpdateDirectoryAlertParams struct {
	Name *string `json:"name,omitempty"`
	// A connection of the organisation, or "" for every connection.
	ConnectionID *string `json:"connection_id,omitempty"`
	Threshold    *int    `json:"threshold,omitempty"`
	Enabled      *bool   `json:"enabled,omitempty"`
}

// DirectoryProfileValue is one item of the list DirectoryService.Profile returns.
type DirectoryProfileValue struct {
	ConnectionID string `json:"connection_id,omitempty"`
	// The connection's name; empty if it was removed.
	ConnectionName string         `json:"connection_name,omitempty"`
	Fields         map[string]any `json:"fields,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at,omitzero"`
}
