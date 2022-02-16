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
var _ gui.CanvasObject = (*Text)(nil)

// Text describes a text primitive in a Bhojpur GUI canvas.
// A text object can have a style set which will apply to the whole string.
// No formatting or text parsing will be performed
type Text struct {
	baseObject
	Alignment gui.TextAlign // The alignment of the text content

	Color     color.Color   // The main text draw color
	Text      string        // The string content of this Text
	TextSize  float32       // Size of the text - if the Canvas scale is 1.0 this will be equivalent to point size
	TextStyle gui.TextStyle // The style of the text content
}

// MinSize returns the minimum size of this text object based on its font size and content.
// This is normally determined by the render implementation.
func (t *Text) MinSize() gui.Size {
	return gui.MeasureText(t.Text, t.TextSize, t.TextStyle)
}

// SetMinSize has no effect as the smallest size this canvas object can be is based on its font size and content.
func (t *Text) SetMinSize(size gui.Size) {
	// no-op
}

// Refresh causes this object to be redrawn in it's current state
func (t *Text) Refresh() {
	Refresh(t)
}

// NewText returns a new Text implementation
func NewText(text string, color color.Color) *Text {
	size := float32(0)
	if gui.CurrentApp() != nil { // nil app possible if app not started
		size = gui.CurrentApp().Settings().Theme().Size("text") // manually name the size to avoid import loop
	}
	return &Text{
		Color:    color,
		Text:     text,
		TextSize: size,
	}
}
