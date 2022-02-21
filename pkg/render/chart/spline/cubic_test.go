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

func TestSimpleCubicSpline(t *testing.T) {
	x := [...]float64{0, 1, 2, 3}
	y := [...]float64{0, 0.5, 2, 1.5}
	s := newSpline(x[:], y[:], CubicSecondDeriv, 0, 0)
	for i := range x {
		if !floatEquals(s.At(x[i]), y[i]) {
			t.Errorf("expected f(%g) = %g, but the result is %g", x[i], y[i], s.At(x[i]))
		}
	}
	if len(s.Range(0, 1, 0.25)) != 5 {
		t.Error("s.Range(0, 1, 0.25) should have 5 elements, but is", s.Range(0, 1, 0.25))
	}
	if len(s.Range(0, 1, 0.3)) != 4 {
		t.Error("s.Range(0, 1, 0.3) should have 4 elements, but is", s.Range(0, 1, 0.3))
	}
}
