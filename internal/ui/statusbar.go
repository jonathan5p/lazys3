package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// StatusBar shows the current status message at the bottom of the screen.
type StatusBar struct {
	message    string
	width      int
	mode       AppMode
	filterText string
}

func newStatusBar(width int) StatusBar {
	return StatusBar{width: width}
}

// SetMessage returns a copy with the given status message.
func (s StatusBar) SetMessage(msg string) StatusBar {
	s.message = msg
	return s
}

// SetError returns a copy with an error message formatted with "Error:" prefix.
func (s StatusBar) SetError(err error) StatusBar {
	s.message = fmt.Sprintf("Error: %v", err)
	return s
}

// SetMode returns a copy with the given mode and filter text.
func (s StatusBar) SetMode(mode AppMode, filterText string) StatusBar {
	s.mode = mode
	s.filterText = filterText
	return s
}

// View renders the status bar at full width.
func (s StatusBar) View() string {
	style := statusBar.Width(s.width)
	left := s.modePrefix() + s.message
	if s.mode == ModeNormal {
		hint := "  ? help"
		return style.Render(left + hint)
	}
	return style.Render(left)
}

func (s StatusBar) modePrefix() string {
	switch s.mode {
	case ModeDownload:
		return "[DOWNLOAD] "
	case ModeMetadata:
		return "[METADATA] "
	case ModeHelp:
		return "[HELP] press any key to close"
	default:
		return ""
	}
}

// Update handles messages for the status bar.
func (s StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	return s, nil
}
