//go:build !mobile
// +build !mobile

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/test"

	"github.com/stretchr/testify/assert"
)

func TestTable_Hovered(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 2, 2 },
		func() gui.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c gui.CanvasObject) {
			c.(*Label).SetText(fmt.Sprintf("Cell %d, %d", id.Row, id.Col))
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(gui.NewSize(180, 180))

	test.MoveMouse(w.Canvas(), gui.NewPos(35, 50))
	test.MoveMouse(w.Canvas(), gui.NewPos(35, 100))

	assert.Nil(t, table.hoveredCell)

	test.AssertRendersToMarkup(t, "table/desktop/hovered_out.xml", w.Canvas())

	table.Length = func() (int, int) { return 3, 5 }
	table.Refresh()

	w.SetContent(table)
	w.Resize(gui.NewSize(180, 180))
	test.MoveMouse(w.Canvas(), gui.NewPos(35, 58))

	assert.Equal(t, 0, table.hoveredCell.Col)
	assert.Equal(t, 1, table.hoveredCell.Row)

	test.AssertRendersToMarkup(t, "table/desktop/hovered.xml", w.Canvas())
}
