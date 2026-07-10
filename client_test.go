package qeetid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// TestNew_WiresEveryService reflects over Client and fails loudly if New()
// ever leaves a *XService field nil — the failure mode of adding a new
// resource to the Client struct literal but forgetting to construct it,
// which would otherwise surface as a runtime nil-pointer panic deep in a
// caller's code instead of here.
func TestNew_WiresEveryService(t *testing.T) {
	c := New(Config{APIKey: "qk_x.y"})

	v := reflect.ValueOf(c).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue // unexported (transport)
		}
		if f.Kind() == reflect.Ptr && f.IsNil() {
			t.Errorf("Client.%s is nil", v.Type().Field(i).Name)
		}
	}
}

// --- Regression tests: every path bug found by the backend-endpoint audit ---
// Each of these targeted a route that either 404s (fictional path) or hits
// the wrong resource (missing/extra tenant segment). Verified directly
// against the qeet-id backend's router source, not just the OpenAPI spec.

func TestRoles_UsesRealTenantScopedPaths(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"id":"r1","name":"editor"}`)

	if _, err := c.Roles.Create(context.Background(), "t1", CreateRoleInput{Name: "editor"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rec.path != "/v1/tenants/t1/roles" {
		t.Fatalf("Create path = %q, want /v1/tenants/t1/roles (no /v1/rbac/roles exists on the backend)", rec.path)
	}

	if err := c.Roles.AssignToUser(context.Background(), "u1", "t1", "r1"); err != nil {
		t.Fatalf("AssignToUser: %v", err)
	}
	if rec.path != "/v1/users/u1/tenants/t1/roles/r1" || rec.method != http.MethodPost {
		t.Fatalf("AssignToUser -> %s %s", rec.method, rec.path)
	}

	if err := c.Roles.GrantPermission(context.Background(), "r1", "p1"); err != nil {
		t.Fatalf("GrantPermission: %v", err)
	}
	if rec.path != "/v1/roles/r1/permissions/p1" {
		t.Fatalf("GrantPermission path = %q", rec.path)
	}
}

func TestPermissions_ListUsesUnprefixedPath(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"items":[{"id":"p1","key":"billing:write"}]}`)

	perms, err := c.Permissions.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.path != "/v1/permissions" {
		t.Fatalf("path = %q, want /v1/permissions (no /v1/rbac/permissions exists on the backend)", rec.path)
	}
	if len(perms) != 1 || perms[0].Key != "billing:write" {
		t.Fatalf("perms = %+v", perms)
	}
}

