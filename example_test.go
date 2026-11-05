package qeetid_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	qeetid "github.com/qeetgroup/qeet-id-go"
)

// fakeQeetID stands in for the Qeet ID API in these examples, answering each
// path with a fixed JSON body.
func fakeQeetID(routes map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		if c := r.URL.Query().Get("cursor"); c != "" {
			key += "?cursor=" + c
		}
		body, ok := routes[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":{"code":"not_found","message":"We couldn't find what you're looking for.","request_id":"req-404"}}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}))
}

func ExampleNewClient() {
	client, err := qeetid.NewClient(
		qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")),
		qeetid.WithTimeout(10*time.Second),
		qeetid.WithMaxRetries(3),
		qeetid.WithLogger(slog.Default()),
		qeetid.WithUserAgentSuffix("billing-worker/2.3.0"),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = client
}

func ExampleWithClientCredentials() {
	// A service principal's tokens are fetched and renewed by the client.
	client, err := qeetid.NewClient(qeetid.WithClientCredentials(
		os.Getenv("QEETID_CLIENT_ID"),
		os.Getenv("QEETID_CLIENT_SECRET"),
		"user.read", "user.write",
	))
	if err != nil {
		log.Fatal(err)
	}
	_ = client
}

func ExampleUsersService_All() {
	api := fakeQeetID(map[string]string{
		"GET /v1/users":           `{"items":[{"id":"u1","email":"ada@brightleaffoods.in"},{"id":"u2","email":"bo@brightleaffoods.in"}],"next_cursor":"c2"}`,
		"GET /v1/users?cursor=c2": `{"items":[{"id":"u3","email":"cy@brightleaffoods.in"}]}`,
	})
	defer api.Close()
	client, err := qeetid.NewClient(qeetid.WithAPIKey("qk_example"), qeetid.WithBaseURL(api.URL))
	if err != nil {
		log.Fatal(err)
	}

	for user, err := range client.Users.All(context.Background(), &qeetid.ListUsersParams{Limit: 2}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.Email)
	}
	// Output:
	// ada@brightleaffoods.in
	// bo@brightleaffoods.in
	// cy@brightleaffoods.in
}

func ExampleAPIError() {
	api := fakeQeetID(nil) // answers every request with 404
	defer api.Close()
	client, err := qeetid.NewClient(qeetid.WithAPIKey("qk_example"), qeetid.WithBaseURL(api.URL))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Users.Get(context.Background(), "u-missing")
	if errors.Is(err, qeetid.ErrNotFound) {
		fmt.Println("no such user")
	}
	if apiErr, ok := errors.AsType[*qeetid.APIError](err); ok {
		fmt.Println(apiErr.StatusCode, apiErr.Code, apiErr.RequestID)
	}
	// Output:
	// no such user
	// 404 not_found req-404
}

func ExampleWithIdempotencyKey() {
	api := fakeQeetID(map[string]string{"POST /v1/users": `{"id":"u1","email":"ada@brightleaffoods.in","status":"active"}`})
	defer api.Close()
	client, err := qeetid.NewClient(qeetid.WithAPIKey("qk_example"), qeetid.WithBaseURL(api.URL))
	if err != nil {
		log.Fatal(err)
	}

	// The same key on a later attempt, even from another process, returns the
	// first answer instead of creating the user twice.
	user, err := client.Users.Create(context.Background(), qeetid.CreateUserParams{
		TenantID: "3f6c2f4e-8a1b-4c7d-9e2f-5a6b7c8d9e01",
		Email:    "ada@brightleaffoods.in",
	}, qeetid.WithIdempotencyKey("hr-import-2026-11-05-row-17"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.ID, user.Status)
	// Output: u1 active
}

func ExampleLDAPService_Authenticate() {
	api := fakeQeetID(map[string]string{
		"POST /ldap/conn-1/authenticate": `{"mfa_required":true,"mfa_token":"mt-1","methods":["totp"],"factors":[{"type":"totp"}]}`,
	})
	defer api.Close()
	client, err := qeetid.NewClient(qeetid.WithBaseURL(api.URL)) // a public endpoint
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.LDAP.Authenticate(context.Background(), "conn-1", qeetid.LDAPAuthenticateParams{
		Username: "ada", Password: "correct horse battery staple",
	})
	if err != nil {
		log.Fatal(err)
	}
	switch r := res.(type) {
	case *qeetid.TokenPair:
		fmt.Println("signed in, session", r.SessionID)
	case *qeetid.MFAChallenge:
		fmt.Println("second factor needed:", r.Methods)
	}
	// Output: second factor needed: [totp]
}

func ExampleOAuthService_ClientCredentials() {
	api := fakeQeetID(map[string]string{
		"POST /v1/oauth/token": `{"access_token":"eyJ...","token_type":"Bearer","expires_in":600,"scope":"user.read"}`,
	})
	defer api.Close()
	client, err := qeetid.NewClient(qeetid.WithBaseURL(api.URL))
	if err != nil {
		log.Fatal(err)
	}

	tok, err := client.OAuth.ClientCredentials(context.Background(), qeetid.ClientCredentialsParams{
		ClientID: "sp-1", ClientSecret: "secret", Scope: "user.read",
	})
	if err != nil {
		var oauthErr *qeetid.OAuthError
		if errors.As(err, &oauthErr) && oauthErr.Code == qeetid.OAuthInvalidClient {
			log.Fatal("wrong client ID or secret")
		}
		log.Fatal(err)
	}
	fmt.Println(tok.TokenType, tok.ExpiresIn, tok.Scope)
	// Output: Bearer 600 user.read
}
