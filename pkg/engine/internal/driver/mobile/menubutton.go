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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type menuButton struct {
	widget.BaseWidget
	win  *window
	menu *gui.MainMenu
}

func (w *window) newMenuButton(menu *gui.MainMenu) *menuButton {
	b := &menuButton{win: w, menu: menu}
	b.ExtendBaseWidget(b)
	return b
}

func (m *menuButton) CreateRenderer() gui.WidgetRenderer {
	return &menuButtonRenderer{btn: widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		m.win.canvas.showMenu(m.menu)
	}), bg: canvas.NewRectangle(theme.BackgroundColor())}
}

type menuButtonRenderer struct {
	btn *widget.Button
	bg  *canvas.Rectangle
}

func (m *menuButtonRenderer) Destroy() {
}

func (m *menuButtonRenderer) Layout(size gui.Size) {
	m.bg.Move(gui.NewPos(theme.Padding()/2, theme.Padding()/2))
	m.bg.Resize(size.Subtract(gui.NewSize(theme.Padding(), theme.Padding())))
	m.btn.Resize(size)
}

func (m *menuButtonRenderer) MinSize() gui.Size {
	return m.btn.MinSize()
}

func (m *menuButtonRenderer) Objects() []gui.CanvasObject {
	return []gui.CanvasObject{m.bg, m.btn}
}

func (m *menuButtonRenderer) Refresh() {
	m.bg.FillColor = theme.BackgroundColor()
}
