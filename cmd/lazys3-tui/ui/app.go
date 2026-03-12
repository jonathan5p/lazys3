package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jonathan5p/lazys3/cmd/lazys3-tui/logger"
	"github.com/jonathan5p/lazys3/cmd/lazys3-tui/s3"
)

type viewport struct {
	idx    int
	offset int
}

func (vp *viewport) scrollDown(height, max int) {
	if vp.idx < max-1 {
		vp.idx++
		if vp.idx-vp.offset >= height {
			vp.offset++
		}
	}
}

func (vp *viewport) scrollUp() {
	if vp.idx > 0 {
		vp.idx--
		if vp.idx < vp.offset {
			vp.offset--
		}
	}
}

func (vp *viewport) reset() {
	vp.idx = 0
	vp.offset = 0
}

func (vp *viewport) visibleSlice(height, total int) (int, int) {
	end := vp.offset + height
	if end > total {
		end = total
	}
	return vp.offset, end
}

type App struct {
	client   *s3.Client
	log      *logger.Logger
	gui      *gocui.Gui
	buckets  []string
	objects  []s3.Object
	bucketVP viewport
	objectVP viewport
	prefix   string
	loaded   bool
}

func NewApp(client *s3.Client, log *logger.Logger) *App {
	return &App{
		client: client,
		log:    log,
	}
}

func (a *App) Run(_ context.Context) error {
	a.log.Info("App starting")
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		return err
	}
	defer g.Close()
	a.gui = g

	g.FgColor = gocui.ColorDefault
	g.BgColor = gocui.ColorDefault

	g.SetLayout(a.layout)

	if err := a.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		a.log.Error("MainLoop error: %v", err)
		return err
	}
	a.log.Info("App exiting normally")
	return nil
}

