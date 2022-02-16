package glfw

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"context"
	"image"
	_ "image/png" // for the icon
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/internal/painter"
	"github.com/bhojpur/gui/pkg/engine/internal/painter/gl"
)

const (
	scrollAccelerateRate   = float64(5)
	scrollAccelerateCutoff = float64(5)
	scrollSpeed            = float32(10)
	doubleClickDelay       = 300 // ms (maximum interval between clicks for double click detection)
	dragMoveThreshold      = 2   // how far can we move before it is a drag
	windowIconSize         = 256
)

var (
	cursorMap    map[desktop.StandardCursor]*glfw.Cursor
	defaultTitle = "Bhojpur Application"
)

func initCursors() {
	cursorMap = map[desktop.StandardCursor]*glfw.Cursor{
		desktop.DefaultCursor:   glfw.CreateStandardCursor(glfw.ArrowCursor),
		desktop.TextCursor:      glfw.CreateStandardCursor(glfw.IBeamCursor),
		desktop.CrosshairCursor: glfw.CreateStandardCursor(glfw.CrosshairCursor),
		desktop.PointerCursor:   glfw.CreateStandardCursor(glfw.HandCursor),
		desktop.HResizeCursor:   glfw.CreateStandardCursor(glfw.HResizeCursor),
		desktop.VResizeCursor:   glfw.CreateStandardCursor(glfw.VResizeCursor),
		desktop.HiddenCursor:    nil,
	}
}

// Declare conformity to Window interface
var _ gui.Window = (*window)(nil)

type window struct {
	common.Window

	viewport   *glfw.Window
	viewLock   sync.RWMutex
	createLock sync.Once
	decorate   bool
	closing    bool
	fixedSize  bool

	cursor       desktop.Cursor
	customCursor *glfw.Cursor
	canvas       *glCanvas
	driver       *gLDriver
	title        string
	icon         gui.Resource
	mainmenu     *gui.MainMenu

	clipboard gui.Clipboard

	master     bool
	fullScreen bool
	centered   bool
	visible    bool

	mouseLock            sync.RWMutex
	mousePos             gui.Position
	mouseDragged         gui.Draggable
	mouseDraggedObjStart gui.Position
	mouseDraggedOffset   gui.Position
	mouseDragPos         gui.Position
	mouseDragStarted     bool
	mouseButton          desktop.MouseButton
	mouseOver            desktop.Hoverable
	mouseLastClick       gui.CanvasObject
	mousePressed         gui.CanvasObject
	mouseClickCount      int
	mouseCancelFunc      context.CancelFunc

	onClosed           func()
	onCloseIntercepted func()

	menuTogglePending       gui.KeyName
	menuDeactivationPending gui.KeyName

	xpos, ypos                      int
	width, height                   int
	requestedWidth, requestedHeight int
	shouldWidth, shouldHeight       int
	shouldExpand                    bool

	pending []func()
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title

	w.runOnMainWhenCreated(func() {
		w.viewport.SetTitle(title)
	})
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

func (w *window) SetFullScreen(full bool) {
	w.fullScreen = full
	if !w.visible {
		return
	}

	runOnMain(func() {
		monitor := w.getMonitorForWindow()
		mode := monitor.GetVideoMode()

		if full {
			w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		} else {
			if w.width == 0 && w.height == 0 { // if we were fullscreen on creation...
				w.width, w.height = w.screenSize(w.canvas.Size())
			}
			w.viewport.SetMonitor(nil, w.xpos, w.ypos, w.width, w.height, 0)
		}
	})
}

func (w *window) CenterOnScreen() {
	w.centered = true

	if w.view() != nil {
		runOnMain(w.doCenterOnScreen)
	}
}

func (w *window) doCenterOnScreen() {
	viewWidth, viewHeight := w.screenSize(w.canvas.size)
	if w.width > viewWidth { // in case our window has not called back to canvas size yet
		viewWidth = w.width
	}
	if w.height > viewHeight {
		viewHeight = w.height
	}

	// get window dimensions in pixels
	monitor := w.getMonitorForWindow()
	monMode := monitor.GetVideoMode()

	// these come into play when dealing with multiple monitors
	monX, monY := monitor.GetPos()

	// math them to the middle
	newX := (monMode.Width / 2) - (viewWidth / 2) + monX
	newY := (monMode.Height / 2) - (viewHeight / 2) + monY

	// set new window coordinates
	w.viewport.SetPos(newX, newY)
}

// minSizeOnScreen gets the padded minimum size of a window content in screen pixels
func (w *window) minSizeOnScreen() (int, int) {
	// get minimum size of content inside the window
	return w.screenSize(w.canvas.MinSize())
}

// screenSize computes the actual output size of the given content size in screen pixels
func (w *window) screenSize(canvasSize gui.Size) (int, int) {
	return internal.ScaleInt(w.canvas, canvasSize.Width), internal.ScaleInt(w.canvas, canvasSize.Height)
}

func (w *window) RequestFocus() {
	if isWayland {
		return
	}

	w.runOnMainWhenCreated(w.viewport.Focus)
}

func (w *window) Resize(size gui.Size) {
	// we cannot perform this until window is prepared as we don't know it's scale!
	bigEnough := size.Max(w.canvas.canvasSize(w.canvas.Content().MinSize()))
	w.runOnMainWhenCreated(func() {
		w.viewLock.Lock()

		width, height := internal.ScaleInt(w.canvas, bigEnough.Width), internal.ScaleInt(w.canvas, bigEnough.Height)
		if w.fixedSize || !w.visible { // fixed size ignores future `resized` and if not visible we may not get the event
			w.shouldWidth, w.shouldHeight = width, height
			w.width, w.height = width, height
		}
		w.viewLock.Unlock()
		w.requestedWidth, w.requestedHeight = width, height
		w.viewport.SetSize(width, height)
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed

	if w.view() != nil {
		w.runOnMainWhenCreated(w.fitContent)
	}
}

func (w *window) Padded() bool {
	return w.canvas.padded
}

func (w *window) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)

	w.runOnMainWhenCreated(w.fitContent)
}

func (w *window) Icon() gui.Resource {
	if w.icon == nil {
		return gui.CurrentApp().Icon()
	}

	return w.icon
}

func (w *window) SetIcon(icon gui.Resource) {
	w.icon = icon
	if icon == nil {
		appIcon := gui.CurrentApp().Icon()
		if appIcon != nil {
			w.SetIcon(appIcon)
		}
		return
	}

	w.runOnMainWhenCreated(func() {
		if w.icon == nil {
			w.viewport.SetIcon(nil)
			return
		}

		var img image.Image
		if painter.IsResourceSVG(w.icon) {
			img = painter.PaintImage(&canvas.Image{Resource: w.icon}, nil, windowIconSize, windowIconSize)
		} else {
			pix, _, err := image.Decode(bytes.NewReader(w.icon.Content()))
			if err != nil {
				gui.LogError("Failed to decode image for window icon", err)
				return
			}
			img = pix
		}

		w.viewport.SetIcon([]image.Image{img})
	})
}

func (w *window) SetMaster() {
	w.master = true
}

func (w *window) MainMenu() *gui.MainMenu {
	return w.mainmenu
}

func (w *window) SetMainMenu(menu *gui.MainMenu) {
	w.mainmenu = menu
	w.runOnMainWhenCreated(func() {
		w.canvas.buildMenu(w, menu)
	})
}

func (w *window) fitContent() {
	if w.canvas.Content() == nil || (w.fullScreen && w.visible) {
		return
	}

	if w.isClosing() {
		return
	}

	minWidth, minHeight := w.minSizeOnScreen()
	w.viewLock.RLock()
	view := w.viewport
	w.viewLock.RUnlock()
	w.shouldWidth, w.shouldHeight = w.width, w.height
	if w.width < minWidth || w.height < minHeight {
		if w.width < minWidth {
			w.shouldWidth = minWidth
		}
		if w.height < minHeight {
			w.shouldHeight = minHeight
		}
		w.viewLock.Lock()
		w.shouldExpand = true // queue the resize to happen on main
		w.viewLock.Unlock()
	}
	if w.fixedSize {
		w.shouldWidth, w.shouldHeight = w.requestedWidth, w.requestedHeight
		view.SetSizeLimits(w.requestedWidth, w.requestedHeight, w.requestedWidth, w.requestedHeight)
	} else {
		view.SetSizeLimits(minWidth, minHeight, glfw.DontCare, glfw.DontCare)
	}
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *window) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	x, y := w.xpos, w.ypos
	if w.fullScreen {
		x, y = w.viewport.GetPos()
	}
	xOff := x + (w.width / 2)
	yOff := y + (w.height / 2)

	for _, monitor := range glfw.GetMonitors() {
		x, y := monitor.GetPos()

		if x > xOff || y > yOff {
			continue
		}
		if x+monitor.GetVideoMode().Width <= xOff || y+monitor.GetVideoMode().Height <= yOff {
			continue
		}

		return monitor
	}

	// try built-in function to detect monitor if above logic didn't succeed
	// if it doesn't work then return primary monitor as default
	monitor := w.viewport.GetMonitor()
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}
	return monitor
}

