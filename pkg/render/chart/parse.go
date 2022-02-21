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
	"strconv"
	"strings"
	"time"
)

// ParseFloats parses a list of floats.
func ParseFloats(values ...string) ([]float64, error) {
	var output []float64
	var parsedValue float64
	var err error
	var cleaned string
	for _, value := range values {
		cleaned = strings.TrimSpace(strings.Replace(value, ",", "", -1))
		if cleaned == "" {
			continue
		}
		if parsedValue, err = strconv.ParseFloat(cleaned, 64); err != nil {
			return nil, err
		}
		output = append(output, parsedValue)
	}
	return output, nil
}

// ParseTimes parses a list of times with a given format.
func ParseTimes(layout string, values ...string) ([]time.Time, error) {
	var output []time.Time
	var parsedValue time.Time
	var err error
	for _, value := range values {
		if parsedValue, err = time.Parse(layout, value); err != nil {
			return nil, err
		}
		output = append(output, parsedValue)
	}
	return output, nil
}
