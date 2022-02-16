package widget

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
	"image/color"
	"sync"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	col "github.com/bhojpur/gui/pkg/engine/internal/color"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

const cursorInterruptTime = 300 * time.Millisecond

type entryCursorAnimation struct {
	mu                *sync.RWMutex
	cursor            *canvas.Rectangle
	anim              *gui.Animation
	lastInterruptTime time.Time

	timeNow func() time.Time // useful for testing
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor, timeNow: time.Now}
	return a
}

// creates Bhojpur GUI animation
func (a *entryCursorAnimation) createAnim(inverted bool) *gui.Animation {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := col.ToNRGBA(theme.PrimaryColor())
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}
	start, end := color.Color(cursorDim), cursorOpaque
	if inverted {
		start, end = cursorOpaque, color.Color(cursorDim)
	}
	interrupted := false
	anim := canvas.NewColorRGBAAnimation(start, end, time.Second/2, func(c color.Color) {
		a.mu.RLock()
		shouldInterrupt := a.timeNow().Sub(a.lastInterruptTime) <= cursorInterruptTime
		a.mu.RUnlock()
		if shouldInterrupt {
			if !interrupted {
				a.cursor.FillColor = cursorOpaque
				a.cursor.Refresh()
				interrupted = true
			}
			return
		}
		if interrupted {
			a.mu.Lock()
			a.anim.Stop()
			if !inverted {
				a.anim = a.createAnim(true)
			}
			interrupted = false
			a.mu.Unlock()
			go func() {
				a.mu.RLock()
				canStart := a.anim != nil
				a.mu.RUnlock()
				if canStart {
					a.anim.Start()
				}
			}()
			return
		}
		a.cursor.FillColor = c
		a.cursor.Refresh()
	})

	anim.RepeatCount = gui.AnimationRepeatForever
	anim.AutoReverse = true
	return anim
}

// starts cursor animation.
func (a *entryCursorAnimation) start() {
	a.mu.Lock()
	isStopped := a.anim == nil
	if isStopped {
		a.anim = a.createAnim(false)
	}
	a.mu.Unlock()
	if isStopped {
		a.anim.Start()
	}
}

// temporarily stops the animation by "cursorInterruptTime".
func (a *entryCursorAnimation) interrupt() {
	a.mu.Lock()
	a.lastInterruptTime = a.timeNow()
	a.mu.Unlock()
}

// stops cursor animation.
func (a *entryCursorAnimation) stop() {
	a.mu.Lock()
	if a.anim != nil {
		a.anim.Stop()
		a.anim = nil
	}
	a.mu.Unlock()
}
