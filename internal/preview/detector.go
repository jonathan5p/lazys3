package preview

import (
	"bytes"
	"strings"
)

// FileType identifies the format of a file for preview purposes.
type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeJSON
	FileTypeParquet
	FileTypeBinary
	FileTypeText
)

// Detect returns the FileType by checking extension first, then magic bytes.
func Detect(key string, header []byte) FileType {
	if ft := detectByExtension(key); ft != FileTypeUnknown {
		return ft
	}
	return detectByMagic(header)
}

func detectByExtension(key string) FileType {
	lower := strings.ToLower(key)
	switch {
	case strings.HasSuffix(lower, ".json"):
		return FileTypeJSON
	case strings.HasSuffix(lower, ".parquet"):
		return FileTypeParquet
	}
	return FileTypeUnknown
}

func detectByMagic(header []byte) FileType {
	if len(header) == 0 {
		return FileTypeUnknown
	}
	if bytes.HasPrefix(header, []byte("PAR1")) {
		return FileTypeParquet
	}
	if header[0] == '{' || header[0] == '[' {
		return FileTypeJSON
	}
	return FileTypeBinary
}
