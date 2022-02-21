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
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestValuesValues(t *testing.T) {
	// replaced new assertions helper

	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).Values()
	testutil.AssertLen(t, values, 7)
	testutil.AssertEqual(t, 10, values[0])
	testutil.AssertEqual(t, 9, values[1])
	testutil.AssertEqual(t, 8, values[2])
	testutil.AssertEqual(t, 7, values[3])
	testutil.AssertEqual(t, 6, values[4])
	testutil.AssertEqual(t, 5, values[5])
	testutil.AssertEqual(t, 2, values[6])
}

func TestValuesValuesNormalized(t *testing.T) {
	// replaced new assertions helper

	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).ValuesNormalized()
	testutil.AssertLen(t, values, 7)
	testutil.AssertEqual(t, 0.2127, values[0])
	testutil.AssertEqual(t, 0.0425, values[6])
}

func TestValuesNormalize(t *testing.T) {
	// replaced new assertions helper

	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).Normalize()
	testutil.AssertLen(t, values, 7)
	testutil.AssertEqual(t, 0.2127, values[0].Value)
	testutil.AssertEqual(t, 0.0425, values[6].Value)
}
