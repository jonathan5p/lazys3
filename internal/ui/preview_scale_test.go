package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestPreviewPaneDefaultColScale100(t *testing.T) {
	p := newPreviewPane(80, 24)
	if p.colScale != 100 {
		t.Errorf("expected default colScale=100, got %d", p.colScale)
	}
}

func TestPreviewPanePlusKeyIncreasesColScale(t *testing.T) {
	p := newPreviewPane(80, 24)
	p, _ = p.Update(tea.KeyPressMsg{Code: '+', Text: "+"})
	if p.colScale != 110 {
		t.Errorf("expected colScale=110 after +, got %d", p.colScale)
	}
}

func TestPreviewPaneEqualKeyIncreasesColScale(t *testing.T) {
	p := newPreviewPane(80, 24)
	p, _ = p.Update(tea.KeyPressMsg{Code: '=', Text: "="})
	if p.colScale != 110 {
		t.Errorf("expected colScale=110 after =, got %d", p.colScale)
	}
}

func TestPreviewPaneMinusKeyDecreasesColScale(t *testing.T) {
	p := newPreviewPane(80, 24)
	p, _ = p.Update(tea.KeyPressMsg{Code: '-', Text: "-"})
	if p.colScale != 90 {
		t.Errorf("expected colScale=90 after -, got %d", p.colScale)
	}
}

func TestPreviewPaneColScaleMinIsTen(t *testing.T) {
	p := newPreviewPane(80, 24)
	p.colScale = 10
	p, _ = p.Update(tea.KeyPressMsg{Code: '-', Text: "-"})
	if p.colScale != 10 {
		t.Errorf("expected colScale min=10, got %d", p.colScale)
	}
}

func TestPreviewPaneColScaleMaxIs300(t *testing.T) {
	p := newPreviewPane(80, 24)
	p.colScale = 300
	p, _ = p.Update(tea.KeyPressMsg{Code: '+', Text: "+"})
	if p.colScale != 300 {
		t.Errorf("expected colScale max=300, got %d", p.colScale)
	}
}

func TestPreviewPanePlusKeyReturnsCmdWhenRawPathSet(t *testing.T) {
	p := newPreviewPane(80, 24)
	p.rawPath = "/tmp/some-file.parquet"
	p.fileKey = "data.parquet"
	_, cmd := p.Update(tea.KeyPressMsg{Code: '+', Text: "+"})
	if cmd == nil {
		t.Error("expected non-nil cmd when rawPath is set and + pressed")
	}
}

func TestPreviewPanePlusKeyReturnsNilCmdWhenNoRawPath(t *testing.T) {
	p := newPreviewPane(80, 24)
	_, cmd := p.Update(tea.KeyPressMsg{Code: '+', Text: "+"})
	if cmd != nil {
		t.Error("expected nil cmd when rawPath is empty")
	}
}
