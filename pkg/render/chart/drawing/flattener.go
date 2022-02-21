package drawing

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

// Liner receive segment definition
type Liner interface {
	// LineTo Draw a line from the current position to the point (x, y)
	LineTo(x, y float64)
}

// Flattener receive segment definition
type Flattener interface {
	// MoveTo Start a New line from the point (x, y)
	MoveTo(x, y float64)
	// LineTo Draw a line from the current position to the point (x, y)
	LineTo(x, y float64)
	// LineJoin add the most recent starting point to close the path to create a polygon
	LineJoin()
	// Close add the most recent starting point to close the path to create a polygon
	Close()
	// End mark the current line as finished so we can draw caps
	End()
}

// Flatten convert curves into straight segments keeping join segments info
func Flatten(path *Path, flattener Flattener, scale float64) {
	// First Point
	var startX, startY float64
	// Current Point
	var x, y float64
	var i int
	for _, cmp := range path.Components {
		switch cmp {
		case MoveToComponent:
			x, y = path.Points[i], path.Points[i+1]
			startX, startY = x, y
			if i != 0 {
				flattener.End()
			}
			flattener.MoveTo(x, y)
			i += 2
		case LineToComponent:
			x, y = path.Points[i], path.Points[i+1]
			flattener.LineTo(x, y)
			flattener.LineJoin()
			i += 2
		case QuadCurveToComponent:
			// we include the previous point for the start of the curve
			TraceQuad(flattener, path.Points[i-2:], 0.5)
			x, y = path.Points[i+2], path.Points[i+3]
			flattener.LineTo(x, y)
			i += 4
		case CubicCurveToComponent:
			TraceCubic(flattener, path.Points[i-2:], 0.5)
			x, y = path.Points[i+4], path.Points[i+5]
			flattener.LineTo(x, y)
			i += 6
		case ArcToComponent:
			x, y = TraceArc(flattener, path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3], path.Points[i+4], path.Points[i+5], scale)
			flattener.LineTo(x, y)
			i += 6
		case CloseComponent:
			flattener.LineTo(startX, startY)
			flattener.Close()
		}
	}
	flattener.End()
}

// SegmentedPath is a path of disparate point sectinos.
type SegmentedPath struct {
	Points []float64
}

// MoveTo implements the path interface.
func (p *SegmentedPath) MoveTo(x, y float64) {
	p.Points = append(p.Points, x, y)
	// TODO need to mark this point as moveto
}

// LineTo implements the path interface.
func (p *SegmentedPath) LineTo(x, y float64) {
	p.Points = append(p.Points, x, y)
}

// LineJoin implements the path interface.
func (p *SegmentedPath) LineJoin() {
	// TODO need to mark the current point as linejoin
}

// Close implements the path interface.
func (p *SegmentedPath) Close() {
	// TODO Close
}

// End implements the path interface.
func (p *SegmentedPath) End() {
	// Nothing to do
}
