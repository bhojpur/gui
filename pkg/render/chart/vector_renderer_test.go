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
	"fmt"
	"strings"
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/drawing"
	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestVectorRendererPath(t *testing.T) {
	// replaced new assertions helper

	vr, err := SVG(100, 100)
	testutil.AssertNil(t, err)

	typed, isTyped := vr.(*vectorRenderer)
	testutil.AssertTrue(t, isTyped)

	typed.MoveTo(0, 0)
	typed.LineTo(100, 100)
	typed.LineTo(0, 100)
	typed.Close()
	typed.FillStroke()

	buffer := bytes.NewBuffer([]byte{})
	err = typed.Save(buffer)
	testutil.AssertNil(t, err)

	raw := string(buffer.Bytes())

	testutil.AssertTrue(t, strings.HasPrefix(raw, "<svg"))
	testutil.AssertTrue(t, strings.HasSuffix(raw, "</svg>"))
}

func TestVectorRendererMeasureText(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	vr, err := SVG(100, 100)
	testutil.AssertNil(t, err)

	vr.SetDPI(DefaultDPI)
	vr.SetFont(f)
	vr.SetFontSize(12.0)

	tb := vr.MeasureText("Ljp")
	testutil.AssertEqual(t, 21, tb.Width())
	testutil.AssertEqual(t, 15, tb.Height())
}

func TestCanvasStyleSVG(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontColor:   drawing.ColorWhite,
		Font:        f,
		Padding:     DefaultBackgroundPadding,
	}

	canvas := &canvas{dpi: DefaultDPI}

	svgString := canvas.styleAsSVG(set)
	testutil.AssertNotEmpty(t, svgString)
	testutil.AssertTrue(t, strings.HasPrefix(svgString, "style=\""))
	testutil.AssertTrue(t, strings.Contains(svgString, "stroke:rgba(255,255,255,1.0)"))
	testutil.AssertTrue(t, strings.Contains(svgString, "stroke-width:5"))
	testutil.AssertTrue(t, strings.Contains(svgString, "fill:rgba(255,255,255,1.0)"))
	testutil.AssertTrue(t, strings.HasSuffix(svgString, "\""))
}

func TestCanvasClassSVG(t *testing.T) {
	set := Style{
		ClassName: "test-class",
	}

	canvas := &canvas{dpi: DefaultDPI}

	testutil.AssertEqual(t, "class=\"test-class\"", canvas.styleAsSVG(set))
}

func TestCanvasCustomInlineStylesheet(t *testing.T) {
	b := strings.Builder{}

	canvas := &canvas{
		w:   &b,
		css: ".background { fill: red }",
	}

	canvas.Start(200, 200)

	testutil.AssertContains(t, b.String(), fmt.Sprintf(`<style type="text/css"><![CDATA[%s]]></style>`, canvas.css))
}

func TestCanvasCustomInlineStylesheetWithNonce(t *testing.T) {
	b := strings.Builder{}

	canvas := &canvas{
		w:     &b,
		css:   ".background { fill: red }",
		nonce: "RAND0MSTRING",
	}

	canvas.Start(200, 200)

	testutil.AssertContains(t, b.String(), fmt.Sprintf(`<style type="text/css" nonce="%s"><![CDATA[%s]]></style>`, canvas.nonce, canvas.css))
}
