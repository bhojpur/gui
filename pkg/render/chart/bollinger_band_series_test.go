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
	"fmt"
	"math"
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestBollingerBandSeries(t *testing.T) {
	// replaced new assertions helper

	s1 := mockValuesProvider{
		X: LinearRange(1.0, 100.0),
		Y: RandomValuesWithMax(100, 1024),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	xvalues := make([]float64, 100)
	y1values := make([]float64, 100)
	y2values := make([]float64, 100)

	for x := 0; x < 100; x++ {
		xvalues[x], y1values[x], y2values[x] = bbs.GetBoundedValues(x)
	}

	for x := bbs.GetPeriod(); x < 100; x++ {
		testutil.AssertTrue(t, y1values[x] > y2values[x], fmt.Sprintf("%v vs. %v", y1values[x], y2values[x]))
	}
}

func TestBollingerBandLastValue(t *testing.T) {
	// replaced new assertions helper

	s1 := mockValuesProvider{
		X: LinearRange(1.0, 100.0),
		Y: LinearRange(1.0, 100.0),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	x, y1, y2 := bbs.GetBoundedLastValues()
	testutil.AssertEqual(t, 100.0, x)
	testutil.AssertEqual(t, 101, math.Floor(y1))
	testutil.AssertEqual(t, 83, math.Floor(y2))
}
