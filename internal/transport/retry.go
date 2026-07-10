package transport

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/qeetgroup/qeet-id-go/internal/constants"
)

// shouldRetry reports whether a response with this status is retry-eligible.
// 429 always is; a 5xx only is if the request was idempotent (the server may
// have already applied a mutation before failing on a non-idempotent call).
func shouldRetry(status int, idempotent bool) bool {
	return status == http.StatusTooManyRequests || (status >= 500 && idempotent)
}

// retryAfterDuration parses the Retry-After header (seconds), or zero if
// absent/unparseable.
func retryAfterDuration(res *http.Response) time.Duration {
	v := res.Header.Get(constants.HeaderRetryAfter)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return time.Duration(n) * time.Second
}

// backoff computes exponential-with-jitter delay for a retry attempt:
// ~250ms, 500ms, 1s, ...
func backoff(attempt int) time.Duration {
	base := 250 * time.Millisecond * time.Duration(1<<attempt)
	return base + time.Duration(rand.Intn(100))*time.Millisecond
}

// sleep waits d, honoring context cancellation.
func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
