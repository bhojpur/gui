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

func TestLinearRegressionSeries(t *testing.T) {
	// replaced new assertions helper

	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: LinearRange(1.0, 100.0),
		YValues: LinearRange(1.0, 100.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	lrx0, lry0 := linRegSeries.GetValues(0)
	testutil.AssertInDelta(t, 1.0, lrx0, 0.0000001)
	testutil.AssertInDelta(t, 1.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	testutil.AssertInDelta(t, 100.0, lrxn, 0.0000001)
	testutil.AssertInDelta(t, 100.0, lryn, 0.0000001)
}

func TestLinearRegressionSeriesDesc(t *testing.T) {
	// replaced new assertions helper

	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: LinearRange(100.0, 1.0),
		YValues: LinearRange(100.0, 1.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	lrx0, lry0 := linRegSeries.GetValues(0)
	testutil.AssertInDelta(t, 100.0, lrx0, 0.0000001)
	testutil.AssertInDelta(t, 100.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	testutil.AssertInDelta(t, 1.0, lrxn, 0.0000001)
	testutil.AssertInDelta(t, 1.0, lryn, 0.0000001)
}

func TestLinearRegressionSeriesWindowAndOffset(t *testing.T) {
	// replaced new assertions helper

	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: LinearRange(100.0, 1.0),
		YValues: LinearRange(100.0, 1.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
		Offset:      10,
		Limit:       10,
	}

	testutil.AssertEqual(t, 10, linRegSeries.Len())

	lrx0, lry0 := linRegSeries.GetValues(0)
	testutil.AssertInDelta(t, 90.0, lrx0, 0.0000001)
	testutil.AssertInDelta(t, 90.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	testutil.AssertInDelta(t, 80.0, lrxn, 0.0000001)
	testutil.AssertInDelta(t, 80.0, lryn, 0.0000001)
}
