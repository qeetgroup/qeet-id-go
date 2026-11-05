package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AuthZENService answers access questions in the OpenID AuthZEN shape, deciding with role- or
// relationship-based rules.
type AuthZENService struct{ c *Client }

// Evaluate answers one access question.
//
// API: POST /v1/tenants/{tenantID}/access/v1/evaluation
func (s *AuthZENService) Evaluate(ctx context.Context, tenantID string, params EvaluationParams, opts ...RequestOption) (*Evaluation, error) {
	return call[Evaluation](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/access/v1/evaluation",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// EvaluateBatch answers many access questions in one request.
//
// API: POST /v1/tenants/{tenantID}/access/v1/evaluations
func (s *AuthZENService) EvaluateBatch(ctx context.Context, tenantID string, params BatchEvaluationParams, opts ...RequestOption) (*BatchEvaluation, error) {
	return call[BatchEvaluation](ctx, s.c, &request{
		method:   http.MethodPost,
		path:     "/v1/tenants/{tenantID}/access/v1/evaluations",
		pathArgs: []string{tenantID},
		body:     params,
	}, opts)
}

// ListDecisions lists recorded access decisions, one page at a time.
//
// API: GET /v1/tenants/{tenantID}/authorization/decisions
func (s *AuthZENService) ListDecisions(ctx context.Context, tenantID string, params *ListDecisionsParams, opts ...RequestOption) (*Page[AccessDecision], error) {
	return call[Page[AccessDecision]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/authorization/decisions",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// AllDecisions returns an iterator over every item ListDecisions returns, across all pages. It fetches
// a page only when the loop needs it, and breaking out of the loop stops paging. A failed request ends
// the iteration with a non-nil error.
func (s *AuthZENService) AllDecisions(ctx context.Context, tenantID string, params *ListDecisionsParams, opts ...RequestOption) iter.Seq2[*AccessDecision, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[AccessDecision], error) {
		return s.ListDecisions(ctx, tenantID, params.withCursor(cursor), opts...)
	})
}

// EvaluationParams is the request body of AuthZENService.Evaluate.
type EvaluationParams struct {
	Subject  EvaluationParamsSubject   `json:"subject"`
	Resource *EvaluationParamsResource `json:"resource,omitempty"`
	Action   EvaluationParamsAction    `json:"action"`
	// Optional. The request's attributes for ABAC policies (context.*), such as ip or mfa. { explain:
	// true } includes the grant path and the decision's trace inline in the response context, and is not
	// an attribute.
	Context map[string]any `json:"context,omitempty"`
}

// EvaluationParamsSubject is the type of EvaluationParams.Subject.
type EvaluationParamsSubject struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	// The subject's attributes for ABAC policies (subject.*).
	Properties map[string]any `json:"properties,omitempty"`
}

// EvaluationParamsResource is the type of EvaluationParams.Resource.
type EvaluationParamsResource struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	// The resource's attributes for ABAC policies (resource.*).
	Properties map[string]any `json:"properties,omitempty"`
}

// EvaluationParamsAction is the type of EvaluationParams.Action.
type EvaluationParamsAction struct {
	Name string `json:"name,omitempty"`
	// The action's attributes for ABAC policies (action.*).
	Properties map[string]any `json:"properties,omitempty"`
}

// Evaluation is the answer to an AuthZEN access question.
type Evaluation struct {
	Decision bool           `json:"decision"`
	Context  map[string]any `json:"context,omitempty"`
}

// BatchEvaluationParams is the request body of AuthZENService.EvaluateBatch.
type BatchEvaluationParams struct {
	Subject     *AuthZENSubject                   `json:"subject,omitempty"`
	Resource    *AuthZENResource                  `json:"resource,omitempty"`
	Action      *AuthZENAction                    `json:"action,omitempty"`
	Context     map[string]any                    `json:"context,omitempty"`
	Evaluations []BatchEvaluationParamsEvaluation `json:"evaluations,omitempty"`
	Options     *BatchEvaluationParamsOptions     `json:"options,omitempty"`
}

// AuthZENSubject is the type of BatchEvaluationParams.Subject.
type AuthZENSubject struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	// The subject's attributes for ABAC policies (subject.*).
	Properties map[string]any `json:"properties,omitempty"`
}

// AuthZENResource is the type of BatchEvaluationParams.Resource.
type AuthZENResource struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	// The resource's attributes for ABAC policies (resource.*).
	Properties map[string]any `json:"properties,omitempty"`
}

// AuthZENAction is the type of BatchEvaluationParams.Action.
type AuthZENAction struct {
	Name string `json:"name,omitempty"`
	// The action's attributes for ABAC policies (action.*).
	Properties map[string]any `json:"properties,omitempty"`
}

