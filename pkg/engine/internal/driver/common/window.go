package common

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

	"github.com/bhojpur/gui/pkg/engine/internal/async"
)

// Window defines common functionality for windows.
type Window struct {
	eventQueue *async.UnboundedFuncChan
}

// DestroyEventQueue destroys the event queue.
func (w *Window) DestroyEventQueue() {
	w.eventQueue.Close()
}

// InitEventQueue initializes the event queue.
func (w *Window) InitEventQueue() {
	// This channel should be closed when the window is closed.
	w.eventQueue = async.NewUnboundedFuncChan()
}

// QueueEvent uses this method to queue up a callback that handles an event. This ensures
// user interaction events for a given window are processed in order.
func (w *Window) QueueEvent(fn func()) {
	w.eventQueue.In() <- fn
}

// RunEventQueue runs the event queue. This should called inside a go routine.
// This function blocks.
func (w *Window) RunEventQueue() {
	for fn := range w.eventQueue.Out() {
		fn()
	}
}

// WaitForEvents wait for all the events.
func (w *Window) WaitForEvents() {
	done := donePool.Get().(chan struct{})
	defer donePool.Put(done)

	w.eventQueue.In() <- func() { done <- struct{}{} }
	<-done
}

var donePool = sync.Pool{
	New: func() interface{} {
		return make(chan struct{})
	},
}
