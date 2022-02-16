package canvas_test

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
	"image"
	"testing"

	"github.com/bhojpur/gui/pkg/engine/canvas"

	"github.com/stretchr/testify/assert"
)

func TestRasterFromImage(t *testing.T) {
	source := image.Rect(2, 2, 4, 4)
	dest := canvas.NewRasterFromImage(source)
	img := dest.Generator(6, 6)

	// image.Rect is a 16 bit color model
	_, _, _, a := img.At(0, 0).RGBA()
	assert.Equal(t, uint32(0x0000), a)
	_, _, _, a = img.At(2, 2).RGBA()
	assert.Equal(t, uint32(0xffff), a)
	_, _, _, a = img.At(4, 4).RGBA()
	assert.Equal(t, uint32(0x0000), a)
}