func (w *window) calculatedScale() float32 {
	return calculateScale(userScale(), gui.CurrentDevice().SystemScaleForWindow(w), w.detectScale())
}

func (w *window) detectScale() float32 {
	monitor := w.getMonitorForWindow()
	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	return calculateDetectedScale(widthMm, widthPx)
}

func (w *window) detectTextureScale() float32 {
	winWidth, _ := w.viewport.GetSize()
	texWidth, _ := w.viewport.GetFramebufferSize()
	return float32(texWidth) / float32(winWidth)
}

func (w *window) Show() {
	go w.doShow()
}

func (w *window) doShow() {
	if w.view() != nil {
		w.doShowAgain()
		return
	}

	for !running() {
		time.Sleep(time.Millisecond * 10)
	}
	w.createLock.Do(w.create)
	if w.view() == nil {
		return
	}

	runOnMain(func() {
		w.viewLock.Lock()
		w.visible = true
		w.viewLock.Unlock()
		w.viewport.SetTitle(w.title)

		if w.centered {
			w.doCenterOnScreen() // lastly center if that was requested
		}
		w.viewport.Show()

		// save coordinates
		w.xpos, w.ypos = w.viewport.GetPos()

		if w.fullScreen { // this does not work if called before viewport.Show()
			go func() {
				time.Sleep(time.Millisecond * 100)
				w.SetFullScreen(true)
			}()
		}
	})

	// show top canvas element
	if w.canvas.Content() != nil {
		w.canvas.Content().Show()
	}
}

func (w *window) Hide() {
	if w.isClosing() {
		return
	}

	runOnMain(func() {
		w.viewLock.Lock()
		w.visible = false
		w.viewport.Hide()
		w.viewLock.Unlock()

		// hide top canvas element
		if w.canvas.Content() != nil {
			w.canvas.Content().Hide()
		}
	})
}

func (w *window) Close() {
	if w.isClosing() {
		return
	}

	// set w.closing flag inside draw thread to ensure we can free textures
	runOnDraw(w, func() {
		w.viewLock.Lock()
		w.closing = true
		w.viewLock.Unlock()
		w.viewport.SetShouldClose(true)
		cache.RangeTexturesFor(w.canvas, func(obj gui.CanvasObject) {
			w.canvas.Painter().Free(obj)
		})
	})

	w.canvas.WalkTrees(nil, func(node *common.RenderCacheNode) {
		if wid, ok := node.Obj().(gui.Widget); ok {
			cache.DestroyRenderer(wid)
		}
	})

	// trigger callbacks
	if w.onClosed != nil {
		w.QueueEvent(w.onClosed)
	}
}

