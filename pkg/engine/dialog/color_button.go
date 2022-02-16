package dialog

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
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	internalwidget "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

var _ gui.Widget = (*colorButton)(nil)
var _ desktop.Hoverable = (*colorButton)(nil)

// colorButton displays a color and triggers the callback when tapped.
type colorButton struct {
	widget.BaseWidget
	color   color.Color
	onTap   func(color.Color)
	hovered bool
}

// newColorButton creates a colorButton with the given color and callback.
func newColorButton(color color.Color, onTap func(color.Color)) *colorButton {
	b := &colorButton{
		color: color,
		onTap: onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (b *colorButton) CreateRenderer() gui.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := newCheckeredBackground()
	rectangle := &canvas.Rectangle{
		FillColor: b.color,
	}
	return &colorButtonRenderer{
		BaseRenderer: internalwidget.NewBaseRenderer([]gui.CanvasObject{background, rectangle}),
		button:       b,
		background:   background,
		rectangle:    rectangle,
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (b *colorButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (b *colorButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *colorButton) MouseMoved(*desktop.MouseEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (b *colorButton) MinSize() gui.Size {
	return b.BaseWidget.MinSize()
}

// SetColor updates the color selected in this color widget
func (b *colorButton) SetColor(color color.Color) {
	if b.color == color {
		return
	}
	b.color = color
	b.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (b *colorButton) Tapped(*gui.PointEvent) {
	writeRecentColor(colorToString(b.color))
	if f := b.onTap; f != nil {
		f(b.color)
	}
}

type colorButtonRenderer struct {
	internalwidget.BaseRenderer
	button     *colorButton
	background *canvas.Raster
	rectangle  *canvas.Rectangle
}

func (r *colorButtonRenderer) Layout(size gui.Size) {
	r.rectangle.Move(gui.NewPos(0, 0))
	r.rectangle.Resize(size)
}

func (r *colorButtonRenderer) MinSize() gui.Size {
	return r.rectangle.MinSize().Max(gui.NewSize(32, 32))
}

func (r *colorButtonRenderer) Refresh() {
	if r.button.hovered {
		r.rectangle.StrokeColor = theme.HoverColor()
		r.rectangle.StrokeWidth = float32(theme.Padding())
	} else {
		r.rectangle.StrokeWidth = 0
	}
	r.rectangle.FillColor = r.button.color
	canvas.Refresh(r.button)
}
