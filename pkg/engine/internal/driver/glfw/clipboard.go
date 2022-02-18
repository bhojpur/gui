//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

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
	"runtime"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Declare conformity with Clipboard interface
var _ gui.Clipboard = (*clipboard)(nil)

// clipboard represents the system clipboard
type clipboard struct {
	window *glfw.Window
}

// Content returns the clipboard content
func (c *clipboard) Content() string {
	// This retry logic is to work around the "Access Denied" error often thrown in windows PR#1679
	if runtime.GOOS != "windows" {
		return c.content()
	}
	for i := 3; i > 0; i-- {
		cb := c.content()
		if cb != "" {
			return cb
		}
		time.Sleep(50 * time.Millisecond)
	}
	//can't log retry as it would alos log errors for an empty clipboard
	return ""
}

func (c *clipboard) content() string {
	content := ""
	runOnMain(func() {
		content = glfw.GetClipboardString()
	})
	return content
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	// This retry logic is to work around the "Access Denied" error often thrown in windows PR#1679
	if runtime.GOOS != "windows" {
		c.setContent(content)
		return
	}
	for i := 3; i > 0; i-- {
		c.setContent(content)
		if c.content() == content {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	gui.LogError("GLFW clipboard set failed", nil)
}

func (c *clipboard) setContent(content string) {
	runOnMain(func() {
		defer func() {
			if r := recover(); r != nil {
				gui.LogError("GLFW clipboard error (details above)", nil)
			}
		}()

		glfw.SetClipboardString(content)
	})
}
