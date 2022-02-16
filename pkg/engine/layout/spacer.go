package layout

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

// SpacerObject is any object that can be used to space out child objects
type SpacerObject interface {
	ExpandVertical() bool
	ExpandHorizontal() bool
}

// Spacer is any simple object that can be used in a box layout to space
// out child objects
type Spacer struct {
	FixHorizontal bool
	FixVertical   bool

	size   gui.Size
	pos    gui.Position
	hidden bool
}

// NewSpacer returns a spacer object which can fill vertical and horizontal
// space. This is primarily used with a box layout.
func NewSpacer() gui.CanvasObject {
	return &Spacer{}
}

// ExpandVertical returns whether or not this spacer expands on the vertical axis
func (s *Spacer) ExpandVertical() bool {
	return !s.FixVertical
}

// ExpandHorizontal returns whether or not this spacer expands on the horizontal axis
func (s *Spacer) ExpandHorizontal() bool {
	return !s.FixHorizontal
}

// Size returns the current size of this Spacer
func (s *Spacer) Size() gui.Size {
	return s.size
}

// Resize sets a new size for the Spacer - this will be called by the layout
func (s *Spacer) Resize(size gui.Size) {
	s.size = size
}

// Position returns the current position of this Spacer
func (s *Spacer) Position() gui.Position {
	return s.pos
}

// Move sets a new position for the Spacer - this will be called by the layout
func (s *Spacer) Move(pos gui.Position) {
	s.pos = pos
}

// MinSize returns a 0 size as a Spacer can shrink to no actual size
func (s *Spacer) MinSize() gui.Size {
	return gui.NewSize(0, 0)
}

// Visible returns true if this spacer should affect the layout
func (s *Spacer) Visible() bool {
	return !s.hidden
}

// Show sets the Spacer to be part of the layout calculations
func (s *Spacer) Show() {
	s.hidden = false
}

// Hide removes this Spacer from layout calculations
func (s *Spacer) Hide() {
	s.hidden = true
}

// Refresh does nothing for a spacer but is part of the CanvasObject definition
func (s *Spacer) Refresh() {
}
