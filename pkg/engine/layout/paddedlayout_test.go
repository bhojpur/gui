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

func TestPaddedLayout(t *testing.T) {
	size := gui.NewSize(100, 100)

	obj := canvas.NewRectangle(color.Black)
	container := &gui.Container{
		Objects: []gui.CanvasObject{obj},
	}
	container.Resize(size)

	layout.NewPaddedLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size().Width, size.Width-2*theme.Padding())
	assert.Equal(t, obj.Size().Height, size.Height-2*theme.Padding())
}

func TestPaddedLayout_MinSize(t *testing.T) {
	text := canvas.NewText("FooBar", color.Black)
	minSize := text.MinSize()

	container := gui.NewContainer(text)
	layoutMin := layout.NewPaddedLayout().MinSize(container.Objects)

	assert.Equal(t, minSize.Width+2*theme.Padding(), layoutMin.Width)
	assert.Equal(t, minSize.Height+2*theme.Padding(), layoutMin.Height)
}

func TestContainer_PaddedLayout_MinSize(t *testing.T) {
	text := canvas.NewText("FooBar", color.Black)
	minSize := text.MinSize()

	container := gui.NewContainer(text)
	container.Layout = layout.NewPaddedLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize.Width+2*theme.Padding(), layoutMin.Width)
	assert.Equal(t, minSize.Height+2*theme.Padding(), layoutMin.Height)
}
