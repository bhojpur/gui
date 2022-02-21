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
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"sort"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/app"
	"github.com/bhojpur/gui/pkg/engine/canvas"
)

//const mmPerPx = 25.4 / 96.0
//const pxPerMm = 96.0 / 25.4
const mmPerPt = 25.4 / 72.0
const ptPerMm = 72.0 / 25.4
const mmPerInch = 25.4
const inchPerMm = 1.0 / 25.4

// Resolution is used for rasterizing. Higher resolutions will result in larger images.
type Resolution float64

// DPMM (dots-per-millimeter) for the resolution of rasterization.
func DPMM(dpmm float64) Resolution {
	return Resolution(dpmm)
}

// DPI (dots-per-inch) for the resolution of rasterization.
func DPI(dpi float64) Resolution {
	return Resolution(dpi * inchPerMm)
}

// DPMM returns the resolution in dots-per-millimeter.
func (res Resolution) DPMM() float64 {
	return float64(res)
}

// DPI returns the resolution in dots-per-inch.
func (res Resolution) DPI() float64 {
	return float64(res) * mmPerInch
}

// DefaultResolution is the default resolution used for font PPEMs and is set to 96 DPI.
const DefaultResolution = Resolution(96.0 * inchPerMm)

// Size defines a size (width and height).
type Size struct {
	Width, Height float64
}

var (
	A0        = Size{841.0, 1189.0}
	A1        = Size{594.0, 841.0}
	A2        = Size{420.0, 594.0}
	A3        = Size{297.0, 420.0}
	A4        = Size{210.0, 297.0}
	A5        = Size{148.0, 210.0}
	A6        = Size{105.0, 148.0}
	A7        = Size{74.0, 105.0}
	A8        = Size{52.0, 74.0}
	B0        = Size{1000.0, 1414.0}
	B1        = Size{707.0, 1000.0}
	B2        = Size{500.0, 707.0}
	B3        = Size{353.0, 500.0}
	B4        = Size{250.0, 353.0}
	B5        = Size{176.0, 250.0}
	B6        = Size{125.0, 176.0}
	B7        = Size{88.0, 125.0}
	B8        = Size{62.0, 88.0}
	B9        = Size{44.0, 62.0}
	B10       = Size{31.0, 44.0}
	C2        = Size{648.0, 458.0}
	C3        = Size{458.0, 324.0}
	C4        = Size{324.0, 229.0}
	C5        = Size{229.0, 162.0}
	C6        = Size{162.0, 114.0}
	D0        = Size{1090.0, 771.0}
	SRA0      = Size{1280.0, 900.0}
	SRA1      = Size{900.0, 640.0}
	SRA2      = Size{640.0, 450.0}
	SRA3      = Size{450.0, 320.0}
	SRA4      = Size{320.0, 225.0}
	RA0       = Size{1220.0, 860.0}
	RA1       = Size{860.0, 610.0}
	RA2       = Size{610.0, 430.0}
	Letter    = Size{215.9, 279.4}
	Legal     = Size{215.9, 355.6}
	Ledger    = Size{279.4, 431.8}
	Tabloid   = Size{431.8, 279.4}
	Executive = Size{184.1, 266.7}
)

////////////////////////////////////////////////////////////////

// Style is the path style that defines how to draw the path. When FillColor
// is transparent it will not fill the path. If StrokeColor is transparent or
// StrokeWidth is zero, it will not stroke the path. If Dashes is an empty array,
// it will not draw dashes but instead a solid stroke line. FillRule determines
// how to fill the path when paths overlap and have certain directions (clockwise,
// counter clockwise).
type Style struct {
	FillColor    color.RGBA
	StrokeColor  color.RGBA
	StrokeWidth  float64
	StrokeCapper Capper
	StrokeJoiner Joiner
	DashOffset   float64
	Dashes       []float64
	FillRule     // TODO: test for all renderers
}

