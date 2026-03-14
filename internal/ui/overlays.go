package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// MetadataOverlay renders a centered box with object metadata.
type MetadataOverlay struct {
	content string
	width   int
	height  int
}

func newMetadataOverlay(width, height int) MetadataOverlay {
	return MetadataOverlay{width: width, height: height}
}

// SetContent returns a copy with the given content.
func (o MetadataOverlay) SetContent(content string) MetadataOverlay {
	o.content = content
	return o
}

// View renders the overlay as a centered lipgloss box.
func (o MetadataOverlay) View() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5")).
		Padding(1, 2).
		Width(o.width / 2).
		Render(lipgloss.NewStyle().Bold(true).Render("Metadata") + "\n\n" + o.content)
	return lipgloss.Place(o.width, o.height, lipgloss.Center, lipgloss.Center, box)
}

// DownloadPrompt renders a centered box for entering a download destination.
type DownloadPrompt struct {
	input  string
	width  int
	height int
}

func newDownloadPrompt(width, height int) DownloadPrompt {
	return DownloadPrompt{width: width, height: height}
}

// SetInput returns a copy with the given input text.
func (p DownloadPrompt) SetInput(s string) DownloadPrompt {
	p.input = s
	return p
}

// View renders the download prompt as a centered lipgloss box.
func (p DownloadPrompt) View() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")).
		Padding(1, 2).
		Width(p.width / 2).
		Render(lipgloss.NewStyle().Bold(true).Render("Download to:") + "\n\n" + p.input + "█")
	return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, box)
}

// HelpOverlay renders a centered keybindings reference box.
type HelpOverlay struct {
	width  int
	height int
}

func newHelpOverlay(width, height int) HelpOverlay {
	return HelpOverlay{width: width, height: height}
}

func (h HelpOverlay) View() string {
	content := helpContent()
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Padding(1, 2).
		Width(h.width / 2).
		Render(lipgloss.NewStyle().Bold(true).Render("Keybindings") + "\n\n" + content)
	return lipgloss.Place(h.width, h.height, lipgloss.Center, lipgloss.Center, box)
}

func helpContent() string {
	lines := []string{
		"Navigation  j/k=up/down  h/l=focus left/right  Enter=select  Esc/Backspace=back",
		"Buckets     r=refresh  /=filter",
		"Objects     p=preview  d=download  m=metadata  /=filter  r=refresh",
		"Preview     j/k=scroll",
		"General     G=bottom  gg=top  ?=toggle help  q=quit",
	}
	result := ""
	for i, l := range lines {
		if i > 0 {
			result += "\n"
		}
		result += l
	}
	return result
}
