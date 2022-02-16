//go:build mobile
// +build mobile

package widget

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
	"github.com/bhojpur/gui/pkg/engine/test"

	"github.com/stretchr/testify/assert"
)

func TestMenu_ItemWithChildTapped(t *testing.T) {
	sub1 := gui.NewMenuItem("sub1", nil)
	sub1.ChildMenu = gui.NewMenu("",
		gui.NewMenuItem("sub1 A", nil),
		gui.NewMenuItem("sub1 B", nil),
	)
	sub2sub := gui.NewMenuItem("sub2sub", nil)
	sub2sub.ChildMenu = gui.NewMenu("",
		gui.NewMenuItem("sub2sub A", nil),
		gui.NewMenuItem("sub2sub B", nil),
	)
	sub2 := gui.NewMenuItem("sub2", nil)
	sub2.ChildMenu = gui.NewMenu("",
		gui.NewMenuItem("sub2 A", nil),
		gui.NewMenuItem("sub2 B", nil),
		sub2sub,
	)
	m := NewMenu(
		gui.NewMenu("",
			gui.NewMenuItem("Foo", nil),
			gui.NewMenuItemSeparator(),
			gui.NewMenuItem("Bar", nil),
			sub1,
			sub2,
		),
	)
	size := m.MinSize()
	m.Resize(size)

	sub1Widget := m.Items[3].(*menuItem)
	assert.Equal(t, sub1, sub1Widget.Item)
	sub2Widget := m.Items[4].(*menuItem)
	assert.Equal(t, sub2, sub2Widget.Item)
	assert.False(t, sub1Widget.child.Visible(), "submenu is invisible initially")
	assert.False(t, sub2Widget.child.Visible(), "submenu is invisible initially")
	test.Tap(sub1Widget)
	assert.True(t, sub1Widget.child.Visible(), "tapping item shows submenu")
	assert.False(t, sub2Widget.child.Visible(), "other child menu stays hidden")
	test.Tap(sub2Widget)
	assert.False(t, sub1Widget.child.Visible(), "tapping other item hides current submenu")
	assert.True(t, sub2Widget.child.Visible(), "other child menu shows up")

	sub2subWidget := sub2Widget.child.Items[2].(*menuItem)
	assert.Equal(t, sub2sub, sub2subWidget.Item)
	assert.False(t, sub2subWidget.child.Visible(), "2nd level submenu is invisible initially")
	test.Tap(sub2subWidget)
	assert.True(t, sub2Widget.child.Visible(), "1st level submenu stays visible")
	assert.True(t, sub2subWidget.child.Visible(), "2nd level submenu shows up")

	test.Tap(sub1Widget)
	assert.False(t, sub2Widget.child.Visible(), "1st level submenu is hidden by other submenu")
	test.Tap(sub2Widget)
	assert.False(t, sub2subWidget.child.Visible(), "2nd level submenu is hidden when re-entering its parent")
}
