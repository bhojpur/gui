//go:build js && wasm
// +build js,wasm

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

import "syscall/js"

type Enum int

type Attrib struct {
	Value int
}

type Program struct {
	js.Value
}

type Shader struct {
	js.Value
}

type Buffer struct {
	js.Value
}

type Framebuffer struct {
	js.Value
}

type Renderbuffer struct {
	js.Value
}

type Texture struct {
	js.Value
}

type Uniform struct {
	js.Value
}

var NoAttrib Attrib
var NoProgram = Program{js.Null()}
var NoShader = Shader{js.Null()}
var NoBuffer = Buffer{js.Null()}
var NoFramebuffer = Framebuffer{js.Null()}
var NoRenderbuffer = Renderbuffer{js.Null()}
var NoTexture = Texture{js.Null()}
var NoUniform = Uniform{js.Null()}

func (v Attrib) IsValid() bool       { return v != NoAttrib }
func (v Program) IsValid() bool      { return !v.Equal(NoProgram.Value) }
func (v Shader) IsValid() bool       { return !v.Equal(NoShader.Value) }
func (v Buffer) IsValid() bool       { return !v.Equal(NoBuffer.Value) }
func (v Framebuffer) IsValid() bool  { return !v.Equal(NoFramebuffer.Value) }
func (v Renderbuffer) IsValid() bool { return !v.Equal(NoRenderbuffer.Value) }
func (v Texture) IsValid() bool      { return !v.Equal(NoTexture.Value) }
func (v Uniform) IsValid() bool      { return !v.Equal(NoUniform.Value) }
