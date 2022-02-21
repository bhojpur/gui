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
	"bytes"
	"errors"
	"image"
	"image/png"
)

// RGBACollector is a render target for a chart.
type RGBACollector interface {
	SetRGBA(i *image.RGBA)
}

// ImageWriter is a special type of io.Writer that produces a final image.
type ImageWriter struct {
	rgba     *image.RGBA
	contents *bytes.Buffer
}

func (ir *ImageWriter) Write(buffer []byte) (int, error) {
	if ir.contents == nil {
		ir.contents = bytes.NewBuffer([]byte{})
	}
	return ir.contents.Write(buffer)
}

// SetRGBA sets a raw version of the image.
func (ir *ImageWriter) SetRGBA(i *image.RGBA) {
	ir.rgba = i
}

// Image returns an *image.Image for the result.
func (ir *ImageWriter) Image() (image.Image, error) {
	if ir.rgba != nil {
		return ir.rgba, nil
	}
	if ir.contents != nil && ir.contents.Len() > 0 {
		return png.Decode(ir.contents)
	}
	return nil, errors.New("no valid sources for image data, cannot continue")
}
