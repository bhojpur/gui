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
	"time"
)

// Interface Assertions.
var (
	_ Series                 = (*TimeSeries)(nil)
	_ FirstValuesProvider    = (*TimeSeries)(nil)
	_ LastValuesProvider     = (*TimeSeries)(nil)
	_ ValueFormatterProvider = (*TimeSeries)(nil)
)

// TimeSeries is a line on a chart.
type TimeSeries struct {
	Name  string
	Style Style

	YAxis YAxisType

	XValues []time.Time
	YValues []float64
}

// GetName returns the name of the time series.
func (ts TimeSeries) GetName() string {
	return ts.Name
}

// GetStyle returns the line style.
func (ts TimeSeries) GetStyle() Style {
	return ts.Style
}

// Len returns the number of elements in the series.
func (ts TimeSeries) Len() int {
	return len(ts.XValues)
}

// GetValues gets x, y values at a given index.
func (ts TimeSeries) GetValues(index int) (x, y float64) {
	x = TimeToFloat64(ts.XValues[index])
	y = ts.YValues[index]
	return
}

// GetFirstValues gets the first values.
func (ts TimeSeries) GetFirstValues() (x, y float64) {
	x = TimeToFloat64(ts.XValues[0])
	y = ts.YValues[0]
	return
}

// GetLastValues gets the last values.
func (ts TimeSeries) GetLastValues() (x, y float64) {
	x = TimeToFloat64(ts.XValues[len(ts.XValues)-1])
	y = ts.YValues[len(ts.YValues)-1]
	return
}

// GetValueFormatters returns value formatter defaults for the series.
func (ts TimeSeries) GetValueFormatters() (x, y ValueFormatter) {
	x = TimeValueFormatter
	y = FloatValueFormatter
	return
}

// GetYAxis returns which YAxis the series draws on.
func (ts TimeSeries) GetYAxis() YAxisType {
	return ts.YAxis
}

// Render renders the series.
func (ts TimeSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	style := ts.Style.InheritFrom(defaults)
	Draw.LineSeries(r, canvasBox, xrange, yrange, style, ts)
}

// Validate validates the series.
func (ts TimeSeries) Validate() error {
	if len(ts.XValues) == 0 {
		return fmt.Errorf("time series must have xvalues set")
	}

	if len(ts.YValues) == 0 {
		return fmt.Errorf("time series must have yvalues set")
	}
	return nil
}
