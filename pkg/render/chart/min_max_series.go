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
)

// MinSeries draws a horizontal line at the minimum value of the inner series.
type MinSeries struct {
	Name        string
	Style       Style
	YAxis       YAxisType
	InnerSeries ValuesProvider

	minValue *float64
}

// GetName returns the name of the time series.
func (ms MinSeries) GetName() string {
	return ms.Name
}

// GetStyle returns the line style.
func (ms MinSeries) GetStyle() Style {
	return ms.Style
}

// GetYAxis returns which YAxis the series draws on.
func (ms MinSeries) GetYAxis() YAxisType {
	return ms.YAxis
}

// Len returns the number of elements in the series.
func (ms MinSeries) Len() int {
	return ms.InnerSeries.Len()
}

// GetValues gets a value at a given index.
func (ms *MinSeries) GetValues(index int) (x, y float64) {
	ms.ensureMinValue()
	x, _ = ms.InnerSeries.GetValues(index)
	y = *ms.minValue
	return
}

// Render renders the series.
func (ms *MinSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	style := ms.Style.InheritFrom(defaults)
	Draw.LineSeries(r, canvasBox, xrange, yrange, style, ms)
}

func (ms *MinSeries) ensureMinValue() {
	if ms.minValue == nil {
		minValue := math.MaxFloat64
		var y float64
		for x := 0; x < ms.InnerSeries.Len(); x++ {
			_, y = ms.InnerSeries.GetValues(x)
			if y < minValue {
				minValue = y
			}
		}
		ms.minValue = &minValue
	}
}

// Validate validates the series.
func (ms *MinSeries) Validate() error {
	if ms.InnerSeries == nil {
		return fmt.Errorf("min series requires InnerSeries to be set")
	}
	return nil
}

// MaxSeries draws a horizontal line at the maximum value of the inner series.
type MaxSeries struct {
	Name        string
	Style       Style
	YAxis       YAxisType
	InnerSeries ValuesProvider

	maxValue *float64
}

// GetName returns the name of the time series.
func (ms MaxSeries) GetName() string {
	return ms.Name
}

// GetStyle returns the line style.
func (ms MaxSeries) GetStyle() Style {
	return ms.Style
}

// GetYAxis returns which YAxis the series draws on.
func (ms MaxSeries) GetYAxis() YAxisType {
	return ms.YAxis
}

// Len returns the number of elements in the series.
func (ms MaxSeries) Len() int {
	return ms.InnerSeries.Len()
}

// GetValues gets a value at a given index.
func (ms *MaxSeries) GetValues(index int) (x, y float64) {
	ms.ensureMaxValue()
	x, _ = ms.InnerSeries.GetValues(index)
	y = *ms.maxValue
	return
}

// Render renders the series.
func (ms *MaxSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	style := ms.Style.InheritFrom(defaults)
	Draw.LineSeries(r, canvasBox, xrange, yrange, style, ms)
}

func (ms *MaxSeries) ensureMaxValue() {
	if ms.maxValue == nil {
		maxValue := -math.MaxFloat64
		var y float64
		for x := 0; x < ms.InnerSeries.Len(); x++ {
			_, y = ms.InnerSeries.GetValues(x)
			if y > maxValue {
				maxValue = y
			}
		}
		ms.maxValue = &maxValue
	}
}

// Validate validates the series.
func (ms *MaxSeries) Validate() error {
	if ms.InnerSeries == nil {
		return fmt.Errorf("max series requires InnerSeries to be set")
	}
	return nil
}
