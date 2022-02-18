package mobile

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
	"runtime"
	"strconv"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/animation"
	intapp "github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/app"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/key"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/lifecycle"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/paint"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/size"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/touch"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/gl"
	"github.com/bhojpur/gui/pkg/engine/internal/painter"
	pgl "github.com/bhojpur/gui/pkg/engine/internal/painter/gl"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

const (
	tapMoveThreshold  = 4.0                    // how far can we move before it is a drag
	tapSecondaryDelay = 300 * time.Millisecond // how long before secondary tap
)

// Configuration is the system information about the current device
type Configuration struct {
	SystemTheme gui.ThemeVariant
}

// ConfiguredDriver is a simple type that allows packages to hook into configuration changes of this driver.
type ConfiguredDriver interface {
	SetOnConfigurationChanged(func(*Configuration))
}

type mobileDriver struct {
	app   app.App
	glctx gl.Context

	windows     []gui.Window
	device      *device
	animation   *animation.Runner
	currentSize size.Event

	theme           gui.ThemeVariant
	onConfigChanged func(*Configuration)
	painting        bool
}

// Declare conformity with Driver
var _ gui.Driver = (*mobileDriver)(nil)
var _ ConfiguredDriver = (*mobileDriver)(nil)

func init() {
	runtime.LockOSThread()
}

func (d *mobileDriver) CreateWindow(title string) gui.Window {
	c := NewCanvas().(*mobileCanvas) // silence lint
	ret := &window{title: title, canvas: c, isChild: len(d.windows) > 0}
	ret.InitEventQueue()
	go ret.RunEventQueue()
	c.setContent(&canvas.Rectangle{FillColor: theme.BackgroundColor()})
	c.SetPainter(pgl.NewPainter(c, ret))
	d.windows = append(d.windows, ret)
	return ret
}

func (d *mobileDriver) AllWindows() []gui.Window {
	return d.windows
}

// currentWindow returns the most recently opened window - we can only show one at a time.
func (d *mobileDriver) currentWindow() *window {
	if len(d.windows) == 0 {
		return nil
	}

	var last *window
	for i := len(d.windows) - 1; i >= 0; i-- {
		last = d.windows[i].(*window)
		if last.visible {
			return last
		}
	}

	return last
}

func (d *mobileDriver) RenderedTextSize(text string, textSize float32, style gui.TextStyle) (size gui.Size, baseline float32) {
	return painter.RenderedTextSize(text, textSize, style)
}

func (d *mobileDriver) CanvasForObject(obj gui.CanvasObject) gui.Canvas {
	if len(d.windows) == 0 {
		return nil
	}

	// TODO figure out how we handle multiple windows...
	return d.currentWindow().Canvas()
}

func (d *mobileDriver) AbsolutePositionForObject(co gui.CanvasObject) gui.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return gui.NewPos(0, 0)
	}

	mc := c.(*mobileCanvas)
	pos := driver.AbsolutePositionForObject(co, mc.ObjectTrees())
	inset, _ := c.InteractiveArea()

	if mc.windowHead != nil {
		if len(mc.windowHead.(*gui.Container).Objects) > 1 {
			topHeight := mc.windowHead.MinSize().Height
			pos = pos.Subtract(gui.NewSize(0, topHeight))
		}
	}
	return pos.Subtract(inset)
}

func (d *mobileDriver) Quit() {
	// Android and iOS guidelines say this should not be allowed!
}

