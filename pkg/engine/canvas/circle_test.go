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

func TestCircle_MinSize(t *testing.T) {
	circle := canvas.NewCircle(color.Black)
	min := circle.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestCircle_FillColor(t *testing.T) {
	c := color.White
	circle := canvas.NewCircle(c)

	assert.Equal(t, c, circle.FillColor)
}

func TestCircle_Resize(t *testing.T) {
	targetWidth := float32(50)
	targetHeight := float32(50)
	circle := canvas.NewCircle(color.White)
	start := circle.Size()
	assert.True(t, start.Height == 0)
	assert.True(t, start.Width == 0)

	circle.Resize(gui.NewSize(targetWidth, targetHeight))
	target := circle.Size()
	assert.True(t, target.Height == targetHeight)
	assert.True(t, target.Width == targetWidth)
}

func TestCircle_Move(t *testing.T) {
	circle := canvas.NewCircle(color.White)
	circle.Resize(gui.NewSize(50, 50))

	start := gui.Position{X: 0, Y: 0}
	assert.True(t, circle.Position() == start)

	target := gui.Position{X: 10, Y: 75}
	circle.Move(target)
	assert.True(t, circle.Position() == target)
}
