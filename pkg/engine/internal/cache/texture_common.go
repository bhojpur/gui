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

var textures = sync.Map{} // map[gui.CanvasObject]*textureInfo

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj gui.CanvasObject) {
	textures.Delete(obj)
}

// GetTexture gets cached texture.
func GetTexture(obj gui.CanvasObject) (TextureType, bool) {
	t, ok := textures.Load(obj)
	if t == nil || !ok {
		return noTexture, false
	}
	texInfo := t.(*textureInfo)
	texInfo.setAlive()
	return texInfo.texture, true
}

// RangeExpiredTexturesFor range over the expired textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeExpiredTexturesFor(canvas gui.Canvas, f func(gui.CanvasObject)) {
	now := timeNow()
	textures.Range(func(key, value interface{}) bool {
		obj, tinfo := key.(gui.CanvasObject), value.(*textureInfo)
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// RangeTexturesFor range over the textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeTexturesFor(canvas gui.Canvas, f func(gui.CanvasObject)) {
	textures.Range(func(key, value interface{}) bool {
		obj, tinfo := key.(gui.CanvasObject), value.(*textureInfo)
		if tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// SetTexture sets cached texture.
func SetTexture(obj gui.CanvasObject, texture TextureType, canvas gui.Canvas) {
	texInfo := &textureInfo{texture: texture}
	texInfo.canvas = canvas
	texInfo.setAlive()
	textures.Store(obj, texInfo)
}

// textureCacheBase defines base texture cache object.
type textureCacheBase struct {
	expiringCacheNoLock
	canvas gui.Canvas
}
