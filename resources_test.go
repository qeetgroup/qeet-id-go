package qeetid

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAllWalksEveryPage(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{
		func(r *http.Request) *http.Response {
			switch r.URL.Query().Get("cursor") {
			case "":
				return response(200, `{"items":[{"id":"u1"},{"id":"u2"}],"next_cursor":"c2","total":3}`)
			case "c2":
				return response(200, `{"items":[{"id":"u3"}],"next_cursor":""}`)
			}
			return response(400, `{"error":{"code":"bad_request","message":"invalid cursor"}}`)
		},
	}}
	c := newFakeClient(t, f)
	var ids []string
	for u, err := range c.Users.All(context.Background(), &ListUsersParams{Limit: 2, Status: UserStatusActive}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, u.ID)
	}
	if strings.Join(ids, ",") != "u1,u2,u3" {
		t.Fatalf("ids = %v", ids)
	}
	calls := f.calls()
	if len(calls) != 2 {
		t.Fatalf("fetched %d pages, want 2", len(calls))
	}
	q, _ := url.ParseQuery(calls[1].query)
	if q.Get("cursor") != "c2" || q.Get("limit") != "2" || q.Get("status") != "active" {
		t.Errorf("second page query = %s", calls[1].query)
	}
}

func TestAllStopsWhenTheLoopBreaks(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(200, `{"items":[{"id":"a"},{"id":"b"}],"next_cursor":"next"}`)}}
	c := newFakeClient(t, f)
	for _, err := range c.Audit.All(context.Background(), "t1", nil) {
		if err != nil {
			t.Fatal(err)
		}
		break
	}
	if n := len(f.calls()); n != 1 {
		t.Fatalf("fetched %d pages after break, want 1", n)
	}
}

func TestAllYieldsErrorOnce(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{
		answer(200, `{"items":[{"id":"a"}],"next_cursor":"c2"}`),
		answer(403, `{"error":{"code":"forbidden","message":"no"}}`),
	}}
	c := newFakeClient(t, f)
	var items, errs int
	for g, err := range c.OAuthGrants.All(context.Background(), "t1", &ListOAuthGrantsParams{ClientID: "app"}) {
		if err != nil {
			errs++
			if g != nil || !errors.Is(err, ErrForbidden) {
				t.Fatalf("(%v, %v), want (nil, forbidden)", g, err)
			}
			continue
		}
		items++
	}
	if items != 1 || errs != 1 {
		t.Fatalf("items = %d, errors = %d", items, errs)
	}
}

func TestAllStopsOnARepeatedCursor(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{
		answer(200, `{"items":[{"id":"a"}],"next_cursor":"c1"}`),
		answer(200, `{"items":[{"id":"b"}],"next_cursor":"c2"}`),
		answer(200, `{"items":[{"id":"c"}],"next_cursor":"c1"}`),
	}}
	c := newFakeClient(t, f)
	n := 0
	for _, err := range c.Webhooks.AllDeliveries(context.Background(), "w1", nil) {
		if err != nil {
			t.Fatal(err)
		}
		n++
	}
	if n != 3 || len(f.calls()) != 3 {
		t.Fatalf("items = %d, pages = %d, want 3 and 3", n, len(f.calls()))
	}
}

func TestCappedListReportsTruncated(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(200, `{"items":[{"id":"a1","name":"triage"}],"truncated":true}`)}}
	c := newFakeClient(t, f)
	page, err := c.Agents.List(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if !page.Truncated || page.HasMore() || len(page.Items) != 1 || page.Items[0].Name != "triage" {
		t.Fatalf("page = %+v", page)
	}
}

func TestInvitationCreateCarriesNoToken(t *testing.T) {
	// The server no longer returns the invitation's token: the accept link goes
	// only to the invited address.
	if _, ok := hasJSONField(InvitationResult{}, "token"); ok {
		t.Fatal("InvitationResult has a token field")
	}
	body := `{"invite":{"id":"i1","tenant_id":"t1","email":"ada@example.com","role_id":null,"status":"pending",
		"expires_at":"2026-11-12T10:00:00Z","created_at":"2026-11-05T10:00:00Z","accepted_at":null,
		"delivery":{"status":"not_sent","code":"email.suppressed","message":"The address bounced before."}}}`
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(201, body)}}
	c := newFakeClient(t, f)
	res, err := c.Invitations.Create(context.Background(), CreateInvitationParams{TenantID: "t1", Email: "ada@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Invite.ID != "i1" || res.Invite.Delivery.Status != "not_sent" || res.Invite.Delivery.Code != "email.suppressed" {
		t.Fatalf("result = %+v", res)
	}
	var sent map[string]any
	if err := json.Unmarshal([]byte(f.calls()[0].body), &sent); err != nil {
		t.Fatal(err)
	}
	if _, ok := sent["role_id"]; ok {
		t.Errorf("role_id sent although unset: %s", f.calls()[0].body)
	}
}

