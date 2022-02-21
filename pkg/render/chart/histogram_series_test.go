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

func TestHistogramSeries(t *testing.T) {
	// replaced new assertions helper

	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 20.0),
		YValues: LinearRange(10.0, -10.0),
	}

	hs := HistogramSeries{
		InnerSeries: cs,
	}

	for x := 0; x < hs.Len(); x++ {
		csx, csy := cs.GetValues(0)
		hsx, hsy1, hsy2 := hs.GetBoundedValues(0)
		testutil.AssertEqual(t, csx, hsx)
		testutil.AssertTrue(t, hsy1 > 0)
		testutil.AssertTrue(t, hsy2 <= 0)
		testutil.AssertTrue(t, csy < 0 || (csy > 0 && csy == hsy1))
		testutil.AssertTrue(t, csy > 0 || (csy < 0 && csy == hsy2))
	}
}
