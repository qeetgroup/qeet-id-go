package qeetid

import "github.com/qeetgroup/qeet-id-go/internal/auth"

// Claims is the verified content of a Qeet-issued token.
type Claims = auth.Claims

// VerifyOptions tightens verification. ClockSkew defaults to 30s.
type VerifyOptions = auth.VerifyOptions