func (w *window) ShowAndRun() {
	w.Show()
	w.driver.Run()
}

// Clipboard returns the system clipboard
func (w *window) Clipboard() gui.Clipboard {
	if w.view() == nil {
		return nil
	}

	if w.clipboard == nil {
		w.clipboard = &clipboard{window: w.viewport}
	}
	return w.clipboard
}

func (w *window) Content() gui.CanvasObject {
	return w.canvas.Content()
}

func (w *window) SetContent(content gui.CanvasObject) {
	w.viewLock.RLock()
	visible := w.visible
	w.viewLock.RUnlock()
	// hide old canvas element
	if visible && w.canvas.Content() != nil {
		w.canvas.Content().Hide()
	}

	w.canvas.SetContent(content)
	// show new canvas element
	if content != nil {
		content.Show()
	}
	w.RescaleContext()
}

func (w *window) Canvas() gui.Canvas {
	return w.canvas
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(false)

	if w.onCloseIntercepted != nil {
		w.QueueEvent(w.onCloseIntercepted)
		return
	}

	w.Close()
}

// destroy this window and, if it's the last window quit the app
func (w *window) destroy(d *gLDriver) {
	w.DestroyEventQueue()
	cache.CleanCanvas(w.canvas)

	if w.master {
		d.Quit()
	} else if runtime.GOOS == "darwin" {
		go d.focusPreviousWindow()
	}
}

func (w *window) moved(_ *glfw.Window, x, y int) {
	if !w.fullScreen { // don't save the move to top left when changing to fullscreen
		// save coordinates
		w.xpos, w.ypos = x, y
	}

	if w.canvas.detectedScale == w.detectScale() {
		return
	}

	w.canvas.detectedScale = w.detectScale()
	go w.canvas.reloadScale()
}

func (w *window) resized(_ *glfw.Window, width, height int) {
	canvasSize := w.computeCanvasSize(width, height)
	if !w.fullScreen {
		w.width = internal.ScaleInt(w.canvas, canvasSize.Width)
		w.height = internal.ScaleInt(w.canvas, canvasSize.Height)
	}

	if !w.visible { // don't redraw if hidden
		w.canvas.Resize(canvasSize)
		return
	}

	if w.fixedSize {
		w.canvas.Resize(canvasSize)
		w.fitContent()
		return
	}

	w.platformResize(canvasSize)
}

func (w *window) frameSized(viewport *glfw.Window, width, height int) {
	if width == 0 || height == 0 || runtime.GOOS != "darwin" {
		return
	}

	winWidth, _ := viewport.GetSize()
	newTexScale := float32(width) / float32(winWidth) // This will be > 1.0 on a HiDPI screen
	w.canvas.RLock()
	texScale := w.canvas.texScale
	w.canvas.RUnlock()
	if texScale != newTexScale {
		w.canvas.Lock()
		w.canvas.texScale = newTexScale
		w.canvas.Unlock()
		w.canvas.Refresh(w.canvas.Content()) // reset graphics to apply texture scale
	}
}

func (w *window) refresh(_ *glfw.Window) {
	refreshWindow(w)
}

func (w *window) findObjectAtPositionMatching(canvas *glCanvas, mouse gui.Position, matches func(object gui.CanvasObject) bool) (gui.CanvasObject, gui.Position, int) {
	return driver.FindObjectAtPositionMatching(mouse, matches, canvas.Overlays().Top(), canvas.menu, canvas.Content())
}

func guiToNativeCursor(cursor desktop.Cursor) (*glfw.Cursor, bool) {
	switch v := cursor.(type) {
	case desktop.StandardCursor:
		ret, ok := cursorMap[v]
		if !ok {
			return cursorMap[desktop.DefaultCursor], false
		}
		return ret, false
	default:
		img, x, y := cursor.Image()
		if img == nil {
			return nil, true
		}
		return glfw.CreateCursor(img, x, y), true
	}
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	w.mouseLock.Lock()
	previousPos := w.mousePos
	w.mousePos = gui.NewPos(internal.UnscaleInt(w.canvas, int(xpos)), internal.UnscaleInt(w.canvas, int(ypos)))
	mousePos := w.mousePos
	mouseButton := w.mouseButton
	mouseDragPos := w.mouseDragPos
	mouseOver := w.mouseOver
	w.mouseLock.Unlock()

	cursor := desktop.Cursor(desktop.DefaultCursor)

	obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, mousePos, func(object gui.CanvasObject) bool {
		if cursorable, ok := object.(desktop.Cursorable); ok {
			cursor = cursorable.Cursor()
		}

		_, hover := object.(desktop.Hoverable)
		return hover
	})

	if w.cursor != cursor {
		// cursor has changed, store new cursor and apply change via glfw
		rawCursor, isCustomCursor := guiToNativeCursor(cursor)
		w.cursor = cursor

		if rawCursor == nil {
			viewport.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		} else {
			viewport.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
			viewport.SetCursor(rawCursor)
		}
		if w.customCursor != nil {
			w.customCursor.Destroy()
			w.customCursor = nil
		}
		if isCustomCursor {
			w.customCursor = rawCursor
		}
	}

	if w.mouseButton != 0 && w.mouseButton != desktop.MouseButtonSecondary && !w.mouseDragStarted {
		obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, previousPos, func(object gui.CanvasObject) bool {
			_, ok := object.(gui.Draggable)
			return ok
		})

		deltaX := mousePos.X - mouseDragPos.X
		deltaY := mousePos.Y - mouseDragPos.Y
		overThreshold := math.Abs(float64(deltaX)) >= dragMoveThreshold || math.Abs(float64(deltaY)) >= dragMoveThreshold

		if wid, ok := obj.(gui.Draggable); ok && overThreshold {
			w.mouseLock.Lock()
			w.mouseDragged = wid
			w.mouseDraggedOffset = previousPos.Subtract(pos)
			w.mouseDraggedObjStart = obj.Position()
			w.mouseDragStarted = true
			w.mouseLock.Unlock()
		}
	}

	w.mouseLock.RLock()
	isObjDragged := w.objIsDragged(obj)
	isMouseOverDragged := w.objIsDragged(mouseOver)
	w.mouseLock.RUnlock()
	if obj != nil && !isObjDragged {
		ev := new(desktop.MouseEvent)
		ev.AbsolutePosition = mousePos
		ev.Position = pos
		ev.Button = mouseButton

		if hovered, ok := obj.(desktop.Hoverable); ok {
			if hovered == mouseOver {
				w.QueueEvent(func() { hovered.MouseMoved(ev) })
			} else {
				w.mouseOut()
				w.mouseIn(hovered, ev)
			}
		} else if mouseOver != nil {
			isChild := false
			driver.WalkCompleteObjectTree(mouseOver.(gui.CanvasObject),
				func(co gui.CanvasObject, p1, p2 gui.Position, s gui.Size) bool {
					if co == obj {
						isChild = true
						return true
					}
					return false
				}, nil)
			if !isChild {
				w.mouseOut()
			}
		}
	} else if mouseOver != nil && !isMouseOverDragged {
		w.mouseOut()
	}

	w.mouseLock.RLock()
	mouseButton = w.mouseButton
	mouseDragged := w.mouseDragged
	mouseDraggedObjStart := w.mouseDraggedObjStart
	mouseDraggedOffset := w.mouseDraggedOffset
	mouseDragPos = w.mouseDragPos
	w.mouseLock.RUnlock()
	if mouseDragged != nil && mouseButton != desktop.MouseButtonSecondary {
		if w.mouseButton > 0 {
			draggedObjDelta := mouseDraggedObjStart.Subtract(mouseDragged.(gui.CanvasObject).Position())
			ev := new(gui.DragEvent)
			ev.AbsolutePosition = mousePos
			ev.Position = mousePos.Subtract(mouseDraggedOffset).Add(draggedObjDelta)
			ev.Dragged = gui.NewDelta(mousePos.X-mouseDragPos.X, mousePos.Y-mouseDragPos.Y)
			wd := mouseDragged
			w.QueueEvent(func() { wd.Dragged(ev) })
		}

		w.mouseLock.Lock()
		w.mouseDragStarted = true
		w.mouseDragPos = mousePos
		w.mouseLock.Unlock()
	}
}

