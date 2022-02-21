package spline

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
)

func TestHermiteSpline(t *testing.T) {
	x := []float64{0, 0.16, 0.42, 0.6425, 0.8575}
	y := []float64{0, 32, 237, 255, 0}
	s := NewMonotoneSpline(x, y)

	// Test the input points are mapped exactly.
	for i := range x {
		if at := s.At(x[i]); !floatEquals(at, y[i]) {
			t.Errorf("interpolate incorrect at %f: %f != %f", x[i], at, y[i])
		}
	}

	const intermediate = 10000

	// Test that intermediate points are monotonic.
	for i := range x[:len(x)-1] {
		diff := y[i+1] - y[i]
		ypos := y[i]
		for j := 0; j < intermediate; j++ {
			xval := x[i] + (x[i+1]-x[i])*float64(j)/intermediate
			yval := s.At(xval)
			jdiff := yval - ypos

			if jdiff*diff < 0 {
				t.Errorf("not monotone at x=%f y=%f", xval, yval)
			}
			ypos = yval
		}
	}
}
