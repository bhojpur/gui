package canvas

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
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// ImageEncoding defines whether the embedded image shall be embedded as
// lossless (typically PNG) or lossy (typically JPG).
type ImageEncoding int

// see ImageEncoding
const (
	Lossless ImageEncoding = iota
	Lossy
)

// Image is a raster image. Keeping the original bytes allows the renderer to
// optimize rendering in some cases.
type Image struct {
	image.Image
	Mimetype string
	Bytes    []byte
}

// NewJPEGImage parses a JPEG image.
func NewJPEGImage(r io.Reader) (Image, error) {
	return newImage("image/jpeg", jpeg.Decode, r)
}

// NewPNGImage parses a PNG image
func NewPNGImage(r io.Reader) (Image, error) {
	return newImage("image/png", png.Decode, r)
}

func newImage(mimetype string, decode func(io.Reader) (image.Image, error), r io.Reader) (Image, error) {
	// TODO: use lazy decoding
	var buffer bytes.Buffer
	r = io.TeeReader(r, &buffer)
	img, err := decode(r)
	return Image{
		Image:    img,
		Bytes:    buffer.Bytes(),
		Mimetype: mimetype,
	}, err
}
