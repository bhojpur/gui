//go:build js && !wasm
// +build js,!wasm

package glfw

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

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

func newContext(canvas *js.Object, ca *contextAttributes) (context *js.Object, err error) {
	if js.Global.Get("WebGLRenderingContext") == js.Undefined {
		return nil, errors.New("Your browser doesn't appear to support WebGL.")
	}

	attrs := map[string]bool{
		"alpha":                           ca.Alpha,
		"depth":                           ca.Depth,
		"stencil":                         ca.Stencil,
		"antialias":                       ca.Antialias,
		"premultipliedAlpha":              ca.PremultipliedAlpha,
		"preserveDrawingBuffer":           ca.PreserveDrawingBuffer,
		"preferLowPowerToHighPerformance": ca.PreferLowPowerToHighPerformance,
		"failIfMajorPerformanceCaveat":    ca.FailIfMajorPerformanceCaveat,
	}

	if gl := canvas.Call("getContext", "webgl", attrs); gl != nil {
		return gl, nil
	} else if gl := canvas.Call("getContext", "experimental-webgl", attrs); gl != nil {
		return gl, nil
	} else {
		return nil, errors.New("Creating a WebGL context has failed.")
	}
}

type contextAttributes struct {
	Alpha                           bool
	Depth                           bool
	Stencil                         bool
	Antialias                       bool
	PremultipliedAlpha              bool
	PreserveDrawingBuffer           bool
	PreferLowPowerToHighPerformance bool
	FailIfMajorPerformanceCaveat    bool
}

func defaultAttributes() *contextAttributes {
	return &contextAttributes{
		Alpha:                 false,
		Depth:                 true,
		Stencil:               false,
		Antialias:             false,
		PremultipliedAlpha:    false,
		PreserveDrawingBuffer: false,
	}
}
