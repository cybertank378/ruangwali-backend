// =========================================================
// File: cmd/api/openapi/embed.go
// =========================================================

package openapi

import (
	_ "embed"
)

//go:embed openapi.yaml
var Specification []byte
