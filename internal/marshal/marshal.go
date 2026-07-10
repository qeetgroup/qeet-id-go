// Package marshal holds small JSON helpers shared by every resource file:
// tolerant decoding of non-JSON error bodies, and the {items:[]}/{data:[]}
// envelope shape several backend list endpoints use interchangeably.
package marshal

import "encoding/json"

// SafeJSON decodes s as JSON; if that fails, it returns s itself so a
// non-JSON error body (e.g. an HTML error page from a proxy) is never lost.
func SafeJSON(s string) any {
	if s == "" {
		return nil
	}
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	return v
}

// Envelope is the common list-response shape: some endpoints key the slice
// as "items", others as "data". Decode into this and call Resolve.
type Envelope[T any] struct {
	Items []T `json:"items"`
	Data  []T `json:"data"`
}

// Resolve returns Items if present, else Data (nil if the response used
// neither key — the caller should treat that as an empty list, not fail).
func (e Envelope[T]) Resolve() []T {
	if e.Items != nil {
		return e.Items
	}
	return e.Data
}
