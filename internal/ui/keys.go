package ui

import (
	"context"
	"io"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/jonathan5p/lazys3/internal/preview"
	s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)

func handleKey(a AppModel, m tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if a.mode == ModeHelp {
		if m.String() == "?" || m.String() == "esc" {
			a.mode = ModeNormal
		}
		return a, nil
	}

	switch m.String() {
	case "q", "ctrl+c":
		return a, tea.Quit
	case "?":
		a.mode = ModeHelp
	case "/":
		a.mode = ModeFilter
		a.filter.active = true
		if a.activePane == PaneObjects && a.objects.prefix != "" {
			a.filter.text = a.objects.prefix
		}
	case "h", "left":
		return handleLeft(a)
	case "l", "right":
		if a.activePane == PaneBuckets {
			return handleEnter(a)
		}
		return handlePaneFocusRight(a)
	case "j", "down":
		return handleCursorDown(a)
	case "k", "up":
		return handleCursorUp(a)
	case "G":
		return handleJumpBottom(a)
	case "g":
		if a.prevKey == "g" {
			a.prevKey = ""
			return handleJumpTop(a)
		}
		a.prevKey = "g"
	case "enter":
		return handleEnter(a)
	case "esc", "backspace":
		return handleBack(a)
	case "r":
		return handleReload(a)
	case "p":
		return handlePreview(a)
	case "d":
		return handleDownload(a)
	case "m":
		return handleMetadata(a)
	case ":":
		if a.activePane == PanePreview && a.stage == StageObjects {
			a.mode = ModeQuery
			a.queryState = QueryState{text: defaultQuery, cursor: len(defaultQuery), active: true}
			return a, nil
		}
	case "+", "=", "-":
		if a.activePane == PanePreview {
			var cmd tea.Cmd
			a.preview, cmd = a.preview.Update(m)
			return a, cmd
		}
	default:
		a.prevKey = ""
	}
	return a, nil
}

func handleLeft(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	if a.activePane == PanePreview {
		a.activePane = PaneObjects
		return a, nil
	}
	if a.stage == StageObjects && a.objects.prefix == "" {
		a.stage = StageBuckets
		a.activePane = PaneBuckets
		return a, nil
	}
	a.activePane = PaneBuckets
	return a, nil
}

func handlePaneFocusRight(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	switch a.activePane {
	case PaneBuckets:
		if a.buckets.SelectedBucket() != "" {
			a.activePane = PaneObjects
		}
	case PaneObjects:
		a.activePane = PanePreview
	}
	return a, nil
}

func handleCursorDown(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	switch a.activePane {
	case PaneBuckets:
		a.buckets, _ = a.buckets.Update(tea.KeyPressMsg{Code: 'j'})
		if a.stage == StageBuckets {
			return a, loadObjectsPreview(a.s3, a.buckets.SelectedBucket())
		}
	case PaneObjects:
		a.objects, _ = a.objects.Update(tea.KeyPressMsg{Code: 'j'})
	case PanePreview:
		a.preview, _ = a.preview.Update(tea.KeyPressMsg{Code: 'j'})
	}
	return a, nil
}

func handleCursorUp(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	switch a.activePane {
	case PaneBuckets:
		a.buckets, _ = a.buckets.Update(tea.KeyPressMsg{Code: 'k'})
		if a.stage == StageBuckets {
			return a, loadObjectsPreview(a.s3, a.buckets.SelectedBucket())
		}
	case PaneObjects:
		a.objects, _ = a.objects.Update(tea.KeyPressMsg{Code: 'k'})
	case PanePreview:
		a.preview, _ = a.preview.Update(tea.KeyPressMsg{Code: 'k'})
	}
	return a, nil
}

func handleJumpBottom(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	msg := tea.KeyPressMsg{Code: 'G'}
	switch a.activePane {
	case PaneBuckets:
		a.buckets, _ = a.buckets.Update(msg)
	case PaneObjects:
		a.objects, _ = a.objects.Update(msg)
	}
	return a, nil
}

func handleJumpTop(a AppModel) (tea.Model, tea.Cmd) {
	switch a.activePane {
	case PaneBuckets:
		a.buckets.cursor = 0
	case PaneObjects:
		a.objects.cursor = 0
	}
	return a, nil
}

