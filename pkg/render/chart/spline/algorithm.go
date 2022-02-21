package spline

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
	"sort"
)

// solves diagonal matrix using Thomas algorithm
// It will change input slices
// the matrix should have full rank
func triThomas(a, b, c, d []float64) []float64 {
	n := len(b)
	if len(a)+1 != n || n != len(c)+1 || n != len(d) {
		panic("invalid input slices")
	}

	d[0] /= b[0]
	if n == 1 {
		return d
	}

	c[0] /= b[0]
	for i := 1; i < n-1; i++ {
		div := b[i] - a[i-1]*c[i-1]
		c[i] /= div
		d[i] = (d[i] - a[i-1]*d[i-1]) / div
	}
	d[n-1] = (d[n-1] - a[n-2]*d[n-2]) / (b[n-1] - a[n-2]*c[n-2])
	for i := n - 2; i >= 0; i-- {
		d[i] -= c[i] * d[i+1]
	}
	return d
}

// Find the segments between the elements in xs in which x resides
// The numbers in xs *must* be in ascending order
// The segments are left inclusive and right exclusive
//
// For example:
//
// findSegment([1, 2, 3, 4, 5], 3) = 2
// since x = 3 is in [3, 4), which is the third segment
//
// findSegment([1, 2, 3, 4, 5], 4.5) = 3
// since x = 4.5 is in [4, 5), which is the fourth segment
//
// If x is not in any of the segments, return the closest one
func findSegment(xs []float64, x float64) int {
	// assert xs in ascending order
	if x <= xs[0] {
		return 0
	}
	if l := len(xs); x >= xs[l-1] {
		return l - 2
	}
	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
}
