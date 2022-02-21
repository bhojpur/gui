package engine

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
	"io"
	"math"

	"github.com/bhojpur/gui/pkg/render/canvas"
	"github.com/bhojpur/gui/pkg/render/canvas/font"
	"github.com/bhojpur/gui/pkg/render/chart"
	"github.com/bhojpur/gui/pkg/render/chart/drawing"
	"github.com/golang/freetype/truetype"
)

type BhojpurChart struct {
	c            *canvas.Canvas
	ctx          *canvas.Context
	height       float64
	writer       canvas.Writer
	dpi          float64
	font         *canvas.FontFamily
	fontSize     float64
	fontColor    drawing.Color
	textRotation float64

	fonts map[string]*canvas.FontFamily
}

// NewBhojpurChart returns a new chart renderer.
func NewBhojpurChart(writer canvas.Writer) func(int, int) (chart.Renderer, error) {
	return func(w, h int) (chart.Renderer, error) {
		cw := int(float64(w) * mmPerPx)
		ch := int(float64(h) * mmPerPx)
		c := canvas.NewCanvas("", cw, ch)
		bhojpurChart := &BhojpurChart{
			c:         &c,
			ctx:       canvas.NewContext(&c),
			height:    float64(h) * mmPerPx,
			writer:    writer,
			dpi:       chart.DefaultDPI,
			fontSize:  12.0, // uses default of github.com/golang/freetype/truetype
			fontColor: drawing.ColorTransparent,
			fonts:     map[string]*canvas.FontFamily{},
		}
		bhojpurChart.ctx.SetFillColor(canvas.Transparent)
		bhojpurChart.ctx.SetStrokeWidth(chart.DefaultStrokeWidth * mmPerPx)

		f, err := chart.GetDefaultFont()
		if err != nil {
			return nil, err
		}
		bhojpurChart.SetFont(f)
		return bhojpurChart, nil
	}
}

// ResetStyle resets any style related settings of the renderer.
func (r *BhojpurChart) ResetStyle() {
	r.ctx.ResetStyle()
	r.textRotation = 0.0
}

// GetDPI gets the DPI for the renderer.
func (r *BhojpurChart) GetDPI() float64 {
	return r.dpi
}

// SetDPI sets the DPI for the renderer.
func (r *BhojpurChart) SetDPI(dpi float64) {
	r.dpi = dpi
}

// SetClassName sets the current class name.
func (r *BhojpurChart) SetClassName(name string) {
	// TODO: SetClassName
}

// SetStrokeColor sets the current stroke color.
func (r *BhojpurChart) SetStrokeColor(col drawing.Color) {
	r.ctx.SetStrokeColor(col)
}

// SetFillColor sets the current fill color.
func (r *BhojpurChart) SetFillColor(col drawing.Color) {
	r.ctx.SetFillColor(col)
}

// SetStrokeWidth sets the stroke width.
func (r *BhojpurChart) SetStrokeWidth(width float64) {
	r.ctx.SetStrokeWidth(width * mmPerPx)
}

// SetStrokeDashArray sets the stroke dash array.
func (r *BhojpurChart) SetStrokeDashArray(dashArray []float64) {
	dashArray2 := make([]float64, len(dashArray))
	for i := 0; i < len(dashArray); i++ {
		dashArray2[i] = dashArray[i] * mmPerPx
	}
	r.ctx.SetDashes(0.0, dashArray2...)
}

// MoveTo moves the cursor to a given point.
func (r *BhojpurChart) MoveTo(x, y int) {
	r.ctx.MoveTo(float64(x)*mmPerPx, r.height-float64(y)*mmPerPx)
}

// LineTo both starts a shape and draws a line to a given point from the previous point.
func (r *BhojpurChart) LineTo(x, y int) {
	r.ctx.LineTo(float64(x)*mmPerPx, r.height-float64(y)*mmPerPx)
}

// QuadCurveTo draws a quad curve. cx and cy represent the Bézier control points.
func (r *BhojpurChart) QuadCurveTo(cx, cy, x, y int) {
	r.ctx.QuadTo(float64(cx)*mmPerPx, r.height-float64(cy)*mmPerPx, float64(x)*mmPerPx, r.height-float64(y)*mmPerPx)
}

