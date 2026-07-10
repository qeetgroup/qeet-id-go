package qeetid

import (
	"context"
	"net/url"

	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// AuthZENSubject, AuthZENResource, AuthZENAction are the standard AuthZEN
// (OpenID unified authorization) request shapes.
type AuthZENSubject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type AuthZENResource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type AuthZENAction struct {
	Name string `json:"name"`
}

// EvaluateInput is an AuthZEN /access/v1/evaluation request. The backend
// routes to RBAC when Resource.Type == "permission", else to ReBAC — so one
// endpoint fronts both authorization models with a single standard shape.
// Set Context["explain"] = true to get a grant-path trace back in the
// response Context.
type EvaluateInput struct {
	Subject  AuthZENSubject  `json:"subject"`
	Resource AuthZENResource `json:"resource"`
	Action   AuthZENAction   `json:"action"`
	Context  map[string]any  `json:"context,omitempty"`
}

// EvaluateResult is the AuthZEN decision response.
type EvaluateResult struct {
	Decision bool           `json:"decision"`
	Context  map[string]any `json:"context,omitempty"`
}

// AuthZENService implements the AuthZEN (OpenID unified authorization)
// /access/v1/evaluation endpoint — a single standard request/response shape
// fronting both RBAC and ReBAC. Exposed as Authorization.Decisions.
type AuthZENService struct{ t *transport.Transport }

func (s *AuthZENService) Evaluate(ctx context.Context, tenantID string, in EvaluateInput) (*EvaluateResult, error) {
	var out EvaluateResult
	err := s.t.Post(ctx, "/v1/tenants/"+url.PathEscape(tenantID)+"/access/v1/evaluation", in, &out)
	return &out, err
}
