//go:build !ci && !no_native_menus && !mobile
// +build !ci,!no_native_menus,!mobile

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
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
)

func TestDarwinMenu(t *testing.T) {
	setExceptionCallback(func(msg string) { t.Error("Obj-C exception:", msg) })
	defer setExceptionCallback(nil)

	resetMainMenu()

	w := createWindow("Test").(*window)

	var lastAction string
	assertNSMenuItem := func(wantTitle, wantAction string, m unsafe.Pointer, i int) {
		item := testNSMenuItemAtIndex(m, i)
		assert.Equal(t, wantTitle, testNSMenuItemTitle(item))
		testNSMenuPerformActionForItemAtIndex(m, i)
		w.WaitForEvents()
		assert.Equal(t, wantAction, lastAction)
	}

	assertNSMenuItemSeparator := func(m unsafe.Pointer, i int) {
		item := testNSMenuItemAtIndex(m, i)
		assert.True(t, testNSMenuItemIsSeparatorItem(item), "item is expected to be a separator")
	}

	itemNew := gui.NewMenuItem("New", func() { lastAction = "new" })
	itemOpen := gui.NewMenuItem("Open", func() { lastAction = "open" })
	itemRecent := gui.NewMenuItem("Recent", nil)
	itemFoo := gui.NewMenuItem("Foo", func() { lastAction = "foo" })
	itemRecent.ChildMenu = gui.NewMenu("", itemFoo)
	menuEdit := gui.NewMenu("File", itemNew, itemOpen, gui.NewMenuItemSeparator(), itemRecent)

	itemHelp := gui.NewMenuItem("Help", func() { lastAction = "Help!!!" })
	itemHelpMe := gui.NewMenuItem("Help Me", func() { lastAction = "Help me!!!" })
	menuHelp := gui.NewMenu("Help", itemHelp, itemHelpMe)

	itemHelloWorld := gui.NewMenuItem("Hello World", func() { lastAction = "Hello World!" })
	itemPrefs := gui.NewMenuItem("Preferences", func() { lastAction = "prefs" })
	itemMore := gui.NewMenuItem("More", func() { lastAction = "more" })
	itemMorePrefs := gui.NewMenuItem("Preferences…", func() { lastAction = "more prefs" })
	menuMore := gui.NewMenu("More Stuff", itemHelloWorld, itemPrefs, itemMore, itemMorePrefs)

	itemSettings := gui.NewMenuItem("Settings", func() { lastAction = "settings" })
	itemMoreSetings := gui.NewMenuItem("Settings…", func() { lastAction = "more settings" })
	menuSettings := gui.NewMenu("Settings", itemSettings, gui.NewMenuItemSeparator(), itemMoreSetings)

	mainMenu := gui.NewMainMenu(menuEdit, menuHelp, menuMore, menuSettings)
	setupNativeMenu(w, mainMenu)
	fmt.Println(lastAction)

	mm := testDarwinMainMenu()
	// The custom “Preferences” menu should be moved to the system app menu completely.
	// -> only three custom menus
	assert.Equal(t, 5, testNSMenuNumberOfItems(mm), "two built-in + three custom")

	m := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 0))
	assert.Equal(t, "", testNSMenuTitle(m), "app menu doesn’t have a title")
	assertNSMenuItemSeparator(m, 1)
	assertNSMenuItem("Preferences", "prefs", m, 2)
	assertNSMenuItem("Preferences…", "more prefs", m, 3)
	assertNSMenuItemSeparator(m, 4)
	assertNSMenuItem("Settings", "settings", m, 5)
	assertNSMenuItem("Settings…", "more settings", m, 6)

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assert.Equal(t, "File", testNSMenuTitle(m))
	assert.Equal(t, 4, testNSMenuNumberOfItems(m))
	assertNSMenuItem("New", "new", m, 0)
	assertNSMenuItem("Open", "open", m, 1)
	assertNSMenuItemSeparator(m, 2)
	i := testNSMenuItemAtIndex(m, 3)
	assert.Equal(t, "Recent", testNSMenuItemTitle(i))
	sm := testNSMenuItemSubmenu(i)
	assert.NotNil(t, sm, "item has submenu")
	assert.Equal(t, 1, testNSMenuNumberOfItems(sm))
	assertNSMenuItem("Foo", "foo", sm, 0)

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 2))
	assert.Equal(t, "More Stuff", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	assertNSMenuItem("Hello World", "Hello World!", m, 0)
	assertNSMenuItem("More", "more", m, 1)

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 3))
	assert.Equal(t, "Window", testNSMenuTitle(m))

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 4))
	assert.Equal(t, "Help", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	assertNSMenuItem("Help", "Help!!!", m, 0)
	assertNSMenuItem("Help Me", "Help me!!!", m, 1)

	// change action works
	itemOpen.Action = func() { lastAction = "new open" }
	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assertNSMenuItem("Open", "new open", m, 1)
}

var initialAppMenuItems []string
var initialMenus []string

func initMainMenu() {
	createWindow("Test").Close() // ensure GLFW has performed [NSApp run]
	mainMenu := testDarwinMainMenu()
	for i := 0; i < testNSMenuNumberOfItems(mainMenu); i++ {
		menu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, i))
		initialMenus = append(initialMenus, testNSMenuTitle(menu))
	}
	appMenu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, 0))
	for i := 0; i < testNSMenuNumberOfItems(appMenu); i++ {
		item := testNSMenuItemAtIndex(appMenu, i)
		initialAppMenuItems = append(initialAppMenuItems, testNSMenuItemTitle(item))
	}
}

func resetMainMenu() {
	mainMenu := testDarwinMainMenu()
	j := 0
	for i := 0; i < testNSMenuNumberOfItems(mainMenu); i++ {
		menu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, i))
		if j < len(initialMenus) && testNSMenuTitle(menu) == initialMenus[j] {
			j++
			continue
		}
		testNSMenuRemoveItemAtIndex(mainMenu, i)
		i--
	}
	j = 0
	appMenu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, 0))
	for i := 0; i < testNSMenuNumberOfItems(appMenu); i++ {
		item := testNSMenuItemAtIndex(appMenu, i)
		if testNSMenuItemTitle(item) == initialAppMenuItems[j] {
			j++
			continue
		}
		testNSMenuRemoveItemAtIndex(appMenu, i)
		i--
	}
}
