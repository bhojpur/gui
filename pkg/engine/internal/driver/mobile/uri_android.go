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
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>

char* contentURIGetFileName(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
*/
import "C"
import (
	"path/filepath"
	"unsafe"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/app"
	"github.com/bhojpur/gui/pkg/engine/storage"
)

type androidURI struct {
	systemURI string
	gui.URI
}

// Override Name on android for content://
func (a *androidURI) Name() string {
	if a.Scheme() == "content" {
		result := contentURIGetFileName(a.systemURI)
		if result != "" {
			return result
		}
	}
	return a.URI.Name()
}

func (a *androidURI) Extension() string {
	return filepath.Ext(a.Name())
}

func contentURIGetFileName(uri string) string {
	uriStr := C.CString(uri)
	defer C.free(unsafe.Pointer(uriStr))

	var filename string
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		fnamePtr := C.contentURIGetFileName(C.uintptr_t(env), C.uintptr_t(ctx), uriStr)
		vPtr := unsafe.Pointer(fnamePtr)
		if vPtr == C.NULL {
			return nil
		}
		filename = C.GoString(fnamePtr)
		C.free(vPtr)

		return nil
	})
	return filename
}

func nativeURI(uri string) gui.URI {
	return &androidURI{URI: storage.NewURI(uri), systemURI: uri}
}
