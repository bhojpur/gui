package spatial

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
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type mapButton struct {
	widget.Button
}

func newMapButton(icon gui.Resource, f func()) *mapButton {
	b := &mapButton{}
	b.ExtendBaseWidget(b)

	b.Icon = icon
	b.OnTapped = f
	return b
}

func (b *mapButton) CreateRenderer() gui.WidgetRenderer {
	return &mapButtonRenderer{WidgetRenderer: b.Button.CreateRenderer(),
		bg: canvas.NewRectangle(theme.ShadowColor())}
}

type mapButtonRenderer struct {
	gui.WidgetRenderer

	bg *canvas.Rectangle
}

func (r *mapButtonRenderer) Layout(s gui.Size) {
	halfPad := theme.Padding() / 2
	r.bg.Move(gui.NewPos(halfPad, halfPad))
	r.bg.Resize(s.Subtract(gui.NewSize(theme.Padding(), theme.Padding())))

	r.WidgetRenderer.Layout(s)
}

func (r *mapButtonRenderer) Objects() []gui.CanvasObject {
	return append([]gui.CanvasObject{r.bg}, r.WidgetRenderer.Objects()...)
}

func (r *mapButtonRenderer) Refresh() {
	r.bg.FillColor = theme.ShadowColor()
	r.WidgetRenderer.Refresh()
}