// HasFill returns true if the style has a fill
func (style Style) HasFill() bool {
	return style.FillColor.A != 0
}

// HasStroke returns true if the style has a stroke
func (style Style) HasStroke() bool {
	return style.StrokeColor.A != 0 && 0.0 < style.StrokeWidth
}

// IsDashed returns true if the style has dashes
func (style Style) IsDashed() bool {
	return 0 < len(style.Dashes)
}

// DefaultStyle is the default style for paths. It fills the path with a
// black color and has no stroke.
var DefaultStyle = Style{
	FillColor:    Black,
	StrokeColor:  Transparent,
	StrokeWidth:  1.0,
	StrokeCapper: ButtCap,
	StrokeJoiner: MiterJoin,
	DashOffset:   0.0,
	Dashes:       []float64{},
	FillRule:     NonZero,
}

// Renderer is an interface that renderers implement. It defines the size of
// the target (in mm) and functions to render paths, text objects and images.
type Renderer interface {
	Size() (float64, float64)
	RenderPath(path *Path, style Style, m Matrix)
	RenderText(text *Text, m Matrix)
	RenderImage(img image.Image, m Matrix)
}

// CoordSystem is the coordinate system, which can be either of the four
// cartesian quadrants. Most useful are the I'th and IV'th quadrants.
// CartesianI is the default quadrant with the zero-point in the bottom-left
// (the default for mathematics). The CartesianII has its zero-point in the
// bottom-right, CartesianIII in the top-right, and CartesianIV in the top-left
// (often used as default for printing devices).
type CoordSystem int

// see CoordSystem
const (
	CartesianI CoordSystem = iota
	CartesianII
	CartesianIII
	CartesianIV
)

// Context maintains the state for the current path, path style, and view transformation matrix.
type Context struct {
	Renderer

	path *Path
	Style
	styleStack     []Style
	view           Matrix
	viewStack      []Matrix
	coordView      Matrix
	coordViewStack []Matrix
}

type layer struct {
	// path, text OR img is set
	path *Path
	text *Text
	img  image.Image

	m     Matrix
	style Style // only for path
}

// Canvas is where the objects are drawn into. It stores all drawing operations
// as layers that can be re-rendered to other renderers.
type Canvas struct {
	Window    gui.Window
	Container *gui.Container
	layers    map[int][]layer
	zindex    int
	Width     float64
	Height    float64
}

// NewCanvas makes a new canvas with width and height in millimeters, that
// records all drawing operations into layers. The canvas can then be rendered
// to any other renderer.
func NewCanvas(name string, w, h int) Canvas {
	c := Canvas{
		Window:    app.New().NewWindow(name),
		Container: gui.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 255})),
		layers:    map[int][]layer{},
		Width:     float64(w),
		Height:    float64(h),
	}
	return c
}

// NewContext returns a new context which is a wrapper around a renderer. Contexts maintain the state of the current path, path style, and view transformation matrix.
func NewContext(r Renderer) *Context {
	return &Context{r, &Path{}, DefaultStyle, nil, Identity, nil, Identity, nil}
}

// Width returns the width of the canvas in millimeters.
func (c *Context) Width() float64 {
	w, _ := c.Size()
	return w
}

// Height returns the height of the canvas in millimeters.
func (c *Context) Height() float64 {
	_, h := c.Size()
	return h
}

// Push saves the current draw state so that it can be popped later on.
func (c *Context) Push() {
	c.styleStack = append(c.styleStack, c.Style)
	c.viewStack = append(c.viewStack, c.view)
	c.coordViewStack = append(c.coordViewStack, c.coordView)
}

