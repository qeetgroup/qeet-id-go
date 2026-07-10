// Package validation holds light, fail-fast input checks — required-field
// and UUID-shape checks that catch obviously-bad input before a network
// round trip, without duplicating the backend's own validation.
//
// Kept intentionally minimal: the backend is the source of truth for
// validity (it returns a proper 422 either way), so this package only
// catches the cases cheap enough to check without a request — an empty ID
// that would otherwise silently build a malformed URL like "/v1/users/".
package validation

import (
	"fmt"
	"regexp"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Required fails if value is empty.
func Required(field, value string) error {
	if value == "" {
		return fmt.Errorf("qeetid: %s is required", field)
	}
	return nil
}

// UUID fails if value isn't a well-formed UUID string.
func UUID(field, value string) error {
	if !uuidPattern.MatchString(value) {
		return fmt.Errorf("qeetid: %s must be a valid UUID, got %q", field, value)
	}
	return nil
}
