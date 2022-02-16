package storage

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
)

// FileFilter is an interface that can be implemented to provide a filter to a file dialog.
type FileFilter interface {
	Matches(gui.URI) bool
}

// ExtensionFileFilter represents a file filter based on the the ending of file names,
// for example ".txt" and ".png".
type ExtensionFileFilter struct {
	Extensions []string
}

// MimeTypeFileFilter represents a file filter based on the files mime type,
// for example "image/*", "audio/mp3".
type MimeTypeFileFilter struct {
	MimeTypes []string
}

// Matches returns true if a file URI has one of the filtered extensions.
func (e *ExtensionFileFilter) Matches(uri gui.URI) bool {
	extension := uri.Extension()
	for _, ext := range e.Extensions {
		if strings.EqualFold(extension, ext) {
			return true
		}
	}
	return false
}

// NewExtensionFileFilter takes a string slice of extensions with a leading . and creates a filter for the file dialog.
// Example: .jpg, .mp3, .txt, .sh
func NewExtensionFileFilter(extensions []string) FileFilter {
	return &ExtensionFileFilter{Extensions: extensions}
}

// Matches returns true if a file URI has one of the filtered mimetypes.
func (mt *MimeTypeFileFilter) Matches(uri gui.URI) bool {
	mimeType, mimeSubType := splitMimeType(uri)
	for _, mimeTypeFull := range mt.MimeTypes {
		mimeTypeSplit := strings.Split(mimeTypeFull, "/")
		if len(mimeTypeSplit) <= 1 {
			continue
		}
		mType := mimeTypeSplit[0]
		mSubType := strings.Split(mimeTypeSplit[1], ";")[0]
		if mType == mimeType {
			if mSubType == mimeSubType || mSubType == "*" {
				return true
			}
		}
	}
	return false
}

// NewMimeTypeFileFilter takes a string slice of mimetypes, including globs, and creates a filter for the file dialog.
// Example: image/*, audio/mp3, text/plain, application/*
func NewMimeTypeFileFilter(mimeTypes []string) FileFilter {
	return &MimeTypeFileFilter{MimeTypes: mimeTypes}
}

func splitMimeType(uri gui.URI) (mimeType, mimeSubType string) {
	mimeTypeFull := uri.MimeType()
	mimeTypeSplit := strings.Split(mimeTypeFull, "/")
	if len(mimeTypeSplit) <= 1 {
		mimeType, mimeSubType = "", ""
		return
	}
	mimeType = mimeTypeSplit[0]
	mimeSubType = mimeTypeSplit[1]

	return
}