// Pop restores the last pushed draw state and uses that as the current draw state. If there are no states on the stack, this will do nothing.
func (c *Context) Pop() {
	if len(c.styleStack) == 0 {
		return
	}
	c.Style = c.styleStack[len(c.styleStack)-1]
	c.styleStack = c.styleStack[:len(c.styleStack)-1]
	c.view = c.viewStack[len(c.viewStack)-1]
	c.viewStack = c.viewStack[:len(c.viewStack)-1]
	c.coordView = c.coordViewStack[len(c.coordViewStack)-1]
	c.coordViewStack = c.coordViewStack[:len(c.coordViewStack)-1]
}

// CoordView returns the current affine transformation matrix through which all operation coordinates will be transformed.
func (c *Context) CoordView() Matrix {
	return c.coordView
}

// SetCoordView sets the current affine transformation matrix through which all operation coordinates will be transformed. See `Matrix` for how transformations work.
func (c *Context) SetCoordView(coordView Matrix) {
	c.coordView = coordView
}

// SetCoordRect sets the current affine transformation matrix through which all operation coordinates will be transformed. It will transform coordinates from (0,0)--(width,height) to the target `rect`.
func (c *Context) SetCoordRect(rect Rect, width, height float64) {
	c.coordView = Identity.Translate(rect.X, rect.Y).Scale(rect.W/width, rect.H/height)
}

// SetCoordSystem sets the current affine transformation matrix through which all operation coordinates will be transformed as a Cartesian coordinate system.
func (c *Context) SetCoordSystem(coordSystem CoordSystem) {
	w, h := c.Size()
	switch coordSystem {
	case CartesianI:
		c.coordView = Identity
	case CartesianII:
		c.coordView = Identity.ReflectXAbout(w / 2.0)
	case CartesianIII:
		c.coordView = Identity.ReflectXAbout(w / 2.0).ReflectYAbout(h / 2.0)
	case CartesianIV:
		c.coordView = Identity.ReflectYAbout(h / 2.0)
	}
}

// View returns the current affine transformation matrix through which all operations will be transformed.
func (c *Context) View() Matrix {
	return c.view
}

// SetView sets the current affine transformation matrix through which all operations will be transformed. See `Matrix` for how transformations work.
func (c *Context) SetView(view Matrix) {
	c.view = view
}

// ResetView resets the current affine transformation matrix to the Identity matrix, ie. no transformations.
func (c *Context) ResetView() {
	c.view = Identity
}

// ComposeView post-multiplies the current affine transformation matrix by the given matrix. This means that any draw action will first be transformed by the new view matrix (parameter) and then by the current view matrix (ie. `Context.View()`). `Context.ComposeView(Identity.ReflectX())` is the same as `Context.ReflectX()`.
func (c *Context) ComposeView(view Matrix) {
	c.view = c.view.Mul(view)
}

// Translate moves the view.
func (c *Context) Translate(x, y float64) {
	c.view = c.view.Mul(Identity.Translate(x, y))
}

// ReflectX inverts the X axis of the view.
func (c *Context) ReflectX() {
	c.view = c.view.Mul(Identity.ReflectX())
}

// ReflectXAbout inverts the X axis of the view about the given X coordinate.
func (c *Context) ReflectXAbout(x float64) {
	c.view = c.view.Mul(Identity.ReflectXAbout(x))
}

// ReflectY inverts the Y axis of the view.
func (c *Context) ReflectY() {
	c.view = c.view.Mul(Identity.ReflectY())
}

// ReflectYAbout inverts the Y axis of the view about the given Y coordinate.
func (c *Context) ReflectYAbout(y float64) {
	c.view = c.view.Mul(Identity.ReflectYAbout(y))
}

// Rotate rotates the view counter clockwise with rot in degrees.
func (c *Context) Rotate(rot float64) {
	c.view = c.view.Mul(Identity.Rotate(rot))
}

// RotateAbout rotates the view counter clockwise around (x,y) with rot in degrees.
func (c *Context) RotateAbout(rot, x, y float64) {
	c.view = c.view.Mul(Identity.RotateAbout(rot, x, y))
}

