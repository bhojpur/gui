package cache

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/stretchr/testify/assert"
)

func TestSvgCacheGet(t *testing.T) {
	ResetThemeCaches()
	img := addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, len(svgs))

	newImg := GetSvg("empty.svg", 25, 25)
	assert.Equal(t, img, newImg)

	miss := GetSvg("missing.svg", 25, 25)
	assert.Nil(t, miss)
	miss = GetSvg("empty.svg", 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheGet_File(t *testing.T) {
	ResetThemeCaches()
	img := addFileToCache("testdata/stroke.svg", 25, 25)
	assert.Equal(t, 1, len(svgs))

	newImg := GetSvg("testdata/stroke.svg", 25, 25)
	assert.Equal(t, img, newImg)

	miss := GetSvg("missing.svg", 25, 25)
	assert.Nil(t, miss)
	miss = GetSvg("testdata/stroke.svg", 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheReset(t *testing.T) {
	ResetThemeCaches()
	_ = addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, len(svgs))

	ResetThemeCaches()
	assert.Equal(t, 0, len(svgs))
}

func addFileToCache(path string, w, h int) image.Image {
	img := canvas.NewImageFromFile(path)
	tex := image.NewNRGBA(image.Rect(0, 0, w, h))
	SetSvg(img.File, tex, w, h)
	return tex
}

func addToCache(name, content string, w, h int) image.Image {
	img := canvas.NewImageFromResource(gui.NewStaticResource(name, []byte(content)))
	tex := image.NewNRGBA(image.Rect(0, 0, w, h))
	SetSvg(img.Resource.Name(), tex, w, h)
	return tex
}
