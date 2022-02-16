package widget_test

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
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/test"

	"github.com/stretchr/testify/assert"
)

var shadowLevel = widget.ElevationLevel(5)

func TestShadow_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(gui.NewSize(50, 50))

	s.Resize(gui.NewSize(30, 30))
	s.Move(gui.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/theme_default.png", w.Canvas().Capture())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "shadow/theme_ugly.png", w.Canvas().Capture())
}

func TestShadow_AroundShadow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(gui.NewSize(50, 50))

	s.Resize(gui.NewSize(30, 30))
	s.Move(gui.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/around.png", w.Canvas().Capture())
}

func TestShadow_Transparency(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(gui.NewSize(50, 50))

	s.Resize(gui.NewSize(40, 20))
	s.Move(gui.NewPos(5, 15))
	s2 := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w.Canvas().Overlays().Add(s2)
	s2.Resize(gui.NewSize(20, 40))
	s2.Move(gui.NewPos(15, 5))
	test.AssertImageMatches(t, "shadow/transparency.png", w.Canvas().Capture())
}

func TestShadow_BottomShadow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowBottom, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(gui.NewSize(50, 50))

	s.Resize(gui.NewSize(30, 30))
	s.Move(gui.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/bottom.png", w.Canvas().Capture())
}

func TestShadow_MinSize(t *testing.T) {
	assert.Equal(t, gui.NewSize(0, 0), widget.NewShadow(widget.ShadowAround, 1).MinSize())
}

func TestShadow_TopShadow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowTop, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(gui.NewSize(50, 50))

	s.Resize(gui.NewSize(30, 30))
	s.Move(gui.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/top.png", w.Canvas().Capture())
}
