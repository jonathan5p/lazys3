package ui

import (
	"testing"

	"github.com/jroimartin/gocui"
)

type fakeView struct {
	name    string
	fgColor gocui.Attribute
}

func (f *fakeView) Name() string                 { return f.name }
func (f *fakeView) setFgColor(c gocui.Attribute) { f.fgColor = c }
func (f *fakeView) getFgColor() gocui.Attribute  { return f.fgColor }

func TestViewHighlighting(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		otherView      string
		highlightColor gocui.Attribute
		defaultColor   gocui.Attribute
	}{
		{
			name:           "highlight view1",
			input:          "testView1",
			otherView:      "testView2",
			highlightColor: gocui.ColorGreen,
			defaultColor:   gocui.ColorDefault,
		},
		{
			name:           "highlight view2",
			input:          "testView2",
			otherView:      "testView1",
			highlightColor: gocui.ColorGreen,
			defaultColor:   gocui.ColorDefault,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			views := []highlightable{
				&fakeView{name: tc.input, fgColor: gocui.ColorDefault},
				&fakeView{name: tc.otherView, fgColor: gocui.ColorDefault},
			}

			if err := applyHighlight(views, tc.input); err != nil {
				t.Fatalf("applyHighlight(%q) returned error: %v", tc.input, err)
			}

			if views[0].getFgColor() != tc.highlightColor {
				t.Errorf("active view %q: expected FgColor %v, got %v", tc.input, tc.highlightColor, views[0].getFgColor())
			}
			if views[1].getFgColor() != tc.defaultColor {
				t.Errorf("inactive view %q: expected FgColor %v, got %v", tc.otherView, tc.defaultColor, views[1].getFgColor())
			}
		})
	}
}
