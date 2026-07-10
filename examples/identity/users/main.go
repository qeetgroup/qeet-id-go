// Command users auto-paginates every user in a tenant using the All
// iterator (Go 1.23+ range-over-func).
//
//	QEETID_API_KEY=qk_… QEETID_TENANT_ID=<uuid> go run ./examples/identity/users
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

	count := 0
	for user, err := range client.Users.All(ctx, qeetid.ListParams{Tenant: os.Getenv("QEETID_TENANT_ID")}) {
		if err != nil {
			log.Fatalf("list users: %v", err)
		}
		fmt.Println(user.Email)
		count++
	}
	fmt.Printf("%d users total\n", count)
}