func hasJSONField(v any, name string) (string, bool) {
	b, _ := json.Marshal(v)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	_, ok := m[name]
	return name, ok
}

func TestLDAPAuthenticate(t *testing.T) {
	pair := `{"access_token":"at","token_type":"Bearer","expires_at":"2026-11-05T10:10:00Z","refresh_token":"rt",
		"session_expires_at":"2026-11-06T10:00:00Z","idle_expires_at":"2026-11-05T11:00:00Z",
		"session_id":"s1","user_id":"u1","tenant_id":"t1"}`
	challenge := `{"mfa_required":true,"mfa_token":"mt","methods":["totp","sms"],
		"factors":[{"type":"totp"},{"type":"sms","id":"f1","destination":"+91 ••••• •4321"}]}`
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(200, pair), answer(200, challenge)}}
	c := newFakeClient(t, f)
	ctx := context.Background()
	params := LDAPAuthenticateParams{Username: "ada", Password: "correct horse"}

	res, err := c.LDAP.Authenticate(ctx, "conn-1", params)
	if err != nil {
		t.Fatal(err)
	}
	tp, ok := res.(*TokenPair)
	if !ok || tp.AccessToken != "at" || tp.RefreshToken != "rt" || tp.TenantID != "t1" {
		t.Fatalf("first result = %#v, want a *TokenPair", res)
	}

	res, err = c.LDAP.Authenticate(ctx, "conn-1", params)
	if err != nil {
		t.Fatal(err)
	}
	ch, ok := res.(*MFAChallenge)
	if !ok || ch.MFAToken != "mt" || len(ch.Factors) != 2 || ch.Factors[1].ID != "f1" {
		t.Fatalf("second result = %#v, want an *MFAChallenge", res)
	}

	got := f.calls()[0]
	if got.path != "/ldap/conn-1/authenticate" || got.header.Get("Authorization") != "" {
		t.Errorf("request = %s %s with Authorization %q", got.method, got.path, got.header.Get("Authorization"))
	}
}

func TestUpdateGroupParentStates(t *testing.T) {
	parent := "g-parent"
	cases := []struct {
		params UpdateGroupParams
		want   string
	}{
		{UpdateGroupParams{Name: new("eng")}, `{"name":"eng"}`},
		{UpdateGroupParams{ParentID: &parent}, `{"parent_id":"g-parent"}`},
		{UpdateGroupParams{MoveToTop: true}, `{"parent_id":null}`},
	}
	for _, tc := range cases {
		b, err := json.Marshal(tc.params)
		if err != nil || string(b) != tc.want {
			t.Errorf("Marshal(%+v) = %s, %v; want %s", tc.params, b, err, tc.want)
		}
	}
	if _, err := json.Marshal(UpdateGroupParams{ParentID: &parent, MoveToTop: true}); err == nil {
		t.Error("ParentID with MoveToTop accepted")
	}
}

func TestUsersImport(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(200, `{"succeeded":2,"failed":0,"created":[{"line":1,"email":"a@x.in","id":"u1"}]}`)}}
	c := newFakeClient(t, f)
	file := []byte("{\"email\":\"a@x.in\"}\n{\"email\":\"b@x.in\"}\n")
	res, err := c.Users.Import(context.Background(), ImportSourceAuth0, file)
	if err != nil {
		t.Fatal(err)
	}
	got := f.calls()[0]
	if got.query != "source=auth0" || got.header.Get("Content-Type") != "application/x-ndjson" || got.body != string(file) {
		t.Errorf("request = %+v", got)
	}
	if res.Succeeded != 2 || len(res.Created) != 1 {
		t.Errorf("result = %+v", res)
	}
	if _, err := c.Users.Import(context.Background(), "okta", file); err == nil {
		t.Error("unknown source accepted")
	}
}

func TestQueryEncoding(t *testing.T) {
	f := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(200, `{"allowed":true}`)}}
	c := newFakeClient(t, f)
	res, err := c.Permissions.Check(context.Background(), &CheckPermissionParams{
		UserID: "u1", TenantID: "t1", Permission: "billing.write", Explain: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Allowed {
		t.Error("Allowed = false")
	}
	q, _ := url.ParseQuery(f.calls()[0].query)
	if q.Get("user_id") != "u1" || q.Get("tenant_id") != "t1" || q.Get("permission") != "billing.write" || q.Get("explain") != "true" {
		t.Errorf("query = %v", q)
	}
	f2 := &fakeAPI{answers: []func(*http.Request) *http.Response{answer(200, `{"items":[]}`)}}
	c2 := newFakeClient(t, f2)
	if _, err := c2.Users.List(context.Background(), &ListUsersParams{EmailVerified: new(false)}); err != nil {
		t.Fatal(err)
	}
	if q := f2.calls()[0].query; q != "email_verified=false" {
		t.Errorf("query = %q, want email_verified=false", q)
	}
}
