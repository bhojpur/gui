package gl

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

import "unsafe"

type call struct {
	args     fnargs
	parg     unsafe.Pointer
	blocking bool
}

type fnargs struct {
	fn glfn

	a0 uintptr
	a1 uintptr
	a2 uintptr
	a3 uintptr
	a4 uintptr
	a5 uintptr
	a6 uintptr
	a7 uintptr
	a8 uintptr
	a9 uintptr
}

type glfn int

const (
	glfnUNDEFINED glfn = iota
	glfnActiveTexture
	glfnAttachShader
	glfnBindBuffer
	glfnBindTexture
	glfnBindVertexArray
	glfnBlendColor
	glfnBlendFunc
	glfnBufferData
	glfnClear
	glfnClearColor
	glfnCompileShader
	glfnCreateProgram
	glfnCreateShader
	glfnDeleteBuffer
	glfnDeleteTexture
	glfnDisable
	glfnDrawArrays
	glfnEnable
	glfnEnableVertexAttribArray
	glfnFlush
	glfnGenBuffer
	glfnGenTexture
	glfnGenVertexArray
	glfnGetAttribLocation
	glfnGetError
	glfnGetShaderInfoLog
	glfnGetShaderSource
	glfnGetShaderiv
	glfnGetTexParameteriv
	glfnGetUniformLocation
	glfnLinkProgram
	glfnReadPixels
	glfnScissor
	glfnShaderSource
	glfnTexImage2D
	glfnTexParameteri
	glfnUniform1f
	glfnUniform4f
	glfnUniform4fv
	glfnUseProgram
	glfnVertexAttribPointer
	glfnViewport
)

func goString(buf []byte) string {
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	panic("buf is not NUL-terminated")
}

func glBoolean(b bool) uintptr {
	if b {
		return True
	}
	return False
}
