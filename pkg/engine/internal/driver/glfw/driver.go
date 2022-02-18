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

// It provides a full Bhojpur GUI desktop application driver that
// uses the system OpenGL libraries. This supports Windows, Mac OS X and Linux
// using the gl and glfw packages from go-gl.

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/animation"
	intapp "github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/internal/painter"
	intRepo "github.com/bhojpur/gui/pkg/engine/internal/repository"
	"github.com/bhojpur/gui/pkg/engine/storage/repository"
)

const mainGoroutineID = 1

var (
	curWindow *window
	isWayland = false
)

// Declare conformity with Driver
var _ gui.Driver = (*gLDriver)(nil)

type gLDriver struct {
	windowLock sync.RWMutex
	windows    []gui.Window
	device     *glDevice
	done       chan interface{}
	drawDone   chan interface{}

	animation *animation.Runner

	drawOnMainThread bool // A workaround on Apple M1, just use 1 thread until fixed upstream
}

func (d *gLDriver) RenderedTextSize(text string, textSize float32, style gui.TextStyle) (size gui.Size, baseline float32) {
	return painter.RenderedTextSize(text, textSize, style)
}

func (d *gLDriver) CanvasForObject(obj gui.CanvasObject) gui.Canvas {
	return common.CanvasForObject(obj)
}

func (d *gLDriver) AbsolutePositionForObject(co gui.CanvasObject) gui.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return gui.NewPos(0, 0)
	}

	glc := c.(*glCanvas)
	return driver.AbsolutePositionForObject(co, glc.ObjectTrees())
}

func (d *gLDriver) Device() gui.Device {
	if d.device == nil {
		d.device = &glDevice{}
	}

	return d.device
}

func (d *gLDriver) Quit() {
	if curWindow != nil {
		curWindow = nil
		gui.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerExitedForeground()
	}
	defer func() {
		recover() // we could be called twice - no safe way to check if d.done is closed
	}()
	close(d.done)
}

func (d *gLDriver) addWindow(w *window) {
	d.windowLock.Lock()
	defer d.windowLock.Unlock()
	d.windows = append(d.windows, w)
}

// a trivial implementation of "focus previous" - return to the most recently opened, or master if set.
// This may not do the right thing if your app has 3 or more windows open, but it was agreed this was not much
// of an issue, and the added complexity to track focus was not needed at this time.
func (d *gLDriver) focusPreviousWindow() {
	d.windowLock.RLock()
	wins := d.windows
	d.windowLock.RUnlock()

	var chosen gui.Window
	for _, w := range wins {
		chosen = w
		if w.(*window).master {
			break
		}
	}

	if chosen == nil || chosen.(*window).view() == nil {
		return
	}
	chosen.RequestFocus()
}

func (d *gLDriver) windowList() []gui.Window {
	d.windowLock.RLock()
	defer d.windowLock.RUnlock()
	return d.windows
}

func (d *gLDriver) initFailed(msg string, err error) {
	gui.LogError(msg, err)

	run.Lock()
	if !run.flag {
		run.Unlock()
		d.Quit()
	} else {
		run.Unlock()
		os.Exit(1)
	}
}

func goroutineID() int {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// string format expects "goroutine X [running..."
	id := strings.Split(strings.TrimSpace(string(b)), " ")[1]

	num, _ := strconv.Atoi(id)
	return num
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() gui.Driver {
	d := new(gLDriver)
	d.done = make(chan interface{})
	d.drawDone = make(chan interface{})
	d.animation = &animation.Runner{}

	repository.Register("file", intRepo.NewFileRepository())

	return d
}
