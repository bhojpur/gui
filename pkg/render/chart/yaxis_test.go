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

func TestYAxisGetTicks(t *testing.T) {
	// replaced new assertions helper

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	ya := YAxis{}
	yr := &ContinuousRange{Min: 10, Max: 100, Domain: 1024}
	styleDefaults := Style{
		Font:     f,
		FontSize: 10.0,
	}
	vf := FloatValueFormatter
	ticks := ya.GetTicks(r, yr, styleDefaults, vf)
	testutil.AssertLen(t, ticks, 32)
}

func TestYAxisGetTicksWithUserDefaults(t *testing.T) {
	// replaced new assertions helper

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	ya := YAxis{
		Ticks: []Tick{{Value: 1.0, Label: "1.0"}},
	}
	yr := &ContinuousRange{Min: 10, Max: 100, Domain: 1024}
	styleDefaults := Style{
		Font:     f,
		FontSize: 10.0,
	}
	vf := FloatValueFormatter
	ticks := ya.GetTicks(r, yr, styleDefaults, vf)
	testutil.AssertLen(t, ticks, 1)
}

func TestYAxisMeasure(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)
	style := Style{
		Font:     f,
		FontSize: 10.0,
	}
	r, err := PNG(100, 100)
	testutil.AssertNil(t, err)
	ticks := []Tick{{Value: 1.0, Label: "1.0"}, {Value: 2.0, Label: "2.0"}, {Value: 3.0, Label: "3.0"}}
	ya := YAxis{}
	yab := ya.Measure(r, NewBox(0, 0, 100, 100), &ContinuousRange{Min: 1.0, Max: 3.0, Domain: 100}, style, ticks)
	testutil.AssertEqual(t, 32, yab.Width())
	testutil.AssertEqual(t, 110, yab.Height())
}

func TestYAxisSecondaryMeasure(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)
	style := Style{
		Font:     f,
		FontSize: 10.0,
	}
	r, err := PNG(100, 100)
	testutil.AssertNil(t, err)
	ticks := []Tick{{Value: 1.0, Label: "1.0"}, {Value: 2.0, Label: "2.0"}, {Value: 3.0, Label: "3.0"}}
	ya := YAxis{AxisType: YAxisSecondary}
	yab := ya.Measure(r, NewBox(0, 0, 100, 100), &ContinuousRange{Min: 1.0, Max: 3.0, Domain: 100}, style, ticks)
	testutil.AssertEqual(t, 32, yab.Width())
	testutil.AssertEqual(t, 110, yab.Height())
}
