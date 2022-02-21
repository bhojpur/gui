package font

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
	"io/ioutil"
	"testing"

	"github.com/bhojpur/gui/pkg/render/canvas/test"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/sfnt"
)

func TestParseTTF(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.ttf")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestParseOTF(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/EBGaramond12-Regular.otf")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(1000))
}

//func TestParseOTF_CFF2(t *testing.T) {
//	b, err := ioutil.ReadFile("../resources/AdobeVFPrototype.otf") // TODO: CFF2
//	test.Error(t, err)
//
//	sfnt, err := ParseFont(b, 0)
//	test.Error(t, err)
//	test.T(t, sfnt.Head.UnitsPerEm, uint16(1000))
//}

func TestParseWOFF(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.woff")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestParseWOFF2(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.woff2")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestParseEOT(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.eot")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestFromGoFreetype(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.ttf")
	test.Error(t, err)

	font, err := truetype.Parse(b)
	test.Error(t, err)

	buf := FromGoFreetype(font)
	sfnt, err := ParseSFNT(buf, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestFromGoSFNT(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.ttf")
	test.Error(t, err)

	font, err := sfnt.Parse(b)
	test.Error(t, err)

	buf := FromGoSFNT(font)
	sfnt, err := ParseSFNT(buf, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}