func (d *mobileDriver) Run() {
	app.Main(func(a app.App) {
		d.app = a
		settingsChange := make(chan gui.Settings)
		gui.CurrentApp().Settings().AddChangeListener(settingsChange)
		draw := time.NewTicker(time.Second / 60)

		for {
			select {
			case <-draw.C:
				d.sendPaintEvent()
			case set := <-settingsChange:
				painter.ClearFontCache()
				cache.ResetThemeCaches()
				intapp.ApplySettingsWithCallback(set, gui.CurrentApp(), func(w gui.Window) {
					c, ok := w.Canvas().(*mobileCanvas)
					if !ok {
						return
					}
					c.applyThemeOutOfTreeObjects()
				})
			case e, ok := <-a.Events():
				if !ok {
					return // events channel closed, app done
				}
				current := d.currentWindow()
				if current == nil {
					continue
				}
				c := current.Canvas().(*mobileCanvas)

				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					d.handleLifecycle(e, current)
				case size.Event:
					if e.WidthPx <= 0 {
						continue
					}
					d.currentSize = e
					currentOrientation = e.Orientation
					currentDPI = e.PixelsPerPt * 72
					d.setTheme(e.DarkMode)

					dev := d.device
					dev.safeTop = e.InsetTopPx
					dev.safeLeft = e.InsetLeftPx
					dev.safeHeight = e.HeightPx - e.InsetTopPx - e.InsetBottomPx
					dev.safeWidth = e.WidthPx - e.InsetLeftPx - e.InsetRightPx
					c.scale = gui.CurrentDevice().SystemScaleForWindow(nil)
					c.Painter().SetFrameBufferScale(1.0)

					// make sure that we paint on the next frame
					c.Content().Refresh()
				case paint.Event:
					d.handlePaint(e, current)
				case touch.Event:
					switch e.Type {
					case touch.TypeBegin:
						d.tapDownCanvas(current, e.X, e.Y, e.Sequence)
					case touch.TypeMove:
						d.tapMoveCanvas(current, e.X, e.Y, e.Sequence)
					case touch.TypeEnd:
						d.tapUpCanvas(current, e.X, e.Y, e.Sequence)
					}
				case key.Event:
					if e.Direction == key.DirPress {
						d.typeDownCanvas(c, e.Rune, e.Code, e.Modifiers)
					} else if e.Direction == key.DirRelease {
						d.typeUpCanvas(c, e.Rune, e.Code, e.Modifiers)
					}
				}
			}
		}
	})
}

func (d *mobileDriver) handleLifecycle(e lifecycle.Event, w gui.Window) {
	c := w.Canvas().(*mobileCanvas)
	switch e.Crosses(lifecycle.StageVisible) {
	case lifecycle.CrossOn:
		d.glctx, _ = e.DrawContext.(gl.Context)
		d.onStart()

		// this is a fix for some android phone to prevent the app from being drawn as a blank screen after being pushed in the background
		c.Content().Refresh()

		d.sendPaintEvent()
	case lifecycle.CrossOff:
		d.onStop()
		d.glctx = nil
	}
	switch e.Crosses(lifecycle.StageFocused) {
	case lifecycle.CrossOn: // foregrounding
		gui.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerEnteredForeground()
	case lifecycle.CrossOff: // will enter background
		if runtime.GOOS == "darwin" {
			if d.glctx == nil {
				return
			}

			s := gui.NewSize(float32(d.currentSize.WidthPx)/c.scale, float32(d.currentSize.HeightPx)/c.scale)
			d.paintWindow(w, s)
			d.app.Publish()
		}
		gui.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerExitedForeground()
	}
}

func (d *mobileDriver) handlePaint(e paint.Event, w gui.Window) {
	c := w.Canvas().(*mobileCanvas)
	d.painting = false
	if d.glctx == nil || e.External {
		return
	}
	if !c.inited {
		c.inited = true
		c.Painter().Init() // we cannot init until the context is set above
	}

	canvasNeedRefresh := c.FreeDirtyTextures() > 0 || c.CheckDirtyAndClear()
	if canvasNeedRefresh {
		newSize := gui.NewSize(float32(d.currentSize.WidthPx)/c.scale, float32(d.currentSize.HeightPx)/c.scale)

		if c.EnsureMinSize() {
			c.sizeContent(newSize) // force resize of content
		} else { // if screen changed
			w.Resize(newSize)
		}

		d.paintWindow(w, newSize)
		d.app.Publish()
	}
	cache.Clean(canvasNeedRefresh)
}

func (d *mobileDriver) onStart() {
	gui.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerStarted()
}

func (d *mobileDriver) onStop() {
	gui.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerStopped()
}

func (d *mobileDriver) paintWindow(window gui.Window, size gui.Size) {
	clips := &internal.ClipStack{}
	c := window.Canvas().(*mobileCanvas)

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	d.glctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	d.glctx.Clear(gl.ColorBufferBit)

	draw := func(node *common.RenderCacheNode, pos gui.Position) {
		obj := node.Obj()
		if _, ok := obj.(gui.Scrollable); ok {
			inner := clips.Push(pos, obj.Size())
			c.Painter().StartClipping(inner.Rect())
		}
		c.Painter().Paint(obj, pos, size)
	}
	afterDraw := func(node *common.RenderCacheNode) {
		if _, ok := node.Obj().(gui.Scrollable); ok {
			c.Painter().StopClipping()
			clips.Pop()
			if top := clips.Top(); top != nil {
				c.Painter().StartClipping(top.Rect())
			}
		}
	}

	c.WalkTrees(draw, afterDraw)
}

