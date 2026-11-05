// Command quickstart creates a user, gives them a role and lists the
// tenant's users with an API key.
//
//	QEETID_API_KEY=qk_... QEETID_TENANT_ID=... go run ./examples/quickstart
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	qeetid "github.com/qeetgroup/qeet-id-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := qeetid.NewClient(qeetid.WithAPIKey(os.Getenv("QEETID_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	tenantID := os.Getenv("QEETID_TENANT_ID")

	user, err := client.Users.Create(ctx, qeetid.CreateUserParams{
		TenantID:    tenantID,
		Email:       "ada@brightleaffoods.in",
		DisplayName: "Ada Lovelace",
	})
	switch {
	case errors.Is(err, qeetid.ErrConflict):
		log.Fatal("a user with this email already exists")
	case err != nil:
		log.Fatal(err)
	}
	fmt.Println("created", user.ID)

	roles, err := client.Roles.List(ctx, tenantID)
	if err != nil {
		log.Fatal(err)
	}
	for _, role := range roles.Items {
		if role.Name == "editor" {
			if err := client.Roles.AssignToUser(ctx, user.ID, tenantID, role.ID); err != nil {
				log.Fatal(err)
			}
			fmt.Println("assigned", role.Name)
		}
	}

	for u, err := range client.Users.All(ctx, &qeetid.ListUsersParams{Limit: 100}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u.Email, u.Status)
	}
}
