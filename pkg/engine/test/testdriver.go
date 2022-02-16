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
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/painter"
	"github.com/bhojpur/gui/pkg/engine/internal/painter/software"
	intRepo "github.com/bhojpur/gui/pkg/engine/internal/repository"
	"github.com/bhojpur/gui/pkg/engine/storage/repository"
)

// SoftwarePainter describes a simple type that can render canvases
type SoftwarePainter interface {
	Paint(gui.Canvas) image.Image
}

type testDriver struct {
	device       *device
	painter      SoftwarePainter
	windows      []gui.Window
	windowsMutex sync.RWMutex
}

// Declare conformity with Driver
var _ gui.Driver = (*testDriver)(nil)

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver() gui.Driver {
	drv := new(testDriver)
	drv.windowsMutex = sync.RWMutex{}
	repository.Register("file", intRepo.NewFileRepository())

	// make a single dummy window for rendering tests
	drv.CreateWindow("")

	return drv
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter SoftwarePainter) gui.Driver {
	drv := new(testDriver)
	drv.painter = painter
	drv.windowsMutex = sync.RWMutex{}

	return drv
}

func (d *testDriver) AbsolutePositionForObject(co gui.CanvasObject) gui.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return gui.NewPos(0, 0)
	}

	tc := c.(*testCanvas)
	return driver.AbsolutePositionForObject(co, tc.objectTrees())
}

func (d *testDriver) AllWindows() []gui.Window {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	return d.windows
}

func (d *testDriver) CanvasForObject(gui.CanvasObject) gui.Canvas {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	// cheating: probably the last created window is meant
	return d.windows[len(d.windows)-1].Canvas()
}

func (d *testDriver) CreateWindow(string) gui.Window {
	canvas := NewCanvas().(*testCanvas)
	if d.painter != nil {
		canvas.painter = d.painter
	} else {
		canvas.painter = software.NewPainter()
	}

	window := &testWindow{canvas: canvas, driver: d}
	window.clipboard = &testClipboard{}

	d.windowsMutex.Lock()
	d.windows = append(d.windows, window)
	d.windowsMutex.Unlock()
	return window
}

func (d *testDriver) Device() gui.Device {
	if d.device == nil {
		d.device = &device{}
	}
	return d.device
}

// RenderedTextSize looks up how bit a string would be if drawn on screen
func (d *testDriver) RenderedTextSize(text string, size float32, style gui.TextStyle) (gui.Size, float32) {
	return painter.RenderedTextSize(text, size, style)
}

func (d *testDriver) Run() {
	// no-op
}

func (d *testDriver) StartAnimation(a *gui.Animation) {
	// currently no animations in test app, we just initialise it and leave
	a.Tick(1.0)
}

func (d *testDriver) StopAnimation(a *gui.Animation) {
	// currently no animations in test app, do nothing
}

func (d *testDriver) Quit() {
	// no-op
}

func (d *testDriver) removeWindow(w *testWindow) {
	d.windowsMutex.Lock()
	i := 0
	for _, window := range d.windows {
		if window == w {
			break
		}
		i++
	}

	d.windows = append(d.windows[:i], d.windows[i+1:]...)
	d.windowsMutex.Unlock()
}
