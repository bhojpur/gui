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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

const (
	fileIconSize       = 64
	fileInlineIconSize = 24
	fileTextSize       = 24
	fileIconCellWidth  = fileIconSize * 1.25
)

type fileDialogItem struct {
	widget.BaseWidget
	picker    *fileDialog
	isCurrent bool

	name     string
	location gui.URI
	dir      bool

	hovered bool
}

func (i *fileDialogItem) MouseIn(*desktop.MouseEvent) {
	i.hovered = true
	i.Refresh()
}

func (i *fileDialogItem) MouseMoved(*desktop.MouseEvent) {
}

func (i *fileDialogItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

func (i *fileDialogItem) Tapped(_ *gui.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogItem) CreateRenderer() gui.WidgetRenderer {
	background := canvas.NewRectangle(theme.PrimaryColor())
	background.Hide()
	text := widget.NewLabelWithStyle(i.name, gui.TextAlignCenter, gui.TextStyle{})
	text.Wrapping = gui.TextTruncate
	icon := widget.NewFileIcon(i.location)

	return &fileItemRenderer{
		item:       i,
		background: background,
		icon:       icon,
		text:       text,
		objects:    []gui.CanvasObject{background, icon, text},
	}
}

func fileName(path gui.URI) (name string) {
	pathstr := path.String()[len(path.Scheme())+3:]
	name = filepath.Base(pathstr)
	ext := filepath.Ext(name[1:])
	name = name[:len(name)-len(ext)]

	return
}

func (i *fileDialogItem) isDirectory() bool {
	return i.dir
}

func (f *fileDialog) newFileItem(location gui.URI, dir bool) *fileDialogItem {
	item := &fileDialogItem{
		picker:   f,
		location: location,
		dir:      dir,
	}

	if dir {
		item.name = location.Name()
	} else {
		item.name = fileName(location)
	}

	item.ExtendBaseWidget(item)
	return item
}

type fileItemRenderer struct {
	item *fileDialogItem

	background *canvas.Rectangle
	icon       *widget.FileIcon
	text       *widget.Label
	objects    []gui.CanvasObject
}

func (s fileItemRenderer) Layout(size gui.Size) {
	s.background.Resize(size)

	if s.item.picker.view == gridView {
		s.icon.Resize(gui.NewSize(fileIconSize, fileIconSize))
		s.icon.Move(gui.NewPos((size.Width-fileIconSize)/2, 0))

		s.text.Alignment = gui.TextAlignCenter
		s.text.Resize(gui.NewSize(size.Width, fileTextSize))
		s.text.Move(gui.NewPos(0, size.Height-s.text.MinSize().Height))
	} else {
		s.icon.Resize(gui.NewSize(fileInlineIconSize, fileInlineIconSize))
		s.icon.Move(gui.NewPos(theme.Padding(), (size.Height-fileInlineIconSize)/2))

		s.text.Alignment = gui.TextAlignLeading
		s.text.Resize(gui.NewSize(size.Width, fileTextSize))
		s.text.Move(gui.NewPos(fileInlineIconSize, (size.Height-s.text.MinSize().Height)/2))
	}
	s.text.Refresh()
}

func (s fileItemRenderer) MinSize() gui.Size {
	var padding gui.Size

	if s.item.picker.view == gridView {
		padding = gui.NewSize(fileIconCellWidth-fileIconSize, theme.Padding())
		return gui.NewSize(fileIconSize, fileIconSize+fileTextSize).Add(padding)
	}

	padding = gui.NewSize(theme.Padding(), theme.Padding()*4)
	return gui.NewSize(fileInlineIconSize+s.text.MinSize().Width, fileTextSize).Add(padding)
}

func (s fileItemRenderer) Refresh() {
	if s.item.isCurrent {
		s.background.FillColor = theme.SelectionColor()
		s.background.Show()
	} else if s.item.hovered {
		s.background.FillColor = theme.HoverColor()
		s.background.Show()
	} else {
		s.background.Hide()
	}
	s.background.Refresh()
	canvas.Refresh(s.item)
}

func (s fileItemRenderer) Objects() []gui.CanvasObject {
	return s.objects
}

func (s fileItemRenderer) Destroy() {
}
