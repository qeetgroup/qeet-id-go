// Command auditlogs runs a free-text search over the audit log and verifies
// the hash chain is intact.
//
//	QEETID_API_KEY=qk_… go run ./examples/administration/auditlogs <tenant-id> <query>
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qeetgroup/qeet-id-go"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal(`usage: auditlogs <tenant-id> "<query>"`)
	}
	tenantID, query := os.Args[1], os.Args[2]

	client := qeetid.New(qeetid.Config{APIKey: os.Getenv("QEETID_API_KEY")})
	ctx := context.Background()

	page, err := client.AuditLogs.List(ctx, tenantID, qeetid.AuditLogListParams{Search: query, Limit: 20})
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	for _, e := range page.Data {
		fmt.Printf("%s %s %s\n", e.CreatedAt, e.Action, e.ResourceType)
	}

	verification, err := client.AuditLogs.Verify(ctx, tenantID)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Printf("chain intact: %v (%d rows checked)\n", verification.OK, verification.RowsChecked)
}
