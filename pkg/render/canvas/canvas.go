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
	"image/color"
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/app"
	"github.com/bhojpur/gui/pkg/engine/canvas"
)

// Canvas is where the objects are drawn into
type Canvas struct {
	Window    gui.Window
	Container *gui.Container
	Width     float64
	Height    float64
}

// NewCanvas makes a new canvas
func NewCanvas(name string, w, h int) Canvas {
	c := Canvas{
		Window:    app.New().NewWindow(name),
		Container: gui.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 255})),
		Width:     float64(w),
		Height:    float64(h),
	}
	return c
}

// MapRange -- given a value between low1 and high1, return the corresponding
// value between low2 and high2
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// Radians converts degrees to radians
func Radians(deg float64) float64 {
	return (deg * math.Pi) / 180.0
}

// Polar returns the euclidian corrdinates from polar coordinates
func Polar(x, y, r, angle float64) (float64, float64) {
	px := (r * math.Cos(Radians(angle))) + x
	py := (r * math.Sin(Radians(angle))) + y
	return px, py
}

// PolarRadians returns the euclidian corrdinates from polar coordinates
func PolarRadians(x, y, r, angle float64) (float64, float64) {
	px := (r * math.Cos(angle)) + x
	py := (r * math.Sin(angle)) + y
	return px, py
}

func pct(p float64, m float64) float64 {
	return ((p / 100.0) * m)
}

// dimen returns canvas dimensions from percentages (converting from x
// increasing left-right, y increasing top-bottom)
func dimen(xp, yp, w, h float64) (float64, float64) {
	return pct(xp, w), pct(100-yp, h)
}

// AbsStart initiates the canvas
func AbsStart(name string, w, h int) (gui.Window, *gui.Container) {
	return app.New().NewWindow(name), gui.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 255}))
}

// EndRun shows the content and runs the app
func (c *Canvas) EndRun() {
	window := c.Window
	window.Resize(gui.NewSize(float32(c.Width), float32(c.Height)))
	window.SetFixedSize(true)
	window.SetPadded(false)
	window.SetContent(c.Container)
	window.ShowAndRun()
}

// AbsEndRun shows the content and runs the app using bare windows and containers
func AbsEndRun(window gui.Window, c *gui.Container, w, h int) {
	window.Resize(gui.NewSize(float32(w), float32(h)))
	window.SetFixedSize(true)
	window.SetPadded(false)
	window.SetContent(c)
	window.ShowAndRun()
}

// iText places text
func iText(x, y int, s string, size int, color color.RGBA) *canvas.Text {
	fx, fy, fsize := float32(x), float32(y), float32(size)
	t := &canvas.Text{Text: s, Color: color, TextSize: fsize}
	adj := fsize / 5
	p := gui.Position{X: fx, Y: fy - (fsize + adj)}
	t.Move(p)
	return t
}

// iTextMid centers text
func iTextMid(x, y int, s string, size int, color color.RGBA) *canvas.Text {
	t := iText(x, y, s, size, color)
	t.Alignment = gui.TextAlignCenter
	return t
}

// iTextEnd end-aligns text
func iTextEnd(x, y int, s string, size int, color color.RGBA) *canvas.Text {
	t := iText(x, y, s, size, color)
	t.Alignment = gui.TextAlignTrailing
	return t
}

// iLine draws a line
func iLine(x1, y1, x2, y2 int, size float32, color color.RGBA) *canvas.Line {
	p1 := gui.Position{X: float32(x1), Y: float32(y1)}
	p2 := gui.Position{X: float32(x2), Y: float32(y2)}
	l := &canvas.Line{StrokeColor: color, StrokeWidth: size, Position1: p1, Position2: p2}
	return l
}

// iCircle draws a circle centered at (x,y)
func iCircle(x, y, r int, color color.RGBA) *canvas.Circle {
	fx, fy, fr := float32(x), float32(y), float32(r)
	p1 := gui.Position{X: fx - fr, Y: fy - fr}
	p2 := gui.Position{X: fx + fr, Y: fy + fr}
	c := &canvas.Circle{FillColor: color, Position1: p1, Position2: p2}
	return c
}

