package s3

import (
	"strings"
	"time"
)

// Object represents an S3 object or virtual directory prefix.
type Object struct {
	Key          string
	Size         int64
	LastModified time.Time
	IsDir        bool
	ContentType  string
}

// DisplayName returns the last path segment; directories get a trailing slash.
func (o Object) DisplayName() string {
	key := strings.TrimSuffix(o.Key, "/")
	parts := strings.Split(key, "/")
	name := parts[len(parts)-1]
	if o.IsDir {
		return name + "/"
	}
	return name
}