func handleEnter(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	if a.activePane == PaneBuckets {
		bucket := a.buckets.SelectedBucket()
		if bucket == "" {
			return a, nil
		}
		a.stage = StageObjects
		a.activePane = PaneObjects
		a.objects.loading = true
		a.objects.prefix = ""
		return a, loadObjects(a.s3, bucket, "")
	}
	if a.activePane == PaneObjects {
		obj := a.objects.SelectedObject()
		if obj == nil {
			return a, nil
		}
		if obj.IsDir {
			a.objects.prefix = obj.Key
			a.objects.loading = true
			bucket := a.buckets.SelectedBucket()
			return a, loadObjects(a.s3, bucket, obj.Key)
		}
		a.activePane = PanePreview
	}
	return a, nil
}

func handleBack(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	if a.activePane == PaneObjects {
		if a.objects.prefix != "" {
			newPrefix := parentPrefix(a.objects.prefix)
			a.objects.prefix = newPrefix
			a.objects.loading = true
			bucket := a.buckets.SelectedBucket()
			return a, loadObjects(a.s3, bucket, newPrefix)
		}
		if a.stage == StageObjects {
			a.stage = StageBuckets
			a.activePane = PaneBuckets
			return a, nil
		}
	}
	if a.activePane == PanePreview {
		a.activePane = PaneObjects
	}
	return a, nil
}

func handleReload(a AppModel) (tea.Model, tea.Cmd) {
	a.prevKey = ""
	if a.activePane == PaneBuckets {
		a.buckets.loading = true
		return a, loadBuckets(a.s3)
	}
	if a.activePane == PaneObjects {
		bucket := a.buckets.SelectedBucket()
		if bucket == "" {
			return a, nil
		}
		a.objects.loading = true
		return a, loadObjects(a.s3, bucket, a.objects.prefix)
	}
	return a, nil
}

func handleMetadata(a AppModel) (tea.Model, tea.Cmd) {
	if a.activePane != PaneObjects {
		a.statusbar = a.statusbar.SetMessage("Metadata: select an object first")
		return a, nil
	}
	obj := a.objects.SelectedObject()
	if obj == nil || obj.IsDir {
		a.statusbar = a.statusbar.SetMessage("Metadata: select a file first")
		return a, nil
	}
	bucket := a.buckets.SelectedBucket()
	return a, loadMetadata(a.s3, bucket, obj.Key)
}

func loadMetadata(client s3pkg.Client, bucket, key string) tea.Cmd {
	if client == nil {
		return nil
	}
	return func() tea.Msg {
		meta, err := client.HeadObject(context.Background(), bucket, key)
		if err != nil {
			return MetadataLoadedMsg{Err: err}
		}
		return MetadataLoadedMsg{Content: formatMetadata(meta)}
	}
}

func formatMetadata(meta map[string]string) string {
	var sb strings.Builder
	for k, v := range meta {
		sb.WriteString(k + "=" + v + "\n")
	}
	return sb.String()
}

func handleMetadataKey(a AppModel, m tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.String() == "esc" {
		a.mode = ModeNormal
	}
	return a, nil
}

func handleDownload(a AppModel) (tea.Model, tea.Cmd) {
	if a.activePane != PaneObjects {
		a.statusbar = a.statusbar.SetMessage("Download: select an object first")
		return a, nil
	}
	obj := a.objects.SelectedObject()
	if obj == nil || obj.IsDir {
		a.statusbar = a.statusbar.SetMessage("Download: select a file first")
		return a, nil
	}
	a.mode = ModeDownload
	a.downloadPrompt = a.downloadPrompt.SetInput(objectFilename(obj.Key))
	return a, nil
}

func objectFilename(key string) string {
	parts := strings.Split(strings.TrimSuffix(key, "/"), "/")
	return parts[len(parts)-1]
}

func handleDownloadKey(a AppModel, m tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch m.String() {
	case "esc":
		a.mode = ModeNormal
		return a, nil
	case "enter":
		a.mode = ModeNormal
		return a, execDownload(a)
	case "backspace":
		inp := a.downloadPrompt.input
		if len(inp) > 0 {
			a.downloadPrompt = a.downloadPrompt.SetInput(inp[:len(inp)-1])
		}
		return a, nil
	default:
		if isPrintable(m) {
			a.downloadPrompt = a.downloadPrompt.SetInput(a.downloadPrompt.input + m.String())
		}
		return a, nil
	}
}

