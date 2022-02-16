package engine

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

// import "github.com/bhojpur/gui/pkg/engine"

// It describes the objects and components available to any Bhojpur GUI application.
// These can all be created, manipulated and tested without rendering (for speed).
// Your main package should use the app package to create an application with
// a default driver that will render your UI.
//
// A simple application may look like this:
//
//   package main
//
//   import "github.com/bhojpur/gui/pkg/engine/app"
//   import "github.com/bhojpur/gui/pkg/engine/container"
//   import "github.com/bhojpur/gui/pkg/engine/widget"
//
//   func main() {
//   	a := app.New()
//   	w := a.NewWindow("Hello")
//
//   	hello := widget.NewLabel("Hello, Bhojpur GUI Developer!")
//   	w.SetContent(container.NewVBox(
//   		hello,
//   		widget.NewButton("Hi!", func() {
//   			hello.SetText("Welcome :)")
//   		}),
//   	))
//
//   	w.ShowAndRun()
//   }
