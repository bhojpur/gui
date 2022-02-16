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

#include <stdbool.h>
#include <stdlib.h>

bool canListURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
bool createListableURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
char *listURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
*/
import "C"
import (
	"errors"
	"strings"
	"unsafe"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/app"
	"github.com/bhojpur/gui/pkg/engine/storage"
)

func canListURI(uri gui.URI) bool {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))
	listable := false

	app.RunOnJVM(func(_, env, ctx uintptr) error {
		listable = bool(C.canListURI(C.uintptr_t(env), C.uintptr_t(ctx), uriStr))
		return nil
	})

	return listable
}

func createListableURI(uri gui.URI) error {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	ok := false
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		ok = bool(C.createListableURI(C.uintptr_t(env), C.uintptr_t(ctx), uriStr))
		return nil
	})

	if ok {
		return nil
	}
	return errors.New("failed to create directory")
}

func listURI(uri gui.URI) ([]gui.URI, error) {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	var str *C.char
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		str = C.listURI(C.uintptr_t(env), C.uintptr_t(ctx), uriStr)
		return nil
	})

	parts := strings.Split(C.GoString(str), "|")
	var list []gui.URI
	for _, part := range parts {
		list = append(list, storage.NewURI(part))
	}
	return list, nil
}
