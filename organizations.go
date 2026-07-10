package qeetid

import (
	"context"
	"iter"
	"net/url"
	"strconv"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/pagination"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
)

// OrganizationsService manages organizations (tenants).
type OrganizationsService struct{ t *transport.Transport }

func (s *OrganizationsService) Create(ctx context.Context, in CreateOrganizationInput) (*Organization, error) {
	var out Organization
	err := s.t.Post(ctx, "/v1/tenants", in, &out)
	return &out, err
}

func (s *OrganizationsService) Get(ctx context.Context, id string) (*Organization, error) {
	var out Organization
	err := s.t.Get(ctx, "/v1/tenants/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *OrganizationsService) Update(ctx context.Context, id string, in UpdateOrganizationInput) (*Organization, error) {
	var out Organization
	err := s.t.Patch(ctx, "/v1/tenants/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *OrganizationsService) Delete(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/tenants/"+url.PathEscape(id), nil)
}

func (s *OrganizationsService) List(ctx context.Context, limit int, cursor string) (*OrganizationPage, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var env struct {
		marshal.Envelope[Organization]
		NextCursor string `json:"next_cursor"`
	}
	if err := s.t.Get(ctx, "/v1/tenants", transport.RequestOptions{Query: q}, &env); err != nil {
		return nil, err
	}
	return &OrganizationPage{Data: env.Resolve(), NextCursor: env.NextCursor}, nil
}

// All iterates every organization across pages.
func (s *OrganizationsService) All(ctx context.Context, limit int) iter.Seq2[Organization, error] {
	return pagination.Paginate(ctx, "", func(ctx context.Context, cursor string) ([]Organization, string, error) {
		page, err := s.List(ctx, limit, cursor)
		if err != nil {
			return nil, "", err
		}
		return page.Data, page.NextCursor, nil
	})
}
