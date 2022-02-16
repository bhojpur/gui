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
	"sync"

	"github.com/bhojpur/gui/pkg/engine/data/binding"
)

// basicBinder stores a DataItem and a function to be called when it changes.
// It provides a convenient way to replace data and callback independently.
type basicBinder struct {
	callbackLock         sync.RWMutex
	callback             func(binding.DataItem) // access guarded by callbackLock
	dataListenerPairLock sync.RWMutex
	dataListenerPair     annotatedListener // access guarded by dataListenerPairLock
}

// Bind replaces the data item whose changes are tracked by the callback function.
func (binder *basicBinder) Bind(data binding.DataItem) {
	listener := binding.NewDataListener(func() { // NB: listener captures `data` but always calls the up-to-date callback
		binder.callbackLock.RLock()
		f := binder.callback
		binder.callbackLock.RUnlock()
		if f != nil {
			f(data)
		}
	})
	data.AddListener(listener)
	listenerInfo := annotatedListener{
		data:     data,
		listener: listener,
	}

	binder.dataListenerPairLock.Lock()
	binder.unbindLocked()
	binder.dataListenerPair = listenerInfo
	binder.dataListenerPairLock.Unlock()
}

// CallWithData passes the currently bound data item as an argument to the
// provided function.
func (binder *basicBinder) CallWithData(f func(data binding.DataItem)) {
	binder.dataListenerPairLock.RLock()
	data := binder.dataListenerPair.data
	binder.dataListenerPairLock.RUnlock()
	f(data)
}

// SetCallback replaces the function to be called when the data changes.
func (binder *basicBinder) SetCallback(f func(data binding.DataItem)) {
	binder.callbackLock.Lock()
	binder.callback = f
	binder.callbackLock.Unlock()
}

// Unbind requests the callback to be no longer called when the previously bound
// data item changes.
func (binder *basicBinder) Unbind() {
	binder.dataListenerPairLock.Lock()
	binder.unbindLocked()
	binder.dataListenerPairLock.Unlock()
}

// unbindLocked expects the caller to hold dataListenerPairLock.
func (binder *basicBinder) unbindLocked() {
	previousListener := binder.dataListenerPair
	binder.dataListenerPair = annotatedListener{nil, nil}

	if previousListener.listener == nil || previousListener.data == nil {
		return
	}
	previousListener.data.RemoveListener(previousListener.listener)
}

type annotatedListener struct {
	data     binding.DataItem
	listener binding.DataListener
}
