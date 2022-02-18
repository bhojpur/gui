//go:build !ci && !mobile && (!darwin || no_native_menus)
// +build !ci
// +build !mobile
// +build !darwin no_native_menus

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"

	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_FocusHandlingWhenActivatingOrDeactivatingTheMenu(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	w.SetMainMenu(
		gui.NewMainMenu(
			gui.NewMenu("test", gui.NewMenuItem("item", func() {})),
			gui.NewMenu("other", gui.NewMenuItem("item", func() {})),
		),
	)
	c := w.Canvas().(*glCanvas)

	ce1 := &focusable{id: "ce1"}
	ce2 := &focusable{id: "ce2"}
	content := container.NewVBox(ce1, ce2)
	w.SetContent(content)

	assert.Nil(t, c.Focused())
	m := c.menu.(*MenuBar)
	assert.False(t, m.IsActive())

	c.FocusPrevious()
	assert.Equal(t, ce2, c.Focused())
	assert.True(t, ce2.focused)

	m.Items[0].(*menuBarItem).Tapped(&gui.PointEvent{})
	assert.True(t, m.IsActive())
	ctxt := "activating the menu changes focus handler and focuses the menu bar item but does not remove focus from content"
	assert.True(t, ce2.focused, ctxt)
	assert.Equal(t, m.Items[0], c.Focused(), ctxt)

	c.FocusNext()
	ctxt = "changing focus with active menu does not affect content focus"
	assert.True(t, ce2.focused, ctxt)
	assert.Equal(t, m.Items[1], c.Focused(), ctxt)

	m.Items[0].(*menuBarItem).Tapped(&gui.PointEvent{})
	assert.False(t, m.IsActive())
	ctxt = "deactivating the menu restores focus handler from content"
	assert.True(t, ce2.focused, ctxt)
	assert.Equal(t, ce2, c.Focused(), ctxt)

	c.FocusPrevious()
	assert.Equal(t, ce1, c.Focused(), ctxt)
	assert.True(t, ce1.focused, ctxt)
	assert.False(t, ce2.focused, ctxt)
}