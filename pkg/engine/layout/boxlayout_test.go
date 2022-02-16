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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

// NewRectangle returns a new Rectangle instance
func NewMinSizeRect(min gui.Size) *canvas.Rectangle {
	rect := &canvas.Rectangle{}
	rect.SetMinSize(min)

	return rect
}

func TestHBoxLayout_Simple(t *testing.T) {
	cellSize := gui.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), gui.NewSize(150+(theme.Padding()*2), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_HiddenItem(t *testing.T) {
	cellSize := gui.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj2.Hide()
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), gui.NewSize(100+(theme.Padding()), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell3Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_Wide(t *testing.T) {
	cellSize := gui.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	container.Resize(gui.NewSize(308, 100))
	assert.Equal(t, gui.NewSize(150+(theme.Padding()*2), 100), container.MinSize())

	assert.Equal(t, float32(50), obj1.Size().Width)
	assert.Equal(t, float32(50), obj2.Size().Width)
	cell2Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_Tall(t *testing.T) {
	cellSize := gui.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), gui.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_Spacer(t *testing.T) {
	cellSize := gui.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), obj1, obj2, obj3)
	container.Resize(gui.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), gui.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, float32(50), obj1.Size().Width)
	assert.Equal(t, float32(50), obj2.Size().Width)
	cell2Pos := gui.NewPos(200-theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_MiddleSpacer(t *testing.T) {
	cellSize := gui.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, layout.NewSpacer(), obj3)
	container.Resize(gui.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), gui.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, float32(50), obj1.Size().Width)
	assert.Equal(t, float32(50), obj2.Size().Width)
	cell2Pos := gui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Simple(t *testing.T) {
	cellSize := gui.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), gui.NewSize(50, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_HiddenItem(t *testing.T) {
	cellSize := gui.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj2.Hide()
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), gui.NewSize(50, 100+(theme.Padding())))

	assert.Equal(t, obj1.Size(), cellSize)
	cell3Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Wide(t *testing.T) {
	cellSize := gui.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), gui.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Tall(t *testing.T) {
	cellSize := gui.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	container.Resize(gui.NewSize(100, 308))
	assert.Equal(t, container.MinSize(), gui.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, float32(50), obj1.Size().Height)
	assert.Equal(t, float32(50), obj2.Size().Height)
	cell2Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Spacer(t *testing.T) {
	cellSize := gui.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), obj1, obj2, obj3)
	container.Resize(gui.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), gui.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, float32(50), obj1.Size().Height)
	assert.Equal(t, float32(50), obj2.Size().Height)
	cell2Pos := gui.NewPos(0, 200-theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_MiddleSpacer(t *testing.T) {
	cellSize := gui.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(gui.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := gui.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, layout.NewSpacer(), obj3)
	container.Resize(gui.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), gui.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, float32(50), obj1.Size().Height)
	assert.Equal(t, float32(50), obj2.Size().Height)
	cell2Pos := gui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := gui.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}
