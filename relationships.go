package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
)

// RelationshipsService manages relationship tuples (relationship-based access control) and checks
// access through them.
type RelationshipsService struct{ c *Client }

// List lists the tuples on an object or naming a subject. Page.Truncated reports a list cut short by
// the server's cap.
//
// API: GET /v1/tenants/{tenantID}/relation-tuples
func (s *RelationshipsService) List(ctx context.Context, tenantID string, params *ListRelationTuplesParams, opts ...RequestOption) (*Page[RelationTuple], error) {
	return call[Page[RelationTuple]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/relation-tuples",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// Create writes a tuple: subject has relation on object.
//
// API: POST /v1/tenants/{tenantID}/relation-tuples
func (s *RelationshipsService) Create(ctx context.Context, tenantID string, params CreateRelationTupleParams, opts ...RequestOption) (*RelationTuple, error) {
	return call[RelationTuple](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/relation-tuples",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// Delete deletes a tuple.
//
// API: DELETE /v1/tenants/{tenantID}/relation-tuples/{id}
func (s *RelationshipsService) Delete(ctx context.Context, tenantID, tupleID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/relation-tuples/{tupleID}",
		pathArgs: []string{tenantID, tupleID},
	}, nil, opts)
}

// Check reports whether a user has a relation on an object, following group and parent relations.
//
// API: POST /v1/tenants/{tenantID}/relation-tuples/check
func (s *RelationshipsService) Check(ctx context.Context, tenantID string, params CheckRelationParams, options *CheckRelationOptions, opts ...RequestOption) (*RelationCheck, error) {
	return call[RelationCheck](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/relation-tuples/check",
		pathArgs: []string{tenantID},
		query:    options.values(),
		body:     params,
	}, opts)
}

// ListTypes lists the object types the tenant's tuples use.
//
// API: GET /v1/tenants/{tenantID}/relation-tuples/types
func (s *RelationshipsService) ListTypes(ctx context.Context, tenantID string, opts ...RequestOption) (*Page[RelationObjectType], error) {
	return call[Page[RelationObjectType]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/relation-tuples/types",
		pathArgs: []string{tenantID},
	}, opts)
}

