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
	"image/color"
	"runtime"
	"testing"
	"time"

	"github.com/bhojpur/gui/pkg/engine/canvas"
	col "github.com/bhojpur/gui/pkg/engine/internal/color"
	_ "github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/stretchr/testify/assert"
)

func TestEntryCursorAnim(t *testing.T) {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := col.ToNRGBA(theme.PrimaryColor())
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}

	alphaEquals := func(color1, color2 color.Color) bool {
		_, _, _, a1 := col.ToNRGBA(color1)
		_, _, _, a2 := col.ToNRGBA(color2)
		return a1 == a2
	}

	cursor := canvas.NewRectangle(color.Black)
	a := newEntryCursorAnimation(cursor)

	a.start()
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorDim, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))

	a.interrupt()
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(0.5)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))

	a.timeNow = func() time.Time {
		return time.Now().Add(cursorInterruptTime)
	}
	// animation should be restarted inverting the colors
	a.anim.Tick(0.0)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond) // ensure go routine for restart animation is executed
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorDim, a.cursor.FillColor))

	a.timeNow = time.Now
	a.interrupt()
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))

	a.timeNow = func() time.Time {
		return time.Now().Add(cursorInterruptTime)
	}
	a.anim.Tick(0.0)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond) // ensure go routine for restart animation is executed
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorDim, a.cursor.FillColor))

	a.stop()
	assert.Nil(t, a.anim)
}
