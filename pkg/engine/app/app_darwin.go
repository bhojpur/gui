//go:build !ci
// +build !ci

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

/*
#cgo CFLAGS: -x objective-c

#cgo LDFLAGS: -framework Foundation

#include <stdbool.h>
#include <stdlib.h>
bool isBundled();
void sendNotification(char *title, char *content);
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	gui "github.com/bhojpur/gui/pkg/engine"
	"golang.org/x/sys/execabs"
)

func (a *bhojpurApp) SendNotification(n *gui.Notification) {
	if C.isBundled() {
		titleStr := C.CString(n.Title)
		defer C.free(unsafe.Pointer(titleStr))
		contentStr := C.CString(n.Content)
		defer C.free(unsafe.Pointer(contentStr))

		C.sendNotification(titleStr, contentStr)
		return
	}

	fallbackNotification(n.Title, n.Content)
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "\\", "\\\\")
	return strings.ReplaceAll(noSlash, "\"", "\\\"")
}

//export fallbackSend
func fallbackSend(cTitle, cContent *C.char) {
	title := C.GoString(cTitle)
	content := C.GoString(cContent)
	fallbackNotification(title, content)
}

func fallbackNotification(title, content string) {
	template := `display notification "%s" with title "%s"`
	script := fmt.Sprintf(template, escapeNotificationString(content), escapeNotificationString(title))

	err := execabs.Command("osascript", "-e", script).Start()
	if err != nil {
		gui.LogError("Failed to launch darwin notify script", err)
	}
}
