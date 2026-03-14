package ui

import "github.com/charmbracelet/lipgloss"

func renderLayout(a AppModel) string {
	sbView := a.statusbar.View()
	filterView := renderFilterBar(a)
	queryView := renderQueryBar(a)
	reserved := lipgloss.Height(sbView) + lipgloss.Height(filterView) + lipgloss.Height(queryView)
	paneHeight := a.height - reserved - 1
	if paneHeight < 1 {
		paneHeight = 1
	}

	a.buckets.height = paneHeight
	a.objects.height = paneHeight
	a.preview.height = paneHeight

	var base string
	if a.stage == StageObjects {
		// Stage 2 — BucketBrowse exited; objects left, preview right (50/50)
		base = renderObjectStage(a, filterView, sbView, queryView)
	} else {
		// Stage 1 — BucketBrowse; buckets left, objects peek right (50/50)
		base = renderBucketStage(a, filterView, sbView, queryView)
	}

	if a.mode == ModeMetadata {
		return a.metaOverlay.View()
	}
	if a.mode == ModeDownload {
		return a.downloadPrompt.View()
	}
	if a.mode == ModeHelp {
		return a.helpOverlay.View()
	}
	return base
}

// renderBucketStage renders Stage 1: Buckets (50%) | Objects peek (50%).
func renderBucketStage(a AppModel, filterView, sbView, queryView string) string {
	leftW := a.width / 2
	rightW := a.width - leftW

	a.buckets.width = leftW
	a.objects.width = rightW

	bucket := a.buckets.SelectedBucket()
	a.objects.title = objectsPeekTitle(bucket)

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		a.buckets.View(a.activePane == PaneBuckets),
		a.objects.View(false),
	)
	return joinWithBars(panes, filterView, queryView, sbView)
}

// renderObjectStage renders Stage 2: Objects (50%) | Preview (50%).
func renderObjectStage(a AppModel, filterView, sbView, queryView string) string {
	leftW := a.width / 2
	rightW := a.width - leftW

	a.objects.width = leftW
	a.preview.width = rightW

	a.objects.title = ""

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		a.objects.View(a.activePane == PaneObjects),
		a.preview.View(a.activePane == PanePreview),
	)
	return joinWithBars(panes, filterView, queryView, sbView)
}

func objectsPeekTitle(bucket string) string {
	if bucket == "" {
		return "Objects"
	}
	return "Objects: " + bucket
}

func joinWithBars(panes, filterView, queryView, sbView string) string {
	parts := []string{panes}
	if filterView != "" {
		parts = append(parts, filterView)
	}
	if queryView != "" {
		parts = append(parts, queryView)
	}
	parts = append(parts, sbView)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func renderFilterBar(a AppModel) string {
	if a.mode != ModeFilter {
		return ""
	}
	text := filterBar.Width(a.width).Render("/ " + a.filter.text + "█")
	return text
}

func renderQueryBar(a AppModel) string {
	if a.mode != ModeQuery {
		return ""
	}
	q := a.queryState
	before := q.text[:q.cursor]
	after := q.text[q.cursor:]
	return queryBar.Width(a.width).Render("SQL> " + before + "█" + after)
}
