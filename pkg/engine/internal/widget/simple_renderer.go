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

import gui "github.com/bhojpur/gui/pkg/engine"

var _ gui.WidgetRenderer = (*SimpleRenderer)(nil)

// SimpleRenderer is a basic renderer that satisfies widget.Renderer interface by wrapping
// a single gui.CanvasObject.
//
// Since: 2.1
type SimpleRenderer struct {
	objects []gui.CanvasObject
}

// NewSimpleRenderer creates a new SimpleRenderer to render a widget using a
// single CanvasObject.
//
// Since: 2.1
func NewSimpleRenderer(object gui.CanvasObject) *SimpleRenderer {
	return &SimpleRenderer{[]gui.CanvasObject{object}}
}

// Destroy does nothing in this implementation.
//
// Implements: gui.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Destroy() {
}

// Layout updates the contained object to be the requested size.
//
// Implements: gui.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Layout(s gui.Size) {
	r.objects[0].Resize(s)
}

// MinSize returns the smallest size that this render can use, returned from the underlying object.
//
// Implements: gui.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) MinSize() gui.Size {
	return r.objects[0].MinSize()
}

// Objects returns the objects that should be rendered.
//
// Implements: gui.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Objects() []gui.CanvasObject {
	return r.objects
}

// Refresh requests the underlying object to redraw.
//
// Implements: gui.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Refresh() {
	r.objects[0].Refresh()
}
