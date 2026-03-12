package s3

import "testing"

func TestObjectFormat(t *testing.T) {
	cases := []struct {
		name     string
		input    Object
		expected string
	}{
		{
			"nested file",
			Object{Key: "photos/cat.jpg", Size: 1024, IsDir: false},
			"cat.jpg",
		},
		{
			"root level file",
			Object{Key: "readme.txt", IsDir: false},
			"readme.txt",
		},
		{
			"directory",
			Object{Key: "photos/", IsDir: true},
			"photos/",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Format()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
