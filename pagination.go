// Every listable resource exposes an All(ctx, ...) method returning an
// iter.Seq2[T, error] (Go 1.23+ range-over-func), walking cursor pages
// lazily via internal/pagination.Paginate. Range over it and check the error
// on each step; breaking out of the loop stops paging immediately:
//
//	for user, err := range client.Users.All(ctx, qeetid.ListParams{Tenant: id}) {
//		if err != nil {
//			return err
//		}
//		fmt.Println(user.Email)
//	}
//
// The individual All() methods live alongside each resource's other methods
// (e.g. UsersService.All in users.go), not in this file — this file exists
// only to hold that shared doc comment in one place.
package qeetid
