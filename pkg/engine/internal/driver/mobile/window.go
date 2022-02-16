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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type window struct {
	common.Window

	title              string
	visible            bool
	onClosed           func()
	onCloseIntercepted func()
	isChild            bool

	clipboard gui.Clipboard
	canvas    *mobileCanvas
	icon      gui.Resource
	menu      *gui.MainMenu
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title
}

func (w *window) FullScreen() bool {
	return true
}

func (w *window) SetFullScreen(bool) {
	// no-op
}

func (w *window) Resize(size gui.Size) {
	w.Canvas().(*mobileCanvas).Resize(size)
}

func (w *window) RequestFocus() {
	// no-op - we cannot change which window is focused
}

func (w *window) FixedSize() bool {
	return true
}

func (w *window) SetFixedSize(bool) {
	// no-op - all windows are fixed size
}

func (w *window) CenterOnScreen() {
	// no-op
}

func (w *window) Padded() bool {
	return w.canvas.padded
}

func (w *window) SetPadded(padded bool) {
	w.canvas.padded = padded
}

func (w *window) Icon() gui.Resource {
	if w.icon == nil {
		return gui.CurrentApp().Icon()
	}

	return w.icon
}

func (w *window) SetIcon(icon gui.Resource) {
	w.icon = icon
}

func (w *window) SetMaster() {
	// no-op on mobile
}

func (w *window) MainMenu() *gui.MainMenu {
	return w.menu
}

func (w *window) SetMainMenu(menu *gui.MainMenu) {
	w.menu = menu
}

func (w *window) SetOnClosed(callback func()) {
	w.onClosed = callback
}

func (w *window) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *window) Show() {
	menu := gui.CurrentApp().Driver().(*mobileDriver).findMenu(w)
	menuButton := w.newMenuButton(menu)
	if menu == nil {
		menuButton.Hide()
	}

	if w.isChild {
		exit := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			w.tryClose()
		})
		title := widget.NewLabel(w.title)
		title.Alignment = gui.TextAlignCenter
		w.canvas.setWindowHead(container.NewHBox(menuButton,
			layout.NewSpacer(), title, layout.NewSpacer(), exit))
		w.canvas.Resize(w.canvas.size)
	} else {
		w.canvas.setWindowHead(container.NewHBox(menuButton))
	}
	w.visible = true

	if w.Content() != nil {
		w.Content().Refresh()
		w.Content().Show()
	}
}

func (w *window) Hide() {
	w.visible = false

	if w.Content() != nil {
		w.Content().Hide()
	}
}

func (w *window) tryClose() {
	if w.onCloseIntercepted != nil {
		w.QueueEvent(w.onCloseIntercepted)
		return
	}

	w.Close()
}

func (w *window) Close() {
	d := gui.CurrentApp().Driver().(*mobileDriver)
	pos := -1
	for i, win := range d.windows {
		if win == w {
			pos = i
		}
	}
	if pos != -1 {
		d.windows = append(d.windows[:pos], d.windows[pos+1:]...)
	}

	cache.RangeTexturesFor(w.canvas, func(obj gui.CanvasObject) {
		w.canvas.Painter().Free(obj)
	})

	w.canvas.WalkTrees(nil, func(node *common.RenderCacheNode) {
		if wid, ok := node.Obj().(gui.Widget); ok {
			cache.DestroyRenderer(wid)
		}
	})

	w.QueueEvent(func() {
		cache.CleanCanvas(w.canvas)
	})

	// Call this in a go routine, because this function could be called
	// inside a button which callback would be queued in this event queue
	// and it will lead to a deadlock if this is performed in the same go
	// routine.
	go w.DestroyEventQueue()

	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) ShowAndRun() {
	w.Show()
	gui.CurrentApp().Driver().Run()
}

func (w *window) Content() gui.CanvasObject {
	return w.canvas.Content()
}

func (w *window) SetContent(content gui.CanvasObject) {
	w.canvas.SetContent(content)
}

func (w *window) Canvas() gui.Canvas {
	return w.canvas
}

func (w *window) Clipboard() gui.Clipboard {
	if w.clipboard == nil {
		w.clipboard = &mobileClipboard{}
	}
	return w.clipboard
}

func (w *window) RunWithContext(f func()) {
	//	ctx, _ = e.DrawContext.(gl.Context)

	f()
}

func (w *window) RescaleContext() {
	// TODO
}

func (w *window) Context() interface{} {
	return gui.CurrentApp().Driver().(*mobileDriver).glctx
}
