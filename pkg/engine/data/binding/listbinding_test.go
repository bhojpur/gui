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
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleList struct {
	listBase
}

func TestListBase_AddListener(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, len(data.listeners))

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, len(data.listeners))

	data.trigger()
	waitForItems()
	assert.True(t, called)
}

func TestListBase_GetItem(t *testing.T) {
	data := &simpleList{}
	f := 0.5
	data.appendItem(BindFloat(&f))
	assert.Equal(t, 1, len(data.items))

	item, err := data.GetItem(0)
	assert.Nil(t, err)
	val, err := item.(Float).Get()
	assert.Nil(t, err)
	assert.Equal(t, f, val)

	_, err = data.GetItem(5)
	assert.NotNil(t, err)
}

func TestListBase_Length(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, data.Length())

	data.appendItem(NewFloat())
	assert.Equal(t, 1, data.Length())
}

func TestListBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleList{}
	data.listeners = []DataListener{fn}

	assert.Equal(t, 1, len(data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, len(data.listeners))

	data.trigger()
	waitForItems()
	assert.False(t, called)
}

func TestNewDataListListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})

	fn.DataChanged()
	assert.True(t, called)
}
