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
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
)

func TestDarwinMenu(t *testing.T) {
	setExceptionCallback(func(msg string) { t.Error("Obj-C exception:", msg) })
	defer setExceptionCallback(nil)

	resetMainMenu()

	w := createWindow("Test").(*window)

	var lastAction string
	assertLastAction := func(wantAction string) {
		w.WaitForEvents()
		assert.Equal(t, wantAction, lastAction)
	}

	assertNSMenuItemSeparator := func(m unsafe.Pointer, i int) {
		item := testNSMenuItemAtIndex(m, i)
		assert.True(t, testNSMenuItemIsSeparatorItem(item), "item is expected to be a separator")
	}

	itemNew := gui.NewMenuItem("New", func() { lastAction = "new" })
	itemNew.Shortcut = &desktop.CustomShortcut{KeyName: gui.KeyN, Modifier: gui.KeyModifierShortcutDefault}
	itemOpen := gui.NewMenuItem("Open", func() { lastAction = "open" })
	itemOpen.Shortcut = &desktop.CustomShortcut{KeyName: gui.KeyO, Modifier: gui.KeyModifierAlt}
	itemRecent := gui.NewMenuItem("Recent", nil)
	itemFoo := gui.NewMenuItem("Foo", func() { lastAction = "foo" })
	itemRecent.ChildMenu = gui.NewMenu("", itemFoo)
	menuEdit := gui.NewMenu("File", itemNew, itemOpen, gui.NewMenuItemSeparator(), itemRecent)

	itemHelp := gui.NewMenuItem("Help", func() { lastAction = "Help!!!" })
	itemHelp.Shortcut = &desktop.CustomShortcut{KeyName: gui.KeyH, Modifier: gui.KeyModifierControl}
	itemHelpMe := gui.NewMenuItem("Help Me", func() { lastAction = "Help me!!!" })
	itemHelpMe.Shortcut = &desktop.CustomShortcut{KeyName: gui.KeyH, Modifier: gui.KeyModifierShift}
	menuHelp := gui.NewMenu("Help", itemHelp, itemHelpMe)

	itemHelloWorld := gui.NewMenuItem("Hello World", func() { lastAction = "Hello World!" })
	itemHelloWorld.Shortcut = &desktop.CustomShortcut{KeyName: gui.KeyH, Modifier: gui.KeyModifierControl | gui.KeyModifierAlt | gui.KeyModifierShift | gui.KeyModifierSuper}
	itemPrefs := gui.NewMenuItem("Preferences", func() { lastAction = "prefs" })
	itemMore := gui.NewMenuItem("More", func() { lastAction = "more" })
	itemMorePrefs := gui.NewMenuItem("Preferences…", func() { lastAction = "more prefs" })
	menuMore := gui.NewMenu("More Stuff", itemHelloWorld, itemPrefs, itemMore, itemMorePrefs)

	itemSettings := gui.NewMenuItem("Settings", func() { lastAction = "settings" })
	itemMoreSetings := gui.NewMenuItem("Settings…", func() { lastAction = "more settings" })
	menuSettings := gui.NewMenu("Settings", itemSettings, gui.NewMenuItemSeparator(), itemMoreSetings)

	mainMenu := gui.NewMainMenu(menuEdit, menuHelp, menuMore, menuSettings)
	setupNativeMenu(w, mainMenu)

	mm := testDarwinMainMenu()
	// The custom “Preferences” menu should be moved to the system app menu completely.
	// -> only three custom menus
	assert.Equal(t, 5, testNSMenuNumberOfItems(mm), "two built-in + three custom")

	m := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 0))
	assert.Equal(t, "", testNSMenuTitle(m), "app menu doesn’t have a title")
	assertNSMenuItemSeparator(m, 1)
	assertNSMenuItem(t, "Preferences", "", 0, m, 2)
	assertLastAction("prefs")
	assertNSMenuItem(t, "Preferences…", "", 0, m, 3)
	assertLastAction("more prefs")
	assertNSMenuItemSeparator(m, 4)
	assertNSMenuItem(t, "Settings", "", 0, m, 5)
	assertLastAction("settings")
	assertNSMenuItem(t, "Settings…", "", 0, m, 6)
	assertLastAction("more settings")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assert.Equal(t, "File", testNSMenuTitle(m))
	assert.Equal(t, 4, testNSMenuNumberOfItems(m))
	// NSEventModifierFlagCommand = 1 << 20
	assertNSMenuItem(t, "New", "n", 0b100000000000000000000, m, 0)
	assertLastAction("new")
	// NSEventModifierFlagOption = 1 << 19
	assertNSMenuItem(t, "Open", "o", 0b10000000000000000000, m, 1)
	assertLastAction("open")
	assertNSMenuItemSeparator(m, 2)
	i := testNSMenuItemAtIndex(m, 3)
	assert.Equal(t, "Recent", testNSMenuItemTitle(i))
	sm := testNSMenuItemSubmenu(i)
	assert.NotNil(t, sm, "item has submenu")
	assert.Equal(t, 1, testNSMenuNumberOfItems(sm))
	assertNSMenuItem(t, "Foo", "", 0, sm, 0)
	assertLastAction("foo")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 2))
	assert.Equal(t, "More Stuff", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	assertNSMenuItem(t, "Hello World", "h", 0b111100000000000000000, m, 0)
	assertLastAction("Hello World!")
	assertNSMenuItem(t, "More", "", 0, m, 1)
	assertLastAction("more")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 3))
	assert.Equal(t, "Window", testNSMenuTitle(m))

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 4))
	assert.Equal(t, "Help", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	// NSEventModifierFlagControl = 1 << 18
	assertNSMenuItem(t, "Help", "h", 0b1000000000000000000, m, 0)
	assertLastAction("Help!!!")
	// NSEventModifierFlagShift = 1 << 17
	assertNSMenuItem(t, "Help Me", "h", 0b100000000000000000, m, 1)
	assertLastAction("Help me!!!")

	// change action works
	itemOpen.Action = func() { lastAction = "new open" }
	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assertNSMenuItem(t, "Open", "", 0, m, 1)
	assertLastAction("new open")
}

