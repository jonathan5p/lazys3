package ui

import (
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
)

// FilterState tracks the current filter text and whether filtering is active.
type FilterState struct {
	text   string
	active bool
}

// Apply returns items whose lowercased form contains the filter text.
func (f FilterState) Apply(items []string) []string {
	if f.text == "" {
		return items
	}
	lower := strings.ToLower(f.text)
	var result []string
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), lower) {
			result = append(result, item)
		}
	}
	return result
}

// s3URL holds a parsed S3 URL.
type s3URL struct {
	bucket string
	key    string
	prefix string // directory portion of key (everything before the last /)
}

// parseS3URL parses an s3:// URL into bucket and key components.
// Returns nil if the input is not a valid S3 URL.
func parseS3URL(input string) *s3URL {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "s3://") {
		return nil
	}
	rest := strings.TrimPrefix(input, "s3://")
	idx := strings.Index(rest, "/")
	if idx < 0 {
		return &s3URL{bucket: rest}
	}
	bucket := rest[:idx]
	key := rest[idx+1:]
	prefix := ""
	if slash := strings.LastIndex(key, "/"); slash >= 0 {
		prefix = key[:slash+1]
	}
	return &s3URL{bucket: bucket, key: key, prefix: prefix}
}

func handleFilterKey(a AppModel, m tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch m.String() {
	case "esc":
		a.filter = FilterState{}
		a.mode = ModeNormal
		applyFilterToApp(&a, "")
		return a, nil
	case "enter":
		return handleFilterEnter(a)
	case "backspace":
		if len(a.filter.text) > 0 {
			a.filter.text = a.filter.text[:len(a.filter.text)-1]
		}
		applyFilterToApp(&a, a.filter.text)
		return a, nil
	default:
		if isPrintable(m) {
			text := m.Text
			if text == "" {
				text = m.String()
			}
			a.filter.text += text
			applyFilterToApp(&a, a.filter.text)
		}
		return a, nil
	}
}

func handleFilterEnter(a AppModel) (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(a.filter.text)
	a.filter = FilterState{}
	a.mode = ModeNormal
	applyFilterToApp(&a, "")

	if url := parseS3URL(text); url != nil {
		return navigateToS3URL(a, url)
	}
	if strings.Contains(text, "/") && a.activePane == PaneObjects {
		return navigateToPrefix(a, text)
	}
	return a, nil
}

func navigateToS3URL(a AppModel, url *s3URL) (tea.Model, tea.Cmd) {
	a.buckets.selectByName(url.bucket)
	a.stage = StageObjects
	a.activePane = PaneObjects
	a.objects.prefix = url.prefix
	a.objects.loading = true
	a.objects.pendingKey = url.key
	return a, loadObjects(a.s3, url.bucket, url.prefix)
}

func navigateToPrefix(a AppModel, prefix string) (tea.Model, tea.Cmd) {
	bucket := a.buckets.SelectedBucket()
	if bucket == "" {
		return a, nil
	}
	a.stage = StageObjects
	a.activePane = PaneObjects
	a.objects.prefix = prefix
	a.objects.loading = true
	return a, loadObjects(a.s3, bucket, prefix)
}

func isPrintable(m tea.KeyPressMsg) bool {
	if m.Text != "" {
		return true
	}
	if m.Code > 0 && m.Code < 32 {
		return false
	}
	s := m.String()
	return len(s) == 1 && unicode.IsPrint(rune(s[0]))
}

func applyFilterToApp(a *AppModel, text string) {
	a.buckets.SetFilter(text)
	a.objects.SetFilter(text)
}
