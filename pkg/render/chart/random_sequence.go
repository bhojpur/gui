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
	"math"
	"math/rand"
	"time"
)

var (
	_ Sequence = (*RandomSeq)(nil)
)

// RandomValues returns an array of random values.
func RandomValues(count int) []float64 {
	return Seq{NewRandomSequence().WithLen(count)}.Values()
}

// RandomValuesWithMax returns an array of random values with a given average.
func RandomValuesWithMax(count int, max float64) []float64 {
	return Seq{NewRandomSequence().WithMax(max).WithLen(count)}.Values()
}

// NewRandomSequence creates a new random seq.
func NewRandomSequence() *RandomSeq {
	return &RandomSeq{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// RandomSeq is a random number seq generator.
type RandomSeq struct {
	rnd *rand.Rand
	max *float64
	min *float64
	len *int
}

// Len returns the number of elements that will be generated.
func (r *RandomSeq) Len() int {
	if r.len != nil {
		return *r.len
	}
	return math.MaxInt32
}

// GetValue returns the value.
func (r *RandomSeq) GetValue(_ int) float64 {
	if r.min != nil && r.max != nil {
		var delta float64

		if *r.max > *r.min {
			delta = *r.max - *r.min
		} else {
			delta = *r.min - *r.max
		}

		return *r.min + (r.rnd.Float64() * delta)
	} else if r.max != nil {
		return r.rnd.Float64() * *r.max
	} else if r.min != nil {
		return *r.min + (r.rnd.Float64())
	}
	return r.rnd.Float64()
}

// WithLen sets a maximum len
func (r *RandomSeq) WithLen(length int) *RandomSeq {
	r.len = &length
	return r
}

// Min returns the minimum value.
func (r RandomSeq) Min() *float64 {
	return r.min
}

// WithMin sets the scale and returns the Random.
func (r *RandomSeq) WithMin(min float64) *RandomSeq {
	r.min = &min
	return r
}

// Max returns the maximum value.
func (r RandomSeq) Max() *float64 {
	return r.max
}

// WithMax sets the average and returns the Random.
func (r *RandomSeq) WithMax(max float64) *RandomSeq {
	r.max = &max
	return r
}
