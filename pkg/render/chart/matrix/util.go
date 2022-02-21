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

import (
	"math"
	"strconv"
)

func minInt(values ...int) int {
	min := math.MaxInt32

	for x := 0; x < len(values); x++ {
		if values[x] < min {
			min = values[x]
		}
	}
	return min
}

func maxInt(values ...int) int {
	max := math.MinInt32

	for x := 0; x < len(values); x++ {
		if values[x] > max {
			max = values[x]
		}
	}
	return max
}

func f64s(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func roundToEpsilon(value, epsilon float64) float64 {
	return math.Nextafter(value, value)
}
