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

import "fmt"

// HistogramSeries is a special type of series that draws as a histogram.
// Some peculiarities; it will always be lower bounded at 0 (at the very least).
// This may alter ranges a bit and generally you want to put a histogram series on it's own y-axis.
type HistogramSeries struct {
	Name        string
	Style       Style
	YAxis       YAxisType
	InnerSeries ValuesProvider
}

// GetName implements Series.GetName.
func (hs HistogramSeries) GetName() string {
	return hs.Name
}

// GetStyle implements Series.GetStyle.
func (hs HistogramSeries) GetStyle() Style {
	return hs.Style
}

// GetYAxis returns which yaxis the series is mapped to.
func (hs HistogramSeries) GetYAxis() YAxisType {
	return hs.YAxis
}

// Len implements BoundedValuesProvider.Len.
func (hs HistogramSeries) Len() int {
	return hs.InnerSeries.Len()
}

// GetValues implements ValuesProvider.GetValues.
func (hs HistogramSeries) GetValues(index int) (x, y float64) {
	return hs.InnerSeries.GetValues(index)
}

// GetBoundedValues implements BoundedValuesProvider.GetBoundedValue
func (hs HistogramSeries) GetBoundedValues(index int) (x, y1, y2 float64) {
	vx, vy := hs.InnerSeries.GetValues(index)

	x = vx

	if vy > 0 {
		y1 = vy
		return
	}

	y2 = vy
	return
}

// Render implements Series.Render.
func (hs HistogramSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	style := hs.Style.InheritFrom(defaults)
	Draw.HistogramSeries(r, canvasBox, xrange, yrange, style, hs)
}

// Validate validates the series.
func (hs HistogramSeries) Validate() error {
	if hs.InnerSeries == nil {
		return fmt.Errorf("histogram series requires InnerSeries to be set")
	}
	return nil
}