func (w *window) objIsDragged(obj interface{}) bool {
	if w.mouseDragged != nil && obj != nil {
		draggedObj, _ := obj.(gui.Draggable)
		return draggedObj == w.mouseDragged
	}
	return false
}

func (w *window) mouseIn(obj desktop.Hoverable, ev *desktop.MouseEvent) {
	w.QueueEvent(func() {
		if obj != nil {
			obj.MouseIn(ev)
		}
		w.mouseLock.Lock()
		w.mouseOver = obj
		w.mouseLock.Unlock()
	})
}

func (w *window) mouseOut() {
	w.QueueEvent(func() {
		w.mouseLock.RLock()
		mouseOver := w.mouseOver
		w.mouseLock.RUnlock()
		if mouseOver != nil {
			mouseOver.MouseOut()
			w.mouseLock.Lock()
			w.mouseOver = nil
			w.mouseLock.Unlock()
		}
	})
}

func (w *window) mouseClicked(_ *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	w.mouseLock.RLock()
	w.mouseDragPos = w.mousePos
	mousePos := w.mousePos
	mouseDragStarted := w.mouseDragStarted
	w.mouseLock.RUnlock()
	if mousePos.IsZero() { // window may not be focused (darwin mostly) and so position callbacks not happening
		xpos, ypos := w.viewport.GetCursorPos()
		w.mouseLock.Lock()
		w.mousePos = gui.NewPos(internal.UnscaleInt(w.canvas, int(xpos)), internal.UnscaleInt(w.canvas, int(ypos)))
		mousePos = w.mousePos
		w.mouseLock.Unlock()
	}

	co, pos, _ := w.findObjectAtPositionMatching(w.canvas, mousePos, func(object gui.CanvasObject) bool {
		switch object.(type) {
		case gui.Tappable, gui.SecondaryTappable, gui.DoubleTappable, gui.Focusable, desktop.Mouseable, desktop.Hoverable:
			return true
		case gui.Draggable:
			if mouseDragStarted {
				return true
			}
		}

		return false
	})
	ev := new(gui.PointEvent)
	ev.Position = pos
	ev.AbsolutePosition = mousePos

	coMouse := co
	button, modifiers := convertMouseButton(btn, mods)
	if wid, ok := co.(desktop.Mouseable); ok {
		mev := new(desktop.MouseEvent)
		mev.Position = ev.Position
		mev.AbsolutePosition = mousePos
		mev.Button = button
		mev.Modifier = modifiers
		w.mouseClickedHandleMouseable(mev, action, wid)
	}

	if wid, ok := co.(gui.Focusable); !ok || wid != w.canvas.Focused() {
		w.canvas.Unfocus()
	}

	w.mouseLock.Lock()
	if action == glfw.Press {
		w.mouseButton |= button
	} else if action == glfw.Release {
		w.mouseButton &= ^button
	}

	mouseDragged := w.mouseDragged
	mouseDragStarted = w.mouseDragStarted
	mouseOver := w.mouseOver
	shouldMouseOut := w.objIsDragged(mouseOver) && !w.objIsDragged(coMouse)
	mousePressed := w.mousePressed
	w.mouseLock.Unlock()

	if action == glfw.Release && mouseDragged != nil {
		if mouseDragStarted {
			w.QueueEvent(mouseDragged.DragEnd)
			w.mouseLock.Lock()
			w.mouseDragStarted = false
			w.mouseLock.Unlock()
		}
		if shouldMouseOut {
			w.mouseOut()
		}
		w.mouseLock.Lock()
		w.mouseDragged = nil
		w.mouseLock.Unlock()
	}

	_, tap := co.(gui.Tappable)
	_, altTap := co.(gui.SecondaryTappable)
	if tap || altTap {
		if action == glfw.Press {
			w.mouseLock.Lock()
			w.mousePressed = co
			w.mouseLock.Unlock()
		} else if action == glfw.Release {
			if co == mousePressed {
				if button == desktop.MouseButtonSecondary && altTap {
					w.QueueEvent(func() { co.(gui.SecondaryTappable).TappedSecondary(ev) })
				}
			}
		}
	}

	// Check for double click/tap on left mouse button
	if action == glfw.Release && button == desktop.MouseButtonPrimary && !mouseDragStarted {
		w.mouseClickedHandleTapDoubleTap(co, ev)
	}
}

