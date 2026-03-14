package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/jonathan5p/lazys3/internal/preview"
)

// PreviewPane renders the content of the selected S3 object.
type PreviewPane struct {
	content  string
	scrollY  int
	loading  bool
	err      error
	width    int
	height   int
	colScale int    // percentage of width to use for columns; default 100
	rawPath  string // temp file path from last preview (for rescale)
	fileKey  string // S3 key for file type detection
}

func newPreviewPane(width, height int) PreviewPane {
	return PreviewPane{width: width, height: height, colScale: 100}
}

// SetContent stores new content and resets the scroll position.
func (p PreviewPane) SetContent(content string) PreviewPane {
	p.content = content
	p.scrollY = 0
	return p
}

// Update handles scroll key events.
func (p PreviewPane) Update(msg tea.Msg) (PreviewPane, tea.Cmd) {
	lines := strings.Split(p.content, "\n")
	maxScroll := len(lines) - (p.height - 4)
	if maxScroll < 0 {
		maxScroll = 0
	}
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch m.String() {
		case "j", "down":
			p.scrollY = clampInt(p.scrollY+1, 0, maxScroll)
		case "k", "up":
			p.scrollY = clampInt(p.scrollY-1, 0, maxScroll)
		case "+", "=":
			p = p.increaseScale()
			return p, p.rescaleCmd()
		case "-":
			p = p.decreaseScale()
			return p, p.rescaleCmd()
		}
	}
	return p, nil
}

func (p PreviewPane) increaseScale() PreviewPane {
	p.colScale = clampInt(p.colScale+10, 10, 300)
	return p
}

func (p PreviewPane) decreaseScale() PreviewPane {
	p.colScale = clampInt(p.colScale-10, 10, 300)
	return p
}

func (p PreviewPane) rescaleCmd() tea.Cmd {
	if p.rawPath == "" {
		return nil
	}
	scaledWidth := p.width * p.colScale / 100
	path := p.rawPath
	key := p.fileKey
	return func() tea.Msg {
		ft := preview.Detect(key, nil)
		msg := formatForType(ft, path, scaledWidth)
		msg.RawPath = path
		msg.FileKey = key
		return msg
	}
}

// View renders the preview pane with the given active state.
func (p PreviewPane) View(active bool) string {
	style := paneStyle(active).Width(p.width - 2).Height(p.height - 2)

	if p.loading {
		return style.Render("Loading...")
	}
	if p.err != nil {
		return style.Render(fmt.Sprintf("Error: %v", p.err))
	}

	header := titleStyle.Render("Preview") + "\n\n"
	lines := strings.Split(p.content, "\n")
	end := p.scrollY + (p.height - 4)
	if end > len(lines) {
		end = len(lines)
	}
	visible := strings.Join(lines[p.scrollY:end], "\n")
	return style.Render(header + visible)
}
