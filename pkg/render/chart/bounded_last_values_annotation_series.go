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

// BoundedLastValuesAnnotationSeries returns a last value annotation series for a bounded values provider.
func BoundedLastValuesAnnotationSeries(innerSeries FullBoundedValuesProvider, vfs ...ValueFormatter) AnnotationSeries {
	lvx, lvy1, lvy2 := innerSeries.GetBoundedLastValues()

	var vf ValueFormatter
	if len(vfs) > 0 {
		vf = vfs[0]
	} else if typed, isTyped := innerSeries.(ValueFormatterProvider); isTyped {
		_, vf = typed.GetValueFormatters()
	} else {
		vf = FloatValueFormatter
	}

	label1 := vf(lvy1)
	label2 := vf(lvy2)

	var seriesName string
	var seriesStyle Style
	if typed, isTyped := innerSeries.(Series); isTyped {
		seriesName = fmt.Sprintf("%s - Last Values", typed.GetName())
		seriesStyle = typed.GetStyle()
	}

	return AnnotationSeries{
		Name:  seriesName,
		Style: seriesStyle,
		Annotations: []Value2{
			{XValue: lvx, YValue: lvy1, Label: label1},
			{XValue: lvx, YValue: lvy2, Label: label2},
		},
	}
}
