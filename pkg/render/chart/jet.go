package chart

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

import "github.com/bhojpur/gui/pkg/render/chart/drawing"

// Jet is a color map provider based on matlab's jet color map.
func Jet(v, vmin, vmax float64) drawing.Color {
	c := drawing.Color{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // white
	var dv float64

	if v < vmin {
		v = vmin
	}
	if v > vmax {
		v = vmax
	}
	dv = vmax - vmin

	if v < (vmin + 0.25*dv) {
		c.R = 0
		c.G = drawing.ColorChannelFromFloat(4 * (v - vmin) / dv)
	} else if v < (vmin + 0.5*dv) {
		c.R = 0
		c.B = drawing.ColorChannelFromFloat(1 + 4*(vmin+0.25*dv-v)/dv)
	} else if v < (vmin + 0.75*dv) {
		c.R = drawing.ColorChannelFromFloat(4 * (v - vmin - 0.5*dv) / dv)
		c.B = 0
	} else {
		c.G = drawing.ColorChannelFromFloat(1 + 4*(vmin+0.75*dv-v)/dv)
		c.B = 0
	}

	return c
}
