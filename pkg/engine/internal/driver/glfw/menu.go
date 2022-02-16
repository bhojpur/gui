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
)

func buildMenuOverlay(menus *gui.MainMenu, w *window) gui.CanvasObject {
	if len(menus.Items) == 0 {
		gui.LogError("Main menu must have at least one child menu", nil)
		return nil
	}

	menus = addMissingQuit(menus, w)
	return NewMenuBar(menus, w.canvas)
}

func addMissingQuit(menus *gui.MainMenu, w *window) *gui.MainMenu {
	var lastItem *gui.MenuItem
	if len(menus.Items[0].Items) > 0 {
		lastItem = menus.Items[0].Items[len(menus.Items[0].Items)-1]
		if lastItem.Label == "Quit" {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit { // make sure the first menu always has a quit option
		quitItem := gui.NewMenuItem("Quit", nil)
		quitItem.IsQuit = true
		menus.Items[0].Items = append(menus.Items[0].Items, gui.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menus.Items[0].Items {
		if item.IsQuit && item.Action == nil {
			item.Action = func() {
				for _, win := range w.driver.AllWindows() {
					if glWin, ok := win.(*window); ok {
						glWin.closed(glWin.view())
					} else {
						win.Close() // for test windows
					}
				}
			}
		}
	}
	return menus
}