func (d *mobileDriver) sendPaintEvent() {
	if d.painting {
		return
	}
	d.app.Send(paint.Event{})
	d.painting = true
}

func (d *mobileDriver) setTheme(dark bool) {
	var mode gui.ThemeVariant
	if dark {
		mode = theme.VariantDark
	} else {
		mode = theme.VariantLight
	}

	if d.theme != mode && d.onConfigChanged != nil {
		d.onConfigChanged(&Configuration{SystemTheme: mode})
	}
	d.theme = mode
}

func (d *mobileDriver) tapDownCanvas(w *window, x, y float32, tapID touch.Sequence) {
	tapX := internal.UnscaleInt(w.canvas, int(x))
	tapY := internal.UnscaleInt(w.canvas, int(y))
	pos := gui.NewPos(tapX, tapY+tapYOffset)

	w.canvas.tapDown(pos, int(tapID))
}

func (d *mobileDriver) tapMoveCanvas(w *window, x, y float32, tapID touch.Sequence) {
	tapX := internal.UnscaleInt(w.canvas, int(x))
	tapY := internal.UnscaleInt(w.canvas, int(y))
	pos := gui.NewPos(tapX, tapY+tapYOffset)

	w.canvas.tapMove(pos, int(tapID), func(wid gui.Draggable, ev *gui.DragEvent) {
		w.QueueEvent(func() { wid.Dragged(ev) })
	})
}

func (d *mobileDriver) tapUpCanvas(w *window, x, y float32, tapID touch.Sequence) {
	tapX := internal.UnscaleInt(w.canvas, int(x))
	tapY := internal.UnscaleInt(w.canvas, int(y))
	pos := gui.NewPos(tapX, tapY+tapYOffset)

	w.canvas.tapUp(pos, int(tapID), func(wid gui.Tappable, ev *gui.PointEvent) {
		w.QueueEvent(func() { wid.Tapped(ev) })
	}, func(wid gui.SecondaryTappable, ev *gui.PointEvent) {
		w.QueueEvent(func() { wid.TappedSecondary(ev) })
	}, func(wid gui.DoubleTappable, ev *gui.PointEvent) {
		w.QueueEvent(func() { wid.DoubleTapped(ev) })
	}, func(wid gui.Draggable) {
		w.QueueEvent(wid.DragEnd)
	})
}

