package qeetid

import (
	"context"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	for _, tc := range routeCases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(http.StatusOK, `{"items":[]}`)}}
			c := newFakeClient(t, f)
			if err := tc.call(context.Background(), c); err != nil {
				t.Fatalf("call: %v", err)
			}
			calls := f.calls()
			if len(calls) != 1 {
				t.Fatalf("sent %d requests, want 1", len(calls))
			}
			got := calls[0]
			if got.method != tc.method || got.path != tc.path {
				t.Errorf("sent %s %s, want %s %s", got.method, got.path, tc.method, tc.path)
			}
			auth := got.header.Get("Authorization")
			switch {
			case tc.public && auth != "":
				t.Errorf("public endpoint got Authorization %q, want none", auth)
			case !tc.public && auth != "ApiKey qk_test.secret":
				t.Errorf("Authorization = %q, want the API key", auth)
			}
			if key := got.header.Get("Idempotency-Key"); (key != "") != tc.idempotent {
				t.Errorf("Idempotency-Key = %q, want one: %v", key, tc.idempotent)
			}
		})
	}
}

// TestRoutesCoverEveryService fails when a service has no route case, which
// happens when a method is added by hand without a test.
func TestRoutesCoverEveryService(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	for _, tc := range routeCases {
		for i := range len(tc.name) {
			if tc.name[i] == '.' {
				covered[tc.name[:i]] = true
				break
			}
		}
	}
	for _, name := range serviceNames(c) {
		if !covered[name] {
			t.Errorf("service %s has no route case", name)
		}
	}
}
