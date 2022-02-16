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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestShowFolderOpen(t *testing.T) {
	var chosen gui.ListableURI
	var openErr error
	win := test.NewWindow(widget.NewLabel("OpenDir"))
	d := NewFolderOpen(func(file gui.ListableURI, err error) {
		chosen = file
		openErr = err
	}, win)
	testData, _ := filepath.Abs("testdata")
	dir, err := storage.ListerForURI(storage.NewFileURI(testData))
	if err != nil {
		t.Error("Failed to open testdata dir", err)
	}
	d.SetLocation(dir)
	d.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*gui.Container)
	title := ui.Objects[1].(*gui.Container).Objects[1].(*widget.Label)
	assert.Equal(t, "Open Folder", title.Text)

	nameLabel := ui.Objects[2].(*gui.Container).Objects[1].(*container.Scroll).Content.(*widget.Label)
	buttons := ui.Objects[2].(*gui.Container).Objects[0].(*gui.Container)
	open := buttons.Objects[1].(*widget.Button)

	files := ui.Objects[0].(*container.Split).Trailing.(*gui.Container).Objects[1].(*container.Scroll).Content.(*gui.Container).Objects[0].(*gui.Container)
	assert.Greater(t, len(files.Objects), 0)

	fileName := files.Objects[0].(*fileDialogItem).name
	assert.Equal(t, "(Parent)", fileName)
	assert.False(t, open.Disabled())

	var target *fileDialogItem
	for _, icon := range files.Objects {
		if icon.(*fileDialogItem).dir {
			target = icon.(*fileDialogItem)
		} else {
			t.Error("Folder dialog should not list files")
		}
	}

	assert.NotNil(t, target, "Failed to find folder in testdata")
	test.Tap(target)
	assert.Equal(t, target.location.Name(), nameLabel.Text)
	assert.False(t, open.Disabled())

	test.Tap(open)
	assert.Nil(t, win.Canvas().Overlays().Top())
	assert.Nil(t, openErr)

	assert.Equal(t, target.location.String(), chosen.String())
}
