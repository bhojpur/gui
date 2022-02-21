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

// LinearCoefficientProvider is a type that returns linear cofficients.
type LinearCoefficientProvider interface {
	Coefficients() (m, b, stdev, avg float64)
}

// LinearCoefficients returns a fixed linear coefficient pair.
func LinearCoefficients(m, b float64) LinearCoefficientSet {
	return LinearCoefficientSet{
		M: m,
		B: b,
	}
}

// NormalizedLinearCoefficients returns a fixed linear coefficient pair.
func NormalizedLinearCoefficients(m, b, stdev, avg float64) LinearCoefficientSet {
	return LinearCoefficientSet{
		M:      m,
		B:      b,
		StdDev: stdev,
		Avg:    avg,
	}
}

// LinearCoefficientSet is the m and b values for the linear equation in the form:
// y = (m*x) + b
type LinearCoefficientSet struct {
	M      float64
	B      float64
	StdDev float64
	Avg    float64
}

// Coefficients returns the coefficients.
func (lcs LinearCoefficientSet) Coefficients() (m, b, stdev, avg float64) {
	m = lcs.M
	b = lcs.B
	stdev = lcs.StdDev
	avg = lcs.Avg
	return
}