// ArcTo draws an arc with a given center (cx,cy) a given set of radii (rx,ry), a startAngle and delta (in radians).
func (r *BhojpurChart) ArcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	startAngle = 2.0*math.Pi - startAngle
	delta = -delta

	start := canvas.EllipsePos(rx*mmPerPx, ry*mmPerPx, 0.0, float64(cx)*mmPerPx, r.height-float64(cy)*mmPerPx, startAngle)
	if r.c.Empty() {
		r.ctx.MoveTo(start.X, start.Y)
	} else {
		r.ctx.LineTo(start.X, start.Y)
	}

	startAngle *= 180.0 / math.Pi
	delta *= 180.0 / math.Pi
	r.ctx.Arc(rx*mmPerPx, ry*mmPerPx, 0.0, startAngle, startAngle+delta)
}

// Close finalizes a shape as drawn by LineTo.
func (r *BhojpurChart) Close() {
	r.ctx.Close()
	r.ctx.MoveTo(0.0, 0.0)
}

// Stroke strokes the path.
func (r *BhojpurChart) Stroke() {
	r.ctx.Stroke()
}

// Fill fills the path, but does not stroke.
func (r *BhojpurChart) Fill() {
	r.ctx.Fill()
}

// FillStroke fills and strokes a path.
func (r *BhojpurChart) FillStroke() {
	r.ctx.FillStroke()
}

// Circle draws a circle at the given coords with a given radius.
func (r *BhojpurChart) Circle(radius float64, x, y int) {
	r.ctx.DrawPath(float64(x)*mmPerPx, r.height-float64(y)*mmPerPx, canvas.Circle(radius*mmPerPx))
}

// SetFont sets a font for a text field.
func (r *BhojpurChart) SetFont(f *truetype.Font) {
	if f == nil {
		r.font = nil
		return
	}

	name := f.Name(truetype.NameIDFontFamily)
	r.font = r.fonts[name]
	if r.font == nil {
		r.font = canvas.NewFontFamily(name)
		if err := r.font.LoadFont(font.FromGoFreetype(f), 0, canvas.FontRegular); err != nil {
			panic(err)
		}
		r.fonts[name] = r.font
	}
}

// SetFontColor sets a font's color.
func (r *BhojpurChart) SetFontColor(col drawing.Color) {
	r.fontColor = col
}

// SetFontSize sets the font size for a text field.
func (r *BhojpurChart) SetFontSize(size float64) {
	r.fontSize = size
}

// Text draws a text blob.
func (r *BhojpurChart) Text(body string, x, y int) {
	if r.font == nil {
		return
	}

	face := r.font.Face(r.fontSize*ptPerMm*mmPerPx*r.dpi/72.0, r.fontColor, canvas.FontRegular, canvas.FontNormal)
	r.ctx.Push()
	r.ctx.ComposeView(canvas.Identity.Rotate(-r.textRotation * 180.0 / math.Pi))
	r.ctx.DrawText(float64(x)*mmPerPx, r.height-float64(y)*mmPerPx, canvas.NewTextLine(face, body, canvas.Left))
	r.ctx.Pop()
}

// MeasureText measures text.
func (r *BhojpurChart) MeasureText(body string) chart.Box {
	if r.font == nil {
		return chart.Box{}
	}

	face := r.font.Face(r.fontSize*ptPerMm*r.dpi/72.0, r.fontColor, canvas.FontRegular, canvas.FontNormal)
	width := face.TextWidth(body)
	return chart.Box{Right: int(math.Ceil(width)), Bottom: int(r.fontSize * r.dpi / 72.0)}
}

// SetTextRotation sets a rotation for drawing elements.
func (r *BhojpurChart) SetTextRotation(radian float64) {
	r.textRotation = radian
}

// ClearTextRotation clears rotation.
func (r *BhojpurChart) ClearTextRotation() {
	r.textRotation = 0.0
}

// Save writes the image to the given writer.
func (r *BhojpurChart) Save(w io.Writer) error {
	return r.writer(w, r.c)
}
