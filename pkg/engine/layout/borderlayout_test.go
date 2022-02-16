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
	_ "github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewBorderContainer(t *testing.T) {
	top := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	top.SetMinSize(gui.NewSize(10, 10))
	right := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	right.SetMinSize(gui.NewSize(10, 10))
	middle := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	c := gui.NewContainerWithLayout(layout.NewBorderLayout(top, nil, nil, right), top, right, middle)
	assert.Equal(t, 3, len(c.Objects))

	c.Resize(gui.NewSize(100, 100))
	assert.Equal(t, float32(0), top.Position().X)
	assert.Equal(t, float32(0), top.Position().Y)
	assert.Equal(t, float32(90), right.Position().X)
	assert.Equal(t, 10+theme.Padding(), right.Position().Y)
	assert.Equal(t, float32(0), middle.Position().X)
	assert.Equal(t, 10+theme.Padding(), middle.Position().Y)
	assert.Equal(t, 90-theme.Padding(), middle.Size().Width)
	assert.Equal(t, 90-theme.Padding(), middle.Size().Height)
}

func TestBorderLayout_Size_Empty(t *testing.T) {
	size := gui.NewSize(100, 100)

	obj := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	container := &gui.Container{
		Objects: []gui.CanvasObject{obj},
	}
	container.Resize(size)

	layout.NewBorderLayout(nil, nil, nil, nil).Layout(container.Objects, size)

	assert.Equal(t, obj.Size(), size)
}

func TestBorderLayout_Size_TopBottom(t *testing.T) {
	size := gui.NewSize(100, 100)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(size)

	layout.NewBorderLayout(obj1, obj2, nil, nil).Layout(container.Objects, size)

	innerSize := gui.NewSize(size.Width, size.Height-obj1.Size().Height-obj2.Size().Height-theme.Padding()*2)
	assert.Equal(t, innerSize, obj3.Size())
	assert.Equal(t, gui.NewPos(0, 0), obj1.Position())
	assert.Equal(t, gui.NewPos(0, size.Height-obj2.Size().Height), obj2.Position())
	assert.Equal(t, gui.NewPos(0, obj1.Size().Height+theme.Padding()), obj3.Position())
}

func TestBorderLayout_Size_LeftRight(t *testing.T) {
	size := gui.NewSize(100, 100)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &gui.Container{
		Objects: []gui.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(size)

	layout.NewBorderLayout(nil, nil, obj1, obj2).Layout(container.Objects, size)

	innerSize := gui.NewSize(size.Width-obj1.Size().Width-obj2.Size().Width-theme.Padding()*2, size.Height)
	assert.Equal(t, innerSize, obj3.Size())
	assert.Equal(t, gui.NewPos(0, 0), obj1.Position())
	assert.Equal(t, gui.NewPos(size.Width-obj2.Size().Width, 0), obj2.Position())
	assert.Equal(t, gui.NewPos(obj1.Size().Width+theme.Padding(), 0), obj3.Position())
}

func TestBorderLayout_MinSize_Center(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := gui.NewContainer(text)
	layoutMin := layout.NewBorderLayout(nil, nil, nil, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderLayout_MinSize_TopBottom(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text2 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text3 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := gui.NewSize(text3.MinSize().Width, text1.MinSize().Height+text2.MinSize().Height+text3.MinSize().Height+theme.Padding()*2)

	container := gui.NewContainer(text1, text2, text3)
	layoutMin := layout.NewBorderLayout(text1, text2, nil, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderLayout_MinSize_TopBottomHidden(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text1.Hide()
	text2 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text2.Hide()
	text3 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})

	container := gui.NewContainer(text1, text2, text3)
	layoutMin := layout.NewBorderLayout(text1, text2, nil, nil).MinSize(container.Objects)

	assert.Equal(t, text1.MinSize(), layoutMin)
}

func TestBorderLayout_MinSize_TopOnly(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := gui.NewSize(text1.MinSize().Width, text1.MinSize().Height+theme.Padding())

	container := gui.NewContainer(text1)
	layoutMin := layout.NewBorderLayout(text1, nil, nil, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderLayout_MinSize_LeftRight(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text2 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text3 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := gui.NewSize(text1.MinSize().Width+text2.MinSize().Width+text3.MinSize().Width+theme.Padding()*2, text3.MinSize().Height)

	container := gui.NewContainer(text1, text2, text3)
	layoutMin := layout.NewBorderLayout(nil, nil, text1, text2).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderLayout_MinSize_LeftRightHidden(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text1.Hide()
	text2 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text2.Hide()
	text3 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})

	container := gui.NewContainer(text1, text2, text3)
	layoutMin := layout.NewBorderLayout(nil, nil, text1, text2).MinSize(container.Objects)

	assert.Equal(t, text3.MinSize(), layoutMin)
}

func TestBorderLayout_MinSize_LeftOnly(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := gui.NewSize(text1.MinSize().Width+theme.Padding(), text1.MinSize().Height)

	container := gui.NewContainer(text1)
	layoutMin := layout.NewBorderLayout(nil, nil, text1, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}
