package widget

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Widget = (*Shadow)(nil)

// Shadow is a widget that renders a shadow.
type Shadow struct {
	Base
	level ElevationLevel
	typ   ShadowType
}

// ElevationLevel is the level of elevation of the shadow casting object.
type ElevationLevel int

// ElevationLevel constants
const (
	BaseLevel             ElevationLevel = 0
	CardLevel             ElevationLevel = 1
	ButtonLevel           ElevationLevel = 2
	MenuLevel             ElevationLevel = 4
	PopUpLevel            ElevationLevel = 8
	SubmergedContentLevel ElevationLevel = 8
	DialogLevel           ElevationLevel = 24
)

// ShadowType specifies the type of the shadow.
type ShadowType int

// ShadowType constants
const (
	ShadowAround ShadowType = iota
	ShadowLeft
	ShadowRight
	ShadowBottom
	ShadowTop
)

// NewShadow create a new Shadow.
func NewShadow(typ ShadowType, level ElevationLevel) *Shadow {
	s := &Shadow{typ: typ, level: level}
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer returns a new renderer for the shadow.
//
// Implements: gui.Widget
func (s *Shadow) CreateRenderer() gui.WidgetRenderer {
	r := &shadowRenderer{s: s}
	r.createShadows()
	return r
}

type shadowRenderer struct {
	BaseRenderer
	b, l, r, t     *canvas.LinearGradient
	bl, br, tl, tr *canvas.RadialGradient
	minSize        gui.Size
	s              *Shadow
}

func (r *shadowRenderer) Layout(size gui.Size) {
	depth := float32(r.s.level)
	if r.tl != nil {
		r.tl.Resize(gui.NewSize(depth, depth))
		r.tl.Move(gui.NewPos(-depth, -depth))
	}
	if r.t != nil {
		r.t.Resize(gui.NewSize(size.Width, depth))
		r.t.Move(gui.NewPos(0, -depth))
	}
	if r.tr != nil {
		r.tr.Resize(gui.NewSize(depth, depth))
		r.tr.Move(gui.NewPos(size.Width, -depth))
	}
	if r.r != nil {
		r.r.Resize(gui.NewSize(depth, size.Height))
		r.r.Move(gui.NewPos(size.Width, 0))
	}
	if r.br != nil {
		r.br.Resize(gui.NewSize(depth, depth))
		r.br.Move(gui.NewPos(size.Width, size.Height))
	}
	if r.b != nil {
		r.b.Resize(gui.NewSize(size.Width, depth))
		r.b.Move(gui.NewPos(0, size.Height))
	}
	if r.bl != nil {
		r.bl.Resize(gui.NewSize(depth, depth))
		r.bl.Move(gui.NewPos(-depth, size.Height))
	}
	if r.l != nil {
		r.l.Resize(gui.NewSize(depth, size.Height))
		r.l.Move(gui.NewPos(-depth, 0))
	}
}

func (r *shadowRenderer) MinSize() gui.Size {
	return r.minSize
}

func (r *shadowRenderer) Refresh() {
	r.refreshShadows()
	r.Layout(r.s.Size())
	canvas.Refresh(r.s)
}

func (r *shadowRenderer) createShadows() {
	switch r.s.typ {
	case ShadowLeft:
		r.l = canvas.NewHorizontalGradient(color.Transparent, theme.ShadowColor())
		r.SetObjects([]gui.CanvasObject{r.l})
	case ShadowRight:
		r.r = canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
		r.SetObjects([]gui.CanvasObject{r.r})
	case ShadowBottom:
		r.b = canvas.NewVerticalGradient(theme.ShadowColor(), color.Transparent)
		r.SetObjects([]gui.CanvasObject{r.b})
	case ShadowTop:
		r.t = canvas.NewVerticalGradient(color.Transparent, theme.ShadowColor())
		r.SetObjects([]gui.CanvasObject{r.t})
	case ShadowAround:
		r.tl = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.tl.CenterOffsetX = 0.5
		r.tl.CenterOffsetY = 0.5
		r.t = canvas.NewVerticalGradient(color.Transparent, theme.ShadowColor())
		r.tr = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.tr.CenterOffsetX = -0.5
		r.tr.CenterOffsetY = 0.5
		r.r = canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
		r.br = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.br.CenterOffsetX = -0.5
		r.br.CenterOffsetY = -0.5
		r.b = canvas.NewVerticalGradient(theme.ShadowColor(), color.Transparent)
		r.bl = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.bl.CenterOffsetX = 0.5
		r.bl.CenterOffsetY = -0.5
		r.l = canvas.NewHorizontalGradient(color.Transparent, theme.ShadowColor())
		r.SetObjects([]gui.CanvasObject{r.tl, r.t, r.tr, r.r, r.br, r.b, r.bl, r.l})
	}
}

func (r *shadowRenderer) refreshShadows() {
	updateShadowEnd(r.l)
	updateShadowStart(r.r)
	updateShadowStart(r.b)
	updateShadowEnd(r.t)

	updateShadowRadial(r.tl)
	updateShadowRadial(r.tr)
	updateShadowRadial(r.bl)
	updateShadowRadial(r.br)
}

func updateShadowEnd(g *canvas.LinearGradient) {
	if g == nil {
		return
	}

	g.EndColor = theme.ShadowColor()
	g.Refresh()
}

func updateShadowRadial(g *canvas.RadialGradient) {
	if g == nil {
		return
	}

	g.StartColor = theme.ShadowColor()
	g.Refresh()
}

func updateShadowStart(g *canvas.LinearGradient) {
	if g == nil {
		return
	}

	g.StartColor = theme.ShadowColor()
	g.Refresh()
}
