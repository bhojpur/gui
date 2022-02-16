package canvas_test

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

	"github.com/stretchr/testify/assert"
)

func TestLine_MinSize(t *testing.T) {
	line := canvas.NewLine(color.Black)
	min := line.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestLine_Move(t *testing.T) {
	line := canvas.NewLine(color.Black)

	line.Resize(gui.NewSize(10, 10))
	assert.Equal(t, gui.NewPos(0, 0), line.Position1)
	assert.Equal(t, gui.NewPos(10, 10), line.Position2)

	line.Move(gui.NewPos(5, 5))
	assert.Equal(t, gui.NewPos(5, 5), line.Position1)
	assert.Equal(t, gui.NewPos(15, 15), line.Position2)

	// rotate
	line.Position1 = gui.NewPos(0, 10)
	line.Position2 = gui.NewPos(10, 0)
	line.Move(gui.NewPos(10, 10))
	assert.Equal(t, gui.NewPos(10, 20), line.Position1)
	assert.Equal(t, gui.NewPos(20, 10), line.Position2)
}

func TestLine_Resize(t *testing.T) {
	line := canvas.NewLine(color.Black)

	line.Resize(gui.NewSize(10, 0))
	size := line.Size()

	assert.Equal(t, float32(10), size.Width)
	assert.Equal(t, float32(0), size.Height)

	// rotate
	line.Position1 = gui.NewPos(0, 10)
	line.Position2 = gui.NewPos(10, 0)
	line.Resize(gui.NewSize(20, 20))
	assert.Equal(t, gui.NewPos(0, 20), line.Position1)
	assert.Equal(t, gui.NewPos(20, 0), line.Position2)
}

func TestLine_StrokeColor(t *testing.T) {
	c := color.White
	line := canvas.NewLine(c)

	assert.Equal(t, c, line.StrokeColor)
}
