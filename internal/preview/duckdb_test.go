//go:build cgo

package preview

import (
	"os"
	"strings"
	"testing"

	parquet "github.com/parquet-go/parquet-go"
)

func TestDuckDBFormatterReturnsNonEmptyResult(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.parquet")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	rows := []testRow{
		{Name: "alice", Age: 30},
		{Name: "bob", Age: 25},
	}
	if err := parquet.WriteFile(f.Name(), rows); err != nil {
		t.Fatalf("write parquet: %v", err)
	}

	d := DuckDBFormatter{Query: "SELECT * FROM data LIMIT 20", Width: 80}
	out, err := d.Format(f.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output from DuckDB formatter")
	}
}

func TestDuckDBFormatterOutputContainsColumnNames(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.parquet")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	rows := []testRow{
		{Name: "charlie", Age: 40},
	}
	if err := parquet.WriteFile(f.Name(), rows); err != nil {
		t.Fatalf("write parquet: %v", err)
	}

	d := DuckDBFormatter{Query: "SELECT * FROM data LIMIT 20", Width: 80}
	out, err := d.Format(f.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out, "name") && !strings.Contains(out, "age") {
		t.Errorf("expected output to contain column names, got:\n%s", out)
	}
}
