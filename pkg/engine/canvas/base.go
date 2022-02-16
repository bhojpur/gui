// It contains all of the primitive CanvasObjects that make up a Bhojpur GUI
//
// The types implemented in this package are used as building blocks in order
// to build higher order functionality. These types are designed to be
// non-interactive, by design. If additional functonality is required,
// it's usually a sign that this type should be used as part of a custom
// Widget.
package canvas // import gui "github.com/bhojpur/gui/pkg/engine/canvas"

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
)

type baseObject struct {
	size     gui.Size     // The current size of the canvas object
	position gui.Position // The current position of the object
	Hidden   bool         // Is this object currently hidden

	min gui.Size // The minimum size this object can be

	propertyLock sync.RWMutex
}

// CurrentSize returns the current size of this canvas object.
func (r *baseObject) Size() gui.Size {
	r.propertyLock.RLock()
	defer r.propertyLock.RUnlock()

	return r.size
}

// Resize sets a new size for the canvas object.
func (r *baseObject) Resize(size gui.Size) {
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()

	r.size = size
}

// CurrentPosition gets the current position of this canvas object, relative to its parent.
func (r *baseObject) Position() gui.Position {
	r.propertyLock.RLock()
	defer r.propertyLock.RUnlock()

	return r.position
}

// Move the object to a new position, relative to its parent.
func (r *baseObject) Move(pos gui.Position) {
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()

	r.position = pos
}

// MinSize returns the specified minimum size, if set, or {1, 1} otherwise.
func (r *baseObject) MinSize() gui.Size {
	r.propertyLock.RLock()
	defer r.propertyLock.RUnlock()

	if r.min.Width == 0 && r.min.Height == 0 {
		return gui.NewSize(1, 1)
	}

	return r.min
}

// SetMinSize specifies the smallest size this object should be.
func (r *baseObject) SetMinSize(size gui.Size) {
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()

	r.min = size
}

// IsVisible returns true if this object is visible, false otherwise.
func (r *baseObject) Visible() bool {
	r.propertyLock.RLock()
	defer r.propertyLock.RUnlock()

	return !r.Hidden
}

// Show will set this object to be visible.
func (r *baseObject) Show() {
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()

	r.Hidden = false
}

// Hide will set this object to not be visible.
func (r *baseObject) Hide() {
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()

	r.Hidden = true
}

// Refresh instructs the containing canvas to refresh the specified obj.
func Refresh(obj gui.CanvasObject) {
	if gui.CurrentApp() == nil || gui.CurrentApp().Driver() == nil {
		return
	}

	c := gui.CurrentApp().Driver().CanvasForObject(obj)
	if c != nil {
		c.Refresh(obj)
	}
}
