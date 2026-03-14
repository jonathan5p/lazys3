package preview

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		header []byte
		want   FileType
	}{
		{"json extension", "data.json", nil, FileTypeJSON},
		{"parquet extension", "data.parquet", nil, FileTypeParquet},
		{"PAR1 magic bytes", "datafile", []byte("PAR1somedata"), FileTypeParquet},
		{"curly brace first byte", "datafile", []byte("{...}"), FileTypeJSON},
		{"square bracket first byte", "datafile", []byte("[...]"), FileTypeJSON},
		{"random bytes", "datafile", []byte("\x00\x01\x02\x03"), FileTypeBinary},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Detect(tc.key, tc.header)
			if got != tc.want {
				t.Errorf("Detect(%q, ...) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}
