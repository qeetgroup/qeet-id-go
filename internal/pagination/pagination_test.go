package pagination

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
)

func TestPaginate_WalksAllPages(t *testing.T) {
	pages := [][]string{{"a", "b"}, {"c"}}
	cursors := []string{"c2", ""}
	call := 0
	fetch := func(ctx context.Context, cursor string) ([]string, string, error) {
		i := call
		call++
		return pages[i], cursors[i], nil
	}

	var got []string
	for v, err := range Paginate(context.Background(), "", fetch) {
		if err != nil {
			t.Fatalf("iteration error: %v", err)
		}
		got = append(got, v)
	}
	if strings.Join(got, ",") != "a,b,c" {
		t.Fatalf("got = %v", got)
	}
}

func TestPaginate_EarlyBreakStopsFetching(t *testing.T) {
	var fetches int32
	fetch := func(ctx context.Context, cursor string) ([]int, string, error) {
		atomic.AddInt32(&fetches, 1)
		return []int{1, 2}, "more", nil // always another page
	}

	count := 0
	for _, err := range Paginate(context.Background(), "", fetch) {
		if err != nil {
			t.Fatalf("iteration error: %v", err)
		}
		count++
		break
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
	if got := atomic.LoadInt32(&fetches); got != 1 {
		t.Fatalf("fetches = %d, want 1 (break must stop paging)", got)
	}
}

func TestPaginate_ErrorPropagatesAndStops(t *testing.T) {
	boom := errors.New("boom")
	fetch := func(ctx context.Context, cursor string) ([]int, string, error) {
		return nil, "", boom
	}

	var gotErr error
	n := 0
	for _, err := range Paginate(context.Background(), "", fetch) {
		n++
		if err != nil {
			gotErr = err
		}
	}
	if n != 1 {
		t.Fatalf("iterations = %d, want exactly 1 (one error yield, then stop)", n)
	}
	if !errors.Is(gotErr, boom) {
		t.Fatalf("gotErr = %v, want %v", gotErr, boom)
	}
}

func TestPaginate_SameCursorGuardStopsInfiniteLoop(t *testing.T) {
	var fetches int32
	fetch := func(ctx context.Context, cursor string) ([]int, string, error) {
		atomic.AddInt32(&fetches, 1)
		return []int{1}, "stuck", nil // echoes a cursor that never changes
	}

	// Seed the iterator so its second fetch would receive "stuck" == "stuck".
	count := 0
	for _, err := range Paginate(context.Background(), "stuck", fetch) {
		if err != nil {
			t.Fatalf("iteration error: %v", err)
		}
		count++
		if count > 5 {
			t.Fatal("did not stop on repeated cursor")
		}
	}
	if got := atomic.LoadInt32(&fetches); got != 1 {
		t.Fatalf("fetches = %d, want 1 (must stop when next==cursor)", got)
	}
}
