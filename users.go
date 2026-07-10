package qeetid

import (
	"context"
	"iter"
	"net/url"
	"strconv"

	"github.com/qeetgroup/qeet-id-go/internal/marshal"
	"github.com/qeetgroup/qeet-id-go/internal/pagination"
	"github.com/qeetgroup/qeet-id-go/internal/transport"
	"github.com/qeetgroup/qeet-id-go/internal/validation"
)

// UsersService manages end-user accounts.
type UsersService struct{ t *transport.Transport }

func (s *UsersService) Create(ctx context.Context, in CreateUserInput) (*User, error) {
	var out User
	err := s.t.Post(ctx, "/v1/users", in, &out)
	return &out, err
}

func (s *UsersService) Get(ctx context.Context, id string) (*User, error) {
	if err := validation.Required("id", id); err != nil {
		return nil, err
	}
	var out User
	err := s.t.Get(ctx, "/v1/users/"+url.PathEscape(id), transport.RequestOptions{}, &out)
	return &out, err
}

func (s *UsersService) Update(ctx context.Context, id string, in UpdateUserInput) (*User, error) {
	var out User
	err := s.t.Patch(ctx, "/v1/users/"+url.PathEscape(id), in, &out)
	return &out, err
}

func (s *UsersService) Delete(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/users/"+url.PathEscape(id), nil)
}

func (s *UsersService) SetPassword(ctx context.Context, id, password string) error {
	body := map[string]string{"password": password}
	return s.t.Post(ctx, "/v1/users/"+url.PathEscape(id)+"/password", body, nil)
}

func (s *UsersService) List(ctx context.Context, params ListParams) (*UserPage, error) {
	q := url.Values{}
	if params.Tenant != "" {
		q.Set("tenant", params.Tenant)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	var env struct {
		marshal.Envelope[User]
		NextCursor string `json:"next_cursor"`
	}
	if err := s.t.Get(ctx, "/v1/users", transport.RequestOptions{Query: q}, &env); err != nil {
		return nil, err
	}
	return &UserPage{Data: env.Resolve(), NextCursor: env.NextCursor}, nil
}

// All iterates every user across pages. Range over it with Go 1.23+
// range-over-func; check the error on each step.
func (s *UsersService) All(ctx context.Context, params ListParams) iter.Seq2[User, error] {
	return pagination.Paginate(ctx, params.Cursor, func(ctx context.Context, cursor string) ([]User, string, error) {
		p := params
		p.Cursor = cursor
		page, err := s.List(ctx, p)
		if err != nil {
			return nil, "", err
		}
		return page.Data, page.NextCursor, nil
	})
}

// BulkCreate creates up to 1000 users in one call (POST /v1/users/bulk).
func (s *UsersService) BulkCreate(ctx context.Context, in BulkCreateInput) (*BulkImportResult, error) {
	var out BulkImportResult
	err := s.t.Post(ctx, "/v1/users/bulk", in, &out)
	return &out, err
}

// BulkImport uploads a vendor export file directly (NDJSON for Auth0, CSV
// for Cognito, or a Microsoft Graph {"value":[...]} JSON document for Azure
// B2C) via POST /v1/users/bulk/import?source=. rawBody is passed through
// unparsed — the backend does the format-specific parsing.
func (s *UsersService) BulkImport(ctx context.Context, source BulkImportSource, contentType string, rawBody []byte) (*BulkImportResult, error) {
	q := url.Values{"source": {string(source)}}
	var out BulkImportResult
	err := s.t.Do(ctx, "POST", "/v1/users/bulk/import", transport.RequestOptions{
		Query:          q,
		RawBody:        rawBody,
		RawContentType: contentType,
	}, &out)
	return &out, err
}

// ListDeleted lists soft-deleted users (the recycle bin).
func (s *UsersService) ListDeleted(ctx context.Context, limit int) ([]User, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var env marshal.Envelope[User]
	err := s.t.Get(ctx, "/v1/users/deleted", transport.RequestOptions{Query: q}, &env)
	return env.Resolve(), err
}

// Restore undoes a soft delete.
func (s *UsersService) Restore(ctx context.Context, id string) (*User, error) {
	var out User
	err := s.t.Post(ctx, "/v1/users/"+url.PathEscape(id)+"/restore", nil, &out)
	return &out, err
}

// Purge permanently (hard) deletes a soft-deleted user. Irreversible.
func (s *UsersService) Purge(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/users/"+url.PathEscape(id)+"/purge", nil)
}

// ResetMFA is an admin-initiated MFA reset for a user who's lost access to
// their factors.
func (s *UsersService) ResetMFA(ctx context.Context, id string) error {
	return s.t.Delete(ctx, "/v1/users/"+url.PathEscape(id)+"/mfa", nil)
}

// VerifyEmailStart sends a verification code to the user's email (or the
// override in in.Email).
func (s *UsersService) VerifyEmailStart(ctx context.Context, id string, in VerifyEmailStartInput) error {
	return s.t.Post(ctx, "/v1/users/"+url.PathEscape(id)+"/verify/email/start", in, nil)
}

// VerifyEmailConfirm redeems the code sent by VerifyEmailStart.
func (s *UsersService) VerifyEmailConfirm(ctx context.Context, id string, in VerifyEmailConfirmInput) error {
	return s.t.Post(ctx, "/v1/users/"+url.PathEscape(id)+"/verify/email/confirm", in, nil)
}

// VerifyPhoneStart sends a verification code to the user's phone (or the
// override in in.Phone).
func (s *UsersService) VerifyPhoneStart(ctx context.Context, id string, in VerifyPhoneStartInput) error {
	return s.t.Post(ctx, "/v1/users/"+url.PathEscape(id)+"/verify/phone/start", in, nil)
}

// VerifyPhoneConfirm redeems the code sent by VerifyPhoneStart.
func (s *UsersService) VerifyPhoneConfirm(ctx context.Context, id string, in VerifyPhoneConfirmInput) error {
	return s.t.Post(ctx, "/v1/users/"+url.PathEscape(id)+"/verify/phone/confirm", in, nil)
}
