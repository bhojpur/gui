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
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestContinuousSeries(t *testing.T) {
	// replaced new assertions helper

	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}

	testutil.AssertEqual(t, "Test Series", cs.GetName())
	testutil.AssertEqual(t, 10, cs.Len())
	x0, y0 := cs.GetValues(0)
	testutil.AssertEqual(t, 1.0, x0)
	testutil.AssertEqual(t, 1.0, y0)

	xn, yn := cs.GetValues(9)
	testutil.AssertEqual(t, 10.0, xn)
	testutil.AssertEqual(t, 10.0, yn)

	xn, yn = cs.GetLastValues()
	testutil.AssertEqual(t, 10.0, xn)
	testutil.AssertEqual(t, 10.0, yn)
}

func TestContinuousSeriesValueFormatter(t *testing.T) {
	// replaced new assertions helper

	cs := ContinuousSeries{
		XValueFormatter: func(v interface{}) string {
			return fmt.Sprintf("%f foo", v)
		},
		YValueFormatter: func(v interface{}) string {
			return fmt.Sprintf("%f bar", v)
		},
	}

	xf, yf := cs.GetValueFormatters()
	testutil.AssertEqual(t, "0.100000 foo", xf(0.1))
	testutil.AssertEqual(t, "0.100000 bar", yf(0.1))
}

func TestContinuousSeriesValidate(t *testing.T) {
	// replaced new assertions helper

	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}
	testutil.AssertNil(t, cs.Validate())

	cs = ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 10.0),
	}
	testutil.AssertNotNil(t, cs.Validate())

	cs = ContinuousSeries{
		Name:    "Test Series",
		YValues: LinearRange(1.0, 10.0),
	}
	testutil.AssertNotNil(t, cs.Validate())
}
