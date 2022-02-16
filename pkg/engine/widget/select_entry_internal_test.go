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
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestSelectEntry_Disableable(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	options := []string{"A", "B", "C"}
	e := NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(gui.NewSize(150, 200))
	e.Resize(e.MinSize().Max(gui.NewSize(130, 0)))
	e.Move(gui.NewPos(10, 10))
	c := w.Canvas()

	assert.False(t, e.Disabled())
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled.xml", c)

	switchPos := gui.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_opened.xml", c)

	test.TapCanvas(c, gui.NewPos(0, 0))
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_tapped_selected.xml", c)

	e.Disable()
	assert.True(t, e.Disabled())
	test.AssertRendersToMarkup(t, "select_entry/disableable_disabled.xml", c)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_disabled.xml", c, "no drop-down when disabled")

	e.Enable()
	assert.False(t, e.Disabled())
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_tapped.xml", c)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_opened.xml", c)
}

func TestSelectEntry_DropDown(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	options := []string{"A", "B", "C"}
	e := NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(gui.NewSize(150, 200))
	e.Resize(e.MinSize().Max(gui.NewSize(130, 0)))
	e.Move(gui.NewPos(10, 10))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "select_entry/dropdown_initial.xml", c)
	assert.Nil(t, c.Overlays().Top())

	switchPos := gui.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)

	test.TapCanvas(c, gui.NewPos(50, 15+2*(theme.Padding()+e.Size().Height)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_tapped_B.xml", c)
	assert.Equal(t, "B", e.Text)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_B_opened.xml", c)

	test.TapCanvas(c, gui.NewPos(50, 15+3*(theme.Padding()+e.Size().Height)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_tapped_C.xml", c)
	assert.Equal(t, "C", e.Text)
}

func TestSelectEntry_DropDownMove(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	e := NewSelectEntry([]string{"one"})
	w := test.NewWindow(e)
	defer w.Close()
	entrySize := e.MinSize()
	w.Resize(entrySize.Add(gui.NewSize(100, 100)))
	e.Resize(entrySize)

	// open the popup
	test.Tap(e.ActionItem.(gui.Tappable))

	// first movement
	e.Move(gui.NewPos(10, 10))
	assert.Equal(t, gui.NewPos(10, 10), e.Entry.Position())
	assert.Equal(t,
		gui.NewPos(10, 10+entrySize.Height-theme.InputBorderSize()),
		e.popUp.Position(),
	)

	// second movement
	e.Move(gui.NewPos(30, 27))
	assert.Equal(t, gui.NewPos(30, 27), e.Entry.Position())
	assert.Equal(t,
		gui.NewPos(30, 27+entrySize.Height-theme.InputBorderSize()),
		e.popUp.Position(),
	)
}

func TestSelectEntry_DropDownResize(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	options := []string{"A", "B", "C"}
	e := NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(gui.NewSize(150, 200))
	e.Resize(e.MinSize().Max(gui.NewSize(130, 0)))
	e.Move(gui.NewPos(10, 10))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "select_entry/dropdown_initial.xml", c)
	assert.Nil(t, c.Overlays().Top())

	switchPos := gui.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)

	e.Resize(e.Size().Subtract(gui.NewSize(20, 0)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened_shrunk.xml", c)

	e.Resize(e.Size().Add(gui.NewSize(20, 0)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)
}

func TestSelectEntry_MinSize(t *testing.T) {
	smallOptions := []string{"A", "B", "C"}
	largeOptions := []string{"Large Option A", "Larger Option B", "Very Large Option C"}
	labelHeight := NewLabel("W").MinSize().Height

	// since we scroll content and don't prop window open with popup all combinations should be the same min
	tests := map[string]struct {
		placeholder string
		value       string
		options     []string
	}{
		"empty": {},
		"empty + small options": {
			options: smallOptions,
		},
		"empty + large options": {
			options: largeOptions,
		},
		"value": {
			value: "foo", // in a scroller
		},
		"large value + small options": {
			value:   "large", // in a scroller
			options: smallOptions,
		},
		"small value + large options": {
			value:   "small", // in a scroller
			options: largeOptions,
		},
	}

	minSize := gui.NewSize(emptyTextWidth()+dropDownIconWidth()+2*theme.Padding(), labelHeight)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := NewSelectEntry(tt.options)
			e.PlaceHolder = tt.placeholder
			e.Text = tt.value
			assert.Equal(t, minSize, e.MinSize())
		})
	}
}

func TestSelectEntry_SetOptions(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	e := NewSelectEntry([]string{"A", "B", "C"})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(gui.NewSize(150, 200))
	e.Resize(e.MinSize().Max(gui.NewSize(130, 0)))
	e.Move(gui.NewPos(10, 10))
	c := w.Canvas()

	switchPos := gui.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)
	test.TapCanvas(c, switchPos)

	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_setopts.xml", c)
}

func TestSelectEntry_SetOptions_Empty(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	e := NewSelectEntry([]string{})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(gui.NewSize(150, 200))
	e.Resize(e.MinSize().Max(gui.NewSize(130, 0)))
	e.Move(gui.NewPos(10, 10))
	c := w.Canvas()

	switchPos := gui.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_setopts.xml", c)
}

func dropDownIconWidth() float32 {
	return theme.IconInlineSize() + theme.Padding()
}

func emptyTextWidth() float32 {
	return NewLabel("M").MinSize().Width
}
