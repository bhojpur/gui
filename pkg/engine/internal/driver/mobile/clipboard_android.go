//go:build android
// +build android

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

/*
#cgo LDFLAGS: -landroid -llog -lEGL -lGLESv2

#include <stdlib.h>

char *getClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
void setClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *content);
*/
import "C"
import (
	"unsafe"

	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/app"
)

// Content returns the clipboard content for Android
func (c *mobileClipboard) Content() string {
	content := ""
	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getClipboardContent(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		if chars == nil {
			return nil
		}

		content = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})
	return content
}

// SetContent sets the clipboard content for Android
func (c *mobileClipboard) SetContent(content string) {
	contentStr := C.CString(content)
	defer C.free(unsafe.Pointer(contentStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.setClipboardContent(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), contentStr)
		return nil
	})
}
