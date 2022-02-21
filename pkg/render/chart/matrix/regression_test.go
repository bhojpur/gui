package matrix

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

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestPoly(t *testing.T) {
	// replaced new assertions helper
	var xGiven = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var yGiven = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}
	var degree = 2

	c, err := Poly(xGiven, yGiven, degree)
	testutil.AssertNil(t, err)
	testutil.AssertLen(t, c, 3)

	testutil.AssertInDelta(t, c[0], 0.999999999, DefaultEpsilon)
	testutil.AssertInDelta(t, c[1], 2, DefaultEpsilon)
	testutil.AssertInDelta(t, c[2], 3, DefaultEpsilon)
}
