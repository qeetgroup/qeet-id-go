# Contributing

## Set up

```bash
git clone https://github.com/qeetgroup/qeet-id-go
cd qeet-id-go
make build test
```

Requires Go 1.26 or 1.27. `make lint` needs golangci-lint v2.15.0; `make vuln` fetches govulncheck.

## Before you open a pull request

```bash
make ci   # gofmt, go vet, go build, go test -race, golangci-lint, govulncheck
```

## Conventions

- **Match the API contract.** Every endpoint, field, type and enum comes from the Qeet ID server's
  OpenAPI documents; where they are loose, follow the server's handler. Never add a method for an
  endpoint that does not exist.
- **Standard library only.** A new dependency needs the maintainers' decision.
- **Method shape:** `context.Context` first, path arguments in path order, a `...Params` value for
  the request body, a pointer for optional query parameters, `opts ...RequestOption` last. Lists
  return `*Page[T]`; cursor lists also get an `All…` iterator.
- **Documentation:** every exported name has a doc comment that starts with its name; each method's
  comment ends with `API: VERB /path`.
- **Tests:** every method has a route case in `routes_cases_test.go` or `routes_test.go`. Use the
  in-process fakes in `helpers_test.go`, not a mock library.
- **Security code** (`auth.go`, `transport.go`, `webhook/`, `jwt/`) needs a second reviewer.
- **No tracker references** in code, comments, tests or docs; describe the behaviour instead.

See [docs/llm/workflows.md](docs/llm/workflows.md) for step-by-step guides and
[ADR-0004](docs/design-decisions/ADR-0004-sdk-design.md) for the design.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org): `feat(scope): …`,
`fix(scope): …`, `docs: …`, `chore(deps): …`. Mark breaking changes with `!` and a
`BREAKING CHANGE:` footer.

## Reporting bugs

Open an issue with the method you called, what you expected, what happened, and the `RequestID` of
the `*qeetid.APIError` if there was one. Never include a credential.
