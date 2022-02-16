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

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/test"
)

func TestMenu_ItemTapped(t *testing.T) {
	tapped := false
	item1 := gui.NewMenuItem("Foo", nil)
	item2 := gui.NewMenuItem("Bar", func() { tapped = true })
	item3 := gui.NewMenuItem("Sub", nil)
	subItem := gui.NewMenuItem("Foo", func() {})
	item3.ChildMenu = gui.NewMenu("", subItem)
	m := NewMenu(gui.NewMenu("", item1, item2, item3))
	size := m.MinSize()
	m.Resize(size)
	dismissed := false
	m.OnDismiss = func() { dismissed = true }

	mi1 := m.Items[0].(*menuItem)
	mi2 := m.Items[1].(*menuItem)
	mi3 := m.Items[2].(*menuItem)
	assert.Equal(t, item1, mi1.Item)
	assert.Equal(t, item2, mi2.Item)
	assert.Equal(t, item3, mi3.Item)

	// tap on item without action does not panic
	test.Tap(mi1)
	assert.False(t, tapped)
	assert.False(t, dismissed, "tap on item w/o action does not dismiss the menu")
	assert.True(t, m.Visible(), "tap on item w/o action does not hide the menu")

	test.Tap(mi2)
	assert.True(t, tapped)
	assert.True(t, dismissed, "tap on item dismisses the menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … OnDismiss is responsible for that")

	dismissed = false // reset
	mi3.MouseIn(nil)
	sm := mi3.child
	smi := sm.Items[0].(*menuItem)
	assert.Equal(t, subItem, smi.Item)
	assert.True(t, sm.Visible(), "submenu is visible")

	test.Tap(smi)
	assert.True(t, dismissed, "tap on child item dismisses the root menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … OnDismiss is responsible for that")
	assert.False(t, sm.Visible(), "tap on child item hides the submenu")

	newActionTapped := false
	item2.Action = func() { newActionTapped = true }
	test.Tap(mi2)
	assert.True(t, newActionTapped, "tap on item performs its current action")
}
