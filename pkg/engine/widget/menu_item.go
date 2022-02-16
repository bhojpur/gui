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
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Widget = (*menuItem)(nil)

// menuItem is a widget for displaying a gui.menuItem.
type menuItem struct {
	widget.Base
	Item   *gui.MenuItem
	Parent *Menu

	alignment gui.TextAlign
	child     *Menu
}

// newMenuItem creates a new menuItem.
func newMenuItem(item *gui.MenuItem, parent *Menu) *menuItem {
	i := &menuItem{Item: item, Parent: parent}
	i.alignment = parent.alignment
	i.ExtendBaseWidget(i)
	return i
}

func (i *menuItem) Child() *Menu {
	if i.Item.ChildMenu != nil && i.child == nil {
		child := NewMenu(i.Item.ChildMenu)
		child.Hide()
		child.OnDismiss = i.Parent.Dismiss
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu item.
//
// Implements: gui.Widget
func (i *menuItem) CreateRenderer() gui.WidgetRenderer {
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	text := canvas.NewText(i.Item.Label, theme.ForegroundColor())
	text.Alignment = i.alignment
	objects := []gui.CanvasObject{background, text}
	var icon *canvas.Image
	if i.Item.ChildMenu != nil {
		icon = canvas.NewImageFromResource(theme.MenuExpandIcon())
		objects = append(objects, icon)
	}
	checkIcon := canvas.NewImageFromResource(theme.ConfirmIcon())
	if !i.Item.Checked {
		checkIcon.Hide()
	}

	objects = append(objects, checkIcon)
	return &menuItemRenderer{
		BaseRenderer: widget.NewBaseRenderer(objects),
		i:            i,
		icon:         icon,
		checkIcon:    checkIcon,
		text:         text,
		background:   background,
	}
}

// MouseIn activates the item which shows the submenu if the item has one.
// The submenu of any sibling of the item will be hidden.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseIn(*desktop.MouseEvent) {
	i.activate()
}

// MouseMoved does nothing.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut deactivates the item unless it has an open submenu.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseOut() {
	if !i.isSubmenuOpen() {
		i.deactivate()
	}
}

// Tapped performs the action of the item and dismisses the menu.
// It does nothing if the item doesn’t have an action.
//
// Implements: gui.Tappable
func (i *menuItem) Tapped(*gui.PointEvent) {
	if i.Item.Disabled {
		return
	}
	if i.Item.Action == nil {
		if gui.CurrentDevice().IsMobile() {
			i.activate()
		}
		return
	}

	i.trigger()
}

func (i *menuItem) activate() {
	if i.Item.Disabled {
		return
	}
	if i.Child() != nil {
		i.Child().Show()
	}
	i.Parent.activateItem(i)
}

func (i *menuItem) activateLastSubmenu() bool {
	if i.Child() == nil {
		return false
	}
	if i.isSubmenuOpen() {
		return i.Child().ActivateLastSubmenu()
	}
	i.Child().Show()
	i.Child().ActivateNext()
	return true
}

func (i *menuItem) deactivate() {
	if i.Child() != nil {
		i.Child().Hide()
	}
	i.Parent.DeactivateChild()
}

func (i *menuItem) deactivateLastSubmenu() bool {
	if !i.isSubmenuOpen() {
		return false
	}
	if !i.Child().DeactivateLastSubmenu() {
		i.Child().DeactivateChild()
		i.Child().Hide()
	}
	return true
}

func (i *menuItem) isActive() bool {
	return i.Parent.activeItem == i
}

func (i *menuItem) isSubmenuOpen() bool {
	return i.Child() != nil && i.Child().Visible()
}

func (i *menuItem) trigger() {
	i.Parent.Dismiss()
	if i.Item.Action != nil {
		i.Item.Action()
	}
}

func (i *menuItem) triggerLast() {
	if i.isSubmenuOpen() {
		i.Child().TriggerLast()
		return
	}
	i.trigger()
}

type menuItemRenderer struct {
	widget.BaseRenderer
	i                *menuItem
	icon             *canvas.Image
	checkIcon        *canvas.Image
	lastThemePadding float32
	minSize          gui.Size
	text             *canvas.Text
	background       *canvas.Rectangle
}

func (r *menuItemRenderer) Layout(size gui.Size) {
	padding := r.itemPadding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.ForegroundColor()
	if r.i.Item.Disabled {
		r.text.Color = theme.DisabledColor()
	}
	r.text.Resize(size.Subtract(gui.NewSize(theme.Padding()*4, theme.Padding()*2)))
	r.text.Move(gui.NewPos(padding.Width/2+r.checkSpace(), padding.Height/2))

	if r.icon != nil {
		r.icon.Resize(gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(gui.NewPos(size.Width-theme.IconInlineSize(), (size.Height-theme.IconInlineSize())/2))
	}
	r.checkIcon.Resize(gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.checkIcon.Move(gui.NewPos(padding.Width/4, (size.Height-theme.IconInlineSize())/2))

	r.background.Resize(size)
}

func (r *menuItemRenderer) checkSpace() float32 {
	if r.i.Parent.containsCheck {
		return theme.IconInlineSize()
	}
	return 0
}

func (r *menuItemRenderer) MinSize() gui.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	minSize := r.text.MinSize().Add(r.itemPadding()).Add(gui.NewSize(r.checkSpace(), 0))
	if r.icon != nil {
		minSize = minSize.Add(gui.NewSize(theme.IconInlineSize(), 0))
	}
	r.minSize = minSize
	return r.minSize
}

func (r *menuItemRenderer) Refresh() {
	if gui.CurrentDevice().IsMobile() {
		r.background.Hide()
	} else if r.i.isActive() {
		r.background.FillColor = theme.FocusColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
	r.text.Alignment = r.i.alignment
	if r.i.Item.Disabled {
		r.text.Color = theme.DisabledColor()
		r.checkIcon.Resource = theme.NewDisabledResource(theme.ConfirmIcon())
	} else {
		r.text.Color = theme.ForegroundColor()
		r.checkIcon.Resource = theme.ConfirmIcon()
	}
	r.text.Refresh()

	if r.i.Item.Checked {
		r.checkIcon.Show()
	} else {
		r.checkIcon.Hide()
	}
	r.checkIcon.Refresh()
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	return !r.minSize.IsZero() &&
		r.text.TextSize == theme.TextSize() &&
		(r.icon == nil || r.icon.Size().Width == theme.IconInlineSize()) &&
		r.lastThemePadding == theme.Padding()
}

func (r *menuItemRenderer) itemPadding() gui.Size {
	return gui.NewSize(theme.Padding()*4, theme.Padding()*2)
}
