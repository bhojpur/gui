package test_test

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
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/bhojpur/gui/pkg/engine/internal/painter"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/gui/pkg/engine/internal/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestAssertImageMatches(t *testing.T) {
	bounds := image.Rect(0, 0, 100, 50)
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(color.White), image.Point{}, draw.Src)

	txtImg := image.NewNRGBA(bounds)
	opts := truetype.Options{Size: 20, DPI: 96}
	f, _ := truetype.Parse(theme.TextFont().Content())
	face := truetype.NewFace(f, &opts)

	d := painter.FontDrawer{}
	d.Dst = txtImg
	d.Src = image.NewUniform(color.Black)
	d.Face = face
	d.Dot = freetype.Pt(0, 50-face.Metrics().Descent.Ceil())

	d.DrawString("Hello!", 4)
	draw.Draw(img, bounds, txtImg, image.Point{}, draw.Over)

	tt := &testing.T{}
	assert.False(t, test.AssertImageMatches(tt, "non_existing_master.png", img), "non existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test failed")
	assert.Equal(t, img, readImage(t, "testdata/failed/non_existing_master.png"), "image was written to disk")

	tt = &testing.T{}
	assert.True(t, test.AssertImageMatches(tt, "master.png", img), "existing master is equal a given image")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertImageMatches(tt, "diffing_master.png", img), "existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test did not fail")
	assert.Equal(t, img, readImage(t, "testdata/failed/diffing_master.png"), "image was written to disk")

	if !t.Failed() {
		os.RemoveAll("testdata/failed")
	}
}

func TestNewCheckedImage(t *testing.T) {
	img := test.NewCheckedImage(10, 10, 5, 2)
	expectedColorValues := [][]uint8{
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
	}
	for y, xv := range expectedColorValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{R: v, G: v, B: v, A: 0xff}, img.At(x, y), fmt.Sprintf("color value at %d,%d", x, y))
		}
	}
}

func readImage(t *testing.T, path string) image.Image {
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()
	raw, _, err := image.Decode(file)
	require.NoError(t, err)
	img := image.NewNRGBA(raw.Bounds())
	draw.Draw(img, img.Bounds(), raw, image.Pt(0, 0), draw.Src)
	return img
}
