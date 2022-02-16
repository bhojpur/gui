//go:build !ci && android
// +build !ci,android

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
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>

void openURL(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *url);
void sendNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *title, char *content);
*/
import "C"
import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"unsafe"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/app"
)

func (a *bhojpurApp) OpenURL(url *url.URL) error {
	urlStr := C.CString(url.String())
	defer C.free(unsafe.Pointer(urlStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.openURL(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), urlStr)
		return nil
	})
	return nil
}

func (a *bhojpurApp) SendNotification(n *gui.Notification) {
	titleStr := C.CString(n.Title)
	defer C.free(unsafe.Pointer(titleStr))
	contentStr := C.CString(n.Content)
	defer C.free(unsafe.Pointer(contentStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.sendNotification(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), titleStr, contentStr)
		return nil
	})
}

func defaultVariant() gui.ThemeVariant {
	return systemTheme
}

func rootConfigDir() string {
	filesDir := os.Getenv("FILESDIR")
	if filesDir == "" {
		log.Println("FILESDIR env was not set by android native code")
		return "/data/data" // probably won't work, but we can't make a better guess
	}

	return filepath.Join(filesDir, "bhojpur")
}