func (w *window) mouseClickedHandleMouseable(mev *desktop.MouseEvent, action glfw.Action, wid desktop.Mouseable) {
	mousePos := mev.AbsolutePosition
	if action == glfw.Press {
		w.QueueEvent(func() { wid.MouseDown(mev) })
	} else if action == glfw.Release {
		w.mouseLock.RLock()
		mouseDragged := w.mouseDragged
		mouseDraggedOffset := w.mouseDraggedOffset
		w.mouseLock.RUnlock()
		if mouseDragged == nil {
			w.QueueEvent(func() { wid.MouseUp(mev) })
		} else {
			if dragged, ok := mouseDragged.(desktop.Mouseable); ok {
				mev.Position = mousePos.Subtract(mouseDraggedOffset)
				w.QueueEvent(func() { dragged.MouseUp(mev) })
			} else {
				w.QueueEvent(func() { wid.MouseUp(mev) })
			}
		}
	}
}

func (w *window) mouseClickedHandleTapDoubleTap(co gui.CanvasObject, ev *gui.PointEvent) {
	_, doubleTap := co.(gui.DoubleTappable)
	if doubleTap {
		w.mouseLock.Lock()
		w.mouseClickCount++
		w.mouseLastClick = co
		mouseCancelFunc := w.mouseCancelFunc
		w.mouseLock.Unlock()
		if mouseCancelFunc != nil {
			mouseCancelFunc()
			return
		}
		go w.waitForDoubleTap(co, ev)
	} else {
		w.mouseLock.Lock()
		if wid, ok := co.(gui.Tappable); ok && co == w.mousePressed {
			w.QueueEvent(func() { wid.Tapped(ev) })
		}
		w.mousePressed = nil
		w.mouseLock.Unlock()
	}
}

func (w *window) waitForDoubleTap(co gui.CanvasObject, ev *gui.PointEvent) {
	var ctx context.Context
	w.mouseLock.Lock()
	ctx, w.mouseCancelFunc = context.WithDeadline(context.TODO(), time.Now().Add(time.Millisecond*doubleClickDelay))
	defer w.mouseCancelFunc()
	w.mouseLock.Unlock()

	<-ctx.Done()

	w.mouseLock.Lock()
	defer w.mouseLock.Unlock()

	if w.mouseClickCount == 2 && w.mouseLastClick == co {
		if wid, ok := co.(gui.DoubleTappable); ok {
			w.QueueEvent(func() { wid.DoubleTapped(ev) })
		}
	} else if co == w.mousePressed {
		if wid, ok := co.(gui.Tappable); ok {
			w.QueueEvent(func() { wid.Tapped(ev) })
		}
	}

	w.mouseClickCount = 0
	w.mousePressed = nil
	w.mouseCancelFunc = nil
	w.mouseLastClick = nil
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff float64, yoff float64) {
	w.mouseLock.RLock()
	mousePos := w.mousePos
	w.mouseLock.RUnlock()
	co, pos, _ := w.findObjectAtPositionMatching(w.canvas, mousePos, func(object gui.CanvasObject) bool {
		_, ok := object.(gui.Scrollable)
		return ok
	})
	switch wid := co.(type) {
	case gui.Scrollable:
		if runtime.GOOS != "darwin" && xoff == 0 &&
			(viewport.GetKey(glfw.KeyLeftShift) == glfw.Press ||
				viewport.GetKey(glfw.KeyRightShift) == glfw.Press) {
			xoff, yoff = yoff, xoff
		}
		if math.Abs(xoff) >= scrollAccelerateCutoff {
			xoff *= scrollAccelerateRate
		}
		if math.Abs(yoff) >= scrollAccelerateCutoff {
			yoff *= scrollAccelerateRate
		}

		ev := &gui.ScrollEvent{}
		ev.Scrolled = gui.NewDelta(float32(xoff)*scrollSpeed, float32(yoff)*scrollSpeed)
		ev.Position = pos
		ev.AbsolutePosition = mousePos
		wid.Scrolled(ev)
	}
}

func convertMouseButton(btn glfw.MouseButton, mods glfw.ModifierKey) (desktop.MouseButton, desktop.Modifier) {
	modifier := desktopModifier(mods)
	var button desktop.MouseButton
	rightClick := false
	if runtime.GOOS == "darwin" {
		if modifier&desktop.ControlModifier != 0 {
			rightClick = true
			modifier &^= desktop.ControlModifier
		}
		if modifier&desktop.SuperModifier != 0 {
			modifier |= desktop.ControlModifier
			modifier &^= desktop.SuperModifier
		}
	}
	switch btn {
	case glfw.MouseButton1:
		if rightClick {
			button = desktop.MouseButtonSecondary
		} else {
			button = desktop.MouseButtonPrimary
		}
	case glfw.MouseButton2:
		button = desktop.MouseButtonSecondary
	case glfw.MouseButton3:
		button = desktop.MouseButtonTertiary
	}
	return button, modifier
}

