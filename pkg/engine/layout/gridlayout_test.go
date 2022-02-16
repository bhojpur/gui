package layout_test

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestGridLayout(t *testing.T) {
	gridSize := gui.NewSize(100+theme.Padding(), 100+theme.Padding())
	cellSize := gui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridLayout(2).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestGridLayoutRounding(t *testing.T) {
	gridSize := gui.NewSize(100+theme.Padding()*2, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridLayout(3).Layout(container.Objects, gridSize)

	assert.Equal(t, gui.NewPos(0, 0), obj1.Position())
	assert.Equal(t, gui.NewSize(33, 50), obj1.Size())
	assert.Equal(t, gui.NewPos(33+theme.Padding(), 0), obj2.Position())
	assert.Equal(t, gui.NewSize(34, 50), obj2.Size())
	assert.Equal(t, gui.NewPos(67+theme.Padding()*2, 0), obj3.Position())
	assert.Equal(t, gui.NewSize(33, 50), obj3.Size())
}

func TestGridLayout_Vertical(t *testing.T) {
	gridSize := gui.NewSize(100+theme.Padding(), 100+theme.Padding())
	cellSize := gui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridLayoutWithRows(2).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestGridLayout_MinSize(t *testing.T) {
	text1 := canvas.NewText("Large Text", color.NRGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("small", color.NRGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(gui.NewSize(0, text2.MinSize().Height+theme.Padding()))

	container := gui.NewContainer(text1, text2)
	layoutMin := layout.NewGridLayout(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestGridLayout_MinSize_Vertical(t *testing.T) {
	text1 := canvas.NewText("Text", color.NRGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("Text", color.NRGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(gui.NewSize(text2.MinSize().Width+theme.Padding(), 0))

	container := gui.NewContainer(text1, text2)
	layoutMin := layout.NewGridLayoutWithRows(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestGridLayout_MinSize_HiddenItem(t *testing.T) {
	text1 := canvas.NewText("Large Text", color.NRGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("hidden", color.NRGBA{0xff, 0, 0, 0})
	text2.Hide()
	text3 := canvas.NewText("small", color.NRGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(gui.NewSize(0, text3.MinSize().Height+theme.Padding()))

	container := gui.NewContainer(text1, text2, text3)
	layoutMin := layout.NewGridLayout(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}
