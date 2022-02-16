//go:build !android && !ios && !mobile
// +build !android,!ios,!mobile

package app

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
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

func TestDefaultTheme(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "windows" { // system defines default for macOS and Windows
		assert.Equal(t, theme.VariantDark, defaultVariant())
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := testPath("guitest")

	ensureDirExists(tmpDir)
	if st, err := os.Stat(tmpDir); err != nil || !st.IsDir() {
		t.Error("Could not ensure directory exists")
	}

	os.Remove(tmpDir)
}

func TestWatchSettings(t *testing.T) {
	settings := &settings{}
	listener := make(chan gui.Settings, 1)
	settings.AddChangeListener(listener)

	settings.fileChanged() // simulate the settings file changing

	select {
	case <-listener:
	case <-time.After(100 * time.Millisecond):
		t.Error("Settings listener was not called")
	}
}

func TestWatchFile(t *testing.T) {
	path := testPath("gui-temp-watch.txt")
	f, _ := os.Create(path)
	f.Close()
	defer os.Remove(path)

	called := make(chan interface{}, 1)
	watchFile(path, func() {
		called <- true
	})
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	file.WriteString(" ")
	file.Close()

	select {
	case <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
}

func TestFileWatcher_FileDeleted(t *testing.T) {
	path := testPath("gui-temp-watch.txt")
	f, _ := os.Create(path)
	f.Close()
	defer os.Remove(path)

	called := make(chan interface{}, 1)
	watcher := watchFile(path, func() {
		called <- true
	})
	if watcher == nil {
		assert.Fail(t, "Could not start watcher")
		return
	}

	defer watcher.Close()
	os.Remove(path)
	f, _ = os.Create(path)

	select {
	case <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
	f.Close()
}

func testPath(child string) string {
	// TMPDIR would be more normal but fsnotify cannot watch that on macOS...
	return filepath.Join("testdata", child)
}
