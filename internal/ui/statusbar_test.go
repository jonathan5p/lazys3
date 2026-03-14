package ui

import (
	"errors"
	"strings"
	"testing"
)

func TestStatusBarSetErrorFormatsPrefix(t *testing.T) {
	s := newStatusBar(80)
	s = s.SetError(errors.New("connection refused"))
	if !strings.Contains(s.message, "Error") {
		t.Errorf("expected error prefix, got %q", s.message)
	}
}

func TestStatusBarViewFullWidth(t *testing.T) {
	s := newStatusBar(80)
	s = s.SetMessage("ready")
	v := s.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestStatusBarSetMessageStored(t *testing.T) {
	s := newStatusBar(80)
	s = s.SetMessage("loading buckets")
	if s.message != "loading buckets" {
		t.Errorf("expected 'loading buckets', got %q", s.message)
	}
}

func TestStatusBarSetModeFilterViewContainsFilter(t *testing.T) {
	a := AppModel{mode: ModeFilter, filter: FilterState{text: "foo"}, width: 80}
	filterView := renderFilterBar(a)
	if !strings.Contains(filterView, "foo") {
		t.Errorf("expected filter bar to contain 'foo', got: %q", filterView)
	}
}

func TestStatusBarSetModeFilterViewContainsText(t *testing.T) {
	a := AppModel{mode: ModeFilter, filter: FilterState{text: "bar"}, width: 80}
	filterView := renderFilterBar(a)
	if !strings.Contains(filterView, "bar") {
		t.Errorf("expected filter bar to contain 'bar', got: %q", filterView)
	}
}

func TestStatusBarSetModeNormalViewNoFilter(t *testing.T) {
	s := newStatusBar(80)
	s = s.SetMode(ModeNormal, "")
	if strings.Contains(s.View(), "FILTER") {
		t.Errorf("expected view NOT to contain 'FILTER', got: %q", s.View())
	}
}
