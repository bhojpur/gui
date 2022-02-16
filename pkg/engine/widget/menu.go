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
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Widget = (*Menu)(nil)
var _ gui.Tappable = (*Menu)(nil)

// Menu is a widget for displaying a gui.Menu.
type Menu struct {
	BaseWidget
	alignment     gui.TextAlign
	Items         []gui.CanvasObject
	OnDismiss     func()
	activeItem    *menuItem
	customSized   bool
	containsCheck bool
}

// NewMenu creates a new Menu.
func NewMenu(menu *gui.Menu) *Menu {
	m := &Menu{}
	m.ExtendBaseWidget(m)
	m.setMenu(menu)
	return m
}

// ActivateLastSubmenu finds the last active menu item traversing through the open submenus
// and activates its submenu if any.
// It returns `true` if there was a submenu and it was activated and `false` elsewhere.
// Activating a submenu does show it and activate its first item.
func (m *Menu) ActivateLastSubmenu() bool {
	if m.activeItem == nil {
		return false
	}
	if !m.activeItem.activateLastSubmenu() {
		return false
	}
	m.Refresh()
	return true
}

// ActivateNext activates the menu item following the currently active menu item.
// If there is no menu item active, it activates the first menu item.
// If there is no menu item after the current active one, it does nothing.
// If a submenu is open, it delegates the activation to this submenu.
func (m *Menu) ActivateNext() {
	if m.activeItem != nil && m.activeItem.isSubmenuOpen() {
		m.activeItem.Child().ActivateNext()
		return
	}

	found := m.activeItem == nil
	for _, item := range m.Items {
		if mItem, ok := item.(*menuItem); ok {
			if found {
				m.activateItem(mItem)
				return
			}
			if mItem == m.activeItem {
				found = true
			}
		}
	}
}

// ActivatePrevious activates the menu item preceding the currently active menu item.
// If there is no menu item active, it activates the last menu item.
// If there is no menu item before the current active one, it does nothing.
// If a submenu is open, it delegates the activation to this submenu.
func (m *Menu) ActivatePrevious() {
	if m.activeItem != nil && m.activeItem.isSubmenuOpen() {
		m.activeItem.Child().ActivatePrevious()
		return
	}

	found := m.activeItem == nil
	for i := len(m.Items) - 1; i >= 0; i-- {
		item := m.Items[i]
		if mItem, ok := item.(*menuItem); ok {
			if found {
				m.activateItem(mItem)
				return
			}
			if mItem == m.activeItem {
				found = true
			}
		}
	}
}

// CreateRenderer returns a new renderer for the menu.
//
// Implements: gui.Widget
func (m *Menu) CreateRenderer() gui.WidgetRenderer {
	m.ExtendBaseWidget(m)
	box := newMenuBox(m.Items)
	scroll := widget.NewVScroll(box)
	scroll.SetMinSize(box.MinSize())
	objects := []gui.CanvasObject{scroll}
	for _, i := range m.Items {
		if item, ok := i.(*menuItem); ok && item.Child() != nil {
			objects = append(objects, item.Child())
		}
	}

	return &menuRenderer{
		widget.NewShadowingRenderer(objects, widget.MenuLevel),
		box,
		m,
		scroll,
	}
}

// DeactivateChild deactivates the active menu item and hides its submenu if any.
func (m *Menu) DeactivateChild() {
	if m.activeItem != nil {
		defer m.activeItem.Refresh()
		if c := m.activeItem.Child(); c != nil {
			c.Hide()
		}
		m.activeItem = nil
	}
}

// DeactivateLastSubmenu finds the last open submenu traversing through the open submenus,
// deactivates its active item and hides it.
// This also deactivates any submenus of the deactivated submenu.
// It returns `true` if there was a submenu open and closed and `false` elsewhere.
func (m *Menu) DeactivateLastSubmenu() bool {
	if m.activeItem == nil {
		return false
	}
	return m.activeItem.deactivateLastSubmenu()
}

// MinSize returns the minimal size of the menu.
//
// Implements: gui.Widget
func (m *Menu) MinSize() gui.Size {
	m.ExtendBaseWidget(m)
	return m.BaseWidget.MinSize()
}

// Refresh updates the menu to reflect changes in the data.
//
// Implements: gui.Widget
func (m *Menu) Refresh() {
	for _, item := range m.Items {
		item.Refresh()
	}
	m.BaseWidget.Refresh()
}

func (m *Menu) getContainsCheck() bool {
	for _, item := range m.Items {
		if mi, ok := item.(*menuItem); ok && mi.Item.Checked {
			return true
		}
	}
	return false
}

// Tapped catches taps on separators and the menu background. It doesn’t perform any action.
//
// Implements: gui.Tappable
func (m *Menu) Tapped(*gui.PointEvent) {
	// Hit a separator or padding -> do nothing.
}

// TriggerLast finds the last active menu item traversing through the open submenus and triggers it.
func (m *Menu) TriggerLast() {
	if m.activeItem == nil {
		m.Dismiss()
		return
	}
	m.activeItem.triggerLast()
}