func TestWebhooks_CreateAndGetAreNotTenantScoped(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"id":"wh1","url":"https://x","events":["user.created"],"secret":"whsec_once"}`)

	wh, err := c.Webhooks.Create(context.Background(), CreateWebhookInput{URL: "https://x", Events: []string{"user.created"}})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rec.path != "/v1/webhooks" {
		t.Fatalf("Create path = %q, want /v1/webhooks (no tenant segment — the previous version wrongly added /tenants/{id})", rec.path)
	}
	if wh.Secret != "whsec_once" {
		t.Fatalf("secret not surfaced on create: %q", wh.Secret)
	}

	if _, err := c.Webhooks.Get(context.Background(), "wh1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.path != "/v1/webhooks/wh1" {
		t.Fatalf("Get path = %q", rec.path)
	}
}

func TestAPIKeys_CreateResponseEnvelopeAndTenantScopedList(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"api_key":{"id":"k1","tenant_id":"t1","name":"ci","prefix":"qk_abc"},"secret":"qk_abc.def","warning":"shown once"}`)

	res, err := c.APIKeys.Create(context.Background(), CreateAPIKeyInput{TenantID: "t1", Name: "ci"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.APIKey.ID != "k1" || res.Secret != "qk_abc.def" {
		t.Fatalf("Create response = %+v (the previous version's json tag \"key\" instead of \"api_key\" would leave APIKey zero-valued)", res)
	}

	if _, err := c.APIKeys.List(context.Background(), "t1"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.path != "/v1/tenants/t1/api-keys" {
		t.Fatalf("List path = %q, want tenant-scoped", rec.path)
	}
}

func TestGroups_TenantScopedListAndPathOnlyAddMember(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"items":[{"id":"g1","name":"eng"}]}`)

	if _, err := c.Groups.List(context.Background(), "t1"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.path != "/v1/tenants/t1/groups" {
		t.Fatalf("List path = %q, want tenant-scoped", rec.path)
	}

	if err := c.Groups.AddMember(context.Background(), "g1", "u1"); err != nil {
		t.Fatalf("AddMember: %v", err)
	}
	if rec.path != "/v1/groups/g1/members/u1" || rec.method != http.MethodPost {
		t.Fatalf("AddMember -> %s %s (userID must be a path segment, not a body field)", rec.method, rec.path)
	}
}

func TestInvitations_TenantScopedList(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"items":[{"id":"i1","email":"a@b.com"}]}`)

	if _, err := c.Invitations.List(context.Background(), "t1"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.path != "/v1/tenants/t1/invites" {
		t.Fatalf("List path = %q, want tenant-scoped", rec.path)
	}
}

func TestAuthHooks_MultiRecordCollectionShape(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"id":"h1","trigger":"pre-login","url":"https://x","enabled":true}`)

	if _, err := c.AuthHooks.Create(context.Background(), "t1", CreateAuthHookInput{URL: "https://x", Secret: "s"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rec.path != "/v1/tenants/t1/auth-hooks" || rec.method != http.MethodPost {
		t.Fatalf("Create -> %s %s", rec.method, rec.path)
	}

	if err := c.AuthHooks.Update(context.Background(), "t1", "h1", UpdateAuthHookInput{Enabled: true, FailOpen: true}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if rec.path != "/v1/tenants/t1/auth-hooks/h1" || rec.method != http.MethodPatch {
		t.Fatalf("Update -> %s %s (must be PATCH on a per-hook ID, not PUT on the tenant)", rec.method, rec.path)
	}
}

func TestOIDC_UsesTenantScopedPaths(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"id":"c1"}`)
	if _, err := c.OIDC.Get(context.Background(), "tenant1", "client1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	want := "/v1/tenants/tenant1/oidc/clients/client1"
	if rec.path != want {
		t.Fatalf("path = %q, want %q (the previous SDK version 404ed here by hitting the top-level /v1/oidc/clients/{id} path)", rec.path, want)
	}

	if _, err := c.OIDC.List(context.Background(), "tenant1"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.path != "/v1/tenants/tenant1/oidc/clients" {
		t.Fatalf("List path = %q", rec.path)
	}

	if _, err := c.OIDC.RotateSecret(context.Background(), "tenant1", "client1"); err != nil {
		t.Fatalf("RotateSecret: %v", err)
	}
	if rec.path != "/v1/tenants/tenant1/oidc/clients/client1/rotate-secret" {
		t.Fatalf("RotateSecret path = %q", rec.path)
	}
}

func TestOAuth_IntrospectAndRevokeUseUnprefixedPath(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"active":true}`)

	if _, err := c.OAuth.Introspect(context.Background(), "tok"); err != nil {
		t.Fatalf("Introspect: %v", err)
	}
	if rec.path != "/oauth/introspect" {
		t.Fatalf("Introspect path = %q, want /oauth/introspect (the previous SDK version incorrectly used /v1/oauth/introspect)", rec.path)
	}

	if err := c.OAuth.Revoke(context.Background(), "tok"); err != nil {
		t.Fatalf("Revoke: %v", err)
	}
	if rec.path != "/oauth/revoke" {
		t.Fatalf("Revoke path = %q, want /oauth/revoke", rec.path)
	}
}

func TestAuditLogs_SendsSearchQueryParam(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"items":[],"next_cursor":""}`)

	_, err := c.AuditLogs.List(context.Background(), "t1", AuditLogListParams{Search: `"failed login" -saml`})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got := rec.query; !strings.Contains(got, "q=") {
		t.Fatalf("query = %q, missing q= param (this was never sent by the previous SDK version)", got)
	}
}

