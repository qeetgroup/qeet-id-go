// Command mfa clears a user's MFA factors — the admin-initiated reset for
// when someone loses their authenticator device. The backend has no
// endpoint to list a user's factors as an admin, only this reset.
//
//	QEETID_API_KEY=qk_… go run ./examples/authentication/mfa <user-id>
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qeetgroup/qeet-id-go"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: mfa <user-id>")
	}
	userID := os.Args[1]

	client := qeetid.New(qeetid.Config{APIKey: os.Getenv("QEETID_API_KEY")})
	ctx := context.Background()

	if err := client.MFA.Reset(ctx, userID); err != nil {
		log.Fatalf("reset: %v", err)
	}
	fmt.Println("MFA factors cleared for", userID)
}
