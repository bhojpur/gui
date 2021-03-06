//go:build !no_native_menus
// +build !no_native_menus

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
	"strings"
	"unsafe"

	gui "github.com/bhojpur/gui/pkg/engine"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

#include <AppKit/AppKit.h>

// Using void* as type for pointers is a workaround.
void        assignDarwinSubmenu(const void*, const void*);
void        completeDarwinMenu(void* menu, bool prepend);
const void* createDarwinMenu(const char* label);
const void* darwinAppMenu();
const void* insertDarwinMenuItem(const void* menu, const char* label, const char* keyEquivalent, unsigned int keyEquivalentModifierMask, int id, int index, bool isSeparator);
void        resetDarwinMenu();

// Used for tests.
const void*   test_darwinMainMenu();
const void*   test_NSMenu_itemAtIndex(const void*, NSInteger);
NSInteger     test_NSMenu_numberOfItems(const void*);
void          test_NSMenu_performActionForItemAtIndex(const void*, NSInteger);
void          test_NSMenu_removeItemAtIndex(const void* m, NSInteger i);
const char*   test_NSMenu_title(const void*);
bool          test_NSMenuItem_isSeparatorItem(const void*);
const char*   test_NSMenuItem_keyEquivalent(const void*);
unsigned long test_NSMenuItem_keyEquivalentModifierMask(const void*);
const void*   test_NSMenuItem_submenu(const void*);
const char*   test_NSMenuItem_title(const void*);
*/
import "C"

type menuCallbacks struct {
	action  func()
	enabled func() bool
	checked func() bool
}

var callbacks []*menuCallbacks
var ecb func(string)
var specialKeys = map[gui.KeyName]string{
	gui.KeyBackspace: "\x08",
	gui.KeyDelete:    "\x7f",
	gui.KeyDown:      "\uf701",
	gui.KeyEnd:       "\uf72b",
	gui.KeyEnter:     "\x03",
	gui.KeyEscape:    "\x1b",
	gui.KeyF10:       "\uf70d",
	gui.KeyF11:       "\uf70e",
	gui.KeyF12:       "\uf70f",
	gui.KeyF1:        "\uf704",
	gui.KeyF2:        "\uf705",
	gui.KeyF3:        "\uf706",
	gui.KeyF4:        "\uf707",
	gui.KeyF5:        "\uf708",
	gui.KeyF6:        "\uf709",
	gui.KeyF7:        "\uf70a",
	gui.KeyF8:        "\uf70b",
	gui.KeyF9:        "\uf70c",
	gui.KeyHome:      "\uf729",
	gui.KeyInsert:    "\uf727",
	gui.KeyLeft:      "\uf702",
	gui.KeyPageDown:  "\uf72d",
	gui.KeyPageUp:    "\uf72c",
	gui.KeyReturn:    "\n",
	gui.KeyRight:     "\uf703",
	gui.KeySpace:     " ",
	gui.KeyTab:       "\t",
	gui.KeyUp:        "\uf700",
}

func addNativeMenu(w *window, menu *gui.Menu, nextItemID int, prepend bool) int {
	menu, nextItemID = handleSpecialItems(w, menu, nextItemID, true)

	containsItems := false
	for _, item := range menu.Items {
		if !item.IsSeparator {
			containsItems = true
			break
		}
	}
	if !containsItems {
		return nextItemID
	}

	nsMenu, nextItemID := createNativeMenu(w, menu, nextItemID)
	C.completeDarwinMenu(nsMenu, C.bool(prepend))
	return nextItemID
}

func addNativeSubmenu(w *window, nsParentMenuItem unsafe.Pointer, menu *gui.Menu, nextItemID int) int {
	nsMenu, nextItemID := createNativeMenu(w, menu, nextItemID)
	C.assignDarwinSubmenu(nsParentMenuItem, nsMenu)
	return nextItemID
}

func clearNativeMenu() {
	C.resetDarwinMenu()
}

func createNativeMenu(w *window, menu *gui.Menu, nextItemID int) (unsafe.Pointer, int) {
	nsMenu := C.createDarwinMenu(C.CString(menu.Label))
	for _, item := range menu.Items {
		nsMenuItem := C.insertDarwinMenuItem(
			nsMenu,
			C.CString(item.Label),
			C.CString(keyEquivalent(item)),
			C.uint(keyEquivalentModifierMask(item)),
			C.int(nextItemID),
			C.int(-1),
			C.bool(item.IsSeparator),
		)
		nextItemID = registerCallback(w, item, nextItemID)
		if item.ChildMenu != nil {
			nextItemID = addNativeSubmenu(w, nsMenuItem, item.ChildMenu, nextItemID)
		}
	}
	return nsMenu, nextItemID
}

//export exceptionCallback
func exceptionCallback(e *C.char) {
	msg := C.GoString(e)
	if ecb == nil {
		panic("unhandled Obj-C exception: " + msg)
	}
	ecb(msg)
}

