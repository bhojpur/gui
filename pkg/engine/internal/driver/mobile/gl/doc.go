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
It implements Go bindings for OpenGL ES 2.0 and ES 3.0.

The GL functions are defined on a Context object that is responsible for
tracking a GL context. Typically, a windowing system package (such as
golang.org/x/exp/shiny/screen) will call NewContext and provide
a gl.Context for a user application.

If the gl package is compiled on a platform capable of supporting ES 3.0,
the gl.Context object also implements gl.Context3.

The bindings are deliberately minimal, staying as close the C API as
possible. The semantics of each function maps onto functions
described in the Khronos documentation:

https://www.khronos.org/opengles/sdk/docs/man/

One notable departure from the C API is the introduction of types
to represent common uses of GLint: Texture, Surface, Buffer, etc.
*/
package gl // import "github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/gl"

/*
Implementation details.

All GL function calls fill out a C.struct_fnargs and drop it on the work
queue. The Start function drains the work queue and hands over a batch
of calls to C.process which runs them. This allows multiple GL calls to
be executed in a single cgo call.

A GL call is marked as blocking if it returns a value, or if it takes a
Go pointer. In this case the call will not return until C.process sends a
value on the retvalue channel.

This implementation ensures any goroutine can make GL calls, but it does
not make the GL interface safe for simultaneous use by multiple goroutines.
For the purpose of analyzing this code for race conditions, picture two
separate goroutines: one blocked on gl.Start, and another making calls to
the gl package exported functions.
*/
