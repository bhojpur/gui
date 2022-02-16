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
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// Declare conformity with CanvasObject interface
var _ gui.CanvasObject = (*Line)(nil)

// Line describes a colored line primitive in a Bhojpur GUI canvas.
// Lines are special as they can have a negative width or height to indicate
// an inverse slope (i.e. slope up vs down).
type Line struct {
	Position1 gui.Position // The current top-left position of the Line
	Position2 gui.Position // The current bottomright position of the Line
	Hidden    bool         // Is this Line currently hidden

	StrokeColor color.Color // The line stroke color
	StrokeWidth float32     // The stroke width of the line
}

// Size returns the current size of bounding box for this line object
func (l *Line) Size() gui.Size {
	return gui.NewSize(float32(math.Abs(float64(l.Position2.X)-float64(l.Position1.X))),
		float32(math.Abs(float64(l.Position2.Y)-float64(l.Position1.Y))))
}

// Resize sets a new bottom-right position for the line object and it will then be refreshed.
func (l *Line) Resize(size gui.Size) {
	if size == l.Size() {
		return
	}

	if l.Position1.X <= l.Position2.X {
		l.Position2.X = l.Position1.X + size.Width
	} else {
		l.Position1.X = l.Position2.X + size.Width
	}
	if l.Position1.Y <= l.Position2.Y {
		l.Position2.Y = l.Position1.Y + size.Height
	} else {
		l.Position1.Y = l.Position2.Y + size.Height
	}
	Refresh(l)
}

// Position gets the current top-left position of this line object, relative to its parent / canvas
func (l *Line) Position() gui.Position {
	return gui.NewPos(gui.Min(l.Position1.X, l.Position2.X), gui.Min(l.Position1.Y, l.Position2.Y))
}

// Move the line object to a new position, relative to its parent / canvas
func (l *Line) Move(pos gui.Position) {
	oldPos := l.Position()
	deltaX := pos.X - oldPos.X
	deltaY := pos.Y - oldPos.Y

	l.Position1 = l.Position1.Add(gui.NewPos(deltaX, deltaY))
	l.Position2 = l.Position2.Add(gui.NewPos(deltaX, deltaY))
}

// MinSize for a Line simply returns Size{1, 1} as there is no
// explicit content
func (l *Line) MinSize() gui.Size {
	return gui.NewSize(1, 1)
}

// Visible returns true if this line// Show will set this circle to be visible is visible, false otherwise
func (l *Line) Visible() bool {
	return !l.Hidden
}

// Show will set this line to be visible
func (l *Line) Show() {
	l.Hidden = false

	l.Refresh()
}

// Hide will set this line to not be visible
func (l *Line) Hide() {
	l.Hidden = true

	l.Refresh()
}

// Refresh causes this object to be redrawn in it's current state
func (l *Line) Refresh() {
	Refresh(l)
}

// NewLine returns a new Line instance
func NewLine(color color.Color) *Line {
	return &Line{
		StrokeColor: color,
		StrokeWidth: 1,
	}
}
