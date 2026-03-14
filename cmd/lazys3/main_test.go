package main

import (
	"testing"
)

func TestRootCommandExists(t *testing.T) {
	cmd := newRootCmd()
	if cmd == nil {
		t.Fatal("expected non-nil root command")
	}
	if cmd.Use != "lazys3" {
		t.Errorf("expected Use=lazys3, got %q", cmd.Use)
	}
}
