package s3

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func writeToFile(rc io.ReadCloser, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	return err
}
