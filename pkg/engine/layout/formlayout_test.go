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

func TestFormLayout(t *testing.T) {
	gridSize := gui.NewSize(125, 125)

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(gui.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(gui.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(gui.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(gui.NewSize(120, 80))

	container := &gui.Container{
		Objects: []gui.CanvasObject{label1, content1, label2, content2},
	}
	container.Resize(gridSize)

	layout.NewFormLayout().Layout(container.Objects, gridSize)

	assert.Equal(t, gui.NewSize(70, 100), label1.Size())
	assert.Equal(t, gui.NewSize(120, 100), content1.Size())
	assert.Equal(t, gui.NewSize(70, 80), label2.Size())
	assert.Equal(t, gui.NewSize(120, 80), content2.Size())
}

func TestFormLayout_Hidden(t *testing.T) {
	gridSize := gui.NewSize(190+theme.Padding(), 125)

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(gui.NewSize(70, 50))
	label1.Hide()
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(gui.NewSize(120, 100))
	content1.Hide()

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(gui.NewSize(50, 30))
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(gui.NewSize(100, 80))

	container := &gui.Container{
		Objects: []gui.CanvasObject{label1, content1, label2, content2},
	}
	container.Resize(gridSize)

	layout.NewFormLayout().Layout(container.Objects, gridSize)

	assert.Equal(t, gui.NewSize(50, 80), label2.Size())
	assert.Equal(t, gui.NewSize(140, 80), content2.Size())
	assert.Equal(t, gui.NewPos(0, 0), label2.Position())
	assert.Equal(t, gui.NewPos(50+theme.Padding(), 0), content2.Position())
}

func TestFormLayout_StretchX(t *testing.T) {
	wideSize := gui.NewSize(150, 50)

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(gui.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(gui.NewSize(50, 50))

	container := &gui.Container{
		Objects: []gui.CanvasObject{label1, content1},
	}
	container.Resize(wideSize)

	layout.NewFormLayout().Layout(container.Objects, wideSize)

	assert.Equal(t, gui.NewSize(50, 50), label1.Size())
	assert.Equal(t, gui.NewSize(wideSize.Width-50-theme.Padding(), 50), content1.Size())
}

func TestFormLayout_MinSize(t *testing.T) {

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(gui.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(gui.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(gui.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(gui.NewSize(120, 80))

	container := &gui.Container{
		Objects: []gui.CanvasObject{label1, content1, label2, content2},
	}

	l := layout.NewFormLayout()
	layoutMin := l.MinSize(container.Objects)
	expectedRowWidth := 70 + 120 + theme.Padding()
	expectedRowHeight := 100 + 80 + theme.Padding()
	assert.Equal(t, gui.NewSize(expectedRowWidth, expectedRowHeight), layoutMin)
}

func TestFormLayout_MinSize_Hidden(t *testing.T) {

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(gui.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(gui.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(gui.NewSize(70, 30))
	label2.Hide()
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(gui.NewSize(120, 80))
	content2.Hide()

	container := &gui.Container{
		Objects: []gui.CanvasObject{label1, content1, label2, content2},
	}

	l := layout.NewFormLayout()
	layoutMin := l.MinSize(container.Objects)
	expectedRowWidth := 50 + 100 + theme.Padding()
	expectedRowHeight := float32(100)
	assert.Equal(t, gui.NewSize(expectedRowWidth, expectedRowHeight), layoutMin)
}
