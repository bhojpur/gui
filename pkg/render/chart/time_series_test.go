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
	"time"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestTimeSeriesGetValue(t *testing.T) {
	// replaced new assertions helper

	ts := TimeSeries{
		Name: "Test",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}

	x0, y0 := ts.GetValues(0)
	testutil.AssertNotZero(t, x0)
	testutil.AssertEqual(t, 1.0, y0)
}

func TestTimeSeriesValidate(t *testing.T) {
	// replaced new assertions helper

	cs := TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	testutil.AssertNil(t, cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
	}
	testutil.AssertNotNil(t, cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	testutil.AssertNotNil(t, cs.Validate())
}
