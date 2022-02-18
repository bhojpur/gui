package render

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
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// CompletionEntry is an Entry with options displayed in a PopUpMenu.
type CompletionEntry struct {
	widget.Entry
	popupMenu     *widget.PopUp
	navigableList *navigableList
	Options       []string
	pause         bool
	itemHeight    float32
}

// NewCompletionEntry creates a new CompletionEntry which creates a popup
// menu that responds to keystrokes to navigate through the items without
// losing the editing ability of the text input.
func NewCompletionEntry(options []string) *CompletionEntry {
	c := &CompletionEntry{Options: options}
	c.ExtendBaseWidget(c)
	return c
}

// HideCompletion hides the completion menu.
func (c *CompletionEntry) HideCompletion() {
	if c.popupMenu != nil {
		c.popupMenu.Hide()
	}
}

// Move changes the relative position of the select entry.
//
// Implements: gui.Widget
func (c *CompletionEntry) Move(pos gui.Position) {
	c.Entry.Move(pos)
	if c.popupMenu != nil {
		c.popupMenu.Resize(c.maxSize())
		c.popupMenu.Move(c.popUpPos())
	}
}

// Refresh the list to update the options to display.
func (c *CompletionEntry) Refresh() {
	c.Entry.Refresh()
	if c.navigableList != nil {
		c.navigableList.SetOptions(c.Options)
	}
}

// SetOptions set the completion list with itemList and update the view.
func (c *CompletionEntry) SetOptions(itemList []string) {
	c.Options = itemList
	c.Refresh()
}

// ShowCompletion displays the completion menu
func (c *CompletionEntry) ShowCompletion() {
	if c.pause {
		return
	}
	if len(c.Options) == 0 {
		c.HideCompletion()
		return
	}

	if c.navigableList == nil {
		c.navigableList = newNavigableList(c.Options, &c.Entry, c.setTextFromMenu, c.HideCompletion)
	}
	holder := gui.CurrentApp().Driver().CanvasForObject(c)

	if c.popupMenu == nil {
		c.popupMenu = widget.NewPopUp(c.navigableList, holder)
	}
	c.popupMenu.Resize(c.maxSize())
	c.popupMenu.ShowAtPosition(c.popUpPos())
	holder.Focus(c.navigableList)
}

// calculate the max size to make the popup to cover everything below the entry
func (c *CompletionEntry) maxSize() gui.Size {
	cnv := gui.CurrentApp().Driver().CanvasForObject(c)

	if c.itemHeight == 0 {
		// set item height to cache
		c.itemHeight = c.navigableList.CreateItem().MinSize().Height
	}

	listheight := float32(len(c.Options))*(c.itemHeight+2*theme.Padding()+theme.SeparatorThicknessSize()) + 2*theme.Padding()
	canvasSize := cnv.Size()
	entrySize := c.Size()
	if canvasSize.Height > listheight {
		return gui.NewSize(entrySize.Width, listheight)
	}

	return gui.NewSize(
		entrySize.Width,
		canvasSize.Height-c.Position().Y-entrySize.Height-theme.InputBorderSize()-theme.Padding())
}

// calculate where the popup should appear
func (c *CompletionEntry) popUpPos() gui.Position {
	entryPos := gui.CurrentApp().Driver().AbsolutePositionForObject(c)
	return entryPos.Add(gui.NewPos(0, c.Size().Height))
}

// Prevent the menu to open when the user validate value from the menu.
func (c *CompletionEntry) setTextFromMenu(s string) {
	c.pause = true
	c.Entry.SetText(s)
	c.Entry.CursorColumn = len([]rune(s))
	c.Entry.Refresh()
	c.pause = false
	c.popupMenu.Hide()
}

type navigableList struct {
	widget.List
	entry           *widget.Entry
	selected        int
	setTextFromMenu func(string)
	hide            func()
	navigating      bool
	items           []string
}

func newNavigableList(items []string, entry *widget.Entry, setTextFromMenu func(string), hide func()) *navigableList {
	n := &navigableList{
		entry:           entry,
		selected:        -1,
		setTextFromMenu: setTextFromMenu,
		hide:            hide,
		items:           items,
	}

	n.List = widget.List{
		Length: func() int {
			return len(n.items)
		},
		CreateItem: func() gui.CanvasObject {
			return widget.NewLabel("")
		},
		UpdateItem: func(i widget.ListItemID, o gui.CanvasObject) {
			o.(*widget.Label).SetText(n.items[i])
		},
		OnSelected: func(id widget.ListItemID) {
			if !n.navigating && id > -1 {
				setTextFromMenu(n.items[id])
			}
			n.navigating = false
		},
	}
	n.ExtendBaseWidget(n)
	return n
}

// Implements: gui.Focusable
func (n *navigableList) FocusGained() {
}

// Implements: gui.Focusable
func (n *navigableList) FocusLost() {
}

func (n *navigableList) SetOptions(items []string) {
	n.Unselect(n.selected)
	n.items = items
	n.Refresh()
	n.selected = -1
}

func (n *navigableList) TypedKey(event *gui.KeyEvent) {
	switch event.Name {
	case gui.KeyDown:
		if n.selected < len(n.items)-1 {
			n.selected++
		} else {
			n.selected = 0
		}
		n.navigating = true
		n.Select(n.selected)

	case gui.KeyUp:
		if n.selected > 0 {
			n.selected--
		} else {
			n.selected = len(n.items) - 1
		}
		n.navigating = true
		n.Select(n.selected)
	case gui.KeyReturn, gui.KeyEnter:
		if n.selected == -1 { // so the user want to submit the entry
			n.hide()
			n.entry.TypedKey(event)
		} else {
			n.navigating = false
			n.OnSelected(n.selected)
		}
	case gui.KeyEscape:
		n.hide()
	default:
		n.entry.TypedKey(event)

	}
}

func (n *navigableList) TypedRune(r rune) {
	n.entry.TypedRune(r)
}
