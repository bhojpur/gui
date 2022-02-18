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
	"image/color"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	intRepo "github.com/bhojpur/gui/pkg/engine/internal/repository"
	"github.com/bhojpur/gui/pkg/engine/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestFromJSON(t *testing.T) {
	repository.Register("file", intRepo.NewFileRepository()) // file uri resolving (avoid test import loop)
	th, err := FromJSON(`{
"Colors": {"background": "#c0c0c0ff"},
"Colors-light": {"foreground": "#ffffffff"},
"Sizes": {"iconInline": 5.0},
"Fonts": {"monospace": "file://./testdata/NotoMono-Regular.ttf"},
"Icons": {"cancel": "file://./testdata/cancel_Paths.svg"}
}`)

	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameBackground, VariantDark))
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, th.Color(ColorNameForeground, VariantLight))
	assert.Equal(t, float32(5), th.Size(SizeNameInlineIcon))
	assert.Equal(t, "NotoMono-Regular.ttf", th.Font(gui.TextStyle{Monospace: true}).Name())
	assert.Equal(t, "cancel_Paths.svg", th.Icon(IconNameCancel).Name())
}

func TestFromTOML_Resource(t *testing.T) {
	r, err := gui.LoadResourceFromPath("./testdata/theme.json")
	assert.Nil(t, err)
	th, err := FromJSON(string(r.Content()))

	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}, th.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, th.Color(ColorNameForeground, VariantDark))
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameForeground, VariantLight))
	assert.Equal(t, float32(10), th.Size(SizeNameInlineIcon))
}

func TestHexColor(t *testing.T) {
	c, err := hexColor("#abc").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}, c)
	c, err = hexColor("abc").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}, c)
	c, err = hexColor("#abcd").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xdd}, c)

	c, err = hexColor("#a1b2c3").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xff}, c)
	c, err = hexColor("a1b2c3").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xff}, c)
	c, err = hexColor("#a1b2c3f4").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xf4}, c)
	c, err = hexColor("a1b2c3f4").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xf4}, c)
}
