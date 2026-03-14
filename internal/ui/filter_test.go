package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

func TestFilterStateApplyEmptyReturnsAll(t *testing.T) {
	fs := FilterState{}
	items := []string{"alpha", "beta", "gamma"}
	got := fs.Apply(items)
	if len(got) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(got))
	}
}

func TestFilterStateApplyMatchesCaseInsensitive(t *testing.T) {
	fs := FilterState{text: "foo"}
	items := []string{"foobar", "baz", "FOOnk"}
	got := fs.Apply(items)
	if len(got) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(got), got)
	}
}

func TestFilterModeEscClearsFilterAndReturnsNormal(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeFilter
	app.filter = FilterState{text: "hello", active: true}

	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	a := updated.(AppModel)

	if a.mode != ModeNormal {
		t.Errorf("expected ModeNormal after esc, got %v", a.mode)
	}
	if a.filter.text != "" {
		t.Errorf("expected empty filter text after esc, got %q", a.filter.text)
	}
}

func TestFilterEnterWithSlashReturnsCmd(t *testing.T) {
	app := NewApp(&mockS3Client{}, 80, 24)
	app.buckets.SetItems([]string{"bucket1"})
	app.objects.SetItems([]s3pkg.Object{{Key: "folder/", IsDir: true}})
	app.mode = ModeFilter
	app.activePane = PaneObjects
	app.filter = FilterState{text: "folder/sub", active: true}

	_, cmd := app.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected non-nil Cmd when filter text contains '/' and pane is objects")
	}
}

func TestFilterEnterWithoutSlashReturnsNilCmd(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeFilter
	app.activePane = PaneObjects
	app.filter = FilterState{text: "noslash", active: true}

	updated, cmd := app.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	a := updated.(AppModel)
	if cmd != nil {
		t.Error("expected nil Cmd when filter text has no '/'")
	}
	if a.mode != ModeNormal {
		t.Errorf("expected ModeNormal, got %v", a.mode)
	}
}

func TestParseS3URL(t *testing.T) {
	cases := []struct {
		input      string
		wantBucket string
		wantKey    string
		wantPrefix string
		wantNil    bool
	}{
		{
			input:      "s3://validdata-clarity-mgmt-pro/jobs/data/file.parquet",
			wantBucket: "validdata-clarity-mgmt-pro",
			wantKey:    "jobs/data/file.parquet",
			wantPrefix: "jobs/data/",
		},
		{
			input:      "s3://my-bucket/deep/path/obj.json",
			wantBucket: "my-bucket",
			wantKey:    "deep/path/obj.json",
			wantPrefix: "deep/path/",
		},
		{
			input:   "not-an-s3-url",
			wantNil: true,
		},
		{
			input:      "s3://bucket-only",
			wantBucket: "bucket-only",
			wantKey:    "",
			wantPrefix: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := parseS3URL(tc.input)
			if tc.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected non-nil result")
			}
			if got.bucket != tc.wantBucket {
				t.Errorf("bucket: got %q want %q", got.bucket, tc.wantBucket)
			}
			if got.key != tc.wantKey {
				t.Errorf("key: got %q want %q", got.key, tc.wantKey)
			}
			if got.prefix != tc.wantPrefix {
				t.Errorf("prefix: got %q want %q", got.prefix, tc.wantPrefix)
			}
		})
	}
}
