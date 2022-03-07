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
It is a Go cross-platform binding for OpenGL, with an OpenGL ES 2-like API.

It supports:
- macOS, Linux and Windows via OpenGL 2.1 backend,
- iOS and Android via OpenGL ES 2.0 backend,
- Modern Browsers (desktop and mobile) via WebGL 1.0 backend.

This is a fork of golang.org/x/mobile/gl package with [CL 8793](https://go-review.googlesource.com/8793)
merged in and Windows support added. It is fully functional, but may eventually become superceded by
the new x/mobile/gl plan. It will exist and be fully supported until it can be safely replaced by a
better package.

Usage

This OpenGL binding has a ContextWatcher, which implements
[glfw.ContextWatcher](https://godoc.org/github.com/bhojpur/gui/pkg/graphic/glfw#ContextWatcher)
interface. Recommended usage is with github.com/bhojpur/gui/pkg/graphic/glfw package, which accepts
a ContextWatcher in its Init, and takes on the responsibility of notifying it when context is made
current or detached.

	if err := glfw.Init(gl.ContextWatcher); err != nil {
		// Handle error.
	}
	defer glfw.Terminate()

If you're not using a ContextWatcher-aware glfw library, you must call methods of
gl.ContextWatcher yourself whenever you make a context current or detached.

	window.MakeContextCurrent()
	gl.ContextWatcher.OnMakeCurrent(nil)

	glfw.DetachCurrentContext()
	gl.ContextWatcher.OnDetach()
*/
package gl // import "github.com/bhojpur/gui/pkg/graphic/gl"
