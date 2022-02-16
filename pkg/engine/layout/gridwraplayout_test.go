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

func TestGridLWrapLayout_Layout(t *testing.T) {
	gridSize := gui.NewSize(125, 125)
	cellSize := gui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridWrapLayout(cellSize).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, obj2.Position(), cell2Pos)
	cell3Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj3.Position(), cell3Pos)
}

func TestGridLWrapLayout_Layout_Min(t *testing.T) {
	cellSize := gui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}

	layout.NewGridWrapLayout(cellSize).Layout(container.Objects, container.MinSize())

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj2.Position(), cell2Pos)
	cell3Pos := gui.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, obj3.Position(), cell3Pos)
}

func TestGridLWrapLayout_Layout_HiddenItem(t *testing.T) {
	gridSize := gui.NewSize(125, 125)
	cellSize := gui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2.Hide()
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj4 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3, obj4},
	}
	container.Resize(gridSize)

	layout.NewGridWrapLayout(cellSize).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj3.Position(), gui.NewPos(50+theme.Padding(), 0))
	assert.Equal(t, obj4.Position(), gui.NewPos(0, 50+theme.Padding()))
}

func TestGridLWrapLayout_MinSize(t *testing.T) {
	cellSize := gui.NewSize(50, 50)
	minSize := cellSize

	container := gui.NewContainer(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	layout := layout.NewGridWrapLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)

	// This has a dynamic minSize so we need to check again after layout!
	layout.Layout(container.Objects, minSize)
	layoutMin = layout.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)
}

func TestGridLWrapLayout_MinSize_Hidden(t *testing.T) {
	cellSize := gui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2.Hide()
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := gui.NewContainer(obj1, obj2, obj3)
	layout := layout.NewGridWrapLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	assert.Equal(t, gui.NewSize(50, 50), layoutMin)

	// This has a dynamic minSize so we need to check again after layout!
	layout.Layout(container.Objects, gui.NewSize(50, 75))
	layoutMin = layout.MinSize(container.Objects)
	assert.Equal(t, gui.NewSize(50, 100+theme.Padding()), layoutMin)
}

func TestGridLWrapLayout_Resize_LessThanMinSize(t *testing.T) {
	cellSize := gui.NewSize(50, 50)
	minSize := cellSize

	container := gui.NewContainer(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	l := layout.NewGridWrapLayout(cellSize)
	container.Resize(gui.NewSize(25, 25))

	layoutMin := l.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)
}
