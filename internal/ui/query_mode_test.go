package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestModeQueryConstantExists(t *testing.T) {
	var m AppMode = ModeQuery
	if m != ModeQuery {
		t.Error("ModeQuery constant should exist")
	}
}

func TestHandleColonKeyEntersModeQuery(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageObjects
	app.activePane = PanePreview

	updated, _ := app.Update(tea.KeyPressMsg{Code: ':', Text: ":"})
	a := updated.(AppModel)
	if a.mode != ModeQuery {
		t.Errorf("expected ModeQuery after :, got %d", a.mode)
	}
}

func TestHandleColonKeyPreFillsQuery(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.stage = StageObjects
	app.activePane = PanePreview

	updated, _ := app.Update(tea.KeyPressMsg{Code: ':', Text: ":"})
	a := updated.(AppModel)
	if a.queryState.text != "SELECT * FROM data LIMIT 20" {
		t.Errorf("expected pre-filled query, got %q", a.queryState.text)
	}
}

func TestHandleQueryKeyEscReturnsModeNormal(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeQuery
	app.queryState.text = "SELECT 1"

	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	a := updated.(AppModel)
	if a.mode != ModeNormal {
		t.Errorf("expected ModeNormal after esc, got %d", a.mode)
	}
}

func TestHandleQueryKeyBackspaceRemovesChar(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeQuery
	app.queryState.text = "SELECT"
	app.queryState.cursor = len("SELECT")

	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	a := updated.(AppModel)
	if a.queryState.text != "SELEC" {
		t.Errorf("expected SELEC after backspace, got %q", a.queryState.text)
	}
}

func TestHandleQueryKeyPrintableAppendsChar(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeQuery
	app.queryState.text = "SELECT"
	app.queryState.cursor = len("SELECT")

	updated, _ := app.Update(tea.KeyPressMsg{Code: '1', Text: "1"})
	a := updated.(AppModel)
	if a.queryState.text != "SELECT1" {
		t.Errorf("expected SELECT1 after typing 1, got %q", a.queryState.text)
	}
}

func TestHandleQueryKeyEnterReturnsNonNilCmdWhenRawPathSet(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeQuery
	app.queryState.text = "SELECT * FROM data LIMIT 20"
	app.preview.rawPath = "/tmp/test.parquet"
	app.preview.fileKey = "test.parquet"

	_, cmd := app.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected non-nil cmd when Enter pressed with rawPath set")
	}
}

func TestModeQueryLayoutContainsSQLPrompt(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeQuery
	app.queryState.text = "SELECT * FROM data"
	app.stage = StageObjects

	view := renderLayout(app)
	if !strings.Contains(view, "SQL>") {
		t.Errorf("expected SQL> in layout when mode is ModeQuery, got: %q", view)
	}
}

func TestPasteMsgInModeQueryAppendsToQueryText(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeQuery
	app.queryState.text = "SELECT "
	app.queryState.cursor = len("SELECT ")

	updated, _ := app.Update(tea.PasteMsg{Content: "name"})
	a := updated.(AppModel)
	if a.queryState.text != "SELECT name" {
		t.Errorf("expected 'SELECT name' after paste, got %q", a.queryState.text)
	}
}
