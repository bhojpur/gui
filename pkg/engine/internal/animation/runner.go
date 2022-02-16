package animation

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
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// Runner is the main driver for animations package
type Runner struct {
	animationMutex    sync.RWMutex
	animations        []*anim
	pendingAnimations []*anim

	runnerStarted bool
}

// Start will register the passed application and initiate its ticking.
func (r *Runner) Start(a *gui.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()

	if !r.runnerStarted {
		r.runnerStarted = true
		r.animations = append(r.animations, newAnim(a))
		r.runAnimations()
	} else {
		r.pendingAnimations = append(r.pendingAnimations, newAnim(a))
	}
}

// Stop causes an animation to stop ticking (if it was still running) and removes it from the runner.
func (r *Runner) Stop(a *gui.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()

	newList := make([]*anim, 0, len(r.animations))
	stopped := false
	for _, item := range r.animations {
		if item.a != a {
			newList = append(newList, item)
		} else {
			item.setStopped()
			stopped = true
		}
	}
	r.animations = newList
	if stopped {
		return
	}

	newList = make([]*anim, 0, len(r.pendingAnimations))
	for _, item := range r.pendingAnimations {
		if item.a != a {
			newList = append(newList, item)
		} else {
			item.setStopped()
		}
	}
	r.pendingAnimations = newList
}

func (r *Runner) runAnimations() {
	draw := time.NewTicker(time.Second / 60)

	go func() {
		for done := false; !done; {
			<-draw.C
			r.animationMutex.Lock()
			oldList := r.animations
			r.animationMutex.Unlock()
			newList := make([]*anim, 0, len(oldList))
			for _, a := range oldList {
				if !a.isStopped() && r.tickAnimation(a) {
					newList = append(newList, a)
				}
			}
			r.animationMutex.Lock()
			r.animations = append(newList, r.pendingAnimations...)
			r.pendingAnimations = nil
			done = len(r.animations) == 0
			r.animationMutex.Unlock()
		}
		r.animationMutex.Lock()
		r.runnerStarted = false
		r.animationMutex.Unlock()
		draw.Stop()
	}()
}

// tickAnimation will process a frame of animation and return true if this should continue animating
func (r *Runner) tickAnimation(a *anim) bool {
	if time.Now().After(a.end) {
		if a.reverse {
			a.a.Tick(0.0)
			if a.repeatsLeft == 0 {
				return false
			}
			a.reverse = false
		} else {
			a.a.Tick(1.0)
			if a.a.AutoReverse {
				a.reverse = true
			}
		}
		if !a.reverse {
			if a.repeatsLeft == 0 {
				return false
			}
			if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
		}

		a.start = time.Now()
		a.end = a.start.Add(a.a.Duration)
		return true
	}

	delta := time.Since(a.start).Nanoseconds() / 1000000 // TODO change this to Milliseconds() when we drop Go 1.12

	val := float32(delta) / float32(a.total)
	curve := a.a.Curve
	if curve == nil {
		curve = gui.AnimationEaseInOut
	}
	if a.reverse {
		a.a.Tick(curve(1 - val))
	} else {
		a.a.Tick(curve(val))
	}

	return true
}
