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

import (
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestRangeTranslate(t *testing.T) {
	// replaced new assertions helper
	values := []float64{1.0, 2.0, 2.5, 2.7, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	r := ContinuousRange{Domain: 1000}
	r.Min, r.Max = MinMax(values...)

	// delta = ~7.0
	// value = ~5.0
	// domain = ~1000
	// 5/8 * 1000 ~=
	testutil.AssertEqual(t, 0, r.Translate(1.0))
	testutil.AssertEqual(t, 1000, r.Translate(8.0))
	testutil.AssertEqual(t, 572, r.Translate(5.0))
}
