// Package version holds the SDK's release version, isolated from the root
// package so internal/transport can reference it (via internal/constants)
// without importing the root package and creating an import cycle.
package version

// SDKVersion is the released version of this SDK. The root package
// re-exports it as qeetid.Version.
const SDKVersion = "0.1.0"
