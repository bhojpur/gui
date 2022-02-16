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
	"fmt"
	"image/color"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestTable_Empty(t *testing.T) {
	table := &Table{}
	table.Resize(gui.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh() // let's not crash :)
}

func TestTable_Cache(t *testing.T) {
	c := test.NewCanvas()
	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	c.SetContent(table)
	c.SetPadded(false)
	c.Resize(gui.NewSize(120, 148))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 6, len(cellRenderer.Objects()))
	assert.Equal(t, "Cell 0, 0", cellRenderer.Objects()[0].(*Label).Text)
	objRef := cellRenderer.Objects()[0].(*Label)

	test.Scroll(c, gui.NewPos(10, 10), -150, -150)
	assert.Equal(t, float32(0), renderer.scroll.Offset.Y) // we didn't scroll as data shorter
	assert.Equal(t, float32(150), renderer.scroll.Offset.X)
	assert.Equal(t, 6, len(cellRenderer.Objects()))
	assert.Equal(t, "Cell 0, 1", cellRenderer.Objects()[0].(*Label).Text)
	assert.NotEqual(t, objRef, cellRenderer.Objects()[0].(*Label)) // we want to re-use visible cells without rewriting them
}

func TestTable_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	content := test.WidgetRenderer(test.WidgetRenderer(table).(*tableRenderer).scroll.Content.(*tableCells)).(*tableCellsRenderer)
	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(180, 180))
	test.AssertImageMatches(t, "table/theme_initial.png", w.Canvas().Capture())
	assert.Equal(t, "Cell 0, 0", content.Objects()[0].(*Label).Text)
	assert.Equal(t, NewLabel("placeholder").MinSize(), content.Objects()[0].(*Label).Size())

	test.WithTestTheme(t, func() {
		test.WidgetRenderer(table).Refresh()
		test.AssertImageMatches(t, "table/theme_changed.png", w.Canvas().Capture())
	})
	assert.Equal(t, NewLabel("placeholder").MinSize(), content.Objects()[0].(*Label).Size())
}

func TestTable_Filled(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() gui.CanvasObject {
			r := canvas.NewRectangle(color.Black)
			r.SetMinSize(gui.NewSize(30, 20))
			r.Resize(gui.NewSize(30, 20))
			return r
		},
		func(TableCellID, gui.CanvasObject) {})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(180, 180))
	test.AssertImageMatches(t, "table/filled.png", w.Canvas().Capture())
}

func TestTable_MinSize(t *testing.T) {
	for name, tt := range map[string]struct {
		cellSize        gui.Size
		expectedMinSize gui.Size
	}{
		"small": {
			gui.NewSize(1, 1),
			gui.NewSize(float32(32), float32(32)),
		},
		"large": {
			gui.NewSize(100, 100),
			gui.NewSize(100, 100),
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMinSize, NewTable(
				func() (int, int) { return 5, 5 },
				func() gui.CanvasObject {
					r := canvas.NewRectangle(color.Black)
					r.SetMinSize(tt.cellSize)
					r.Resize(tt.cellSize)
					return r
				},
				func(TableCellID, gui.CanvasObject) {}).MinSize())
		})
	}
}

func TestTable_Unselect(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	unselectedRow, unselectedColumn := -1, -1
	table.OnUnselected = func(id TableCellID) {
		unselectedRow = id.Row
		unselectedColumn = id.Col
	}
	table.selectedCell = &TableCellID{1, 1}
	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(180, 180))

	table.Unselect(*table.selectedCell)
	assert.Equal(t, 1, unselectedRow)
	assert.Equal(t, 1, unselectedColumn)
	test.AssertImageMatches(t, "table/theme_initial.png", w.Canvas().Capture())

	unselectedRow, unselectedColumn = -1, -1
	table.Select(TableCellID{2, 2})
	table.Unselect(TableCellID{1, 1})
	assert.Equal(t, -1, unselectedRow)
	assert.Equal(t, -1, unselectedColumn)

	table.UnselectAll()
	assert.Equal(t, 2, unselectedRow)
	assert.Equal(t, 2, unselectedColumn)
}

func TestTable_Refresh(t *testing.T) {
	displayText := "placeholder"
	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() gui.CanvasObject {
			return NewLabel("template")
		},
		func(_ TableCellID, obj gui.CanvasObject) {
			obj.(*Label).SetText(displayText)
		})
	table.Resize(gui.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	assert.Equal(t, "placeholder", cellRenderer.(*tableCellsRenderer).Objects()[7].(*Label).Text)

	displayText = "replaced"
	table.Refresh()
	assert.Equal(t, "replaced", cellRenderer.(*tableCellsRenderer).Objects()[7].(*Label).Text)
}