var keyCodeMap = map[key.Code]gui.KeyName{
	// non-printable
	key.CodeEscape:          gui.KeyEscape,
	key.CodeReturnEnter:     gui.KeyReturn,
	key.CodeTab:             gui.KeyTab,
	key.CodeDeleteBackspace: gui.KeyBackspace,
	key.CodeInsert:          gui.KeyInsert,
	key.CodePageUp:          gui.KeyPageUp,
	key.CodePageDown:        gui.KeyPageDown,
	key.CodeHome:            gui.KeyHome,
	key.CodeEnd:             gui.KeyEnd,

	key.CodeF1:  gui.KeyF1,
	key.CodeF2:  gui.KeyF2,
	key.CodeF3:  gui.KeyF3,
	key.CodeF4:  gui.KeyF4,
	key.CodeF5:  gui.KeyF5,
	key.CodeF6:  gui.KeyF6,
	key.CodeF7:  gui.KeyF7,
	key.CodeF8:  gui.KeyF8,
	key.CodeF9:  gui.KeyF9,
	key.CodeF10: gui.KeyF10,
	key.CodeF11: gui.KeyF11,
	key.CodeF12: gui.KeyF12,

	key.CodeKeypadEnter: gui.KeyEnter,

	// printable
	key.CodeA:       gui.KeyA,
	key.CodeB:       gui.KeyB,
	key.CodeC:       gui.KeyC,
	key.CodeD:       gui.KeyD,
	key.CodeE:       gui.KeyE,
	key.CodeF:       gui.KeyF,
	key.CodeG:       gui.KeyG,
	key.CodeH:       gui.KeyH,
	key.CodeI:       gui.KeyI,
	key.CodeJ:       gui.KeyJ,
	key.CodeK:       gui.KeyK,
	key.CodeL:       gui.KeyL,
	key.CodeM:       gui.KeyM,
	key.CodeN:       gui.KeyN,
	key.CodeO:       gui.KeyO,
	key.CodeP:       gui.KeyP,
	key.CodeQ:       gui.KeyQ,
	key.CodeR:       gui.KeyR,
	key.CodeS:       gui.KeyS,
	key.CodeT:       gui.KeyT,
	key.CodeU:       gui.KeyU,
	key.CodeV:       gui.KeyV,
	key.CodeW:       gui.KeyW,
	key.CodeX:       gui.KeyX,
	key.CodeY:       gui.KeyY,
	key.CodeZ:       gui.KeyZ,
	key.Code0:       gui.Key0,
	key.CodeKeypad0: gui.Key0,
	key.Code1:       gui.Key1,
	key.CodeKeypad1: gui.Key1,
	key.Code2:       gui.Key2,
	key.CodeKeypad2: gui.Key2,
	key.Code3:       gui.Key3,
	key.CodeKeypad3: gui.Key3,
	key.Code4:       gui.Key4,
	key.CodeKeypad4: gui.Key4,
	key.Code5:       gui.Key5,
	key.CodeKeypad5: gui.Key5,
	key.Code6:       gui.Key6,
	key.CodeKeypad6: gui.Key6,
	key.Code7:       gui.Key7,
	key.CodeKeypad7: gui.Key7,
	key.Code8:       gui.Key8,
	key.CodeKeypad8: gui.Key8,
	key.Code9:       gui.Key9,
	key.CodeKeypad9: gui.Key9,

	key.CodeSemicolon: gui.KeySemicolon,
	key.CodeEqualSign: gui.KeyEqual,

	key.CodeSpacebar:           gui.KeySpace,
	key.CodeApostrophe:         gui.KeyApostrophe,
	key.CodeComma:              gui.KeyComma,
	key.CodeHyphenMinus:        gui.KeyMinus,
	key.CodeKeypadHyphenMinus:  gui.KeyMinus,
	key.CodeFullStop:           gui.KeyPeriod,
	key.CodeKeypadFullStop:     gui.KeyPeriod,
	key.CodeSlash:              gui.KeySlash,
	key.CodeLeftSquareBracket:  gui.KeyLeftBracket,
	key.CodeBackslash:          gui.KeyBackslash,
	key.CodeRightSquareBracket: gui.KeyRightBracket,
	key.CodeGraveAccent:        gui.KeyBackTick,
}

func keyToName(code key.Code) gui.KeyName {
	ret, ok := keyCodeMap[code]
	if !ok {
		return ""
	}

	return ret
}

func runeToPrintable(r rune) rune {
	if strconv.IsPrint(r) {
		return r
	}

	return 0
}

func (d *mobileDriver) typeDownCanvas(canvas *mobileCanvas, r rune, code key.Code, mod key.Modifiers) {
	keyName := keyToName(code)
	switch keyName {
	case gui.KeyTab:
		capture := false
		if ent, ok := canvas.Focused().(gui.Tabbable); ok {
			capture = ent.AcceptsTab()
		}
		if !capture {
			switch mod {
			case 0:
				canvas.FocusNext()
				return
			case key.ModShift:
				canvas.FocusPrevious()
				return
			}
		}
	}

	r = runeToPrintable(r)
	keyEvent := &gui.KeyEvent{Name: keyName}

	if canvas.Focused() != nil {
		if keyName != "" {
			canvas.Focused().TypedKey(keyEvent)
		}
		if r > 0 {
			canvas.Focused().TypedRune(r)
		}
	} else if canvas.onTypedKey != nil {
		if keyName != "" {
			canvas.onTypedKey(keyEvent)
		}
		if r > 0 {
			canvas.onTypedRune(r)
		}
	}
}

func (d *mobileDriver) typeUpCanvas(_ *mobileCanvas, _ rune, _ key.Code, _ key.Modifiers) {
}

func (d *mobileDriver) Device() gui.Device {
	if d.device == nil {
		d.device = &device{}
	}

	return d.device
}

func (d *mobileDriver) SetOnConfigurationChanged(f func(*Configuration)) {
	d.onConfigChanged = f
}

// NewGoMobileDriver sets up a new Driver instance implemented using the Go
// Mobile extension and OpenGL bindings.
func NewGoMobileDriver() gui.Driver {
	d := new(mobileDriver)
	d.theme = gui.ThemeVariant(2) // unspecified
	d.animation = &animation.Runner{}

	registerRepository(d)
	return d
}
