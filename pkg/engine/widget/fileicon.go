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
	"strings"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

const (
	ratioDown     = 0.45
	ratioTextSize = 0.22
)

// FileIcon is an adaption of widget.Icon for showing files and folders
//
// Since: 1.4
type FileIcon struct {
	BaseWidget

	Selected bool
	URI      gui.URI

	resource  gui.Resource
	extension string
}

// NewFileIcon takes a filepath and creates an icon with an overlaid label using the detected mimetype and extension
//
// Since: 1.4
func NewFileIcon(uri gui.URI) *FileIcon {
	i := &FileIcon{URI: uri}
	i.ExtendBaseWidget(i)
	return i
}

// SetURI changes the URI and makes the icon reflect a different file
func (i *FileIcon) SetURI(uri gui.URI) {
	i.URI = uri
	i.Refresh()
}

func (i *FileIcon) setURI(uri gui.URI) {
	if uri == nil {
		i.resource = theme.FileIcon()
		return
	}

	i.URI = uri
	i.resource = i.lookupIcon(i.URI)
	i.extension = trimmedExtension(uri)
}

// MinSize returns the size that this widget should not shrink below
func (i *FileIcon) MinSize() gui.Size {
	i.ExtendBaseWidget(i)
	return i.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (i *FileIcon) CreateRenderer() gui.WidgetRenderer {
	i.ExtendBaseWidget(i)
	i.propertyLock.Lock()
	i.setURI(i.URI)
	i.propertyLock.Unlock()

	i.propertyLock.RLock()
	defer i.propertyLock.RUnlock()

	// TODO remove background when `SetSelected` is gone.
	background := canvas.NewRectangle(theme.SelectionColor())
	background.Hide()

	s := &fileIconRenderer{file: i, background: background}
	s.img = canvas.NewImageFromResource(s.file.resource)
	s.img.FillMode = canvas.ImageFillContain
	s.ext = canvas.NewText(s.file.extension, theme.BackgroundColor())
	s.ext.Alignment = gui.TextAlignCenter

	s.SetObjects([]gui.CanvasObject{s.background, s.img, s.ext})

	return s
}

// SetSelected makes the file look like it is selected.
//
// Deprecated: Selection is now handled externally.
func (i *FileIcon) SetSelected(selected bool) {
	i.Selected = selected
	i.Refresh()
}

func (i *FileIcon) lookupIcon(uri gui.URI) gui.Resource {
	if i.isDir(uri) {
		return theme.FolderIcon()
	}

	switch splitMimeType(uri) {
	case "application":
		return theme.FileApplicationIcon()
	case "audio":
		return theme.FileAudioIcon()
	case "image":
		return theme.FileImageIcon()
	case "text":
		return theme.FileTextIcon()
	case "video":
		return theme.FileVideoIcon()
	default:
		return theme.FileIcon()
	}
}

func (i *FileIcon) isDir(uri gui.URI) bool {
	if _, ok := uri.(gui.ListableURI); ok {
		return true
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		return false
	}

	i.URI = listable // Avoid having to call storage.ListerForURI(uri) the next time.
	return true
}

type fileIconRenderer struct {
	widget.BaseRenderer

	file *FileIcon

	background *canvas.Rectangle
	ext        *canvas.Text
	img        *canvas.Image
}

func (s *fileIconRenderer) MinSize() gui.Size {
	size := theme.IconInlineSize()
	return gui.NewSize(size, size)
}

func (s *fileIconRenderer) Layout(size gui.Size) {
	isize := gui.Min(size.Width, size.Height)

	xoff := float32(0)
	yoff := (size.Height - isize) / 2

	if size.Width > size.Height {
		xoff = (size.Width - isize) / 2
	}
	yoff += isize * ratioDown

	oldSize := s.ext.TextSize
	s.ext.TextSize = float32(int(isize * ratioTextSize))
	s.ext.Resize(gui.NewSize(isize, s.ext.MinSize().Height))
	s.ext.Move(gui.NewPos(xoff, yoff))
	if oldSize != s.ext.TextSize {
		s.ext.Refresh()
	}

	s.Objects()[0].Resize(size)
	s.Objects()[1].Resize(size)
}

func (s *fileIconRenderer) Refresh() {
	s.file.propertyLock.Lock()
	s.file.setURI(s.file.URI)
	s.file.propertyLock.Unlock()

	s.file.propertyLock.RLock()
	s.img.Resource = s.file.resource
	s.ext.Text = s.file.extension
	s.file.propertyLock.RUnlock()

	if s.file.Selected {
		s.background.Show()
		s.ext.Color = theme.SelectionColor()
		if _, ok := s.img.Resource.(*theme.InvertedThemedResource); !ok {
			s.img.Resource = theme.NewInvertedThemedResource(s.img.Resource)
		}
	} else {
		s.background.Hide()
		s.ext.Color = theme.BackgroundColor()
		if res, ok := s.img.Resource.(*theme.InvertedThemedResource); ok {
			s.img.Resource = res.Original()
		}
	}

	canvas.Refresh(s.file.super())
	canvas.Refresh(s.ext)
}

func trimmedExtension(uri gui.URI) string {
	ext := uri.Extension()
	if len(ext) > 5 {
		ext = ext[:5]
	}
	return ext
}

func splitMimeType(uri gui.URI) string {
	mimeTypeSplit := strings.Split(uri.MimeType(), "/")
	if len(mimeTypeSplit) <= 1 {
		return ""
	}
	return mimeTypeSplit[0]
}
