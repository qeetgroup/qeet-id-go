// Package pagination holds the generic cursor-walking iterator shared by
// every resource's All() method. Kept pure and dependency-free so it's
// unit-tested in isolation from any real HTTP transport.
package pagination

import (
	"context"
	"iter"
)

// Paginate walks cursor-paginated pages lazily, yielding one item at a time.
// fetchPage returns a page's items and the cursor for the next page ("" when
// the collection is exhausted). On error it yields a single (zero, err) and
// stops. The next==cursor guard defends against a server that echoes the
// same cursor forever.
func Paginate[T any](
	ctx context.Context,
	startCursor string,
	fetchPage func(ctx context.Context, cursor string) ([]T, string, error),
) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		cursor := startCursor
		for {
			items, next, err := fetchPage(ctx, cursor)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, it := range items {
				if !yield(it, nil) {
					return
				}
			}
			if next == "" || next == cursor {
				return
			}
			cursor = next
		}
	}
}
