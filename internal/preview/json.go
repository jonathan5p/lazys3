package preview

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/pretty"
)

const maxPreviewLines = 500

// JSONFormatter formats a JSON file with colorized pretty-printing.
type JSONFormatter struct{}

// Format reads the file at path, validates it as JSON, and returns pretty output.
func (f JSONFormatter) Format(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	if !json.Valid(data) {
		return "", fmt.Errorf("invalid JSON in %s", path)
	}
	result := string(pretty.Color(pretty.Pretty(data), nil))
	return truncateLines(result, maxPreviewLines), nil
}

func truncateLines(s string, max int) string {
	lines := strings.SplitN(s, "\n", max+2)
	if len(lines) > max {
		lines = lines[:max]
	}
	return strings.Join(lines, "\n")
}
