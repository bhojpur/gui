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
	"io"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/app"
	"github.com/bhojpur/gui/pkg/engine/storage"
)

type fileOpen struct {
	io.ReadCloser
	uri  gui.URI
	done func()
}

func (f *fileOpen) URI() gui.URI {
	return f.uri
}

func fileReaderForURI(u gui.URI) (gui.URIReadCloser, error) {
	file := &fileOpen{uri: u}
	read, err := nativeFileOpen(file)
	if read == nil {
		return nil, err
	}
	file.ReadCloser = read
	return file, err
}

func mobileFilter(filter storage.FileFilter) *app.FileFilter {
	mobile := &app.FileFilter{}

	if f, ok := filter.(*storage.MimeTypeFileFilter); ok {
		mobile.MimeTypes = f.MimeTypes
	} else if f, ok := filter.(*storage.ExtensionFileFilter); ok {
		mobile.Extensions = f.Extensions
	} else if filter != nil {
		gui.LogError("Custom filter types not supported on mobile", nil)
	}

	return mobile
}

type hasOpenPicker interface {
	ShowFileOpenPicker(func(string, func()), *app.FileFilter)
}

// ShowFileOpenPicker loads the native file open dialog and returns the chosen file path via the callback func.
func ShowFileOpenPicker(callback func(gui.URIReadCloser, error), filter storage.FileFilter) {
	drv := gui.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasOpenPicker); ok {
		a.ShowFileOpenPicker(func(uri string, closer func()) {
			if uri == "" {
				callback(nil, nil)
				return
			}
			f, err := fileReaderForURI(nativeURI(uri))
			if f != nil {
				f.(*fileOpen).done = closer
			}
			callback(f, err)
		}, mobileFilter(filter))
	}
}

// ShowFolderOpenPicker loads the native folder open dialog and calls back the chosen directory path as a ListableURI.
func ShowFolderOpenPicker(callback func(gui.ListableURI, error)) {
	filter := storage.NewMimeTypeFileFilter([]string{"application/x-directory"})
	drv := gui.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasOpenPicker); ok {
		a.ShowFileOpenPicker(func(path string, _ func()) {
			if path == "" {
				callback(nil, nil)
				return
			}

			uri, err := storage.ParseURI(path)
			if err != nil {
				callback(nil, err)
				return
			}

			callback(listerForURI(uri))
		}, mobileFilter(filter))
	}
}

type fileSave struct {
	io.WriteCloser
	uri  gui.URI
	done func()
}

func (f *fileSave) URI() gui.URI {
	return f.uri
}

func fileWriterForURI(u gui.URI) (gui.URIWriteCloser, error) {
	file := &fileSave{uri: u}
	write, err := nativeFileSave(file)
	if write == nil {
		return nil, err
	}
	file.WriteCloser = write
	return file, err
}

type hasSavePicker interface {
	ShowFileSavePicker(func(string, func()), *app.FileFilter, string)
}

// ShowFileSavePicker loads the native file save dialog and returns the chosen file path via the callback func.
func ShowFileSavePicker(callback func(gui.URIWriteCloser, error), filter storage.FileFilter, filename string) {
	drv := gui.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasSavePicker); ok {
		a.ShowFileSavePicker(func(path string, closer func()) {
			if path == "" {
				callback(nil, nil)
				return
			}

			uri, err := storage.ParseURI(path)
			if err != nil {
				callback(nil, err)
				return
			}

			f, err := fileWriterForURI(uri)
			if f != nil {
				f.(*fileSave).done = closer
			}
			callback(f, err)
		}, mobileFilter(filter), filename)
	}
}
