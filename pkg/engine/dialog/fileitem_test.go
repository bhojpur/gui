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

	"github.com/bhojpur/gui/pkg/engine/storage"

	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestFileItem_Name(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem(storage.NewFileURI("/path/to/filename.txt"), false)
	assert.Equal(t, "filename", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/MyFile.jpeg"), false)
	assert.Equal(t, "MyFile", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/.maybeHidden.txt"), false)
	assert.Equal(t, ".maybeHidden", item.name)
}

func TestFileItem_FolderName(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem(storage.NewFileURI("/path/to/foldername/"), true)
	assert.Equal(t, "foldername", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/myapp.app/"), true)
	assert.Equal(t, "myapp.app", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/.maybeHidden/"), true)
	assert.Equal(t, ".maybeHidden", item.name)
}

func TestNewFileItem(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	item := f.newFileItem(storage.NewFileURI("/path/to/filename.txt"), false)

	assert.Equal(t, "filename", item.name)

	test.Tap(item)
	assert.True(t, item.isCurrent)
	assert.Equal(t, item, f.selected)
}

func TestNewFileItem_Folder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	currentLister, err := storage.ListerForURI(storage.NewFileURI(currentDir))
	if err != nil {
		t.Error(err)
	}

	parentDir := storage.NewFileURI(filepath.Dir(currentDir))
	parentLister, err := storage.ListerForURI(parentDir)
	if err != nil {
		t.Error(err)
	}
	f.setLocation(parentLister)
	item := f.newFileItem(currentLister, true)

	assert.Equal(t, filepath.Base(currentDir), item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, storage.NewFileURI(""+currentDir).String(), f.dir.String())
}

func TestNewFileItem_ParentFolder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	currentLister, err := storage.ListerForURI(storage.NewFileURI(currentDir))
	if err != nil {
		t.Error(err)
	}
	parentDir := storage.NewFileURI(filepath.Dir(currentDir))
	f.setLocation(currentLister)

	item := &fileDialogItem{picker: f, name: "(Parent)", location: parentDir, dir: true}
	item.ExtendBaseWidget(item)

	assert.Equal(t, "(Parent)", item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, parentDir.String(), f.dir.String())
}
