package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

func TestObjectPaneSelectedObjectEmpty(t *testing.T) {
	p := newObjectPane(80, 24)
	if got := p.SelectedObject(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestObjectPaneSelectedObjectAfterSetItems(t *testing.T) {
	p := newObjectPane(80, 24)
	objs := []s3pkg.Object{
		{Key: "prefix/file1.txt"},
		{Key: "prefix/file2.txt"},
	}
	p.SetItems(objs)
	got := p.SelectedObject()
	if got == nil {
		t.Fatal("expected non-nil object")
	}
	if got.Key != "prefix/file1.txt" {
		t.Errorf("expected prefix/file1.txt, got %q", got.Key)
	}
}

func TestObjectPanePrefixDefault(t *testing.T) {
	p := newObjectPane(80, 24)
	if p.Prefix() != "" {
		t.Errorf("expected empty prefix, got %q", p.Prefix())
	}
}

func TestObjectPaneUpdateJMovesDown(t *testing.T) {
	p := newObjectPane(80, 24)
	p.SetItems([]s3pkg.Object{{Key: "a/b.txt"}, {Key: "a/c.txt"}})
	msg := tea.KeyPressMsg{Code: 'j'}
	updated, _ := p.Update(msg)
	obj := updated.SelectedObject()
	if obj == nil || obj.Key != "a/c.txt" {
		t.Errorf("expected a/c.txt after j, got %v", obj)
	}
}
