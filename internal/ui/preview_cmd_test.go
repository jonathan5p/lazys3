package ui

import (
	"context"
	"io"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

type mockS3Client struct{}

func (m *mockS3Client) ListBuckets(_ context.Context) ([]string, error) {
	return []string{"bucket1"}, nil
}

func (m *mockS3Client) ListObjects(_ context.Context, _, _ string) ([]s3pkg.Object, error) {
	return []s3pkg.Object{{Key: "file.json"}}, nil
}

func (m *mockS3Client) HeadObject(_ context.Context, _, _ string) (map[string]string, error) {
	return map[string]string{"ContentType": "application/json"}, nil
}

func (m *mockS3Client) GetObject(_ context.Context, _, _ string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(`{"test":true}`)), nil
}

func (m *mockS3Client) DownloadObject(_ context.Context, _, _, _ string) error {
	return nil
}

func TestPressPOnNonDirObjectReturnsCmdNotNil(t *testing.T) {
	app := NewApp(&mockS3Client{}, 120, 40)
	app.buckets.SetItems([]string{"bucket1"})
	app.objects.SetItems([]s3pkg.Object{{Key: "file.json", IsDir: false}})
	app.activePane = PaneObjects

	_, cmd := app.Update(tea.KeyPressMsg{Code: 'p'})
	if cmd == nil {
		t.Fatal("expected non-nil Cmd when pressing p on a non-dir object")
	}
}
