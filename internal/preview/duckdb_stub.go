//go:build !cgo

package preview

import "errors"

// DuckDBFormatter is a no-op stub when CGO is disabled.
type DuckDBFormatter struct {
	Query string
	Width int
}

// Format always returns an error when CGO is not available.
func (f DuckDBFormatter) Format(_ string) (string, error) {
	return "", errors.New("DuckDB not available (CGO disabled)")
}
