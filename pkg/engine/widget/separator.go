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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Widget = (*Separator)(nil)

// Separator is a widget for displaying a separator with themeable color.
//
// Since: 1.4
type Separator struct {
	BaseWidget
}

// NewSeparator creates a new separator.
//
// Since: 1.4
func NewSeparator() *Separator {
	s := &Separator{}
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer returns a new renderer for the separator.
//
// Implements: gui.Widget
func (s *Separator) CreateRenderer() gui.WidgetRenderer {
	s.ExtendBaseWidget(s)
	bar := canvas.NewRectangle(theme.DisabledColor())
	return &separatorRenderer{
		WidgetRenderer: NewSimpleRenderer(bar),
		bar:            bar,
		d:              s,
	}
}

// MinSize returns the minimal size of the separator.
//
// Implements: gui.Widget
func (s *Separator) MinSize() gui.Size {
	s.ExtendBaseWidget(s)
	t := theme.SeparatorThicknessSize()
	return gui.NewSize(t, t)
}

var _ gui.WidgetRenderer = (*separatorRenderer)(nil)

type separatorRenderer struct {
	gui.WidgetRenderer
	bar *canvas.Rectangle
	d   *Separator
}

func (r *separatorRenderer) MinSize() gui.Size {
	t := theme.SeparatorThicknessSize()
	return gui.NewSize(t, t)
}

func (r *separatorRenderer) Refresh() {
	r.bar.FillColor = theme.DisabledColor()
	canvas.Refresh(r.d)
}
