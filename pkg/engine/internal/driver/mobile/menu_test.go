//go:build !windows || !ci
// +build !windows !ci

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	internalWidget "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func TestMobileCanvas_DismissBar(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(canvas.NewRectangle(theme.BackgroundColor()))
	menu := gui.NewMainMenu(
		gui.NewMenu("Test"))
	c.showMenu(menu)
	c.Resize(gui.NewSize(100, 100))

	assert.NotNil(t, c.menu)
	// simulate tap as the test util does not know about our menu...
	c.tapDown(gui.NewPos(80, 20), 1)
	c.tapUp(gui.NewPos(80, 20), 1, nil, nil, nil, nil)
	assert.Nil(t, c.menu)
}

func TestMobileCanvas_DismissMenu(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(canvas.NewRectangle(theme.BackgroundColor()))
	menu := gui.NewMainMenu(
		gui.NewMenu("Test", gui.NewMenuItem("TapMe", func() {})))
	c.showMenu(menu)
	c.Resize(gui.NewSize(100, 100))

	assert.NotNil(t, c.menu)
	menuObj := c.menu.(*gui.Container).Objects[1].(*gui.Container).Objects[1].(*menuLabel)
	point := &gui.PointEvent{Position: gui.NewPos(10, 10)}
	menuObj.Tapped(point)

	tapMeItem := c.Overlays().Top().(*internalWidget.OverlayContainer).Content.(*widget.PopUpMenu).Items[0].(gui.Tappable)
	tapMeItem.Tapped(point)
	assert.Nil(t, c.menu)
}

func TestMobileCanvas_Menu(t *testing.T) {
	c := &mobileCanvas{}
	labels := []string{"File", "Edit"}
	menu := gui.NewMainMenu(
		gui.NewMenu(labels[0]),
		gui.NewMenu(labels[1]))

	c.showMenu(menu)
	menuObjects := c.menu.(*gui.Container).Objects[1].(*gui.Container)
	assert.Equal(t, 3, len(menuObjects.Objects))
	header, ok := menuObjects.Objects[0].(*gui.Container)
	assert.True(t, ok)
	closed, ok := header.Objects[0].(*widget.Button)
	assert.True(t, ok)
	assert.Equal(t, theme.CancelIcon(), closed.Icon)

	for i := 1; i < 3; i++ {
		item, ok := menuObjects.Objects[i].(*menuLabel)
		assert.True(t, ok)
		assert.Equal(t, labels[i-1], item.menu.Label)
	}
}

func dummyWin(d *mobileDriver, title string) *window {
	ret := &window{title: title}
	d.windows = append(d.windows, ret)

	return ret
}

func TestMobileDriver_FindMenu(t *testing.T) {
	m1 := gui.NewMainMenu(gui.NewMenu("1"))
	m2 := gui.NewMainMenu(gui.NewMenu("2"))

	d := NewGoMobileDriver().(*mobileDriver)
	w1 := dummyWin(d, "top")
	w1.SetMainMenu(m1)
	assert.Equal(t, m1, d.findMenu(w1))

	w2 := dummyWin(d, "child")
	assert.Equal(t, m1, d.findMenu(w2))

	w2.SetMainMenu(m2)
	assert.Equal(t, m2, d.findMenu(w2))
}
