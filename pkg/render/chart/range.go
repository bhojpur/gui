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

// NameProvider is a type that returns a name.
type NameProvider interface {
	GetName() string
}

// StyleProvider is a type that returns a style.
type StyleProvider interface {
	GetStyle() Style
}

// IsZeroable is a type that returns if it's been set or not.
type IsZeroable interface {
	IsZero() bool
}

// Stringable is a type that has a string representation.
type Stringable interface {
	String() string
}

// Range is a common interface for a range of values.
type Range interface {
	Stringable
	IsZeroable

	GetMin() float64
	SetMin(min float64)

	GetMax() float64
	SetMax(max float64)

	GetDelta() float64

	GetDomain() int
	SetDomain(domain int)

	IsDescending() bool

	// Translate the range to the domain.
	Translate(value float64) int
}
