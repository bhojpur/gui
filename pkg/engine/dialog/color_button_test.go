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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

func Test_colorButton_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		color   color.Color
		hovered bool
	}{
		"primary": {
			color: theme.PrimaryColor(),
		},
		"primary_hovered": {
			color:   theme.PrimaryColor(),
			hovered: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := newColorButton(tt.color, nil)

			if tt.hovered {
				color.MouseIn(nil)
			}

			window := test.NewWindow(container.NewCenter(color))
			window.Resize(color.MinSize().Max(gui.NewSize(50, 50)))

			test.AssertImageMatches(t, "color/button_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
