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

var (
	fontSizeCache = map[fontSizeEntry]fontMetric{}
	fontSizeLock  = sync.RWMutex{}
)

type fontMetric struct {
	size     gui.Size
	baseLine float32
}

type fontSizeEntry struct {
	text  string
	size  float32
	style gui.TextStyle
}

// GetFontMetrics looks up a calculated size and baseline required for the specified text parameters.
func GetFontMetrics(text string, fontSize float32, style gui.TextStyle) (size gui.Size, base float32) {
	ent := fontSizeEntry{text, fontSize, style}
	fontSizeLock.RLock()
	ret, ok := fontSizeCache[ent]
	fontSizeLock.RUnlock()
	if !ok {
		return gui.Size{Width: 0, Height: 0}, 0
	}
	return ret.size, ret.baseLine
}

// SetFontMetrics stores a calculated font size and baseline for parameters that were missing from the cache.
func SetFontMetrics(text string, fontSize float32, style gui.TextStyle, size gui.Size, base float32) {
	ent := fontSizeEntry{text, fontSize, style}
	fontSizeLock.Lock()
	fontSizeCache[ent] = fontMetric{size: size, baseLine: base}
	fontSizeLock.Unlock()
}
