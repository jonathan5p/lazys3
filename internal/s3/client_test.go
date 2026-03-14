package s3_test

import (
	"context"
	"testing"
	"time"

	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

func TestObjectDisplayNameFile(t *testing.T) {
	obj := s3pkg.Object{Key: "folder/subfolder/file.txt"}
	if got := obj.DisplayName(); got != "file.txt" {
		t.Errorf("expected file.txt, got %q", got)
	}
}

func TestObjectDisplayNameDir(t *testing.T) {
	obj := s3pkg.Object{Key: "folder/subfolder/", IsDir: true}
	if got := obj.DisplayName(); got != "subfolder/" {
		t.Errorf("expected subfolder/, got %q", got)
	}
}

func TestObjectDisplayNameTopLevel(t *testing.T) {
	obj := s3pkg.Object{Key: "file.csv"}
	if got := obj.DisplayName(); got != "file.csv" {
		t.Errorf("expected file.csv, got %q", got)
	}
}

type fakeS3API struct {
	buckets []string
	objects []s3pkg.Object
}

func (f *fakeS3API) ListBuckets(ctx context.Context) ([]string, error) {
	return f.buckets, nil
}

func (f *fakeS3API) ListObjects(ctx context.Context, bucket, prefix string) ([]s3pkg.Object, error) {
	return f.objects, nil
}

func TestListBucketsFromFake(t *testing.T) {
	fake := &fakeS3API{buckets: []string{"alpha", "beta"}}
	buckets, err := fake.ListBuckets(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(buckets))
	}
}

func TestListObjectsFromFake(t *testing.T) {
	fake := &fakeS3API{
		objects: []s3pkg.Object{
			{Key: "prefix/file.txt", Size: 100, LastModified: time.Now()},
			{Key: "prefix/sub/", IsDir: true},
		},
	}
	objects, err := fake.ListObjects(context.Background(), "bucket", "prefix/")
	if err != nil {
		t.Fatal(err)
	}
	if len(objects) != 2 {
		t.Errorf("expected 2 objects, got %d", len(objects))
	}
	if objects[1].DisplayName() != "sub/" {
		t.Errorf("expected sub/, got %q", objects[1].DisplayName())
	}
}