// Scale scales the view.
func (c *Context) Scale(sx, sy float64) {
	c.view = c.view.Mul(Identity.Scale(sx, sy))
}

// ScaleAbout scales the view around (x,y).
func (c *Context) ScaleAbout(sx, sy, x, y float64) {
	c.view = c.view.Mul(Identity.ScaleAbout(sx, sy, x, y))
}

// Shear shear stretches the view.
func (c *Context) Shear(sx, sy float64) {
	c.view = c.view.Mul(Identity.Shear(sx, sy))
}

// ShearAbout shear stretches the view around (x,y).
func (c *Context) ShearAbout(sx, sy, x, y float64) {
	c.view = c.view.Mul(Identity.ShearAbout(sx, sy, x, y))
}

// SetFillColor sets the color to be used for filling operations.
func (c *Context) SetFillColor(col color.Color) {
	r, g, b, a := col.RGBA()
	// RGBA returns an alpha-premultiplied color so that c <= a. We silently correct the color by clipping r,g,b to a
	if a < r {
		r = a
	}
	if a < g {
		g = a
	}
	if a < b {
		b = a
	}
	c.Style.FillColor = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// SetStrokeColor sets the color to be used for stroking operations.
func (c *Context) SetStrokeColor(col color.Color) {
	r, g, b, a := col.RGBA()
	// RGBA returns an alpha-premultiplied color so that c <= a. We silently correct the color by clipping r,g,b to a
	if a < r {
		r = a
	}
	if a < g {
		g = a
	}
	if a < b {
		b = a
	}
	c.Style.StrokeColor = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// SetStrokeWidth sets the width in millimeters for stroking operations.
func (c *Context) SetStrokeWidth(width float64) {
	c.Style.StrokeWidth = width
}

// SetStrokeCapper sets the line cap function to be used for stroke end points.
func (c *Context) SetStrokeCapper(capper Capper) {
	c.Style.StrokeCapper = capper
}

// SetStrokeJoiner sets the line join function to be used for stroke mid points.
func (c *Context) SetStrokeJoiner(joiner Joiner) {
	c.Style.StrokeJoiner = joiner
}

// SetDashes sets the dash pattern to be used for stroking operations. The dash offset denotes the offset into the dash array in millimeters from where to start. Negative values are allowed.
func (c *Context) SetDashes(offset float64, dashes ...float64) {
	c.Style.DashOffset = offset
	c.Style.Dashes = dashes
}

// SetFillRule sets the fill rule to be used for filling paths.
func (c *Context) SetFillRule(rule FillRule) {
	c.Style.FillRule = rule
}

// ResetStyle resets the draw state to its default (colors, stroke widths, dashes, ...).
func (c *Context) ResetStyle() {
	c.Style = DefaultStyle
}

// SetZIndex sets the z-index. This will call the renderer's `SetZIndex` function only if it exists (in this case only for `Canvas`).
func (c *Context) SetZIndex(zindex int) {
	if zindexer, ok := c.Renderer.(interface{ SetZIndex(int) }); ok {
		zindexer.SetZIndex(zindex)
	}
}

// Pos returns the current position of the path, which is the end point of the last command.
func (c *Context) Pos() (float64, float64) {
	return c.path.Pos().X, c.path.Pos().Y
}

// MoveTo moves the path to (x,y) without connecting with the previous path. It starts a new independent subpath. Multiple subpaths can be useful when negating parts of a previous path by overlapping it with a path in the opposite direction. The behaviour of overlapping paths depends on the FillRule.
func (c *Context) MoveTo(x, y float64) {
	c.path.MoveTo(x, y)
}

// LineTo adds a linear path to (x,y).
func (c *Context) LineTo(x, y float64) {
	c.path.LineTo(x, y)
}

// QuadTo adds a quadratic Bézier path with control point (cpx,cpy) and end point (x,y).
func (c *Context) QuadTo(cpx, cpy, x, y float64) {
	c.path.QuadTo(cpx, cpy, x, y)
}

// CubeTo adds a cubic Bézier path with control points (cpx1,cpy1) and (cpx2,cpy2) and end point (x,y).
func (c *Context) CubeTo(cpx1, cpy1, cpx2, cpy2, x, y float64) {
	c.path.CubeTo(cpx1, cpy1, cpx2, cpy2, x, y)
}

// ArcTo adds an arc with radii rx and ry, with rot the counter clockwise rotation with respect to the coordinate system in degrees, large and sweep booleans (see https://developer.mozilla.org/en-US/docs/Web/SVG/Tutorial/Paths#Arcs), and (x,y) the end position of the pen. The start position of the pen was given by a previous command's end point.
func (c *Context) ArcTo(rx, ry, rot float64, large, sweep bool, x, y float64) {
	c.path.ArcTo(rx, ry, rot, large, sweep, x, y)
}

// Arc adds an elliptical arc with radii rx and ry, with rot the counter clockwise rotation in degrees, and theta0 and theta1 the angles in degrees of the ellipse (before rot is applied) between which the arc will run. If theta0 < theta1, the arc will run in a CCW direction. If the difference between theta0 and theta1 is bigger than 360 degrees, one full circle will be drawn and the remaining part of diff % 360, e.g. a difference of 810 degrees will draw one full circle and an arc over 90 degrees.
func (c *Context) Arc(rx, ry, rot, theta0, theta1 float64) {
	c.path.Arc(rx, ry, rot, theta0, theta1)
}

// Close closes the current path.
func (c *Context) Close() {
	c.path.Close()
}

// Fill fills the current path and resets the path.
func (c *Context) Fill() {
	style := c.Style
	style.StrokeColor = Transparent
	c.RenderPath(c.path, style, c.view)
	c.path = &Path{}
}

// Stroke strokes the current path and resets the path.
func (c *Context) Stroke() {
	style := c.Style
	style.FillColor = Transparent
	c.RenderPath(c.path, style, c.view)
	c.path = &Path{}
}

// FillStroke fills and then strokes the current path and resets the path.
func (c *Context) FillStroke() {
	c.RenderPath(c.path, c.Style, c.view)
	c.path = &Path{}
}

// DrawPath draws a path at position (x,y) using the current draw state.
func (c *Context) DrawPath(x, y float64, paths ...*Path) {
	if !c.Style.HasFill() && !c.Style.HasStroke() {
		return
	}

	coord := c.coordView.Dot(Point{x, y})
	m := c.view.Translate(coord.X, coord.Y)
	for _, path := range paths {
		var dashes []float64
		path, dashes = path.checkDash(c.Style.DashOffset, c.Style.Dashes)
		if path.Empty() {
			continue
		}
		style := c.Style
		style.Dashes = dashes
		c.RenderPath(path, style, m)
	}
}

// DrawText draws text at position (x,y) using the current draw state.
func (c *Context) DrawText(x, y float64, texts ...*Text) {
	coord := c.coordView.Dot(Point{x, y})
	m := c.view.Translate(coord.X, coord.Y)
	for _, text := range texts {
		if text.Empty() {
			continue
		}
		c.RenderText(text, m)
	}
}

// DrawImage draws an image at position (x,y) using the current draw state and the given resolution in pixels-per-millimeter. A higher resolution will draw a smaller image (ie. more image pixels per millimeter of document).
func (c *Context) DrawImage(x, y float64, img image.Image, resolution Resolution) {
	if img.Bounds().Size().Eq(image.Point{}) {
		return
	}

	coord := c.coordView.Dot(Point{x, y})
	m := c.view.Translate(coord.X, coord.Y).Scale(1.0/resolution.DPMM(), 1.0/resolution.DPMM())
	c.RenderImage(img, m)
}

// NewCanvasFromSize returns a new canvas of given size in millimeters, that records all
// drawing operations into layers. The canvas can then be rendered to any other renderer.
func NewCanvasFromSize(name string, size Size) Canvas {
	w := int(size.Width)
	h := int(size.Height)
	return NewCanvas(name, w, h)
}

// Size returns the size of the canvas in millimeters.
func (c *Canvas) Size() (float64, float64) {
	return c.Width, c.Height
}

// RenderPath renders a path to the canvas using a style and a transformation matrix.
func (c *Canvas) RenderPath(path *Path, style Style, m Matrix) {
	path = path.Copy()
	c.layers[c.zindex] = append(c.layers[c.zindex], layer{path: path, m: m, style: style})
}

// RenderText renders a text object to the canvas using a transformation matrix.
func (c *Canvas) RenderText(text *Text, m Matrix) {
	c.layers[c.zindex] = append(c.layers[c.zindex], layer{text: text, m: m})
}

// RenderImage renders an image to the canvas using a transformation matrix.
func (c *Canvas) RenderImage(img image.Image, m Matrix) {
	c.layers[c.zindex] = append(c.layers[c.zindex], layer{img: img, m: m})
}

// Empty return true if the canvas is empty.
func (c *Canvas) Empty() bool {
	return len(c.layers) == 0
}

// Reset empties the canvas.
func (c *Canvas) Reset() {
	c.layers = map[int][]layer{}
}

// SetZIndex sets the z-index.
func (c *Canvas) SetZIndex(zindex int) {
	c.zindex = zindex
}

// Fit shrinks the canvas' size so all elements fit with a given margin in millimeters.
func (c *Canvas) Fit(margin float64) {
	rect := Rect{}
	// TODO: slow when we have many paths (see Graph example)
	for _, layers := range c.layers {
		for i, l := range layers {
			bounds := Rect{}
			if l.path != nil {
				bounds = l.path.Bounds()
				if l.style.StrokeColor.A != 0 && 0.0 < l.style.StrokeWidth {
					bounds.X -= l.style.StrokeWidth / 2.0
					bounds.Y -= l.style.StrokeWidth / 2.0
					bounds.W += l.style.StrokeWidth
					bounds.H += l.style.StrokeWidth
				}
			} else if l.text != nil {
				bounds = l.text.Bounds()
			} else if l.img != nil {
				size := l.img.Bounds().Size()
				bounds = Rect{0.0, 0.0, float64(size.X), float64(size.Y)}
			}
			bounds = bounds.Transform(l.m)
			if i == 0 {
				rect = bounds
			} else {
				rect = rect.Add(bounds)
			}
		}
	}
	for _, layers := range c.layers {
		for i := range layers {
			layers[i].m = Identity.Translate(-rect.X+margin, -rect.Y+margin).Mul(layers[i].m)
		}
	}
	c.Width = rect.W + 2*margin
	c.Height = rect.H + 2*margin
}

// Render renders the accumulated canvas drawing operations to another renderer.
func (c *Canvas) Render(r Renderer) {
	view := Identity
	if viewer, ok := r.(interface{ View() Matrix }); ok {
		view = viewer.View()
	}

	zindices := []int{}
	for zindex := range c.layers {
		zindices = append(zindices, zindex)
	}
	sort.Ints(zindices)

	for _, zindex := range zindices {
		for _, l := range c.layers[zindex] {
			m := view.Mul(l.m)
			if l.path != nil {
				r.RenderPath(l.path, l.style, m)
			} else if l.text != nil {
				r.RenderText(l.text, m)
			} else if l.img != nil {
				r.RenderImage(l.img, m)
			}
		}
	}
}

// Writer can write a canvas to a writer.
type Writer func(w io.Writer, c *Canvas) error

// WriteFile writes the canvas to a file named by filename using the given writer.
func (c *Canvas) WriteFile(filename string, w Writer) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if err = w(f, c); err != nil {
		f.Close()
		return err
	}
	return f.Close()
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
