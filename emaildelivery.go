package qeetid

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// EmailDeliveryService reports on the emails sent for a tenant and manages suppressed addresses.
type EmailDeliveryService struct{ c *Client }

// ListMessages lists the emails sent to an address and whether it is suppressed.
//
// API: GET /v1/tenants/{tenantID}/email/messages
func (s *EmailDeliveryService) ListMessages(ctx context.Context, tenantID string, params *ListEmailMessagesParams, opts ...RequestOption) (*EmailMessages, error) {
	return call[EmailMessages](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/email/messages",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// ListSuppressions lists addresses that bounced or reported spam.
//
// API: GET /v1/tenants/{tenantID}/email/suppressions
func (s *EmailDeliveryService) ListSuppressions(ctx context.Context, tenantID string, params *ListEmailSuppressionsParams, opts ...RequestOption) (*Page[EmailSuppression], error) {
	return call[Page[EmailSuppression]](ctx, s.c, &request{
		method:   http.MethodGet,
		path:     "/v1/tenants/{tenantID}/email/suppressions",
		pathArgs: []string{tenantID},
		query:    params.values(),
	}, opts)
}

// DeleteSuppression clears a suppression so email is sent to the address again.
//
// API: DELETE /v1/tenants/{tenantID}/email/suppressions/{id}
func (s *EmailDeliveryService) DeleteSuppression(ctx context.Context, tenantID string, suppressionID string, opts ...RequestOption) error {
	return s.c.do(ctx, &request{
		method:   http.MethodDelete,
		path:     "/v1/tenants/{tenantID}/email/suppressions/{suppressionID}",
		pathArgs: []string{tenantID, suppressionID},
	}, nil, opts)
}

// ListEmailMessagesParams holds the optional query parameters of EmailDeliveryService.ListMessages. A
// nil *ListEmailMessagesParams sends none.
type ListEmailMessagesParams struct {
	To    string
	Limit int
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
	return q
}

// EmailMessages is the result of EmailDeliveryService.ListMessages.
type EmailMessages struct {
	Items       []EmailMessage            `json:"items,omitempty"`
	Suppression *EmailMessagesSuppression `json:"suppression,omitempty"`
}

// EmailMessage is an email sent for a tenant.
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
	// This organisation's mail caused the suppression; the other fields are present only then.
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
}

func (p *ListEmailSuppressionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.FormatInt(int64(p.Limit), 10))
	}
	return q
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
