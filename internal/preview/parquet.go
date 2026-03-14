package preview

import (
	"fmt"
	"io"
	"os"
	"strings"

	parquet "github.com/parquet-go/parquet-go"
)

// ParquetFormatter formats a Parquet file as an ASCII table.
type ParquetFormatter struct {
	MaxRows int
	Width   int // available display width; 0 = no wrapping
}

// Format opens the parquet file and renders first MaxRows rows as a table.
func (f ParquetFormatter) Format(path string) (string, error) {
	maxRows := f.MaxRows
	if maxRows <= 0 {
		maxRows = 20
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file: %w", err)
	}

	pf, err := parquet.OpenFile(file, stat.Size())
	if err != nil {
		return "", fmt.Errorf("open parquet: %w", err)
	}

	cols := leafColumnNames(pf)
	rows, err := readFirstRows(pf, maxRows)
	if err != nil {
		return "", err
	}
	return buildTable(cols, rows, f.Width), nil
}

func leafColumnNames(pf *parquet.File) []string {
	schema := pf.Schema()
	if schema == nil {
		return nil
	}
	var names []string
	for _, path := range schema.Columns() {
		names = append(names, strings.Join(path, "."))
	}
	return names
}

func readFirstRows(pf *parquet.File, max int) ([]parquet.Row, error) {
	rgs := pf.RowGroups()
	if len(rgs) == 0 {
		return nil, nil
	}
	reader := parquet.NewRowGroupRowReader(rgs[0])
	defer reader.Close()

	var all []parquet.Row
	buf := make([]parquet.Row, max)
	n, err := reader.ReadRows(buf)
	for i := 0; i < n; i++ {
		all = append(all, buf[i])
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read rows: %w", err)
	}
	return all, nil
}

func buildTable(cols []string, rows []parquet.Row, width int) string {
	// Collect all cell values as strings.
	cells := make([][]string, len(rows))
	for i, row := range rows {
		cells[i] = make([]string, len(row))
		for j, v := range row {
			cells[i][j] = v.String()
		}
	}

	colWidths := computeColWidths(cols, cells, width)

	var sb strings.Builder
	sb.WriteString(formatRow(cols, colWidths) + "\n")
	sb.WriteString(strings.Repeat("-", totalWidth(colWidths)) + "\n")
	for _, row := range cells {
		sb.WriteString(formatRow(row, colWidths) + "\n")
	}
	return sb.String()
}

func computeColWidths(cols []string, cells [][]string, availWidth int) []int {
	n := len(cols)
	widths := make([]int, n)
	for i, c := range cols {
		widths[i] = len(c)
	}
	for _, row := range cells {
		for i, v := range row {
			if i < n && len(v) > widths[i] {
				widths[i] = len(v)
			}
		}
	}
	if availWidth <= 0 {
		return widths
	}
	// Cap each column so total fits in availWidth (leaving 3 chars per col separator).
	budget := availWidth - (n * 3)
	if budget < n {
		budget = n
	}
	return capWidths(widths, budget)
}

func capWidths(widths []int, budget int) []int {
	result := make([]int, len(widths))
	copy(result, widths)
	for {
		total := 0
		for _, w := range result {
			total += w
		}
		if total <= budget {
			break
		}
		maxIdx, maxVal := 0, 0
		for i, w := range result {
			if w > maxVal {
				maxVal = w
				maxIdx = i
			}
		}
		if maxVal <= 1 {
			break
		}
		result[maxIdx]--
	}
	return result
}

func formatRow(cells []string, widths []int) string {
	parts := make([]string, len(widths))
	for i, w := range widths {
		v := ""
		if i < len(cells) {
			v = cells[i]
		}
		if len(v) > w {
			if w > 1 {
				v = v[:w-1] + "…"
			} else {
				v = "…"
			}
		}
		parts[i] = fmt.Sprintf("%-*s", w, v)
	}
	return "| " + strings.Join(parts, " | ") + " |"
}

func totalWidth(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w + 3
	}
	return total + 1
}
