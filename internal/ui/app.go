package ui

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

// AppStage represents the two-stage navigation layout of the TUI.
type AppStage int

const (
	// StageBuckets — Stage 1: user is browsing buckets; objects pane is a read-only peek.
	StageBuckets AppStage = iota
	// StageObjects — Stage 2: user has entered a bucket; objects pane is fully navigable.
	StageObjects
)

// AppMode represents the current interaction mode of the TUI.
type AppMode int

const (
	ModeNormal AppMode = iota
	ModeFilter
	ModeDownload
	ModeHelp
	ModeMetadata
	ModeQuery
)

// QueryState tracks the current SQL query input with cursor position.
type QueryState struct {
	text   string
	cursor int // byte position within text
	active bool
}

// PaneID identifies which pane currently has focus.
type PaneID int

const (
	PaneBuckets PaneID = iota
	PaneObjects
	PanePreview
)

// AppModel is the root Bubble Tea model.
type AppModel struct {
	mode           AppMode
	stage          AppStage
	activePane     PaneID
	width          int
	height         int
	buckets        BucketPane
	objects        ObjectPane
	preview        PreviewPane
	statusbar      StatusBar
	s3             s3pkg.Client
	prevKey        string
	filter         FilterState
	metaOverlay    MetadataOverlay
	downloadPrompt DownloadPrompt
	helpOverlay    HelpOverlay
	queryState     QueryState
}

// NewApp constructs the root AppModel with sensible defaults.
func NewApp(s3Client s3pkg.Client, width, height int) AppModel {
	return AppModel{
		mode:           ModeNormal,
		activePane:     PaneBuckets,
		width:          width,
		height:         height,
		buckets:        newBucketPane(width/3, height-1),
		objects:        newObjectPane(width/3, height-1),
		preview:        newPreviewPane(width/3, height-1),
		statusbar:      newStatusBar(width),
		s3:             s3Client,
		metaOverlay:    newMetadataOverlay(width, height),
		downloadPrompt: newDownloadPrompt(width, height),
		helpOverlay:    newHelpOverlay(width, height),
	}
}

func (a AppModel) Init() tea.Cmd {
	if a.s3 == nil {
		return nil
	}
	return func() tea.Msg {
		buckets, err := a.s3.ListBuckets(context.Background())
		return BucketsLoadedMsg{Buckets: buckets, Err: err}
	}
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		a = a.handleResize(m)
	case BucketsLoadedMsg:
		a, cmd = a.handleBucketsLoaded(m)
	case ObjectsLoadedMsg:
		a = a.handleObjectsLoaded(m)
	case PreviewReadyMsg:
		a = a.handlePreviewReady(m)
	case MetadataLoadedMsg:
		a = a.handleMetadataLoaded(m)
	case DownloadDoneMsg:
		a = a.handleDownloadDone(m)
	case tea.KeyPressMsg:
		var model tea.Model
		if a.mode == ModeFilter {
			model, cmd = handleFilterKey(a, m)
		} else if a.mode == ModeMetadata {
			model, cmd = handleMetadataKey(a, m)
		} else if a.mode == ModeDownload {
			model, cmd = handleDownloadKey(a, m)
		} else if a.mode == ModeQuery {
			model, cmd = handleQueryKey(a, m)
		} else {
			model, cmd = handleKey(a, m)
		}
		a = model.(AppModel)
	case tea.PasteMsg:
		if a.mode == ModeFilter {
			a.filter.text += m.Content
			applyFilterToApp(&a, a.filter.text)
		} else if a.mode == ModeDownload {
			a.downloadPrompt = a.downloadPrompt.SetInput(a.downloadPrompt.input + m.Content)
		} else if a.mode == ModeQuery {
			a.queryState = queryInsert(a.queryState, m.Content)
		}
	default:
		return a, nil
	}
	a.statusbar = a.statusbar.SetMode(a.mode, a.filter.text)
	if a.stage == StageObjects {
		bucket := a.buckets.SelectedBucket()
		a.statusbar = a.statusbar.SetMessage("s3://" + bucket + "/" + a.objects.prefix)
	}
	return a, cmd
}

func (a AppModel) handleResize(m tea.WindowSizeMsg) AppModel {
	a.width = m.Width
	a.height = m.Height
	a.statusbar.width = m.Width
	paneH := m.Height - 1
	a.buckets.width = m.Width / 3
	a.buckets.height = paneH
	a.objects.width = m.Width / 3
	a.objects.height = paneH
	a.preview.width = m.Width - (m.Width/3)*2
	a.preview.height = paneH
	return a
}

func (a AppModel) handleBucketsLoaded(m BucketsLoadedMsg) (AppModel, tea.Cmd) {
	a.buckets.loading = false
	if m.Err != nil {
		a.statusbar = a.statusbar.SetError(m.Err)
		return a, nil
	}
	a.buckets.SetItems(m.Buckets)
	a.statusbar = a.statusbar.SetMessage(fmt.Sprintf("%d buckets loaded", len(m.Buckets)))
	return a, loadObjectsPreview(a.s3, a.buckets.SelectedBucket())
}

func (a AppModel) handleObjectsLoaded(m ObjectsLoadedMsg) AppModel {
	a.objects.loading = false
	if m.Err != nil {
		a.statusbar = a.statusbar.SetError(m.Err)
		return a
	}
	a.objects.SetItems(m.Objects)
	a.statusbar = a.statusbar.SetMessage(fmt.Sprintf("%d objects", len(m.Objects)))
	return a
}

func (a AppModel) handlePreviewReady(m PreviewReadyMsg) AppModel {
	if a.preview.rawPath != "" && a.preview.rawPath != m.RawPath {
		os.Remove(a.preview.rawPath)
	}
	a.preview.loading = false
	if m.Err != nil {
		a.statusbar = a.statusbar.SetError(m.Err)
		return a
	}
	a.preview = a.preview.SetContent(m.Content)
	if m.RawPath != "" {
		a.preview.rawPath = m.RawPath
		a.preview.fileKey = m.FileKey
	}
	a.activePane = PanePreview
	return a
}

func (a AppModel) handleMetadataLoaded(m MetadataLoadedMsg) AppModel {
	if m.Err != nil {
		a.statusbar = a.statusbar.SetError(m.Err)
		return a
	}
	a.metaOverlay = a.metaOverlay.SetContent(m.Content)
	a.mode = ModeMetadata
	return a
}

func (a AppModel) handleDownloadDone(m DownloadDoneMsg) AppModel {
	if m.Err != nil {
		a.statusbar = a.statusbar.SetError(m.Err)
		return a
	}
	a.statusbar = a.statusbar.SetMessage("Downloaded: " + m.Path)
	return a
}

func (a AppModel) View() tea.View {
	v := tea.NewView(renderLayout(a))
	v.AltScreen = true
	return v
}