// iCornerRect makes a rectangle
func iCornerRect(x, y, w, h int, color color.RGBA) *canvas.Rectangle {
	r := &canvas.Rectangle{FillColor: color}
	r.Move(gui.Position{X: float32(x), Y: float32(y)})
	r.Resize(gui.Size{Width: float32(w), Height: float32(h)})
	return r
}

// IRect makes a rectangle centered at x,y
func iRect(x, y, w, h int, color color.RGBA) *canvas.Rectangle {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	r := &canvas.Rectangle{FillColor: color}
	r.Move(gui.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	r.Resize(gui.Size{Width: fw, Height: fh})
	return r
}

// iImage places the image centered at x, y
func iImage(x, y, w, h int, name string) *canvas.Image {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	i := canvas.NewImageFromFile(name)
	i.Move(gui.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	i.Resize(gui.Size{Width: fw, Height: fh})
	return i
}

// iCornerImage places the image centered at x, y
func iCornerImage(x, y, w, h int, name string) *canvas.Image {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	i := canvas.NewImageFromFile(name)
	i.Move(gui.Position{X: fx, Y: fy})
	i.Resize(gui.Size{Width: fw, Height: fh})
	return i
}

// container methods, Absoulte coordinates

// AbsText places text within a container
func AbsText(cont *gui.Container, x, y int, s string, size int, color color.RGBA) {
	fx, fy, fsize := float32(x), float32(y), float32(size)
	t := &canvas.Text{Text: s, Color: color, TextSize: fsize}
	adj := fsize / 5
	p := gui.Position{X: fx, Y: fy - (fsize + adj)}
	t.Move(p)
	cont.AddObject(t)
}

// AbsTextMid centers text within a container
func AbsTextMid(cont *gui.Container, x, y int, s string, size int, color color.RGBA) {
	t := iText(x, y, s, size, color)
	t.Alignment = gui.TextAlignCenter
	cont.AddObject(t)
}

// AbsTextEnd end-aligns text within a container
func AbsTextEnd(cont *gui.Container, x, y int, s string, size int, color color.RGBA) {
	t := iText(x, y, s, size, color)
	t.Alignment = gui.TextAlignTrailing
	cont.AddObject(t)
}

// AbsLine draws a line within a container
func AbsLine(cont *gui.Container, x1, y1, x2, y2 int, size float32, color color.RGBA) {

	//	currently there is a cap of StrokeWidth > 10 for straight lines, so make rectangles
	//	TODO: remove this special case when the bug is fixed.
	// if x1 == x2 && size > 10 { // vertical line
	// 	lineLength := y2 - y1
	// 	AbsRect(cont, x1, y1+(lineLength/2), int(size), lineLength, color)
	// 	return
	// }
	// if y1 == y2 && size > 10 { // horizontal line
	// 	lineLength := x2 - x1
	// 	AbsRect(cont, x1+(lineLength/2), y1, lineLength, int(size), color)
	// 	return
	// }
	p1 := gui.Position{X: float32(x1), Y: float32(y1)}
	p2 := gui.Position{X: float32(x2), Y: float32(y2)}
	cont.AddObject(&canvas.Line{StrokeColor: color, StrokeWidth: size, Position1: p1, Position2: p2})
}

// AbsCircle is a containerized circle within a container
func AbsCircle(cont *gui.Container, x, y, r int, color color.RGBA) {
	fx, fy, fr := float32(x), float32(y), float32(r)
	p1 := gui.Position{X: fx - fr, Y: fy - fr}
	p2 := gui.Position{X: fx + fr, Y: fy + fr}
	cont.AddObject(&canvas.Circle{FillColor: color, Position1: p1, Position2: p2})
}

// AbsCornerRect makes a rectangle within a container
func AbsCornerRect(cont *gui.Container, x, y, w, h int, color color.RGBA) {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	r := &canvas.Rectangle{FillColor: color}
	r.Move(gui.Position{X: fx, Y: fy})
	r.Resize(gui.Size{Width: fw, Height: fh})
	cont.AddObject(r)
}

// AbsRect makes a rectangle centered at x,y within a container
func AbsRect(cont *gui.Container, x, y, w, h int, color color.RGBA) {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	r := &canvas.Rectangle{FillColor: color}
	r.Move(gui.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	r.Resize(gui.Size{Width: fw, Height: fh})
	cont.AddObject(r)
}

// AbsImage places the image centered at x, y within a container
func AbsImage(cont *gui.Container, x, y, w, h int, name string) {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	i := canvas.NewImageFromFile(name)
	i.Move(gui.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	i.Resize(gui.Size{Width: fw, Height: fh})
	cont.AddObject(i)
}

// AbsCornerImage places the image centered at x, y within a container
func AbsCornerImage(cont *gui.Container, x, y, w, h int, name string) {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	i := canvas.NewImageFromFile(name)
	i.Move(gui.Position{X: fx, Y: fy})
	i.Resize(gui.Size{Width: fw, Height: fh})
	cont.AddObject(i)
}

//
// container methods, Percent coordinates
//

// TextWidth returns the width of a string
func (c *Canvas) TextWidth(s string, size float64) float64 {
	t := &canvas.Text{Text: s, TextSize: float32(pct(size, c.Width))}
	return (float64(t.MinSize().Width) / float64(c.Width)) * 100
}

// Text places text within a container, using percent coordinates
func (c *Canvas) Text(x, y float64, size float64, s string, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	AbsText(c.Container, int(x), int(y), s, int(size), color)
}

// CText places centered text using percent coordinates
func (c *Canvas) CText(x, y float64, size float64, s string, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	AbsTextMid(c.Container, int(x), int(y), s, int(size), color)
}

// EText places end-aligned text within a container, using percent coordinates
func (c *Canvas) EText(x, y float64, size float64, s string, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	AbsTextEnd(c.Container, int(x), int(y), s, int(size), color)
}

// Circle places a circle within a container, using percent coordinates
func (c *Canvas) Circle(x, y, r float64, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	r = pct(r, c.Width)
	AbsCircle(c.Container, int(x), int(y), int(r/2), color)
}

// Line places a line within a container, using percent coordinates
func (c *Canvas) Line(x1, y1, x2, y2, size float64, color color.RGBA) {
	x1, y1 = dimen(x1, y1, c.Width, c.Height)
	x2, y2 = dimen(x2, y2, c.Width, c.Height)
	lsize := pct(size, c.Width)
	AbsLine(c.Container, int(x1), int(y1), int(x2), int(y2), float32(lsize), color)

}

// Rect places a rectangle centered on (x,y) within a container, using percent coordinates
func (c *Canvas) Rect(x, y, w, h float64, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, float64(c.Width))
	h = pct(h, float64(c.Height))
	AbsCornerRect(c.Container, int(x-(w/2)), int(y-(h/2)), int(w), int(h), color)
}

// CornerRect places a rectangle with upper left corner  on (x,y) within a container,
// using percent coordinates
func (c *Canvas) CornerRect(x, y, w, h float64, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, float64(c.Width))
	h = pct(h, float64(c.Height))
	AbsCornerRect(c.Container, int(x), int(y), int(w), int(h), color)
}

// Image places an image centered at (x, y) within a container, using percent coordinates
func (c *Canvas) Image(x, y float64, w, h int, name string) {
	x, y = dimen(x, y, c.Width, c.Height)
	AbsImage(c.Container, int(x), int(y), w, h, name)
}

// ArcLine makes a stroked arc centered at (x, y), with radius r
func (c *Canvas) ArcLine(x, y, r, a1, a2, size float64, color color.RGBA) {
	step := (a2 - a1) / 100
	x1, y1 := Polar(x, y, r, a1)
	for t := a1 + step; t <= a2; t += step {
		x2, y2 := PolarRadians(x, y, r, t)
		c.Line(x1, y1, x2, y2, size, color)
		x1 = x2
		y1 = y2
	}
}
