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
	"image"
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// Declare conformity with Canvas interface
var _ gui.Canvas = (*glCanvas)(nil)

type glCanvas struct {
	common.Canvas

	content gui.CanvasObject
	menu    gui.CanvasObject
	padded  bool
	size    gui.Size

	onTypedRune func(rune)
	onTypedKey  func(*gui.KeyEvent)
	onKeyDown   func(*gui.KeyEvent)
	onKeyUp     func(*gui.KeyEvent)
	// shortcut    gui.ShortcutHandler

	scale, detectedScale, texScale float32

	context driver.WithContext
}

func (c *glCanvas) Capture() image.Image {
	var img image.Image
	runOnDraw(c.context.(*window), func() {
		img = c.Painter().Capture(c)
	})
	return img
}

func (c *glCanvas) Content() gui.CanvasObject {
	c.RLock()
	retval := c.content
	c.RUnlock()
	return retval
}

func (c *glCanvas) DismissMenu() bool {
	c.RLock()
	menu := c.menu
	c.RUnlock()
	if menu != nil && menu.(*MenuBar).IsActive() {
		menu.(*MenuBar).Toggle()
		return true
	}
	return false
}

func (c *glCanvas) InteractiveArea() (gui.Position, gui.Size) {
	return gui.Position{}, c.Size()
}

func (c *glCanvas) MinSize() gui.Size {
	c.RLock()
	defer c.RUnlock()
	return c.canvasSize(c.content.MinSize())
}

func (c *glCanvas) OnKeyDown() func(*gui.KeyEvent) {
	return c.onKeyDown
}

func (c *glCanvas) OnKeyUp() func(*gui.KeyEvent) {
	return c.onKeyUp
}

func (c *glCanvas) OnTypedKey() func(*gui.KeyEvent) {
	return c.onTypedKey
}

func (c *glCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *glCanvas) Padded() bool {
	return c.padded
}

func (c *glCanvas) PixelCoordinateForPosition(pos gui.Position) (int, int) {
	c.RLock()
	texScale := c.texScale
	c.RUnlock()
	multiple := c.Scale() * texScale
	scaleInt := func(x float32) int {
		return int(math.Round(float64(x * multiple)))
	}

	return scaleInt(pos.X), scaleInt(pos.Y)
}

func (c *glCanvas) Resize(size gui.Size) {
	c.Lock()
	c.size = size
	c.Unlock()

	for _, overlay := range c.Overlays().List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(size)
		}
	}

	c.RLock()
	c.content.Resize(c.contentSize(size))
	c.content.Move(c.contentPos())

	if c.menu != nil {
		c.menu.Refresh()
		c.menu.Resize(gui.NewSize(size.Width, c.menu.MinSize().Height))
	}
	c.RUnlock()
}

func (c *glCanvas) Scale() float32 {
	c.RLock()
	defer c.RUnlock()
	return c.scale
}

func (c *glCanvas) SetContent(content gui.CanvasObject) {
	c.Lock()
	c.setContent(content)

	c.content.Resize(c.content.MinSize()) // give it the space it wants then calculate the real min
	// the pass above makes some layouts wide enough to wrap, so we ask again what the true min is.
	newSize := c.size.Max(c.canvasSize(c.content.MinSize()))
	c.Unlock()

	c.Resize(newSize)
	c.SetDirty()
}

func (c *glCanvas) SetOnKeyDown(typed func(*gui.KeyEvent)) {
	c.onKeyDown = typed
}

func (c *glCanvas) SetOnKeyUp(typed func(*gui.KeyEvent)) {
	c.onKeyUp = typed
}

