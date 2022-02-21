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

func TestPolyline(t *testing.T) {
	p := &Polyline{}
	p.Add(10, 0)
	p.Add(20, 10)
	test.T(t, len(p.Coords()), 2)
	test.T(t, p.Coords()[0], Point{10, 0})
	test.T(t, p.Coords()[1], Point{20, 10})

	test.T(t, (&Polyline{}).ToPath(), MustParseSVG(""))
	test.T(t, (&Polyline{}).Add(10, 0).ToPath(), MustParseSVG(""))
	test.T(t, (&Polyline{}).Add(10, 0).Add(20, 10).ToPath(), MustParseSVG("M10 0L20 10"))
	test.T(t, (&Polyline{}).Add(10, 0).Add(20, 10).Add(10, 0).ToPath(), MustParseSVG("M10 0L20 10z"))

	test.That(t, (&Polyline{}).Add(10, 0).Add(20, 10).Add(10, 10).Add(10, 0).Interior(12, 5, NonZero))
	test.That(t, !(&Polyline{}).Add(10, 0).Add(20, 10).Add(10, 10).Add(10, 0).Interior(5, 5, NonZero))

	test.That(t, (&Polyline{}).Add(10, 0).Add(20, 10).Add(10, 10).Add(10, 0).Interior(12, 5, EvenOdd))
	test.That(t, !(&Polyline{}).Add(10, 0).Add(20, 10).Add(10, 10).Add(10, 0).Interior(5, 5, EvenOdd))
}

func TestPolylineSmoothen(t *testing.T) {
	test.T(t, (&Polyline{}).Smoothen(), MustParseSVG(""))
	test.T(t, (&Polyline{}).Add(0, 0).Add(10, 0).Smoothen(), MustParseSVG("M0 0L10 0"))
	test.T(t, (&Polyline{}).Add(0, 0).Add(5, 10).Add(10, 0).Add(5, -10).Smoothen(), MustParseSVG("M0 0C1.444444 5.111111 2.888889 10.22222 5 10C7.111111 9.777778 9.888889 4.222222 10 0C10.11111 -4.222222 7.555556 -7.111111 5 -10"))
	test.T(t, (&Polyline{}).Add(0, 0).Add(5, 10).Add(10, 0).Add(5, -10).Add(0, 0).Smoothen(), MustParseSVG("M0 0C0 5 2.5 10 5 10C7.5 10 10 5 10 0C10 -5 7.5 -10 5 -10C2.5 -10 0 -5 0 0z"))
}
