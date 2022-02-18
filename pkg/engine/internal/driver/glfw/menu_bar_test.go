//go:build !mobile && (!ci || !windows)
// +build !mobile
// +build !ci !windows

package glfw_test

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
	"strconv"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/glfw"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenuBar(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var lastAction string

	m1i3i3i1 := gui.NewMenuItem("Old 1", func() { lastAction = "old 1" })
	m1i3i3i2 := gui.NewMenuItem("Old 2", func() { lastAction = "old 2" })
	m1i3i1 := gui.NewMenuItem("File 1", func() { lastAction = "file 1" })
	m1i3i2 := gui.NewMenuItem("File 2", func() { lastAction = "file 2" })
	m1i3i3 := gui.NewMenuItem("Older", nil)
	m1i3i3.ChildMenu = gui.NewMenu("", m1i3i3i1, m1i3i3i2)
	m1i1 := gui.NewMenuItem("New", func() { lastAction = "new" })
	m1i2 := gui.NewMenuItem("Open", func() { lastAction = "open" })
	m1i3 := gui.NewMenuItem("Recent", nil)
	m1i3.ChildMenu = gui.NewMenu("", m1i3i1, m1i3i2, m1i3i3)
	// TODO: remove useless separators: trailing, leading & double
	// m1 := gui.NewMenu("File", m1i1, m1i2, gui.newMenuItemSeparator(), m1i3)
	m1 := gui.NewMenu("File", m1i1, m1i2, m1i3)

	m2i1 := gui.NewMenuItem("Copy", func() { lastAction = "copy" })
	m2i2 := gui.NewMenuItem("Paste", func() { lastAction = "paste" })
	m2 := gui.NewMenu("Edit", m2i1, m2i2)

	m3i1 := gui.NewMenuItem("Help!", func() { lastAction = "help" })
	m3 := gui.NewMenu("Help", m3i1)

	menu := gui.NewMainMenu(m1, m2, m3)

	t.Run("mouse control and basic behaviour", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(gui.NewSize(300, 300))
		c := w.Canvas()

		menuBar := glfw.NewMenuBar(menu, c)
		themeCounter := 0
		button := newNotFocusableButton("Button", func() {
			switch themeCounter % 2 {
			case 0:
				test.ApplyTheme(t, test.NewTheme())
			case 1:
				test.ApplyTheme(t, test.Theme())
			}
			themeCounter++
		})
		container := container.NewWithoutLayout(button, menuBar)
		w.SetContent(container)
		w.Resize(gui.NewSize(300, 300))
		button.Resize(button.MinSize())
		button.Move(gui.NewPos(100, 50))
		menuBar.Resize(gui.NewSize(300, 0).Max(menuBar.MinSize()))

		buttonPos := gui.NewPos(110, 60)
		fileMenuPos := gui.NewPos(20, 10)
		fileNewPos := gui.NewPos(20, 50)
		fileOpenPos := gui.NewPos(20, 70)
		fileRecentPos := gui.NewPos(20, 100)
		fileRecentOlderPos := gui.NewPos(120, 170)
		fileRecentOlderOld1Pos := gui.NewPos(200, 170)
		editMenuPos := gui.NewPos(70, 10)
		helpMenuPos := gui.NewPos(120, 10)
		type action struct {
			typ string
			pos gui.Position
		}
		type step struct {
			actions    []action
			wantImage  string
			wantAction string
		}
		for name, tt := range map[string]struct {
			steps []step
		}{
			"switch theme": {
				steps: []step{
					{
						actions:   []action{{"move", buttonPos}},
						wantImage: "menu_bar_hovered_content.png",
					},
					{
						actions:   []action{{"tap", buttonPos}},
						wantImage: "menu_bar_hovered_content_test_theme.png",
					},
					{
						actions:   []action{{"tap", buttonPos}},
						wantImage: "menu_bar_hovered_content.png",
					},
				},
			},
			"activate and deactivate menu": {
				[]step{
					{
						actions:   []action{{"move", fileMenuPos}},
						wantImage: "menu_bar_inactive_file.png",
					},
					{
						actions:   []action{{"tap", fileMenuPos}},
						wantImage: "menu_bar_active_file.png",
					},
					{
						actions:   []action{{"tap", fileMenuPos}},
						wantImage: "menu_bar_inactive_file.png",
					},
					{
						actions:   []action{{"move", buttonPos}},
						wantImage: "menu_bar_hovered_content.png",
					},
				},
			},
			"active menu deactivates content": {
				[]step{
					{
						actions:   []action{{"move", fileMenuPos}},
						wantImage: "menu_bar_inactive_file.png",
					},
					{
						actions:   []action{{"tap", fileMenuPos}},
						wantImage: "menu_bar_active_file.png",
					},
					{
						actions:   []action{{"move", buttonPos}},
						wantImage: "menu_bar_content_not_hoverable_with_active_menu.png",
					},
					{
						// TODO: should be hovered over content here
						// -> is this canvas logic?
						// -> it would be for overlays (probably the same issue present)
						// -> it would not be for current menu implementation (menu is not an overlay, canvas does not know about activation state)
						actions:   []action{{"tap", buttonPos}},
						wantImage: "menu_bar_tap_content_with_active_menu_does_not_trigger_action_but_dismisses_menu.png",
					},
					{
						// menu bar is inactive again (menu not shown at hover)
						actions:   []action{{"move", fileMenuPos}},
						wantImage: "menu_bar_inactive_file.png",
					},
					{
						actions:   []action{{"move", buttonPos}},
						wantImage: "menu_bar_hovered_content.png",
					},
				},
			},
			"menu action File->New": {
				[]step{
					{
						actions:   []action{{"tap", fileMenuPos}},
						wantImage: "menu_bar_active_file.png",
					},
					{
						actions:   []action{{"move", fileNewPos}},
						wantImage: "menu_bar_hovered_file_new.png",
					},
					{
						actions:    []action{{"tap", fileNewPos}},
						wantAction: "new",
						wantImage:  "menu_bar_initial.png",
					},
				},
			},
			"menu action File->Open": {
				[]step{
					{
						actions: []action{
							{"tap", fileMenuPos},
							{"move", fileOpenPos},
						},
						wantImage: "menu_bar_hovered_file_open.png",
					},
					{
						actions:    []action{{"tap", fileOpenPos}},
						wantAction: "open",
						wantImage:  "menu_bar_initial.png",
					},
				},
			},
			"menu action File->Recent->Older->Old 1": {
				[]step{
					{
						actions: []action{
							{"tap", fileMenuPos},
							{"move", fileRecentPos},
						},
						wantImage: "menu_bar_hovered_file_recent.png",
					},
					{
						actions:   []action{{"move", fileRecentOlderPos}},
						wantImage: "menu_bar_hovered_file_recent_older.png",
					},
					{
						actions:   []action{{"move", fileRecentOlderOld1Pos}},
						wantImage: "menu_bar_hovered_file_recent_older_old1.png",
					},
					{
						actions:    []action{{"tap", fileRecentOlderOld1Pos}},
						wantAction: "old 1",
						wantImage:  "menu_bar_initial.png",
					},
				},
			},
			"move mouse outside does not hide menu": {
				[]step{
					{
						actions: []action{
							{"tap", fileMenuPos},
							{"move", fileRecentPos},
							{"move", fileRecentOlderPos},
							{"move", fileRecentOlderOld1Pos},
						},
						wantImage: "menu_bar_hovered_file_recent_older_old1.png",
					},
					{
						actions:   []action{{"move", buttonPos}},
						wantImage: "menu_bar_hovered_file_recent_older.png",
					},
				},
			},
			"hover other menu item hides previous submenus": {
				[]step{
					{
						actions: []action{
							{"tap", fileMenuPos},
							{"move", fileRecentPos},
							{"move", fileRecentOlderPos},
							{"move", fileRecentOlderOld1Pos},
						},
						wantImage: "menu_bar_hovered_file_recent_older_old1.png",
					},
					{
						actions:   []action{{"move", fileNewPos}},
						wantImage: "menu_bar_hovered_file_new.png",
					},
					{
						actions:   []action{{"move", fileRecentPos}},
						wantImage: "menu_bar_hovered_file_recent.png",
					},
				},
			},
			"hover other menu bar item changes active menu": {
				[]step{
					{
						actions:   []action{{"tap", fileMenuPos}},
						wantImage: "menu_bar_active_file.png",
					},
					{
						actions:   []action{{"move", editMenuPos}},
						wantImage: "menu_bar_active_edit.png",
					},
					{
						actions:   []action{{"move", helpMenuPos}},
						wantImage: "menu_bar_active_help.png",
					},
				},
			},
			"hover other menu bar item hides previous submenus": {
				[]step{
					{
						actions: []action{
							{"tap", fileMenuPos},
							{"move", fileRecentPos},
							{"move", fileRecentOlderPos},
							{"move", fileRecentOlderOld1Pos},
						},
						wantImage: "menu_bar_hovered_file_recent_older_old1.png",
					},
					{
						actions:   []action{{"move", helpMenuPos}},
						wantImage: "menu_bar_active_help.png",
					},
					{
						actions:   []action{{"move", fileMenuPos}},
						wantImage: "menu_bar_active_file.png",
					},
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				test.MoveMouse(c, gui.NewPos(0, 0))
				test.TapCanvas(c, gui.NewPos(0, 0))
				if test.AssertImageMatches(t, "menu_bar_initial.png", c.Capture()) {
					for i, s := range tt.steps {
						t.Run("step "+strconv.Itoa(i+1), func(t *testing.T) {
							lastAction = ""
							for _, a := range s.actions {
								switch a.typ {
								case "move":
									test.MoveMouse(c, a.pos)
								case "tap":
									test.MoveMouse(c, a.pos)
									test.TapCanvas(c, a.pos)
								}
							}
							test.AssertImageMatches(t, s.wantImage, c.Capture())
							assert.Equal(t, s.wantAction, lastAction, "last action should match expected")
						})
					}
				}
			})
		}
	})

	t.Run("keyboard control", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(gui.NewSize(300, 300))
		c := w.Canvas()

		menuBar := glfw.NewMenuBar(menu, c)
		themeCounter := 0
		button := newNotFocusableButton("Button", func() {
			switch themeCounter % 2 {
			case 0:
				test.ApplyTheme(t, test.NewTheme())
			case 1:
				test.ApplyTheme(t, test.Theme())
			}
			themeCounter++
		})
		container := container.NewWithoutLayout(button, menuBar)
		w.SetContent(container)
		w.Resize(gui.NewSize(300, 300))
		button.Resize(button.MinSize())
		button.Move(gui.NewPos(100, 50))
		menuBar.Resize(gui.NewSize(300, 0).Max(menuBar.MinSize()))

		fileMenuPos := gui.NewPos(20, 10)
		for name, tt := range map[string]struct {
			keys       []gui.KeyName
			wantAction string
		}{
			"traverse_menu_bar_items_right_1": {
				keys: []gui.KeyName{gui.KeyRight},
			},
			"traverse_menu_bar_items_right_2": {
				keys: []gui.KeyName{gui.KeyRight, gui.KeyRight},
			},
			"traverse_menu_bar_items_right_3": {
				keys: []gui.KeyName{gui.KeyRight, gui.KeyRight, gui.KeyRight},
			},
			"traverse_menu_bar_items_left_1": {
				keys: []gui.KeyName{gui.KeyLeft},
			},
			"traverse_menu_bar_items_left_2": {
				keys: []gui.KeyName{gui.KeyLeft, gui.KeyLeft},
			},
			"traverse_menu_bar_items_left_3": {
				keys: []gui.KeyName{gui.KeyLeft, gui.KeyLeft, gui.KeyLeft},
			},
			"traverse_menu_down_1": {
				keys: []gui.KeyName{gui.KeyDown},
			},
			"traverse_menu_down_2": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown},
			},
			"traverse_menu_down_3": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown},
			},
			"traverse_menu_up_1": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyUp},
			},
			"traverse_menu_up_2": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyUp, gui.KeyUp},
			},
			"open_submenu_1": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyRight},
			},
			"open_submenu_2": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyDown, gui.KeyDown, gui.KeyRight},
			},
			"close_submenu_1": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyLeft},
			},
			"close_submenu_2": {
				keys: []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyLeft, gui.KeyLeft},
			},
			"trigger_with_enter": {
				keys:       []gui.KeyName{gui.KeyDown, gui.KeyEnter},
				wantAction: "new",
			},
			"trigger_with_return": {
				keys:       []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyReturn},
				wantAction: "open",
			},
			"trigger_with_space": {
				keys:       []gui.KeyName{gui.KeyRight, gui.KeyDown, gui.KeySpace},
				wantAction: "copy",
			},
			"trigger_submenu_item": {
				keys:       []gui.KeyName{gui.KeyDown, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyDown, gui.KeyDown, gui.KeyRight, gui.KeyReturn},
				wantAction: "old 1",
			},
			"trigger_without_active_item": {
				keys:       []gui.KeyName{gui.KeyEnter},
				wantAction: "",
			},
		} {
			t.Run(name, func(t *testing.T) {
				test.MoveMouse(c, gui.NewPos(0, 0))
				test.TapCanvas(c, gui.NewPos(0, 0))
				test.TapCanvas(c, fileMenuPos) // activate menu
				require.Equal(t, menuBar.Items[0], c.Focused())
				if test.AssertImageMatches(t, "menu_bar_active_file.png", c.Capture()) {
					lastAction = ""
					for _, key := range tt.keys {
						c.Focused().TypedKey(&gui.KeyEvent{
							Name: key,
						})
					}
					test.AssertRendersToMarkup(t, "menu_bar_kbdctrl_"+name+".xml", c)
					assert.Equal(t, tt.wantAction, lastAction, "last action should match expected")
				}
			})
		}

		t.Run("moving mouse over unfocused item moves focus", func(t *testing.T) {
			test.MoveMouse(c, gui.NewPos(0, 0))
			test.TapCanvas(c, gui.NewPos(0, 0))
			test.MoveMouse(c, fileMenuPos)
			test.TapCanvas(c, fileMenuPos) // activate menu
			require.Equal(t, menuBar.Items[0], c.Focused())
			c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
			c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
			require.Equal(t, menuBar.Items[2], c.Focused())

			test.MoveMouse(c, fileMenuPos.Add(gui.NewPos(1, 0)))
			assert.Equal(t, menuBar.Items[0], c.Focused())
			test.AssertImageMatches(t, "menu_bar_active_file.png", c.Capture())
		})
	})
}

