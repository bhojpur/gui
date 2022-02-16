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
	"github.com/bhojpur/gui/pkg/engine/theme"
	publicWidget "github.com/bhojpur/gui/pkg/engine/widget"
)

var _ desktop.Hoverable = (*menuBarItem)(nil)
var _ gui.Focusable = (*menuBarItem)(nil)
var _ gui.Widget = (*menuBarItem)(nil)

// menuBarItem is a widget for displaying an item for a gui.Menu in a MenuBar.
type menuBarItem struct {
	widget.Base
	Menu   *gui.Menu
	Parent *MenuBar

	active  bool
	child   *publicWidget.Menu
	hovered bool
}

func (i *menuBarItem) Child() *publicWidget.Menu {
	if i.child == nil {
		child := publicWidget.NewMenu(i.Menu)
		child.Hide()
		child.OnDismiss = i.Parent.deactivate
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu bar item.
//
// Implements: gui.Widget
func (i *menuBarItem) CreateRenderer() gui.WidgetRenderer {
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	text := canvas.NewText(i.Menu.Label, theme.ForegroundColor())
	objects := []gui.CanvasObject{background, text}

	return &menuBarItemRenderer{
		widget.NewBaseRenderer(objects),
		i,
		text,
		background,
	}
}

func (i *menuBarItem) FocusGained() {
	i.active = true
	if i.Parent.active {
		i.Parent.activateChild(i)
	}
	i.Refresh()
}

func (i *menuBarItem) FocusLost() {
	i.active = false
	i.Refresh()
}

func (i *menuBarItem) Focused() bool {
	return i.active
}

// MouseIn activates the item and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
//
// If the bar is not active, the item will be hovered.
//
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseIn(_ *desktop.MouseEvent) {
	i.hovered = true
	if i.Parent.active {
		i.Parent.canvas.Focus(i)
	}
	i.Refresh()
}

// MouseMoved activates the item and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
// This might have an effect when mouse and keyboard control are mixed.
// Changing the active menu with the keyboard will make the hovered menu bar item inactive.
// On the next mouse move the hovered item is activated again.
//
// If the bar is not active, this will do nothing.
//
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseMoved(_ *desktop.MouseEvent) {
	if i.Parent.active {
		i.Parent.canvas.Focus(i)
	}
}

// MouseOut does nothing if the bar is active.
//
// IF the bar is not active, it changes the item to not be hovered.
//
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Tapped toggles the activation state of the menu bar.
// It shows the item’s menu if the bar is activated and hides it if the bar is deactivated.
//
// Implements: gui.Tappable
func (i *menuBarItem) Tapped(*gui.PointEvent) {
	i.Parent.toggle(i)
}

func (i *menuBarItem) TypedKey(event *gui.KeyEvent) {
	switch event.Name {
	case gui.KeyLeft:
		if !i.Child().DeactivateLastSubmenu() {
			i.Parent.canvas.FocusPrevious()
		}
	case gui.KeyRight:
		if !i.Child().ActivateLastSubmenu() {
			i.Parent.canvas.FocusNext()
		}
	case gui.KeyDown:
		i.Child().ActivateNext()
	case gui.KeyUp:
		i.Child().ActivatePrevious()
	case gui.KeyEnter, gui.KeyReturn, gui.KeySpace:
		i.Child().TriggerLast()
	}
}

func (i *menuBarItem) TypedRune(_ rune) {
}

type menuBarItemRenderer struct {
	widget.BaseRenderer
	i          *menuBarItem
	text       *canvas.Text
	background *canvas.Rectangle
}

func (r *menuBarItemRenderer) Layout(size gui.Size) {
	padding := r.padding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.ForegroundColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(gui.NewPos(padding.Width/2, padding.Height/2))

	r.background.Resize(size)
}

func (r *menuBarItemRenderer) MinSize() gui.Size {
	return r.text.MinSize().Add(r.padding())
}

func (r *menuBarItemRenderer) Refresh() {
	if r.i.active && r.i.Parent.active {
		r.background.FillColor = theme.FocusColor()
		r.background.Show()
	} else if r.i.hovered && !r.i.Parent.active {
		r.background.FillColor = theme.HoverColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
	canvas.Refresh(r.i)
}

func (r *menuBarItemRenderer) padding() gui.Size {
	return gui.NewSize(theme.Padding()*4, theme.Padding()*2)
}