var keyCodeMap = map[glfw.Key]gui.KeyName{
	// non-printable
	glfw.KeyEscape:    gui.KeyEscape,
	glfw.KeyEnter:     gui.KeyReturn,
	glfw.KeyTab:       gui.KeyTab,
	glfw.KeyBackspace: gui.KeyBackspace,
	glfw.KeyInsert:    gui.KeyInsert,
	glfw.KeyDelete:    gui.KeyDelete,
	glfw.KeyRight:     gui.KeyRight,
	glfw.KeyLeft:      gui.KeyLeft,
	glfw.KeyDown:      gui.KeyDown,
	glfw.KeyUp:        gui.KeyUp,
	glfw.KeyPageUp:    gui.KeyPageUp,
	glfw.KeyPageDown:  gui.KeyPageDown,
	glfw.KeyHome:      gui.KeyHome,
	glfw.KeyEnd:       gui.KeyEnd,

	glfw.KeySpace:   gui.KeySpace,
	glfw.KeyKPEnter: gui.KeyEnter,

	// functions
	glfw.KeyF1:  gui.KeyF1,
	glfw.KeyF2:  gui.KeyF2,
	glfw.KeyF3:  gui.KeyF3,
	glfw.KeyF4:  gui.KeyF4,
	glfw.KeyF5:  gui.KeyF5,
	glfw.KeyF6:  gui.KeyF6,
	glfw.KeyF7:  gui.KeyF7,
	glfw.KeyF8:  gui.KeyF8,
	glfw.KeyF9:  gui.KeyF9,
	glfw.KeyF10: gui.KeyF10,
	glfw.KeyF11: gui.KeyF11,
	glfw.KeyF12: gui.KeyF12,

	// numbers - lookup by code to avoid AZERTY using the symbol name instead of number
	glfw.Key0:   gui.Key0,
	glfw.KeyKP0: gui.Key0,
	glfw.Key1:   gui.Key1,
	glfw.KeyKP1: gui.Key1,
	glfw.Key2:   gui.Key2,
	glfw.KeyKP2: gui.Key2,
	glfw.Key3:   gui.Key3,
	glfw.KeyKP3: gui.Key3,
	glfw.Key4:   gui.Key4,
	glfw.KeyKP4: gui.Key4,
	glfw.Key5:   gui.Key5,
	glfw.KeyKP5: gui.Key5,
	glfw.Key6:   gui.Key6,
	glfw.KeyKP6: gui.Key6,
	glfw.Key7:   gui.Key7,
	glfw.KeyKP7: gui.Key7,
	glfw.Key8:   gui.Key8,
	glfw.KeyKP8: gui.Key8,
	glfw.Key9:   gui.Key9,
	glfw.KeyKP9: gui.Key9,

	// desktop
	glfw.KeyLeftShift:    desktop.KeyShiftLeft,
	glfw.KeyRightShift:   desktop.KeyShiftRight,
	glfw.KeyLeftControl:  desktop.KeyControlLeft,
	glfw.KeyRightControl: desktop.KeyControlRight,
	glfw.KeyLeftAlt:      desktop.KeyAltLeft,
	glfw.KeyRightAlt:     desktop.KeyAltRight,
	glfw.KeyLeftSuper:    desktop.KeySuperLeft,
	glfw.KeyRightSuper:   desktop.KeySuperRight,
	glfw.KeyMenu:         desktop.KeyMenu,
	glfw.KeyPrintScreen:  desktop.KeyPrintScreen,
	glfw.KeyCapsLock:     desktop.KeyCapsLock,
}

var keyNameMap = map[string]gui.KeyName{
	"'": gui.KeyApostrophe,
	",": gui.KeyComma,
	"-": gui.KeyMinus,
	".": gui.KeyPeriod,
	"/": gui.KeySlash,
	"*": gui.KeyAsterisk,
	"`": gui.KeyBackTick,

	";": gui.KeySemicolon,
	"+": gui.KeyPlus,
	"=": gui.KeyEqual,

	"a": gui.KeyA,
	"b": gui.KeyB,
	"c": gui.KeyC,
	"d": gui.KeyD,
	"e": gui.KeyE,
	"f": gui.KeyF,
	"g": gui.KeyG,
	"h": gui.KeyH,
	"i": gui.KeyI,
	"j": gui.KeyJ,
	"k": gui.KeyK,
	"l": gui.KeyL,
	"m": gui.KeyM,
	"n": gui.KeyN,
	"o": gui.KeyO,
	"p": gui.KeyP,
	"q": gui.KeyQ,
	"r": gui.KeyR,
	"s": gui.KeyS,
	"t": gui.KeyT,
	"u": gui.KeyU,
	"v": gui.KeyV,
	"w": gui.KeyW,
	"x": gui.KeyX,
	"y": gui.KeyY,
	"z": gui.KeyZ,

	"[":  gui.KeyLeftBracket,
	"\\": gui.KeyBackslash,
	"]":  gui.KeyRightBracket,
}

func keyToName(code glfw.Key, scancode int) gui.KeyName {
	if runtime.GOOS == "darwin" && scancode == 0x69 { // TODO remove once fixed upstream
		code = glfw.KeyPrintScreen
	}

	ret, ok := keyCodeMap[code]
	if ok {
		return ret
	}

	keyName := glfw.GetKeyName(code, scancode)
	ret, ok = keyNameMap[keyName]
	if !ok {
		return gui.KeyUnknown
	}

	return ret
}

func (w *window) capturesTab(modifier desktop.Modifier) bool {
	captures := false

	if ent, ok := w.canvas.Focused().(gui.Tabbable); ok {
		captures = ent.AcceptsTab()
	}
	if !captures {
		switch modifier {
		case 0:
			w.QueueEvent(w.canvas.FocusNext)
			return false
		case desktop.ShiftModifier:
			w.QueueEvent(w.canvas.FocusPrevious)
			return false
		}
	}

	return captures
}

