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
	"testing"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestProgressBarInfinite_Creation(t *testing.T) {
	bar := NewProgressBarInfinite()
	// ticker should start automatically
	assert.True(t, bar.Running())
}

func TestProgressBarInfinite_Destroy(t *testing.T) {
	bar := NewProgressBarInfinite()
	assert.True(t, cache.IsRendered(bar))
	assert.True(t, bar.Running())

	// check that it stopped
	cache.DestroyRenderer(bar)
	assert.False(t, bar.Running())

	// and that the cache was removed
	assert.False(t, cache.IsRendered(bar))
}

func TestProgressBarInfinite_Reshown(t *testing.T) {
	bar := NewProgressBarInfinite()

	assert.True(t, bar.Running())
	bar.Hide()
	assert.False(t, bar.Running())

	// make sure it restarts when re-shown
	bar.Show()
	// Show() starts a goroutine, so pause for it to initialize
	time.Sleep(10 * time.Millisecond)
	assert.True(t, bar.Running())
	bar.Hide()
	assert.False(t, bar.Running())
}

func TestInfiniteProgressRenderer_Layout(t *testing.T) {
	bar := NewProgressBarInfinite()
	width := float32(100.0)
	bar.Resize(gui.NewSize(width, 10))

	render := test.WidgetRenderer(bar).(*infProgressRenderer)

	render.updateBar(0.0)
	// start at the smallest size
	assert.Equal(t, width*minProgressBarInfiniteWidthRatio, render.bar.Size().Width)

	// make sure the inner progress bar grows in size
	// call updateBar() enough times to grow the inner bar
	maxWidth := width * maxProgressBarInfiniteWidthRatio
	render.updateBar(0.5)
	assert.Equal(t, maxWidth, render.bar.Size().Width)

	render.updateBar(1.0)
	// ends at the smallest size again
	assert.Equal(t, width*minProgressBarInfiniteWidthRatio, render.bar.Size().Width)

	bar.Hide()
}
