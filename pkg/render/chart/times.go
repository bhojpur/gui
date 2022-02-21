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
	"sort"
	"time"
)

// Assert types implement interfaces.
var (
	_ Sequence       = (*Times)(nil)
	_ sort.Interface = (*Times)(nil)
)

// Times are an array of times.
// It wraps the array with methods that implement `seq.Provider`.
type Times []time.Time

// Array returns the times to an array.
func (t Times) Array() []time.Time {
	return []time.Time(t)
}

// Len returns the length of the array.
func (t Times) Len() int {
	return len(t)
}

// GetValue returns a value at an index as a time.
func (t Times) GetValue(index int) float64 {
	return ToFloat64(t[index])
}

// Swap implements sort.Interface.
func (t Times) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less implements sort.Interface.
func (t Times) Less(i, j int) bool {
	return t[i].Before(t[j])
}

// ToFloat64 returns a float64 representation of a time.
func ToFloat64(t time.Time) float64 {
	return float64(t.UnixNano())
}