func (w *window) keyPressed(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key, scancode)
	keyEvent := &gui.KeyEvent{Name: keyName, Physical: gui.HardwareKey{ScanCode: scancode}}

	keyDesktopModifier := desktopModifier(mods)
	pendingMenuToggle := w.menuTogglePending
	pendingMenuDeactivation := w.menuDeactivationPending
	w.menuTogglePending = desktop.KeyNone
	w.menuDeactivationPending = desktop.KeyNone
	switch action {
	case glfw.Release:
		if action == glfw.Release && keyName != "" {
			switch keyName {
			case pendingMenuToggle:
				w.canvas.ToggleMenu()
			case pendingMenuDeactivation:
				if w.canvas.DismissMenu() {
					return
				}
			}
		}

		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.QueueEvent(func() { focused.KeyUp(keyEvent) })
			}
		} else if w.canvas.onKeyUp != nil {
			w.QueueEvent(func() { w.canvas.onKeyUp(keyEvent) })
		}
		return // ignore key up in other core events
	case glfw.Press:
		switch keyName {
		case desktop.KeyAltLeft, desktop.KeyAltRight:
			// compensate for GLFW modifiers bug
			if (runtime.GOOS == "linux" && keyDesktopModifier == 0) || (runtime.GOOS != "linux" && keyDesktopModifier == desktop.AltModifier) {
				w.menuTogglePending = keyName
			}
		case gui.KeyEscape:
			w.menuDeactivationPending = keyName
		}
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.QueueEvent(func() { focused.KeyDown(keyEvent) })
			}
		} else if w.canvas.onKeyDown != nil {
			w.QueueEvent(func() { w.canvas.onKeyDown(keyEvent) })
		}
	default:
		// key repeat will fall through to TypedKey and TypedShortcut
	}

	if (keyName == gui.KeyTab && !w.capturesTab(keyDesktopModifier)) || w.triggersShortcut(keyName, keyDesktopModifier) {
		return
	}

	// No shortcut detected, pass down to TypedKey
	focused := w.canvas.Focused()
	if focused != nil {
		w.QueueEvent(func() { focused.TypedKey(keyEvent) })
	} else if w.canvas.onTypedKey != nil {
		w.QueueEvent(func() { w.canvas.onTypedKey(keyEvent) })
	}
}

func desktopModifier(mods glfw.ModifierKey) desktop.Modifier {
	var m desktop.Modifier
	if (mods & glfw.ModShift) != 0 {
		m |= desktop.ShiftModifier
	}
	if (mods & glfw.ModControl) != 0 {
		m |= desktop.ControlModifier
	}
	if (mods & glfw.ModAlt) != 0 {
		m |= desktop.AltModifier
	}
	if (mods & glfw.ModSuper) != 0 {
		m |= desktop.SuperModifier
	}
	return m
}

// charInput defines the character with modifiers callback which is called when a
// Unicode character is input.
//
// Characters do not map 1:1 to physical keys, as a key may produce zero, one or more characters.
func (w *window) charInput(_ *glfw.Window, char rune) {
	if focused := w.canvas.Focused(); focused != nil {
		w.QueueEvent(func() { focused.TypedRune(char) })
	} else if w.canvas.onTypedRune != nil {
		w.QueueEvent(func() { w.canvas.onTypedRune(char) })
	}
}

func (w *window) focused(_ *glfw.Window, focus bool) {
	if focus {
		if curWindow == nil {
			gui.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerEnteredForeground()
		}
		curWindow = w
		w.canvas.FocusGained()
	} else {
		w.canvas.FocusLost()
		w.mouseLock.Lock()
		w.mousePos = gui.Position{}
		w.mouseLock.Unlock()

		go func() { // check whether another window was focused or not
			time.Sleep(time.Millisecond * 100)
			if curWindow != w {
				return
			}

			curWindow = nil
			gui.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerExitedForeground()
		}()
	}
}

func (w *window) triggersShortcut(keyName gui.KeyName, modifier desktop.Modifier) bool {
	var shortcut gui.Shortcut
	ctrlMod := desktop.ControlModifier
	if runtime.GOOS == "darwin" {
		ctrlMod = desktop.SuperModifier
	}
	if modifier == ctrlMod {
		switch keyName {
		case gui.KeyV:
			// detect paste shortcut
			shortcut = &gui.ShortcutPaste{
				Clipboard: w.Clipboard(),
			}
		case gui.KeyC, gui.KeyInsert:
			// detect copy shortcut
			shortcut = &gui.ShortcutCopy{
				Clipboard: w.Clipboard(),
			}
		case gui.KeyX:
			// detect cut shortcut
			shortcut = &gui.ShortcutCut{
				Clipboard: w.Clipboard(),
			}
		case gui.KeyA:
			// detect selectAll shortcut
			shortcut = &gui.ShortcutSelectAll{}
		}
	}

	if modifier == desktop.ShiftModifier {
		switch keyName {
		case gui.KeyInsert:
			// detect paste shortcut
			shortcut = &gui.ShortcutPaste{
				Clipboard: w.Clipboard(),
			}
		case gui.KeyDelete:
			// detect cut shortcut
			shortcut = &gui.ShortcutCut{
				Clipboard: w.Clipboard(),
			}
		}
	}

	if shortcut == nil && modifier != 0 && !isKeyModifier(keyName) && modifier != desktop.ShiftModifier {
		shortcut = &desktop.CustomShortcut{
			KeyName:  keyName,
			Modifier: modifier,
		}
	}

	if shortcut != nil {
		if focused, ok := w.canvas.Focused().(gui.Shortcutable); ok {
			shouldRunShortcut := true
			type selectableText interface {
				gui.Disableable
				SelectedText() string
			}
			if selectableTextWid, ok := focused.(selectableText); ok && selectableTextWid.Disabled() {
				shouldRunShortcut = shortcut.ShortcutName() == "Copy"
			}
			if shouldRunShortcut {
				w.QueueEvent(func() { focused.TypedShortcut(shortcut) })
			}
			return shouldRunShortcut
		}
		w.QueueEvent(func() { w.canvas.TypedShortcut(shortcut) })
		return true
	}

	return false
}

func (w *window) RunWithContext(f func()) {
	if w.isClosing() {
		return
	}
	w.viewport.MakeContextCurrent()

	f()

	glfw.DetachCurrentContext()
}

