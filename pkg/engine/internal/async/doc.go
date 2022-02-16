//go:generate go run gen.go

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

// It provides unbounded channel data structures that are
// designed for caching unlimited number of a concrete type. For better
// performance, a given type should be less or euqal than 16 bytes.
//
// Delicate dance: One must aware that an unbounded channel may lead to
// OOM when the consuming speed of the buffer is lower than the producing
// speed constantly. However, such a channel may be fairly used for event
// delivering if the consumer of the channel consumes the incoming
// forever, such as even processing.
//
// One must close such a channel via Close() method, closing the input
// channel via close() built-in method can leads to memory leak.
//
// To support a new type, one may add the required data in the gen.go,
// for instances:
//
// 	types := map[string]data{
// 		"bhojpur_canvasobject.go": data{
// 			Type: "gui.CanvasObject",
// 			Name: "CanvasObject",
// 			Imports: `import gui "github.com/bhojpur/gui/pkg/engine"`,
// 		},
// 		"func.go": data{
// 			Type:    "func()",
// 			Name:    "Func",
// 			Imports: "",
// 		},
// 	}
//
// then run: `go generate ./...` to generate more desired unbounded channels.
package async
