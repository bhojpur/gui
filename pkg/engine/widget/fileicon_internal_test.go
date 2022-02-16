//go:build !windows
// +build !windows

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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// Simulate being rendered by calling CreateRenderer() to update icon
func newRenderedFileIcon(uri gui.URI) *FileIcon {
	f := NewFileIcon(uri)
	f.CreateRenderer()
	return f
}

func TestFileIcon_NewFileIcon(t *testing.T) {
	item := newRenderedFileIcon(storage.NewFileURI("/path/to/filename.zip"))
	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.mp3"))
	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.png"))
	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.txt"))
	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.mp4"))
	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, theme.FileVideoIcon(), item.resource)
}

func TestFileIcon_NewFileIcon_NoExtension(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		gui.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	binFileWithNoExt := filepath.Join(workingDir, "testdata/bin")
	textFileWithNoExt := filepath.Join(workingDir, "testdata/text")

	item := newRenderedFileIcon(storage.NewFileURI(binFileWithNoExt))
	assert.Equal(t, "", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI(textFileWithNoExt))
	assert.Equal(t, "", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)
}

func TestFileIcon_NewURI_WithFolder(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		gui.LogError("Could not get current working directory", err)
		t.FailNow()
	}

	dir := filepath.Join(workingDir, "testdata")
	item := newRenderedFileIcon(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)

	item.SetURI(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)
}

func TestFileIcon_NewFileIcon_Rendered(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	workingDir, err := os.Getwd()
	if err != nil {
		gui.LogError("Could not get current working directory", err)
		t.FailNow()
	}

	icon := NewFileIcon(nil)

	w := test.NewWindow(icon)
	w.Resize(gui.NewSize(150, 150))

	test.AssertImageMatches(t, "fileicon/fileicon_nil.png", w.Canvas().Capture())

	text := filepath.Join(workingDir, "testdata/text")
	icon2 := NewFileIcon(storage.NewFileURI(text))

	w.SetContent(icon2)
	w.Resize(gui.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_text.png", w.Canvas().Capture())

	text += ".txt"
	icon3 := NewFileIcon(storage.NewFileURI(text))

	w.SetContent(icon3)
	w.Resize(gui.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_text_txt.png", w.Canvas().Capture())

	bin := filepath.Join(workingDir, "testdata/bin")
	icon4 := NewFileIcon(storage.NewFileURI(bin))

	w.SetContent(icon4)
	w.Resize(gui.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_bin.png", w.Canvas().Capture())

	dir := filepath.Join(workingDir, "testdata")
	icon5 := NewFileIcon(storage.NewFileURI(dir))

	w.SetContent(icon5)
	w.Resize(gui.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_folder.png", w.Canvas().Capture())

	w.Close()
}

func TestFileIcon_SetURI(t *testing.T) {
	item := newRenderedFileIcon(storage.NewFileURI("/path/to/filename.zip"))
	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.mp3"))
	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.png"))
	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.txt"))
	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.mp4"))
	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, theme.FileVideoIcon(), item.resource)
}

func TestFileIcon_SetURI_WithFolder(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		gui.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	dir := filepath.Join(workingDir, "testdata")

	item := newRenderedFileIcon(nil)
	assert.Empty(t, item.extension)

	item.SetURI(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)

	item.SetURI(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)
}

func TestFileIcon_DirURIUpdated(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		gui.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	testDir := filepath.Join(workingDir, "testdata", "notCreatedYet")

	uri := storage.NewFileURI(testDir)
	item := newRenderedFileIcon(uri)

	// The directory has not been created. It can not be listed yet.
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	err = os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDir)

	item.Refresh()
	assert.Equal(t, theme.FolderIcon(), item.resource)
}
