//go:build cgo

package preview

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/marcboeker/go-duckdb/v2"
)

// DuckDBFormatter executes a SQL query against a parquet/json/csv file.
type DuckDBFormatter struct {
	Query string // SQL with 'data' as table placeholder
	Width int
}

// Format runs the query against the file at path and returns an ASCII table.
func (f DuckDBFormatter) Format(path string) (string, error) {
	query := replaceDataSource(f.Query, path)
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return "", fmt.Errorf("open duckdb: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	return renderQueryResult(rows, f.Width)
}

func replaceDataSource(query, path string) string {
	lower := strings.ToLower(query)
	reader := pickReader(path)
	src := reader + "('" + path + "')"
	if idx := strings.Index(lower, "from data"); idx >= 0 {
		return query[:idx+5] + " " + src + query[idx+9:]
	}
	return query
}

func pickReader(path string) string {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".json"):
		return "read_json_auto"
	case strings.HasSuffix(lower, ".csv"):
		return "read_csv_auto"
	default:
		return "read_parquet"
	}
}

func renderQueryResult(rows *sql.Rows, width int) (string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("columns: %w", err)
	}

	cells, err := scanRows(rows, len(cols))
	if err != nil {
		return "", err
	}

	colWidths := computeColWidths(cols, cells, width)
	return buildGenericTable(cols, cells, colWidths), nil
}

func scanRows(rows *sql.Rows, colCount int) ([][]string, error) {
	const maxRows = 1000
	var cells [][]string
	for rows.Next() && len(cells) < maxRows {
		vals := make([]interface{}, colCount)
		ptrs := make([]interface{}, colCount)
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		row := make([]string, colCount)
		for i, v := range vals {
			row[i] = fmt.Sprintf("%v", v)
		}
		cells = append(cells, row)
	}
	return cells, rows.Err()
}

func buildGenericTable(cols []string, cells [][]string, colWidths []int) string {
	var sb strings.Builder
	sb.WriteString(formatRow(cols, colWidths) + "\n")
	sb.WriteString(strings.Repeat("-", totalWidth(colWidths)) + "\n")
	for _, row := range cells {
		sb.WriteString(formatRow(row, colWidths) + "\n")
	}
	return sb.String()
}
