package dialog

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

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestShowCustom_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))

	label := widget.NewLabel("Content")
	label.Alignment = gui.TextAlignCenter

	d := NewCustom("Title", "OK", label, w)
	shadowPad := float32(50)
	w.Resize(d.MinSize().Add(gui.NewSize(shadowPad, shadowPad)))

	d.Show()
	test.AssertImageMatches(t, "dialog-custom-default.png", w.Canvas().Capture())

	test.ApplyTheme(t, test.NewTheme())
	w.Resize(d.MinSize().Add(gui.NewSize(shadowPad, shadowPad)))
	d.Resize(d.MinSize()) // TODO remove once #707 is resolved
	test.AssertImageMatches(t, "dialog-custom-ugly.png", w.Canvas().Capture())
}

func TestShowCustom_Resize(t *testing.T) {
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(gui.NewSize(300, 300))

	label := widget.NewLabel("Content")
	label.Alignment = gui.TextAlignCenter
	d := NewCustom("Title", "OK", label, w)

	size := gui.NewSize(200, 200)
	d.Resize(size)
	d.Show()
	assert.Equal(t, size, d.(*dialog).win.Content.Size().Add(gui.NewSize(theme.Padding()*2, theme.Padding()*2)))
}

func TestCustom_ApplyThemeOnShow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(gui.NewSize(200, 300))

	label := widget.NewLabel("Content")
	label.Alignment = gui.TextAlignCenter
	d := NewCustom("Title", "OK", label, w)

	test.ApplyTheme(t, test.Theme())
	d.Show()
	test.AssertImageMatches(t, "dialog-onshow-theme-default.png", w.Canvas().Capture())
	d.Hide()

	test.ApplyTheme(t, test.NewTheme())
	d.Show()
	test.AssertImageMatches(t, "dialog-onshow-theme-changed.png", w.Canvas().Capture())
	d.Hide()

	test.ApplyTheme(t, test.Theme())
	d.Show()
	test.AssertImageMatches(t, "dialog-onshow-theme-default.png", w.Canvas().Capture())
	d.Hide()
}

func TestCustom_ResizeOnShow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	size := gui.NewSize(200, 300)
	w.Resize(size)

	label := widget.NewLabel("Content")
	label.Alignment = gui.TextAlignCenter
	d := NewCustom("Title", "OK", label, w).(*dialog)

	d.Show()
	assert.Equal(t, size, d.win.Size())
	d.Hide()

	size = gui.NewSize(500, 500)
	w.Resize(size)
	d.Show()
	assert.Equal(t, size, d.win.Size())
	d.Hide()
}
