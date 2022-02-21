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
	"time"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestTimeValueFormatterWithFormat(t *testing.T) {
	// replaced new assertions helper

	d := time.Now()
	di := TimeToFloat64(d)
	df := float64(di)

	s := formatTime(d, DefaultDateFormat)
	si := formatTime(di, DefaultDateFormat)
	sf := formatTime(df, DefaultDateFormat)
	testutil.AssertEqual(t, s, si)
	testutil.AssertEqual(t, s, sf)

	sd := TimeValueFormatter(d)
	sdi := TimeValueFormatter(di)
	sdf := TimeValueFormatter(df)
	testutil.AssertEqual(t, s, sd)
	testutil.AssertEqual(t, s, sdi)
	testutil.AssertEqual(t, s, sdf)
}

func TestFloatValueFormatter(t *testing.T) {
	// replaced new assertions helper
	testutil.AssertEqual(t, "1234.00", FloatValueFormatter(1234.00))
}

func TestFloatValueFormatterWithFloat32Input(t *testing.T) {
	// replaced new assertions helper
	testutil.AssertEqual(t, "1234.00", FloatValueFormatter(float32(1234.00)))
}

func TestFloatValueFormatterWithIntegerInput(t *testing.T) {
	// replaced new assertions helper
	testutil.AssertEqual(t, "1234.00", FloatValueFormatter(1234))
}

func TestFloatValueFormatterWithInt64Input(t *testing.T) {
	// replaced new assertions helper
	testutil.AssertEqual(t, "1234.00", FloatValueFormatter(int64(1234)))
}

func TestFloatValueFormatterWithFormat(t *testing.T) {
	// replaced new assertions helper

	v := 123.456
	sv := FloatValueFormatterWithFormat(v, "%.3f")
	testutil.AssertEqual(t, "123.456", sv)
	testutil.AssertEqual(t, "123.000", FloatValueFormatterWithFormat(123, "%.3f"))
}
