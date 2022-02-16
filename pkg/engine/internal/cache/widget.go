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

var renderersLock sync.RWMutex
var renderers = map[gui.Widget]*rendererInfo{}

type isBaseWidget interface {
	ExtendBaseWidget(gui.Widget)
	super() gui.Widget
}

// Renderer looks up the render implementation for a widget
func Renderer(wid gui.Widget) gui.WidgetRenderer {
	if wid == nil {
		return nil
	}

	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}

	renderersLock.RLock()
	rinfo, ok := renderers[wid]
	renderersLock.RUnlock()
	if !ok {
		rinfo = &rendererInfo{renderer: wid.CreateRenderer()}
		renderersLock.Lock()
		renderers[wid] = rinfo
		renderersLock.Unlock()
	}

	if rinfo == nil {
		return nil
	}

	rinfo.setAlive()

	return rinfo.renderer
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid gui.Widget) {
	renderersLock.RLock()
	rinfo, ok := renderers[wid]
	renderersLock.RUnlock()
	if !ok {
		return
	}
	if rinfo != nil {
		rinfo.renderer.Destroy()
	}
	renderersLock.Lock()
	delete(renderers, wid)
	renderersLock.Unlock()
}

// IsRendered returns true of the widget currently has a renderer.
// One will be created the first time a widget is shown but may be removed after it is hidden.
func IsRendered(wid gui.Widget) bool {
	renderersLock.RLock()
	_, found := renderers[wid]
	renderersLock.RUnlock()
	return found
}

type rendererInfo struct {
	expiringCache
	renderer gui.WidgetRenderer
}
