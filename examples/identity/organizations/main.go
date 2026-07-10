// Command organizations creates an organization (tenant) and lists the
// first page of existing ones.
//
//	QEETID_API_KEY=qk_… go run ./examples/identity/organizations
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

	org, err := client.Organizations.Create(ctx, qeetid.CreateOrganizationInput{
		Name: "Acme Corp", Slug: "acme",
	})
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	fmt.Println("created", org.ID, org.Slug)

	page, err := client.Organizations.List(ctx, 20, "")
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	for _, o := range page.Data {
		fmt.Println("-", o.Name)
	}
}
