package widget

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
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
)

// Base provides a helper that handles basic widget behaviours.
type Base struct {
	hidden   bool
	position gui.Position
	size     gui.Size

	impl         gui.Widget
	propertyLock sync.RWMutex
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *Base) ExtendBaseWidget(wid gui.Widget) {
	impl := w.super()
	if impl != nil {
		return
	}

	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.impl = wid
}

// Size gets the current size of this widget.
func (w *Base) Size() gui.Size {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return w.size
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *Base) Resize(size gui.Size) {
	w.propertyLock.RLock()
	baseSize := w.size
	impl := w.impl
	w.propertyLock.RUnlock()
	if baseSize == size {
		return
	}

	w.propertyLock.Lock()
	w.size = size
	w.propertyLock.Unlock()

	if impl == nil {
		return
	}
	cache.Renderer(impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *Base) Position() gui.Position {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return w.position
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *Base) Move(pos gui.Position) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()

	w.position = pos
}

// MinSize for the widget - it should never be resized below this value.
func (w *Base) MinSize() gui.Size {
	impl := w.super()

	r := cache.Renderer(impl)
	if r == nil {
		return gui.NewSize(0, 0)
	}

	return r.MinSize()
}

// Visible returns whether or not this widget should be visible.
// Note that this may not mean it is currently visible if a parent has been hidden.
func (w *Base) Visible() bool {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return !w.hidden
}

// Show this widget so it becomes visible
func (w *Base) Show() {
	if w.Visible() {
		return
	}

	w.setFieldsAndRefresh(func() {
		w.hidden = false
	})
}

// Hide this widget so it is no longer visible
func (w *Base) Hide() {
	if !w.Visible() {
		return
	}

	w.propertyLock.Lock()
	w.hidden = true
	impl := w.impl
	w.propertyLock.Unlock()

	if impl == nil {
		return
	}
	canvas.Refresh(impl)
}

// Refresh causes this widget to be redrawn in it's current state
func (w *Base) Refresh() {
	impl := w.super()
	if impl == nil {
		return
	}

	render := cache.Renderer(impl)
	render.Refresh()
}

// setFieldsAndRefresh helps to make changes to a widget that should be followed by a refresh.
// This method is a guaranteed thread-safe way of directly manipulating widget fields.
func (w *Base) setFieldsAndRefresh(f func()) {
	w.propertyLock.Lock()
	f()
	impl := w.impl
	w.propertyLock.Unlock()

	if impl == nil {
		return
	}
	impl.Refresh()
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is nil.
func (w *Base) super() gui.Widget {
	w.propertyLock.RLock()
	impl := w.impl
	w.propertyLock.RUnlock()
	return impl
}
