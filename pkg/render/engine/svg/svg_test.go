package svg

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
	"testing"
)

func TestSVGText(t *testing.T) {
	//dejaVuSerif := NewFontFamily("dejavu-serif")
	//dejaVuSerif.LoadFontFile("font/DejaVuSerif.ttf", FontRegular)

	//ebGaramond := NewFontFamily("eb-garamond")
	//ebGaramond.LoadFontFile("font/EBGaramond12-Regular.otf", FontRegular)

	//dejaVu8 := dejaVuSerif.Face(8.0*ptPerMm, Black, FontRegular, FontNormal)
	//dejaVu12 := dejaVuSerif.Face(12.0*ptPerMm, Red, FontItalic, FontNormal, FontUnderline)
	//dejaVu12sub := dejaVuSerif.Face(12.0*ptPerMm, Black, FontRegular, FontSubscript)
	//garamond10 := ebGaramond.Face(10.0*ptPerMm, Black, FontBold, FontNormal)

	//rt := NewRichText()
	//rt.Add(dejaVu8, "dejaVu8")
	//rt.Add(dejaVu12, " glyphspacing")
	//rt.Add(dejaVu12sub, " dejaVu12sub")
	//rt.Add(garamond10, " garamond10")
	//text := rt.ToText(dejaVu12.TextWidth("glyphspacing")+float64(len("glyphspacing")-1), 100.0, Justify, Top, 0.0, 0.0)

	//buf := &bytes.Buffer{}
	//svg := newSVGWriter(buf, 0.0, 0.0)
	//buf.Reset()
	//textLayer{text, Identity}.WriteSVG(svg)
	//s := regexp.MustCompile(`base64,.+'`).ReplaceAllString(buf.String(), "base64,'") // remove embedded font
	//test.String(t, s, `<style>`+"\n"+`@font-face{font-family:'dejavu-serif';src:url('data:font/truetype;base64,');}`+"\n"+`@font-face{font-family:'eb-garamond';src:url('data:font/opentype;base64,');}`+"\n"+`</style><text x="0" y="0" style="font: 12px dejavu-serif"><tspan x="0" y="7.421875" style="font:8px dejavu-serif">dejaVu8</tspan><tspan x="0" y="20.453125" letter-spacing="1" style="font-style:italic;fill:#f00">glyphspacing</tspan><tspan x="0" y="33.725625" style="font:700 6.996px dejavu-serif">dejaVu12sub</tspan><tspan x="0" y="38.5" style="font:700 10px eb-garamond">garamond10</tspan></text><path d="M0 22.703125H91.71875V21.803125H0z" fill="#f00"/>`)
}
