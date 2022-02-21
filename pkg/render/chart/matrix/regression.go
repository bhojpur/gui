package matrix

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

import "errors"

var (
	// ErrPolyRegArraysSameLength is a common error.
	ErrPolyRegArraysSameLength = errors.New("polynomial array inputs must be the same length")
)

// Poly returns the polynomial regress of a given degree over the given values.
func Poly(xvalues, yvalues []float64, degree int) ([]float64, error) {
	if len(xvalues) != len(yvalues) {
		return nil, ErrPolyRegArraysSameLength
	}

	m := len(yvalues)
	n := degree + 1
	y := New(m, 1, yvalues...)
	x := Zero(m, n)

	for i := 0; i < m; i++ {
		ip := float64(1)
		for j := 0; j < n; j++ {
			x.Set(i, j, ip)
			ip *= xvalues[i]
		}
	}

	q, r := x.QR()
	qty, err := q.Transpose().Times(y)
	if err != nil {
		return nil, err
	}

	c := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		c[i] = qty.Get(i, 0)
		for j := i + 1; j < n; j++ {
			c[i] -= c[j] * r.Get(i, j)
		}
		c[i] /= r.Get(i, i)
	}

	return c, nil
}
