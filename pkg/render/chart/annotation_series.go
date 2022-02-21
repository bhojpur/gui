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

// Interface Assertions.
var (
	_ Series = (*AnnotationSeries)(nil)
)

// AnnotationSeries is a series of labels on the chart.
type AnnotationSeries struct {
	Name        string
	Style       Style
	YAxis       YAxisType
	Annotations []Value2
}

// GetName returns the name of the time series.
func (as AnnotationSeries) GetName() string {
	return as.Name
}

// GetStyle returns the line style.
func (as AnnotationSeries) GetStyle() Style {
	return as.Style
}

// GetYAxis returns which YAxis the series draws on.
func (as AnnotationSeries) GetYAxis() YAxisType {
	return as.YAxis
}

func (as AnnotationSeries) annotationStyleDefaults(defaults Style) Style {
	return Style{
		FontColor:   DefaultTextColor,
		Font:        defaults.Font,
		FillColor:   DefaultAnnotationFillColor,
		FontSize:    DefaultAnnotationFontSize,
		StrokeColor: defaults.StrokeColor,
		StrokeWidth: defaults.StrokeWidth,
		Padding:     DefaultAnnotationPadding,
	}
}

// Measure returns a bounds box of the series.
func (as AnnotationSeries) Measure(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) Box {
	box := Box{
		Top:    math.MaxInt32,
		Left:   math.MaxInt32,
		Right:  0,
		Bottom: 0,
	}
	if !as.Style.Hidden {
		seriesStyle := as.Style.InheritFrom(as.annotationStyleDefaults(defaults))
		for _, a := range as.Annotations {
			style := a.Style.InheritFrom(seriesStyle)
			lx := canvasBox.Left + xrange.Translate(a.XValue)
			ly := canvasBox.Bottom - yrange.Translate(a.YValue)
			ab := Draw.MeasureAnnotation(r, canvasBox, style, lx, ly, a.Label)
			box.Top = MinInt(box.Top, ab.Top)
			box.Left = MinInt(box.Left, ab.Left)
			box.Right = MaxInt(box.Right, ab.Right)
			box.Bottom = MaxInt(box.Bottom, ab.Bottom)
		}
	}
	return box
}

// Render draws the series.
func (as AnnotationSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	if !as.Style.Hidden {
		seriesStyle := as.Style.InheritFrom(as.annotationStyleDefaults(defaults))
		for _, a := range as.Annotations {
			style := a.Style.InheritFrom(seriesStyle)
			lx := canvasBox.Left + xrange.Translate(a.XValue)
			ly := canvasBox.Bottom - yrange.Translate(a.YValue)
			Draw.Annotation(r, canvasBox, style, lx, ly, a.Label)
		}
	}
}

// Validate validates the series.
func (as AnnotationSeries) Validate() error {
	if len(as.Annotations) == 0 {
		return fmt.Errorf("annotation series requires annotations to be set and not empty")
	}
	return nil
}