func TestAuditLogs_VerifyIsWholeChainNotPerEntry(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"ok":true,"rows_checked":42}`)

	res, err := c.AuditLogs.Verify(context.Background(), "t1")
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if rec.method != http.MethodGet {
		t.Fatalf("method = %q, want GET (the previous SDK version incorrectly POSTed to a nonexistent per-entry endpoint)", rec.method)
	}
	if rec.path != "/v1/tenants/t1/audit/verify" {
		t.Fatalf("path = %q", rec.path)
	}
	if !res.OK || res.RowsChecked != 42 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

// --- Representative resource-shape coverage -----------------------------

func TestUsersRequestShapes(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"id":"u1","email":"a@b.com","status":"active"}`)

	if _, err := c.Users.Create(context.Background(), CreateUserInput{Email: "a@b.com"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rec.method != http.MethodPost || rec.path != "/v1/users" {
		t.Fatalf("Create -> %s %s", rec.method, rec.path)
	}

	if err := c.Users.ResetMFA(context.Background(), "u1"); err != nil {
		t.Fatalf("ResetMFA: %v", err)
	}
	if rec.method != http.MethodDelete || rec.path != "/v1/users/u1/mfa" {
		t.Fatalf("ResetMFA -> %s %s", rec.method, rec.path)
	}
}

func TestPermissions_CheckAndExplain(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"allowed":true}`)

	ok, err := c.Permissions.Check(context.Background(), PermissionCheck{User: "u1", Tenant: "t1", Permission: "billing:write"})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatal("expected allowed")
	}
	if rec.path != "/v1/check" || strings.Contains(rec.query, "explain") {
		t.Fatalf("Check -> %s?%s", rec.path, rec.query)
	}

	if _, err := c.Permissions.Explain(context.Background(), PermissionCheck{User: "u1", Tenant: "t1", Permission: "billing:write"}); err != nil {
		t.Fatalf("Explain: %v", err)
	}
	if !strings.Contains(rec.query, "explain=true") {
		t.Fatalf("Explain query = %q, missing explain=true", rec.query)
	}
}

func TestRelationships_GraphShape(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"nodes":[{"id":"user:u1","type":"user","label":"u1"}],"edges":[]}`)

	g, err := c.Relationships.Graph(context.Background(), "t1", "document:readme", "viewer", 5)
	if err != nil {
		t.Fatalf("Graph: %v", err)
	}
	if rec.path != "/v1/tenants/t1/relation-tuples/graph" {
		t.Fatalf("path = %q", rec.path)
	}
	if !strings.Contains(rec.query, "depth=5") {
		t.Fatalf("query = %q, missing depth=5", rec.query)
	}
	if len(g.Nodes) != 1 {
		t.Fatalf("Nodes = %+v", g.Nodes)
	}
}

func TestSAMLProviders_TenantScopedShape(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"id":"sp1","name":"Acme","entity_id":"urn:acme","acs_url":"https://acme/acs","status":"active"}`)

	if _, err := c.SAMLProviders.Create(context.Background(), "t1", CreateSAMLServiceProviderInput{
		Name: "Acme", EntityID: "urn:acme", ACSURL: "https://acme/acs",
	}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rec.path != "/v1/tenants/t1/saml-providers" {
		t.Fatalf("path = %q", rec.path)
	}
}

func TestIPRules_CheckAndSetEnforcement(t *testing.T) {
	rec := &recorder{}
	c := recordingClient(t, rec, `{"enabled":true,"allowed":false,"reason":"denylisted"}`)

	res, err := c.IPRules.Check(context.Background(), "t1", "1.2.3.4")
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if rec.path != "/v1/tenants/t1/ip-rules/check" || res.Allowed {
		t.Fatalf("Check -> %s, result %+v", rec.path, res)
	}

	if err := c.IPRules.SetEnforcement(context.Background(), "t1", true); err != nil {
		t.Fatalf("SetEnforcement: %v", err)
	}
	if rec.path != "/v1/tenants/t1/ip-rules/config" || rec.method != http.MethodPut {
		t.Fatalf("SetEnforcement -> %s %s", rec.method, rec.path)
	}
}

// --- Shared test helpers -------------------------------------------------

type recorder struct {
	method string
	path   string
	body   string
	query  string
}

func recordingClient(t *testing.T, rec *recorder, response string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method, rec.path, rec.query = r.Method, r.URL.Path, r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)
	return New(Config{APIKey: "qk_x.y", BaseURL: srv.URL})
}
