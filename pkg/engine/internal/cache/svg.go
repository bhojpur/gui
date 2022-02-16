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
	"image"
	"sync"
	"time"
)

var svgLock sync.RWMutex
var svgs = make(map[string]*svgInfo)

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, w int, h int) *image.NRGBA {
	svgLock.RLock()
	sinfo, ok := svgs[name]
	svgLock.RUnlock()
	if !ok || sinfo == nil || sinfo.w != w || sinfo.h != h {
		return nil
	}
	sinfo.setAlive()
	return sinfo.pix
}

// SetSvg sets a svg into the cache map.
func SetSvg(name string, pix *image.NRGBA, w int, h int) {
	sinfo := &svgInfo{
		pix: pix,
		w:   w,
		h:   h,
	}
	sinfo.setAlive()
	svgLock.Lock()
	svgs[name] = sinfo
	svgLock.Unlock()
}

type svgInfo struct {
	expiringCacheNoLock
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	expiredSvgs := make([]string, 0, 20)
	svgLock.RLock()
	for s, sinfo := range svgs {
		if sinfo.isExpired(now) {
			expiredSvgs = append(expiredSvgs, s)
		}
	}
	svgLock.RUnlock()
	if len(expiredSvgs) > 0 {
		svgLock.Lock()
		for _, exp := range expiredSvgs {
			delete(svgs, exp)
		}
		svgLock.Unlock()
	}
}
