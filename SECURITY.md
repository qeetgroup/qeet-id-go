# Security Policy

## Reporting a vulnerability

Please report security vulnerabilities privately — do not open a public GitHub issue. Email
**security@qeet.in** with:

- a description of the vulnerability and its impact;
- steps to reproduce (a minimal program using this SDK);
- the SDK version (`qeetid.Version`) and Go version.

We aim to acknowledge reports within 2 business days.

## Scope

This policy covers the `qeet-id-go` module. The most sensitive code is:

- `webhook/` — webhook signature and timestamp verification;
- `jwt/` — access token verification and the signing-key cache;
- `auth.go` and `transport.go` — credentials, the client credentials token cache, and the retry
  rules that keep writes from being applied twice.

Vulnerabilities in the Qeet ID service itself should be reported through Qeet ID's own security
channel at the same address.

## Supported versions

Until v1.0.0, only the latest minor version receives security fixes. From v1.0.0, the two most
recent minor versions do. The SDK supports the two most recent Go releases.

## Disclosure

We credit reporters, unless they ask otherwise, in the release notes once a fix ships.