func TestTable_ScrollTo(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	// for this test the separator thickness is 0
	test.ApplyTheme(t, &separatorThicknessZeroTheme{test.Theme()})

	// we will test a 20 row x 5 column table where each cell is 50x50
	const (
		maxRows int     = 20
		maxCols int     = 5
		width   float32 = 50
		height  float32 = 50
	)

	templ := canvas.NewRectangle(color.Gray16{})
	templ.SetMinSize(gui.Size{Width: width, Height: height})

	table := NewTable(
		func() (int, int) { return maxRows, maxCols },
		func() gui.CanvasObject { return templ },
		func(TableCellID, gui.CanvasObject) {})

	w := test.NewWindow(table)
	defer w.Close()

	// these position expectations have a built-in assumption that the window
	// is smaller than or equal to the size of a single table cell.
	expectedOffset := func(row, col float32) gui.Position {
		return gui.Position{
			X: col * width,
			Y: row * height,
		}
	}

	tt := []struct {
		name string
		in   TableCellID
		want gui.Position
	}{
		{
			"row 0, col 0",
			TableCellID{},
			expectedOffset(0, 0),
		},
		{
			"row 0, col 1",
			TableCellID{Row: 0, Col: 1},
			expectedOffset(0, 1),
		},
		{
			"row 1, col 0",
			TableCellID{Row: 1, Col: 0},
			expectedOffset(1, 0),
		},
		{
			"row 1, col 1",
			TableCellID{Row: 1, Col: 1},
			expectedOffset(1, 1),
		},
		{
			"second last element",
			TableCellID{Row: maxRows - 2, Col: maxCols - 2},
			expectedOffset(float32(maxRows)-2, float32(maxCols)-2),
		},
		{
			"last element",
			TableCellID{Row: maxRows - 1, Col: maxCols - 1},
			expectedOffset(float32(maxRows)-1, float32(maxCols)-1),
		},
		{
			"row 0, col 0 (scrolling backwards)",
			TableCellID{},
			expectedOffset(0, 0),
		},
		{
			"row 99, col 99 (scrolling beyond the end)",
			TableCellID{Row: 99, Col: 99},
			expectedOffset(float32(maxRows)-1, float32(maxCols)-1),
		},
		{
			"row -1, col -1 (scrolling before the start)",
			TableCellID{Row: -1, Col: -1},
			expectedOffset(0, 0),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			table.ScrollTo(tc.in)
			assert.Equal(t, tc.want, table.offset)
			assert.Equal(t, tc.want, table.scroll.Offset)
		})
	}
}

func TestTable_ScrollToBottom(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	const (
		maxRows int     = 20
		maxCols int     = 5
		width   float32 = 50
		height  float32 = 50
	)

	templ := canvas.NewRectangle(color.Gray16{})
	templ.SetMinSize(gui.NewSize(width, height))

	table := NewTable(
		func() (int, int) { return maxRows, maxCols },
		func() gui.CanvasObject { return templ },
		func(TableCellID, gui.CanvasObject) {})

	w := test.NewWindow(table)
	defer w.Close()

	w.Resize(gui.NewSize(200, 200))

	table.ScrollTo(TableCellID{19, 2})
	want := table.offset

	table.ScrollTo(TableCellID{2, 2})
	table.ScrollToBottom()

	assert.Equal(t, want, table.offset)
	assert.Equal(t, want, table.scroll.Offset)
}

func TestTable_ScrollToLeading(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})

	w := test.NewWindow(table)
	defer w.Close()

	table.ScrollTo(TableCellID{Row: 8, Col: 4})
	prev := table.offset
	table.ScrollToLeading()

	want := gui.Position{X: 0, Y: prev.Y}
	assert.Equal(t, want, table.offset)
	assert.Equal(t, want, table.scroll.Offset)
}

func TestTable_ScrollToTop(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	const (
		maxRows int     = 6
		maxCols int     = 10
		width   float32 = 50
		height  float32 = 50
	)

	templ := canvas.NewRectangle(color.Gray16{})
	templ.SetMinSize(gui.Size{Width: width, Height: height})

	table := NewTable(
		func() (int, int) { return maxRows, maxCols },
		func() gui.CanvasObject { return templ },
		func(TableCellID, gui.CanvasObject) {})

	w := test.NewWindow(table)
	defer w.Close()

	table.ScrollTo(TableCellID{12, 3})
	prev := table.offset

	table.ScrollToTop()

	want := gui.Position{X: prev.X, Y: 0}
	assert.Equal(t, want, table.offset)
	assert.Equal(t, want, table.scroll.Offset)
}

