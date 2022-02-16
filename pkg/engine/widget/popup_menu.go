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
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
)

var _ gui.Widget = (*PopUpMenu)(nil)
var _ gui.Focusable = (*PopUpMenu)(nil)

// PopUpMenu is a Menu which displays itself in an OverlayContainer.
type PopUpMenu struct {
	*Menu
	canvas  gui.Canvas
	overlay *widget.OverlayContainer
}

// NewPopUpMenu creates a new, reusable popup menu. You can show it using ShowAtPosition.
//
// Since: 2.0
func NewPopUpMenu(menu *gui.Menu, c gui.Canvas) *PopUpMenu {
	m := &Menu{}
	m.setMenu(menu)
	p := &PopUpMenu{Menu: m, canvas: c}
	p.ExtendBaseWidget(p)
	p.Menu.Resize(p.Menu.MinSize())
	p.Menu.customSized = true
	o := widget.NewOverlayContainer(p, c, p.Dismiss)
	o.Resize(o.MinSize())
	p.overlay = o
	p.OnDismiss = func() {
		p.Hide()
	}
	return p
}

// ShowPopUpMenuAtPosition creates a PopUp menu populated with items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func ShowPopUpMenuAtPosition(menu *gui.Menu, c gui.Canvas, pos gui.Position) {
	m := NewPopUpMenu(menu, c)
	m.ShowAtPosition(pos)
}

// FocusGained is triggered when the object gained focus. For the pop-up menu it does nothing.
//
// Implements: gui.Focusable
func (p *PopUpMenu) FocusGained() {}

// FocusLost is triggered when the object lost focus. For the pop-up menu it does nothing.
//
// Implements: gui.Focusable
func (p *PopUpMenu) FocusLost() {}

// Hide hides the pop-up menu.
//
// Implements: gui.Widget
func (p *PopUpMenu) Hide() {
	p.overlay.Hide()
	p.Menu.Hide()
}

// Move moves the pop-up menu.
// The position is absolute because pop-up menus are shown in an overlay which covers the whole canvas.
//
// Implements: gui.Widget
func (p *PopUpMenu) Move(pos gui.Position) {
	p.BaseWidget.Move(p.adjustedPosition(pos, p.Size()))
}

// Resize changes the size of the pop-up menu.
//
// Implements: gui.Widget
func (p *PopUpMenu) Resize(size gui.Size) {
	p.BaseWidget.Move(p.adjustedPosition(p.Position(), size))
	p.Menu.Resize(size)
}

// Show makes the pop-up menu visible.
//
// Implements: gui.Widget
func (p *PopUpMenu) Show() {
	p.Menu.alignment = p.alignment
	p.Menu.Refresh()

	p.overlay.Show()
	p.Menu.Show()
	if !gui.CurrentDevice().IsMobile() {
		p.canvas.Focus(p)
	}
}

// ShowAtPosition shows the pop-up menu at the specified position.
func (p *PopUpMenu) ShowAtPosition(pos gui.Position) {
	p.Move(pos)
	p.Show()
}

// TypedKey handles key events. It allows keyboard control of the pop-up menu.
//
// Implements: gui.Focusable
func (p *PopUpMenu) TypedKey(e *gui.KeyEvent) {
	switch e.Name {
	case gui.KeyDown:
		p.ActivateNext()
	case gui.KeyEnter, gui.KeyReturn, gui.KeySpace:
		p.TriggerLast()
	case gui.KeyEscape:
		p.Dismiss()
	case gui.KeyLeft:
		p.DeactivateLastSubmenu()
	case gui.KeyRight:
		p.ActivateLastSubmenu()
	case gui.KeyUp:
		p.ActivatePrevious()
	}
}

// TypedRune handles text events. For pop-up menus this does nothing.
//
// Implements: gui.Focusable
func (p *PopUpMenu) TypedRune(rune) {}

func (p *PopUpMenu) adjustedPosition(pos gui.Position, size gui.Size) gui.Position {
	x := pos.X
	y := pos.Y
	if x+size.Width > p.canvas.Size().Width {
		x = p.canvas.Size().Width - size.Width
		if x < 0 {
			x = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if y+size.Height > p.canvas.Size().Height {
		y = p.canvas.Size().Height - size.Height
		if y < 0 {
			y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	return gui.NewPos(x, y)
}