func handleSpecialItems(w *window, menu *gui.Menu, nextItemID int, addSeparator bool) (*gui.Menu, int) {
	for i, item := range menu.Items {
		if item.Label == "Settings" || item.Label == "Settings???" || item.Label == "Preferences" || item.Label == "Preferences???" {
			items := make([]*gui.MenuItem, 0, len(menu.Items)-1)
			items = append(items, menu.Items[:i]...)
			items = append(items, menu.Items[i+1:]...)
			menu, nextItemID = handleSpecialItems(w, gui.NewMenu(menu.Label, items...), nextItemID, false)

			C.insertDarwinMenuItem(
				C.darwinAppMenu(),
				C.CString(item.Label),
				C.CString(keyEquivalent(item)),
				C.uint(keyEquivalentModifierMask(item)),
				C.int(nextItemID),
				C.int(1),
				C.bool(false),
			)
			if addSeparator {
				C.insertDarwinMenuItem(
					C.darwinAppMenu(),
					C.CString(""),
					C.CString(""),
					C.uint(0),
					C.int(nextItemID),
					C.int(1),
					C.bool(true),
				)
			}
			nextItemID = registerCallback(w, item, nextItemID)
			break
		}
	}
	return menu, nextItemID
}

func keyEquivalent(item *gui.MenuItem) (key string) {
	if s, ok := item.Shortcut.(gui.KeyboardShortcut); ok {
		if key = specialKeys[s.Key()]; key == "" {
			if len(s.Key()) > 1 {
				gui.LogError(fmt.Sprintf("unsupported key ???%s??? for menu shortcut", s.Key()), nil)
			}
			key = strings.ToLower(string(s.Key()))
		}
	}
	return
}

func keyEquivalentModifierMask(item *gui.MenuItem) (mask uint) {
	if s, ok := item.Shortcut.(gui.KeyboardShortcut); ok {
		if (s.Mod() & gui.KeyModifierShift) != 0 {
			mask |= 1 << 17 // NSEventModifierFlagShift
		}
		if (s.Mod() & gui.KeyModifierAlt) != 0 {
			mask |= 1 << 19 // NSEventModifierFlagOption
		}
		if (s.Mod() & gui.KeyModifierControl) != 0 {
			mask |= 1 << 18 // NSEventModifierFlagControl
		}
		if (s.Mod() & gui.KeyModifierSuper) != 0 {
			mask |= 1 << 20 // NSEventModifierFlagCommand
		}
	}
	return
}

func registerCallback(w *window, item *gui.MenuItem, nextItemID int) int {
	if !item.IsSeparator {
		callbacks = append(callbacks, &menuCallbacks{
			action: func() {
				if item.Action != nil {
					w.QueueEvent(item.Action)
				}
			},
			enabled: func() bool {
				return !item.Disabled
			},
			checked: func() bool {
				return item.Checked
			},
		})
		nextItemID++
	}
	return nextItemID
}

func setExceptionCallback(cb func(string)) {
	ecb = cb
}

func hasNativeMenu() bool {
	return true
}

//export menuCallback
func menuCallback(id int) {
	callbacks[id].action()
}

//export menuEnabled
func menuEnabled(id int) bool {
	return callbacks[id].enabled()
}

//export menuChecked
func menuChecked(id int) bool {
	return callbacks[id].checked()
}

func setupNativeMenu(w *window, main *gui.MainMenu) {
	clearNativeMenu()
	nextItemID := 0
	callbacks = []*menuCallbacks{}
	var helpMenu *gui.Menu
	for i := len(main.Items) - 1; i >= 0; i-- {
		menu := main.Items[i]
		if menu.Label == "Help" {
			helpMenu = menu
			continue
		}
		nextItemID = addNativeMenu(w, menu, nextItemID, true)
	}
	if helpMenu != nil {
		addNativeMenu(w, helpMenu, nextItemID, false)
	}
}

//
// Test support methods
// These are needed because CGo is not supported inside test files.
//

func testDarwinMainMenu() unsafe.Pointer {
	return C.test_darwinMainMenu()
}

func testNSMenuItemAtIndex(m unsafe.Pointer, i int) unsafe.Pointer {
	return C.test_NSMenu_itemAtIndex(m, C.long(i))
}

func testNSMenuNumberOfItems(m unsafe.Pointer) int {
	return int(C.test_NSMenu_numberOfItems(m))
}

func testNSMenuPerformActionForItemAtIndex(m unsafe.Pointer, i int) {
	C.test_NSMenu_performActionForItemAtIndex(m, C.long(i))
}

func testNSMenuRemoveItemAtIndex(m unsafe.Pointer, i int) {
	C.test_NSMenu_removeItemAtIndex(m, C.long(i))
}

func testNSMenuTitle(m unsafe.Pointer) string {
	return C.GoString(C.test_NSMenu_title(m))
}

func testNSMenuItemIsSeparatorItem(i unsafe.Pointer) bool {
	return bool(C.test_NSMenuItem_isSeparatorItem(i))
}

func testNSMenuItemKeyEquivalent(i unsafe.Pointer) string {
	return C.GoString(C.test_NSMenuItem_keyEquivalent(i))
}

func testNSMenuItemKeyEquivalentModifierMask(i unsafe.Pointer) uint64 {
	return uint64(C.ulong(C.test_NSMenuItem_keyEquivalentModifierMask(i)))
}

func testNSMenuItemSubmenu(i unsafe.Pointer) unsafe.Pointer {
	return C.test_NSMenuItem_submenu(i)
}

func testNSMenuItemTitle(i unsafe.Pointer) string {
	return C.GoString(C.test_NSMenuItem_title(i))
}
