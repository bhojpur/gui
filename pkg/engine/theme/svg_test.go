package theme

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
	"bytes"
	"encoding/xml"
	"image/color"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSVG_ReplaceFillColor(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/cancel_Paths.svg")
	if err != nil {
		t.Fatal(err)
	}
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.replaceFillColor(red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(src), string(res))
	assert.True(t, strings.Contains(string(res), "#ff0000"))
}

func TestSVG_ReplaceFillColor_Ellipse(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/ellipse.svg")
	if err != nil {
		t.Fatal(err)
	}
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.replaceFillColor(red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(src), string(res))
	assert.True(t, strings.Contains(string(res), "#ff0000"))
}
