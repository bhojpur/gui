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
	"testing"

	"github.com/bhojpur/gui/pkg/render/canvas/test"
)

func TestShapes(t *testing.T) {
	Epsilon = 0.01
	test.T(t, Rectangle(0.0, 10.0), &Path{})
	test.T(t, Rectangle(5.0, 10.0), MustParseSVG("H5V10H0z"))
	test.T(t, RoundedRectangle(0.0, 10.0, 0.0), &Path{})
	test.T(t, RoundedRectangle(5.0, 10.0, 0.0), MustParseSVG("H5V10H0z"))
	test.T(t, RoundedRectangle(5.0, 10.0, 2.0), MustParseSVG("M0 2A2 2 0 0 1 2 0L3 0A2 2 0 0 1 5 2L5 8A2 2 0 0 1 3 10L2 10A2 2 0 0 1 0 8z"))
	test.T(t, RoundedRectangle(5.0, 10.0, -2.0), MustParseSVG("M0 2A2 2 0 0 0 2 0L3 0A2 2 0 0 0 5 2L5 8A2 2 0 0 0 3 10L2 10A2 2 0 0 0 0 8z"))
	test.T(t, BeveledRectangle(0.0, 10.0, 0.0), &Path{})
	test.T(t, BeveledRectangle(5.0, 10.0, 0.0), MustParseSVG("H5V10H0z"))
	test.T(t, BeveledRectangle(5.0, 10.0, 2.0), MustParseSVG("M0 2 2 0 3 0 5 2 5 8 3 10 2 10 0 8z"))
	test.T(t, Circle(0.0), &Path{})
	test.T(t, Circle(2.0), MustParseSVG("M2 0A2 2 0 0 1 -2 0A2 2 0 0 1 2 0z"))
	test.T(t, RegularPolygon(2, 2.0, true), &Path{})
	test.T(t, RegularPolygon(4, 0.0, true), &Path{})
	test.T(t, RegularPolygon(4, 2.0, true), MustParseSVG("M0 2 -2 0 0 -2 2 0z"))
	test.T(t, RegularPolygon(3, 2.0, false), MustParseSVG("M-1.7321 1L0 -2L1.7321 1z"))
	test.T(t, StarPolygon(2, 4.0, 2.0, true), &Path{})
	test.T(t, StarPolygon(4, 4.0, 2.0, true), MustParseSVG("M0 4 -1.41 1.41 -4 0 -1.41 -1.41 0 -4 1.41 -1.41 4 0 1.41 1.41z"))
	test.T(t, StarPolygon(3, 4.0, 2.0, false), MustParseSVG("M-3.4641 2L-1.7321 -1L0 -4L1.7321 -1L3.4641 2L0 2z"))
}
