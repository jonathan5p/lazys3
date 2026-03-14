package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

func TestMetadataOverlayViewContainsTitle(t *testing.T) {
	o := newMetadataOverlay(80, 24)
	o = o.SetContent("key=value")
	v := o.View()
	if !strings.Contains(v, "Metadata") {
		t.Errorf("expected view to contain 'Metadata', got: %q", v)
	}
}

func TestMetadataModeEscReturnsNormal(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeMetadata

	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	a := updated.(AppModel)

	if a.mode != ModeNormal {
		t.Errorf("expected ModeNormal after esc, got %v", a.mode)
	}
}

func TestPressMOnNonDirObjectReturnsCmdNotNil(t *testing.T) {
	app := NewApp(&mockS3Client{}, 120, 40)
	app.buckets.SetItems([]string{"bucket1"})
	app.objects.SetItems([]s3pkg.Object{{Key: "file.json", IsDir: false}})
	app.activePane = PaneObjects

	_, cmd := app.Update(tea.KeyPressMsg{Code: 'm'})
	if cmd == nil {
		t.Fatal("expected non-nil Cmd when pressing m on a non-dir object")
	}
}

func TestDownloadPromptViewContainsTitle(t *testing.T) {
	p := newDownloadPrompt(80, 24)
	p = p.SetInput("file.json")
	v := p.View()
	if !strings.Contains(v, "Download to:") {
		t.Errorf("expected view to contain 'Download to:', got: %q", v)
	}
}

func TestDownloadModeEscReturnsNormal(t *testing.T) {
	app := NewApp(nil, 80, 24)
	app.mode = ModeDownload

	updated, _ := app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	a := updated.(AppModel)

	if a.mode != ModeNormal {
		t.Errorf("expected ModeNormal after esc, got %v", a.mode)
	}
}

func TestDownloadModeEnterReturnsCmdNotNil(t *testing.T) {
	app := NewApp(&mockS3Client{}, 120, 40)
	app.buckets.SetItems([]string{"bucket1"})
	app.objects.SetItems([]s3pkg.Object{{Key: "file.json", IsDir: false}})
	app.mode = ModeDownload
	app.downloadPrompt = app.downloadPrompt.SetInput("file.json")

	_, cmd := app.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected non-nil Cmd when pressing enter in ModeDownload with s3 client")
	}
}

func TestHelpOverlayViewContainsKeybindings(t *testing.T) {
	h := newHelpOverlay(80, 24)
	if !strings.Contains(h.View(), "Keybindings") {
		t.Error("expected HelpOverlay.View() to contain 'Keybindings'")
	}
}

func TestHelpOverlayViewContainsJK(t *testing.T) {
	h := newHelpOverlay(80, 24)
	if !strings.Contains(h.View(), "j/k") {
		t.Error("expected HelpOverlay.View() to contain 'j/k'")
	}
}

func TestHelpOverlayViewContainsQuit(t *testing.T) {
	h := newHelpOverlay(80, 24)
	if !strings.Contains(h.View(), "quit") {
		t.Error("expected HelpOverlay.View() to contain 'quit'")
	}
}
