// Command verify-session verifies a Qeet ID session token locally (against
// the published JWKS) and runs a permission check.
//
//	QEETID_API_KEY=qk_… QEETID_TOKEN=<jwt> go run ./examples/authentication/verify-session
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qeetgroup/qeet-id-go"
)

func main() {
	client := qeetid.New(qeetid.Config{APIKey: os.Getenv("QEETID_API_KEY")})
	ctx := context.Background()

	claims, err := client.Sessions.Verify(ctx, os.Getenv("QEETID_TOKEN"))
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Printf("user=%s tenant=%s scope=%q\n", claims.UserID, claims.TenantID, claims.Scope)

	ok, err := client.Permissions.Check(ctx, qeetid.PermissionCheck{
		User:       claims.UserID,
		Tenant:     claims.TenantID,
		Permission: "billing:write",
	})
	if err != nil {
		log.Fatalf("check: %v", err)
	}
	fmt.Println("billing:write allowed:", ok)
}
