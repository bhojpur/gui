package mobile

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type menuLabel struct {
	widget.BaseWidget

	menu   *gui.Menu
	bar    *gui.Container
	canvas *mobileCanvas
}

func (m *menuLabel) Tapped(*gui.PointEvent) {
	pos := gui.CurrentApp().Driver().AbsolutePositionForObject(m)
	menu := widget.NewPopUpMenu(m.menu, m.canvas)
	menu.ShowAtPosition(gui.NewPos(pos.X+m.Size().Width, pos.Y))

	menuDismiss := menu.OnDismiss // this dismisses the menu stack
	menu.OnDismiss = func() {
		menuDismiss()
		m.bar.Hide() // dismiss the overlay menu bar
		m.canvas.setMenu(nil)
	}
}

func (m *menuLabel) CreateRenderer() gui.WidgetRenderer {
	label := widget.NewLabel(m.menu.Label)
	box := container.NewHBox(layout.NewSpacer(), label, layout.NewSpacer(), widget.NewIcon(theme.MenuExpandIcon()))

	return &menuLabelRenderer{menu: m, content: box}
}

func newMenuLabel(item *gui.Menu, parent *gui.Container, c *mobileCanvas) *menuLabel {
	l := &menuLabel{menu: item, bar: parent, canvas: c}
	l.ExtendBaseWidget(l)
	return l
}

func (c *mobileCanvas) showMenu(menu *gui.MainMenu) {
	var panel *gui.Container
	top := container.NewHBox(widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		panel.Hide()
		c.setMenu(nil)
	}))
	panel = container.NewVBox(top)
	for _, item := range menu.Items {
		panel.Add(newMenuLabel(item, panel, c))
	}

	bg := canvas.NewRectangle(theme.BackgroundColor())
	shadow := canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)

	safePos, safeSize := c.InteractiveArea()
	bg.Move(safePos)
	bg.Resize(gui.NewSize(panel.MinSize().Width+theme.Padding(), safeSize.Height))
	panel.Move(safePos)
	panel.Resize(gui.NewSize(panel.MinSize().Width+theme.Padding(), safeSize.Height))
	shadow.Resize(gui.NewSize(theme.Padding()/2, safeSize.Height))
	shadow.Move(gui.NewPos(panel.Size().Width+safePos.X, safePos.Y))

	c.setMenu(container.NewWithoutLayout(bg, panel, shadow))
}

func (d *mobileDriver) findMenu(win *window) *gui.MainMenu {
	if win.menu != nil {
		return win.menu
	}

	matched := false
	for x := len(d.windows) - 1; x >= 0; x-- {
		w := d.windows[x]
		if !matched {
			if w == win {
				matched = true
			}
			continue
		}

		if w.(*window).menu != nil {
			return w.(*window).menu
		}
	}

	return nil
}

type menuLabelRenderer struct {
	menu    *menuLabel
	content *gui.Container
}

func (m *menuLabelRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (m *menuLabelRenderer) Destroy() {
}

func (m *menuLabelRenderer) Layout(size gui.Size) {
	m.content.Resize(size)
}

func (m *menuLabelRenderer) MinSize() gui.Size {
	return m.content.MinSize()
}

func (m *menuLabelRenderer) Objects() []gui.CanvasObject {
	return []gui.CanvasObject{m.content}
}

func (m *menuLabelRenderer) Refresh() {
	m.content.Refresh()
}