func (c *glCanvas) SetOnTypedKey(typed func(*gui.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *glCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *glCanvas) SetPadded(padded bool) {
	c.Lock()
	content := c.content
	c.padded = padded
	pos := c.contentPos()
	c.Unlock()

	content.Move(pos)
}

func (c *glCanvas) reloadScale() {
	w := c.context.(*window)
	w.viewLock.RLock()
	windowVisible := w.visible
	w.viewLock.RUnlock()
	if !windowVisible {
		return
	}

	c.Lock()
	c.scale = c.context.(*window).calculatedScale()
	c.Unlock()
	c.SetDirty()

	c.context.RescaleContext()
}

func (c *glCanvas) Size() gui.Size {
	c.RLock()
	defer c.RUnlock()
	return c.size
}

func (c *glCanvas) ToggleMenu() {
	c.RLock()
	menu := c.menu
	c.RUnlock()
	if menu != nil {
		menu.(*MenuBar).Toggle()
	}
}

func (c *glCanvas) buildMenu(w *window, m *gui.MainMenu) {
	c.Lock()
	defer c.Unlock()
	c.setMenuOverlay(nil)
	if m == nil {
		return
	}
	if hasNativeMenu() {
		setupNativeMenu(w, m)
	} else {
		c.setMenuOverlay(buildMenuOverlay(m, w))
	}
}

// canvasSize computes the needed canvas size for the given content size
func (c *glCanvas) canvasSize(contentSize gui.Size) gui.Size {
	canvasSize := contentSize.Add(gui.NewSize(0, c.menuHeight()))
	if c.Padded() {
		pad := theme.Padding() * 2
		canvasSize = canvasSize.Add(gui.NewSize(pad, pad))
	}
	return canvasSize
}

func (c *glCanvas) contentPos() gui.Position {
	contentPos := gui.NewPos(0, c.menuHeight())
	if c.Padded() {
		contentPos = contentPos.Add(gui.NewPos(theme.Padding(), theme.Padding()))
	}
	return contentPos
}

func (c *glCanvas) contentSize(canvasSize gui.Size) gui.Size {
	contentSize := gui.NewSize(canvasSize.Width, canvasSize.Height-c.menuHeight())
	if c.Padded() {
		pad := theme.Padding() * 2
		contentSize = contentSize.Subtract(gui.NewSize(pad, pad))
	}
	return contentSize
}

func (c *glCanvas) menuHeight() float32 {
	switch c.menu {
	case nil:
		// no menu or native menu -> does not consume space on the canvas
		return 0
	default:
		return c.menu.MinSize().Height
	}
}

func (c *glCanvas) overlayChanged() {
	c.SetDirty()
}

func (c *glCanvas) paint(size gui.Size) {
	clips := &internal.ClipStack{}
	if c.Content() == nil {
		return
	}
	c.Painter().Clear()

	paint := func(node *common.RenderCacheNode, pos gui.Position) {
		obj := node.Obj()
		if _, ok := obj.(gui.Scrollable); ok {
			inner := clips.Push(pos, obj.Size())
			c.Painter().StartClipping(inner.Rect())
		}
		c.Painter().Paint(obj, pos, size)
	}
	afterPaint := func(node *common.RenderCacheNode) {
		if _, ok := node.Obj().(gui.Scrollable); ok {
			clips.Pop()
			if top := clips.Top(); top != nil {
				c.Painter().StartClipping(top.Rect())
			} else {
				c.Painter().StopClipping()

			}
		}
	}

	c.WalkTrees(paint, afterPaint)
}

func (c *glCanvas) setContent(content gui.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

func (c *glCanvas) setMenuOverlay(b gui.CanvasObject) {
	c.menu = b
	c.SetMenuTreeAndFocusMgr(b)

	if c.menu != nil && !c.size.IsZero() {
		c.content.Resize(c.contentSize(c.size))
		c.content.Move(c.contentPos())

		c.menu.Refresh()
		c.menu.Resize(gui.NewSize(c.size.Width, c.menu.MinSize().Height))
	}
}

func (c *glCanvas) applyThemeOutOfTreeObjects() {
	c.RLock()
	menu := c.menu
	padded := c.padded
	c.RUnlock()
	if menu != nil {
		app.ApplyThemeTo(menu, c) // Ensure our menu gets the theme change message as it's out-of-tree
	}

	c.SetPadded(padded) // refresh the padding for potential theme differences
}

func newCanvas() *glCanvas {
	c := &glCanvas{scale: 1.0, texScale: 1.0}
	c.Initialize(c, c.overlayChanged)
	c.setContent(&canvas.Rectangle{FillColor: theme.BackgroundColor()})
	c.padded = true
	return c
}
