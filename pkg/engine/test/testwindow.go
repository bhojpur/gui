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
	gui "github.com/bhojpur/gui/pkg/engine"
)

type testWindow struct {
	title              string
	fullScreen         bool
	fixedSize          bool
	focused            bool
	onClosed           func()
	onCloseIntercepted func()

	canvas    *testCanvas
	clipboard gui.Clipboard
	driver    *testDriver
	menu      *gui.MainMenu
}

// NewWindow creates and registers a new window for test purposes
func NewWindow(content gui.CanvasObject) gui.Window {
	window := gui.CurrentApp().NewWindow("")
	window.SetContent(content)
	return window
}

func (w *testWindow) Canvas() gui.Canvas {
	return w.canvas
}

func (w *testWindow) CenterOnScreen() {
	// no-op
}

func (w *testWindow) Clipboard() gui.Clipboard {
	return w.clipboard
}

func (w *testWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

func (w *testWindow) Content() gui.CanvasObject {
	return w.Canvas().Content()
}

func (w *testWindow) FixedSize() bool {
	return w.fixedSize
}

func (w *testWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *testWindow) Hide() {
	w.focused = false
}

func (w *testWindow) Icon() gui.Resource {
	return gui.CurrentApp().Icon()
}

func (w *testWindow) MainMenu() *gui.MainMenu {
	return w.menu
}

func (w *testWindow) Padded() bool {
	return w.canvas.Padded()
}

func (w *testWindow) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*testWindow).focused = false
	}

	w.focused = true
}

func (w *testWindow) Resize(size gui.Size) {
	w.canvas.Resize(size)
}

func (w *testWindow) SetContent(obj gui.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *testWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *testWindow) SetIcon(_ gui.Resource) {
	// no-op
}

func (w *testWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *testWindow) SetMainMenu(menu *gui.MainMenu) {
	w.menu = menu
}

func (w *testWindow) SetMaster() {
	// no-op
}

func (w *testWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *testWindow) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *testWindow) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

func (w *testWindow) SetTitle(title string) {
	w.title = title
}

func (w *testWindow) Show() {
	w.RequestFocus()
}

func (w *testWindow) ShowAndRun() {
	w.Show()
}

func (w *testWindow) Title() string {
	return w.title
}