func TestTable_ScrollToTrailing(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 24, 24 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})

	w := test.NewWindow(table)
	defer w.Close()

	w.Resize(gui.NewSize(200, 200))

	table.ScrollTo(TableCellID{Row: 7, Col: 23})
	want := table.offset

	table.ScrollTo(TableCellID{Row: 7})
	table.ScrollToTrailing()

	assert.Equal(t, want, table.offset)
	assert.Equal(t, want, table.scroll.Offset)
}

func TestTable_Selection(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	assert.Nil(t, table.selectedCell)

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnSelected = func(id TableCellID) {
		selectedCol = id.Col
		selectedRow = id.Row
	}
	test.TapCanvas(w.Canvas(), gui.NewPos(35, 58))
	assert.Equal(t, 0, table.selectedCell.Col)
	assert.Equal(t, 1, table.selectedCell.Row)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)

	test.AssertRendersToMarkup(t, "table/selected.xml", w.Canvas())

	// check out of bounds col
	w.Resize(gui.NewSize(680, 180))
	test.TapCanvas(w.Canvas(), gui.NewPos(575, 58))
	assert.Equal(t, 0, selectedCol)

	// check out of bounds row
	w.Resize(gui.NewSize(180, 480))
	test.TapCanvas(w.Canvas(), gui.NewPos(35, 428))
	assert.Equal(t, 1, selectedRow)
}

func TestTable_Select(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnSelected = func(id TableCellID) {
		selectedCol = id.Col
		selectedRow = id.Row
	}
	table.Select(TableCellID{1, 0})
	assert.Equal(t, 0, table.selectedCell.Col)
	assert.Equal(t, 1, table.selectedCell.Row)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)
	test.AssertRendersToMarkup(t, "table/selected.xml", w.Canvas())

	table.Select(TableCellID{4, 3})
	assert.Equal(t, 3, table.selectedCell.Col)
	assert.Equal(t, 4, table.selectedCell.Row)
	assert.Equal(t, 3, selectedCol)
	assert.Equal(t, 4, selectedRow)
	test.AssertRendersToMarkup(t, "table/selected_scrolled.xml", w.Canvas())
}

func TestTable_SetColumnWidth(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, obj gui.CanvasObject) {
			if id.Col == 0 {
				obj.(*Label).Text = "p"
			} else {
				obj.(*Label).Text = "placeholder"
			}
			obj.Refresh()
		})
	table.SetColumnWidth(0, 32)
	table.Resize(gui.NewSize(120, 120))
	table.Select(TableCellID{1, 0})

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 10, len(cellRenderer.Objects()))
	assert.Equal(t, float32(32), cellRenderer.(*tableCellsRenderer).Objects()[0].Size().Width)
	cell1Offset := theme.SeparatorThicknessSize()
	assert.Equal(t, float32(32)+cell1Offset, cellRenderer.(*tableCellsRenderer).Objects()[1].Position().X)

	table.SetColumnWidth(0, 16)
	assert.Equal(t, float32(16), cellRenderer.(*tableCellsRenderer).Objects()[0].Size().Width)
	assert.Equal(t, float32(16)+cell1Offset, cellRenderer.(*tableCellsRenderer).Objects()[1].Position().X)

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(120+(2*theme.Padding()), 120+(2*theme.Padding())))
	test.AssertImageMatches(t, "table/col_size.png", w.Canvas().Capture())
}

func TestTable_ShowVisible(t *testing.T) {
	table := NewTable(
		func() (int, int) { return 50, 50 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(TableCellID, gui.CanvasObject) {})
	table.Resize(gui.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 10, len(cellRenderer.Objects()))
}

func TestTable_SeparatorThicknessZero_NotPanics(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	test.ApplyTheme(t, &separatorThicknessZeroTheme{test.Theme()})

	table := NewTable(
		func() (int, int) { return 500, 150 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(TableCellID, gui.CanvasObject) {})

	assert.NotPanics(t, func() {
		table.Resize(gui.NewSize(400, 644))
	})
}

type separatorThicknessZeroTheme struct {
	gui.Theme
}

func (t *separatorThicknessZeroTheme) Size(n gui.ThemeSizeName) float32 {
	if n == theme.SizeNameSeparatorThickness {
		return 0
	}
	return t.Theme.Size(n)
}
