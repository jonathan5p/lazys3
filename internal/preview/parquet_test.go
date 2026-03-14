package preview

import (
	"os"
	"strings"
	"testing"

	parquet "github.com/parquet-go/parquet-go"
)

type testRow struct {
	Name string `parquet:"name"`
	Age  int32  `parquet:"age"`
}

func TestParquetFormatterOutputContainsColumns(t *testing.T) {
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
		t.Fatalf("write parquet file: %v", err)
	}

	out, err := ParquetFormatter{MaxRows: 20}.Format(f.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out, "name") {
		t.Errorf("expected output to contain 'name', got:\n%s", out)
	}
	if !strings.Contains(out, "age") {
		t.Errorf("expected output to contain 'age', got:\n%s", out)
	}
}
