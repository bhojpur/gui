//go:build linux || darwin || windows
// +build linux darwin windows

package app

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
	"github.com/bhojpur/gui/pkg/engine/internal/async"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/lifecycle"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/size"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/gl"

	// Initialize necessary mobile functionality, such as logging.
	_ "github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/mobileinit"
)

// Main is called by the main.main function to run the mobile application.
//
// It calls f on the App, in a separate goroutine, as some OS-specific
// libraries require being on 'the main thread'.
func Main(f func(App)) {
	main(f)
}

// App is how a Bhojpur GUI mobile application interacts with the OS.
type App interface {
	// Events returns the events channel. It carries events from the system to
	// the app. The type of such events include:
	//  - lifecycle.Event
	//  - mouse.Event
	//  - paint.Event
	//  - size.Event
	//  - touch.Event
	// from the golang.org/x/mobile/event/etc packages. Other packages may
	// define other event types that are carried on this channel.
	Events() <-chan interface{}

	// Send sends an event on the events channel. It does not block.
	Send(event interface{})

	// Publish flushes any pending drawing commands, such as OpenGL calls, and
	// swaps the back buffer to the screen.
	Publish() PublishResult

	// TODO: replace filters (and the Events channel) with a NextEvent method?

	// Filter calls each registered event filter function in sequence.
	Filter(event interface{}) interface{}

	// RegisterFilter registers a event filter function to be called by Filter. The
	// function can return a different event, or return nil to consume the event,
	// but the function can also return its argument unchanged, where its purpose
	// is to trigger a side effect rather than modify the event.
	RegisterFilter(f func(interface{}) interface{})

	ShowVirtualKeyboard(KeyboardType)
	HideVirtualKeyboard()
	ShowFileOpenPicker(func(string, func()), *FileFilter)
	ShowFileSavePicker(func(string, func()), *FileFilter, string)
}

// FileFilter is a filter of files.
type FileFilter struct {
	Extensions []string
	MimeTypes  []string
}

// PublishResult is the result of an App.Publish call.
type PublishResult struct {
	// BackBufferPreserved is whether the contents of the back buffer was
	// preserved. If false, the contents are undefined.
	BackBufferPreserved bool
}

var theApp = &app{
	events:         async.NewUnboundedInterfaceChan(),
	lifecycleStage: lifecycle.StageDead,
	publish:        make(chan struct{}),
	publishResult:  make(chan PublishResult),
}

func init() {
	theApp.glctx, theApp.worker = gl.NewContext()
}

func (a *app) sendLifecycle(to lifecycle.Stage) {
	if a.lifecycleStage == to {
		return
	}
	a.events.In() <- lifecycle.Event{
		From:        a.lifecycleStage,
		To:          to,
		DrawContext: a.glctx,
	}
	a.lifecycleStage = to
}

type app struct {
	filters []func(interface{}) interface{}

	events         *async.UnboundedInterfaceChan
	lifecycleStage lifecycle.Stage
	publish        chan struct{}
	publishResult  chan PublishResult

	glctx  gl.Context
	worker gl.Worker
}

func (a *app) Events() <-chan interface{} {
	return a.events.Out()
}

func (a *app) Send(event interface{}) {
	a.events.In() <- event
}

func (a *app) Publish() PublishResult {
	// gl.Flush is a lightweight (on modern GL drivers) blocking call
	// that ensures all GL functions pending in the gl package have
	// been passed onto the GL driver before the app package attempts
	// to swap the screen buffer.
	//
	// This enforces that the final receive (for this paint cycle) on
	// gl.WorkAvailable happens before the send on endPaint.
	a.glctx.Flush()
	a.publish <- struct{}{}
	return <-a.publishResult
}

func (a *app) Filter(event interface{}) interface{} {
	for _, f := range a.filters {
		event = f(event)
	}
	return event
}

func (a *app) RegisterFilter(f func(interface{}) interface{}) {
	a.filters = append(a.filters, f)
}

func (a *app) ShowVirtualKeyboard(keyboard KeyboardType) {
	driverShowVirtualKeyboard(keyboard)
}

func (a *app) HideVirtualKeyboard() {
	driverHideVirtualKeyboard()
}

func (a *app) ShowFileOpenPicker(callback func(string, func()), filter *FileFilter) {
	driverShowFileOpenPicker(callback, filter)
}
func (a *app) ShowFileSavePicker(callback func(string, func()), filter *FileFilter, filename string) {
	driverShowFileSavePicker(callback, filter, filename)
}

// TODO: do this for all build targets, not just linux (x11 and Android)? If
// so, should package gl instead of this package call RegisterFilter??
//
// TODO: does Android need this?? It seems to work without it (Nexus 7,
// KitKat). If only x11 needs this, should we move this to x11.go??
func (a *app) registerGLViewportFilter() {
	a.RegisterFilter(func(e interface{}) interface{} {
		if e, ok := e.(size.Event); ok {
			a.glctx.Viewport(0, 0, e.WidthPx, e.HeightPx)
		}
		return e
	})
}

func screenOrientation(width, height int) size.Orientation {
	if width > height {
		return size.OrientationLandscape
	}

	return size.OrientationPortrait
}
