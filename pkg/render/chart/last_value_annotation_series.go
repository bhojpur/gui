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

// LastValueAnnotationSeries returns an annotation series of just the last value of a value provider.
func LastValueAnnotationSeries(innerSeries ValuesProvider, vfs ...ValueFormatter) AnnotationSeries {
	var vf ValueFormatter
	if len(vfs) > 0 {
		vf = vfs[0]
	} else if typed, isTyped := innerSeries.(ValueFormatterProvider); isTyped {
		_, vf = typed.GetValueFormatters()
	} else {
		vf = FloatValueFormatter
	}

	var lastValue Value2
	if typed, isTyped := innerSeries.(LastValuesProvider); isTyped {
		lastValue.XValue, lastValue.YValue = typed.GetLastValues()
		lastValue.Label = vf(lastValue.YValue)
	} else {
		lastValue.XValue, lastValue.YValue = innerSeries.GetValues(innerSeries.Len() - 1)
		lastValue.Label = vf(lastValue.YValue)
	}

	var seriesName string
	var seriesStyle Style
	if typed, isTyped := innerSeries.(Series); isTyped {
		seriesName = fmt.Sprintf("%s - Last Value", typed.GetName())
		seriesStyle = typed.GetStyle()
	}

	return AnnotationSeries{
		Name:        seriesName,
		Style:       seriesStyle,
		Annotations: []Value2{lastValue},
	}
}
