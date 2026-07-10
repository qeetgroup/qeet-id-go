package transport

import "time"

// Logger is an optional hook invoked once per request, after the retry loop
// settles (success or final failure). It never blocks or alters the
// request/response — implement it for structured logging or metrics, not
// for control flow. A nil Logger (the default) is a no-op.
type Logger interface {
	LogRequest(method, path string, status int, duration time.Duration, requestID string)
}

// logRequest calls the configured Logger, if any. Safe to call with a nil
// Logger — every call site does this rather than nil-checking inline.
func (t *Transport) logRequest(method, path string, status int, start time.Time, requestID string) {
	if t.logger == nil {
		return
	}
	t.logger.LogRequest(method, path, status, time.Since(start), requestID)
}
