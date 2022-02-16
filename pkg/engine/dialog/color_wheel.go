package dialog

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
	"image/draw"
	"math"
	"math/cmplx"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	internalwidget "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

var _ gui.Widget = (*colorWheel)(nil)
var _ gui.Tappable = (*colorWheel)(nil)
var _ gui.Draggable = (*colorWheel)(nil)

// colorWheel displays a circular color gradient and triggers the callback when tapped.
type colorWheel struct {
	widget.BaseWidget
	generator func(w, h int) image.Image
	cache     draw.Image
	onChange  func(int, int, int, int)

	Hue                   int // Range 0-360 (degrees)
	Saturation, Lightness int // Range 0-100 (percent)
	Alpha                 int // Range 0-255
}

// newColorWheel returns a new color area that triggers the given onChange callback when tapped.
func newColorWheel(onChange func(int, int, int, int)) *colorWheel {
	a := &colorWheel{
		onChange: onChange,
	}
	a.generator = func(w, h int) image.Image {
		if a.cache == nil || a.cache.Bounds().Dx() != w || a.cache.Bounds().Dy() != h {
			rect := image.Rect(0, 0, w, h)
			a.cache = image.NewRGBA(rect)
		}
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if c := a.colorAt(x, y, w, h); c != nil {
					a.cache.Set(x, y, c)
				}
			}
		}
		return a.cache
	}
	a.ExtendBaseWidget(a)
	return a
}

// Cursor returns the cursor type of this widget.
func (a *colorWheel) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer.
func (a *colorWheel) CreateRenderer() gui.WidgetRenderer {
	raster := &canvas.Raster{
		Generator: a.generator,
	}
	x := canvas.NewLine(color.Black)
	y := canvas.NewLine(color.Black)
	return &colorWheelRenderer{
		BaseRenderer: internalwidget.NewBaseRenderer([]gui.CanvasObject{raster, x, y}),
		area:         a,
		raster:       raster,
		x:            x,
		y:            y,
	}
}

// MinSize returns the size that this widget should not shrink below.
func (a *colorWheel) MinSize() gui.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// SetHSLA updates the selected color in the wheel.
func (a *colorWheel) SetHSLA(hue, saturation, lightness, alpha int) {
	if a.Hue == hue && a.Saturation == saturation && a.Lightness == lightness && a.Alpha == alpha {
		return
	}
	a.Hue = hue
	a.Saturation = saturation
	a.Lightness = lightness
	a.Alpha = alpha
	a.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler.
func (a *colorWheel) Tapped(event *gui.PointEvent) {
	a.trigger(event.Position)
}

// Dragged is called when a pointer drag event is captured and triggers any change handler
func (a *colorWheel) Dragged(event *gui.DragEvent) {
	a.trigger(event.Position)
}

// DragEnd is called when a pointer drag ends
func (a *colorWheel) DragEnd() {
}

func (a *colorWheel) colorAt(x, y, w, h int) color.Color {
	width, height := float64(w), float64(h)
	dx := float64(x) - (width / 2.0)
	dy := float64(y) - (height / 2.0)
	radius, radians := cmplx.Polar(complex(dx, dy))
	limit := math.Min(width, height) / 2.0
	if radius > limit {
		// Out of bounds
		return theme.BackgroundColor()
	}
	degrees := radians * (180.0 / math.Pi)
	hue := wrapHue(int(degrees))
	saturation := int(radius / limit * 100.0)
	red, green, blue := hslToRgb(hue, saturation, a.Lightness)
	return &color.NRGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(a.Alpha),
	}
}

func (a *colorWheel) locationForPosition(pos gui.Position) (x, y int) {
	can := gui.CurrentApp().Driver().CanvasForObject(a)
	x, y = int(pos.X), int(pos.Y)
	if can != nil {
		x, y = can.PixelCoordinateForPosition(pos)
	}
	return
}

func (a *colorWheel) selection(width, height float32) (float32, float32) {
	w, h := float64(width), float64(height)
	radius := float64(a.Saturation) / 100.0 * math.Min(w, h) / 2.0
	degrees := float64(a.Hue)
	radians := degrees * math.Pi / 180.0
	c := cmplx.Rect(radius, radians)
	return float32(real(c) + w/2.0), float32(imag(c) + h/2.0)
}

func (a *colorWheel) trigger(pos gui.Position) {
	x, y := a.locationForPosition(pos)
	if c, f := a.cache, a.onChange; c != nil && f != nil {
		b := c.Bounds()
		width, height := float64(b.Dx()), float64(b.Dy())
		dx := float64(x) - (width / 2)
		dy := float64(y) - (height / 2)
		radius, radians := cmplx.Polar(complex(dx, dy))
		limit := math.Min(width, height) / 2.0
		if radius > limit {
			// Out of bounds
			return
		}
		degrees := radians * (180.0 / math.Pi)
		a.Hue = wrapHue(int(degrees))
		a.Saturation = int(radius / limit * 100.0)
		f(a.Hue, a.Saturation, a.Lightness, a.Alpha)
	}
	a.Refresh()
}

type colorWheelRenderer struct {
	internalwidget.BaseRenderer
	area   *colorWheel
	raster *canvas.Raster
	x, y   *canvas.Line
}

func (r *colorWheelRenderer) Layout(size gui.Size) {
	if f := r.area.selection; f != nil {
		x, y := f(size.Width, size.Height)
		r.x.Position1 = gui.NewPos(0, y)
		r.x.Position2 = gui.NewPos(size.Width, y)
		r.y.Position1 = gui.NewPos(x, 0)
		r.y.Position2 = gui.NewPos(x, size.Height)
	}
	r.raster.Move(gui.NewPos(0, 0))
	r.raster.Resize(size)
}

func (r *colorWheelRenderer) MinSize() gui.Size {
	return r.raster.MinSize().Max(gui.NewSize(128, 128))
}

func (r *colorWheelRenderer) Refresh() {
	s := r.area.Size()
	if s.IsZero() {
		r.area.Resize(r.area.MinSize())
	} else {
		r.Layout(s)
	}
	r.x.StrokeColor = theme.ForegroundColor()
	r.x.Refresh()
	r.y.StrokeColor = theme.ForegroundColor()
	r.y.Refresh()
	r.raster.Refresh()
	canvas.Refresh(r.area)
}