func execDownload(a AppModel) tea.Cmd {
	if a.s3 == nil {
		return nil
	}
	obj := a.objects.SelectedObject()
	if obj == nil {
		return nil
	}
	bucket := a.buckets.SelectedBucket()
	dest := a.downloadPrompt.input
	return func() tea.Msg {
		err := a.s3.DownloadObject(context.Background(), bucket, obj.Key, dest)
		return DownloadDoneMsg{Path: dest, Err: err}
	}
}

func handlePreview(a AppModel) (tea.Model, tea.Cmd) {
	if a.activePane != PaneObjects {
		a.statusbar = a.statusbar.SetMessage("Preview: select an object first")
		return a, nil
	}
	obj := a.objects.SelectedObject()
	if obj == nil || obj.IsDir {
		a.statusbar = a.statusbar.SetMessage("Preview: select a file first")
		return a, nil
	}
	a.preview.loading = true
	bucket := a.buckets.SelectedBucket()
	return a, loadPreview(a.s3, bucket, obj.Key, a.preview.width)
}

func loadPreview(client s3pkg.Client, bucket, key string, width int) tea.Cmd {
	if client == nil {
		return nil
	}
	return func() tea.Msg {
		rc, err := client.GetObject(context.Background(), bucket, key)
		if err != nil {
			return PreviewReadyMsg{Err: err}
		}
		defer rc.Close()
		return writeAndFormatPreview(rc, key, width)
	}
}

func writeAndFormatPreview(rc io.ReadCloser, key string, width int) PreviewReadyMsg {
	tmp, err := os.CreateTemp("", "lazys3-preview-*")
	if err != nil {
		return PreviewReadyMsg{Err: err}
	}
	defer tmp.Close()

	header := make([]byte, 512)
	n, _ := rc.Read(header)
	header = header[:n]
	if _, err := tmp.Write(header); err != nil {
		os.Remove(tmp.Name())
		return PreviewReadyMsg{Err: err}
	}
	if _, err := io.Copy(tmp, rc); err != nil {
		os.Remove(tmp.Name())
		return PreviewReadyMsg{Err: err}
	}
	tmp.Close()

	ft := preview.Detect(key, header)
	msg := formatForType(ft, tmp.Name(), width)
	msg.RawPath = tmp.Name()
	msg.FileKey = key
	return msg
}

func formatForType(ft preview.FileType, path string, width int) PreviewReadyMsg {
	switch ft {
	case preview.FileTypeJSON:
		out, err := preview.JSONFormatter{}.Format(path)
		return PreviewReadyMsg{Content: out, Err: err}
	case preview.FileTypeParquet:
		out, err := preview.ParquetFormatter{MaxRows: 20, Width: width}.Format(path)
		return PreviewReadyMsg{Content: out, Err: err}
	default:
		data, err := os.ReadFile(path)
		return PreviewReadyMsg{Content: string(data), Err: err}
	}
}

func loadBuckets(client s3pkg.Client) tea.Cmd {
	if client == nil {
		return nil
	}
	return func() tea.Msg {
		buckets, err := client.ListBuckets(context.Background())
		return BucketsLoadedMsg{Buckets: buckets, Err: err}
	}
}

func loadObjects(client s3pkg.Client, bucket, prefix string) tea.Cmd {
	if client == nil {
		return nil
	}
	return func() tea.Msg {
		objs, err := client.ListObjects(context.Background(), bucket, prefix)
		return ObjectsLoadedMsg{Objects: objs, Err: err}
	}
}

// loadObjectsPreview fires a ListObjects for the given bucket to populate the
// read-only peek pane in Stage 1 (BucketBrowse). It reuses ObjectsLoadedMsg;
// handleObjectsLoaded only updates items when stage == StageBuckets.
func loadObjectsPreview(client s3pkg.Client, bucket string) tea.Cmd {
	if client == nil || bucket == "" {
		return nil
	}
	return func() tea.Msg {
		objs, err := client.ListObjects(context.Background(), bucket, "")
		return ObjectsLoadedMsg{Objects: objs, Err: err}
	}
}

func parentPrefix(prefix string) string {
	if prefix == "" {
		return ""
	}
	trimmed := prefix[:len(prefix)-1]
	for i := len(trimmed) - 1; i >= 0; i-- {
		if trimmed[i] == '/' {
			return trimmed[:i+1]
		}
	}
	return ""
}
