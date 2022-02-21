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

func TestGenerateContinuousTicks(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    0.0,
		Max:    10.0,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	testutil.AssertNotEmpty(t, ticks)
	testutil.AssertLen(t, ticks, 11)
	testutil.AssertEqual(t, 0.0, ticks[0].Value)
	testutil.AssertEqual(t, 10, ticks[len(ticks)-1].Value)
}

func TestGenerateContinuousTicksDescending(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:        0.0,
		Max:        10.0,
		Domain:     256,
		Descending: true,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	testutil.AssertNotEmpty(t, ticks)
	testutil.AssertLen(t, ticks, 11)
	testutil.AssertEqual(t, 10.0, ticks[0].Value)
	testutil.AssertEqual(t, 9.0, ticks[1].Value)
	testutil.AssertEqual(t, 1.0, ticks[len(ticks)-2].Value)
	testutil.AssertEqual(t, 0.0, ticks[len(ticks)-1].Value)
}
