package qeetid

import (
	"context"
	"net/url"
	"strconv"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// Tuple is a ReBAC relationship assertion: "object relation subject". object
// is "type:id"; subject is "user:id" for a direct grant, or "type:id#relation"
// for a userset (e.g. group:eng#member).
type Tuple struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Relation string `json:"relation"`
	Subject  string `json:"subject"`
}

type CreateTupleInput struct {
	Object   string `json:"object"`
	Relation string `json:"relation"`
	Subject  string `json:"subject"`
}

// CheckRelationInput is a ReBAC check — does user_id have relation on object,
// resolving usersets recursively?
type CheckRelationInput struct {
	Object   string `json:"object"`
	Relation string `json:"relation"`
	UserID   string `json:"user_id"`
}

// RelationCheckResult is the response shape; Path is populated only when
// Check is called with explain=true.
type RelationCheckResult struct {
	Allowed bool               `json:"allowed"`
	Path    []RelationPathStep `json:"path,omitempty"`
}

// RelationPathStep is one hop in a Check(explain=true) grant path.
type RelationPathStep struct {
	Object   string `json:"object"`
	Relation string `json:"relation"`
	Subject  string `json:"subject"`
	Depth    int    `json:"depth"`
}

// GraphNode and GraphEdge describe the identity-graph shape returned by
// Graph — a BFS expansion of every subject reachable from an object+relation.
type GraphNode struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type GraphEdge struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Relation string `json:"relation"`
}

type RelationGraph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// RelationTuplesService manages ReBAC (Zanzibar-style) relationship tuples,
// their recursive Check, and the Graph visualization. Exposed as
// Authorization.Relationships.
type RelationTuplesService struct{ t *transport.Transport }

func (s *RelationTuplesService) Create(ctx context.Context, tenantID string, in CreateTupleInput) (*Tuple, error) {
	var out Tuple
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/relation-tuples", in, &out)
	return &out, err
}

func (s *RelationTuplesService) Delete(ctx context.Context, tenantID, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/relation-tuples/"+url.PathEscape(id), nil)
}

// ListByObject lists every tuple on an object (e.g. "document:readme").
func (s *RelationTuplesService) ListByObject(ctx context.Context, tenantID, object string) ([]Tuple, error) {
	q := url.Values{"object": {object}}
	var env marshal.Envelope[Tuple]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/relation-tuples", transport.RequestOptions{Query: q}, &env)
	return env.Resolve(), err
}

// ListBySubject is the reverse lookup: every tuple naming this subject
// (e.g. "user:123" or "group:eng#member").
func (s *RelationTuplesService) ListBySubject(ctx context.Context, tenantID, subject string) ([]Tuple, error) {
	q := url.Values{"subject": {subject}}
	var env marshal.Envelope[Tuple]
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/relation-tuples", transport.RequestOptions{Query: q}, &env)
	return env.Resolve(), err
}

// Check resolves whether in.UserID has in.Relation on in.Object, recursively
// through usersets. Pass explain=true to get the grant path back in Path.
func (s *RelationTuplesService) Check(ctx context.Context, tenantID string, in CheckRelationInput, explain bool) (*RelationCheckResult, error) {
	q := url.Values{}
	if explain {
		q.Set("explain", "true")
	}
	var out RelationCheckResult
	err := s.t.Do(ctx, "POST", "/v1/tenants/"+url.PathEscape(tenantID)+"/relation-tuples/check", transport.RequestOptions{Query: q, Body: in}, &out)
	return &out, err
}

// Graph expands every subject reachable from object+relation (BFS, capped at
// depth 1-10, default 10) — the data behind the console's Identity Graph
// visualization.
func (s *RelationTuplesService) Graph(ctx context.Context, tenantID, object, relation string, depth int) (*RelationGraph, error) {
	q := url.Values{"object": {object}, "relation": {relation}}
	if depth > 0 {
		q.Set("depth", strconv.Itoa(depth))
	}
	var out RelationGraph
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/relation-tuples/graph", transport.RequestOptions{Query: q}, &out)
	return &out, err
}
