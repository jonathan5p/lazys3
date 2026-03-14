package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestHandleKeyQuitReturnsQuitCmd(t *testing.T) {
	app := NewApp(nil, 80, 24)
	_, cmd := app.Update(tea.KeyPressMsg{Code: 'q'})
	if cmd == nil {
		t.Fatal("expected non-nil cmd for q key")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestHandleKeyJMovesDownInBuckets(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.buckets.SetItems([]string{"a", "b", "c"})
	app.activePane = PaneBuckets
	updated, _ := app.Update(tea.KeyPressMsg{Code: 'j'})
	a := updated.(AppModel)
	if a.buckets.cursor != 1 {
		t.Errorf("expected cursor=1 after j, got %d", a.buckets.cursor)
	}
}

func TestHandleKeyHSetsFocusToBuckets(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.activePane = PaneObjects
	updated, _ := app.Update(tea.KeyPressMsg{Code: 'h'})
	a := updated.(AppModel)
	if a.activePane != PaneBuckets {
		t.Errorf("expected PaneBuckets after h, got %d", a.activePane)
	}
}

func TestHandleEnterOnBucketSetsStageObjects(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageBuckets
	app.activePane = PaneBuckets
	app.buckets.SetItems([]string{"my-bucket"})
	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	a := updated.(AppModel)
	if a.stage != StageObjects {
		t.Errorf("expected StageObjects after Enter on bucket, got %d", a.stage)
	}
	if a.activePane != PaneObjects {
		t.Errorf("expected PaneObjects after Enter on bucket, got %d", a.activePane)
	}
}

func TestHandleEscOnObjectsAtRootReturnsToStageBuckets(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageObjects
	app.activePane = PaneObjects
	app.objects.prefix = ""
	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	a := updated.(AppModel)
	if a.stage != StageBuckets {
		t.Errorf("expected StageBuckets after Esc at root, got %d", a.stage)
	}
	if a.activePane != PaneBuckets {
		t.Errorf("expected PaneBuckets after Esc at root, got %d", a.activePane)
	}
}

func TestHandleHOnObjectsAtRootReturnsToStageBuckets(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageObjects
	app.activePane = PaneObjects
	app.objects.prefix = ""
	updated, _ := app.Update(tea.KeyPressMsg{Code: 'h'})
	a := updated.(AppModel)
	if a.stage != StageBuckets {
		t.Errorf("expected StageBuckets after h at root, got %d", a.stage)
	}
}

func TestHandleEscOnObjectsWithPrefixNavigatesUp(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageObjects
	app.activePane = PaneObjects
	app.buckets.SetItems([]string{"my-bucket"})
	app.objects.prefix = "folder/"
	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	a := updated.(AppModel)
	if a.stage != StageObjects {
		t.Errorf("expected StageObjects to remain when esc with prefix, got %d", a.stage)
	}
	if a.objects.prefix != "" {
		t.Errorf("expected prefix to go to root, got %q", a.objects.prefix)
	}
}

func TestHandleCursorDownInStageBucketsFiresPreviewCmd(t *testing.T) {
	app := NewApp(&mockS3Client{}, 80, 24)
	app.stage = StageBuckets
	app.activePane = PaneBuckets
	app.buckets.SetItems([]string{"bucket-a", "bucket-b"})
	_, cmd := app.Update(tea.KeyPressMsg{Code: 'j'})
	if cmd == nil {
		t.Error("expected a non-nil cmd after j in StageBuckets (preview load)")
	}
}

func TestHandleCursorUpInStageBucketsFiresPreviewCmd(t *testing.T) {
	app := NewApp(&mockS3Client{}, 80, 24)
	app.stage = StageBuckets
	app.activePane = PaneBuckets
	app.buckets.SetItems([]string{"bucket-a", "bucket-b"})
	app.buckets.cursor = 1
	_, cmd := app.Update(tea.KeyPressMsg{Code: 'k'})
	if cmd == nil {
		t.Error("expected a non-nil cmd after k in StageBuckets (preview load)")
	}
}

func TestHandleCursorDownInStageObjectsDoesNotFirePreviewCmd(t *testing.T) {
	app := NewApp(&mockS3Client{}, 80, 24)
	app.stage = StageObjects
	app.activePane = PaneBuckets
	app.buckets.SetItems([]string{"bucket-a", "bucket-b"})
	_, cmd := app.Update(tea.KeyPressMsg{Code: 'j'})
	if cmd != nil {
		t.Error("expected nil cmd after j in StageObjects on buckets pane")
	}
}

func TestHandleWindowSizeUpdatesWidth(t *testing.T) {
	app := NewApp(nil, 80, 24)
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 160, Height: 40})
	a := updated.(AppModel)
	if a.width != 160 {
		t.Errorf("expected width=160, got %d", a.width)
	}
	if a.height != 40 {
		t.Errorf("expected height=40, got %d", a.height)
	}
}

func TestHandleKeyGGJumpsToTop(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.buckets.SetItems([]string{"a", "b", "c"})
	app.buckets.cursor = 2
	app.activePane = PaneBuckets

	first, _ := app.Update(tea.KeyPressMsg{Code: 'g'})
	a1 := first.(AppModel)
	if a1.prevKey != "g" {
		t.Errorf("expected prevKey=g, got %q", a1.prevKey)
	}

	second, _ := a1.Update(tea.KeyPressMsg{Code: 'g'})
	a2 := second.(AppModel)
	if a2.buckets.cursor != 0 {
		t.Errorf("expected cursor=0 after gg, got %d", a2.buckets.cursor)
	}
}
