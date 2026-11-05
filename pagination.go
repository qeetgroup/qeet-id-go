package qeetid

import (
	"context"
	"iter"
)

// Page is one page of a list.
//
// Lists the server pages with cursors return a NextCursor; pass it back as the
// Cursor parameter to fetch the next page, or use the matching All method,
// which does that for you. Lists the server returns whole come back as a
// single Page with an empty NextCursor. Some of those are capped: when a list
// is longer than the cap, the server returns the first items and sets
// Truncated.
type Page[T any] struct {
	// Items holds the page's items.
	Items []T `json:"items"`
	// NextCursor fetches the next page. It is empty on the last page.
	NextCursor string `json:"next_cursor,omitempty"`
	// Truncated reports that the server cut a capped list short: there are
	// more items than Items holds and no cursor to fetch them with. Narrow the
	// request with a filter to see the rest.
	Truncated bool `json:"truncated,omitempty"`
	// Total is how many items match, for lists that report it; nil otherwise.
	Total *int `json:"total,omitempty"`
	// TotalCapped reports that the server stopped counting, so the real total
	// is larger than Total.
	TotalCapped bool `json:"total_capped,omitempty"`
}

// HasMore reports whether another page can be fetched with NextCursor.
func (p *Page[T]) HasMore() bool { return p != nil && p.NextCursor != "" }

// paginate walks a cursor-paginated list, fetching each page when the loop
// asks for its first item. It stops at the last page, when the loop breaks,
// when the context ends, or after yielding the first error. A server that
// repeats a cursor ends the walk rather than looping forever.
func paginate[T any](ctx context.Context, cursor string, fetch func(ctx context.Context, cursor string) (*Page[T], error)) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		seen := map[string]bool{}
		for {
			if err := ctx.Err(); err != nil {
				yield(nil, err)
				return
			}
			page, err := fetch(ctx, cursor)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range page.Items {
				if !yield(&page.Items[i], nil) {
					return
				}
			}
			next := page.NextCursor
			if next == "" || next == cursor || seen[next] {
				return
			}
			seen[cursor] = true
			cursor = next
		}
	}
}
