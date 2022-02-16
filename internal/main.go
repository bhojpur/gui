package main

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
	"log"
	"net/url"

	"github.com/bhojpur/gui/internal/tutorials"

	// demotheme "github.com/bhojpur/gui/internal/theme"
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/app"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow gui.Window

func main() {
	// initilize the Bhojpur GUI application
	demo := app.NewWithID("net.bhojpur.gui.demo")
	demo.SetIcon(theme.BhojpurLogo())
	wm := demo.NewWindow("Bhojpur Demo Application")
	topWindow = wm

	wm.SetMainMenu(makeMenu(demo, wm))
	wm.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = gui.TextWrapWord
	setTutorial := func(t tutorials.Tutorial) {
		if gui.CurrentDevice().IsMobile() {
			child := demo.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = wm
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []gui.CanvasObject{t.View(wm)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if gui.CurrentDevice().IsMobile() {
		wm.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		wm.SetContent(split)
	}
	wm.Resize(gui.NewSize(640, 460))
	wm.ShowAndRun()

	// wm.SetContent(c)
	// wm.Show()
	// wm.Resize(gui.NewSize(400, 300))
	// demo.Settings().SetTheme(demotheme.MyTheme{})
	// demo.Run()
}

func logLifecycle(a gui.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeMenu(a gui.App, w gui.Window) *gui.MainMenu {
	newItem := gui.NewMenuItem("New", nil)
	checkedItem := gui.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := gui.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := gui.NewMenuItem("Other", nil)
	otherItem.ChildMenu = gui.NewMenu("",
		gui.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		gui.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") }),
	)
	newItem.ChildMenu = gui.NewMenu("",
		gui.NewMenuItem("File", func() { fmt.Println("Menu New->File") }),
		gui.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") }),
		otherItem,
	)
	settingsItem := gui.NewMenuItem("Settings", func() {
		w := a.NewWindow("gui Settings")
		//w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(gui.NewSize(480, 480))
		w.Show()
	})

	cutItem := gui.NewMenuItem("Cut", func() {
		shortcutFocused(&gui.ShortcutCut{
			Clipboard: w.Clipboard(),
		}, w)
	})
	copyItem := gui.NewMenuItem("Copy", func() {
		shortcutFocused(&gui.ShortcutCopy{
			Clipboard: w.Clipboard(),
		}, w)
	})
	pasteItem := gui.NewMenuItem("Paste", func() {
		shortcutFocused(&gui.ShortcutPaste{
			Clipboard: w.Clipboard(),
		}, w)
	})
	findItem := gui.NewMenuItem("Find", func() { fmt.Println("Menu Find") })

	helpMenu := gui.NewMenu("Help",
		gui.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://docs.bhojpur.net")
			_ = a.OpenURL(u)
		}),
		gui.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://desk.bhojpur-consulting.com/")
			_ = a.OpenURL(u)
		}),
		gui.NewMenuItemSeparator(),
		gui.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://www.bhojpur-consulting.com")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := gui.NewMenu("File", newItem, checkedItem, disabledItem)
	if !gui.CurrentDevice().IsMobile() {
		file.Items = append(file.Items, gui.NewMenuItemSeparator(), settingsItem)
	}
	return gui.NewMainMenu(
		file,
		gui.NewMenu("Edit", cutItem, copyItem, pasteItem, gui.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
}

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) gui.CanvasObject {
	a := gui.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) gui.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj gui.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				gui.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := tutorials.Tutorials[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := gui.NewContainerWithLayout(layout.NewGridLayout(2),
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s gui.Shortcut, w gui.Window) {
	if focused, ok := w.Canvas().Focused().(gui.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func toolbar() *widget.Toolbar {
	var r gui.Resource
	r, _ = gui.LoadResourceFromPath("image/add.png")
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(r, func() {
			fmt.Println("create")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(r, func() {
			fmt.Println("create")
		}),
	)

	return toolbar
}

func menuMain1() *gui.MainMenu {
	menuFileNewFile := gui.NewMenuItem("New File", nil)
	menuFileNewWindow := gui.NewMenuItem("New Window", nil)
	menuFileOpen := gui.NewMenuItem("Open...", nil)
	menuFileOpenFolder := gui.NewMenuItem("Open Folder...", nil)
	menuFileSave := gui.NewMenuItem("Save", nil)
	menuFileSaveAs := gui.NewMenuItem("Save As", nil)
	menuFileSaveAll := gui.NewMenuItem("Save All", nil)
	menuFileClose := gui.NewMenuItem("Close", nil)
	menuFileCloseAll := gui.NewMenuItem("Close All", nil)
	menuFileExit := gui.NewMenuItem("Exit", nil)
	menuFileExit.IsQuit = true
	menuFile := gui.NewMenu("File",
		menuFileNewFile,
		menuFileNewWindow,
		gui.NewMenuItemSeparator(),
		menuFileOpen,
		menuFileOpenFolder,
		gui.NewMenuItemSeparator(),
		menuFileSave,
		menuFileSaveAs,
		menuFileSaveAll,
		gui.NewMenuItemSeparator(),
		menuFileClose,
		menuFileCloseAll,
		gui.NewMenuItemSeparator(),
		menuFileExit,
	)

	menuEditUndo := gui.NewMenuItem("Undo", nil)
	menuEditRedo := gui.NewMenuItem("Redo", nil)
	menuEdit := gui.NewMenu("Edit",
		menuEditUndo,
		menuEditRedo,
	)

	menuSelectSelectAll := gui.NewMenuItem("Select All", nil)
	menuSelect := gui.NewMenu("Select",
		menuSelectSelectAll,
	)

	menuViewCommandPallette := gui.NewMenuItem("Command Pallette", nil)
	menuView := gui.NewMenu("View",
		menuViewCommandPallette,
	)

	menuHelpGetStarted := gui.NewMenuItem("Get Started", nil)
	menuHelp := gui.NewMenu("Help",
		menuHelpGetStarted,
	)

	menu := gui.NewMainMenu(
		menuFile,
		menuEdit,
		menuSelect,
		menuView,
		menuHelp,
	)

	return menu
}

func menu2() *gui.MainMenu {
	// new menu items
	//first parameter is label, 2nd is function
	item1 := gui.NewMenuItem("edit", nil)
	item2 := gui.NewMenuItem("details", nil)
	item3 := gui.NewMenuItem("home", nil)
	item4 := gui.NewMenuItem("run", nil)
	// child menu
	item1.ChildMenu = gui.NewMenu(
		"",                           // leave label blank
		gui.NewMenuItem("copy", nil), // child menu items
		gui.NewMenuItem("cut", nil),
		gui.NewMenuItem("paste", nil),
	)
	// create child menu for 2nd item
	item2.ChildMenu = gui.NewMenu(
		"",                            // leave label blank
		gui.NewMenuItem("books", nil), // child menu items
		gui.NewMenuItem("magzine", nil),
		gui.NewMenuItem("notebook", nil),
	)
	// create child menu for third item
	item3.ChildMenu = gui.NewMenu(
		"",                             // leave label blank
		gui.NewMenuItem("school", nil), // child menu items
		gui.NewMenuItem("college", nil),
		gui.NewMenuItem("university", nil),
	)
	NewMenu1 := gui.NewMenu("File", item1, item2, item3, item4)
	NewMenu2 := gui.NewMenu("Help", item1, item2, item3, item4)
	// main menu
	menu := gui.NewMainMenu(NewMenu1, NewMenu2)
	// setup menu
	return menu
}
