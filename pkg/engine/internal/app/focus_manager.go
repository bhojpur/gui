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
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
)

// FocusManager represents a standard manager of input focus for a canvas
type FocusManager struct {
	sync.RWMutex

	content gui.CanvasObject
	focused gui.Focusable
}

// NewFocusManager returns a new instance of the standard focus manager for a canvas.
func NewFocusManager(c gui.CanvasObject) *FocusManager {
	return &FocusManager{content: c}
}

// Focus focuses the given obj.
func (f *FocusManager) Focus(obj gui.Focusable) bool {
	f.Lock()
	defer f.Unlock()
	if obj != nil {
		var hiddenAncestor gui.CanvasObject
		hidden := false
		found := driver.WalkCompleteObjectTree(
			f.content,
			func(object gui.CanvasObject, _, _ gui.Position, _ gui.Size) bool {
				if hiddenAncestor == nil && !object.Visible() {
					hiddenAncestor = object
				}
				if object == obj.(gui.CanvasObject) {
					hidden = hiddenAncestor != nil
					return true
				}
				return false
			},
			func(object, _ gui.CanvasObject) {
				if hiddenAncestor == object {
					hiddenAncestor = nil
				}
			},
		)
		if !found {
			return false
		}
		if hidden {
			return true
		}
		if dis, ok := obj.(gui.Disableable); ok && dis.Disabled() {
			type selectableText interface {
				SelectedText() string
			}
			if _, isSelectableText := obj.(selectableText); !isSelectableText || gui.CurrentDevice().IsMobile() {
				return true
			}
		}
	}
	f.focus(obj)
	return true
}

// Focused returns the currently focused object or nil if none.
func (f *FocusManager) Focused() gui.Focusable {
	f.RLock()
	defer f.RUnlock()
	return f.focused
}

// FocusGained signals to the manager that its content got focus (due to window/overlay switch for instance).
func (f *FocusManager) FocusGained() {
	if focused := f.Focused(); focused != nil {
		focused.FocusGained()
	}
}

// FocusLost signals to the manager that its content lost focus (due to window/overlay switch for instance).
func (f *FocusManager) FocusLost() {
	if focused := f.Focused(); focused != nil {
		focused.FocusLost()
	}
}

// FocusNext will find the item after the current that can be focused and focus it.
// If current is nil then the first focusable item in the canvas will be focused.
func (f *FocusManager) FocusNext() {
	f.Lock()
	defer f.Unlock()
	f.focus(f.nextInChain(f.focused))
}

// FocusPrevious will find the item before the current that can be focused and focus it.
// If current is nil then the last focusable item in the canvas will be focused.
func (f *FocusManager) FocusPrevious() {
	f.Lock()
	defer f.Unlock()
	f.focus(f.previousInChain(f.focused))
}

func (f *FocusManager) focus(obj gui.Focusable) {
	if f.focused == obj {
		return
	}

	if f.focused != nil {
		f.focused.FocusLost()
	}
	f.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (f *FocusManager) nextInChain(current gui.Focusable) gui.Focusable {
	return f.nextWithWalker(current, driver.WalkVisibleObjectTree)
}

func (f *FocusManager) nextWithWalker(current gui.Focusable, walker walkerFunc) gui.Focusable {
	var next gui.Focusable
	found := current == nil // if we have no starting point then pretend we matched already
	walker(f.content, func(obj gui.CanvasObject, _ gui.Position, _ gui.Position, _ gui.Size) bool {
		if w, ok := obj.(gui.Disableable); ok && w.Disabled() {
			// disabled widget cannot receive focus
			return false
		}

		focus, ok := obj.(gui.Focusable)
		if !ok {
			return false
		}

		if found {
			next = focus
			return true
		}
		if next == nil {
			next = focus
		}

		if obj == current.(gui.CanvasObject) {
			found = true
		}

		return false
	}, nil)

	return next
}

func (f *FocusManager) previousInChain(current gui.Focusable) gui.Focusable {
	return f.nextWithWalker(current, driver.ReverseWalkVisibleObjectTree)
}

type walkerFunc func(
	gui.CanvasObject,
	func(gui.CanvasObject, gui.Position, gui.Position, gui.Size) bool,
	func(gui.CanvasObject, gui.CanvasObject),
) bool
