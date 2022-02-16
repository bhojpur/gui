//go:build mobile
// +build mobile

package widget_test

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
	"image/color"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	internalWidget "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestMenu_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	item1 := gui.NewMenuItem("A", nil)
	item2 := gui.NewMenuItem("B (long)", nil)
	sep := gui.NewMenuItemSeparator()
	item3 := gui.NewMenuItem("C", nil)
	subItem1 := gui.NewMenuItem("subitem A", nil)
	subItem2 := gui.NewMenuItem("subitem B", nil)
	subItem3 := gui.NewMenuItem("subitem C (long)", nil)
	subsubItem1 := gui.NewMenuItem("subsubitem A (long)", nil)
	subsubItem2 := gui.NewMenuItem("subsubitem B", nil)
	subItem3.ChildMenu = gui.NewMenu("", subsubItem1, subsubItem2)
	item3.ChildMenu = gui.NewMenu("", subItem1, subItem2, subItem3)
	menu := gui.NewMenu("", item1, sep, item2, item3)

	for name, tt := range map[string]struct {
		windowSize   gui.Size
		menuPos      gui.Position
		tapPositions []gui.Position
		useTestTheme bool
		want         string
	}{
		"normal": {
			windowSize: gui.NewSize(500, 300),
			menuPos:    gui.NewPos(10, 10),
			want:       "menu/mobile/layout_normal.xml",
		},
		"normal with submenus": {
			windowSize: gui.NewSize(500, 300),
			menuPos:    gui.NewPos(10, 10),
			tapPositions: []gui.Position{
				gui.NewPos(30, 100),
				gui.NewPos(100, 170),
			},
			want: "menu/mobile/layout_normal_with_submenus.xml",
		},
		"background of active submenu parents resets if sibling is hovered": {
			windowSize: gui.NewSize(500, 300),
			menuPos:    gui.NewPos(10, 10),
			tapPositions: []gui.Position{
				gui.NewPos(30, 100),  // open submenu
				gui.NewPos(100, 170), // open subsubmenu
				gui.NewPos(300, 170), // hover subsubmenu item
				gui.NewPos(30, 60),   // hover sibling of submenu parent
			},
			want: "menu/mobile/layout_background_reset.xml",
		},
		"no space on right side for submenu": {
			windowSize: gui.NewSize(500, 300),
			menuPos:    gui.NewPos(410, 10),
			tapPositions: []gui.Position{
				gui.NewPos(430, 100), // open submenu
				gui.NewPos(300, 170), // open subsubmenu
			},
			want: "menu/mobile/layout_no_space_on_right.xml",
		},
		"no space on left & right side for submenu": {
			windowSize: gui.NewSize(200, 300),
			menuPos:    gui.NewPos(10, 10),
			tapPositions: []gui.Position{
				gui.NewPos(30, 100),  // open submenu
				gui.NewPos(100, 170), // open subsubmenu
			},
			want: "menu/mobile/layout_no_space_on_both_sides.xml",
		},
		"window too short for submenu": {
			windowSize: gui.NewSize(500, 150),
			menuPos:    gui.NewPos(10, 10),
			tapPositions: []gui.Position{
				gui.NewPos(30, 100),  // open submenu
				gui.NewPos(100, 130), // open subsubmenu
			},
			want: "menu/mobile/layout_window_too_short_for_submenu.xml",
		},
		"theme change": {
			windowSize:   gui.NewSize(500, 300),
			menuPos:      gui.NewPos(10, 10),
			useTestTheme: true,
			want:         "menu/mobile/layout_theme_changed.xml",
		},
		"window too short for menu": {
			windowSize: gui.NewSize(100, 50),
			menuPos:    gui.NewPos(10, 10),
			want:       "menu/mobile/layout_window_too_short.xml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			w.Resize(tt.windowSize)
			m := widget.NewMenu(menu)
			o := internalWidget.NewOverlayContainer(m, c, nil)
			c.Overlays().Add(o)
			defer c.Overlays().Remove(o)
			m.Move(tt.menuPos)
			m.Resize(m.MinSize())
			for _, pos := range tt.tapPositions {
				test.TapCanvas(c, pos)
			}
			test.AssertRendersToMarkup(t, tt.want, w.Canvas())
			if tt.useTestTheme {
				test.AssertImageMatches(t, "menu/layout_normal.png", c.Capture())
				test.WithTestTheme(t, func() {
					test.AssertImageMatches(t, "menu/layout_theme_changed.png", c.Capture())
				})
			}
		})
	}
}

func TestMenu_Dragging(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	menu := gui.NewMenu("",
		gui.NewMenuItem("A", nil),
		gui.NewMenuItem("B", nil),
		gui.NewMenuItem("C", nil),
		gui.NewMenuItem("D", nil),
		gui.NewMenuItem("E", nil),
		gui.NewMenuItem("F", nil),
	)

	w.Resize(gui.NewSize(100, 100))
	m := widget.NewMenu(menu)
	o := internalWidget.NewOverlayContainer(m, c, nil)
	c.Overlays().Add(o)
	defer c.Overlays().Remove(o)
	m.Move(gui.NewPos(10, 10))
	m.Resize(m.MinSize())
	maxDragDistance := m.MinSize().Height - 90
	test.AssertRendersToMarkup(t, "menu/mobile/drag_top.xml", w.Canvas())

	test.Drag(c, gui.NewPos(20, 20), 0, -50)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_middle.xml", w.Canvas())

	test.Drag(c, gui.NewPos(20, 20), 0, -maxDragDistance)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_bottom.xml", w.Canvas())

	test.Drag(c, gui.NewPos(20, 20), 0, maxDragDistance-50)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_middle.xml", w.Canvas())

	test.Drag(c, gui.NewPos(20, 20), 0, 50)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_top.xml", w.Canvas())
}
