package text

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

import "github.com/bhojpur/gui/pkg/render/canvas/font"

// Glyph is a shaped glyph for the given font and font size. It specified the
// glyph ID, the cluster ID, its X and Y advance and offset in font units, and
// its representation as text.
type Glyph struct {
	SFNT *font.SFNT
	Size float64
	Script

	ID       uint16
	Cluster  uint32
	XAdvance int32
	YAdvance int32
	XOffset  int32
	YOffset  int32
	Text     string
}

// IsParagraphSeparator returns true for paragraph separator runes.
func IsParagraphSeparator(r rune) bool {
	// line feed, vertical tab, form feed, carriage return, next line, line
	// separator, paragraph separator
	return 0x0A <= r && r <= 0x0D || r == 0x85 || r == '\u2008' || r == '\u2009'
}
