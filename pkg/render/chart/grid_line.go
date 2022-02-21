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

// GridLineProvider is a type that provides grid lines.
type GridLineProvider interface {
	GetGridLines(ticks []Tick, isVertical bool, majorStyle, minorStyle Style) []GridLine
}

// GridLine is a line on a graph canvas.
type GridLine struct {
	IsMinor bool
	Style   Style
	Value   float64
}

// Major returns if the gridline is a `major` line.
func (gl GridLine) Major() bool {
	return !gl.IsMinor
}

// Minor returns if the gridline is a `minor` line.
func (gl GridLine) Minor() bool {
	return gl.IsMinor
}

// Render renders the gridline
func (gl GridLine) Render(r Renderer, canvasBox Box, ra Range, isVertical bool, defaults Style) {
	r.SetStrokeColor(gl.Style.GetStrokeColor(defaults.GetStrokeColor()))
	r.SetStrokeWidth(gl.Style.GetStrokeWidth(defaults.GetStrokeWidth()))
	r.SetStrokeDashArray(gl.Style.GetStrokeDashArray(defaults.GetStrokeDashArray()))

	if isVertical {
		lineLeft := canvasBox.Left + ra.Translate(gl.Value)
		lineBottom := canvasBox.Bottom
		lineTop := canvasBox.Top

		r.MoveTo(lineLeft, lineBottom)
		r.LineTo(lineLeft, lineTop)
		r.Stroke()
	} else {
		lineLeft := canvasBox.Left
		lineRight := canvasBox.Right
		lineHeight := canvasBox.Bottom - ra.Translate(gl.Value)

		r.MoveTo(lineLeft, lineHeight)
		r.LineTo(lineRight, lineHeight)
		r.Stroke()
	}
}

// GenerateGridLines generates grid lines.
func GenerateGridLines(ticks []Tick, majorStyle, minorStyle Style) []GridLine {
	var gl []GridLine
	isMinor := false

	if len(ticks) < 3 {
		return gl
	}

	for _, t := range ticks[1 : len(ticks)-1] {
		s := majorStyle
		if isMinor {
			s = minorStyle
		}
		gl = append(gl, GridLine{
			Style:   s,
			IsMinor: isMinor,
			Value:   t.Value,
		})
		isMinor = !isMinor
	}
	return gl
}
