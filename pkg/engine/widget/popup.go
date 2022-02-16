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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// PopUp is a widget that can float above the user interface.
// It wraps any standard elements with padding and a shadow.
// If it is modal then the shadow will cover the entire canvas it hovers over and block interactions.
type PopUp struct {
	BaseWidget

	Content gui.CanvasObject
	Canvas  gui.Canvas

	innerPos     gui.Position
	innerSize    gui.Size
	modal        bool
	overlayShown bool
}

// Hide this widget, if it was previously visible
func (p *PopUp) Hide() {
	if p.overlayShown {
		p.Canvas.Overlays().Remove(p)
		p.overlayShown = false
	}
	p.BaseWidget.Hide()
}

// Move the widget to a new position. A PopUp position is absolute to the top, left of its canvas.
// For PopUp this actually moves the content so checking Position() will not return the same value as is set here.
func (p *PopUp) Move(pos gui.Position) {
	if p.modal {
		return
	}
	p.innerPos = pos
	p.Refresh()
}

// Resize changes the size of the PopUp's content.
// PopUps always have the size of their canvas, but this call updates the
// size of the content portion.
//
// Implements: gui.Widget
func (p *PopUp) Resize(size gui.Size) {
	p.innerSize = size
	// The canvas size might not have changed and therefore the Resize won't trigger a layout.
	// Until we have a widget.Relayout() or similar, the renderer's refresh will do the re-layout.
	p.Refresh()
}

// Show this pop-up as overlay if not already shown.
func (p *PopUp) Show() {
	if !p.overlayShown {
		p.Canvas.Overlays().Add(p)
		p.overlayShown = true
	}
	p.Refresh()
	p.BaseWidget.Show()
}

// ShowAtPosition shows this pop-up at the given position.
func (p *PopUp) ShowAtPosition(pos gui.Position) {
	p.Move(pos)
	p.Show()
}

// Tapped is called when the user taps the popUp background - if not modal then dismiss this widget
func (p *PopUp) Tapped(_ *gui.PointEvent) {
	if !p.modal {
		p.Hide()
	}
}

// TappedSecondary is called when the user right/alt taps the background - if not modal then dismiss this widget
func (p *PopUp) TappedSecondary(_ *gui.PointEvent) {
	if !p.modal {
		p.Hide()
	}
}

// MinSize returns the size that this widget should not shrink below
func (p *PopUp) MinSize() gui.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (p *PopUp) CreateRenderer() gui.WidgetRenderer {
	p.ExtendBaseWidget(p)
	background := canvas.NewRectangle(theme.BackgroundColor())
	if p.modal {
		underlay := canvas.NewRectangle(theme.ShadowColor())
		objects := []gui.CanvasObject{underlay, background, p.Content}
		return &modalPopUpRenderer{
			widget.NewShadowingRenderer(objects, widget.DialogLevel),
			popUpBaseRenderer{popUp: p, background: background},
			underlay,
		}
	}
	objects := []gui.CanvasObject{background, p.Content}
	return &popUpRenderer{
		widget.NewShadowingRenderer(objects, widget.PopUpLevel),
		popUpBaseRenderer{popUp: p, background: background},
	}
}

// ShowPopUpAtPosition creates a new popUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
func ShowPopUpAtPosition(content gui.CanvasObject, canvas gui.Canvas, pos gui.Position) {
	newPopUp(content, canvas).ShowAtPosition(pos)
}

func newPopUp(content gui.CanvasObject, canvas gui.Canvas) *PopUp {
	ret := &PopUp{Content: content, Canvas: canvas, modal: false}
	ret.ExtendBaseWidget(ret)
	return ret
}

// NewPopUp creates a new popUp for the specified content and displays it on the passed canvas.
func NewPopUp(content gui.CanvasObject, canvas gui.Canvas) *PopUp {
	return newPopUp(content, canvas)
}

// ShowPopUp creates a new popUp for the specified content and displays it on the passed canvas.
func ShowPopUp(content gui.CanvasObject, canvas gui.Canvas) {
	newPopUp(content, canvas).Show()
}

func newModalPopUp(content gui.CanvasObject, canvas gui.Canvas) *PopUp {
	p := &PopUp{Content: content, Canvas: canvas, modal: true}
	p.ExtendBaseWidget(p)
	return p
}

// NewModalPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// A modal PopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func NewModalPopUp(content gui.CanvasObject, canvas gui.Canvas) *PopUp {
	return newModalPopUp(content, canvas)
}

// ShowModalPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// A modal PopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func ShowModalPopUp(content gui.CanvasObject, canvas gui.Canvas) {
	p := newModalPopUp(content, canvas)
	p.Show()
}

type popUpBaseRenderer struct {
	popUp      *PopUp
	background *canvas.Rectangle
}

func (r *popUpBaseRenderer) padding() gui.Size {
	return gui.NewSize(theme.Padding()*2, theme.Padding()*2)
}

func (r *popUpBaseRenderer) offset() gui.Position {
	return gui.NewPos(theme.Padding(), theme.Padding())
}

type popUpRenderer struct {
	*widget.ShadowingRenderer
	popUpBaseRenderer
}

func (r *popUpRenderer) Layout(_ gui.Size) {
	innerSize := r.popUp.innerSize.Max(r.popUp.MinSize())
	r.popUp.Content.Resize(innerSize.Subtract(r.padding()))

	innerPos := r.popUp.innerPos
	if innerPos.X+innerSize.Width > r.popUp.Canvas.Size().Width {
		innerPos.X = r.popUp.Canvas.Size().Width - innerSize.Width
		if innerPos.X < 0 {
			innerPos.X = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if innerPos.Y+innerSize.Height > r.popUp.Canvas.Size().Height {
		innerPos.Y = r.popUp.Canvas.Size().Height - innerSize.Height
		if innerPos.Y < 0 {
			innerPos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	r.popUp.Content.Move(innerPos.Add(r.offset()))

	r.background.Resize(innerSize)
	r.background.Move(innerPos)
	r.LayoutShadow(innerSize, innerPos)
}

func (r *popUpRenderer) MinSize() gui.Size {
	return r.popUp.Content.MinSize().Add(r.padding())
}

func (r *popUpRenderer) Refresh() {
	r.background.FillColor = theme.BackgroundColor()
	expectedContentSize := r.popUp.innerSize.Max(r.popUp.MinSize()).Subtract(r.padding())
	shouldRelayout := r.popUp.Content.Size() != expectedContentSize

	if r.background.Size() != r.popUp.innerSize || r.background.Position() != r.popUp.innerPos || shouldRelayout {
		r.Layout(r.popUp.Size())
	}
	if r.popUp.Canvas.Size() != r.popUp.BaseWidget.Size() {
		r.popUp.BaseWidget.Resize(r.popUp.Canvas.Size())
	}
	r.popUp.Content.Refresh()
	r.background.Refresh()
	r.ShadowingRenderer.RefreshShadow()
}

type modalPopUpRenderer struct {
	*widget.ShadowingRenderer
	popUpBaseRenderer
	underlay *canvas.Rectangle
}

func (r *modalPopUpRenderer) Layout(canvasSize gui.Size) {
	r.underlay.Resize(canvasSize)

	padding := r.padding()
	innerSize := r.popUp.innerSize.Max(r.popUp.Content.MinSize().Add(padding))

	requestedSize := innerSize.Subtract(padding)
	size := r.popUp.Content.MinSize().Max(requestedSize)
	size = size.Min(canvasSize.Subtract(padding))
	pos := gui.NewPos((canvasSize.Width-size.Width)/2, (canvasSize.Height-size.Height)/2)
	r.popUp.Content.Move(pos)
	r.popUp.Content.Resize(size)

	innerPos := pos.Subtract(r.offset())
	r.background.Move(innerPos)
	r.background.Resize(size.Add(padding))
	r.LayoutShadow(innerSize, innerPos)
}

func (r *modalPopUpRenderer) MinSize() gui.Size {
	return r.popUp.Content.MinSize().Add(r.padding())
}

func (r *modalPopUpRenderer) Refresh() {
	r.underlay.FillColor = theme.ShadowColor()
	r.background.FillColor = theme.BackgroundColor()
	expectedContentSize := r.popUp.innerSize.Max(r.popUp.MinSize()).Subtract(r.padding())
	shouldLayout := r.popUp.Content.Size() != expectedContentSize

	if r.background.Size() != r.popUp.innerSize || shouldLayout {
		r.Layout(r.popUp.Size())
	}
	if r.popUp.Canvas.Size() != r.popUp.BaseWidget.Size() {
		r.popUp.BaseWidget.Resize(r.popUp.Canvas.Size())
	}
	r.popUp.Content.Refresh()
	r.background.Refresh()
}
