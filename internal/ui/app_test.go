package ui

import (
	"strings"
	"testing"

	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

func TestNewAppDefaultMode(t *testing.T) {
	app := NewApp(nil, 0, 0)
	if app.mode != ModeNormal {
		t.Errorf("expected default mode ModeNormal (%d), got %d", ModeNormal, app.mode)
	}
}

func TestAppUpdateBucketsLoadedMsg(t *testing.T) {
	app := NewApp(nil, 80, 24)
	msg := BucketsLoadedMsg{Buckets: []string{"a", "b"}}
	updated, _ := app.Update(msg)
	a := updated.(AppModel)
	if a.buckets.SelectedBucket() != "a" {
		t.Errorf("expected first bucket 'a', got %q", a.buckets.SelectedBucket())
	}
}

func TestAppUpdateBucketsLoadedMsgCount(t *testing.T) {
	app := NewApp(nil, 80, 24)
	msg := BucketsLoadedMsg{Buckets: []string{"x", "y", "z"}}
	updated, _ := app.Update(msg)
	a := updated.(AppModel)
	if len(a.buckets.items) != 3 {
		t.Errorf("expected 3 buckets, got %d", len(a.buckets.items))
	}
}

func TestAppUpdateObjectsLoadedMsg(t *testing.T) {
	app := NewApp(nil, 80, 24)
	objs := []s3pkg.Object{{Key: "buck/file.txt"}}
	msg := ObjectsLoadedMsg{Objects: objs}
	updated, _ := app.Update(msg)
	a := updated.(AppModel)
	if a.objects.SelectedObject() == nil {
		t.Error("expected non-nil selected object")
	}
}

func TestNewAppDefaultStageBuckets(t *testing.T) {
	app := NewApp(nil, 0, 0)
	if app.stage != StageBuckets {
		t.Errorf("expected default stage StageBuckets (%d), got %d", StageBuckets, app.stage)
	}
}

func TestRenderLayoutStageBucketsContainsBuckets(t *testing.T) {
	app := NewApp(nil, 120, 30)
	app.stage = StageBuckets
	app.buckets.SetItems([]string{"my-bucket"})
	out := renderLayout(app)
	if !strings.Contains(out, "Buckets") {
		t.Error("expected 'Buckets' in StageBuckets layout")
	}
}

func TestRenderLayoutStageObjectsHidesBuckets(t *testing.T) {
	app := NewApp(nil, 120, 30)
	app.stage = StageObjects
	app.buckets.SetItems([]string{"my-bucket"})
	app.objects.SetItems([]s3pkg.Object{{Key: "file.txt"}})
	out := renderLayout(app)
	if strings.Contains(out, "Buckets") {
		t.Error("expected 'Buckets' pane to be hidden in StageObjects layout")
	}
}

func TestRenderLayoutStageBucketsObjectsPaneTitleShowsBucket(t *testing.T) {
	app := NewApp(nil, 120, 30)
	app.stage = StageBuckets
	app.buckets.SetItems([]string{"my-bucket"})
	out := renderLayout(app)
	if !strings.Contains(out, "Objects: my-bucket") {
		t.Error("expected objects pane title to show 'Objects: my-bucket' in StageBuckets")
	}
}

func TestHandleObjectsLoadedInStageBucketsResetsObjectsCursorOnly(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageBuckets
	app.objects.cursor = 5
	msg := ObjectsLoadedMsg{Objects: []s3pkg.Object{{Key: "a.txt"}, {Key: "b.txt"}}}
	updated, _ := app.Update(msg)
	a := updated.(AppModel)
	if a.stage != StageBuckets {
		t.Errorf("expected stage to stay StageBuckets on preview load, got %d", a.stage)
	}
	if a.objects.cursor != 0 {
		t.Errorf("expected cursor reset to 0 on preview load, got %d", a.objects.cursor)
	}
}

func TestBucketsLoadedFiresPreviewForFirstBucket(t *testing.T) {
	app := NewApp(&mockS3Client{}, 80, 24)
	msg := BucketsLoadedMsg{Buckets: []string{"first-bucket", "second-bucket"}}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected non-nil cmd after BucketsLoadedMsg (preview of first bucket)")
	}
}
