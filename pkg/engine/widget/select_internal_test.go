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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestSelectRenderer_TapAnimation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	test.ApplyTheme(t, test.NewTheme())
	sel := NewSelect([]string{"one"}, func(s string) {})
	w := test.NewWindow(sel)
	defer w.Close()
	w.Resize(sel.MinSize().Add(gui.NewSize(10, 10)))
	sel.Resize(sel.MinSize())
	sel.Refresh()

	render1 := test.WidgetRenderer(sel).(*selectRenderer)
	test.Tap(sel)
	sel.popUp.Hide()
	sel.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "select/tap_animation.png", w.Canvas().Capture())

	cache.DestroyRenderer(sel)
	sel.Refresh()

	render2 := test.WidgetRenderer(sel).(*selectRenderer)

	assert.NotEqual(t, render1, render2)

	test.Tap(sel)
	sel.popUp.Hide()
	sel.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "select/tap_animation.png", w.Canvas().Capture())
}
