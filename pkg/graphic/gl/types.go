//go:build !wasm
// +build !wasm

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

func (v Attrib) IsValid() bool       { return v != NoAttrib }
func (v Program) IsValid() bool      { return v != NoProgram }
func (v Shader) IsValid() bool       { return v != NoShader }
func (v Buffer) IsValid() bool       { return v != NoBuffer }
func (v Framebuffer) IsValid() bool  { return v != NoFramebuffer }
func (v Renderbuffer) IsValid() bool { return v != NoRenderbuffer }
func (v Texture) IsValid() bool      { return v != NoTexture }
func (v Uniform) IsValid() bool      { return v != NoUniform }
