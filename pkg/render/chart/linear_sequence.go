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

// LinearRange returns an array of values representing the range from start to end, incremented by 1.0.
func LinearRange(start, end float64) []float64 {
	return Seq{NewLinearSequence().WithStart(start).WithEnd(end).WithStep(1.0)}.Values()
}

// LinearRangeWithStep returns the array values of a linear seq with a given start, end and optional step.
func LinearRangeWithStep(start, end, step float64) []float64 {
	return Seq{NewLinearSequence().WithStart(start).WithEnd(end).WithStep(step)}.Values()
}

// NewLinearSequence returns a new linear generator.
func NewLinearSequence() *LinearSeq {
	return &LinearSeq{step: 1.0}
}

// LinearSeq is a stepwise generator.
type LinearSeq struct {
	start float64
	end   float64
	step  float64
}

// Start returns the start value.
func (lg LinearSeq) Start() float64 {
	return lg.start
}

// End returns the end value.
func (lg LinearSeq) End() float64 {
	return lg.end
}

// Step returns the step value.
func (lg LinearSeq) Step() float64 {
	return lg.step
}

// Len returns the number of elements in the seq.
func (lg LinearSeq) Len() int {
	if lg.start < lg.end {
		return int((lg.end-lg.start)/lg.step) + 1
	}
	return int((lg.start-lg.end)/lg.step) + 1
}

// GetValue returns the value at a given index.
func (lg LinearSeq) GetValue(index int) float64 {
	fi := float64(index)
	if lg.start < lg.end {
		return lg.start + (fi * lg.step)
	}
	return lg.start - (fi * lg.step)
}

// WithStart sets the start and returns the linear generator.
func (lg *LinearSeq) WithStart(start float64) *LinearSeq {
	lg.start = start
	return lg
}

// WithEnd sets the end and returns the linear generator.
func (lg *LinearSeq) WithEnd(end float64) *LinearSeq {
	lg.end = end
	return lg
}

// WithStep sets the step and returns the linear generator.
func (lg *LinearSeq) WithStep(step float64) *LinearSeq {
	lg.step = step
	return lg
}
