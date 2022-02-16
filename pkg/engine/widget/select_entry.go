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
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// SelectEntry is an input field which supports selecting from a fixed set of options.
type SelectEntry struct {
	Entry
	dropDown *gui.Menu
	popUp    *PopUpMenu
	options  []string
}

// NewSelectEntry creates a SelectEntry.
func NewSelectEntry(options []string) *SelectEntry {
	e := &SelectEntry{options: options}
	e.ExtendBaseWidget(e)
	e.Wrapping = gui.TextTruncate
	return e
}

// CreateRenderer returns a new renderer for this select entry.
//
// Implements: gui.Widget
func (e *SelectEntry) CreateRenderer() gui.WidgetRenderer {
	e.ExtendBaseWidget(e)
	e.SetOptions(e.options)
	return e.Entry.CreateRenderer()
}

// Enable this widget, updating any style or features appropriately.
//
// Implements: gui.DisableableWidget
func (e *SelectEntry) Enable() {
	if e.ActionItem != nil {
		e.ActionItem.(gui.Disableable).Enable()
	}
	e.Entry.Enable()
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
//
// Implements: gui.DisableableWidget
func (e *SelectEntry) Disable() {
	if e.ActionItem != nil {
		e.ActionItem.(gui.Disableable).Disable()
	}
	e.Entry.Disable()
}

// MinSize returns the minimal size of the select entry.
//
// Implements: gui.Widget
func (e *SelectEntry) MinSize() gui.Size {
	e.ExtendBaseWidget(e)
	return e.Entry.MinSize()
}

// Move changes the relative position of the select entry.
//
// Implements: gui.Widget
func (e *SelectEntry) Move(pos gui.Position) {
	e.Entry.Move(pos)
	if e.popUp != nil {
		e.popUp.Move(e.popUpPos())
	}
}

// Resize changes the size of the select entry.
//
// Implements: gui.Widget
func (e *SelectEntry) Resize(size gui.Size) {
	e.Entry.Resize(size)
	if e.popUp != nil {
		e.popUp.Resize(gui.NewSize(size.Width, e.popUp.Size().Height))
	}
}

// SetOptions sets the options the user might select from.
func (e *SelectEntry) SetOptions(options []string) {
	e.options = options
	items := make([]*gui.MenuItem, len(options))
	for i, option := range options {
		option := option // capture
		items[i] = gui.NewMenuItem(option, func() { e.SetText(option) })
	}
	e.dropDown = gui.NewMenu("", items...)

	if e.ActionItem == nil {
		e.ActionItem = e.setupDropDown()
		if e.Disabled() {
			e.ActionItem.(gui.Disableable).Disable()
		}
	}
}

func (e *SelectEntry) popUpPos() gui.Position {
	entryPos := gui.CurrentApp().Driver().AbsolutePositionForObject(e.super())
	return entryPos.Add(gui.NewPos(0, e.Size().Height-theme.InputBorderSize()))
}

func (e *SelectEntry) setupDropDown() *Button {
	dropDownButton := NewButton("", func() {
		c := gui.CurrentApp().Driver().CanvasForObject(e.super())

		e.popUp = NewPopUpMenu(e.dropDown, c)
		e.popUp.ShowAtPosition(e.popUpPos())
		e.popUp.Resize(gui.NewSize(e.Size().Width, e.popUp.MinSize().Height))
	})
	dropDownButton.Importance = LowImportance
	dropDownButton.SetIcon(theme.MenuDropDownIcon())
	return dropDownButton
}
