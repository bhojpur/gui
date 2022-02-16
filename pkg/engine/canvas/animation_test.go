package canvas

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
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/stretchr/testify/assert"
)

func TestNewColorAnimation(t *testing.T) {
	start := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	stop := color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	var current color.Color
	anim := NewColorRGBAAnimation(start, stop, time.Second, func(c color.Color) {
		current = c
	})

	anim.Tick(0.0)
	assert.Equal(t, color.RGBA{R: 0, G: 0, B: 0, A: 0xFF}, current)
	anim.Tick(0.5)
	assert.Equal(t, color.RGBA{R: 0x7F, G: 0x7F, B: 0x7F, A: 0xFF}, current)
	anim.Tick(1.0)
	assert.Equal(t, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, current)
}

func TestNewPositionAnimation(t *testing.T) {
	start := gui.NewPos(110, 10)
	stop := gui.NewPos(10, 110)

	var current gui.Position
	anim := NewPositionAnimation(start, stop, time.Second, func(p gui.Position) {
		current = p
	})

	anim.Tick(0.0)
	assert.Equal(t, start, current)
	anim.Tick(0.5)
	assert.Equal(t, gui.NewPos(60, 60), current)
	anim.Tick(1.0)
	assert.Equal(t, stop, current)
}

func TestNewSizeAnimation(t *testing.T) {
	start := gui.NewSize(110, 10)
	stop := gui.NewSize(10, 110)

	var current gui.Size
	anim := NewSizeAnimation(start, stop, time.Second, func(s gui.Size) {
		current = s
	})

	anim.Tick(0.0)
	assert.Equal(t, start, current)
	anim.Tick(0.5)
	assert.Equal(t, gui.NewSize(60, 60), current)
	anim.Tick(1.0)
	assert.Equal(t, stop, current)
}
