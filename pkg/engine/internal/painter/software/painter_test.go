package software_test

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/internal/painter/software"
	internalTest "github.com/bhojpur/gui/pkg/engine/internal/test"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func makeTestImage(w, h int) image.Image {
	return internalTest.NewCheckedImage(w, h, w, h)
}

func TestPainter_paintCircle(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewCircle(color.Black)

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(gui.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_circle.png", p.Paint(c))
}

func TestPainter_paintCircleStroke(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewCircle(color.White)
	obj.StrokeColor = color.Black
	obj.StrokeWidth = 4

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(gui.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_circle_stroke.png", p.Paint(c))
}

func TestPainter_paintGradient_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	g := canvas.NewRadialGradient(color.NRGBA{R: 200, A: 255}, color.NRGBA{B: 200, A: 255})
	g.SetMinSize(gui.NewSize(100, 100))
	scroll := container.NewScroll(g)
	scroll.Move(gui.NewPos(10, 10))
	scroll.Resize(gui.NewSize(50, 50))
	scroll.Scrolled(&gui.ScrollEvent{Scrolled: gui.NewDelta(-30, -30)})
	cont := gui.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(gui.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_gradient_clipped.png", p.Paint(c))
}

func TestPainter_paintImage(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestPainter_paintImage_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	img := canvas.NewImageFromImage(makeTestImage(5, 5))
	img.ScaleMode = canvas.ImageScalePixels
	img.SetMinSize(gui.NewSize(100, 100))
	scroll := container.NewScroll(img)
	scroll.Move(gui.NewPos(10, 10))
	scroll.Resize(gui.NewSize(50, 50))
	scroll.Scrolled(&gui.ScrollEvent{Scrolled: gui.NewDelta(-15, -15)})
	cont := gui.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(gui.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_image_clipped.png", p.Paint(c))
}

func TestPainter_paintImage_scalePixels(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_ImageScalePixels.png", target)
}

func TestPainter_paintImage_scaleSmooth(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScaleMode = canvas.ImageScaleSmooth

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_ImageScaleSmooth.png", target)
}

func TestPainter_paintImage_scaleFastest(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScaleMode = canvas.ImageScaleFastest

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_ImageScaleFastest.png", target)
}

func TestPainter_paintImage_stretchX(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage(3, 3)))
	c.Resize(gui.NewSize(100, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchx.png", target)
}

func TestPainter_paintImage_stretchY(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage(3, 3)))
	c.Resize(gui.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchy.png", target)
}

func TestPainter_paintImage_contain(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_contain.png", target)
}

func TestPainter_paintImage_containX(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	img := canvas.NewImageFromImage(makeTestImage(3, 4))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(100, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_containx.png", target)
}

func TestPainter_paintImage_containY(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	img := canvas.NewImageFromImage(makeTestImage(4, 3))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_containy.png", target)
}

func TestPainter_paintLine(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewLine(color.Black)
	obj.StrokeWidth = 6

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(gui.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_line.png", p.Paint(c))
}

func TestPainter_paintLine_thin(t *testing.T) {
	c := test.NewCanvas()
	lines := [5]*canvas.Line{}
	sws := []float32{4, 2, 1, 0.5, 0.3}
	for i, sw := range sws {
		lines[i] = canvas.NewLine(color.RGBA{255, 0, 0, 255})
		lines[i].StrokeWidth = sw
		x := float32(i * 20)
		lines[i].Position1 = gui.NewPos(x, 10)
		lines[i].Position2 = gui.NewPos(x+15, 10)
	}
	c.SetContent(container.NewWithoutLayout(lines[0], lines[1], lines[2], lines[3], lines[4]))
	c.Resize(gui.NewSize(109, 28))

	p := software.NewPainter()
	test.AssertImageMatches(t, "draw_line_thin.png", p.Paint(c))
}

func TestPainter_paintRaster(t *testing.T) {
	img := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		x = x / 5
		y = y / 5
		if x%2 == y%2 {
			return color.White
		}
		return color.Black
	})

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(gui.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_raster.png", target)
}

func TestPainter_paintRaster_scaled(t *testing.T) {
	img := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		x = x / 5
		y = y / 5
		if x%2 == y%2 {
			return color.White
		}
		return color.Black
	})

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.SetScale(5.0)
	c.Resize(gui.NewSize(5, 5))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_raster_scale.png", target)
}

func TestPainter_paintRectangle_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	red1 := canvas.NewRectangle(color.NRGBA{R: 200, A: 255})
	red1.SetMinSize(gui.NewSize(20, 20))
	red2 := canvas.NewRectangle(color.NRGBA{R: 150, A: 255})
	red2.SetMinSize(gui.NewSize(20, 20))
	red3 := canvas.NewRectangle(color.NRGBA{R: 100, A: 255})
	red3.SetMinSize(gui.NewSize(20, 20))
	reds := container.NewHBox(red1, red2, red3)
	green1 := canvas.NewRectangle(color.NRGBA{G: 200, A: 255})
	green1.SetMinSize(gui.NewSize(20, 20))
	green2 := canvas.NewRectangle(color.NRGBA{G: 150, A: 255})
	green2.SetMinSize(gui.NewSize(20, 20))
	green3 := canvas.NewRectangle(color.NRGBA{G: 100, A: 255})
	green3.SetMinSize(gui.NewSize(20, 20))
	greens := container.NewHBox(green1, green2, green3)
	blue1 := canvas.NewRectangle(color.NRGBA{B: 200, A: 255})
	blue1.SetMinSize(gui.NewSize(20, 20))
	blue2 := canvas.NewRectangle(color.NRGBA{B: 150, A: 255})
	blue2.SetMinSize(gui.NewSize(20, 20))
	blue3 := canvas.NewRectangle(color.NRGBA{B: 100, A: 255})
	blue3.SetMinSize(gui.NewSize(20, 20))
	blues := container.NewHBox(blue1, blue2, blue3)
	box := container.NewVBox(reds, greens, blues)
	scroll := container.NewScroll(box)
	scroll.Move(gui.NewPos(10, 10))
	scroll.Resize(gui.NewSize(50, 50))
	scroll.Scrolled(&gui.ScrollEvent{Scrolled: gui.NewDelta(-10, -10)})
	cont := gui.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(gui.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_rect_clipped.png", p.Paint(c))
}

func TestPainter_paintRectangle_stroke(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewRectangle(color.Black)
	obj.StrokeWidth = 5
	obj.StrokeColor = &color.RGBA{R: 0xFF, G: 0x33, B: 0x33, A: 0xFF}

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(gui.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_rectangle_stroke.png", p.Paint(c))
}

func TestPainter_paintText_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	scroll := container.NewScroll(widget.NewLabel("some text\nis here\nand here"))
	scroll.Move(gui.NewPos(10, 10))
	scroll.Resize(gui.NewSize(50, 50))
	scroll.Scrolled(&gui.ScrollEvent{Scrolled: gui.NewDelta(-10, -10)})
	cont := gui.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(gui.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_text_clipped.png", p.Paint(c))
}
