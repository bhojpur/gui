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

import (
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// Declare conformity with Layout interface
var _ gui.Layout = (*gridWrapLayout)(nil)

type gridWrapLayout struct {
	CellSize gui.Size
	colCount int
	rowCount int
}

// NewGridWrapLayout returns a new GridWrapLayout instance
func NewGridWrapLayout(size gui.Size) gui.Layout {
	return &gridWrapLayout{size, 1, 1}
}

// Layout is called to pack all child objects into a specified size.
// For a GridWrapLayout this will attempt to lay all the child objects in a row
// and wrap to a new row if the size is not large enough.
func (g *gridWrapLayout) Layout(objects []gui.CanvasObject, size gui.Size) {
	g.colCount = 1
	g.rowCount = 1

	if size.Width > g.CellSize.Width {
		g.colCount = int(math.Floor(float64(size.Width+theme.Padding()) / float64(g.CellSize.Width+theme.Padding())))
	}

	i, x, y := 0, float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		child.Move(gui.NewPos(x, y))
		child.Resize(g.CellSize)

		if (i+1)%g.colCount == 0 {
			x = 0
			y += g.CellSize.Height + theme.Padding()
			if i > 0 {
				g.rowCount++
			}
		} else {
			x += g.CellSize.Width + theme.Padding()
		}
		i++
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridWrapLayout this is simply the specified cellsize as a single column
// layout has no padding. The returned size does not take into account the number
// of columns as this layout re-flows dynamically.
func (g *gridWrapLayout) MinSize(objects []gui.CanvasObject) gui.Size {
	return gui.NewSize(g.CellSize.Width,
		(g.CellSize.Height*float32(g.rowCount))+(float32(g.rowCount-1)*theme.Padding()))
}
