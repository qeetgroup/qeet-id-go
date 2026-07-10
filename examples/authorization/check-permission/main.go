// Command check-permission runs an RBAC check and prints the grant path
// that decided it (Explain), not just the boolean.
//
//	QEETID_API_KEY=qk_… go run ./examples/authorization/check-permission <user> <tenant> <permission>
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qeetgroup/qeet-id-go"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("usage: check-permission <user-id> <tenant-id> <permission>")
	}
	check := qeetid.PermissionCheck{User: os.Args[1], Tenant: os.Args[2], Permission: os.Args[3]}

	client := qeetid.New(qeetid.Config{APIKey: os.Getenv("QEETID_API_KEY")})
	ctx := context.Background()

	explanation, err := client.Permissions.Explain(ctx, check)
	if err != nil {
		log.Fatalf("explain: %v", err)
	}

	if !explanation.Allowed {
		fmt.Println("denied:", explanation.Reason)
		return
	}
	fmt.Println("allowed via:")
	for _, step := range explanation.Paths {
		fmt.Printf("  %s granted by %s (%s)\n", step.Permission, step.GrantedBy, step.Via)
	}
}
