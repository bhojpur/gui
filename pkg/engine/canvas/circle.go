package canvas

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
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// Declare conformity with CanvasObject interface
var _ gui.CanvasObject = (*Circle)(nil)

// Circle describes a colored circle primitive in a Bhojpur GUI canvas
type Circle struct {
	Position1 gui.Position // The current top-left position of the Circle
	Position2 gui.Position // The current bottomright position of the Circle
	Hidden    bool         // Is this circle currently hidden

	FillColor   color.Color // The circle fill color
	StrokeColor color.Color // The circle stroke color
	StrokeWidth float32     // The stroke width of the circle
}

// Size returns the current size of bounding box for this circle object
func (l *Circle) Size() gui.Size {
	return gui.NewSize(l.Position2.X-l.Position1.X, l.Position2.Y-l.Position1.Y)
}

// Resize sets a new bottom-right position for the circle object
// If it has a stroke width this will cause it to Refresh.
func (l *Circle) Resize(size gui.Size) {
	if size == l.Size() {
		return
	}

	l.Position2 = gui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)

	Refresh(l)
}

// Position gets the current top-left position of this circle object, relative to its parent / canvas
func (l *Circle) Position() gui.Position {
	return l.Position1
}

// Move the circle object to a new position, relative to its parent / canvas
func (l *Circle) Move(pos gui.Position) {
	size := l.Size()
	l.Position1 = pos
	l.Position2 = gui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// MinSize for a Circle simply returns Size{1, 1} as there is no
// explicit content
func (l *Circle) MinSize() gui.Size {
	return gui.NewSize(1, 1)
}

// Visible returns true if this circle is visible, false otherwise
func (l *Circle) Visible() bool {
	return !l.Hidden
}

// Show will set this circle to be visible
func (l *Circle) Show() {
	l.Hidden = false

	l.Refresh()
}

// Hide will set this circle to not be visible
func (l *Circle) Hide() {
	l.Hidden = true

	l.Refresh()
}

// Refresh causes this object to be redrawn in it's current state
func (l *Circle) Refresh() {
	Refresh(l)
}

// NewCircle returns a new Circle instance
func NewCircle(color color.Color) *Circle {
	return &Circle{
		FillColor: color,
	}
}
