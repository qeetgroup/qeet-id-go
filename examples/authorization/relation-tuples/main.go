// Command relation-tuples grants a ReBAC relationship and then expands the
// identity graph rooted at it — the same data the console's Identity Graph
// visualization renders.
//
//	QEETID_API_KEY=qk_… go run ./examples/authorization/relation-tuples <tenant-id>
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
		log.Fatal("usage: relation-tuples <tenant-id>")
	}
	tenantID := os.Args[1]

	client := qeetid.New(qeetid.Config{APIKey: os.Getenv("QEETID_API_KEY")})
	ctx := context.Background()
	rel := client.Relationships

	if _, err := rel.Create(ctx, tenantID, qeetid.CreateTupleInput{
		Object: "document:readme", Relation: "viewer", Subject: "group:eng#member",
	}); err != nil {
		log.Fatalf("create tuple: %v", err)
	}

	result, err := rel.Check(ctx, tenantID, qeetid.CheckRelationInput{
		Object: "document:readme", Relation: "viewer", UserID: "user-in-eng-group",
	}, true)
	if err != nil {
		log.Fatalf("check: %v", err)
	}
	fmt.Println("allowed:", result.Allowed)

	graph, err := rel.Graph(ctx, tenantID, "document:readme", "viewer", 5)
	if err != nil {
		log.Fatalf("graph: %v", err)
	}
	for _, n := range graph.Nodes {
		fmt.Printf("node: %s (%s)\n", n.Label, n.Type)
	}
}
