package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

// BucketPane displays the list of S3 buckets.
type BucketPane struct {
	items         []string
	filteredItems []string
	filterText    string
	cursor        int
	loading       bool
	err           error
	width         int
	height        int
}

// selectByName sets the cursor to the bucket with the given name.
// If not found, cursor stays at 0.
func (p *BucketPane) selectByName(name string) {
	for i, b := range p.items {
		if b == name {
			p.cursor = i
			return
		}
	}
	p.cursor = 0
}

func newBucketPane(width, height int) BucketPane {
	return BucketPane{width: width, height: height}
}

// SetItems replaces the current bucket list and resets the cursor.
func (p *BucketPane) SetItems(buckets []string) {
	p.items = buckets
	p.cursor = 0
	p.applyFilter()
}

// SetFilter stores the filter text and re-filters the displayed items.
func (p *BucketPane) SetFilter(text string) {
	p.filterText = text
	p.cursor = 0
	p.applyFilter()
}

func (p *BucketPane) applyFilter() {
	fs := FilterState{text: p.filterText}
	p.filteredItems = fs.Apply(p.items)
}

func (p BucketPane) visibleItems() []string {
	if p.filteredItems != nil {
		return p.filteredItems
	}
	return p.items
}

// SelectedBucket returns the name of the currently highlighted bucket.
func (p BucketPane) SelectedBucket() string {
	items := p.visibleItems()
	if len(items) == 0 {
		return ""
	}
	return items[p.cursor]
}

func (p BucketPane) clampCursor(n int) int {
	return clampInt(n, 0, len(p.visibleItems())-1)
}

// Update handles messages for the bucket pane.
func (p BucketPane) Update(msg tea.Msg) (BucketPane, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch m.String() {
		case "j", "down":
			p.cursor = p.clampCursor(p.cursor + 1)
		case "k", "up":
			p.cursor = p.clampCursor(p.cursor - 1)
		case "G":
			p.cursor = p.clampCursor(len(p.items) - 1)
		}
	}
	return p, nil
}

// View renders the bucket pane with the given active state.
func (p BucketPane) View(active bool) string {
	style := inactivePane
	if active {
		style = activePane
	}
	style = style.Width(p.width - 2).Height(p.height - 2)

	if p.loading {
		return style.Render("Loading...")
	}
	if p.err != nil {
		return style.Render(fmt.Sprintf("Error: %v", p.err))
	}

	content := titleStyle.Render("Buckets") + "\n\n"
	items := p.visibleItems()
	visible := visibleRange(p.cursor, p.height-4, len(items))
	for i := visible.start; i < visible.end; i++ {
		line := " " + items[i]
		if i == p.cursor {
			line = selectedItem.Render("> " + items[i])
		}
		content += line + "\n"
	}
	return style.Render(content)
}

type visRange struct{ start, end int }

func visibleRange(cursor, maxVisible, total int) visRange {
	if total == 0 {
		return visRange{}
	}
	start := cursor - maxVisible/2
	if start < 0 {
		start = 0
	}
	end := start + maxVisible
	if end > total {
		end = total
		start = end - maxVisible
		if start < 0 {
			start = 0
		}
	}
	return visRange{start: start, end: end}
}

func paneStyle(active bool) lipgloss.Style {
	if active {
		return activePane
	}
	return inactivePane
}
