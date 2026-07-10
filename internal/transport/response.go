package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/qeetgroup/qeet-id-go/internal/constants"
)

// Error is returned by every failed API call. Inspect Status or use the Is*
// helpers; errors.As(err, &qeetid.Error{}) unwraps it (qeetid.Error is a type
// alias for this struct — see the root package's errors.go).
type Error struct {
	Status            int
	Code              string
	Message           string
	RequestID         string
	RetryAfterSeconds int // set on 429 when the server provided Retry-After
}

func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("qeetid: %s (status %d, code %q, request %s)", e.Message, e.Status, e.Code, e.RequestID)
	}
	return fmt.Sprintf("qeetid: %s (status %d, code %q)", e.Message, e.Status, e.Code)
}

func (e *Error) IsUnauthorized() bool { return e.Status == 401 }
func (e *Error) IsForbidden() bool    { return e.Status == 403 }
func (e *Error) IsNotFound() bool     { return e.Status == 404 }
func (e *Error) IsRateLimited() bool  { return e.Status == 429 }

// parseErrorBody maps a non-2xx response body to an *Error. The backend's
// error envelope is {"error":{"code":"...","message":"..."}}; a body that
// doesn't match (e.g. a proxy's HTML error page) still produces a usable
// Error with a generic code/message.
func parseErrorBody(status int, body []byte, requestID string, retryAfterSec int) *Error {
	var env struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.Unmarshal(body, &env)
	code := env.Error.Code
	if code == "" {
		code = "http_" + strconv.Itoa(status)
	}
	msg := env.Error.Message
	if msg == "" {
		msg = "request failed with status " + strconv.Itoa(status)
	}
	return &Error{Status: status, Code: code, Message: msg, RequestID: requestID, RetryAfterSeconds: retryAfterSec}
}

// readResponse reads (capped) the body, and on success unmarshals into out
// (if non-nil and the body is non-empty); on failure it returns a mapped
// *Error. 204 short-circuits to a nil error with out left untouched.
func readResponse(res *http.Response, out any) error {
	requestID := res.Header.Get(constants.HeaderRequestID)
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	data, _ := io.ReadAll(io.LimitReader(res.Body, constants.MaxResponseBytes))
	if res.StatusCode >= 300 {
		return parseErrorBody(res.StatusCode, data, requestID, int(retryAfterDuration(res).Seconds()))
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("qeetid: decode response: %w", err)
		}
	}
	return nil
}
