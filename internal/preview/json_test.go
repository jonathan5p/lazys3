package preview

import (
	"os"
	"strings"
	"testing"
)

func TestJSONFormatterValidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{"name":"alice","age":30}`)
	f.Close()

	out, err := JSONFormatter{}.Format(f.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out, "{") {
		t.Errorf("expected output to contain '{', got %q", out)
	}
}

func TestJSONFormatterInvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`not json at all`)
	f.Close()

	_, err = JSONFormatter{}.Format(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
