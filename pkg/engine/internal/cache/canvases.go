package cache

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
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
)

var canvasesLock sync.RWMutex
var canvases = make(map[gui.CanvasObject]*canvasInfo, 1024)

// GetCanvasForObject returns the canvas for the specified object.
func GetCanvasForObject(obj gui.CanvasObject) gui.Canvas {
	canvasesLock.RLock()
	cinfo, ok := canvases[obj]
	canvasesLock.RUnlock()
	if cinfo == nil || !ok {
		return nil
	}
	cinfo.setAlive()
	return cinfo.canvas
}

// SetCanvasForObject sets the canvas for the specified object.
func SetCanvasForObject(obj gui.CanvasObject, canvas gui.Canvas) {
	cinfo := &canvasInfo{canvas: canvas}
	cinfo.setAlive()
	canvasesLock.Lock()
	canvases[obj] = cinfo
	canvasesLock.Unlock()
}

type canvasInfo struct {
	expiringCache
	canvas gui.Canvas
}
