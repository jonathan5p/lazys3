package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

// ObjectPane displays the list of S3 objects for the selected bucket/prefix.
type ObjectPane struct {
	items         []s3pkg.Object
	filteredItems []s3pkg.Object
	filterText    string
	cursor        int
	loading       bool
	err           error
	prefix        string
	title         string // optional override for the pane title
	pendingKey    string // if set, cursor jumps to this key after next SetItems
	width         int
	height        int
}

func newObjectPane(width, height int) ObjectPane {
	return ObjectPane{width: width, height: height}
}

// SetItems replaces the current object list and resets the cursor.
// If pendingKey is set, the cursor jumps to that key.
func (p *ObjectPane) SetItems(objects []s3pkg.Object) {
	p.items = objects
	p.cursor = 0
	if p.pendingKey != "" {
		for i, o := range objects {
			if o.Key == p.pendingKey {
				p.cursor = i
				break
			}
		}
		p.pendingKey = ""
	}
	p.applyFilter()
}

// SetFilter stores the filter text and re-filters the displayed items.
func (p *ObjectPane) SetFilter(text string) {
	p.filterText = text
	p.cursor = 0
	p.applyFilter()
}

func (p *ObjectPane) applyFilter() {
	fs := FilterState{text: p.filterText}
	names := make([]string, len(p.items))
	for i, o := range p.items {
		names[i] = o.DisplayName()
	}
	filtered := fs.Apply(names)
	filteredSet := make(map[string]bool, len(filtered))
	for _, n := range filtered {
		filteredSet[n] = true
	}
	p.filteredItems = nil
	for _, o := range p.items {
		if filteredSet[o.DisplayName()] {
			p.filteredItems = append(p.filteredItems, o)
		}
	}
}

func (p ObjectPane) visibleItems() []s3pkg.Object {
	if p.filteredItems != nil {
		return p.filteredItems
	}
	return p.items
}

// SelectedObject returns a pointer to the currently highlighted object.
func (p ObjectPane) SelectedObject() *s3pkg.Object {
	items := p.visibleItems()
	if len(items) == 0 {
		return nil
	}
	obj := items[p.cursor]
	return &obj
}

// Prefix returns the current navigation prefix.
func (p ObjectPane) Prefix() string {
	return p.prefix
}

// Update handles messages for the object pane.
func (p ObjectPane) Update(msg tea.Msg) (ObjectPane, tea.Cmd) {
	items := p.visibleItems()
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch m.String() {
		case "j", "down":
			p.cursor = clampInt(p.cursor+1, 0, len(items)-1)
		case "k", "up":
			p.cursor = clampInt(p.cursor-1, 0, len(items)-1)
		case "G":
			p.cursor = clampInt(len(items)-1, 0, len(items)-1)
		}
	}
	return p, nil
}

// View renders the object pane with the given active state.
func (p ObjectPane) View(active bool) string {
	style := paneStyle(active).Width(p.width - 2).Height(p.height - 2)

	if p.loading {
		return style.Render("Loading...")
	}
	if p.err != nil {
		return style.Render(fmt.Sprintf("Error: %v", p.err))
	}

	heading := titleStyle.Render(p.paneTitle()) + "\n\n"
	items := p.visibleItems()
	visible := visibleRange(p.cursor, p.height-4, len(items))
	content := heading
	for i := visible.start; i < visible.end; i++ {
		line := " " + items[i].DisplayName()
		if i == p.cursor {
			line = selectedItem.Render("> " + items[i].DisplayName())
		}
		content += line + "\n"
	}
	return style.Render(content)
}

func (p ObjectPane) paneTitle() string {
	if p.title != "" {
		return p.title
	}
	return "Objects"
}

func clampInt(v, lo, hi int) int {
	if lo > hi {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
