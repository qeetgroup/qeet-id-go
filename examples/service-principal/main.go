// Command service-principal calls the management API as a service principal,
// with short-lived tokens from the OAuth client credentials grant instead of a
// long-lived API key.
//
//	QEETID_CLIENT_ID=... QEETID_CLIENT_SECRET=... QEETID_TENANT_ID=... go run ./examples/service-principal
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	qeetid "github.com/qeetgroup/qeet-id-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := qeetid.NewClient(
		qeetid.WithClientCredentials(os.Getenv("QEETID_CLIENT_ID"), os.Getenv("QEETID_CLIENT_SECRET"), "group.read"),
		qeetid.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
	if err != nil {
		log.Fatal(err)
	}

	groups, err := client.Groups.List(ctx, os.Getenv("QEETID_TENANT_ID"), nil)
	if oauthErr, ok := errors.AsType[*qeetid.OAuthError](err); ok {
		log.Fatalf("the token endpoint refused the service principal: %s", oauthErr.Code)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, g := range groups.Items {
		fmt.Println(g.Name)
	}
	if groups.Truncated {
		fmt.Println("(the list was cut short; narrow it with a filter)")
	}
}
