# Security Policy

## Reporting a vulnerability

Please report security vulnerabilities privately — do not open a public
GitHub issue. Email **security@qeet.in** with:

- A description of the vulnerability and its impact.
- Steps to reproduce (a minimal code sample against this SDK, not the
  Qeet ID backend).
- The SDK version (`qeetid.Version`) and Go version affected.

We aim to acknowledge reports within 2 business days.

## Scope

This policy covers the `qeet-id-go` SDK itself — the JWT/JWKS verification
logic (`internal/auth`), webhook HMAC verification (`webhooks_verify.go`),
and the HTTP transport (`internal/transport`) are the highest-sensitivity
surfaces. Vulnerabilities in the Qeet ID platform/backend itself should be
reported through Qeet ID's own security channel, not here.

## Supported versions

Only the latest minor version receives security fixes pre-1.0. Once the SDK
reaches 1.0, the most recent two minor versions will be supported.

## Disclosure

We'll credit reporters (unless anonymity is requested) in the release notes
once a fix ships.
