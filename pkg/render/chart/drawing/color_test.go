package drawing

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

	"image/color"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestColorFromHex(t *testing.T) {
	// replaced new assertions helper

	white := ColorFromHex("FFFFFF")
	testutil.AssertEqual(t, ColorWhite, white)

	shortWhite := ColorFromHex("FFF")
	testutil.AssertEqual(t, ColorWhite, shortWhite)

	black := ColorFromHex("000000")
	testutil.AssertEqual(t, ColorBlack, black)

	shortBlack := ColorFromHex("000")
	testutil.AssertEqual(t, ColorBlack, shortBlack)

	red := ColorFromHex("FF0000")
	testutil.AssertEqual(t, ColorRed, red)

	shortRed := ColorFromHex("F00")
	testutil.AssertEqual(t, ColorRed, shortRed)

	green := ColorFromHex("00FF00")
	testutil.AssertEqual(t, ColorGreen, green)

	shortGreen := ColorFromHex("0F0")
	testutil.AssertEqual(t, ColorGreen, shortGreen)

	blue := ColorFromHex("0000FF")
	testutil.AssertEqual(t, ColorBlue, blue)

	shortBlue := ColorFromHex("00F")
	testutil.AssertEqual(t, ColorBlue, shortBlue)
}

func TestColorFromAlphaMixedRGBA(t *testing.T) {
	// replaced new assertions helper

	black := ColorFromAlphaMixedRGBA(color.Black.RGBA())
	testutil.AssertTrue(t, black.Equals(ColorBlack), black.String())

	white := ColorFromAlphaMixedRGBA(color.White.RGBA())
	testutil.AssertTrue(t, white.Equals(ColorWhite), white.String())
}