// ListObjects lists the objects the tenant's tuples name, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/relation-tuples/objects
func (s *RelationshipsService) ListObjects(ctx context.Context, tenantID string, params *ListRelationObjectsParams, opts ...RequestOption) (*Page[RelationObject], error) {
	return call[Page[RelationObject]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/relation-tuples/objects",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AllObjects returns an iterator over every item ListObjects returns, across all pages. It fetches a
// page only when the loop needs it, and breaking out of the loop stops paging. A failed request ends
// the iteration with a non-nil error.
func (s *RelationshipsService) AllObjects(ctx context.Context, tenantID string, params *ListRelationObjectsParams, opts ...RequestOption) iter.Seq2[*RelationObject, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[RelationObject], error) {
		return s.ListObjects(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// Graph expands every subject that can reach an object through a relation.
//
// API: GET /v1/tenants/{tenantID}/relation-tuples/graph
func (s *RelationshipsService) Graph(ctx context.Context, tenantID string, params *RelationGraphParams, opts ...RequestOption) (*RelationGraph, error) {
	return call[RelationGraph](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/relation-tuples/graph",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// ListRelationTuplesParams holds the optional query parameters of RelationshipsService.List. A nil
// *ListRelationTuplesParams sends none.
type ListRelationTuplesParams struct {
	// Object ref "type:id", e.g. document:readme. Mutually exclusive with subject.
	Object string
	// Subject ref "type:id" or "type:id#relation", e.g. user:abc123. Returns all objects the subject
	// appears on. Mutually exclusive with object.
	Subject string
}

func (p *ListRelationTuplesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Object != "" {
		q.Set("object", p.Object)
	}
	if p.Subject != "" {
		q.Set("subject", p.Subject)
	}
	return q
}

// RelationTuple states that a subject has a relation on an object, such as user:alice is viewer of
// document:readme.
type RelationTuple struct {
	ID       string `json:"id,omitempty"`
	Object   string `json:"object,omitempty"`
	Relation string `json:"relation,omitempty"`
	Subject  string `json:"subject,omitempty"`
	// A user subject's person, by display name or address, when they belong to the organisation.
	SubjectName string `json:"subject_name,omitempty"`
}

// CreateRelationTupleParams is the request body of RelationshipsService.Create.
type CreateRelationTupleParams struct {
	// "type:id"
	Object   string `json:"object,omitempty"`
	Relation string `json:"relation,omitempty"`
	// "user:<id>" or "type:id#relation" (userset)
	Subject string `json:"subject,omitempty"`
}

// CheckRelationParams is the request body of RelationshipsService.Check.
type CheckRelationParams struct {
	Object   string `json:"object,omitempty"`
	Relation string `json:"relation,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// CheckRelationOptions holds the optional query parameters of RelationshipsService.Check. A nil
// *CheckRelationOptions sends none.
type CheckRelationOptions struct {
	// When "true", returns the root-to-leaf chain of relation tuples that granted (or would grant) access
	// instead of just the boolean.
	Explain bool
}

func (p *CheckRelationOptions) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Explain {
		q.Set("explain", "true")
	}
	return q
}

// RelationCheck is the answer to a relationship-based access check.
type RelationCheck struct {
	Allowed   bool               `json:"allowed"`
	DecidedBy AuthzDecidedBy     `json:"decided_by,omitzero"`
	Trace     []string           `json:"trace,omitempty"`
	Path      []RelationPathStep `json:"path,omitempty"`
}

// RelationPathStep is one element of RelationCheck.Path.
type RelationPathStep struct {
	Object   string `json:"object,omitempty"`
	Relation string `json:"relation,omitempty"`
	Subject  string `json:"subject,omitempty"`
	Depth    int    `json:"depth"`
}

// RelationObjectType is one item of the list RelationshipsService.ListTypes returns.
type RelationObjectType struct {
	Type    string `json:"type,omitempty"`
	Objects int    `json:"objects"`
	Tuples  int    `json:"tuples"`
}

// ListRelationObjectsParams holds the optional query parameters of RelationshipsService.ListObjects. A
// nil *ListRelationObjectsParams sends none.
type ListRelationObjectsParams struct {
	// The type of object, such as document.
	Type string
	// Page size, 1 to 500; 100 by default.
	Limit int
	// The next_cursor of the previous page.
	Cursor string
}

func (p *ListRelationObjectsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Type != "" {
		q.Set("type", p.Type)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ListRelationObjectsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListRelationObjectsParams) withCursor(cursor string) *ListRelationObjectsParams {
	var c ListRelationObjectsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// RelationObject is one item of the list RelationshipsService.ListObjects returns.
type RelationObject struct {
	// "type:id"
	Object string `json:"object,omitempty"`
	ID     string `json:"id,omitempty"`
	Tuples int    `json:"tuples"`
}

// RelationGraphParams holds the optional query parameters of RelationshipsService.Graph. A nil
// *RelationGraphParams sends none.
type RelationGraphParams struct {
	// Root object ref "type:id", e.g. document:security-runbook
	Object string
	// Relation to expand from the root, e.g. viewer
	Relation string
	// Maximum hops to follow (default and max 10)
	Depth int
}

func (p *RelationGraphParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Object != "" {
		q.Set("object", p.Object)
	}
	if p.Relation != "" {
		q.Set("relation", p.Relation)
	}
	if p.Depth != 0 {
		q.Set("depth", strconv.FormatInt(int64(p.Depth), 10))
	}
	return q
}

// RelationGraph is the result of RelationshipsService.Graph.
type RelationGraph struct {
	Nodes []RelationGraphNode `json:"nodes,omitempty"`
	Edges []RelationGraphEdge `json:"edges,omitempty"`
}

// RelationGraphNode is one element of RelationGraph.Nodes.
type RelationGraphNode struct {
	// "type:id"
	ID    string `json:"id,omitempty"`
	Type  string `json:"type,omitempty"`
	Label string `json:"label,omitempty"`
	// A user node's person, by display name or address, when they belong to the organisation.
	Name string `json:"name,omitempty"`
}

// RelationGraphEdge is one element of RelationGraph.Edges.
type RelationGraphEdge struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Relation string `json:"relation,omitempty"`
}
