package glfw

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
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Widget = (*MenuBar)(nil)

// MenuBar is a widget for displaying a gui.MainMenu in a bar.
type MenuBar struct {
	widget.Base
	Items []gui.CanvasObject

	active     bool
	activeItem *menuBarItem
	canvas     gui.Canvas
}

// NewMenuBar creates a menu bar populated with items from the passed main menu structure.
func NewMenuBar(mainMenu *gui.MainMenu, canvas gui.Canvas) *MenuBar {
	items := make([]gui.CanvasObject, len(mainMenu.Items))
	b := &MenuBar{Items: items, canvas: canvas}
	b.ExtendBaseWidget(b)
	for i, menu := range mainMenu.Items {
		barItem := &menuBarItem{Menu: menu, Parent: b}
		barItem.ExtendBaseWidget(barItem)
		items[i] = barItem
	}
	return b
}

// CreateRenderer returns a new renderer for the menu bar.
//
// Implements: gui.Widget
func (b *MenuBar) CreateRenderer() gui.WidgetRenderer {
	cont := gui.NewContainerWithLayout(layout.NewHBoxLayout(), b.Items...)
	background := canvas.NewRectangle(theme.ButtonColor())
	underlay := &menuBarUnderlay{action: b.deactivate}
	underlay.ExtendBaseWidget(underlay)
	objects := []gui.CanvasObject{underlay, background, cont}
	for _, item := range b.Items {
		objects = append(objects, item.(*menuBarItem).Child())
	}
	return &menuBarRenderer{
		widget.NewShadowingRenderer(objects, widget.MenuLevel),
		b,
		background,
		underlay,
		cont,
	}
}

// IsActive returns whether the menu bar is active or not.
// An active menu bar shows the current selected menu and should have the focus.
func (b *MenuBar) IsActive() bool {
	return b.active
}

// Toggle changes the activation state of the menu bar.
// On activation, the first item will become active.
func (b *MenuBar) Toggle() {
	b.toggle(b.Items[0].(*menuBarItem))
}

func (b *MenuBar) activateChild(item *menuBarItem) {
	if !b.active {
		b.active = true
	}
	if item.Child() != nil {
		item.Child().DeactivateChild()
	}
	if b.activeItem == item {
		return
	}

	if b.activeItem != nil {
		if c := b.activeItem.Child(); c != nil {
			c.Hide()
		}
		b.activeItem.Refresh()
	}
	b.activeItem = item
	if item == nil {
		return
	}

	item.Refresh()
	item.Child().Show()
	b.Refresh()
}

func (b *MenuBar) deactivate() {
	if !b.active {
		return
	}

	b.active = false
	if b.activeItem != nil {
		if c := b.activeItem.Child(); c != nil {
			defer c.Dismiss()
			c.Hide()
		}
		b.activeItem.Refresh()
		b.activeItem = nil
	}
	b.Refresh()
}

func (b *MenuBar) toggle(item *menuBarItem) {
	if b.active {
		b.canvas.Unfocus()
		b.deactivate()
	} else {
		b.activateChild(item)
		b.canvas.Focus(item)
	}
}

type menuBarRenderer struct {
	*widget.ShadowingRenderer
	b          *MenuBar
	background *canvas.Rectangle
	underlay   *menuBarUnderlay
	cont       *gui.Container
}

func (r *menuBarRenderer) Layout(size gui.Size) {
	r.LayoutShadow(size, gui.NewPos(0, 0))
	minSize := r.MinSize()
	if size.Height != minSize.Height || size.Width < minSize.Width {
		r.b.Resize(gui.NewSize(gui.Max(size.Width, minSize.Width), minSize.Height))
		return
	}

	if r.b.active {
		r.underlay.Resize(r.b.canvas.Size())
	} else {
		r.underlay.Resize(gui.NewSize(0, 0))
	}
	r.cont.Resize(gui.NewSize(size.Width-2*theme.Padding(), size.Height))
	r.cont.Move(gui.NewPos(theme.Padding(), 0))
	if item := r.b.activeItem; item != nil {
		if item.Child().Size().IsZero() {
			item.Child().Resize(item.Child().MinSize())
		}
		item.Child().Move(gui.NewPos(item.Position().X+theme.Padding(), item.Size().Height))
	}
	r.background.Resize(size)
}

func (r *menuBarRenderer) MinSize() gui.Size {
	return r.cont.MinSize().Add(gui.NewSize(theme.Padding()*2, 0))
}

func (r *menuBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	r.background.FillColor = theme.ButtonColor()
	r.background.Refresh()
	r.ShadowingRenderer.RefreshShadow()
	canvas.Refresh(r.b)
}

// Transparent underlay shown as soon as menu is active.
// It catches mouse events outside the menu's objects.
type menuBarUnderlay struct {
	widget.Base
	action func()
}

var _ gui.Widget = (*menuBarUnderlay)(nil)
var _ gui.Tappable = (*menuBarUnderlay)(nil)      // deactivate menu on click outside
var _ desktop.Hoverable = (*menuBarUnderlay)(nil) // block hover events on main content

func (u *menuBarUnderlay) CreateRenderer() gui.WidgetRenderer {
	return &menuUnderlayRenderer{}
}

func (u *menuBarUnderlay) MouseIn(*desktop.MouseEvent) {
}

func (u *menuBarUnderlay) MouseOut() {
}

func (u *menuBarUnderlay) MouseMoved(*desktop.MouseEvent) {
}

func (u *menuBarUnderlay) Tapped(*gui.PointEvent) {
	u.action()
}

type menuUnderlayRenderer struct {
	widget.BaseRenderer
}

var _ gui.WidgetRenderer = (*menuUnderlayRenderer)(nil)

func (r *menuUnderlayRenderer) Layout(gui.Size) {
}

func (r *menuUnderlayRenderer) MinSize() gui.Size {
	return gui.NewSize(0, 0)
}

func (r *menuUnderlayRenderer) Refresh() {
}