func TestDarwinMenu_specialKeyShortcuts(t *testing.T) {
	setExceptionCallback(func(msg string) { t.Error("Obj-C exception:", msg) })
	defer setExceptionCallback(nil)

	for name, tt := range map[string]struct {
		key     gui.KeyName
		wantKey string
	}{
		"Backspace": {
			key:     gui.KeyBackspace,
			wantKey: "\x08", // NSBackspaceCharacter
		},
		"Delete": {
			key:     gui.KeyDelete,
			wantKey: "\x7f", // NSDeleteCharacter
		},
		"Down": {
			key:     gui.KeyDown,
			wantKey: "\uf701", // NSDownArrowFunctionKey
		},
		"End": {
			key:     gui.KeyEnd,
			wantKey: "\uf72b", // NSEndFunctionKey
		},
		"Enter": {
			key:     gui.KeyEnter,
			wantKey: "\x03", // NSEnterCharacter
		},
		"Escape": {
			key:     gui.KeyEscape,
			wantKey: "\x1b", // escape
		},
		"F10": {
			key:     gui.KeyF10,
			wantKey: "\uf70d", // NSF10FunctionKey
		},
		"F11": {
			key:     gui.KeyF11,
			wantKey: "\uf70e", // NSF11FunctionKey
		},
		"F12": {
			key:     gui.KeyF12,
			wantKey: "\uf70f", // NSF12FunctionKey
		},
		"F1": {
			key:     gui.KeyF1,
			wantKey: "\uf704", // NSF1FunctionKey
		},
		"F2": {
			key:     gui.KeyF2,
			wantKey: "\uf705", // NSF2FunctionKey
		},
		"F3": {
			key:     gui.KeyF3,
			wantKey: "\uf706", // NSF3FunctionKey
		},
		"F4": {
			key:     gui.KeyF4,
			wantKey: "\uf707", // NSF4FunctionKey
		},
		"F5": {
			key:     gui.KeyF5,
			wantKey: "\uf708", // NSF5FunctionKey
		},
		"F6": {
			key:     gui.KeyF6,
			wantKey: "\uf709", // NSF6FunctionKey
		},
		"F7": {
			key:     gui.KeyF7,
			wantKey: "\uf70a", // NSF7FunctionKey
		},
		"F8": {
			key:     gui.KeyF8,
			wantKey: "\uf70b", // NSF8FunctionKey
		},
		"F9": {
			key:     gui.KeyF9,
			wantKey: "\uf70c", // NSF9FunctionKey
		},
		"Home": {
			key:     gui.KeyHome,
			wantKey: "\uf729", // NSHomeFunctionKey
		},
		"Insert": {
			key:     gui.KeyInsert,
			wantKey: "\uf727", // NSInsertFunctionKey
		},
		"Left": {
			key:     gui.KeyLeft,
			wantKey: "\uf702", // NSLeftArrowFunctionKey
		},
		"PageDown": {
			key:     gui.KeyPageDown,
			wantKey: "\uf72d", // NSPageDownFunctionKey
		},
		"PageUp": {
			key:     gui.KeyPageUp,
			wantKey: "\uf72c", // NSPageUpFunctionKey
		},
		"Return": {
			key:     gui.KeyReturn,
			wantKey: "\n",
		},
		"Right": {
			key:     gui.KeyRight,
			wantKey: "\uf703", // NSRightArrowFunctionKey
		},
		"Space": {
			key:     gui.KeySpace,
			wantKey: " ",
		},
		"Tab": {
			key:     gui.KeyTab,
			wantKey: "\t",
		},
		"Up": {
			key:     gui.KeyUp,
			wantKey: "\uf700", // NSUpArrowFunctionKey
		},
	} {
		t.Run(name, func(t *testing.T) {
			resetMainMenu()
			w := createWindow("Test").(*window)
			item := gui.NewMenuItem("Special", func() {})
			item.Shortcut = &desktop.CustomShortcut{KeyName: tt.key, Modifier: gui.KeyModifierShortcutDefault}
			menu := gui.NewMenu("Special", item)
			mainMenu := gui.NewMainMenu(menu)
			setupNativeMenu(w, mainMenu)

			mm := testDarwinMainMenu()
			m := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
			assertNSMenuItem(t, "Special", tt.wantKey, 0b100000000000000000000, m, 0)
		})
	}
}

var initialAppMenuItems []string
var initialMenus []string

func assertNSMenuItem(t *testing.T, wantTitle, wantKey string, wantModifier uint64, m unsafe.Pointer, i int) {
	item := testNSMenuItemAtIndex(m, i)
	assert.Equal(t, wantTitle, testNSMenuItemTitle(item))
	if wantKey != "" {
		assert.Equal(t, wantKey, testNSMenuItemKeyEquivalent(item))
		assert.Equal(t, wantModifier, testNSMenuItemKeyEquivalentModifierMask(item))
	}
	testNSMenuPerformActionForItemAtIndex(m, i)
}

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
