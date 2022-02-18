package engine

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

// Menu stores the information required for a standard menu.
// A menu can pop down from a MainMenu or could be a pop out menu.
type Menu struct {
	Label string
	Items []*MenuItem
}

// NewMenu creates a new menu given the specified label (to show in a MainMenu) and list of items to display.
func NewMenu(label string, items ...*MenuItem) *Menu {
	return &Menu{Label: label, Items: items}
}

// MenuItem is a single item within any menu, it contains a display Label and Action function that is called when tapped.
type MenuItem struct {
	ChildMenu *Menu
	// Since: 2.1
	IsQuit      bool
	IsSeparator bool
	Label       string
	Action      func()
	// Since: 2.1
	Disabled bool
	// Since: 2.1
	Checked bool
	// Since: 2.2
	Shortcut Shortcut
}

// NewMenuItem creates a new menu item from the passed label and action parameters.
func NewMenuItem(label string, action func()) *MenuItem {
	return &MenuItem{Label: label, Action: action}
}

// NewMenuItemSeparator creates a menu item that is to be used as a separator.
func NewMenuItemSeparator() *MenuItem {
	return &MenuItem{IsSeparator: true, Action: func() {}}
}

// MainMenu defines the data required to show a menu bar (desktop) or other appropriate top level menu.
type MainMenu struct {
	Items []*Menu
}

// NewMainMenu creates a top level menu structure used by gui.Window for displaying a menubar
// (or appropriate equivalent).
func NewMainMenu(items ...*Menu) *MainMenu {
	return &MainMenu{items}
}
