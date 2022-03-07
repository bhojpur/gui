package gl

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
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
)

type captureImage struct {
	pix           []uint8
	width, height int
}

func (c *captureImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *captureImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height)
}

func (c *captureImage) At(x, y int) color.Color {
	start := ((c.height-y-1)*c.width + x) * 4
	return color.RGBA{R: c.pix[start], G: c.pix[start+1], B: c.pix[start+2], A: c.pix[start+3]}
}

func (p *glPainter) Capture(c gui.Canvas) image.Image {
	pos := gui.NewPos(c.Size().Width, c.Size().Height)
	width, height := c.PixelCoordinateForPosition(pos)
	pixels := make([]uint8, width*height*4)

	p.context.RunWithContext(func() {
		p.glCapture(int32(width), int32(height), &pixels)
	})

	return &captureImage{
		pix:    pixels,
		width:  width,
		height: height,
	}
}
