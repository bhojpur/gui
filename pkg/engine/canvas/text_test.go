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
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestText_MinSize(t *testing.T) {
	text := canvas.NewText("Test", color.NRGBA{0, 0, 0, 0xff})
	min := text.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)

	text = canvas.NewText("Test2", color.NRGBA{0, 0, 0, 0xff})
	min2 := text.MinSize()
	assert.True(t, min2.Width > min.Width)
}

func TestText_MinSize_NoMultiLine(t *testing.T) {
	text := canvas.NewText("Break", color.NRGBA{0, 0, 0, 0xff})
	min := text.MinSize()

	text = canvas.NewText("Bre\nak", color.NRGBA{0, 0, 0, 0xff})
	min2 := text.MinSize()
	assert.True(t, min2.Width > min.Width)
	assert.True(t, min2.Height == min.Height)
}

func TestText_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		text  string
		align gui.TextAlign
		size  gui.Size
	}{
		"short_leading_small": {
			text:  "abc",
			align: gui.TextAlignLeading,
			size:  gui.NewSize(1, 1),
		},
		"short_leading_large": {
			text:  "abc",
			align: gui.TextAlignLeading,
			size:  gui.NewSize(500, 101),
		},
		"long_leading_small": {
			text:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			align: gui.TextAlignLeading,
			size:  gui.NewSize(1, 1),
		},
		"long_leading_large": {
			text:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			align: gui.TextAlignLeading,
			size:  gui.NewSize(500, 101),
		},
		"short_center_small": {
			text:  "abc",
			align: gui.TextAlignCenter,
			size:  gui.NewSize(1, 1),
		},
		"short_center_large": {
			text:  "abc",
			align: gui.TextAlignCenter,
			size:  gui.NewSize(500, 101),
		},
		"long_center_small": {
			text:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			align: gui.TextAlignCenter,
			size:  gui.NewSize(1, 1),
		},
		"long_center_large": {
			text:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			align: gui.TextAlignCenter,
			size:  gui.NewSize(500, 101),
		},
		"short_trailing_small": {
			text:  "abc",
			align: gui.TextAlignTrailing,
			size:  gui.NewSize(1, 1),
		},
		"short_trailing_large": {
			text:  "abc",
			align: gui.TextAlignTrailing,
			size:  gui.NewSize(500, 101),
		},
		"long_trailing_small": {
			text:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			align: gui.TextAlignTrailing,
			size:  gui.NewSize(1, 1),
		},
		"long_trailing_large": {
			text:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			align: gui.TextAlignTrailing,
			size:  gui.NewSize(500, 101),
		},
	} {
		t.Run(name, func(t *testing.T) {
			text := canvas.NewText(tt.text, theme.ForegroundColor())
			text.Alignment = tt.align

			window := test.NewWindow(text)
			window.SetPadded(false)
			window.Resize(text.MinSize().Max(tt.size))

			test.AssertImageMatches(t, "text/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
