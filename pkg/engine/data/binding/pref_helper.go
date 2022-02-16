package binding

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

type preferenceItem interface {
	checkForChange()
}

type preferenceBindings struct {
	lock  sync.RWMutex
	items map[string]preferenceItem
}

func (b *preferenceBindings) getItem(key string) preferenceItem {
	b.lock.RLock()
	item := b.items[key]
	b.lock.RUnlock()
	return item
}

func (b *preferenceBindings) list() []preferenceItem {
	b.lock.RLock()
	allItems := b.items
	b.lock.RUnlock()
	ret := make([]preferenceItem, 0, len(allItems))
	for _, i := range allItems {
		ret = append(ret, i)
	}
	return ret
}

func (b *preferenceBindings) setItem(key string, item preferenceItem) {
	b.lock.Lock()
	b.items[key] = item
	b.lock.Unlock()
}

type preferencesMap struct {
	lock  sync.RWMutex
	prefs map[gui.Preferences]*preferenceBindings
}

func newPreferencesMap() *preferencesMap {
	return &preferencesMap{
		prefs: make(map[gui.Preferences]*preferenceBindings),
	}
}

func (m *preferencesMap) ensurePreferencesAttached(p gui.Preferences) *preferenceBindings {
	m.lock.RLock()
	binds := m.prefs[p]
	m.lock.RUnlock()

	if binds != nil {
		return binds
	}

	m.lock.Lock()
	m.prefs[p] = &preferenceBindings{
		items: make(map[string]preferenceItem),
	}
	binds = m.prefs[p]
	m.lock.Unlock()

	p.AddChangeListener(func() {
		m.preferencesChanged(p)
	})
	return binds
}

func (m *preferencesMap) getBindings(p gui.Preferences) *preferenceBindings {
	m.lock.RLock()
	binds := m.prefs[p]
	m.lock.RUnlock()
	return binds
}

func (m *preferencesMap) preferencesChanged(p gui.Preferences) {
	binds := m.getBindings(p)
	if binds == nil {
		return
	}
	for _, item := range binds.list() {
		item.checkForChange()
	}
}
