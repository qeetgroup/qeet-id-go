# ADR-0001: Single package, file-per-resource

**Status:** Accepted

## Context

The SDK needed a coverage expansion from ~24 to ~90 methods. Before adding
more, the folder structure came under review: should each resource be its
own importable Go package (`qeetid/users`, `qeetid/oauth`, ...), or should
everything stay in one package?

## Decision

One package (`qeetid`), one file per resource (splitting into
`_models`/`_requests`/`_responses` companions once a resource gets large).

## Why

- **No import cycles.** A root `Client` aggregator needs to reference every
  resource. If each resource were its own package, the root package would
  import all of them — and none of them could import the root package back
  (for the `Config`/`Error` types) without a cycle. The workaround (push
  shared types into `internal/` and re-export via type aliases) is exactly
  the layering this SDK already does *within* one package for HTTP/crypto
  machinery — doing it *again* at the public-package boundary for zero
  additional benefit isn't worth the complexity.
- **WorkOS precedent.** WorkOS's Go SDK covers ~17 namespaces in one flat
  package and is a well-regarded server-SDK design at a comparable surface
  size to this one.
- **No consumer benefit foregone.** Nothing about per-resource packages
  gives a caller anything they don't already get from `client.Users.Create(...)`.

## Rejected alternative: package-per-resource (Clerk v2 style)

Clerk v2 splits each resource into its own package
(`github.com/clerk/clerk-sdk-go/v2/user`, `.../organization`, ...) with no
aggregate client — callers import only the packages they use, at the cost of
more import statements for multi-resource code and the cycle-avoidance
machinery above. This was implemented and then reverted during this SDK's
development, before any release, once the cycle cost became concrete rather
than theoretical.