func (a *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	vBuckets, err := g.SetView("buckets", 0, 0, maxX/3-1, maxY-3)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	vBuckets.Title = "Buckets"

	vObjects, err := g.SetView("objects", maxX/3, 0, maxX-1, maxY-3)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	vObjects.Title = "Objects"

	vStatus, err := g.SetView("status", 0, maxY-3, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	vStatus.Frame = false

	if !a.loaded {
		a.loaded = true
		if err := g.SetCurrentView("buckets"); err != nil {
			return err
		}
		ctx := context.Background()
		a.log.Debug("layout: first render, loading buckets")
		if err := a.loadBuckets(ctx); err != nil {
			a.log.Error("layout: loadBuckets failed: %v", err)
			fmt.Fprintln(vStatus, "Error:", err)
			return nil
		}
	}

	if a.bucketVP.idx < len(a.buckets) {
		fmt.Fprintln(vStatus, "Bucket:", a.buckets[a.bucketVP.idx], "| Prefix:", a.prefix)
	} else {
		fmt.Fprintln(vStatus, "No buckets found")
	}

	return nil
}

func (a *App) loadBuckets(ctx context.Context) error {
	a.log.Info("loadBuckets: listing buckets")
	buckets, err := a.client.ListBuckets(ctx)
	if err != nil {
		return err
	}
	// TODO: Remove this limit and implement proper pagination
	a.buckets = buckets[:50]
	a.log.Info("loadBuckets: got %d buckets", len(buckets))

	if len(buckets) > 0 {
		if err := a.loadObjects(ctx, buckets[0], ""); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) loadObjects(ctx context.Context, bucket, prefix string) error {
	a.log.Info("loadObjects: bucket=%q prefix=%q", bucket, prefix)
	objects, err := a.client.ListObjects(ctx, bucket, prefix)
	if err != nil {
		a.log.Error("loadObjects failed: %v", err)
		return err
	}
	a.objects = objects
	a.prefix = prefix
	a.objectVP.reset()
	a.log.Debug("loadObjects: %d objects loaded, rendering", len(objects))
	return a.render()
}

func (a *App) render() error {
	a.log.Debug("render: buckets=%d objects=%d bucketIdx=%d objIdx=%d", len(a.buckets), len(a.objects), a.bucketVP.idx, a.objectVP.idx)
	vBuckets, err := a.gui.View("buckets")
	if err != nil {
		return err
	}
	_, vBucketsY := vBuckets.Size()
	vBuckets.Clear()
	start, end := a.bucketVP.visibleSlice(vBucketsY, len(a.buckets))
	for i, b := range a.buckets[start:end] {
		if start+i == a.bucketVP.idx {
			fmt.Fprintln(vBuckets, "> "+b)
		} else {
			fmt.Fprintln(vBuckets, "  "+b)
		}
	}

	vObjects, err := a.gui.View("objects")
	if err != nil {
		return err
	}
	_, vObjectsY := vObjects.Size()
	vObjects.Clear()
	oStart, oEnd := a.objectVP.visibleSlice(vObjectsY, len(a.objects))
	for i, obj := range a.objects[oStart:oEnd] {
		prefix := "  "
		if obj.IsDir {
			prefix = "> "
		}
		if oStart+i == a.objectVP.idx {
			prefix = ">> "
		}
		fmt.Fprintln(vObjects, prefix+obj.Format())
	}
	return nil
}

func (a *App) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("buckets", gocui.KeyArrowUp, gocui.ModNone, a.bucketUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("buckets", gocui.KeyArrowDown, gocui.ModNone, a.bucketDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("buckets", gocui.KeyEnter, gocui.ModNone, a.selectBucket); err != nil {
		return err
	}
	if err := g.SetKeybinding("buckets", gocui.KeyArrowRight, gocui.ModNone, a.switchFocus); err != nil {
		return err
	}

	if err := g.SetKeybinding("objects", gocui.KeyArrowUp, gocui.ModNone, a.objUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("objects", gocui.KeyArrowDown, gocui.ModNone, a.objDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("objects", gocui.KeyEnter, gocui.ModNone, a.enterObject); err != nil {
		return err
	}
	if err := g.SetKeybinding("objects", gocui.KeyEsc, gocui.ModNone, a.goBack); err != nil {
		return err
	}
	if err := g.SetKeybinding("objects", gocui.KeyArrowLeft, gocui.ModNone, a.switchFocusBack); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, a.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, a.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gocui.ModNone, a.download); err != nil {
		return err
	}

	return nil
}

func (a *App) bucketUp(g *gocui.Gui, v *gocui.View) error {
	a.log.Debug("key: bucket up (idx=%d)", a.bucketVP.idx)
	a.bucketVP.scrollUp()
	return a.selectBucket(g, v)
}

func (a *App) bucketDown(g *gocui.Gui, v *gocui.View) error {
	a.log.Debug("key: bucket down (idx=%d)", a.bucketVP.idx)
	_, vY := v.Size()
	a.bucketVP.scrollDown(vY, len(a.buckets))
	return a.selectBucket(g, v)
}

func (a *App) selectBucket(g *gocui.Gui, v *gocui.View) error {
	a.log.Info("key: select bucket idx=%d", a.bucketVP.idx)
	if len(a.buckets) == 0 {
		return nil
	}
	ctx := context.Background()
	a.prefix = ""
	return a.loadObjects(ctx, a.buckets[a.bucketVP.idx], "")
}

func (a *App) objUp(g *gocui.Gui, v *gocui.View) error {
	a.log.Debug("key: obj up (idx=%d)", a.objectVP.idx)
	a.objectVP.scrollUp()
	return a.render()
}

func (a *App) objDown(g *gocui.Gui, v *gocui.View) error {
	a.log.Debug("key: obj down (idx=%d)", a.objectVP.idx)
	_, vY := v.Size()
	a.objectVP.scrollDown(vY, len(a.objects))
	return a.render()
}

func (a *App) enterObject(g *gocui.Gui, v *gocui.View) error {
	a.log.Debug("key: enter object idx=%d", a.objectVP.idx)
	if len(a.objects) == 0 || len(a.buckets) == 0 {
		return nil
	}
	obj := a.objects[a.objectVP.idx]
	if obj.IsDir {
		ctx := context.Background()
		return a.loadObjects(ctx, a.buckets[a.bucketVP.idx], obj.Key)
	}
	return nil
}

func (a *App) goBack(g *gocui.Gui, v *gocui.View) error {
	a.log.Debug("key: go back prefix=%q", a.prefix)
	if a.prefix == "" || len(a.buckets) == 0 {
		return nil
	}
	ctx := context.Background()
	parts := strings.Split(strings.TrimSuffix(a.prefix, "/"), "/")
	if len(parts) > 1 {
		newPrefix := strings.Join(parts[:len(parts)-1], "/") + "/"
		return a.loadObjects(ctx, a.buckets[a.bucketVP.idx], newPrefix)
	}
	return a.loadObjects(ctx, a.buckets[a.bucketVP.idx], "")
}

func (a *App) download(g *gocui.Gui, v *gocui.View) error {
	if len(a.objects) == 0 || len(a.buckets) == 0 {
		return nil
	}
	obj := a.objects[a.objectVP.idx]
	if obj.IsDir {
		return nil
	}
	ctx := context.Background()
	filename := obj.Key
	parts := strings.Split(obj.Key, "/")
	if len(parts) > 0 {
		filename = parts[len(parts)-1]
	}
	err := a.client.DownloadObject(ctx, a.buckets[a.bucketVP.idx], obj.Key, filename)
	if err != nil {
		vStatus, _ := g.View("status")
		fmt.Fprintln(vStatus, "Error downloading:", err)
		return nil
	}
	vStatus, _ := g.View("status")
	fmt.Fprintln(vStatus, "Downloaded:", filename)
	return nil
}

func (a *App) switchFocus(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView("objects")
}

func (a *App) switchFocusBack(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView("buckets")
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	a.log.Info("key: quit")
	return gocui.ErrQuit
}
