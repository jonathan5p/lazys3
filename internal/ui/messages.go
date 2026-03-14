package ui

import s3pkg "github.com/jonathan5p/lazys3/internal/s3"

// BucketsLoadedMsg carries the result of a ListBuckets call.
type BucketsLoadedMsg struct {
	Buckets []string
	Err     error
}

// ObjectsLoadedMsg carries the result of a ListObjects call.
type ObjectsLoadedMsg struct {
	Objects []s3pkg.Object
	Err     error
}

// PreviewReadyMsg carries the content of a fetched object for preview.
type PreviewReadyMsg struct {
	Content string
	RawPath string // temp file path kept for rescale (caller must clean up)
	FileKey string // S3 key used for file type detection
	Err     error
}

// MetadataLoadedMsg carries the formatted metadata for an object.
type MetadataLoadedMsg struct {
	Content string
	Err     error
}

// DownloadDoneMsg reports the result of a completed download.
type DownloadDoneMsg struct {
	Path string
	Err  error
}
