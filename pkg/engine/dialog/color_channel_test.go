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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/test"
)

func Test_colorChannel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	min := 0
	max := 100
	size := gui.NewSize(250, 50)

	for name, tt := range map[string]struct {
		name  string
		value int
	}{
		"foobar_0": {
			name:  "foobar",
			value: 0,
		},
		"foobar_50": {
			name:  "foobar",
			value: 50,
		},
		"foobar_100": {
			name:  "foobar",
			value: 100,
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := newColorChannel(tt.name, min, max, tt.value, nil)
			color.Resize(size)

			window := test.NewWindow(color)

			test.AssertImageMatches(t, "color/channel_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
