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

type mockValuesProvider struct {
	X []float64
	Y []float64
}

func (m mockValuesProvider) Len() int {
	return MinInt(len(m.X), len(m.Y))
}

func (m mockValuesProvider) GetValues(index int) (x, y float64) {
	if index < 0 {
		panic("negative index at GetValue()")
	}
	if index >= MinInt(len(m.X), len(m.Y)) {
		panic("index is outside the length of m.X or m.Y")
	}
	x = m.X[index]
	y = m.Y[index]
	return
}

func TestSMASeriesGetValue(t *testing.T) {
	// replaced new assertions helper

	mockSeries := mockValuesProvider{
		LinearRange(1.0, 10.0),
		LinearRange(10, 1.0),
	}
	testutil.AssertEqual(t, 10, mockSeries.Len())

	mas := &SMASeries{
		InnerSeries: mockSeries,
		Period:      10,
	}

	var yvalues []float64
	for x := 0; x < mas.Len(); x++ {
		_, y := mas.GetValues(x)
		yvalues = append(yvalues, y)
	}

	testutil.AssertEqual(t, 10.0, yvalues[0])
	testutil.AssertEqual(t, 9.5, yvalues[1])
	testutil.AssertEqual(t, 9.0, yvalues[2])
	testutil.AssertEqual(t, 8.5, yvalues[3])
	testutil.AssertEqual(t, 8.0, yvalues[4])
	testutil.AssertEqual(t, 7.5, yvalues[5])
	testutil.AssertEqual(t, 7.0, yvalues[6])
	testutil.AssertEqual(t, 6.5, yvalues[7])
	testutil.AssertEqual(t, 6.0, yvalues[8])
}

func TestSMASeriesGetLastValueWindowOverlap(t *testing.T) {
	// replaced new assertions helper

	mockSeries := mockValuesProvider{
		LinearRange(1.0, 10.0),
		LinearRange(10, 1.0),
	}
	testutil.AssertEqual(t, 10, mockSeries.Len())

	mas := &SMASeries{
		InnerSeries: mockSeries,
		Period:      15,
	}

	var yvalues []float64
	for x := 0; x < mas.Len(); x++ {
		_, y := mas.GetValues(x)
		yvalues = append(yvalues, y)
	}

	lx, ly := mas.GetLastValues()
	testutil.AssertEqual(t, 10.0, lx)
	testutil.AssertEqual(t, 5.5, ly)
	testutil.AssertEqual(t, yvalues[len(yvalues)-1], ly)
}

func TestSMASeriesGetLastValue(t *testing.T) {
	// replaced new assertions helper

	mockSeries := mockValuesProvider{
		LinearRange(1.0, 100.0),
		LinearRange(100, 1.0),
	}
	testutil.AssertEqual(t, 100, mockSeries.Len())

	mas := &SMASeries{
		InnerSeries: mockSeries,
		Period:      10,
	}

	var yvalues []float64
	for x := 0; x < mas.Len(); x++ {
		_, y := mas.GetValues(x)
		yvalues = append(yvalues, y)
	}

	lx, ly := mas.GetLastValues()
	testutil.AssertEqual(t, 100.0, lx)
	testutil.AssertEqual(t, 6, ly)
	testutil.AssertEqual(t, yvalues[len(yvalues)-1], ly)
}
