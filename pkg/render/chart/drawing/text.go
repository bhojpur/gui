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

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

// DrawContour draws the given closed contour at the given sub-pixel offset.
func DrawContour(path PathBuilder, ps []truetype.Point, dx, dy float64) {
	if len(ps) == 0 {
		return
	}
	startX, startY := pointToF64Point(ps[0])
	path.MoveTo(startX+dx, startY+dy)
	q0X, q0Y, on0 := startX, startY, true
	for _, p := range ps[1:] {
		qX, qY := pointToF64Point(p)
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				path.LineTo(qX+dx, qY+dy)
			} else {
				path.QuadCurveTo(q0X+dx, q0Y+dy, qX+dx, qY+dy)
			}
		} else if !on0 {
			midX := (q0X + qX) / 2
			midY := (q0Y + qY) / 2
			path.QuadCurveTo(q0X+dx, q0Y+dy, midX+dx, midY+dy)
		}
		q0X, q0Y, on0 = qX, qY, on
	}
	// Close the curve.
	if on0 {
		path.LineTo(startX+dx, startY+dy)
	} else {
		path.QuadCurveTo(q0X+dx, q0Y+dy, startX+dx, startY+dy)
	}
}

// FontExtents contains font metric information.
type FontExtents struct {
	// Ascent is the distance that the text
	// extends above the baseline.
	Ascent float64

	// Descent is the distance that the text
	// extends below the baseline.  The descent
	// is given as a negative value.
	Descent float64

	// Height is the distance from the lowest
	// descending point to the highest ascending
	// point.
	Height float64
}

// Extents returns the FontExtents for a font.
// TODO needs to read this https://developer.apple.com/fonts/TrueType-Reference-Manual/RM02/Chap2.html#intro
func Extents(font *truetype.Font, size float64) FontExtents {
	bounds := font.Bounds(fixed.Int26_6(font.FUnitsPerEm()))
	scale := size / float64(font.FUnitsPerEm())
	return FontExtents{
		Ascent:  float64(bounds.Max.Y) * scale,
		Descent: float64(bounds.Min.Y) * scale,
		Height:  float64(bounds.Max.Y-bounds.Min.Y) * scale,
	}
}
