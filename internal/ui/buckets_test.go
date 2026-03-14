package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestBucketPaneSelectedBucketEmpty(t *testing.T) {
	p := newBucketPane(80, 24)
	if got := p.SelectedBucket(); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestBucketPaneSelectedBucketAfterSetItems(t *testing.T) {
	p := newBucketPane(80, 24)
	p.SetItems([]string{"alpha", "beta", "gamma"})
	if got := p.SelectedBucket(); got != "alpha" {
		t.Errorf("expected alpha, got %q", got)
	}
}

func TestBucketPaneViewNotEmpty(t *testing.T) {
	p := newBucketPane(80, 24)
	p.SetItems([]string{"alpha"})
	v := p.View(true)
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestBucketPaneUpdateJMovesDown(t *testing.T) {
	p := newBucketPane(80, 24)
	p.SetItems([]string{"alpha", "beta"})
	msg := tea.KeyPressMsg{Code: 'j'}
	updated, _ := p.Update(msg)
	if updated.SelectedBucket() != "beta" {
		t.Errorf("expected beta after j, got %q", updated.SelectedBucket())
	}
}
