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
	"bytes"
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestDonutChart(t *testing.T) {
	// replaced new assertions helper

	pie := DonutChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 10, Label: "Blue"},
			{Value: 9, Label: "Green"},
			{Value: 8, Label: "Gray"},
			{Value: 7, Label: "Orange"},
			{Value: 6, Label: "HEANG"},
			{Value: 5, Label: "??"},
			{Value: 2, Label: "!!"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	pie.Render(PNG, b)
	testutil.AssertNotZero(t, b.Len())
}

func TestDonutChartDropsZeroValues(t *testing.T) {
	// replaced new assertions helper

	pie := DonutChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 0, Label: "Gray"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	testutil.AssertNil(t, err)
}

func TestDonutChartAllZeroValues(t *testing.T) {
	// replaced new assertions helper

	pie := DonutChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 0, Label: "Blue"},
			{Value: 0, Label: "Green"},
			{Value: 0, Label: "Gray"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	testutil.AssertNotNil(t, err)
}
