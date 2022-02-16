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

var _ gui.Widget = (*radioItem)(nil)
var _ desktop.Hoverable = (*radioItem)(nil)
var _ gui.Tappable = (*radioItem)(nil)
var _ gui.Focusable = (*radioItem)(nil)

func newRadioItem(label string, onTap func(*radioItem)) *radioItem {
	i := &radioItem{Label: label, onTap: onTap}
	i.ExtendBaseWidget(i)
	return i
}

// radioItem is a single radio item to be used by RadioGroup.
type radioItem struct {
	DisableableWidget

	Label    string
	Selected bool

	focused bool
	hovered bool
	onTap   func(item *radioItem)
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer.
//
// Implements: gui.Widget
func (i *radioItem) CreateRenderer() gui.WidgetRenderer {
	focusIndicator := canvas.NewCircle(theme.BackgroundColor())
	icon := canvas.NewImageFromResource(theme.RadioButtonIcon())
	label := canvas.NewText(i.Label, theme.ForegroundColor())
	label.Alignment = gui.TextAlignLeading
	r := &radioItemRenderer{
		BaseRenderer:   widget.NewBaseRenderer([]gui.CanvasObject{focusIndicator, icon, label}),
		focusIndicator: focusIndicator,
		icon:           icon,
		item:           i,
		label:          label,
	}
	r.update()
	return r
}

// FocusGained is called when this item gained the focus.
//
// Implements: gui.Focusable
func (i *radioItem) FocusGained() {
	i.focused = true
	i.Refresh()
}

// FocusLost is called when this item lost the focus.
//
// Implements: gui.Focusable
func (i *radioItem) FocusLost() {
	i.focused = false
	i.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
//
// Implements: desktop.Hoverable
func (i *radioItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Disabled() {
		return
	}

	i.hovered = true
	i.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
//
// Implements: desktop.Hoverable
func (i *radioItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
//
// Implements: desktop.Hoverable
func (i *radioItem) MouseOut() {
	if i.Disabled() {
		return
	}

	i.hovered = false
	i.Refresh()
}

// SetSelected sets whether this radio item is selected or not.
func (i *radioItem) SetSelected(selected bool) {
	if i.Disabled() {
		return
	}
	if i.Selected == selected {
		return
	}
	i.Selected = selected
	i.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
//
// Implements: gui.Tappable
func (i *radioItem) Tapped(_ *gui.PointEvent) {
	if !i.focused && !gui.CurrentDevice().IsMobile() {
		gui.CurrentApp().Driver().CanvasForObject(i.super()).Focus(i.super().(gui.Focusable))
	}
	i.toggle()
}

// TypedKey is called when this item receives a key event.
//
// Implements: gui.Focusable
func (i *radioItem) TypedKey(_ *gui.KeyEvent) {
}

// TypedRune is called when this item receives a char event.
//
// Implements: gui.Focusable
func (i *radioItem) TypedRune(r rune) {
	if r == ' ' {
		i.toggle()
	}
}

func (i *radioItem) toggle() {
	if i.Disabled() {
		return
	}
	if i.onTap == nil {
		return
	}

	i.onTap(i)
}

type radioItemRenderer struct {
	widget.BaseRenderer

	focusIndicator *canvas.Circle
	icon           *canvas.Image
	item           *radioItem
	label          *canvas.Text
}

func (r *radioItemRenderer) Layout(size gui.Size) {
	labelSize := gui.NewSize(size.Width, size.Height)
	focusIndicatorSize := gui.NewSize(theme.IconInlineSize()+theme.Padding()*2, theme.IconInlineSize()+theme.Padding())

	r.focusIndicator.Resize(focusIndicatorSize)
	r.focusIndicator.Move(gui.NewPos(theme.Padding()*0.5, (size.Height-focusIndicatorSize.Height)/2))

	r.label.Resize(labelSize)
	r.label.Move(gui.NewPos(focusIndicatorSize.Width+theme.Padding(), 0))

	r.icon.Resize(gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(gui.NewPos(theme.Padding()*1.5, (labelSize.Height-theme.IconInlineSize())/2))
}

func (r *radioItemRenderer) MinSize() gui.Size {
	pad4 := theme.Padding() * 4

	return r.label.MinSize().
		Add(gui.NewSize(pad4, pad4)).
		Add(gui.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
}

func (r *radioItemRenderer) Refresh() {
	r.update()
	canvas.Refresh(r.item.super())
}

func (r *radioItemRenderer) update() {
	r.label.Text = r.item.Label
	r.label.Color = theme.ForegroundColor()
	r.label.TextSize = theme.TextSize()
	if r.item.Disabled() {
		r.label.Color = theme.DisabledColor()
	}

	res := theme.RadioButtonIcon()
	if r.item.Selected {
		res = theme.RadioButtonCheckedIcon()
	}
	if r.item.Disabled() {
		res = theme.NewDisabledResource(res)
	}
	r.icon.Resource = res

	if r.item.Disabled() {
		r.focusIndicator.FillColor = theme.BackgroundColor()
	} else if r.item.focused {
		r.focusIndicator.FillColor = theme.FocusColor()
	} else if r.item.hovered {
		r.focusIndicator.FillColor = theme.HoverColor()
	} else {
		r.focusIndicator.FillColor = theme.BackgroundColor()
	}
}
