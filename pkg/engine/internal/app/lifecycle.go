package app

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
	"sync/atomic"

	gui "github.com/bhojpur/gui/pkg/engine"
)

var _ gui.Lifecycle = (*Lifecycle)(nil)

// Lifecycle represents the various phases that an app can transition through.
//
// Since: 2.1
type Lifecycle struct {
	onForeground atomic.Value // func()
	onBackground atomic.Value // func()
	onStarted    atomic.Value // func()
	onStopped    atomic.Value // func()
}

// SetOnEnteredForeground hooks into the the app becoming foreground.
func (l *Lifecycle) SetOnEnteredForeground(f func()) {
	l.onForeground.Store(f)
}

// SetOnExitedForeground hooks into the app having moved to the background.
// Depending on the platform it may still be  visible but will not receive keyboard events.
// On some systems hover or desktop mouse move events may still occur.
func (l *Lifecycle) SetOnExitedForeground(f func()) {
	l.onBackground.Store(f)
}

// SetOnStarted hooks into an event that says the app is now running.
func (l *Lifecycle) SetOnStarted(f func()) {
	l.onStarted.Store(f)
}

// SetOnStopped hooks into an event that says the app is no longer running.
func (l *Lifecycle) SetOnStopped(f func()) {
	l.onStopped.Store(f)
}

// TriggerEnteredForeground will call the focus gained hook, if one is registered.
func (l *Lifecycle) TriggerEnteredForeground() {
	f := l.onForeground.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}

// TriggerExitedForeground will call the focus lost hook, if one is registered.
func (l *Lifecycle) TriggerExitedForeground() {
	f := l.onBackground.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}

// TriggerStarted will call the started hook, if one is registered.
func (l *Lifecycle) TriggerStarted() {
	f := l.onStarted.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}

// TriggerStopped will call the stopped hook, if one is registered.
func (l *Lifecycle) TriggerStopped() {
	f := l.onStopped.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}
