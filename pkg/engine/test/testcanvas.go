package test

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
	"image/draw"
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var (
	dummyCanvas gui.Canvas
)

// WindowlessCanvas provides functionality for a canvas to operate without a window
type WindowlessCanvas interface {
	gui.Canvas

	Padded() bool
	Resize(gui.Size)
	SetPadded(bool)
	SetScale(float32)
}

type testCanvas struct {
	size  gui.Size
	scale float32

	content  gui.CanvasObject
	overlays *internal.OverlayStack
	focusMgr *app.FocusManager
	hovered  desktop.Hoverable
	padded   bool

	onTypedRune func(rune)
	onTypedKey  func(*gui.KeyEvent)

	gui.ShortcutHandler
	painter      SoftwarePainter
	propertyLock sync.RWMutex
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() gui.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvas()
	}

	return dummyCanvas
}

// NewCanvas returns a single use in-memory canvas used for testing
func NewCanvas() WindowlessCanvas {
	c := &testCanvas{
		focusMgr: app.NewFocusManager(nil),
		padded:   true,
		scale:    1.0,
		size:     gui.NewSize(10, 10),
	}
	c.overlays = &internal.OverlayStack{Canvas: c}
	return c
}

// NewCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	canvas := NewCanvas().(*testCanvas)
	canvas.painter = painter

	return canvas
}

func (c *testCanvas) Capture() image.Image {
	bounds := image.Rect(0, 0, internal.ScaleInt(c, c.Size().Width), internal.ScaleInt(c, c.Size().Height))
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(theme.BackgroundColor()), image.Point{}, draw.Src)

	if c.painter != nil {
		draw.Draw(img, bounds, c.painter.Paint(c), image.Point{}, draw.Over)
	}

	return img
}

func (c *testCanvas) Content() gui.CanvasObject {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.content
}

func (c *testCanvas) Focus(obj gui.Focusable) {
	c.focusManager().Focus(obj)
}

func (c *testCanvas) FocusNext() {
	c.focusManager().FocusNext()
}

func (c *testCanvas) FocusPrevious() {
	c.focusManager().FocusPrevious()
}

func (c *testCanvas) Focused() gui.Focusable {
	return c.focusManager().Focused()
}

func (c *testCanvas) InteractiveArea() (gui.Position, gui.Size) {
	return gui.Position{}, c.Size()
}

func (c *testCanvas) OnTypedKey() func(*gui.KeyEvent) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedKey
}

func (c *testCanvas) OnTypedRune() func(rune) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedRune
}

func (c *testCanvas) Overlays() gui.OverlayStack {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.overlays
}

func (c *testCanvas) Padded() bool {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.padded
}

func (c *testCanvas) PixelCoordinateForPosition(pos gui.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *testCanvas) Refresh(gui.CanvasObject) {
}

func (c *testCanvas) Resize(size gui.Size) {
	c.propertyLock.Lock()
	content := c.content
	overlays := c.overlays
	padded := c.padded
	c.size = size
	c.propertyLock.Unlock()

	if content == nil {
		return
	}

	// Ensure testcanvas mimics real canvas.Resize behavior
	for _, overlay := range overlays.List() {
		type popupWidget interface {
			gui.CanvasObject
			ShowAtPosition(gui.Position)
		}
		if p, ok := overlay.(popupWidget); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(size)
		}
	}

	if padded {
		content.Resize(size.Subtract(gui.NewSize(theme.Padding()*2, theme.Padding()*2)))
		content.Move(gui.NewPos(theme.Padding(), theme.Padding()))
	} else {
		content.Resize(size)
		content.Move(gui.NewPos(0, 0))
	}
}

func (c *testCanvas) Scale() float32 {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.scale
}

func (c *testCanvas) SetContent(content gui.CanvasObject) {
	c.propertyLock.Lock()
	c.content = content
	c.focusMgr = app.NewFocusManager(c.content)
	c.propertyLock.Unlock()

	if content == nil {
		return
	}

	padding := gui.NewSize(0, 0)
	if c.padded {
		padding = gui.NewSize(theme.Padding()*2, theme.Padding()*2)
	}
	c.Resize(content.MinSize().Add(padding))
}

func (c *testCanvas) SetOnTypedKey(handler func(*gui.KeyEvent)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedKey = handler
}

func (c *testCanvas) SetOnTypedRune(handler func(rune)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedRune = handler
}

func (c *testCanvas) SetPadded(padded bool) {
	c.propertyLock.Lock()
	c.padded = padded
	c.propertyLock.Unlock()

	c.Resize(c.Size())
}

func (c *testCanvas) SetScale(scale float32) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.scale = scale
}

func (c *testCanvas) Size() gui.Size {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.size
}

func (c *testCanvas) Unfocus() {
	c.focusManager().Focus(nil)
}

func (c *testCanvas) focusManager() *app.FocusManager {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	return c.focusMgr
}

func (c *testCanvas) objectTrees() []gui.CanvasObject {
	trees := make([]gui.CanvasObject, 0, len(c.Overlays().List())+1)
	if c.content != nil {
		trees = append(trees, c.content)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

func layoutAndCollect(objects []gui.CanvasObject, o gui.CanvasObject, size gui.Size) []gui.CanvasObject {
	objects = append(objects, o)
	switch c := o.(type) {
	case gui.Widget:
		r := c.CreateRenderer()
		r.Layout(size)
		for _, child := range r.Objects() {
			objects = layoutAndCollect(objects, child, child.Size())
		}
	case *gui.Container:
		if c.Layout != nil {
			c.Layout.Layout(c.Objects, size)
		}
		for _, child := range c.Objects {
			objects = layoutAndCollect(objects, child, child.Size())
		}
	}
	return objects
}
