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
var _ gui.Layout = (*gridLayout)(nil)

type gridLayout struct {
	Cols            int
	vertical, adapt bool
}

// NewAdaptiveGridLayout returns a new grid layout which uses columns when horizontal but rows when vertical.
func NewAdaptiveGridLayout(rowcols int) gui.Layout {
	return &gridLayout{Cols: rowcols, adapt: true}
}

// NewGridLayout returns a grid layout arranged in a specified number of columns.
// The number of rows will depend on how many children are in the container that uses this layout.
func NewGridLayout(cols int) gui.Layout {
	return NewGridLayoutWithColumns(cols)
}

// NewGridLayoutWithColumns returns a new grid layout that specifies a column count and wrap to new rows when needed.
func NewGridLayoutWithColumns(cols int) gui.Layout {
	return &gridLayout{Cols: cols}
}

// NewGridLayoutWithRows returns a new grid layout that specifies a row count that creates new rows as required.
func NewGridLayoutWithRows(rows int) gui.Layout {
	return &gridLayout{Cols: rows, vertical: true}
}

func (g *gridLayout) horizontal() bool {
	if g.adapt {
		return gui.IsHorizontal(gui.CurrentDevice().Orientation())
	}

	return !g.vertical
}

func (g *gridLayout) countRows(objects []gui.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(g.Cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getLeading(size float64, offset int) float32 {
	ret := (size + float64(theme.Padding())) * float64(offset)

	return float32(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int) float32 {
	return getLeading(size, offset+1) - theme.Padding()
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *gridLayout) Layout(objects []gui.CanvasObject, size gui.Size) {
	rows := g.countRows(objects)

	padWidth := float32(g.Cols-1) * theme.Padding()
	padHeight := float32(rows-1) * theme.Padding()
	cellWidth := float64(size.Width-padWidth) / float64(g.Cols)
	cellHeight := float64(size.Height-padHeight) / float64(rows)

	if !g.horizontal() {
		padWidth, padHeight = padHeight, padWidth
		cellWidth = float64(size.Width-padWidth) / float64(rows)
		cellHeight = float64(size.Height-padHeight) / float64(g.Cols)
	}

	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		child.Move(gui.NewPos(x1, y1))
		child.Resize(gui.NewSize(x2-x1, y2-y1))

		if g.horizontal() {
			if (i+1)%g.Cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%g.Cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridLayout this is the size of the largest child object multiplied by
// the required number of columns and rows, with appropriate padding between
// children.
func (g *gridLayout) MinSize(objects []gui.CanvasObject) gui.Size {
	rows := g.countRows(objects)
	minSize := gui.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	if g.horizontal() {
		minContentSize := gui.NewSize(minSize.Width*float32(g.Cols), minSize.Height*float32(rows))
		return minContentSize.Add(gui.NewSize(theme.Padding()*gui.Max(float32(g.Cols-1), 0), theme.Padding()*gui.Max(float32(rows-1), 0)))
	}

	minContentSize := gui.NewSize(minSize.Width*float32(rows), minSize.Height*float32(g.Cols))
	return minContentSize.Add(gui.NewSize(theme.Padding()*gui.Max(float32(rows-1), 0), theme.Padding()*gui.Max(float32(g.Cols-1), 0)))
}
