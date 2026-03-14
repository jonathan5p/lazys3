package ui

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/jonathan5p/lazys3/internal/preview"
)

const defaultQuery = "SELECT * FROM data LIMIT 20"

func handleQueryKey(a AppModel, m tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch m.String() {
	case "esc":
		a.mode = ModeNormal
		a.queryState = QueryState{}
		return a, nil
	case "enter":
		return executeQuery(a)
	case "left":
		a.queryState = queryCursorLeft(a.queryState)
		return a, nil
	case "right":
		a.queryState = queryCursorRight(a.queryState)
		return a, nil
	case "ctrl+a", "home":
		a.queryState.cursor = 0
		return a, nil
	case "ctrl+e", "end":
		a.queryState.cursor = len(a.queryState.text)
		return a, nil
	case "backspace":
		a.queryState = queryDeleteBefore(a.queryState)
		return a, nil
	case "delete":
		a.queryState = queryDeleteAfter(a.queryState)
		return a, nil
	default:
		if isPrintable(m) {
			text := m.Text
			if text == "" {
				text = m.String()
			}
			a.queryState = queryInsert(a.queryState, text)
		}
		return a, nil
	}
}

func queryCursorLeft(q QueryState) QueryState {
	if q.cursor > 0 {
		q.cursor--
	}
	return q
}

func queryCursorRight(q QueryState) QueryState {
	if q.cursor < len(q.text) {
		q.cursor++
	}
	return q
}

func queryInsert(q QueryState, text string) QueryState {
	q.text = q.text[:q.cursor] + text + q.text[q.cursor:]
	q.cursor += len(text)
	return q
}

func queryDeleteBefore(q QueryState) QueryState {
	if q.cursor == 0 {
		return q
	}
	q.text = q.text[:q.cursor-1] + q.text[q.cursor:]
	q.cursor--
	return q
}

func queryDeleteAfter(q QueryState) QueryState {
	if q.cursor >= len(q.text) {
		return q
	}
	q.text = q.text[:q.cursor] + q.text[q.cursor+1:]
	return q
}

func executeQuery(a AppModel) (tea.Model, tea.Cmd) {
	a.mode = ModeNormal
	if a.preview.rawPath == "" {
		return a, nil
	}
	query := a.queryState.text
	path := a.preview.rawPath
	width := a.preview.width
	key := a.preview.fileKey
	a.queryState = QueryState{}
	return a, runQueryCmd(query, path, key, width)
}

func runQueryCmd(query, path, key string, width int) tea.Cmd {
	return func() tea.Msg {
		ft := preview.Detect(key, nil)
		formatter := buildQueryFormatter(ft, query, width)
		out, err := formatter.Format(path)
		if err != nil {
			return PreviewReadyMsg{Err: err, RawPath: path, FileKey: key}
		}
		stat, _ := os.Stat(path)
		rawPath := ""
		if stat != nil {
			rawPath = path
		}
		return PreviewReadyMsg{Content: out, RawPath: rawPath, FileKey: key}
	}
}

func buildQueryFormatter(ft preview.FileType, query string, width int) preview.Formatter {
	return preview.DuckDBFormatter{Query: query, Width: width}
}

func removeLastChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}