// Dismiss dismisses the menu by dismissing and hiding the active child and performing OnDismiss.
func (m *Menu) Dismiss() {
	if m.activeItem != nil {
		if m.activeItem.Child() != nil {
			defer m.activeItem.Child().Dismiss()
		}
		m.DeactivateChild()
	}
	if m.OnDismiss != nil {
		m.OnDismiss()
	}
}

func (m *Menu) activateItem(item *menuItem) {
	if item.Child() != nil {
		item.Child().DeactivateChild()
	}
	if m.activeItem == item {
		return
	}

	m.DeactivateChild()
	m.activeItem = item
	m.activeItem.Refresh()
	m.Refresh()
}

func (m *Menu) setMenu(menu *gui.Menu) {
	m.Items = make([]gui.CanvasObject, len(menu.Items))
	for i, item := range menu.Items {
		if item.IsSeparator {
			m.Items[i] = NewSeparator()
		} else {
			m.Items[i] = newMenuItem(item, m)
		}
	}
	m.containsCheck = m.getContainsCheck()
}

type menuRenderer struct {
	*widget.ShadowingRenderer
	box    *menuBox
	m      *Menu
	scroll *widget.Scroll
}

func (r *menuRenderer) Layout(s gui.Size) {
	minSize := r.MinSize()
	var boxSize gui.Size
	if r.m.customSized {
		boxSize = minSize.Max(s)
	} else {
		boxSize = minSize
	}
	scrollSize := boxSize
	if c := gui.CurrentApp().Driver().CanvasForObject(r.m.super()); c != nil {
		ap := gui.CurrentApp().Driver().AbsolutePositionForObject(r.m.super())
		pos, size := c.InteractiveArea()
		bottomPad := c.Size().Height - pos.Y - size.Height
		if ah := c.Size().Height - bottomPad - ap.Y; ah < boxSize.Height {
			scrollSize = gui.NewSize(boxSize.Width, ah)
		}
	}
	if scrollSize != r.m.Size() {
		r.m.Resize(scrollSize)
		return
	}

	r.LayoutShadow(scrollSize, gui.NewPos(0, 0))
	r.scroll.Resize(scrollSize)
	r.box.Resize(boxSize)
	r.layoutActiveChild()
}

func (r *menuRenderer) MinSize() gui.Size {
	return r.box.MinSize()
}

func (r *menuRenderer) Refresh() {
	r.layoutActiveChild()
	r.ShadowingRenderer.RefreshShadow()

	for _, i := range r.m.Items {
		if txt, ok := i.(*menuItem); ok {
			txt.alignment = r.m.alignment
			txt.Refresh()
		}
	}

	canvas.Refresh(r.m)
}

func (r *menuRenderer) layoutActiveChild() {
	item := r.m.activeItem
	if item == nil || item.Child() == nil {
		return
	}

	if item.Child().Size().IsZero() {
		item.Child().Resize(item.Child().MinSize())
	}

	itemSize := item.Size()
	cp := gui.NewPos(itemSize.Width, item.Position().Y)
	d := gui.CurrentApp().Driver()
	c := d.CanvasForObject(item)
	if c != nil {
		absPos := d.AbsolutePositionForObject(item)
		childSize := item.Child().Size()
		if absPos.X+itemSize.Width+childSize.Width > c.Size().Width {
			if absPos.X-childSize.Width >= 0 {
				cp.X = -childSize.Width
			} else {
				cp.X = c.Size().Width - absPos.X - childSize.Width
			}
		}
		requiredHeight := childSize.Height - theme.Padding()
		availableHeight := c.Size().Height - absPos.Y
		missingHeight := requiredHeight - availableHeight
		if missingHeight > 0 {
			cp.Y -= missingHeight
		}
	}
	item.Child().Move(cp)
}

type menuBox struct {
	BaseWidget
	items []gui.CanvasObject
}

var _ gui.Widget = (*menuBox)(nil)

func newMenuBox(items []gui.CanvasObject) *menuBox {
	b := &menuBox{items: items}
	b.ExtendBaseWidget(b)
	return b
}

func (b *menuBox) CreateRenderer() gui.WidgetRenderer {
	background := canvas.NewRectangle(theme.BackgroundColor())
	cont := gui.NewContainerWithLayout(layout.NewVBoxLayout(), b.items...)
	return &menuBoxRenderer{
		BaseRenderer: widget.NewBaseRenderer([]gui.CanvasObject{background, cont}),
		b:            b,
		background:   background,
		cont:         cont,
	}
}

type menuBoxRenderer struct {
	widget.BaseRenderer
	b          *menuBox
	background *canvas.Rectangle
	cont       *gui.Container
}

var _ gui.WidgetRenderer = (*menuBoxRenderer)(nil)

func (r *menuBoxRenderer) Layout(size gui.Size) {
	s := gui.NewSize(size.Width, size.Height+2*theme.Padding())
	r.background.Resize(s)
	r.cont.Resize(s)
	r.cont.Move(gui.NewPos(0, theme.Padding()))
}

func (r *menuBoxRenderer) MinSize() gui.Size {
	return r.cont.MinSize().Add(gui.NewSize(0, 2*theme.Padding()))
}

func (r *menuBoxRenderer) Refresh() {
	r.background.FillColor = theme.BackgroundColor()
	r.background.Refresh()
	canvas.Refresh(r.b)
}
