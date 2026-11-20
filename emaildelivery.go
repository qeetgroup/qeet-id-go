package qeetid

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// EmailDeliveryService reports on the emails sent for an organization and manages suppressed addresses.
type EmailDeliveryService struct{ c *Client }

// ListMessages lists the emails sent to an address, newest first, one page at a time, and whether it is
// suppressed.
//
// API: GET /v1/organizations/{organizationID}/email/messages
func (s *EmailDeliveryService) ListMessages(ctx context.Context, organizationID string, params *ListEmailMessagesParams, opts ...RequestOption) (*EmailMessages, error) {
	return call[EmailMessages](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/email/messages",
		pathArgs: []string{organizationID},
		query:    params.values(),
	}, opts)
}

// ListSuppressions lists addresses that bounced or reported spam, newest first, one page at a time.
//
// API: GET /v1/organizations/{organizationID}/email/suppressions
func (s *EmailDeliveryService) ListSuppressions(ctx context.Context, organizationID string, params *ListEmailSuppressionsParams, opts ...RequestOption) (*Page[EmailSuppression], error) {
	return call[Page[EmailSuppression]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/organizations/{organizationID}/email/suppressions",
		pathArgs: []string{organizationID},
		query:    params.values(),
	}, opts)
}

// AllSuppressions returns an iterator over every suppression ListSuppressions returns, across all pages.
// It fetches a page only when the loop needs it, and breaking out of the loop stops paging. A failed
// request ends the iteration with a non-nil error.
func (s *EmailDeliveryService) AllSuppressions(ctx context.Context, organizationID string, params *ListEmailSuppressionsParams, opts ...RequestOption) iter.Seq2[*EmailSuppression, error] {
	return paginate(ctx, params.cursor(), func(ctx context.Context, cursor string) (*Page[EmailSuppression], error) {
		return s.ListSuppressions(ctx, organizationID, params.withCursor(cursor), opts...)
	})
}

// DeleteSuppression clears a suppression so email is sent to the address again.
//
// API: DELETE /v1/organizations/{organizationID}/email/suppressions/{id}
func (s *EmailDeliveryService) DeleteSuppression(ctx context.Context, organizationID, suppressionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/organizations/{organizationID}/email/suppressions/{suppressionID}",
		pathArgs: []string{organizationID, suppressionID},
	}, nil, opts)
}

// ListEmailMessagesParams holds the optional query parameters of EmailDeliveryService.ListMessages. A
// nil *ListEmailMessagesParams sends none.
type ListEmailMessagesParams struct {
	To    string
	Limit int
	// The previous page's next_cursor.
	Cursor string
}

func (p *ListEmailMessagesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

// EmailMessages is the result of EmailDeliveryService.ListMessages.
type EmailMessages struct {
	Items []EmailMessage `json:"items,omitempty"`
	// Present when there is a next page: send it back as Cursor.
	NextCursor  string                    `json:"next_cursor,omitempty"`
	Suppression *EmailMessagesSuppression `json:"suppression,omitempty"`
}

// EmailMessage is an email sent for an organization.
type EmailMessage struct {
	ID string `json:"id,omitempty"`
	// The template the message was sent from (invite
	Template string             `json:"template,omitempty"`
	Status   EmailMessageStatus `json:"status,omitempty"`
	// The provider's reason for a bounce, complaint or failure (for example Permanent/General: 550 no
	// such user).
	Reason    string    `json:"reason,omitempty"`
	SentAt    time.Time `json:"sent_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

// EmailMessageStatus is a string enumeration used by EmailMessage.Status.
type EmailMessageStatus string

// Values of EmailMessageStatus. The server may add values; handle an unknown one gracefully.
const (
	EmailMessageStatusSent       EmailMessageStatus = "sent"
	EmailMessageStatusDelivered  EmailMessageStatus = "delivered"
	EmailMessageStatusBounced    EmailMessageStatus = "bounced"
	EmailMessageStatusComplained EmailMessageStatus = "complained"
	EmailMessageStatusSuppressed EmailMessageStatus = "suppressed"
	EmailMessageStatusFailed     EmailMessageStatus = "failed"
)

// EmailMessagesSuppression is the type of EmailMessages.Suppression.
type EmailMessagesSuppression struct {
	// This organization's mail caused the suppression; the other fields are present only then.
	Own       bool                   `json:"own"`
	ID        string                 `json:"id,omitempty"`
	Reason    EmailSuppressionReason `json:"reason,omitempty"`
	Detail    string                 `json:"detail,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitzero"`
}

// EmailSuppressionReason is a string enumeration used by EmailMessagesSuppression.Reason.
type EmailSuppressionReason string

// Values of EmailSuppressionReason. The server may add values; handle an unknown one gracefully.
const (
	EmailSuppressionReasonBounce    EmailSuppressionReason = "bounce"
	EmailSuppressionReasonComplaint EmailSuppressionReason = "complaint"
)

// ListEmailSuppressionsParams holds the optional query parameters of
// EmailDeliveryService.ListSuppressions. A nil *ListEmailSuppressionsParams sends none.
type ListEmailSuppressionsParams struct {
	Limit int
	// The previous page's next_cursor.
	Cursor string
}

func (p *ListEmailSuppressionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ListEmailSuppressionsParams) cursor() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

func (p *ListEmailSuppressionsParams) withCursor(cursor string) *ListEmailSuppressionsParams {
	var c ListEmailSuppressionsParams
	if p != nil {
		c = *p
	}
	c.Cursor = cursor
	return &c
}

// EmailSuppression is an address email is not sent to, because it bounced or reported spam.
type EmailSuppression struct {
	ID string `json:"id,omitempty"`
	// The address masked (first character and domain).
	Address   string                 `json:"address,omitempty"`
	Reason    EmailSuppressionReason `json:"reason,omitempty"`
	Detail    string                 `json:"detail,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitzero"`
}
