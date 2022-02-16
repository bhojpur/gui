package dialog

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
	"github.com/bhojpur/gui/pkg/engine/widget"
	"github.com/stretchr/testify/assert"
)

func TestProgressInfiniteDialog_MinSize(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()
	d := NewProgressInfinite("title", "message", window)

	dialogContent := d.win.Content.MinSize()
	progressBar := d.bar.MinSize()

	assert.Less(t, progressBar.Width, dialogContent.Width)
}

func TestProgressInfiniteDialog_Resize(t *testing.T) {
	window := test.NewWindow(nil)
	window.Resize(gui.NewSize(600, 400))
	defer window.Close()
	d := NewProgressInfinite("title", "message", window)
	theDialog := d.dialog
	d.dialog.Show() // we cannot check window size if not shown

	//Test resize - normal size scenario
	size := gui.NewSize(300, 180) //normal size to fit (600,400)
	theDialog.Resize(size)
	expectedWidth := float32(300)
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight := float32(180)
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height+theme.Padding()*2)
	//Test resize - normal size scenario again
	size = gui.NewSize(310, 280) //normal size to fit (600,400)
	theDialog.Resize(size)
	expectedWidth = 310
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 280
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height+theme.Padding()*2)

	//Test resize - greater than max size scenario
	size = gui.NewSize(800, 600)
	theDialog.Resize(size)
	expectedWidth = 600                                        //since win width only 600
	assert.Equal(t, expectedWidth, theDialog.win.Size().Width) //max, also work
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 400                                         //since win heigh only 400
	assert.Equal(t, expectedHeight, theDialog.win.Size().Height) //max, also work
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height+theme.Padding()*2)

	//Test again - extreme small size
	size = gui.NewSize(1, 1)
	theDialog.Resize(size)
	expectedWidth = theDialog.win.Content.MinSize().Width
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width)
	expectedHeight = theDialog.win.Content.MinSize().Height
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height)
}

func TestProgressInfiniteDialog_Content(t *testing.T) {
	title := "title"
	message := "message"

	window := test.NewWindow(nil)
	defer window.Close()
	d := NewProgressInfinite(title, message, window)

	assert.Equal(t, d.title, title)
	assert.Equal(t, d.content.(*widget.Label).Text, message)
}

func TestProgressInfiniteDialog_Show(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()
	d := NewProgressInfinite("title", "message", window)

	d.Show()

	assert.False(t, d.win.Hidden)
	assert.True(t, d.bar.Running())

	d.Hide()

	assert.True(t, d.win.Hidden)
	assert.False(t, d.bar.Running())
}