// BatchEvaluationParamsEvaluation is one element of BatchEvaluationParams.Evaluations.
type BatchEvaluationParamsEvaluation struct {
	Subject  *AuthZENSubject  `json:"subject,omitempty"`
	Resource *AuthZENResource `json:"resource,omitempty"`
	Action   *AuthZENAction   `json:"action,omitempty"`
	Context  map[string]any   `json:"context,omitempty"`
}

// BatchEvaluationParamsOptions is the type of BatchEvaluationParams.Options.
type BatchEvaluationParamsOptions struct {
	EvaluationsSemantic BatchEvaluationParamsEvaluationsSemantic `json:"evaluations_semantic,omitempty"`
}

// BatchEvaluationParamsEvaluationsSemantic is a string enumeration used by
// BatchEvaluationParamsOptions.EvaluationsSemantic.
type BatchEvaluationParamsEvaluationsSemantic string

// Values of BatchEvaluationParamsEvaluationsSemantic. The server may add values; handle an unknown one gracefully.
const (
	BatchEvaluationParamsEvaluationsSemanticExecuteAll          BatchEvaluationParamsEvaluationsSemantic = "execute_all"
	BatchEvaluationParamsEvaluationsSemanticDenyOnFirstDeny     BatchEvaluationParamsEvaluationsSemantic = "deny_on_first_deny"
	BatchEvaluationParamsEvaluationsSemanticPermitOnFirstPermit BatchEvaluationParamsEvaluationsSemantic = "permit_on_first_permit"
)

// BatchEvaluation is the result of AuthZENService.EvaluateBatch.
type BatchEvaluation struct {
	Evaluations []BatchEvaluationEvaluation `json:"evaluations,omitempty"`
	Decision    bool                        `json:"decision"`
	Context     map[string]any              `json:"context,omitempty"`
}

// BatchEvaluationEvaluation is one element of BatchEvaluation.Evaluations.
type BatchEvaluationEvaluation struct {
	Decision bool           `json:"decision"`
	Context  map[string]any `json:"context,omitempty"`
}

// DecisionSource is a string enumeration used by ListDecisionsParams.Source.
type DecisionSource string

// Values of DecisionSource. The server may add values; handle an unknown one gracefully.
const (
	DecisionSourceCheck         DecisionSource = "check"
	DecisionSourceRelationCheck DecisionSource = "relation_check"
	DecisionSourceAuthZEN       DecisionSource = "authzen"
	DecisionSourceABACEvaluate  DecisionSource = "abac_evaluate"
	DecisionSourceABACDraft     DecisionSource = "abac_draft"
)

// ListDecisionsParams holds the optional query parameters of AuthZENService.ListDecisions. A nil
// *ListDecisionsParams sends none.
type ListDecisionsParams struct {
	// Only decisions about this subject (a user id, or type:id).
	Subject string
	Source  DecisionSource
	// Only allowed (true) or refused (false) decisions.
	Allowed *bool
	// Decisions of the last 1 to 30 days; 30 by default.
	Days int
	// Page size, 1 to 200; 50 by default.
	Limit int
	// The next_cursor of the previous page.
	Cursor string
}

func (p *ListDecisionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Subject != "" {
		q.Set("subject", p.Subject)
	}
	if p.Source != "" {
		q.Set("source", string(p.Source))
	}
	if p.Allowed != nil {
		q.Set("allowed", strconv.FormatBool(*p.Allowed))
	}
	if p.Days != 0 {
		q.Set("days", strconv.FormatInt(int64(p.Days), 10))
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ListDecisionsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListDecisionsParams) withCursor(cursor string) *ListDecisionsParams {
	var c ListDecisionsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// AccessDecision is one item of the list AuthZENService.ListDecisions returns.
type AccessDecision struct {
	ID        string         `json:"id,omitempty"`
	DecidedAt time.Time      `json:"decided_at,omitzero"`
	Source    DecisionSource `json:"source,omitempty"`
	// user, api_key, service or agent.
	CallerType string `json:"caller_type,omitempty"`
	CallerID   string `json:"caller_id,omitempty"`
	// The application whose token asked, when one did.
	CallerClient string `json:"caller_client,omitempty"`
	Subject      string `json:"subject,omitempty"`
	// "permission" for a permission check, else type:id.
	Resource string `json:"resource,omitempty"`
	// The permission key, relation or action.
	Action  string `json:"action,omitempty"`
	Allowed bool   `json:"allowed"`
	// rbac, rebac, abac or none: the model that decided.
	Model string `json:"model,omitempty"`
	// The role, relation or ABAC policy that decided.
	Rule       string `json:"rule,omitempty"`
	ReasonCode string `json:"reason_code,omitempty"`
}
