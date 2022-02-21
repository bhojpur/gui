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

// ConcatSeries is a special type of series that concatenates its `InnerSeries`.
type ConcatSeries []Series

// Len returns the length of the concatenated set of series.
func (cs ConcatSeries) Len() int {
	total := 0
	for _, s := range cs {
		if typed, isValuesProvider := s.(ValuesProvider); isValuesProvider {
			total += typed.Len()
		}
	}

	return total
}

// GetValue returns the value at the (meta) index (i.e 0 => totalLen-1)
func (cs ConcatSeries) GetValue(index int) (x, y float64) {
	cursor := 0
	for _, s := range cs {
		if typed, isValuesProvider := s.(ValuesProvider); isValuesProvider {
			len := typed.Len()
			if index < cursor+len {
				x, y = typed.GetValues(index - cursor) //FENCEPOSTS.
				return
			}
			cursor += typed.Len()
		}
	}
	return
}

// Validate validates the series.
func (cs ConcatSeries) Validate() error {
	var err error
	for _, s := range cs {
		err = s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