func TestMenuBar_Toggle(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	m1i1 := gui.NewMenuItem("New", nil)
	m1i2 := gui.NewMenuItem("Open", nil)
	m1 := gui.NewMenu("File", m1i1, m1i2)

	m2i1 := gui.NewMenuItem("Copy", nil)
	m2i2 := gui.NewMenuItem("Paste", nil)
	m2 := gui.NewMenu("Edit", m2i1, m2i2)

	menu := gui.NewMainMenu(m1, m2)

	t.Run("when menu bar is inactive", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(gui.NewSize(300, 300))
		c := w.Canvas()
		menuBar := glfw.NewMenuBar(menu, c)
		w.SetContent(container.NewWithoutLayout(menuBar))
		w.Resize(gui.NewSize(300, 300))
		menuBar.Resize(gui.NewSize(300, 0).Max(menuBar.MinSize()))

		require.False(t, menuBar.IsActive())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_deactivated.xml", c)

		menuBar.Toggle()
		assert.True(t, menuBar.IsActive())
		assert.Equal(t, c.Focused(), menuBar.Items[0])
		test.AssertRendersToMarkup(t, "menu_bar_toggle_first_item_active.xml", c)
	})

	t.Run("when menu bar is active (first menu item active)", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(gui.NewSize(300, 300))
		c := w.Canvas()
		menuBar := glfw.NewMenuBar(menu, c)
		w.SetContent(container.NewWithoutLayout(menuBar))
		w.Resize(gui.NewSize(300, 300))
		menuBar.Resize(gui.NewSize(300, 0).Max(menuBar.MinSize()))

		menuBar.Toggle()
		require.True(t, menuBar.IsActive())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_first_item_active.xml", c)

		menuBar.Toggle()
		assert.False(t, menuBar.IsActive())
		assert.Nil(t, c.Focused())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_deactivated.xml", c)
	})

	t.Run("when menu bar is active (second menu item active)", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(gui.NewSize(300, 300))
		c := w.Canvas()
		menuBar := glfw.NewMenuBar(menu, c)
		w.SetContent(container.NewWithoutLayout(menuBar))
		w.Resize(gui.NewSize(300, 300))
		menuBar.Resize(gui.NewSize(300, 0).Max(menuBar.MinSize()))

		menuBar.Toggle()
		c.(test.WindowlessCanvas).FocusNext()
		require.True(t, menuBar.IsActive())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_second_item_active.xml", c)

		menuBar.Toggle()
		assert.False(t, menuBar.IsActive())
		assert.Nil(t, c.Focused())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_deactivated.xml", c)
	})
}

type notFocusableButton struct {
	widget.Label
	f func()
}

func newNotFocusableButton(l string, f func()) *notFocusableButton {
	n := &notFocusableButton{f: f}
	n.ExtendBaseWidget(n)
	n.Label.Text = l
	return n
}

func (n *notFocusableButton) Tapped(e *gui.PointEvent) {
	n.f()
}
