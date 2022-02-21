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

func TestConcatSeries(t *testing.T) {
	// replaced new assertions helper

	s1 := ContinuousSeries{
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}

	s2 := ContinuousSeries{
		XValues: LinearRange(11, 20.0),
		YValues: LinearRange(10.0, 1.0),
	}

	s3 := ContinuousSeries{
		XValues: LinearRange(21, 30.0),
		YValues: LinearRange(1.0, 10.0),
	}

	cs := ConcatSeries([]Series{s1, s2, s3})
	testutil.AssertEqual(t, 30, cs.Len())

	x0, y0 := cs.GetValue(0)
	testutil.AssertEqual(t, 1.0, x0)
	testutil.AssertEqual(t, 1.0, y0)

	xm, ym := cs.GetValue(19)
	testutil.AssertEqual(t, 20.0, xm)
	testutil.AssertEqual(t, 1.0, ym)

	xn, yn := cs.GetValue(29)
	testutil.AssertEqual(t, 30.0, xn)
	testutil.AssertEqual(t, 10.0, yn)
}
