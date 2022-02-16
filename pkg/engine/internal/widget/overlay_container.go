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
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
)

var _ gui.Widget = (*OverlayContainer)(nil)
var _ gui.Tappable = (*OverlayContainer)(nil)

// OverlayContainer is a transparent widget containing one gui.CanvasObject and meant to be used as overlay.
type OverlayContainer struct {
	Base
	Content gui.CanvasObject

	canvas    gui.Canvas
	onDismiss func()
	shown     bool
}

// NewOverlayContainer creates an OverlayContainer.
func NewOverlayContainer(c gui.CanvasObject, canvas gui.Canvas, onDismiss func()) *OverlayContainer {
	o := &OverlayContainer{canvas: canvas, Content: c, onDismiss: onDismiss}
	o.ExtendBaseWidget(o)
	return o
}

// CreateRenderer returns a new renderer for the overlay container.
//
// Implements: gui.Widget
func (o *OverlayContainer) CreateRenderer() gui.WidgetRenderer {
	return &overlayRenderer{BaseRenderer{[]gui.CanvasObject{o.Content}}, o}
}

// Hide hides the overlay container.
//
// Implements: gui.Widget
func (o *OverlayContainer) Hide() {
	if o.shown {
		o.canvas.Overlays().Remove(o)
		o.shown = false
	}
	o.Base.Hide()
}

// MouseIn catches mouse-in events not handled by the container’s content. It does nothing.
//
// Implements: desktop.Hoverable
func (o *OverlayContainer) MouseIn(*desktop.MouseEvent) {
}

// MouseMoved catches mouse-moved events not handled by the container’s content. It does nothing.
//
// Implements: desktop.Hoverable
func (o *OverlayContainer) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut catches mouse-out events not handled by the container’s content. It does nothing.
//
// Implements: desktop.Hoverable
func (o *OverlayContainer) MouseOut() {
}

// Show makes the overlay container visible.
//
// Implements: gui.Widget
func (o *OverlayContainer) Show() {
	if !o.shown {
		o.canvas.Overlays().Add(o)
		o.shown = true
	}
	o.Base.Show()
}

// Tapped catches tap events not handled by the container’s content.
// It performs the overlay container’s dismiss action.
//
// Implements: gui.Tappable
func (o *OverlayContainer) Tapped(*gui.PointEvent) {
	if o.onDismiss != nil {
		o.onDismiss()
	}
}

type overlayRenderer struct {
	BaseRenderer
	o *OverlayContainer
}

var _ gui.WidgetRenderer = (*overlayRenderer)(nil)

func (r *overlayRenderer) Layout(gui.Size) {
}

func (r *overlayRenderer) MinSize() gui.Size {
	return r.o.canvas.Size()
}

func (r *overlayRenderer) Refresh() {
}