func (w *window) RescaleContext() {
	runOnMain(func() {
		w.rescaleOnMain()
	})
}

func (w *window) rescaleOnMain() {
	if w.isClosing() {
		return
	}
	w.fitContent()

	if w.fullScreen {
		w.width, w.height = w.viewport.GetSize()
		scaledFull := gui.NewSize(
			internal.UnscaleInt(w.canvas, w.width),
			internal.UnscaleInt(w.canvas, w.height))
		w.canvas.Resize(scaledFull)
		return
	}

	size := w.canvas.size.Max(w.canvas.MinSize())
	newWidth, newHeight := w.screenSize(size)
	w.viewport.SetSize(newWidth, newHeight)
}

func (w *window) Context() interface{} {
	return nil
}

func (w *window) runOnMainWhenCreated(fn func()) {
	if w.viewport != nil {
		runOnMain(fn)
		return
	}

	w.pending = append(w.pending, fn)
}

func (d *gLDriver) CreateWindow(title string) gui.Window {
	return d.createWindow(title, true)
}

func (d *gLDriver) createWindow(title string, decorate bool) gui.Window {
	var ret *window
	if title == "" {
		title = defaultTitle
	}
	runOnMain(func() {
		d.initGLFW()

		ret = &window{title: title, decorate: decorate, driver: d}
		// This queue is destroyed when the window is closed.
		ret.InitEventQueue()
		go ret.RunEventQueue()

		ret.canvas = newCanvas()
		ret.canvas.context = ret
		ret.SetIcon(ret.icon)
		d.addWindow(ret)
	})
	return ret
}

func (w *window) create() {
	runOnMain(func() {
		if !isWayland {
			// make the window hidden, we will set it up and then show it later
			glfw.WindowHint(glfw.Visible, glfw.False)
		}
		if w.decorate {
			glfw.WindowHint(glfw.Decorated, glfw.True)
		} else {
			glfw.WindowHint(glfw.Decorated, glfw.False)
		}
		if w.fixedSize {
			glfw.WindowHint(glfw.Resizable, glfw.False)
		} else {
			glfw.WindowHint(glfw.Resizable, glfw.True)
		}
		glfw.WindowHint(glfw.AutoIconify, glfw.False)
		initWindowHints()

		pixWidth, pixHeight := w.screenSize(w.canvas.size)
		pixWidth = int(gui.Max(float32(pixWidth), float32(w.width)))
		if pixWidth == 0 {
			pixWidth = 10
		}
		pixHeight = int(gui.Max(float32(pixHeight), float32(w.height)))
		if pixHeight == 0 {
			pixHeight = 10
		}

		win, err := glfw.CreateWindow(pixWidth, pixHeight, w.title, nil, nil)
		if err != nil {
			w.driver.initFailed("window creation error", err)
			return
		}

		w.viewLock.Lock()
		w.viewport = win
		w.viewLock.Unlock()
	})
	if w.view() == nil { // something went wrong above, it will have been logged
		return
	}

	// run the GL init on the draw thread
	runOnDraw(w, func() {
		w.canvas.SetPainter(gl.NewPainter(w.canvas, w))
		w.canvas.Painter().Init()
	})

	runOnMain(func() {
		w.setDarkMode()

		win := w.view()
		win.SetCloseCallback(w.closed)
		win.SetPosCallback(w.moved)
		win.SetSizeCallback(w.resized)
		win.SetFramebufferSizeCallback(w.frameSized)
		win.SetRefreshCallback(w.refresh)
		win.SetCursorPosCallback(w.mouseMoved)
		win.SetMouseButtonCallback(w.mouseClicked)
		win.SetScrollCallback(w.mouseScrolled)
		win.SetKeyCallback(w.keyPressed)
		win.SetCharCallback(w.charInput)
		win.SetFocusCallback(w.focused)

		w.canvas.detectedScale = w.detectScale()
		w.canvas.scale = w.calculatedScale()
		w.canvas.texScale = w.detectTextureScale()
		// update window size now we have scaled detected
		w.fitContent()

		for _, fn := range w.pending {
			fn()
		}

		w.requestedWidth, w.requestedHeight = w.width, w.height

		if w.fixedSize { // as the window will not be sized later we may need to pack menus etc
			w.canvas.Resize(w.canvas.Size())
		}
		// order of operation matters so we do these last items in order
		w.viewport.SetSize(w.shouldWidth, w.shouldHeight) // ensure we requested latest size
	})
}

func (w *window) doShowAgain() {
	if w.isClosing() {
		return
	}

	runOnMain(func() {
		// show top canvas element
		if w.canvas.Content() != nil {
			w.canvas.Content().Show()
		}

		w.viewport.SetPos(w.xpos, w.ypos)
		w.viewport.Show()
		w.viewLock.Lock()
		w.visible = true
		w.viewLock.Unlock()
	})
}

func (w *window) isClosing() bool {
	w.viewLock.RLock()
	closing := w.closing || w.viewport == nil
	w.viewLock.RUnlock()
	return closing
}

func (w *window) view() *glfw.Window {
	w.viewLock.RLock()
	defer w.viewLock.RUnlock()

	if w.closing {
		return nil
	}
	return w.viewport
}

func (d *gLDriver) CreateSplashWindow() gui.Window {
	win := d.createWindow("", false)
	win.SetPadded(false)
	win.CenterOnScreen()
	return win
}

func (d *gLDriver) AllWindows() []gui.Window {
	return d.windows
}

func isKeyModifier(keyName gui.KeyName) bool {
	return keyName == desktop.KeyShiftLeft || keyName == desktop.KeyShiftRight ||
		keyName == desktop.KeyControlLeft || keyName == desktop.KeyControlRight ||
		keyName == desktop.KeyAltLeft || keyName == desktop.KeyAltRight ||
		keyName == desktop.KeySuperLeft || keyName == desktop.KeySuperRight
}
