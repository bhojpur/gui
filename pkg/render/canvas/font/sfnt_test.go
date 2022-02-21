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
)

func TestSFNTDejaVuSerifTTF(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.ttf")
	test.Error(t, err)

	sfnt, err := ParseSFNT(b, 0)
	test.Error(t, err)

	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
	test.T(t, sfnt.Hhea.Ascender, int16(1901))
	test.T(t, sfnt.Hhea.Descender, int16(-483))
	test.T(t, sfnt.OS2.SCapHeight, int16(1493)) // height of H glyph
	test.T(t, sfnt.Head.XMin, int16(-1576))
	test.T(t, sfnt.Head.YMin, int16(-710))
	test.T(t, sfnt.Head.XMax, int16(4312))
	test.T(t, sfnt.Head.YMax, int16(2272))

	id := sfnt.GlyphIndex(' ')
	contour, err := sfnt.Glyf.Contour(id, 0)
	test.Error(t, err)
	test.T(t, contour.GlyphID, id)
	test.T(t, len(contour.XCoordinates), 0)
}

func TestSFNTWrite(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.ttf")
	test.Error(t, err)

	sfnt, err := ParseSFNT(b, 0)
	test.Error(t, err)

	b2 := sfnt.Write()
	sfnt2, err := ParseSFNT(b2, 0)
	test.Error(t, err)

	test.T(t, sfnt2.GlyphIndex('A'), sfnt.GlyphIndex('A'))
	test.T(t, sfnt2.GlyphIndex('B'), sfnt.GlyphIndex('B'))
	test.T(t, sfnt2.GlyphIndex('C'), sfnt.GlyphIndex('C'))

	//ioutil.WriteFile("out.otf", subset, 0644)
}

func TestSFNTSubset(t *testing.T) {
	b, err := ioutil.ReadFile("../resources/DejaVuSerif.ttf")
	test.Error(t, err)

	sfnt, err := ParseSFNT(b, 0)
	test.Error(t, err)

	subset, glyphIDs := sfnt.Subset([]uint16{0, 3, 6, 36, 37, 38, 55, 131}) // .notdef, space, #, A, B, C, T, Á
	sfntSubset, err := ParseSFNT(subset, 0)
	test.Error(t, err)

	test.T(t, len(glyphIDs), 9) // Á is a composite glyph containing two simple glyphs: 36 and 3452
	test.T(t, glyphIDs[8], uint16(3452))

	test.T(t, sfntSubset.GlyphIndex('A'), uint16(3))
	test.T(t, sfntSubset.GlyphIndex('B'), uint16(4))
	test.T(t, sfntSubset.GlyphIndex('C'), uint16(5))

	//ioutil.WriteFile("out.otf", subset, 0644)
}
