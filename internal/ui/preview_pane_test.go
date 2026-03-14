package ui

import (
	"strings"
	"testing"
)

func TestPreviewPaneSetContentStored(t *testing.T) {
	p := newPreviewPane(80, 24)
	p2 := p.SetContent("hello world")
	if p2.content != "hello world" {
		t.Errorf("expected content stored, got %q", p2.content)
	}
}

func TestPreviewPaneViewContainsBorder(t *testing.T) {
	p := newPreviewPane(80, 24)
	p = p.SetContent("test content")
	v := p.View(false)
	if !strings.Contains(v, "test content") {
		t.Errorf("expected view to contain content, got: %q", v)
	}
}

func TestPreviewPaneViewEmpty(t *testing.T) {
	p := newPreviewPane(80, 24)
	v := p.View(false)
	if v == "" {
		t.Error("expected non-empty view even when content is empty")
	}
}
