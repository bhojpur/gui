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

import "github.com/bhojpur/gui/pkg/render/chart/drawing"

// ValuesProvider is a type that produces values.
type ValuesProvider interface {
	Len() int
	GetValues(index int) (float64, float64)
}

// BoundedValuesProvider allows series to return a range.
type BoundedValuesProvider interface {
	Len() int
	GetBoundedValues(index int) (x, y1, y2 float64)
}

// FirstValuesProvider is a special type of value provider that can return it's (potentially computed) first value.
type FirstValuesProvider interface {
	GetFirstValues() (x, y float64)
}

// LastValuesProvider is a special type of value provider that can return it's (potentially computed) last value.
type LastValuesProvider interface {
	GetLastValues() (x, y float64)
}

// BoundedLastValuesProvider is a special type of value provider that can return it's (potentially computed) bounded last value.
type BoundedLastValuesProvider interface {
	GetBoundedLastValues() (x, y1, y2 float64)
}

// FullValuesProvider is an interface that combines `ValuesProvider` and `LastValuesProvider`
type FullValuesProvider interface {
	ValuesProvider
	LastValuesProvider
}

// FullBoundedValuesProvider is an interface that combines `BoundedValuesProvider` and `BoundedLastValuesProvider`
type FullBoundedValuesProvider interface {
	BoundedValuesProvider
	BoundedLastValuesProvider
}

// SizeProvider is a provider for integer size.
type SizeProvider func(xrange, yrange Range, index int, x, y float64) float64

// ColorProvider is a general provider for color ranges based on values.
type ColorProvider func(v, vmin, vmax float64) drawing.Color

// DotColorProvider is a provider for dot color.
type DotColorProvider func(xrange, yrange Range, index int, x, y float64) drawing.Color
